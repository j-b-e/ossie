package shell

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/j-b-e/ossie/internal/config"
	"github.com/j-b-e/ossie/internal/model"
	"golang.org/x/sys/unix"
)

type ShellHandler interface {
	Spawn(model.Cloud)
	Update(model.Cloud)
	Prev() *string // returns previous session or nil if not found
	fmt.Stringer
}

func SpawnEnv(cloud model.Cloud) {
	shell := DetectShell()
	if shell == nil {
		fmt.Println("Shell not supported.")
		return
	}
	shell.Spawn(cloud)
}

func UpdateEnv(cloud model.Cloud) {
	shell := DetectShell()
	if shell == nil {
		fmt.Println("Shell not supported.")
		return
	}
	shell.Update(cloud)
}

func envToExport(cloud model.Cloud) string {
	var export string
	for k, v := range cloud.Env {
		export += "export " + k + "=\"" + v + "\"\n"
	}
	return export
}

func generatePrompt(cloud model.Cloud) string {
	replacer := strings.NewReplacer(
		"%n", "$"+bashCurrentSessionKey,
		"%r", "$OS_REGION_NAME",
		"%d", "$OS_DOMAIN_NAME",
		"%p", "$OS_PROJECT_NAME",
		"%u", "$OS_USERNAME",
	)
	prompt := replacer.Replace(config.Global.Prompt)
	return prompt
}

// Creates an in-memory tmpfile only accessible via file-descriptor
type Tmpfile struct {
	path string
	fd   int
}

func NewTempfile() (Tmpfile, error) {
	fd, err := unix.MemfdCreate("ossie_tmp", 0)
	if err != nil {
		return Tmpfile{}, err
	}
	fp := fmt.Sprintf("/proc/%d/fd/%d", os.Getpid(), fd)
	return Tmpfile{path: fp, fd: fd}, nil
}

func (t Tmpfile) Path() (path string) {
	return t.path
}

func (t Tmpfile) Write(content []byte) error {
	err := unix.Ftruncate(t.fd, int64(len(content)))
	if err != nil {
		return err
	}

	data, err := unix.Mmap(t.fd, 0, len(content), unix.PROT_READ|unix.PROT_WRITE, unix.MAP_SHARED)
	if err != nil {
		return err
	}

	copy(data, content)

	err = unix.Munmap(data)
	if err != nil {
		return err
	}
	return nil
}

// returns parent pid and parent cmd of a pid
func parentPid(pid int) (int, string) {
	path := path.Join("/proc", strconv.Itoa(pid), "stat")
	content, err := os.ReadFile(path)
	if err != nil {
		return 1, "error"
	}
	stat := strings.Split(string(content), " ")
	ppid, err := strconv.Atoi(stat[3])
	if err != nil {
		return 1, "error"
	}
	return ppid, stat[1]
}

// walkPidTree finds first ancestor which is a known shell
func walkPidTree(pid int) ShellHandler {
	ppid, cmd := parentPid(pid)

	switch {
	case strings.Contains(cmd, "bash"):
		return &Bash{}
	default:
		if ppid != 1 {
			return walkPidTree(ppid)
		}
	}
	return nil
}

func DetectShell() ShellHandler {
	return walkPidTree(os.Getpid())
}

func replaceFileContent(filePath, content string) error {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer func() { err = errors.Join(err, file.Close()) }()

	_, err = file.WriteString(content)
	if err != nil {
		return err
	}

	return nil
}

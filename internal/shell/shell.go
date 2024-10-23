package shell

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/j-b-e/ossie/internal/config"
	"github.com/j-b-e/ossie/internal/model"
	"golang.org/x/sys/unix"
)

type Shell string

const (
	Bash Shell = "bash"
)

func SpawnEnv(cloud model.Cloud) {
	switch detectShell() {
	case Bash:
		SpawnBash(cloud)
	default:
		fmt.Println("Shell not supported.")
	}
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
		"%n", cloud.Name,
		"%r", "$OS_REGION_NAME",
		"%d", "$OS_DOMAIN_NAME",
		"%p", "$OS_PROJECT_NAME",
		"%u", "$OS_USERNAME",
	)
	prompt := replacer.Replace(config.Global.Prompt)
	return prompt
}

// Creates an in-memory tmpfile only accessible via file-descriptor
type tmpfile []byte

// Returns path to tmpfile
func (t tmpfile) Path() (path string, fd int) {
	hash := sha256.Sum256(t)
	hashStr := hex.EncodeToString(hash[:])
	fd, err := unix.MemfdCreate(hashStr, 0)
	if err != nil {
		return
	}

	err = unix.Ftruncate(fd, int64(len(t)))
	if err != nil {
		return
	}

	data, err := unix.Mmap(fd, 0, len(t), unix.PROT_READ|unix.PROT_WRITE, unix.MAP_SHARED)
	if err != nil {
		return
	}

	copy(data, t)

	err = unix.Munmap(data)
	if err != nil {
		return
	}

	fp := fmt.Sprintf("/proc/self/fd/%d", fd)
	return fp, fd
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
func walkPidTree(pid int) Shell {
	ppid, cmd := parentPid(pid)

	switch {
	case strings.Contains(cmd, "bash"):
		return Bash
	default:
		if ppid != 1 {
			return walkPidTree(ppid)
		}
	}
	return Shell("unknown")
}

func detectShell() Shell {
	return walkPidTree(os.Getpid())
}

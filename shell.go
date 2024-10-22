package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"

	"golang.org/x/sys/unix"
)

type Shell string

const (
	Bash Shell = "bash"
)

func SpawnEnv(cloud Cloud) {
	switch detectShell() {
	case Bash:
		SpawnBash(cloud)
	default:
		fmt.Println("Shell not supported.")
	}
}

func envToExport(cloud Cloud) string {
	var export string
	for k, v := range cloud.Env {
		export += "export " + k + "=" + v + "\n"
	}
	return export
}

func generatePrompt(cloud Cloud) string {
	prompt := strings.ReplaceAll(gConf.Prompt, "%n", cloud.Name)
	prompt = strings.ReplaceAll(prompt, "%r", "$OS_REGION_NAME")
	prompt = strings.ReplaceAll(prompt, "%d", "$OS_DOMAIN_NAME")
	prompt = strings.ReplaceAll(prompt, "%p", "$OS_PROJECT_NAME")
	prompt = strings.ReplaceAll(prompt, "%u", "$OS_USERNAME")
	return prompt
}

func SpawnBash(cloud Cloud) {
	export := envToExport(cloud)

	ossierc, fd := tmpfile([]byte(strings.Join([]string{
		`
if [[ -f "/etc/bash.bashrc" ]] ; then
source "/etc/bash.bashrc"
fi

if [[ -f "$HOME/.bashrc" ]] ; then
	source "$HOME/.bashrc"
fi
` + export + `
function _ossie_exec_() {
	` + func(p bool) string {
			if p {
				return "unset ${!OS_*}\n" + export
			}
			return ""
		}(gConf.ProtectEnv) + `
	export ` + nestedEnvKey + `=` + nestedEnvVal + `
}
` +
			func(a bool) string {
				if a {
					return "alias os=openstack; alias o=openstack;"
				} else {
					return ""
				}
			}(gConf.Aliases) + `

function osenv() {
	echo -n '` + strings.ReplaceAll(export, "OS_PASSWORD="+cloud.Env["OS_PASSWORD"], "OS_PASSWORD=****") + `'
}

trap '_ossie_exec_' DEBUG
OLDPS="$PS1"
PS1="[` + generatePrompt(cloud) + `]$OLDPS"
`,
	}, "\n"))).Path()
	defer unix.Close(fd)
	newEnv := []string{
		fmt.Sprintf("%s=%s", nestedEnvKey, nestedEnvVal),
	}

	newEnv = append(newEnv, os.Environ()...)
	newCmd := []string{
		"bash",
		"--noprofile",
		"--rcfile",
		ossierc,
		"-i",
	}

	cmd := exec.Command(newCmd[0], newCmd[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		fmt.Println(err)
		return
	}
	err = cmd.Wait()
	if err != nil {
		fmt.Println(err)
		return
	}
}

// Creates an in-memory tmpfile only accessible via file-descriptor
type tmpfile []byte

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

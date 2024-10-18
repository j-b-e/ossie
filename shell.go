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

var (
	Bash Shell = "bash"
)

func SpawnEnv(cloud string) {
	s := detectShell()
	switch s {
	case Bash:
		SpawnBash(cloud)
	default:
		fmt.Printf("Shell \"%s\" not supported.\n", s)
	}
}

func SpawnBash(cloud string) {
	ossierc, fd := tmpfile([]byte(strings.Join([]string{
		`
if [[ -f "/etc/bash.bashrc" ]] ; then
source "/etc/bash.bashrc"
fi

if [[ -f "$HOME/.bashrc" ]] ; then
	source "$HOME/.bashrc"
fi
function _ossie_exec_() {
    export OS_CLOUD="` + cloud + `"
	export __OSSIE_SPAWNED=righty
}

trap '_ossie_exec_' DEBUG
OLDPS="$PS1"
PS1="[` + cloud + `]$OLDPS"
`,
		// fmt.Sprintf("PS1=\"[%s]\\h$ \"", cloud),
		// fmt.Sprintf("export OS_CLOUD2=%s", cloud),
		// "",
	}, "\n"))).Path()
	defer unix.Close(fd)
	newEnv := []string{
		fmt.Sprintf("OS_CLOUD=%s", cloud),
		fmt.Sprintf("__OSSIE_SPAWNED=righto"),
	}
	ps1 := os.Getenv("PS1")
	os.Setenv("PS1", fmt.Sprintf("[%s]%s", cloud, ps1))
	newEnv = append(newEnv, os.Environ()...)
	newCmd := []string{
		"bash",
		"--noprofile",
		"--rcfile",
		ossierc,
		"-i",
	}
	// if err := syscall.Exec(newCmd[0], newCmd, newEnv); err != nil {
	// 	fmt.Println(err)
	// }
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

type tmpfile []byte

func (t tmpfile) Path() (path string, fd int) {
	hash := sha256.Sum256(t)
	hashStr := hex.EncodeToString(hash[:])
	fd, err := unix.MemfdCreate(hashStr, 0)
	if err != nil {
		return //0, fmt.Errorf("MemfdCreate: %v", err)
	}

	err = unix.Ftruncate(fd, int64(len(t)))
	if err != nil {
		return //0, fmt.Errorf("Ftruncate: %v", err)
	}

	data, err := unix.Mmap(fd, 0, len(t), unix.PROT_READ|unix.PROT_WRITE, unix.MAP_SHARED)
	if err != nil {
		return //0, fmt.Errorf("Mmap: %v", err)
	}

	copy(data, t)

	err = unix.Munmap(data)
	if err != nil {
		return //0, fmt.Errorf("Munmap: %v", err)
	}

	fp := fmt.Sprintf("/proc/self/fd/%d", fd)
	return fp, fd
}

type envars []string

func (e envars) Add(env string) {
	e = append(e, env)
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

// walk pid tree to find first ancestor which is a known shell
func detectShell() Shell {
	return walkPidTree(os.Getpid())
}

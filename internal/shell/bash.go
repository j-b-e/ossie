package shell

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/j-b-e/ossie/internal/config"
	"github.com/j-b-e/ossie/internal/model"
	"golang.org/x/sys/unix"
)

func bashRC(cloud model.Cloud) string {
	export := envToExport(cloud)
	var env strings.Builder

	env.WriteString(`
if [[ -f "/etc/bash.bashrc" ]] ; then
  source "/etc/bash.bashrc"
fi
if [[ -f "$HOME/.bashrc" ]] ; then
  source "$HOME/.bashrc"
fi
`)
	env.WriteString(export)
	env.WriteString("function _ossie_exec_ () {\n")
	if config.Global.ProtectEnv {
		env.WriteString("unset ${!OS_*}\n")
		env.WriteString(export)
	}
	env.WriteString("export " + config.NestedEnvKey + "=" + config.NestedEnvVal + "\n}\n")
	if config.Global.Aliases {
		env.WriteString("alias os=openstack; alias o=openstack;\n")
	}
	env.WriteString(`function osenv () {
	echo -n '` + strings.ReplaceAll(export, "OS_PASSWORD="+cloud.Env["OS_PASSWORD"], "OS_PASSWORD=****") + `'
}
trap '_ossie_exec_' DEBUG
_ossie_OLDPS="$PS1"
PS1="[` + generatePrompt(cloud) + `]$_ossie_OLDPS"
`)
	return env.String()
}

func SpawnBash(cloud model.Cloud) {
	ossierc, fd := tmpfile([]byte(bashRC(cloud))).Path()
	defer unix.Close(fd)
	newEnv := []string{
		fmt.Sprintf("%s=%s", config.NestedEnvKey, config.NestedEnvVal),
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

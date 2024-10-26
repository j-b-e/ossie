package shell

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/j-b-e/ossie/internal/config"
	"github.com/j-b-e/ossie/internal/model"
	"golang.org/x/sys/unix"
)

func bashRC(cloud model.Cloud) string {
	const rcTempl = `
if [[ -f "/etc/bash.bashrc" ]] ; then
  source "/etc/bash.bashrc"
fi
if [[ -f "$HOME/.bashrc" ]] ; then
  source "$HOME/.bashrc"
fi
{{ .export }}
function _ossie_exec_ () {
{{ if .protectenv -}}
unset ${!OS_*}
{{ .export }}
{{ end -}}
export {{ .nested_marker }}
}
{{ if .aliases -}}
alias os=openstack
alias o=openstack
{{ end -}}
function osenv () {
	echo -n '{{ .export_safe }}'
}
trap '_ossie_exec_' DEBUG
_ossie_OLDPS="$PS1"
PS1="[{{ .prompt }}]$_ossie_OLDPS"
`
	var out strings.Builder
	t := template.Must(template.New("rc").Parse(rcTempl))
	env := envToExport(cloud)
	data := map[string]any{
		"export":        env,
		"export_safe":   strings.ReplaceAll(env, cloud.Env["OS_PASSWORD"], "****"),
		"nested_marker": config.NestedEnvKey + "=" + config.NestedEnvVal,
		"prompt":        generatePrompt(cloud),
		"protectenv":    config.Global.ProtectEnv,
		"aliases":       config.Global.Aliases,
	}
	err := t.Execute(&out, data)
	if err != nil {
		panic(err)
	}
	return out.String()
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

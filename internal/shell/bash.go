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

type Bash struct{}

const (
	bashOsEnvFileKey   = "__OSSIE_OS_ENV_FILE"
	bashPromptFileKey  = "__OSSIE_PROMPT_FILE"
	bashSessionFileKey = "__OSSIE_SESSION_FILE"
)

func bashRC(cloud model.Cloud, osrc string, promptfile string) string {
	const rcTempl = `
if [[ -f "/etc/bash.bashrc" ]] ; then
  source "/etc/bash.bashrc"
fi
if [[ -f "$HOME/.bashrc" ]] ; then
  source "$HOME/.bashrc"
fi
export ` + bashOsEnvFileKey + `="{{ .osrc }}"
export ` + bashPromptFileKey + `="{{ .promptfile }}"

set -o allexport
source $` + bashOsEnvFileKey + `
set +o allexport
function _ossie_exec_ () {
  export {{ .nested_marker }}
  export ` + bashOsEnvFileKey + `="{{ .osrc }}"
{{ if .protectenv }}
  unset ${!OS_*}
  set -o allexport
  source $` + bashOsEnvFileKey + `
  set +o allexport
{{ end -}}
}

{{ if .aliases -}}
alias os=openstack
alias o=openstack
{{ end -}}
function osenv () {
  while IFS= read -r line; do
    if [[ "$line" == *"OS_PASSWORD"* ]]; then
      echo 'OS_PASSWORD="****"'
    else
      echo ${line/export/}
    fi
  done < "$` + bashOsEnvFileKey + `"
}
trap '_ossie_exec_' DEBUG
_ossie_OLDPS="$PS1"

PS1="[$(<{{ .promptfile }})]$_ossie_OLDPS"
`
	var out strings.Builder
	t := template.Must(template.New("rc").Parse(rcTempl))
	data := map[string]any{
		"nested_marker": config.NestedEnvKey + "=" + config.NestedEnvVal,
		"protectenv":    config.Global.ProtectEnv,
		"aliases":       config.Global.Aliases,
		"osrc":          osrc,
		"promptfile":    promptfile,
	}
	err := t.Execute(&out, data)
	if err != nil {
		panic(err)
	}
	return out.String()
}

func (b *Bash) Spawn(cloud model.Cloud) {
	osrc, err := NewTempfile()
	if err != nil {
		panic(err)
	}
	defer unix.Close(osrc.fd)
	osrc.Write([]byte(envToExport(cloud)))

	prompt, err := NewTempfile()
	if err != nil {
		panic(err)
	}
	defer unix.Close(prompt.fd)
	prompt.Write([]byte(generatePrompt(cloud)))

	ossierc, err := NewTempfile()
	if err != nil {
		panic(err)
	}
	defer unix.Close(ossierc.fd)
	err = ossierc.Write([]byte(bashRC(cloud, osrc.Path(), prompt.Path())))

	newEnv := []string{
		fmt.Sprintf("%s=%s", config.NestedEnvKey, config.NestedEnvVal),
	}

	newEnv = append(newEnv, os.Environ()...)
	newCmd := []string{
		"bash",
		"--noprofile",
		"--rcfile",
		ossierc.Path(),
		"-i",
	}

	cmd := exec.Command(newCmd[0], newCmd[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	err = cmd.Start()
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

func (b *Bash) Update(cloud model.Cloud) {
	osEnvFile := os.Getenv(bashOsEnvFileKey)
	err := replaceFileContent(osEnvFile, envToExport(cloud))
	if err != nil {
		panic(err)
	}
	promptFile := os.Getenv(bashPromptFileKey)
	err = replaceFileContent(promptFile, generatePrompt(cloud))
	if err != nil {
		panic(err)
	}
}

func (b Bash) String() string {
	return "Bash"
}

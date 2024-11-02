package shell

import (
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"
	"text/template"

	"github.com/j-b-e/ossie/internal/config"
	"github.com/j-b-e/ossie/internal/model"
	"golang.org/x/sys/unix"
)

type bash struct{}

const (
	bashOsEnvFileKey      = "__OSSIE_OS_ENV_FILE"
	bashPromptFileKey     = "__OSSIE_PROMPT_FILE"
	bashSessionFileKey    = "__OSSIE_SESSION_FILE"
	bashCurrentSessionKey = "__OSSIE_CURRENT_SESSION_"
	bashPrevSessionKey    = "__OSSIE_PREV_SESSION_"
)

func bashRC(osrc string, promptfile string, sessionfile string) string {
	const rcTempl = `
if [[ -f "/etc/bash.bashrc" ]] ; then
  source "/etc/bash.bashrc"
fi
if [[ -f "$HOME/.bashrc" ]] ; then
  source "$HOME/.bashrc"
fi
export ` + bashOsEnvFileKey + `="{{ .osrc }}"
export ` + bashPromptFileKey + `="{{ .promptfile }}"
export ` + bashSessionFileKey + `="{{ .sessionfile }}"


set -o allexport
source $` + bashOsEnvFileKey + `
set +o allexport
function _ossie_exec_ () {
  export {{ .nested_marker }}
  export ` + bashOsEnvFileKey + `="{{ .osrc }}"
  export ` + bashSessionFileKey + `="{{ .sessionfile }}"
  set -o allexport
  source $` + bashSessionFileKey + `
  set +o allexport
{{- if .protectenv }}
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
		"sessionfile":   sessionfile,
	}
	err := t.Execute(&out, data)
	if err != nil {
		panic(err)
	}
	return out.String()
}

// Spawn executes new Shell environment
func (b *bash) Spawn(cloud model.Cloud) {
	osrc, err := NewTempfile()
	if err != nil {
		panic(err)
	}
	defer unix.Close(osrc.fd)
	_, err = osrc.Write([]byte(envToExport(cloud)))
	if err != nil {
		panic(err)
	}

	prompt, err := NewTempfile()
	if err != nil {
		panic(err)
	}
	defer unix.Close(prompt.fd)
	_, err = prompt.Write([]byte(generatePrompt()))
	if err != nil {
		panic(err)
	}

	sessionfile, err := NewTempfile()
	if err != nil {
		panic(err)
	}
	defer unix.Close(sessionfile.fd)
	_, err = sessionfile.Write([]byte(
		"export " + bashCurrentSessionKey + "=\"" + cloud.Name + "\"\nexport " + bashPrevSessionKey + "=\"-\"",
	))
	if err != nil {
		panic(err)
	}

	ossierc, err := NewTempfile()
	if err != nil {
		panic(err)
	}
	defer unix.Close(ossierc.fd)
	_, err = ossierc.Write([]byte(bashRC(osrc.Path(), prompt.Path(), sessionfile.Path())))
	if err != nil {
		panic(err)
	}

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

// Update sets new env without spawning a New shell
func (b *bash) Update(cloud model.Cloud) {
	currentSession := os.Getenv(bashCurrentSessionKey)

	osEnvFile := os.Getenv(bashOsEnvFileKey)
	err := replaceFileContent(osEnvFile, envToExport(cloud))
	if err != nil {
		panic(err)
	}
	promptFile := os.Getenv(bashPromptFileKey)
	err = replaceFileContent(promptFile, generatePrompt())
	if err != nil {
		panic(err)
	}
	sessionFile := os.Getenv(bashSessionFileKey)
	err = replaceFileContent(sessionFile,
		"export "+bashCurrentSessionKey+"=\""+cloud.Name+"\"\nexport "+bashPrevSessionKey+"=\""+currentSession+"\"",
	)
	if err != nil {
		panic(err)
	}
}

func (b bash) String() string {
	return "Bash"
}

// Prev returns previous session name
func (b *bash) Prev() *string {
	osPrevSession := os.Getenv(bashPrevSessionKey)
	if slices.Contains([]string{"-", ""}, osPrevSession) {
		return nil
	}
	return &osPrevSession
}

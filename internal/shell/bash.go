package shell

import (
	_ "embed"
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

//go:embed bashrc.j2
var rcTempl string

type Bash struct{}

const (
	bashOsEnvFileKey      = "__OSSIE_OS_ENV_FILE"
	bashPromptFileKey     = "__OSSIE_PROMPT_FILE"
	bashSessionFileKey    = "__OSSIE_SESSION_FILE"
	bashCurrentSessionKey = "__OSSIE_CURRENT_SESSION_"
	bashPrevSessionKey    = "__OSSIE_PREV_SESSION_"
)

func bashRC(cloud model.Cloud, osrc string, promptfile string, sessionfile string) string {

	var out strings.Builder
	t := template.Must(template.New("rc").Parse(rcTempl))
	data := map[string]any{
		"nested_marker":  config.NestedEnvKey + "=" + config.NestedEnvVal,
		"protectenv":     config.Global.ProtectEnv,
		"aliases":        config.Global.Aliases,
		"osrc":           osrc,
		"promptfile":     promptfile,
		"sessionfile":    sessionfile,
		"OsEnvFileKey":   bashOsEnvFileKey,
		"PromptFileKey":  bashPromptFileKey,
		"SessionFileKey": bashSessionFileKey,
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
	defer func() {
		_ = unix.Close(osrc.fd)
	}()
	err = osrc.Write([]byte(envToExport(cloud)))
	if err != nil {
		panic(err)
	}

	prompt, err := NewTempfile()
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = unix.Close(prompt.fd)
	}()
	err = prompt.Write([]byte(generatePrompt(cloud)))
	if err != nil {
		panic(err)
	}

	sessionfile, err := NewTempfile()
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = unix.Close(sessionfile.fd)
	}()
	err = sessionfile.Write([]byte(
		"export " + bashCurrentSessionKey + "=\"" + cloud.Name + "\"\nexport " + bashPrevSessionKey + "=\"-\"",
	))
	if err != nil {
		panic(err)
	}

	ossierc, err := NewTempfile()
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = unix.Close(ossierc.fd)
	}()
	err = ossierc.Write([]byte(bashRC(cloud, osrc.Path(), prompt.Path(), sessionfile.Path())))
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

func (b *Bash) Update(cloud model.Cloud) {
	currentSession := os.Getenv(bashCurrentSessionKey)

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
	sessionFile := os.Getenv(bashSessionFileKey)
	err = replaceFileContent(sessionFile,
		"export "+bashCurrentSessionKey+"=\""+cloud.Name+"\"\nexport "+bashPrevSessionKey+"=\""+currentSession+"\"",
	)
	if err != nil {
		panic(err)
	}
}

func (b Bash) String() string {
	return "Bash"
}

func (b *Bash) Prev() *string {
	osPrevSession := os.Getenv(bashPrevSessionKey)
	if slices.Contains([]string{"-", ""}, osPrevSession) {
		return nil
	}
	return &osPrevSession
}

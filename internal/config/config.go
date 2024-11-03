package config

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/j-b-e/ossie/internal/load"
	"github.com/j-b-e/ossie/internal/model"
	"github.com/urfave/cli/v3"
)

var configDefaultPath = "~/.config/openstack/ossie.toml"

const (
	NestedEnvKey = "__OSSIE_SPAWNED"
	NestedEnvVal = "righto"
)

// Config is the ossie configuration
type Config struct {
	// Comments will be used for generating ossie.toml.example
	RCPath     string // Path to RC Files, clouds.yaml from standard Paths will always be loaded
	Prompt     string // Customize prompt: %n = Name, %r = Region, %d = Domain, %p = Project, %u = User
	ProtectEnv bool   // Ensures that "OS_" envars cant be overwritten in the spawned Shell
	Aliases    bool   // Sets up Shell aliases 'os' and 'o'
	Clouds     model.Clouds
}

// Global ossie configuration
var Global = Config{
	// The default values
	RCPath:     "~/.config/openstack",
	Prompt:     "%n:%r",
	ProtectEnv: true,
	Aliases:    false,
}

// Replace ~ with the home directory path
func expandHomedir(path string) string {
	if strings.HasPrefix(path, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		path = filepath.Join(homeDir, path[1:])
	}
	return path
}

// SetupConfig loads configfiles and sets up global ossie configuration
func SetupConfig(_ context.Context, c *cli.Command) error {
	configfile := expandHomedir(c.String("config"))
	if configfile == "" {
		configfile = expandHomedir(configDefaultPath)
	}
	_, err := toml.DecodeFile(configfile, &Global)
	if err != nil {
		return err
	}
	Global.RCPath = expandHomedir(Global.RCPath)

	Global.Clouds = load.Clouds(Global.RCPath)
	if len(Global.Clouds) == 0 {
		return fmt.Errorf("No clouds found")
	}
	return nil
}

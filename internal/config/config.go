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

type Config struct {
	RCPath     string // Path to openstack rc files
	Prompt     string // Prompt definiton
	ProtectEnv bool   // Protect OS_ env against accidental modfication
	Aliases    bool   // setup shell aliases o and os
	Clouds     model.Clouds
}

var Global Config

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

func SetupConfig(_ context.Context, c *cli.Command) error {
	configfile := expandHomedir(c.String("config"))
	if configfile == "" {
		configfile = expandHomedir(configDefaultPath)
	}
	Global = Config{
		// The default values
		RCPath:     "~/.config/openstack/",
		Prompt:     "%n:%r",
		ProtectEnv: true,
		Aliases:    false,
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

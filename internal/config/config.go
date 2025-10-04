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
	RCPath     string  `toml:"rcpath"`     // Path to openstack rc files
	Prompt     string  `toml:"prompt"`     // Prompt definiton
	ProtectEnv bool    `toml:"protectenv"` // Protect OS_ env against accidental modfication
	Aliases    bool    `toml:"aliases"`    // setup shell aliases o and os
	Shell      *string `toml:"shell"`      // which shell to use (disables autodetect)
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

func SetupConfig(ctx context.Context, c *cli.Command) (context.Context, error) {
	configfile := expandHomedir(c.String("config"))
	if configfile == "" {
		configfile = expandHomedir(configDefaultPath)
	}
	Global = Config{
		// Set default values
		RCPath:     "~/.config/openstack/",
		Prompt:     "%n:%r",
		ProtectEnv: true,
		Aliases:    false,
	}
	_, err := toml.DecodeFile(configfile, &Global)
	if err != nil {
		return ctx, err
	}
	Global.RCPath = expandHomedir(Global.RCPath)

	Global.Clouds = load.Clouds(Global.RCPath)
	if len(Global.Clouds) == 0 {
		return ctx, fmt.Errorf("no clouds found")
	}
	return ctx, nil
}

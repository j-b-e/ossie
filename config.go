package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/urfave/cli/v3"
)

var configDefaultPath = "~/.config/openstack/ossie.toml"

type Config struct {
	RCPath     string
	Prompt     string
	ProtectEnv bool
	clouds     []Cloud
}

var gConf Config

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
	gConf = Config{
		// Set default values
		RCPath:     "~/.config/openstack/",
		Prompt:     "%n:%r",
		ProtectEnv: true,
	}
	toml.DecodeFile(configfile, &gConf)
	gConf.RCPath = expandHomedir(gConf.RCPath)

	gConf.clouds = loadClouds(gConf.RCPath)
	if len(gConf.clouds) == 0 {
		return fmt.Errorf("No clouds found")
	}
	return nil
}

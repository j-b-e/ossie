package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

var configDefaultPath = "ossie.toml"

type Config struct {
	RCPath     string
	Prompt     string
	ProtectEnv bool
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

func SetupConfig(configfile string) (Config, error) {
	if configfile == "" {
		configfile = configDefaultPath
	}
	c := Config{
		// Set default values
		RCPath:     "~/.config/openstack/",
		Prompt:     "%n:%r",
		ProtectEnv: true,
	}
	toml.DecodeFile(configfile, &c)
	c.RCPath = expandHomedir(c.RCPath)

	return c, nil
}

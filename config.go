package main

import (
	"io"
	"os"

	"github.com/BurntSushi/toml"
)

var configDefaultPath = "ossie.toml"

type Config struct {
	RCPath string
}

func SetupConfig(configfile string) (Config, error) {
	if configfile == "" {
		configfile = configDefaultPath
	}
	c := Config{
		// Set default values
		RCPath: "~/.config/openstack/",
	}
	f, _ := os.Open(configfile)
	bytes, _ := io.ReadAll(f)
	_ = toml.Unmarshal(bytes, &c)
	return c, nil
}

package main

import (
	"fmt"
	"io"
	"os"
	"path"

	"gopkg.in/yaml.v3"
)

func loadCloudsYaml() []Clouds {
	var t map[string]any
	home := os.Getenv("HOME")
	f, err := os.Open(path.Join(home, ".config", "openstack", "clouds.yaml"))
	if err != nil {
		fmt.Println(err)
		return nil
	}
	bytes, _ := io.ReadAll(f)
	_ = yaml.Unmarshal(bytes, &t)
	clouds := []Clouds{}
	tree, ok := t["clouds"].(map[string]any)
	if !ok {
		fmt.Println("No clouds found.")
		return nil
	}
	for k := range tree {
		clouds = append(clouds, Clouds{name: k})

	}
	return clouds
}

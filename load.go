package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"gopkg.in/yaml.v3"
)

func loadClouds(rcPath string) []Cloud {
	cloudyml := loadCloudsYaml()
	cloudrc := loadRCClouds(rcPath)
	return append(cloudyml, cloudrc...)
}

func extractCloudYamlEnv(input map[string]any) map[string]string {
	result := make(map[string]string)
	for key, value := range input {
		switch v := value.(type) {
		case string:
			result["OS_"+strings.ToUpper(key)] = v
		case map[string]any:
			for subKey, subValue := range v {
				if strVal, ok := subValue.(string); ok {
					result["OS_"+strings.ToUpper(subKey)] = strVal
				}
			}
		}
	}

	return result
}

func loadCloudsYaml() []Cloud {
	var t map[string]any
	home := os.Getenv("HOME")
	f, err := os.Open(path.Join(home, ".config", "openstack", "clouds.yaml"))
	if err != nil {
		fmt.Println(err)
		return nil
	}
	bytes, _ := io.ReadAll(f)
	_ = yaml.Unmarshal(bytes, &t)
	clouds := []Cloud{}
	tree, ok := t["clouds"].(map[string]any)
	if !ok {
		fmt.Println("No clouds found.")
		return nil
	}
	for k, v := range tree {
		cloud := Cloud{Name: k}
		cloud.Env = extractCloudYamlEnv(v.(map[string]any))
		clouds = append(clouds, cloud)
	}
	return clouds
}

func loadRCClouds(rcPath string) []Cloud {
	files, _ := os.ReadDir(rcPath)
	var clouds []Cloud
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if cloud := loadRC(path.Join(rcPath, file.Name())); cloud.Name != "" {
			clouds = append(clouds, cloud)
		}
	}
	return clouds
}

func loadRC(filePath string) Cloud {
	file, err := os.Open(filePath)
	if err != nil {
		return Cloud{}
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	name := filePath[strings.LastIndex(filePath, "/")+1:]
	cloud := Cloud{Name: name, Env: make(map[string]string)}
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			continue
		}
		if !strings.HasPrefix(line, "export OS_") {
			continue
		}
		entry, _ := strings.CutPrefix(line, "export ")
		split := strings.Split(entry, "=")
		cloud.Env[split[0]] = split[1]
	}
	if err := scanner.Err(); err != nil {
		return Cloud{}
	}
	if len(cloud.Env) == 0 {
		return Cloud{}
	}
	return cloud
}

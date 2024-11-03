package load

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path"
	"reflect"
	"slices"
	"strconv"
	"strings"

	"github.com/j-b-e/ossie/internal/model"
	"gopkg.in/yaml.v3"
)

func Clouds(rcPath string) model.Clouds {
	cloudyml := loadCloudsYaml()
	cloudrc := loadRCClouds(rcPath)
	return append(cloudyml, cloudrc...)
}

func extractCloudYamlEnv(input map[string]any) (map[string]string, error) {
	result := make(map[string]string)
	for key, value := range input {
		if slices.Contains([]string{"log_file", "log_level", "operation_log", "cloud"}, key) {
			//TODO: "cloud" subkey should merge from clouds-public.yml
			// filter out keys
			continue
		}
		switch v := value.(type) {
		case string:
			result["OS_"+strings.ToUpper(key)] = v
		case map[string]any:
			for subKey, subValue := range v {
				switch sv := subValue.(type) {
				case string:
					result["OS_"+strings.ToUpper(subKey)] = sv
				case int:
					result["OS_"+strings.ToUpper(subKey)] = strconv.Itoa(sv)
				default:
					return nil, fmt.Errorf("unexpected type encountered: %v:%s", sv, reflect.TypeOf(sv))
				}
			}
		default:
			return nil, fmt.Errorf("unexpected type encountered: %v:%s", v, reflect.TypeOf(v))
		}
	}
	return result, nil
}

func loadCloudsYaml() model.Clouds {
	var t map[string]any
	home := os.Getenv("HOME")
	f, err := os.Open(path.Join(home, ".config", "openstack", "clouds.yaml"))
	if err != nil {
		fmt.Println(err)
		return nil
	}
	bytes, _ := io.ReadAll(f)
	_ = yaml.Unmarshal(bytes, &t)
	tree, ok := t["clouds"].(map[string]any)
	if !ok {
		fmt.Println("No clouds found.")
		return nil
	}
	clouds := model.Clouds{}
	for k, v := range tree {
		cloud := model.Cloud{Name: k, Source: "~/.config/openstack/clouds.yaml"}
		cloud.Env, err = extractCloudYamlEnv(v.(map[string]any))
		if err != nil {
			fmt.Printf("error: could not load cloud \"%s\": %s\n", k, err)
			continue
		}
		clouds = append(clouds, cloud)
	}
	return clouds
}

func loadRCClouds(rcPath string) model.Clouds {
	files, _ := os.ReadDir(rcPath)
	var clouds model.Clouds
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

func loadRC(filePath string) model.Cloud {
	file, err := os.Open(filePath)
	if err != nil {
		return model.Cloud{}
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	name := filePath[strings.LastIndex(filePath, "/")+1:]
	cloud := model.Cloud{Name: name, Env: make(map[string]string), Source: filePath}
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			continue
		}
		if !strings.HasPrefix(line, "export OS_") {
			continue
		}
		entry, _ := strings.CutPrefix(line, "export ")
		key, val, _ := strings.Cut(entry, "=")
		cloud.Env[key] = val
	}
	if err := scanner.Err(); err != nil {
		return model.Cloud{}
	}
	if len(cloud.Env) == 0 {
		return model.Cloud{}
	}
	return cloud
}

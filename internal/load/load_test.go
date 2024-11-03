package load

import (
	"reflect"
	"testing"

	"gopkg.in/yaml.v3"
)

const (
	sample1 = `
    auth:
      auth_url: http://192.168.122.10:5000/
      project_name: demo
      username: demo
      password: 0penstack
    region_name: RegionOne
`
	sample2 = `
    cloud: rackspace
    auth:
      project_id: 275610
      username: openstack
      password: xyzpdq!lazydog
    region_name: DFW,ORD,IAD
    interface: internal
`
	sample3 = `
    interface: public
    operation_log:
      logging: TRUE
      file: /tmp/openstackclient_demo.log
      level: info
    log_file: /tmp/openstackclient_admin.log
    log_level: debug
`
)

func sampleToInput(sample string) (result map[string]any) {
	yamlBytes := []byte(sample)
	err := yaml.Unmarshal(yamlBytes, &result)
	if err != nil {
		panic(err)
	}
	return

}

func Test_extractCloudYamlEnv(t *testing.T) {

	tests := []struct {
		name  string
		input map[string]any
		want  map[string]string
	}{
		{
			name:  "offical sample 1",
			input: sampleToInput(sample1),
			want: map[string]string{
				"OS_AUTH_URL":     "http://192.168.122.10:5000/",
				"OS_PASSWORD":     "0penstack",
				"OS_PROJECT_NAME": "demo",
				"OS_USERNAME":     "demo",
				"OS_REGION_NAME":  "RegionOne",
			},
		},
		{
			name:  "offical sample 2",
			input: sampleToInput(sample2),
			want: map[string]string{
				"OS_INTERFACE":   "internal",
				"OS_PASSWORD":    "xyzpdq!lazydog",
				"OS_PROJECT_ID":  "275610",
				"OS_REGION_NAME": "DFW,ORD,IAD",
				"OS_USERNAME":    "openstack",
			},
		},
		{
			name:  "sample 3",
			input: sampleToInput(sample3),
			want: map[string]string{
				"OS_INTERFACE": "public",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := extractCloudYamlEnv(tt.input)
			if err != nil {
				t.Errorf("extractCloudYamlEnv() unexpected error %v, want  %v", err, tt.want)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("extractCloudYamlEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}

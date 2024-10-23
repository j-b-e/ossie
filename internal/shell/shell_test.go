package shell

import (
	"testing"

	"github.com/j-b-e/ossie/internal/config"
	"github.com/j-b-e/ossie/internal/model"
)

func Test_generatePrompt(t *testing.T) {
	type args struct {
		cloud model.Cloud
	}
	var testCloud = model.Cloud{
		Name: "testcloud",
		Env: map[string]string{
			"OS_PROJECT_NAME": "Retinentals",
			"OS_DOMAIN_NAME":  "Socks",
			"OS_REGION_NAME":  "region1",
			"OS_USERNAME":     "belzebub",
		},
	}
	tests := []struct {
		name   string
		args   args
		want   string
		prompt string
	}{
		{
			name:   "test-basic",
			args:   args{cloud: testCloud},
			prompt: "%n:%r",
			want:   testCloud.Name + ":$OS_REGION_NAME",
		},
		{
			name:   "test-basic",
			args:   args{cloud: testCloud},
			prompt: "%n:%r:%u:%u:%p:%%d%d",
			want:   testCloud.Name + ":$OS_REGION_NAME:$OS_USERNAME:$OS_USERNAME:$OS_PROJECT_NAME:%$OS_DOMAIN_NAME$OS_DOMAIN_NAME",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config.Global.Prompt = tt.prompt
			if got := generatePrompt(tt.args.cloud); got != tt.want {
				t.Errorf("generatePrompt() = %v, want %v", got, tt.want)
			}
		})
	}
}

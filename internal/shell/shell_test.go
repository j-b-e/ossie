package shell

import (
	"testing"

	"github.com/j-b-e/ossie/internal/config"
)

func Test_generatePrompt(t *testing.T) {

	tests := []struct {
		name   string
		want   string
		prompt string
	}{
		{
			name: "test-basic",

			prompt: "%n:%r",
			want:   "$" + bashCurrentSessionKey + ":$OS_REGION_NAME",
		},
		{
			name: "test-basic",

			prompt: "%n:%r:%u:%u:%p:%%d%d",
			want:   "$" + bashCurrentSessionKey + ":$OS_REGION_NAME:$OS_USERNAME:$OS_USERNAME:$OS_PROJECT_NAME:%$OS_DOMAIN_NAME$OS_DOMAIN_NAME",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config.Global.Prompt = tt.prompt
			if got := generatePrompt(); got != tt.want {
				t.Errorf("generatePrompt() = %v, want %v", got, tt.want)
			}
		})
	}
}

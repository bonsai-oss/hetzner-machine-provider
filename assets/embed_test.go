package assets_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/bonsai-oss/hetzner-machine-provider/assets"
)

func TestCloudInitTemplate(t *testing.T) {
	for _, testCase := range []struct {
		name      string
		input     map[string]any
		checkFunc func(t *testing.T, output *bytes.Buffer)
	}{
		{
			name: "set ssh_authorized_keys",
			input: map[string]any{
				"ssh_authorized_keys": []string{
					"ssh-rsa AAAAB3NzaC1yc2E...",
					"ssh-rsa AAAAB3NzaC1yc2E...",
				},
			},
			checkFunc: func(t *testing.T, output *bytes.Buffer) {
				if !strings.Contains(output.String(), "ssh_authorized_keys:") {
					t.Fatalf("template output does not contain ssh_authorized_keys")
				}
			},
		},
		{
			name: "set ssh_authorized_keys with empty slice",
			input: map[string]any{
				"ssh_authorized_keys": []string{},
			},
			checkFunc: func(t *testing.T, output *bytes.Buffer) {
				if strings.Contains(output.String(), "ssh_authorized_keys:") {
					t.Fatalf("template output contains ssh_authorized_keys")
				}
			},
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			buf := &bytes.Buffer{}

			err := assets.CloudInitTemplate.Execute(buf, testCase.input)
			if err != nil {
				t.Fatalf("failed to execute template: %s", err)
			}

			testCase.checkFunc(t, buf)
		})
	}
}

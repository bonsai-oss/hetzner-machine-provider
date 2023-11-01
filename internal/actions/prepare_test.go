package actions

import (
	"testing"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

func TestDetermineArchitectureString(t *testing.T) {
	for _, testCase := range []struct {
		architecture hcloud.Architecture
		expected     string
	}{
		{hcloud.ArchitectureX86, "amd64"},
		{hcloud.ArchitectureARM, "arm64"},
	} {
		architecture := determineArchitectureString(testCase.architecture)
		if architecture != testCase.expected {
			t.Errorf("expected %s, got %s", testCase.expected, architecture)
		}
	}
}

func TestGetServerClassification(t *testing.T) {
	for _, testCase := range []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "server name with number",
			input:    "cx11",
			expected: "11",
		},
		{
			name:     "server name without number",
			input:    "cx",
			expected: "",
		},
		{
			name:     "server name with number and suffix",
			input:    "cx11-suffix",
			expected: "11",
		},
		{
			name:     "server name with number and prefix",
			input:    "prefix-cx11",
			expected: "11",
		},
		{
			name:     "server name with number and prefix and suffix",
			input:    "prefix-cx11-suffix",
			expected: "11",
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			classification := getServerClassification(testCase.input)
			if classification != testCase.expected {
				t.Errorf("expected %s, got %s", testCase.expected, classification)
			}
		})
	}
}

func TestGetAvailableServerTypesPerLocation(t *testing.T) {
	t.Skip("TODO")
}

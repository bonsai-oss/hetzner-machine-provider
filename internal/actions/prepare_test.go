package actions

import (
	"testing"

	"github.com/hetznercloud/hcloud-go/hcloud"
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

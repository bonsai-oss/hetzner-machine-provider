package actions

import (
	"testing"
	"time"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"
	"github.com/stretchr/testify/assert"

	"github.com/bonsai-oss/hetzner-machine-provider/internal/helper"
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

func TestImageSelection(t *testing.T) {
	// Prepare some test images
	images := []*hcloud.Image{
		{
			ID:        1,
			Name:      "ubuntu-18.04",
			OSVersion: "18.04",
			Created:   time.Now().Add(-72 * time.Hour),
		},
		{
			ID:        2,
			Name:      "ubuntu-20.04",
			OSVersion: "20.04",
			Created:   time.Now().Add(-24 * time.Hour),
		},
		{
			ID:        3,
			Name:      "ubuntu-21.04",
			OSVersion: "21.04",
			Created:   time.Now(),
		},
		{
			ID:      4,
			Name:    "testing-snapshot-d92",
			Type:    hcloud.ImageTypeSnapshot,
			Created: time.Now().Add(-48 * time.Hour),
			Labels: map[string]string{
				"test-environment": "development",
			},
		},
		{
			ID:      5,
			Name:    "testing-snapshot-222e",
			Type:    hcloud.ImageTypeSnapshot,
			Created: time.Now().Add(-40 * time.Hour),
			Labels: map[string]string{
				"test-environment": "development",
			},
		},
	}

	for _, testCase := range []struct {
		name           string
		imageSelector  string
		expectedImage  *hcloud.Image
		expectingError bool
	}{
		{
			name:          "snapshot image with specific label",
			imageSelector: "label#test-environment=development",
			expectedImage: images[4],
		},
		{
			name:          "select specific image by name",
			imageSelector: "ubuntu-20.04",
			expectedImage: images[1],
		},
		{
			name:          "select image by name with latest suffix",
			imageSelector: "ubuntu:latest",
			expectedImage: images[2],
		},
		{
			name:          "select image by name with latest suffix (more precise)",
			imageSelector: "ubuntu-1:latest",
			expectedImage: images[0],
		},
		{
			name: "no image found",
			// This image does not exist
			imageSelector:  "non-existing-image",
			expectingError: true,
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			image, err := imageSelection(helper.Filter(images, func(image *hcloud.Image) bool {
				if isLabelSelector(testCase.imageSelector) {
					return image.Labels != nil
				} else {
					return true
				}
			}), testCase.imageSelector)
			if testCase.expectingError {
				assert.Error(t, err)
				assert.Nil(t, image)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.expectedImage.Name, image.Name)
				assert.Equal(t, testCase.expectedImage.ID, image.ID)
			}
		})
	}
}

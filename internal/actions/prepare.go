package actions

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/hetznercloud/hcloud-go/hcloud"

	"github.com/bonsai-oss/hetzner-machine-provider/internal/helper"
)

type VMParams struct {
	Image    string
	Type     string
	Location string
}

const userData = `
#cloud-config
packages:
  - git
  - git-lfs
  - curl
runcmd:
  - curl -L --output /usr/local/bin/gitlab-runner "https://gitlab-runner-downloads.s3.amazonaws.com/latest/binaries/gitlab-runner-linux-amd64" && chmod +x /usr/local/bin/gitlab-runner
  - sed -i 's/#Port 22/Port 2222/g' /etc/ssh/sshd_config
  - systemctl restart sshd
  - echo -n "\n--- CI Server is ready ---" > /dev/tty1
`

func Prepare(client *hcloud.Client, jobID string, params VMParams) error {
	privateKey, pub, generateSSHKeyError := helper.GenerateSSHKeyPair()
	if generateSSHKeyError != nil {
		return generateSSHKeyError
	}

	fmt.Println("🔧 Create SSH key pair")
	// Create SSH key
	hcloudSSHKey, _, keyCreateError := client.SSHKey.Create(context.Background(), hcloud.SSHKeyCreateOpts{
		Name:      helper.ResourceName(jobID),
		PublicKey: pub,
		Labels: map[string]string{
			"managed-by": "hmp",
		},
	})
	if keyCreateError != nil {
		return keyCreateError
	}
	defer client.SSHKey.Delete(context.Background(), hcloudSSHKey)

	// Assign server labels from environment variables
	labels := map[string]string{"managed-by": "hmp"}
	assignLabels(labels, map[string]string{
		"commit-ref":  "CUSTOM_ENV_CI_COMMIT_REF_NAME",
		"commit-sha":  "CUSTOM_ENV_CI_COMMIT_SHA",
		"job-id":      "CUSTOM_ENV_CI_JOB_ID",
		"pipeline-id": "CUSTOM_ENV_CI_PIPELINE_ID",
		"project-id":  "CUSTOM_ENV_CI_PROJECT_ID",
		"tag":         "CUSTOM_ENV_CI_COMMIT_TAG",
	})

	image := hcloud.Image{}

	// if the image selector is a label selector, we need to get the image ID
	if strings.Contains(params.Image, "=") {
		images, _, imageGetError := client.Image.List(context.Background(), hcloud.ImageListOpts{
			Type: []hcloud.ImageType{hcloud.ImageTypeSnapshot, hcloud.ImageTypeSystem},
			ListOpts: hcloud.ListOpts{
				LabelSelector: params.Image,
			},
			Status: []hcloud.ImageStatus{hcloud.ImageStatusAvailable},
		})
		if imageGetError != nil {
			return imageGetError
		}
		if len(images) == 0 {
			return fmt.Errorf("no images found for label selector %+q", params.Image)
		}
		sort.SliceStable(images, func(i, j int) bool {
			return images[i].Created.After(images[j].Created)
		})
		image.ID = images[0].ID
	} else {
		image.Name = params.Image
	}

	fmt.Println("🔧 Create ci server")
	createResult, _, serverCreateError := client.Server.Create(context.Background(), hcloud.ServerCreateOpts{
		Name: helper.ResourceName(jobID),
		ServerType: &hcloud.ServerType{
			Name: params.Type,
		},
		Labels: labels,
		SSHKeys: []*hcloud.SSHKey{
			hcloudSSHKey,
		},
		Location: &hcloud.Location{
			Name: params.Location,
		},
		Image:    &image,
		UserData: userData,
	})

	if serverCreateError != nil {
		return serverCreateError
	}

	if createResult.Server == nil {
		return fmt.Errorf("server is not found")
	}

	state := helper.State{
		ServerAddress: createResult.Server.PublicNet.IPv4.IP.String(),
		SSHPrivateKey: privateKey,
	}

	return state.WriteToFile(helper.StatePath)
}

// assignLabels assigns values from environment variables to server labels
func assignLabels(labels map[string]string, labelEnvironmentVariableMapping map[string]string) {
	for label, environmentVariable := range labelEnvironmentVariableMapping {
		if value, variableIsSet := os.LookupEnv(environmentVariable); variableIsSet {
			labels[label] = value
		}
	}
}

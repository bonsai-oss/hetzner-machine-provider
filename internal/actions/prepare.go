package actions

import (
	"context"
	"fmt"

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
package_update: false
package_upgrade: false
runcmd:
  - systemctl stop sshd
  - apt update && apt install -y ca-certificates git git-lfs curl
  - curl -L --output /usr/local/bin/gitlab-runner "https://gitlab-runner-downloads.s3.amazonaws.com/latest/binaries/gitlab-runner-linux-amd64"
  - chmod +x /usr/local/bin/gitlab-runner
  - reboot
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

	fmt.Println("🔧 Create ci server")
	createResult, _, serverCreateError := client.Server.Create(context.Background(), hcloud.ServerCreateOpts{
		Name: helper.ResourceName(jobID),
		ServerType: &hcloud.ServerType{
			Name: params.Type,
		},
		Labels: map[string]string{"managed-by": "hmp"},
		SSHKeys: []*hcloud.SSHKey{
			hcloudSSHKey,
		},
		Location: &hcloud.Location{
			Name: params.Location,
		},
		Image: &hcloud.Image{
			Name: params.Image,
		},
		UserData: userData,
	})

	if serverCreateError != nil {
		return serverCreateError
	}

	if createResult.Server == nil {
		return fmt.Errorf("server is nil")
	}

	state := helper.State{
		ServerAddress: createResult.Server.PublicNet.IPv4.IP.String(),
		SSHPrivateKey: privateKey,
	}

	return state.WriteToFile(helper.StatePath)
}

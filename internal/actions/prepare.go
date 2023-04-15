package actions

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/hetznercloud/hcloud-go/hcloud"
	"golang.org/x/crypto/ssh"

	"github.com/bonsai-oss/hetzner-machine-provider/assets"
	"github.com/bonsai-oss/hetzner-machine-provider/internal/helper"
)

type VMParams struct {
	Image    string
	Type     string
	Location string
}

type PrepareOptions struct {
	JobID                    string
	WaitDeadline             time.Duration
	AdditionalAuthorizedKeys string
}

func Prepare(client *hcloud.Client, options PrepareOptions, params VMParams) error {
	privateKey, pub, generateSSHKeyError := helper.GenerateSSHKeyPair()
	if generateSSHKeyError != nil {
		return generateSSHKeyError
	}

	fmt.Println("🔐 Create SSH key pair")
	if pk, pkParseError := ssh.ParsePrivateKey([]byte(privateKey)); pkParseError != nil {
		return pkParseError
	} else {
		fmt.Printf("\t\tFingerprint: %+v\n\n", ssh.FingerprintLegacyMD5(pk.PublicKey()))
	}

	// Create SSH key
	hcloudSSHKey, _, keyCreateError := client.SSHKey.Create(context.Background(), hcloud.SSHKeyCreateOpts{
		Name:      helper.ResourceName(options.JobID),
		PublicKey: pub,
		Labels: map[string]string{
			"managed-by": "hmp",
		},
	})
	if keyCreateError != nil {
		return keyCreateError
	}
	defer client.SSHKey.Delete(context.Background(), hcloudSSHKey)

	fmt.Println("📠 Create CI server")
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

	serverType, _, serverTypeGetError := client.ServerType.GetByName(context.Background(), params.Type)
	if serverTypeGetError != nil {
		fmt.Printf("❌ Cannot fetch server information: %+q\n", serverTypeGetError)
		return serverTypeGetError
	}
	fmt.Printf(
		"\t\tType:  %+v [%s]\n\t\tImage: %+v\n",
		serverType.Description,
		determineArchitectureString(serverType.Architecture),
		params.Image,
	)

	userDataBuffer := &bytes.Buffer{}
	userData := map[string]any{
		"ssh_authorized_keys": strings.Split(options.AdditionalAuthorizedKeys, "\n"),
		"architecture":        determineArchitectureString(serverType.Architecture),
	}
	if userdataRenderError := assets.CloudInitTemplate.Execute(userDataBuffer, userData); userdataRenderError != nil {
		return userdataRenderError
	}

	createResult, _, serverCreateError := client.Server.Create(context.Background(), hcloud.ServerCreateOpts{
		Name:       helper.ResourceName(options.JobID),
		ServerType: serverType,
		Labels:     labels,
		SSHKeys: []*hcloud.SSHKey{
			hcloudSSHKey,
		},
		Location: &hcloud.Location{
			Name: params.Location,
		},
		Image:    &image,
		UserData: userDataBuffer.String(),
	})
	if serverCreateError != nil {
		fmt.Println("❌ Server creation failed")
		return serverCreateError
	}

	if createResult.Server == nil {
		fmt.Println("❌ Server creation failed")
		return fmt.Errorf("server is not found")
	}

	fmt.Printf("⏳ Waiting %s for server to be ready\n", options.WaitDeadline)

	waitDeadlineContext, cancel := context.WithTimeout(context.Background(), options.WaitDeadline)
	defer cancel()
	if waitReachableError := helper.WaitReachable(waitDeadlineContext, privateKey, createResult.Server.PublicNet.IPv4.IP.String()); waitReachableError != nil {
		return waitReachableError
	}
	fmt.Println("✅ Server created")

	state := helper.State{
		ServerAddress: createResult.Server.PublicNet.IPv4.IP.String(),
		SSHPrivateKey: privateKey,
	}

	return state.WriteToFile(helper.StatePath)
}

func determineArchitectureString(serverArchitecture hcloud.Architecture) string {
	switch serverArchitecture {
	case hcloud.ArchitectureX86:
		return "amd64"
	case hcloud.ArchitectureARM:
		return "arm64"
	}

	return "amd64"
}

// assignLabels assigns values from environment variables to server labels
func assignLabels(labels map[string]string, labelEnvironmentVariableMapping map[string]string) {
	for label, environmentVariable := range labelEnvironmentVariableMapping {
		if value, variableIsSet := os.LookupEnv(environmentVariable); variableIsSet {
			labelValid, labelValidationError := hcloud.ValidateResourceLabels(map[string]any{label: value})
			if labelValidationError != nil {
				fmt.Printf("\t\t⚠️ Label validation failed: %+q\n", labelValidationError)
				continue
			}
			if !labelValid {
				continue
			}

			labels[label] = value
		}
	}
}

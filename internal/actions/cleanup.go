package actions

import (
	"context"
	"fmt"
	"os"

	"github.com/hetznercloud/hcloud-go/v2/hcloud"

	"github.com/bonsai-oss/hetzner-machine-provider/internal/helper"
)

func Cleanup(client *hcloud.Client, jobID string) error {
	server, _, getServerError := client.Server.GetByName(context.Background(), helper.ResourceName(jobID))
	if getServerError != nil {
		return getServerError
	}
	if server == nil {
		return fmt.Errorf("server is not found")
	}

	if _, _, err := client.Server.DeleteWithResult(context.Background(), server); err != nil {
		return err
	}

	return os.Remove(helper.StatePath)
}

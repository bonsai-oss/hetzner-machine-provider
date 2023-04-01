package actions

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/avast/retry-go/v4"

	"hcloud-machine-provider/internal/helper"
)

func Exec(cmdFile, stageName string) error {
	state, readStateError := helper.ReadStateFromFile(helper.StatePath)
	if readStateError != nil {
		return readStateError
	}

	if state.SSHPrivateKey == "" || state.ServerAddress == "" {
		return fmt.Errorf("incomplete state")
	}

	finalError := retry.Do(
		func() error {
			sshClient, sshClientError := helper.NewSSHClient(state.SSHPrivateKey, state.ServerAddress, 22)
			if sshClientError != nil {
				return sshClientError
			}
			defer sshClient.Close()
			return sshClient.RunCommand(context.Background(), "true")
		},
		retry.OnRetry(func(n uint, err error) {
			fmt.Printf("⏳ retrying (%d): %s\n", n, err.Error())
		}),
		retry.Attempts(10),
		retry.Delay(5*time.Second),
		retry.LastErrorOnly(true),
	)

	if finalError != nil {
		return finalError
	}

	sshClient, sshClientError := helper.NewSSHClient(state.SSHPrivateKey, state.ServerAddress, 22)
	if sshClientError != nil {
		return sshClientError
	}
	defer sshClient.Close()

	content, _ := os.ReadFile(cmdFile)
	errod := sshClient.RunCommand(context.Background(), string(content))
	if errod != nil {
		return errod
	}

	return nil
}

package actions

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/avast/retry-go/v4"

	"github.com/bonsai-oss/hetzner-machine-provider/internal/helper"
)

func Exec(cmdFile, stageName string) error {
	state, readStateError := helper.ReadStateFromFile(helper.StatePath)
	if readStateError != nil {
		return readStateError
	}

	if state.SSHPrivateKey == "" || state.ServerAddress == "" {
		return fmt.Errorf("incomplete state")
	}

	waitDeadlineContext, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	if err := helper.WaitReachable(waitDeadlineContext, state.SSHPrivateKey, state.ServerAddress); err != nil {
		return err
	}

	var sshClient *helper.SSHClient

	clientConnectError := retry.Do(
		func() error {
			var sshClientError error
			sshClient, sshClientError = helper.NewSSHClient(state.SSHPrivateKey, state.ServerAddress, helper.CustomSSHPort)
			return sshClientError
		},
		retry.Attempts(3),
		retry.Delay(10*time.Second),
		retry.LastErrorOnly(true),
		retry.OnRetry(func(n uint, err error) {
			fmt.Printf("❌ failed to connect, retrying (%d): %s\n", n, err.Error())
		}),
	)
	if clientConnectError != nil {
		return clientConnectError
	}
	defer sshClient.Close()

	scriptContent, _ := os.ReadFile(cmdFile)
	return sshClient.RunCommand(context.Background(), string(scriptContent))
}

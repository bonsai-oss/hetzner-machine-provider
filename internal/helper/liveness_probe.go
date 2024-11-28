package helper

import (
	"context"
	"fmt"
	"time"

	"github.com/avast/retry-go/v4"
)

func CheckLivenessSSH(privateKey, serverAddress string) error {
	sshClient, sshClientError := NewSSHClient(privateKey, serverAddress, CustomSSHPort)
	if sshClientError != nil {
		return sshClientError
	}
	defer sshClient.Close()
	return sshClient.RunCommand(context.Background(), "true")
}

func WaitReachable(ctx context.Context, privateKey, serverAddress string) error {
	deadline, _ := ctx.Deadline()
	return retry.Do(
		func() error {
			return CheckLivenessSSH(privateKey, serverAddress)
		},
		retry.OnRetry(func(n uint, err error) {
			fmt.Printf("\t\tServer not ready yet: %+q ... retrying (%s remaining)\n", err.Error(), time.Until(deadline).Round(time.Second))
		}),
		retry.Attempts(0),
		retry.Delay(5*time.Second),
		retry.DelayType(retry.FixedDelay),
		retry.LastErrorOnly(true),
		retry.Context(ctx),
	)
}

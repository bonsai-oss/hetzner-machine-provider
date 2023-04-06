package helper

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"golang.org/x/crypto/ssh"
)

type SSHClient struct {
	client *ssh.Client
}

const CustomSSHPort = 2222

func NewSSHClient(privateKeyStr, serverIP string, port uint16) (*SSHClient, error) {
	client, err := connectSSH(privateKeyStr, serverIP, port)
	if err != nil {
		return nil, err
	}

	return &SSHClient{client}, nil
}

func connectSSH(privateKeyStr, serverIP string, port uint16) (*ssh.Client, error) {
	// Parse the private key
	privateKey, err := ssh.ParsePrivateKey([]byte(privateKeyStr))
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	// Create an SSH client configuration
	config := &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(privateKey),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         500 * time.Millisecond,
	}

	// Connect to the SSH server
	client, err := ssh.Dial("tcp", serverIP+":"+strconv.Itoa(int(port)), config)
	if err != nil {
		return nil, fmt.Errorf("cannot connect to SSH server: %w", err)
	}

	return client, nil
}

func (c *SSHClient) RunCommand(ctx context.Context, command string) error {
	// Create a session
	session, err := c.client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Close()

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	// Run the command
	err = session.Start(command)
	if err != nil {
		return fmt.Errorf("failed to run command: %w", err)
	}
	var waiter = make(chan error)

	// Wait for the command to finish or the context to be cancelled
	go func() {
		sessionError := session.Wait()
		waiter <- sessionError
	}()
	select {
	case <-ctx.Done():
		session.Signal(ssh.SIGTERM)
		fmt.Println("")
		return ctx.Err()
	case err := <-waiter:
		if err != nil {
			return fmt.Errorf("failed to wait for session: %w", err)
		}
	}

	return nil
}

func (c *SSHClient) Close() error {
	return c.client.Close()
}

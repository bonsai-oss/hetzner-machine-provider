package helper

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"golang.org/x/crypto/ssh"
)

func GenerateSSHKeyPair() (string, string, error) {
	// Generate the ECDSA key pair with P-256 curve
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate ecdsa key pair: %v", err)
	}

	// Marshal the private key into PKCS8 format
	pkcs8Key, err := x509.MarshalPKCS8PrivateKey(privKey)
	if err != nil {
		return "", "", fmt.Errorf("failed to marshal private key into PKCS8: %v", err)
	}

	// PEM encode the private key
	privKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: pkcs8Key,
	})

	// Convert the public key to SSH format
	sshPubKey, err := ssh.NewPublicKey(&privKey.PublicKey)
	if err != nil {
		return "", "", fmt.Errorf("failed to convert public key to ssh format: %v", err)
	}

	return string(privKeyPEM), string(ssh.MarshalAuthorizedKey(sshPubKey)), nil
}

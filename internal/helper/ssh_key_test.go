package helper_test

import (
	"fmt"
	"testing"

	"hcloud-machine-provider/internal/helper"
)

func TestGenerateSSHKeyPair(t *testing.T) {
	privKey, pubKey, err := helper.GenerateSSHKeyPair()
	if err != nil {
		t.Fatal(err)
	}

	if privKey == "" {
		t.Fatal("private key is empty")
	}

	if pubKey == "" {
		t.Fatal("public key is empty")
	}

	fmt.Println(privKey)
	fmt.Println(pubKey)
}

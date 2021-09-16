package commissioner

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"
)

var tokenName = "Vote"

func generateMintingScript(
	invalidHereafter int,
	policyVerificationKeyPath string,
) (string, error) {
	cmd := exec.Command(
		"cardano-cli", "address", "key-hash",
		"--payment-verification-key-file", policyVerificationKeyPath,
	)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}
	if err := cmd.Start(); err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(stdout)
	if err != nil {
		return "", err
	}
	if err := cmd.Wait(); err != nil {
		return "", fmt.Errorf("cardano-cli: %v", err)
	}
	keyHash := strings.TrimSpace(string(body))
	return fmt.Sprintf(
		`{"type": "sig", "keyHash": "%s"}`,
		keyHash,
	), nil
	return fmt.Sprintf(
		`{
	"type": "all",
	"scripts": [{
		"type": "before",
		"slot": %d,
	}, {
		"type": "sig",
		"keyHash": "%s"
	}]
}`,
		invalidHereafter,
		keyHash,
	), nil
}

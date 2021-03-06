package commissioner

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/thavlik/bvs/components/commissioner/pkg/api"
	"github.com/thavlik/bvs/components/commissioner/pkg/storage"
)

func (s *Server) CreateElection(
	ctx context.Context,
	req api.CreateElectionRequest,
) (*api.CreateElectionResponse, error) {
	fmt.Printf("CreateElection %#v\n", req)
	id := uuid.New().String()
	dir := filepath.Join("/tmp/policy/", id)
	if err := os.MkdirAll(dir, 0644); err != nil {
		return nil, err
	}
	defer func() {
		if err := os.RemoveAll(dir); err != nil {
			fmt.Printf("failed to delete policy files: %v\n", err)
		}
	}()
	protocolJsonPath := filepath.Join(dir, "protocol.json")
	if err := queryProtocol(protocolJsonPath); err != nil {
		return nil, fmt.Errorf("queryProtocol: %v", err)
	}
	protocolJson, err := ioutil.ReadFile(protocolJsonPath)
	if err != nil {
		return nil, err
	}
	pubKeyPath := filepath.Join(dir, "key.pub")
	privKeyPath := filepath.Join(dir, "key.priv")
	if err := (exec.Command(
		"cardano-cli", "address", "key-gen",
		"--verification-key-file", pubKeyPath,
		"--signing-key-file", privKeyPath,
	)).Run(); err != nil {
		return nil, err
	}
	signingKey, err := ioutil.ReadFile(privKeyPath)
	if err != nil {
		return nil, err
	}
	verificationKey, err := ioutil.ReadFile(pubKeyPath)
	if err != nil {
		return nil, err
	}
	currentSlot, err := getCurrentSlot()
	if err != nil {
		return nil, err
	}
	invalidHereafter := currentSlot + 31557600 // one year
	mintingScript, err := generateMintingScript(invalidHereafter, pubKeyPath)
	if err != nil {
		return nil, fmt.Errorf("generateMintingScript: %v", err)
	}
	scriptPath := filepath.Join(dir, "policy.script")
	if err := ioutil.WriteFile(
		scriptPath,
		[]byte(mintingScript),
		0644,
	); err != nil {
		return nil, err
	}
	policyIDPath := filepath.Join(dir, "policyID")
	if _, err := Exec(
		"bash", "-c",
		fmt.Sprintf(
			"cardano-cli transaction policyid --script-file %s >> %s",
			scriptPath,
			policyIDPath,
		),
	); err != nil {
		return nil, err
	}
	policyIDBytes, err := ioutil.ReadFile(policyIDPath)
	if err != nil {
		return nil, err
	}
	policyID := strings.TrimSpace(string(policyIDBytes))
	if err := s.storage.StoreElection(&storage.Election{
		ID:               id,
		SigningKey:       string(signingKey),
		VerificationKey:  string(verificationKey),
		PolicyID:         policyID,
		MintingScript:    mintingScript,
		InvalidHereafter: invalidHereafter,
		Protocol:         string(protocolJson),
	}); err != nil {
		return nil, fmt.Errorf("storage: %v", err)
	}
	var vkey api.Key
	if err := json.Unmarshal(verificationKey, &vkey); err != nil {
		return nil, fmt.Errorf("unmarshal: %v", err)
	}
	fmt.Printf("Created election %s\n", id)
	return &api.CreateElectionResponse{
		ID:              id,
		PolicyID:        policyID,
		VerificationKey: vkey,
	}, nil
}

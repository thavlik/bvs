package commissioner

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/thavlik/bvs/components/commissioner/pkg/api"
	"github.com/thavlik/bvs/components/commissioner/pkg/storage"
)

func (s *Server) CreateElection(
	ctx context.Context,
	req api.CreateElectionRequest,
) (*api.CreateElectionResponse, error) {
	fmt.Printf("CreateElection %#v\n", req)
	dir := filepath.Join("/tmp/policy/", req.Name)
	if err := os.MkdirAll(dir, 0644); err != nil {
		return nil, err
	}
	defer func() {
		if err := os.RemoveAll(dir); err != nil {
			fmt.Printf("failed to delete policy files: %v\n", err)
		}
	}()
	pubKeyPath := filepath.Join(dir, "key.pub")
	privKeyPath := filepath.Join(dir, "key.priv")
	cmd := exec.Command(
		"cardano-cli", "address", "key-gen",
		"--verification-key-file", pubKeyPath,
		"--signing-key-file", privKeyPath,
	)
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	signingKey, err := ioutil.ReadFile(pubKeyPath)
	if err != nil {
		return nil, err
	}
	verificationKey, err := ioutil.ReadFile(pubKeyPath)
	if err != nil {
		return nil, err
	}
	// Store info about election in database
	if err := s.storage.StoreElection(&storage.Election{
		ID:              req.Name,
		Deadline:        time.Now().Add(365 * 24 * time.Hour),
		SigningKey:      string(signingKey),
		VerificationKey: string(verificationKey),
	}); err != nil {
		return nil, fmt.Errorf("storage: %v", err)
	}
	var vkey api.Key
	if err := json.Unmarshal(verificationKey, &vkey); err != nil {
		return nil, fmt.Errorf("unmarshal: %v", err)
	}
	fmt.Printf("Created election, id=%s\n", req.Name)
	return &api.CreateElectionResponse{
		ID:              req.Name,
		VerificationKey: vkey,
	}, nil
}

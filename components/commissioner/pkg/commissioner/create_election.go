package commissioner

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/thavlik/bvs/components/commissioner/pkg/api"
	"go.uber.org/zap"
)

func (s *Server) CreateElection(
	ctx context.Context,
	req api.CreateElectionRequest,
) (*api.CreateElectionResponse, error) {
	s.log.Info("CreateElection",
		zap.String("req.Name", req.Name),
		zap.Int("req.Deadline", int(req.Deadline)))
	id := uuid.New().String()
	dir := filepath.Join("/tmp/policy/", id)
	if err := os.MkdirAll(dir, 0644); err != nil {
		return nil, err
	}
	defer func() {
		if err := os.RemoveAll(dir); err != nil {
			s.log.Error("failed to delete policy files", zap.Error(err))
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

	// TODO: store info about election in database
	if err := s.storage.Store(id, string(signingKey)); err != nil {
		return nil, fmt.Errorf("storage: %v", err)
	}

	verificationKeyData, err := ioutil.ReadFile(pubKeyPath)
	if err != nil {
		return nil, err
	}
	var vkey api.VerificationKey
	if err := json.Unmarshal(verificationKeyData, &vkey); err != nil {
		return nil, fmt.Errorf("unmarshal: %v", err)
	}
	return &api.CreateElectionResponse{
		ID:              id,
		VerificationKey: vkey,
	}, nil
}

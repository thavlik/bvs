package commissioner

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/thavlik/bvs/components/commissioner/pkg/api"
)

func (s *Server) CreateVoter(
	ctx context.Context,
	req api.CreateVoterRequest,
) (*api.CreateVoterResponse, error) {
	id := uuid.New().String()
	dir := filepath.Join("/tmp/%s", id)
	if err := os.MkdirAll(dir, 0777); err != nil {
		return nil, err
	}
	defer func() {
		if err := os.RemoveAll(dir); err != nil {
			fmt.Printf("failed to clean up temp dir: %v\n", err)
		}
	}()
	vkeyPath := filepath.Join(dir, "payment.vkey")
	skeyPath := filepath.Join(dir, "payment.skey")
	addrPath := filepath.Join(dir, "payment.addr")
	if _, err := Exec(
		"cardano-cli", "address", "key-gen",
		"--verification-key-file", vkeyPath,
		"--signing-key-file", skeyPath,
	); err != nil {
		return nil, err
	}
	if _, err := Exec(
		"cardano-cli", "address", "build",
		"--payment-verification-key-file", vkeyPath,
		"--out-file", addrPath,
		"--testnet-magic", fmt.Sprintf("%d", CardanoTestNetMagic),
	); err != nil {
		return nil, err
	}
	vkeyBytes, err := ioutil.ReadFile(vkeyPath)
	if err != nil {
		return nil, err
	}
	skeyBytes, err := ioutil.ReadFile(skeyPath)
	if err != nil {
		return nil, err
	}
	addrBytes, err := ioutil.ReadFile(addrPath)
	if err != nil {
		return nil, err
	}
	var vkey api.Key
	var skey api.Key
	if err := json.Unmarshal(vkeyBytes, &vkey); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(skeyBytes, &skey); err != nil {
		return nil, err
	}
	fmt.Printf("Created voter %s\n", id)
	return &api.CreateVoterResponse{
		SigningKey:      skey,
		VerificationKey: vkey,
		Address:         strings.TrimSpace(string(addrBytes)),
	}, nil
}

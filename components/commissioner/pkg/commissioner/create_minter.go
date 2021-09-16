package commissioner

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/thavlik/bvs/components/commissioner/pkg/api"
	"github.com/thavlik/bvs/components/commissioner/pkg/storage"
)

func (s *Server) CreateMinter(
	ctx context.Context,
	req api.CreateMinterRequest,
) (*api.CreateMinterResponse, error) {
	id := uuid.New().String()
	skey, err := json.Marshal(&req.SigningKey)
	if err != nil {
		return nil, err
	}
	vkey, err := json.Marshal(&req.VerificationKey)
	if err != nil {
		return nil, err
	}
	if err := s.storage.StoreMinter(&storage.Minter{
		ID:              id,
		SigningKey:      string(skey),
		VerificationKey: string(vkey),
		Address:         req.Address,
	}); err != nil {
		return nil, fmt.Errorf("storage: %v", err)
	}
	fmt.Printf("Created Minter %s\n", id)
	return &api.CreateMinterResponse{
		ID: id,
	}, nil
}

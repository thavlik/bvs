package node

import (
	"context"
	"sync/atomic"

	"github.com/thavlik/bvs/components/node/pkg/api"
)

func (s *Server) ProbeReady(ctx context.Context, req api.ProbeReadyRequest) (*api.ProbeReadyResponse, error) {
	return &api.ProbeReadyResponse{
		IsReady: atomic.LoadInt32(&s.isReady) == 1,
	}, nil
}

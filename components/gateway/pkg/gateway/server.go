package chillpill

import (
	"fmt"

	"go.uber.org/zap"
)

type Server struct {
	log *zap.Logger
}

func NewServer(
	log *zap.Logger,
) *Server {
	return &Server{
		log,
	}
}

func (s *Server) Listen(port int) error {
	return fmt.Errorf("unimplemented")
}

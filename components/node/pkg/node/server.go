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

func (s *Server) Start(
	port int,
	config,
	databasePath,
	socketPath,
	hostAddr,
	topology string,
) error {
	// TODO: start up cardano-node
	return fmt.Errorf("unimplemented")
}

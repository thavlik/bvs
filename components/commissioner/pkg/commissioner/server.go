package commissioner

import (
	"fmt"

	"github.com/thavlik/bvs/components/commissioner/pkg/storage"
	"go.uber.org/zap"
)

var CardanoTestNetMagic = 1097911063

// cardano-cli query tip --testnet-magic 1097911063

type Server struct {
	tokenName string
	storage   storage.Storage
	log       *zap.Logger
}

func NewServer(
	tokenName string,
	storage storage.Storage,
	log *zap.Logger,
) *Server {
	return &Server{
		tokenName,
		storage,
		log,
	}
}

func (s *Server) Start(
	port int,
	nodePort int,
	nodeConfig,
	nodeDatabasePath,
	nodeSocketPath,
	nodeHostAddr,
	nodeTopology string,
) error {
	nodeDone := make(chan error, 1)
	go func() {
		nodeDone <- s.startNode(
			nodePort,
			nodeConfig,
			nodeDatabasePath,
			nodeSocketPath,
			nodeHostAddr,
			nodeTopology,
		)
		close(nodeDone)
	}()
	serverDone := make(chan error, 1)
	go func() {
		serverDone <- s.listen(port)
		close(serverDone)
	}()
	select {
	case err := <-nodeDone:
		return fmt.Errorf("cardano-node: %v", err)
	case err := <-serverDone:
		return fmt.Errorf("listen: %v", err)
	}
}

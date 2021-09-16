package commissioner

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/thavlik/bvs/components/commissioner/pkg/storage"
)

var CardanoTestNetMagic = 1097911063

// cardano-cli query tip --testnet-magic 1097911063

type Server struct {
	tokenName string
	storage   storage.Storage
}

func NewServer(
	tokenName string,
	storage storage.Storage,
) *Server {
	return &Server{
		tokenName,
		storage,
	}
}

func (s *Server) Start(
	port int,
	metricsPort int,
) error {
	proxyDone := make(chan error, 1)
	go func() {
		proxyDone <- s.startProxyClient("bvs-node:2100")
		close(proxyDone)
	}()
	serverDone := make(chan error, 1)
	go func() {
		serverDone <- s.listen(port)
		close(serverDone)
	}()
	metricsDone := make(chan error, 1)
	if metricsPort == 0 {
		// Metrics are disabled
		defer close(metricsDone)
	} else {
		go func() {
			http.Handle("/metrics", promhttp.Handler())
			fmt.Printf("Metrics server listening on %d\n", metricsPort)
			metricsDone <- http.ListenAndServe(fmt.Sprintf(":%d", metricsPort), nil)
			close(metricsDone)
		}()
	}
	select {
	case err := <-proxyDone:
		return fmt.Errorf("proxy: %v", err)
	case err := <-serverDone:
		return fmt.Errorf("listen: %v", err)
	case err := <-metricsDone:
		return fmt.Errorf("metrics: %v", err)
	}
}

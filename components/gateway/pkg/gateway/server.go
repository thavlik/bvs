package gateway

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Start(
	port int,
	metricsPort int,
) error {
	serverDone := make(chan error, 1)
	go func() {
		fmt.Printf("Listening on %d\n", port)
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
	case err := <-serverDone:
		return fmt.Errorf("listen: %v", err)
	case err := <-metricsDone:
		return fmt.Errorf("metrics: %v", err)
	}
}

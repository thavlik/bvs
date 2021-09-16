package node

import (
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var CardanoTestNetMagic = 1097911063

// cardano-cli query tip --testnet-magic 1097911063

type Server struct {
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Start(
	proxyPort int,
	metricsPort int,
	nodePort int,
	nodeConfig,
	nodeDatabasePath,
	nodeSocketPath,
	nodeHostAddr,
	nodeTopology string,
) error {
	databaseLoaded := make(chan int, 1)
	startProxy := make(chan int, 1)
	nodeDone := make(chan error, 1)
	go func() {
		nodeDone <- s.startNode(
			nodePort,
			nodeConfig,
			nodeDatabasePath,
			nodeSocketPath,
			nodeHostAddr,
			nodeTopology,
			databaseLoaded,
		)
		close(nodeDone)
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
	fmt.Println("Waiting for cardano-node to fully synchronize...")
	fullySynced := make(chan error, 1)
	go func() {
		if err := waitForReady(databaseLoaded); err != nil {
			fullySynced <- fmt.Errorf("waitForReady: %v", err)
		} else {
			fullySynced <- nil
			startProxy <- 1
			close(startProxy)
		}
		close(fullySynced)
	}()
	stop := make(chan int, 1)
	go func() {
		<-fullySynced
		for {
			timer := time.After(10 * time.Second)
			select {
			case <-timer:
				tip, err := queryTip()
				if err != nil {
					fmt.Printf("queryTip error: %v\n", err)
				} else {
					if tip.SyncProgress == "100.00" {
						fmt.Println("cardano-node is currently in sync with the network")
					} else {
						fmt.Printf("cardano-node is %s%% synchronized with the network\n", tip.SyncProgress)
					}
				}
			case <-stop:
				return
			}
		}
	}()
	defer func() {
		stop <- 0
		close(stop)
	}()
	proxyDone := make(chan error, 1)
	go func() {
		<-startProxy
		//proxyDone <- s.startProxyServer(proxyPort)
		//close(proxyDone)
	}()
	select {
	case err := <-proxyDone:
		return fmt.Errorf("proxy: %v", err)
	case err := <-nodeDone:
		return fmt.Errorf("cardano-node: %v", err)
	case err := <-metricsDone:
		return fmt.Errorf("metrics: %v", err)
	}
}

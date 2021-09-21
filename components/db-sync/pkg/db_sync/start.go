package db_sync

import (
	"fmt"
	"time"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

func Start(nodeAddr string, cl client.Client) error {
	if err := WaitForNode(cl); err != nil {
		return fmt.Errorf("WaitForNode: %v", err)
	}
	//socatDone := make(chan error, 1)
	//go func() {
	//	fmt.Println("Starting socat")
	//	socatDone <- StartSocat(nodeAddr)
	//	close(socatDone)
	//}()
	dbSyncDone := make(chan error, 1)
	go func() {
		time.Sleep(3 * time.Second)
		fmt.Println("Starting cardano-db-sync...")
		dbSyncDone <- StartCardanoDBSync()
		close(dbSyncDone)
	}()
	select {
	case err := <-dbSyncDone:
		return fmt.Errorf("cardano-db-sync: %v", err)
		//case err := <-socatDone:
		//	return fmt.Errorf("socat: %v", err)
	}
}

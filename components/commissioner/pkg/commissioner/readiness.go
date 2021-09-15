package commissioner

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"
)

func waitForReady(databaseLoaded <-chan int) error {
	stop := make(chan int, 1)
	progress := "?"
	var progressL sync.Mutex
	message := "Waiting for cardano-node to load the blockchain data from local disk"
	go func() {
		start := time.Now()
		for {
			time.Sleep(10 * time.Second)
			progressL.Lock()
			m := message
			p := progress
			progressL.Unlock()
			fmt.Printf(
				"%s (%s%% complete, %s elapsed)\n",
				m,
				p,
				time.Since(start).Round(time.Second).String())
		}
	}()
	defer func() {
		stop <- 0
		close(stop)
	}()
	<-databaseLoaded
	progressL.Lock()
	message = "Waiting for cardano-node to fully synchronize with the network"
	progressL.Unlock()
	for {
		tip, err := queryTip()
		if err != nil {
			fmt.Printf("queryTip error: %v\n", err)
		} else {
			progressL.Lock()
			progress = tip.SyncProgress
			progressL.Unlock()
			if tip.SyncProgress == "100.00" {
				if err := signalReady(); err != nil {
					return err
				}
				return nil
			}
		}
		time.Sleep(5 * time.Second)
	}
}

func signalReady() error {
	probePath, ok := os.LookupEnv("READY_SIGNAL_PATH")
	if !ok {
		probePath = "/ready"
	}
	if err := os.MkdirAll(filepath.Dir(probePath), 0777); err != nil {
		return err
	}
	if err := ioutil.WriteFile(probePath, []byte("1"), 0644); err != nil {
		return err
	}
	return nil
}

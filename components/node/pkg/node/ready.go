package node

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

func waitForReady(databaseLoaded <-chan int) error {
	start := time.Now()
	stop := make(chan int, 1)
	var l sync.Mutex
	message := "Waiting for cardano-node to load the blockchain data from local disk. This usually takes 3-4 minutes."
	var epoch, slot, block int
	progress := "?"
	go func() {
		for {
			timer := time.After(10 * time.Second)
			select {
			case <-stop:
				return
			case <-timer:
				l.Lock()
				m := message
				e := epoch
				s := slot
				b := block
				p := progress
				l.Unlock()
				fmt.Printf(
					"%s (epoch %d, slot %d, block %d, %s%% complete, %s elapsed)\n",
					m, e, s, b, p,
					time.Since(start).Round(time.Second).String())
			}
		}
	}()
	defer func() {
		stop <- 0
		close(stop)
	}()
	<-databaseLoaded
	l.Lock()
	message = "Waiting for cardano-node to fully synchronize with the network"
	l.Unlock()
	for {
		tip, err := queryTip()
		if err != nil {
			fmt.Printf("queryTip error: %v\n", err)
		} else {
			l.Lock()
			progress = tip.SyncProgress
			epoch = tip.Epoch
			slot = tip.Slot
			block = tip.Block
			l.Unlock()
			if strings.HasPrefix(tip.SyncProgress, "100") {
				if err := signalReady(); err != nil {
					return err
				}
				fmt.Printf("Startup and synchronization took %v\n", time.Since(start).Round(time.Second).String())
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

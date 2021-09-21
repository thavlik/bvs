package db_sync

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"strings"
)

func StartCardanoDBSync() error {
	cmd := exec.Command(
		"cardano-db-sync",
		"--socket-path", socketPath,
		"--config", "/etc/db-sync/configs/testnet-config.yaml",
		"--state-dir", "/etc/db-sync/state/testnet",
		"--schema-dir", "/etc/db-sync/schema",
	)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	stdoutDone := make(chan error, 1)
	go func() {
		defer close(stdoutDone)
		stdoutDone <- func() error {
			scanner := bufio.NewScanner(stdout)
			for scanner.Scan() {
				fmt.Println(scanner.Text())
			}
			return scanner.Err()
		}()
	}()
	stderrDone := make(chan error, 1)
	go func() {
		defer close(stderrDone)
		stderrDone <- func() error {
			scanner := bufio.NewScanner(stderr)
			for scanner.Scan() {
				fmt.Println(scanner.Text())
			}
			return scanner.Err()
		}()
	}()
	if err := cmd.Wait(); err != nil {
		printErrorLog()
		return fmt.Errorf("cardano-db-sync: %v", err)
	}
	<-stdoutDone
	<-stderrDone
	return nil
}

func printErrorLog() {
	fmt.Println("cardano-db-sync exited unexpectedly. Retrieving error log...")
	files, err := ioutil.ReadDir("/tmp")
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".log") {
			path := filepath.Join("/tmp", file.Name())
			body, err := ioutil.ReadFile(path)
			if err != nil {
				panic(err)
			}
			fmt.Printf("Found error log at %s:\n", path)
			fmt.Println(string(body))
		}
	}
}

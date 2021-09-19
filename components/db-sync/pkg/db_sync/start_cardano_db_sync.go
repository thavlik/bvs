package db_sync

import (
	"bufio"
	"fmt"
	"os/exec"
)

func StartCardanoDBSync() error {
	cmd := exec.Command(
		"cardano-db-sync",
		"--socket-path", socketPath,
		"--config", "/configs/db-sync/mainnet-config.yaml",
		"--state-dir", "/etc/db-sync/state",
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
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("cardano-db-sync: %v", err)
	}
	<-stdoutDone
	<-stderrDone
	return nil
}

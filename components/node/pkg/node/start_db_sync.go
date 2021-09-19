package node

import (
	"bufio"
	"fmt"
	"os/exec"
)

func (s *Server) startDBSync() error {
	fmt.Println("Starting cardano-db-sync...")
	cmd := exec.Command(
		"cardano-db-sync", "run",
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
	if err := cmd.Start(); err != nil {
		return err
	}
	exit := make(chan error, 1)
	go func() {
		exit <- cmd.Wait()
		close(exit)
	}()
	select {
	case err := <-exit:
		return fmt.Errorf("cardano-db-sync: %v", err)
	case err := <-stdoutDone:
		return fmt.Errorf("cardano-db-sync stdout: %v", err)
	case err := <-stderrDone:
		return fmt.Errorf("cardano-db-sync stderr: %v", err)
	}
}

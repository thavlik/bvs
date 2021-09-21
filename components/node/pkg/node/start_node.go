package node

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"
)

func (s *Server) startNode(
	nodePort int,
	nodeConfig,
	nodeDatabasePath,
	nodeSocketPath,
	nodeHostAddr,
	nodeTopology string,
	databaseLoaded chan<- int,
) error {
	fmt.Println("Starting cardano-node...")
	cmd := exec.Command(
		"cardano-node", "run",
		"--config", nodeConfig,
		"--database-path", nodeDatabasePath,
		"--socket-path", nodeSocketPath,
		"--host-addr", nodeHostAddr,
		"--port", fmt.Sprintf("%d", nodePort),
		"--topology", nodeTopology,
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
	if err := cmd.Start(); err != nil {
		return err
	}
	go func() {
		defer close(stdoutDone)
		stdoutDone <- func() error {
			scanner := bufio.NewScanner(stdout)
			isReady := false
			for scanner.Scan() {
				text := scanner.Text()
				if !isReady {
					if strings.Contains(text, "Chain extended, new tip:") {
						isReady = true
						databaseLoaded <- 0
						close(databaseLoaded)
					}
				}
				fmt.Println(text)
			}
			if err := scanner.Err(); err != nil {
				return err
			}
			return nil
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
			if err := scanner.Err(); err != nil {
				return err
			}
			return nil
		}()
	}()
	//exit := make(chan error, 1)
	//go func() {
	//	exit <- cmd.Wait()
	//	close(exit)
	//}()
	//select {
	//case err := <-stdoutDone:
	//	return fmt.Errorf("stdout: %v", err)
	//case err := <-stderrDone:
	//	return fmt.Errorf("stderr: %v", err)
	//case err := <-exit:
	err = cmd.Wait()
	<-stdoutDone
	<-stderrDone
	return fmt.Errorf("exited prematurely: %v", err)
	//}
}

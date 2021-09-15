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
	isReady := false
	go func() {
		defer close(stdoutDone)
		stdoutDone <- func() error {
			scanner := bufio.NewScanner(stdout)
			for {
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
				fmt.Printf("WARNING: stdout error: %v\n", scanner.Err())
			}
		}()
	}()
	stderrDone := make(chan error, 1)
	go func() {
		defer close(stderrDone)
		stderrDone <- func() error {
			scanner := bufio.NewScanner(stderr)
			for {
				for scanner.Scan() {
					fmt.Println(scanner.Text())
				}
				fmt.Printf("WARNING: stderr error: %v\n", scanner.Err())
			}
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
	case err := <-stdoutDone:
		return fmt.Errorf("stdout: %v", err)
	case err := <-stderrDone:
		return fmt.Errorf("stderr: %v", err)
	case err := <-exit:
		return fmt.Errorf("exited prematurely: %v", err)
	}
}

package commissioner

import (
	"bufio"
	"fmt"
	"os/exec"
)

func (s *Server) startNode(
	nodePort int,
	nodeConfig,
	nodeDatabasePath,
	nodeSocketPath,
	nodeHostAddr,
	nodeTopology string,
) error {
	s.log.Info("Starting cardano-node...")
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
	started := make(chan int, 1)
	stdoutDone := make(chan error)
	go func() {
		defer close(stdoutDone)
		stdoutDone <- func() error {
			scanner := bufio.NewScanner(stdout)
			for scanner.Scan() {
				fmt.Println(scanner.Text())
			}
			return nil
		}()
	}()
	stderrDone := make(chan error)
	go func() {
		defer close(stderrDone)
		stderrDone <- func() error {
			scanner := bufio.NewScanner(stderr)
			for scanner.Scan() {
				fmt.Println(scanner.Text())
			}
			return fmt.Errorf("early exit")
		}()
	}()
	if err := cmd.Start(); err != nil {
		return err
	}
	exit := make(chan error)
	go func() {
		exit <- cmd.Wait()
	}()
	<-started
	select {
	case err := <-stdoutDone:
		return fmt.Errorf("stdout: %v", err)
	case err := <-stderrDone:
		return fmt.Errorf("stderr: %v", err)
	case err := <-exit:
		return fmt.Errorf("exited prematurely: %v", err)
	}
}

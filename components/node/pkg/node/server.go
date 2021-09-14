package chillpill

import (
	"bufio"
	"fmt"
	"os/exec"
	"path/filepath"

	"go.uber.org/zap"
)

type Server struct {
	log *zap.Logger
}

func NewServer(
	log *zap.Logger,
) *Server {
	return &Server{
		log,
	}
}

func (s *Server) Start(
	port int,
	config,
	databasePath,
	socketPath,
	hostAddr,
	topology string,
) error {
	s.log.Info("Starting cardano-node...")
	cmd := exec.Command(
		"cardano-node", "run",
		"--config", filepath.Join(config, "testnet-config.json"),
		"--database-path", databasePath,
		"--socket-path", filepath.Join(databasePath, "node.socket"),
		"--host-addr", "172.17.0.2",
		"--port", "1337",
		"--topology", filepath.Join(config, "testnet-topology.json"),
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
				text := scanner.Text()
				s.log.Info(text)
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
				s.log.Error(scanner.Text())
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

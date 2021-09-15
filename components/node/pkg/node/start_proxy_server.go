package node

import (
	"bufio"
	"fmt"
	"os/exec"
)

func (s *Server) startProxyServer(port int) error {
	fmt.Printf("Starting TCP proxy on port %d\n", port)
	cmd := exec.Command(
		"gocat", "unix-to-tcp",
		"--src", "/mnt/db/node.socket",
		"--dst", fmt.Sprintf("0.0.0.0:%d", port),
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
				fmt.Printf("[gocat.stdout] %s\n", scanner.Text())
			}
			return fmt.Errorf("%v", scanner.Err())
		}()
	}()
	stderrDone := make(chan error, 1)
	go func() {
		defer close(stderrDone)
		stderrDone <- func() error {
			scanner := bufio.NewScanner(stderr)
			for scanner.Scan() {
				fmt.Printf("[gocat.stderr] %s\n", scanner.Text())
			}
			return fmt.Errorf("%v", scanner.Err())
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

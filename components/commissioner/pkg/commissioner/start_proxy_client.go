package commissioner

import (
	"bufio"
	"fmt"
	"os/exec"
)

func (s *Server) startProxyClient(addr string) error {
	fmt.Printf("Starting TCP proxy client for %s\n", addr)
	cmd := exec.Command(
		"bash",
		"-c",
		fmt.Sprintf(
			"socat -t 100000 -v UNIX-LISTEN:/tmp/node.socket,unlink-early,mode=777,fork TCP:%s",
			addr,
		),
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
			return fmt.Errorf("%v", scanner.Err())
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

package node

import (
	"fmt"
	"os/exec"
)

func (s *Server) startProxyServer(port int) error {
	fmt.Printf("Starting TCP proxy on port %d\n", port)
	cmd := exec.Command(
		"bash",
		"-c",
		fmt.Sprintf(
			"socat -d TCP-LISTEN:%d,reuseaddr,fork UNIX-CLIENT:/shared/node.socket",
			port,
		),
	)
	// This command works so reliably that I've taken out logging completely
	if err := cmd.Run(); err != nil {
		return err
	}
	return fmt.Errorf("exited prematurely")
	/*
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
					//fmt.Println(scanner.Text())
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
	*/
	//if err := cmd.Start(); err != nil {
	//	return err
	//}
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
	//	return fmt.Errorf("exited prematurely: %v", err)
	//}

}

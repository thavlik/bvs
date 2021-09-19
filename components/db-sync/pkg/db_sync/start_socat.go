package db_sync

import (
	"bufio"
	"fmt"
	"os/exec"
)

var socketPath = "/tmp/node.socket"

func StartSocat(nodeAddr string) error {
	fmt.Printf("Starting proxy from unix://%s to tcp://%s\n", socketPath, nodeAddr)
	cmd := exec.Command(
		"bash",
		"-c",
		fmt.Sprintf(
			"socat -t 100000 -v UNIX-LISTEN:%s,unlink-early,mode=777,fork TCP:%s",
			socketPath,
			nodeAddr,
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
		return err
	}
	return fmt.Errorf("exited prematurely")
}

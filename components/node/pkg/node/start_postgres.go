package node

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
)

func getPostgresPath() (string, error) {
	v, ok := os.LookupEnv("POSTGRES_BIN_PATH")
	if !ok {
		return "", fmt.Errorf("missing POSTGRES_BIN_PATH")
	}
	return v, nil
}

func (s *Server) startPostgres(dbPath string, postgresPort int) error {
	fmt.Printf("Starting postgres on port %d...\n", postgresPort)
	binPath, err := getPostgresPath()
	if err != nil {
		return err
	}
	cmd := exec.Command(
		"runuser",
		"-l", "postgres",
		"-c", fmt.Sprintf(
			"%s -D %s -p %d",
			binPath,
			dbPath,
			postgresPort,
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
				fmt.Println("[postgres|stdout] " + scanner.Text())
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
				fmt.Println("[postgres|stderr] " + scanner.Text())
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
		return fmt.Errorf("postgres: %v", err)
	case err := <-stdoutDone:
		return fmt.Errorf("postgres stdout: %v", err)
	case err := <-stderrDone:
		return fmt.Errorf("postgres stderr: %v", err)
	}
}

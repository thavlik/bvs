package commissioner

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"
)

func Exec(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return "", err
	}
	if err := cmd.Start(); err != nil {
		return "", err
	}
	stdoutBytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		return "", err
	}
	stderrBytes, err := ioutil.ReadAll(stderr)
	if err != nil {
		return "", err
	}
	if err := cmd.Wait(); err != nil {
		fmt.Println(string(stdoutBytes))
		fmt.Println(string(stderrBytes))
		return "", fmt.Errorf("%s: %v", command, err)
	}
	return strings.TrimSpace(string(stdoutBytes)), nil
}

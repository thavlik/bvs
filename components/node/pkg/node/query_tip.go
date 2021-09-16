package node

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/exec"
)

type Tip struct {
	Epoch        int    `json:"epoch"`
	Hash         string `json:"hash"`
	Slot         int    `json:"slot"`
	Block        int    `json:"block"`
	Era          string `json:"era"`
	SyncProgress string `json:"syncProgress"`
}

func queryTip() (*Tip, error) {
	cmd := exec.Command(
		"cardano-cli", "query", "tip",
		"--testnet-magic", fmt.Sprintf("%d", CardanoTestNetMagic),
	)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("start cardano-cli: %v", err)
	}
	stderrBytes, err := ioutil.ReadAll(stderr)
	if err != nil {
		return nil, err
	}
	stdoutBytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		return nil, err
	}
	if err := cmd.Wait(); err != nil {
		fmt.Println(string(stdoutBytes))
		fmt.Println(string(stderrBytes))
		return nil, fmt.Errorf("cardano-cli: %v", err)
	}
	tip := &Tip{}
	if err := json.Unmarshal(stdoutBytes, tip); err != nil {
		return nil, fmt.Errorf("json: %v", err)
	}
	return tip, nil
}

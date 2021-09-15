package node

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"strconv"
)

func calculateFee(
	rawTxPath string,
) (int, error) {
	cmd := exec.Command(
		"bash", "-c",
		fmt.Sprintf(
			`cardano-cli transaction calculate-min-fee \
				--tx-body-file %s \
				--tx-in-count 1 \
				--tx-out-count 1 \
				--witness-count 1 \
				--mainnet \
				--protocol-params-file protocol.json \
			| cut -d " " -f1)`,
			rawTxPath,
		),
	)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return 0, err
	}
	if err := cmd.Run(); err != nil {
		return 0, fmt.Errorf("cardano-cli: %v", err)
	}
	output, err := ioutil.ReadAll(stdout)
	if err != nil {
		return 0, err
	}
	fee, err := strconv.ParseInt(string(output), 10, 64)
	if err != nil {
		return 0, err
	}
	return int(fee), nil
}

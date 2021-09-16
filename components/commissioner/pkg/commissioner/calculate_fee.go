package commissioner

import (
	"fmt"
	"strconv"
	"strings"
)

func calculateFee(rawTxPath, protocolJsonPath string) (int, error) {
	stdout, err := Exec(
		"bash", "-c",
		fmt.Sprintf(
			`cardano-cli transaction calculate-min-fee \
				--tx-body-file %s \
				--tx-in-count 1 \
				--tx-out-count 1 \
				--witness-count 1 \
				--testnet-magic %d \
				--protocol-params-file %s`,
			rawTxPath,
			CardanoTestNetMagic,
			protocolJsonPath,
		),
	)
	if err != nil {
		return 0, err
	}
	fee, err := strconv.ParseInt(strings.Split(stdout, " ")[0], 10, 64)
	if err != nil {
		return 0, err
	}
	return int(fee), nil
}

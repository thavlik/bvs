package commissioner

import "fmt"

func queryProtocol(outPath string) error {
	if _, err := Exec(
		"cardano-cli", "query", "protocol-parameters",
		"--testnet-magic", fmt.Sprintf("%d", CardanoTestNetMagic),
		"--out-file", outPath,
	); err != nil {
		return err
	}
	return nil
}

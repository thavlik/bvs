package commissioner

import (
	"context"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strconv"
	"strings"

	"github.com/thavlik/bvs/components/commissioner/pkg/api"
)

var errNoTransactions = fmt.Errorf("no previous transactions")

type addressInfo struct {
	txHash   string
	txIx     int
	lovelace int
	balance  string
}

func queryAddress(address string) ([]*addressInfo, error) {
	cmd := exec.Command(
		"bash",
		"-c",
		fmt.Sprintf(`cardano-cli query utxo --address %s --testnet-magic %d`, address, CardanoTestNetMagic),
	)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(stdout)
	if err != nil {
		return nil, err
	}
	if err := cmd.Wait(); err != nil {
		return nil, fmt.Errorf("cardano-cli: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(string(body)), "\n")
	if len(lines) < 3 {
		return nil, errNoTransactions
	}
	var utxos []*addressInfo
	for _, tx := range lines[2:] {
		parts := strings.Split(tx, " ")
		var filtered []string
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if len(part) > 0 {
				filtered = append(filtered, part)
			}
		}
		txHash := filtered[0]
		txIx, err := strconv.ParseInt(filtered[1], 10, 64)
		if err != nil {
			return nil, err
		}
		lovelace, err := strconv.ParseInt(filtered[2], 10, 64)
		if err != nil {
			return nil, err
		}
		// If the unit is displayed, make sure it's lovelace, otherwise just assume
		if len(filtered) > 3 && filtered[3] != "lovelace" {
			return nil, fmt.Errorf("error querying address, expected 'lovelace', got '%s'", filtered[3])
		}
		var balance string
		for _, item := range filtered[2:] {
			if len(balance) > 0 {
				balance += " "
			}
			balance += item
		}
		utxos = append(utxos, &addressInfo{
			txHash:   txHash,
			txIx:     int(txIx),
			lovelace: int(lovelace),
			balance:  balance,
		})
	}
	return utxos, nil
}

func (s *Server) QueryAddress(ctx context.Context, req api.QueryAddressRequest) (*api.QueryAddressResponse, error) {
	utxos, err := queryAddress(req.Address)
	if err != nil {
		return nil, err
	}
	var result []*api.UnspentTransaction
	for _, utxo := range utxos {
		result = append(result, &api.UnspentTransaction{
			TxHash:   utxo.txHash,
			TxIx:     utxo.txIx,
			Lovelace: utxo.lovelace,
			Balance:  utxo.balance,
		})
	}
	return &api.QueryAddressResponse{
		UnspentTransactions: result,
	}, nil
}

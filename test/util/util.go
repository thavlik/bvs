package util

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/thavlik/bvs/components/commissioner/pkg/api"

	"github.com/stretchr/testify/require"

	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

var CardanoNetworkTimeout = 5 * time.Minute

func CreateClient(t *testing.T) client.Client {
	s := scheme.Scheme
	cl, err := client.New(config.GetConfigOrDie(), client.Options{Scheme: s})
	require.NoError(t, err)
	return cl
}

func WaitForUTXO(com api.Commissioner, address string) ([]*api.UnspentTransaction, error) {
	start := time.Now()
	for time.Since(start) < CardanoNetworkTimeout {
		resp, err := com.QueryAddress(context.TODO(), api.QueryAddressRequest{
			Address: address,
		})
		if err != nil {
			return nil, err
		}
		if len(resp.UnspentTransactions) > 0 {
			return resp.UnspentTransactions, nil
		}
		time.Sleep(5 * time.Second)
	}
	return nil, fmt.Errorf("failed to observe utxo after %v", CardanoNetworkTimeout)
}

func CountVotes(com api.Commissioner, address string, policyID string) (int, error) {
	utxos, err := com.QueryAddress(context.TODO(), api.QueryAddressRequest{
		Address: address,
	})
	if err != nil {
		return 0, fmt.Errorf("commissioner: %v", err)
	}
	count := 0
	for _, utxo := range utxos.UnspentTransactions {
		idx := strings.Index(utxo.Balance, policyID+".Vote")
		if idx != -1 {
			count++
		}
	}
	return count, nil
}

func WaitForUTXOChange(com api.Commissioner, oldUtxos []*api.UnspentTransaction, address string) error {
	start := time.Now()
	for time.Since(start) < CardanoNetworkTimeout {
		newUtxos, err := WaitForUTXO(com, address)
		if err != nil {
			return err
		}
		if len(newUtxos) != len(oldUtxos) {
			return nil
		}
		for i, newUtxo := range newUtxos {
			if newUtxo.Balance != oldUtxos[i].Balance || newUtxo.TxHash != oldUtxos[i].TxHash {
				return nil
			}
		}
		time.Sleep(5 * time.Second)
	}
	return fmt.Errorf("failed to observe utxo change after %v", CardanoNetworkTimeout)
}

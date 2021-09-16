package util

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/thavlik/bvs/components/commissioner/pkg/api"

	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"

	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

func CreateClient(t *testing.T) client.Client {
	s := scheme.Scheme
	cl, err := client.New(config.GetConfigOrDie(), client.Options{Scheme: s})
	require.NoError(t, err)
	return cl
}

type Wallet struct {
	SigningKey      string
	VerificationKey string
	Address         string
}

func GetWallet(t *testing.T) *Wallet {
	cl := CreateClient(t)
	secret := &corev1.Secret{}
	require.NoError(t, cl.Get(context.TODO(), types.NamespacedName{
		Name:      "test-wallet-cred",
		Namespace: "default",
	}, secret))
	skey := string(secret.Data["payment.skey"])
	vkey := string(secret.Data["payment.vkey"])
	addr := string(secret.Data["addr"])
	return &Wallet{
		SigningKey:      skey,
		VerificationKey: vkey,
		Address:         addr,
	}
}

func WaitForBalance(com api.Commissioner, address string) ([]*api.UnspentTransaction, error) {
	timeout := 3 * time.Minute
	start := time.Now()
	for {
		resp, err := com.QueryAddress(context.TODO(), api.QueryAddressRequest{
			Address: address,
		})
		if err != nil {
			if strings.Contains(err.Error(), "no previous transactions") {
				if time.Since(start) > timeout {
					return nil, fmt.Errorf("failed to witness transaction history after %s", timeout.String())
				}
				time.Sleep(5 * time.Second)
				continue
			}
			return nil, fmt.Errorf("commissioner: %v", err)
		}
		// We have an unspent transaction we can use
		return resp.UnspentTransactions, nil
	}
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

func WaitForBalanceChange(com api.Commissioner, oldUtxos []*api.UnspentTransaction, address string) error {
	timeout := 5 * time.Minute
	start := time.Now()
	for time.Since(start) < timeout {
		newUtxos, err := WaitForBalance(com, address)
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
		//fmt.Printf("Waiting for balance in %s to change\n", address)
		time.Sleep(5 * time.Second)
	}
	return fmt.Errorf("failed to observe balance change after %s", timeout.String())
}

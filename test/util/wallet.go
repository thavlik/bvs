package util

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

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

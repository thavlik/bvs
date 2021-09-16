package util

import (
	"context"
	"testing"

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

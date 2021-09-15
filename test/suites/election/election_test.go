package cli

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/thavlik/bvs/components/commissioner/pkg/api"
)

func CreateTestCommissioner(t *testing.T) api.Commissioner {
	uri, ok := os.LookupEnv("COMMISSIONER_URI")
	require.Truef(t, ok, "missing COMMISSIONER_URI")
	return api.NewCommissionerClient(uri, "", "", 30*time.Second)
}

func TestElectionCreate(t *testing.T) {
	com := CreateTestCommissioner(t)
	name := "bvstest" + uuid.New().String()[:5]
	resp, err := com.CreateElection(context.TODO(), api.CreateElectionRequest{
		Name:     name,
		Deadline: time.Now().Add(24 * time.Hour).UnixNano(),
	})
	require.NoError(t, err)
	fmt.Println(resp.ID)
	fmt.Println(resp.VerificationKey)
	require.NotEmpty(t, resp.ID)
	require.NotEmpty(t, resp.VerificationKey)
}

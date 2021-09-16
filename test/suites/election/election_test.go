package cli

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/thavlik/bvs/components/commissioner/pkg/api"
	"github.com/thavlik/bvs/test/util"
)

func CreateTestCommissioner(t *testing.T) api.Commissioner {
	uri, ok := os.LookupEnv("COMMISSIONER_URI")
	require.Truef(t, ok, "missing COMMISSIONER_URI")
	return api.NewCommissionerClient(uri, "", "", 30*time.Second)
}

func TestElectionCreate(t *testing.T) {
	com := CreateTestCommissioner(t)
	name := "bvstest" + uuid.New().String()[:5]
	_, err := com.CreateElection(context.TODO(), api.CreateElectionRequest{
		Name:     name,
		Deadline: time.Now().Add(24 * time.Hour).UnixNano(),
	})
	require.NoError(t, err)
	w := util.GetWallet(t)
	req := api.CreateMinterRequest{}
	require.NoError(t, json.Unmarshal([]byte(w.VerificationKey), &req.VerificationKey))
	require.NoError(t, json.Unmarshal([]byte(w.SigningKey), &req.SigningKey))
	resp, err := com.CreateMinter(context.TODO(), req)
	require.NoError(t, err)
	_, err = com.MintVote(context.TODO(), api.MintVoteRequest{
		Election: name,
		Voter:    "",
		Ident:    "",
		Auditor: api.Auditor{
			Agent:     resp.ID,
			Timestamp: time.Now().Unix(),
		},
	})
	require.NoError(t, err)
}

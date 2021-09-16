package cli

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

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
	election, err := com.CreateElection(context.TODO(), api.CreateElectionRequest{})
	require.NoError(t, err)
	w := util.GetWallet(t)
	req := api.CreateMinterRequest{Address: w.Address}
	require.NoError(t, json.Unmarshal([]byte(w.VerificationKey), &req.VerificationKey))
	require.NoError(t, json.Unmarshal([]byte(w.SigningKey), &req.SigningKey))
	minter, err := com.CreateMinter(context.TODO(), req)
	require.NoError(t, err)
	voter, err := com.CreateVoter(context.TODO(), api.CreateVoterRequest{})
	require.NoError(t, err)
	//body, err := json.Marshal(&voter.VerificationKey)
	//require.NoError(t, err)
	//fmt.Println(string(body))
	//body, err = json.Marshal(&voter.SigningKey)
	//require.NoError(t, err)
	//fmt.Println(string(body))
	//fmt.Println(voter.Address)
	_, err = com.MintVote(context.TODO(), api.MintVoteRequest{
		Election: election.ID,
		Voter:    voter.Address,
		Ident:    "TODO",
		Auditor: api.Auditor{
			Agent:     minter.ID,
			Timestamp: time.Now().Unix(),
			Proof:     "TODO",
		},
	})
	require.NoError(t, err)
}

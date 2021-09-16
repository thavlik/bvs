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

var tokenName = "Vote" // TODO: replace globally

func CreateTestCommissioner(t *testing.T) api.Commissioner {
	uri, ok := os.LookupEnv("COMMISSIONER_URI")
	require.Truef(t, ok, "missing COMMISSIONER_URI")
	return api.NewCommissionerClient(uri, "", "", 30*time.Second)
}

func TestSingleVote(t *testing.T) {
	com := CreateTestCommissioner(t)
	election, err := com.CreateElection(context.TODO(), api.CreateElectionRequest{})
	require.NoError(t, err)
	w := util.GetWallet(t)
	req := api.CreateMinterRequest{Address: w.Address}
	require.NoError(t, json.Unmarshal([]byte(w.VerificationKey), &req.VerificationKey))
	require.NoError(t, json.Unmarshal([]byte(w.SigningKey), &req.SigningKey))
	minter, err := com.CreateMinter(context.TODO(), req)
	require.NoError(t, err)
	candidate, err := com.CreateVoter(context.TODO(), api.CreateVoterRequest{})
	require.NoError(t, err)
	voter, err := com.CreateVoter(context.TODO(), api.CreateVoterRequest{})
	require.NoError(t, err)
	_, err = com.MintVote(context.TODO(), api.MintVoteRequest{
		Election: election.ID,
		Voter:    voter.Address,
		Minter:   minter.ID,
	})
	require.NoError(t, err)
	_, err = util.WaitForBalance(com, voter.Address)
	require.NoError(t, err)
	_, err = com.CastVote(context.TODO(), api.CastVoteRequest{
		Election:   election.ID,
		Voter:      voter.Address,
		SigningKey: voter.SigningKey,
		Candidate:  candidate.Address,
	})
	require.NoError(t, err)
	_, err = util.WaitForBalance(com, candidate.Address)
	require.NoError(t, err)
}

func TestMultipleVotes(t *testing.T) {
	com := CreateTestCommissioner(t)
	election, err := com.CreateElection(context.TODO(), api.CreateElectionRequest{})
	require.NoError(t, err)
	w := util.GetWallet(t)
	req := api.CreateMinterRequest{Address: w.Address}
	require.NoError(t, json.Unmarshal([]byte(w.VerificationKey), &req.VerificationKey))
	require.NoError(t, json.Unmarshal([]byte(w.SigningKey), &req.SigningKey))
	minter, err := com.CreateMinter(context.TODO(), req)
	require.NoError(t, err)
	candidate, err := com.CreateVoter(context.TODO(), api.CreateVoterRequest{})
	require.NoError(t, err)
	numVoters := 2
	for i := 0; i < numVoters; i++ {
		voter, err := com.CreateVoter(context.TODO(), api.CreateVoterRequest{})
		require.NoError(t, err)
		_, err = com.MintVote(context.TODO(), api.MintVoteRequest{
			Election: election.ID,
			Voter:    voter.Address,
			Minter:   minter.ID,
		})
		require.NoError(t, err)
		oldUtxos, err := util.WaitForBalance(com, voter.Address)
		require.NoError(t, err)
		_, err = com.CastVote(context.TODO(), api.CastVoteRequest{
			Election:   election.ID,
			Voter:      voter.Address,
			SigningKey: voter.SigningKey,
			Candidate:  candidate.Address,
		})
		require.NoError(t, err)
		require.NoError(t, util.WaitForBalanceChange(com, oldUtxos, voter.Address))
	}
	numVotes, err := util.CountVotes(com, candidate.Address, election.PolicyID)
	require.NoError(t, err)
	require.Equal(t, numVoters, numVotes)
}

func TestParallelVotes(t *testing.T) {
	com := CreateTestCommissioner(t)
	election, err := com.CreateElection(context.TODO(), api.CreateElectionRequest{})
	require.NoError(t, err)
	w := util.GetWallet(t)
	req := api.CreateMinterRequest{Address: w.Address}
	require.NoError(t, json.Unmarshal([]byte(w.VerificationKey), &req.VerificationKey))
	require.NoError(t, json.Unmarshal([]byte(w.SigningKey), &req.SigningKey))
	minter, err := com.CreateMinter(context.TODO(), req)
	require.NoError(t, err)
	candidate, err := com.CreateVoter(context.TODO(), api.CreateVoterRequest{})
	require.NoError(t, err)
	numVoters := 3
	dones := make([]chan error, numVoters)
	for i := 0; i < numVoters; i++ {
		voter, err := com.CreateVoter(context.TODO(), api.CreateVoterRequest{})
		require.NoError(t, err)
		_, err = com.MintVote(context.TODO(), api.MintVoteRequest{
			Election: election.ID,
			Voter:    voter.Address,
			Minter:   minter.ID,
		})
		require.NoError(t, err)
		oldUtxos, err := util.WaitForBalance(com, voter.Address)
		require.NoError(t, err)
		done := make(chan error, 1)
		dones[i] = done
		go func(oldUtxos []*api.UnspentTransaction, done chan<- error) {
			defer close(done)
			done <- func() error {
				if _, err = com.CastVote(context.TODO(), api.CastVoteRequest{
					Election:   election.ID,
					Voter:      voter.Address,
					SigningKey: voter.SigningKey,
					Candidate:  candidate.Address,
				}); err != nil {
					return err
				}
				if err := util.WaitForBalanceChange(com, oldUtxos, voter.Address); err != nil {
					return err
				}
				return nil
			}()
		}(oldUtxos, done)
	}
	for _, done := range dones {
		require.NoError(t, <-done)
	}
	numVotes, err := util.CountVotes(com, candidate.Address, election.PolicyID)
	require.NoError(t, err)
	require.Equal(t, numVoters, numVotes)
}

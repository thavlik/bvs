package commissioner

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/thavlik/bvs/components/commissioner/pkg/api"
)

func (s *Server) CastVote(ctx context.Context, req api.CastVoteRequest) (*api.CastVoteResponse, error) {
	// Create temporary directory for interacting with cardano-cli
	id := uuid.New().String()

	election, err := s.storage.RetrieveElection(req.Election)
	if err != nil {
		return nil, fmt.Errorf("storage: %v", err)
	}

	rootDir := filepath.Join("/tmp", id)
	if err := os.MkdirAll(rootDir, 0777); err != nil {
		return nil, err
	}
	defer func() {
		if err := os.RemoveAll(rootDir); err != nil {
			fmt.Printf("failed to clean up temp dir: %v\n", err)
		}
	}()
	paymentSigningKey, err := json.Marshal(&req.SigningKey)
	if err != nil {
		return nil, err
	}
	paymentSigningKeyPath := filepath.Join(rootDir, "payment.skey")
	if err := ioutil.WriteFile(
		paymentSigningKeyPath,
		paymentSigningKey,
		0644,
	); err != nil {
		return nil, err
	}
	protocolJsonPath := filepath.Join(rootDir, "protocol.json")
	if err := ioutil.WriteFile(
		protocolJsonPath,
		[]byte(election.Protocol),
		0644,
	); err != nil {
		return nil, err
	}
	utxos, err := queryAddress(req.Voter)
	if err != nil {
		return nil, fmt.Errorf("queryAddress: %v", err)
	}
	walletInfo := utxos[0]
	sendAmount := 2 * minUTxOValue
	walletInfo.lovelace -= sendAmount
	rawTxPath := filepath.Join(rootDir, "matx.raw")
	if _, err := Exec(
		"cardano-cli", "transaction", "build-raw",
		"--fee", "0",
		"--tx-in", fmt.Sprintf("%s#%d", walletInfo.txHash, walletInfo.txIx),
		"--tx-out", fmt.Sprintf(`%s + %d`, req.Voter, walletInfo.lovelace),
		"--tx-out", fmt.Sprintf(`%s + %d + 1 %s.%s`, req.Candidate, sendAmount, election.PolicyID, tokenName),
		"--out-file", rawTxPath,
	); err != nil {
		return nil, fmt.Errorf("draft tx: %v", err)
	}
	fee, err := calculateFee(rawTxPath, protocolJsonPath)
	if err != nil {
		return nil, fmt.Errorf("calculateFee: %v", err)
	}
	if walletInfo.lovelace-fee < minUTxOValue {
		return nil, fmt.Errorf("voter has insufficient funds (has %d, needed >=%d)", walletInfo.lovelace, fee+minUTxOValue)
	}
	walletInfo.lovelace -= fee
	// Rebuild the transaction with the calculated fee
	if _, err := Exec(
		"cardano-cli", "transaction", "build-raw",
		"--fee", fmt.Sprintf("%d", fee),
		"--tx-in", fmt.Sprintf("%s#%d", walletInfo.txHash, walletInfo.txIx),
		"--tx-out", fmt.Sprintf(`%s + %d`, req.Voter, walletInfo.lovelace),
		"--tx-out", fmt.Sprintf(`%s + %d + 1 %s.%s`, req.Candidate, sendAmount, election.PolicyID, tokenName),
		"--out-file", rawTxPath,
	); err != nil {
		return nil, fmt.Errorf("draft tx: %v", err)
	}
	// Sign the transaction
	signedTxPath := filepath.Join(rootDir, "matx.signed")
	if _, err := Exec(
		"cardano-cli", "transaction", "sign",
		"--signing-key-file", paymentSigningKeyPath,
		"--mainnet",
		"--tx-body-file", rawTxPath,
		"--out-file", signedTxPath,
	); err != nil {
		return nil, fmt.Errorf("sign: %v", err)
	}
	// Submit the transaction
	if _, err := Exec(
		"cardano-cli", "transaction", "submit",
		"--tx-file", signedTxPath,
		"--testnet-magic", fmt.Sprintf("%d", CardanoTestNetMagic),
	); err != nil {
		return nil, fmt.Errorf("submit: %v", err)
	}
	fmt.Printf("Casted vote %s -> %s\n", req.Voter, req.Candidate)
	return &api.CastVoteResponse{}, nil
}

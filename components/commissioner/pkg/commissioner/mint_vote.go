package commissioner

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/thavlik/bvs/components/commissioner/pkg/api"
)

var minUTxOValue = 1000000

func (s *Server) MintVote(ctx context.Context, req api.MintVoteRequest) (*api.MintVoteResponse, error) {
	// Look up the election (policy) signing key
	election, err := s.storage.RetrieveElection(req.Election)
	if err != nil {
		return nil, fmt.Errorf("storage: %v", err)
	}
	policySigningKey := election.SigningKey
	policyVerificationKey := election.VerificationKey

	// Look up minter signing key
	minter, err := s.storage.RetrieveMinter(req.Minter)
	if err != nil {
		return nil, fmt.Errorf("storage: %v", err)
	}
	paymentSigningKey := minter.SigningKey

	// Create temporary directory for interacting with cardano-cli
	id := uuid.New().String()
	rootDir := filepath.Join("/tmp", id)
	if err := os.MkdirAll(rootDir, 0777); err != nil {
		return nil, err
	}
	defer func() {
		if err := os.RemoveAll(rootDir); err != nil {
			fmt.Printf("failed to clean up temp dir: %v\n", err)
		}
	}()

	// Write key files
	paymentSigningKeyPath := filepath.Join(rootDir, "payment.skey")
	if err := ioutil.WriteFile(
		paymentSigningKeyPath,
		[]byte(paymentSigningKey),
		0644,
	); err != nil {
		return nil, err
	}
	policyVerificationKeyPath := filepath.Join(rootDir, "policy.vkey")
	if err := ioutil.WriteFile(
		policyVerificationKeyPath,
		[]byte(policyVerificationKey),
		0644,
	); err != nil {
		return nil, err
	}
	policySigningKeyPath := filepath.Join(rootDir, "policy.skey")
	if err := ioutil.WriteFile(
		policySigningKeyPath,
		[]byte(policySigningKey),
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
	policyID := election.PolicyID
	minterAddress := minter.Address
	voter := req.Voter
	tokenAmount := 1
	mintingScript := election.MintingScript
	invalidHereafter := election.InvalidHereafter
	mintingScriptPath := filepath.Join(rootDir, "policy.script")
	if err := ioutil.WriteFile(
		mintingScriptPath,
		[]byte(mintingScript),
		0644,
	); err != nil {
		return nil, err
	}

	// Get info about the minter's last transaction
	utxos, err := queryAddress(minterAddress)
	if err != nil {
		return nil, fmt.Errorf("queryAddress: %v", err)
	}
	minterInfo := utxos[0]
	txHash := minterInfo.txHash
	txIx := minterInfo.txIx
	output := minterInfo.lovelace
	timestamp := time.Now()
	metadata := metadataJson(id, policyID, timestamp.UnixNano())
	metadataJsonPath := filepath.Join(rootDir, "metadata.json")
	if err := ioutil.WriteFile(
		metadataJsonPath,
		[]byte(metadata),
		0644,
	); err != nil {
		return nil, err
	}

	// Build the transaction without specifying a fee
	rawTxPath := filepath.Join(rootDir, "matx.raw")
	giftAmount := 6 * minUTxOValue
	output -= giftAmount
	if _, err := Exec(
		"cardano-cli", "transaction", "build-raw",
		"--fee", "0",
		"--tx-in", fmt.Sprintf("%s#%d", txHash, txIx),
		"--tx-out", fmt.Sprintf(`%s + %d`, minterAddress, output),
		"--tx-out", fmt.Sprintf(`%s + %d + %d %s.%s`, voter, giftAmount, tokenAmount, policyID, tokenName),
		"--mint", fmt.Sprintf("%d %s.%s", tokenAmount, policyID, tokenName),
		"--minting-script-file", mintingScriptPath,
		"--metadata-json-file", metadataJsonPath,
		"--invalid-hereafter", fmt.Sprintf("%d", invalidHereafter),
		"--out-file", rawTxPath,
	); err != nil {
		return nil, fmt.Errorf("draft tx: %v", err)
	}

	// Calculate the fee
	fee, err := calculateFee(rawTxPath, protocolJsonPath)
	if err != nil {
		return nil, fmt.Errorf("calculateFee: %v", err)
	}

	// Gift the minimum amount plus a healthy padding to pay fee
	padding := fee * 2
	giftAmount += padding
	output -= padding

	// Sanity check: make sure we have enough money left over
	if output-fee < minUTxOValue {
		return nil, fmt.Errorf("minter has insufficient funds (has %d, needed >=%d)", output, fee+minUTxOValue)
	}

	// Subtract the fee from the output ADA balance
	// Unspent assets are offered as network fee
	output -= fee

	// Rebuild the transaction with the calculated fee
	if _, err := Exec(
		"cardano-cli", "transaction", "build-raw",
		"--fee", fmt.Sprintf("%d", fee),
		"--tx-in", fmt.Sprintf("%s#%d", txHash, txIx),
		"--tx-out", fmt.Sprintf(`%s + %d`, minterAddress, output),
		"--tx-out", fmt.Sprintf(`%s + %d + %d %s.%s`, voter, giftAmount, tokenAmount, policyID, tokenName),
		"--mint", fmt.Sprintf("%d %s.%s", tokenAmount, policyID, tokenName),
		"--minting-script-file", mintingScriptPath,
		"--metadata-json-file", metadataJsonPath,
		"--invalid-hereafter", fmt.Sprintf("%d", invalidHereafter),
		"--out-file", rawTxPath,
	); err != nil {
		return nil, fmt.Errorf("rebuild tx: %v", err)
	}

	// Sign the transaction
	signedTxPath := filepath.Join(rootDir, "matx.signed")
	if _, err := Exec(
		"cardano-cli", "transaction", "sign",
		"--signing-key-file", paymentSigningKeyPath,
		"--signing-key-file", policySigningKeyPath,
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

	fmt.Printf("Minted vote %s to %s\n", id, voter)

	// TODO: get minted asset ID
	return &api.MintVoteResponse{
		ID:    id,
		Asset: "",
	}, nil
}

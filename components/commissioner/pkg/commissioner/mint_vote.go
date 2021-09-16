package commissioner

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/thavlik/bvs/components/commissioner/pkg/api"
)

func metadataJson(policyID string) string {
	return fmt.Sprintf(`{
	"721": {
		"%s": {
			"Vote": {
				"description": "This is my first NFT thanks to the Cardano foundation",
				"name": "Cardano foundation NFT guide token",
				"id": 1
			}
		}
	}
}`, policyID)
}

func generateMintingScript(
	invalidHereafter int,
	policyVerificationKeyPath string,
) (string, error) {
	cmd := exec.Command(
		"cardano-cli", "address", "key-hash",
		"--payment-verification-key-file", policyVerificationKeyPath,
	)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}
	if err := cmd.Start(); err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(stdout)
	if err != nil {
		return "", err
	}
	if err := cmd.Wait(); err != nil {
		return "", fmt.Errorf("cardano-cli: %v", err)
	}
	keyHash := strings.TrimSpace(string(body))
	return fmt.Sprintf(
		`{
	"type": "all",
	"scripts": [{
		"type": "before",
		"slot": %d,
	}, {
		"type": "sig",
		"keyHash": "%s"
	}]
}`,
		invalidHereafter,
		keyHash,
	), nil
}

func getCurrentSlot() (int, error) {
	tip, err := queryTip()
	if err != nil {
		return 0, fmt.Errorf("queryTip: %v", err)
	}
	return tip.Slot, nil
}

var errNoTransactions = fmt.Errorf("no previous transactions")

type addressInfo struct {
	txHash   string
	txIx     int
	lovelace int
}

func queryAddress(address string) (*addressInfo, error) {
	cmd := exec.Command(
		"bash",
		"-c",
		fmt.Sprintf(`cardano-cli query utxo --address %s --testnet-magic %d`, address, CardanoTestNetMagic),
	)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(stdout)
	if err != nil {
		return nil, err
	}
	if err := cmd.Wait(); err != nil {
		return nil, fmt.Errorf("cardano-cli: %v", err)
	}
	lines := strings.Split(string(body), "\n")
	if len(lines) < 3 {
		return nil, errNoTransactions
	}
	parts := strings.Split(lines[2], " ")
	var filtered []string
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if len(part) > 0 {
			filtered = append(filtered, part)
		}
	}
	txHash := filtered[0]
	txIx, err := strconv.ParseInt(filtered[1], 10, 64)
	if err != nil {
		return nil, err
	}
	lovelace, err := strconv.ParseInt(filtered[2], 10, 64)
	if err != nil {
		return nil, err
	}
	// If the unit is displayed, make sure it's lovelace, otherwise just assume
	if len(filtered) > 3 && filtered[3] != "lovelace" {
		return nil, fmt.Errorf("error querying address, expected 'lovelace', got '%s'", filtered[3])
	}
	return &addressInfo{
		txHash:   txHash,
		txIx:     int(txIx),
		lovelace: int(lovelace),
	}, nil
}

func (s *Server) MintVote(ctx context.Context, req api.MintVoteRequest) (*api.MintVoteResponse, error) {
	// TODO: make sure minter has enough ADA and return specific error if not

	// Look up the election (policy) signing key
	election, err := s.storage.RetrieveElection(req.Election)
	if err != nil {
		return nil, fmt.Errorf("storage: %v", err)
	}
	policySigningKey := election.SigningKey
	policyVerificationKey := election.VerificationKey

	// Look up auditor (minter) signing key
	minter, err := s.storage.RetrieveMinter(req.Auditor.Agent)
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

	// https://developers.cardano.org/docs/native-tokens/minting-nfts/
	policyID := election.PolicyID
	address := "addr_test1vqd42cm0w0v7z9rve320kecjjsl7fdy8xxlqm28mg970d3c5ss752"
	tokenName := "Vote"
	tokenAmount := 1

	mintingScript := election.MintingScript
	invalidHereafter := election.InvalidHereafter
	mintingScriptFile := filepath.Join(rootDir, "policy.script")
	if err := ioutil.WriteFile(
		mintingScriptFile,
		[]byte(mintingScript),
		0644,
	); err != nil {
		return nil, err
	}

	// Get info about the minter's last transaction
	// TODO: get payment address
	addressInfo, err := queryAddress(address)
	if err != nil {
		return nil, fmt.Errorf("queryAddress: %v", err)
	}
	txHash := addressInfo.txHash
	txIx := addressInfo.txIx
	output := addressInfo.lovelace

	metadata := metadataJson(policyID)
	metadataJsonFile := filepath.Join(rootDir, "metadata.json")
	if err := ioutil.WriteFile(
		metadata,
		[]byte(metadata),
		0644,
	); err != nil {
		return nil, err
	}

	// Build the transaction without specifying a fee
	rawTxPath := filepath.Join(rootDir, "matx.raw")
	txOut := fmt.Sprintf(`%s+%d+%d %s.%s`, address, output, tokenAmount, policyID, tokenName)
	fmt.Printf("txout: %s\n", txOut)
	cmd := exec.Command(
		"cardano-cli", "transaction", "build-raw",
		"--fee", "300000",
		"--tx-in", fmt.Sprintf("%s#%d", txHash, txIx),
		"--tx-out", txOut,
		fmt.Sprintf(`--mint="%d %s.%s"`, tokenAmount, policyID, tokenName),
		"--minting-script-file", mintingScriptFile,
		"--metadata-json-file", metadataJsonFile,
		"--invalid-hereafter", fmt.Sprintf("%d", invalidHereafter),
		"--out-file", rawTxPath,
	)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	stdoutBytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		return nil, err
	}
	stderrBytes, err := ioutil.ReadAll(stderr)
	if err != nil {
		return nil, err
	}
	if err := cmd.Wait(); err != nil {
		fmt.Println(string(stdoutBytes))
		fmt.Println(string(stderrBytes))
		return nil, fmt.Errorf("cardano-cli draft tx: %v", err)
	}

	// Calculate the fee
	fee, err := calculateFee(rawTxPath)
	if err != nil {
		return nil, fmt.Errorf("calculateFee: %v", err)
	}

	// Subtract the fee from the output ADA balance
	// Unspent assets are offered as network fee
	output -= fee

	// Rebuild the transaction with the calculated fee
	if err := (exec.Command(
		"cardano-cli", "transaction", "build-raw",
		"--fee", fmt.Sprintf("%d", fee),
		"--tx-in", fmt.Sprintf("%s#%d", txHash, txIx),
		"--tx-out", fmt.Sprintf(`%s+%d+%d %s.%s`, address, output, tokenAmount, policyID, tokenName),
		fmt.Sprintf(`--mint="%d %s.%s"`, tokenAmount, policyID, tokenName),
		"--minting-script-file", mintingScriptFile,
		"--metadata-json-file", metadataJsonFile,
		"--invalid-hereafter", fmt.Sprintf("%d", invalidHereafter),
		"--out-file", rawTxPath,
	)).Run(); err != nil {
		return nil, fmt.Errorf("cardano-cli rebuild tx: %v", err)
	}

	// Sign the transaction
	signedTxPath := filepath.Join(rootDir, "matx.signed")
	if err := (exec.Command(
		"cardano-cli", "transaction", "sign",
		"--signing-key-file", paymentSigningKeyPath,
		"--signing-key-file", policySigningKeyPath,
		"--mainnet",
		"--tx-body-file", rawTxPath,
		"--out-file", signedTxPath,
	)).Run(); err != nil {
		return nil, fmt.Errorf("cardano-cli sign: %v", err)
	}

	// Submit the transaction
	if err := (exec.Command(
		"cardano-cli", "transaction", "submit",
		"--tx-file", signedTxPath,
		"--mainnet",
	)).Run(); err != nil {
		return nil, fmt.Errorf("cardano-cli submit: %v", err)
	}

	// TODO: get minted asset ID
	return &api.MintVoteResponse{
		ID:    id,
		Asset: "",
	}, nil
}

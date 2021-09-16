package definitions

type Commissioner interface {
	CreateElection(CreateElectionRequest) CreateElectionResponse
	CreateMinter(CreateMinterRequest) CreateMinterResponse
	CreateVoter(CreateVoterRequest) CreateVoterResponse
	MintVote(MintVoteRequest) MintVoteResponse
	CastVote(CastVoteRequest) CastVoteResponse
	QueryAddress(QueryAddressRequest) QueryAddressResponse
}

type CreateElectionRequest struct {
}

type Key struct {
	Type        string `json:"type"`
	Description string `json:"description"`
	CborHex     string `json:"cborHex"`
}

type CreateMinterRequest struct {
	SigningKey      Key    `json:"signingKey"`      // contents of signing key file
	VerificationKey Key    `json:"verificationKey"` // contents of verification key file
	Address         string `json:"address"`
}

type CreateVoterRequest struct {
}

type CreateVoterResponse struct {
	SigningKey      Key    `json:"signingKey"`      // contents of signing key file
	VerificationKey Key    `json:"verificationKey"` // contents of verification key file
	Address         string `json:"address"`         // payment address
}

type CreateMinterResponse struct {
	ID string `json:"id"`
}

type CreateElectionResponse struct {
	ID              string `json:"id"`              // generated UUID
	PolicyID        string `json:"policyID"`        // https://developers.cardano.org/docs/native-tokens/minting-nfts/
	VerificationKey Key    `json:"verificationKey"` // contents of verification key file
}

type MintVoteRequest struct {
	Election string `json:"election"` // Corresponds to CreateElectionResponse.ID
	Voter    string `json:"voter"`    // Public address of the voter that will receive the NFT
	Minter   string `json:"minter"`   // Minter UUID
}

type MintVoteResponse struct {
	ID    string `json:"id"`    // Request UUID
	Asset string `json:"asset"` // Address of the asset; Concatenation of the policy_id and hex-encoded asset_name; https://docs.blockfrost.io/#tag/Cardano-Assets/paths/~1assets~1{asset}/get
}

type CastVoteRequest struct {
	Election   string `json:"election"`   // Corresponds to CreateElectionResponse.ID
	Voter      string `json:"voter"`      // Input wallet address
	SigningKey Key    `json:"signingKey"` // contents of signing key file for input wallet
	Candidate  string `json:"candidate"`  // Candidate payment address
}

type CastVoteResponse struct {
	ID string `json:"id"` // Request UUID
}

type QueryAddressRequest struct {
	Address string `json:"address"`
}

type UnspentTransaction struct {
	TxHash   string `json:"txHash"`
	TxIx     int    `json:"txIx"`
	Lovelace int    `json:"lovelace"`
	Balance  string `json:"balance"`
}

type QueryAddressResponse struct {
	UnspentTransactions []*UnspentTransaction `json:"unspentTransactions"`
}

package definitions

type Commissioner interface {
	CreateElection(CreateElectionRequest) CreateElectionResponse
	MintVote(MintVoteRequest) MintVoteResponse
}

type CreateElectionRequest struct {
	Name     string `json:"name"`
	Deadline int64  `json:"deadline"`
}

type VerificationKey struct {
	Type        string `json:"type"`
	Description string `json:"description"`
	CborHex     string `json:"cborHex"`
}

type CreateElectionResponse struct {
	ID              string          `json:"id"`              // policyID; https://developers.cardano.org/docs/native-tokens/minting-nfts/
	VerificationKey VerificationKey `json:"verificationKey"` // contents of verification key file
}

type Auditor struct {
	Agent     string `json:"agent"`     // Public address of the minting poll worker
	Timestamp int64  `json:"timestamp"` // 64-bit nanosecond count since unix epoch
	Proof     string `json:"proof"`     // Signature of the vote from the auditor's corresponding private key
}

type MintVoteRequest struct {
	Election string  `json:"election"` // Corresponds to CreateElectionRequest.Name
	Voter    string  `json:"voter"`    // Public address of the voter that will receive the NFT
	Ident    string  `json:"ident"`    // Personally identifying fingerprint
	Auditor  Auditor `json:"auditor"`  // Info about poll worker who did ID check
}

type MintVoteResponse struct {
	ID    string `json:"id"`    // Request UUID
	Asset string `json:"asset"` // Address of the asset; Concatenation of the policy_id and hex-encoded asset_name; https://docs.blockfrost.io/#tag/Cardano-Assets/paths/~1assets~1{asset}/get
}

package definitions

type Node interface {
	ProbeReady(ProbeReadyRequest) ProbeReadyResponse
}

type ProbeReadyRequest struct {
}

type ProbeReadyResponse struct {
	IsReady bool `json:"isReady"`
}

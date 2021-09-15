// Code generated by oto; DO NOT EDIT.

package api

import (
	"context"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/pacedotdev/oto/otohttp"
)

var (
	commissionerCreateElectionTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "commissioner_create_election_total",
		Help: "Auto-generated metric incremented on every call to Commissioner.CreateElection",
	})
	commissionerCreateElectionSuccessTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "commissioner_create_election_success_total",
		Help: "Auto-generated metric incremented on every call to Commissioner.CreateElection that does not return with an error",
	})

	commissionerMintVoteTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "commissioner_mint_vote_total",
		Help: "Auto-generated metric incremented on every call to Commissioner.MintVote",
	})
	commissionerMintVoteSuccessTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "commissioner_mint_vote_success_total",
		Help: "Auto-generated metric incremented on every call to Commissioner.MintVote that does not return with an error",
	})
)

type Commissioner interface {
	CreateElection(context.Context, CreateElectionRequest) (*CreateElectionResponse, error)
	MintVote(context.Context, MintVoteRequest) (*MintVoteResponse, error)
}

type commissionerServer struct {
	server       *otohttp.Server
	commissioner Commissioner
}

func RegisterCommissioner(server *otohttp.Server, commissioner Commissioner) {
	handler := &commissionerServer{
		server:       server,
		commissioner: commissioner,
	}
	server.Register("Commissioner", "CreateElection", handler.handleCreateElection)
	server.Register("Commissioner", "MintVote", handler.handleMintVote)
}

func (s *commissionerServer) handleCreateElection(w http.ResponseWriter, r *http.Request) {
	commissionerCreateElectionTotal.Inc()
	var request CreateElectionRequest
	if err := otohttp.Decode(r, &request); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	response, err := s.commissioner.CreateElection(r.Context(), request)
	if err != nil {
		log.Println("TODO: oto service error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := otohttp.Encode(w, r, http.StatusOK, response); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	commissionerCreateElectionSuccessTotal.Inc()
}

func (s *commissionerServer) handleMintVote(w http.ResponseWriter, r *http.Request) {
	commissionerMintVoteTotal.Inc()
	var request MintVoteRequest
	if err := otohttp.Decode(r, &request); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	response, err := s.commissioner.MintVote(r.Context(), request)
	if err != nil {
		log.Println("TODO: oto service error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := otohttp.Encode(w, r, http.StatusOK, response); err != nil {
		s.server.OnErr(w, r, err)
		return
	}
	commissionerMintVoteSuccessTotal.Inc()
}

type Auditor struct {
	Agent     string `json:"agent"`
	Timestamp int64  `json:"timestamp"`
	Proof     string `json:"proof"`
}

type CreateElectionRequest struct {
	Name     string `json:"name"`
	Deadline int64  `json:"deadline"`
}

type CreateElectionResponse struct {
	ID              string `json:"id"`
	VerificationKey string `json:"verificationKey"`
	Error           string `json:"error,omitempty"`
}

type MintVoteRequest struct {
	Election string  `json:"election"`
	Voter    string  `json:"voter"`
	Ident    string  `json:"ident"`
	Auditor  Auditor `json:"auditor"`
}

type MintVoteResponse struct {
	ID    string `json:"id"`
	Asset string `json:"asset"`
	Error string `json:"error,omitempty"`
}

package gateway

import (
	"fmt"
	"net/http"
)

func (s *Server) listen(port int) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/elections/get", s.handleGetElection())
	mux.HandleFunc("/elections/list", s.handleListElections())
	mux.HandleFunc("/elections/create", s.handleCreateElection())
	mux.HandleFunc("/elections/delete", s.handleDeleteElection())
	mux.HandleFunc("/election/candidates", s.handleSetElectionCandidates())
	mux.HandleFunc("/minter/get", s.handleGetMinter())
	mux.HandleFunc("/minter/list", s.handleListMinters())
	mux.HandleFunc("/minter/create", s.handleCreateMinter())
	mux.HandleFunc("/minter/delete", s.handleDeleteMinter())
	mux.HandleFunc("/vote/get", s.handleGetVote())
	mux.HandleFunc("/vote/list", s.handleListVotes())
	mux.HandleFunc("/vote/mint", s.handleMintVote())
	mux.HandleFunc("/vote/cast", s.handleCastVote())
	mux.HandleFunc("/vote/void", s.handleVoidVote())
	return http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), mux)
}

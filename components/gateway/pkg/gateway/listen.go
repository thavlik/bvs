package gateway

import (
	"fmt"
	"net/http"
)

func (s *Server) listen(port int) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/vote/get", s.handleGetVote())
	mux.HandleFunc("/vote/list", s.handleListVotes())
	mux.HandleFunc("/vote/mint", s.handleMintVote())
	mux.HandleFunc("/vote/cast", s.handleCastVote())
	return http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), mux)
}

func (s *Server) handleGetVote() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
	}
}

func (s *Server) handleListVotes() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
	}
}

func (s *Server) handleMintVote() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
	}
}

func (s *Server) handleCastVote() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
	}
}

package gateway

import (
	"net/http"
)

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

func (s *Server) handleVoidVote() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
	}
}

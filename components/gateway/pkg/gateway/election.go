package gateway

import (
	"net/http"
)

func (s *Server) handleGetElection() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
	}
}

func (s *Server) handleListElections() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
	}
}

func (s *Server) handleCreateElection() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
	}
}

func (s *Server) handleDeleteElection() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
	}
}

func (s *Server) handleSetElectionCandidates() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
	}
}

package kfca

import (
	"encoding/json"
	"log"
	"net/http"
)

type ReadyzResult struct {
	Ready bool `json:"ready"`
}

type readyzsrvc struct {
	logger *log.Logger
}

func NewReadyzHandler(logger *log.Logger) http.HandlerFunc {
	s := &readyzsrvc{logger: logger}
	return s.ServeHTTP
}

func (s *readyzsrvc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	res := &ReadyzResult{Ready: true}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

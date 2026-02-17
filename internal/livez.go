package kfca

import (
	"encoding/json"
	"log"
	"net/http"
)

type LivezResult struct {
	Live bool `json:"live"`
}

type livezsrvc struct {
	logger *log.Logger
}

func NewLivezHandler(logger *log.Logger) http.HandlerFunc {
	s := &livezsrvc{logger: logger}
	return s.ServeHTTP
}

func (s *livezsrvc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	res := &LivezResult{Live: true}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

package kfca

import (
	"context"
	livez "k8s-federated-credential-api/gen/livez"
	"log"
)

// livez service example implementation.
// The example methods log the requests and return zero values.
type livezsrvc struct {
	logger *log.Logger
}

// NewLivez returns the livez service implementation.
func NewLivez(logger *log.Logger) livez.Service {
	return &livezsrvc{logger}
}

// Livez implements livez.
func (s *livezsrvc) Livez(ctx context.Context) (res *livez.LivezResult, err error) {
	res = &livez.LivezResult{}
	res.Live = true
	return res, nil
}

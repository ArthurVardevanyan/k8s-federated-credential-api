package kfca

import (
	"context"
	readyz "k8s-federated-credential-api/gen/readyz"
	"log"
)

// readyz service example implementation.
// The example methods log the requests and return zero values.
type readyzsrvc struct {
	logger *log.Logger
}

// NewReadyz returns the readyz service implementation.
func NewReadyz(logger *log.Logger) readyz.Service {
	return &readyzsrvc{logger}
}

// Readyz implements readyz.
func (s *readyzsrvc) Readyz(ctx context.Context) (res *readyz.ReadyzResult, err error) {
	res = &readyz.ReadyzResult{}
	res.Ready = true

	return res, nil
}

package kfcc

import (
	"context"
	tokenexchange "kubernetes-federated-credential-controller/gen/token_exchange"
	"log"
)

// tokenExchange service example implementation.
// The example methods log the requests and return zero values.
type tokenExchangesrvc struct {
	logger *log.Logger
}

// NewTokenExchange returns the tokenExchange service implementation.
func NewTokenExchange(logger *log.Logger) tokenexchange.Service {
	return &tokenExchangesrvc{logger}
}

// ExchangeToken implements exchangeToken.
func (s *tokenExchangesrvc) ExchangeToken(ctx context.Context, p *tokenexchange.ExchangeTokenPayload) (res string, err error) {
	s.logger.Print("tokenExchange.exchangeToken")
	return
}

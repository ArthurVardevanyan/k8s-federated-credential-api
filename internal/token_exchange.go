package kfcc

import (
	"context"
	tokenexchange "kubernetes-federated-credential-controller/gen/token_exchange"
	"log"

	"github.com/coreos/go-oidc/v3/oidc"
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

	provider, err := oidc.NewProvider(context.TODO(), "")
	if err != nil {
		log.Fatal(err)
	}

	oidcConfig := &oidc.Config{
		SkipClientIDCheck: true,
	}
	verifier := provider.Verifier(oidcConfig)

	_, err = verifier.Verify(context.TODO(), *p.JWT)
	if err != nil {
		return "Failed to parse the JWT.\nError: %s", err
	}
	return "The token is valid.", nil

}

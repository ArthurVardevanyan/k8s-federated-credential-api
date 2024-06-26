// Code generated by goa v3.16.1, DO NOT EDIT.
//
// tokenExchange service
//
// Command:
// $ goa gen k8s-federated-credential-api/design

package tokenexchange

import (
	"context"

	goa "goa.design/goa/v3/pkg"
	"goa.design/goa/v3/security"
)

// Exchange an incoming Kubernetes Token for Another Kubernetes Token
type Service interface {
	// ExchangeToken implements exchangeToken.
	ExchangeToken(context.Context, *ExchangeTokenPayload) (res *StatusResult, err error)
}

// Auther defines the authorization functions to be implemented by the service.
type Auther interface {
	// JWTAuth implements the authorization logic for the JWT security scheme.
	JWTAuth(ctx context.Context, token string, schema *security.JWTScheme) (context.Context, error)
}

// APIName is the name of the API as defined in the design.
const APIName = "kfca"

// APIVersion is the version of the API as defined in the design.
const APIVersion = "0.0.1"

// ServiceName is the name of the service as defined in the design. This is the
// same value that is set in the endpoint request contexts under the ServiceKey
// key.
const ServiceName = "tokenExchange"

// MethodNames lists the service method names as defined in the design. These
// are the same values that are set in the endpoint request contexts under the
// MethodKey key.
var MethodNames = [1]string{"exchangeToken"}

// ExchangeTokenPayload is the payload type of the tokenExchange service
// exchangeToken method.
type ExchangeTokenPayload struct {
	// The JWT Token from the impersonating service account
	Authorization string
	// The target namespace for impersonation
	Namespace string
	// The target serviceAccount
	ServiceAccountName string
}

// Status with a token.
type Status struct {
	// The status token
	Token string
}

// StatusResult is the result type of the tokenExchange service exchangeToken
// method.
type StatusResult struct {
	// The status information with a token
	Status *Status
}

// MakeInternalError builds a goa.ServiceError from an error.
func MakeInternalError(err error) *goa.ServiceError {
	return goa.NewServiceError(err, "internal_error", false, false, false)
}

// MakeForbidden builds a goa.ServiceError from an error.
func MakeForbidden(err error) *goa.ServiceError {
	return goa.NewServiceError(err, "forbidden", false, false, false)
}

// MakeNotFound builds a goa.ServiceError from an error.
func MakeNotFound(err error) *goa.ServiceError {
	return goa.NewServiceError(err, "not_found", false, false, false)
}

// MakeNotAcceptable builds a goa.ServiceError from an error.
func MakeNotAcceptable(err error) *goa.ServiceError {
	return goa.NewServiceError(err, "not_acceptable", false, false, false)
}

// MakeUnauthorized builds a goa.ServiceError from an error.
func MakeUnauthorized(err error) *goa.ServiceError {
	return goa.NewServiceError(err, "unauthorized", false, false, false)
}

// MakeBadRequestError builds a goa.ServiceError from an error.
func MakeBadRequestError(err error) *goa.ServiceError {
	return goa.NewServiceError(err, "bad_request_error", false, false, false)
}

// MakeUnsupportedMediaType builds a goa.ServiceError from an error.
func MakeUnsupportedMediaType(err error) *goa.ServiceError {
	return goa.NewServiceError(err, "unsupported_media_type", false, false, false)
}

// MakeTooManyRequests builds a goa.ServiceError from an error.
func MakeTooManyRequests(err error) *goa.ServiceError {
	return goa.NewServiceError(err, "too_many_requests", false, false, false)
}

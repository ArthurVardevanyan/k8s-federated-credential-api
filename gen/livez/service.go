// Code generated by goa v3.16.1, DO NOT EDIT.
//
// livez service
//
// Command:
// $ goa gen k8s-federated-credential-api/design

package livez

import (
	"context"
)

// Service is the livez service interface.
type Service interface {
	// Livez implements livez.
	Livez(context.Context) (res *LivezResult, err error)
}

// APIName is the name of the API as defined in the design.
const APIName = "kfca"

// APIVersion is the version of the API as defined in the design.
const APIVersion = "0.0.1"

// ServiceName is the name of the service as defined in the design. This is the
// same value that is set in the endpoint request contexts under the ServiceKey
// key.
const ServiceName = "livez"

// MethodNames lists the service method names as defined in the design. These
// are the same values that are set in the endpoint request contexts under the
// MethodKey key.
var MethodNames = [1]string{"livez"}

// LivezResult is the result type of the livez service livez method.
type LivezResult struct {
	Live bool
}
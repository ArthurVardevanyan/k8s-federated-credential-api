// Code generated by goa v3.16.1, DO NOT EDIT.
//
// livez HTTP client types
//
// Command:
// $ goa gen k8s-federated-credential-api/design

package client

import (
	livez "k8s-federated-credential-api/gen/livez"

	goa "goa.design/goa/v3/pkg"
)

// LivezResponseBody is the type of the "livez" service "livez" endpoint HTTP
// response body.
type LivezResponseBody struct {
	Live *bool `form:"live,omitempty" json:"live,omitempty" xml:"live,omitempty"`
}

// NewLivezResultOK builds a "livez" service "livez" endpoint result from a
// HTTP "OK" response.
func NewLivezResultOK(body *LivezResponseBody) *livez.LivezResult {
	v := &livez.LivezResult{
		Live: *body.Live,
	}

	return v
}

// ValidateLivezResponseBody runs the validations defined on LivezResponseBody
func ValidateLivezResponseBody(body *LivezResponseBody) (err error) {
	if body.Live == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("live", "body"))
	}
	return
}

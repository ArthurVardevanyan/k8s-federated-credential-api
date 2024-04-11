package design

import (
	. "goa.design/goa/v3/dsl"
)

// Service describes a service
var _ = Service("tokenExchange", func() {
	Description("Exchange an incoming Kubernetes Token for Another Kubernetes Token")
	// Method describes a service method (endpoint)
	Method("exchangeToken", func() {
		// Payload describes the method payload
		// Here the payload is an object that consists of two fields
		Payload(func() {

			// Attribute describes an object field
			Field(0, "jwt", String, "The JWT Token from the impersonating service account")
			Field(1, "namespace", String, "The target namespace for impersonation")
			Field(2, "serviceAccount", String, "The target serviceAccount")
			//Required("jwt", "namespace", "serviceAccount" )
		})
		// Result describes the method result
		// Here the result is a simple integer value
		Result(String)
		// HTTP describes the HTTP transport mapping
		HTTP(func() {
			// Requests to the service consist of HTTP GET requests
			// The payload fields are encoded as path parameters
			POST("/exchangeToken")
			// Responses use a "200 OK" HTTP status
			// The result is encoded in the response body
			Response(StatusOK)
		})
	})
})

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
			Field(2, "serviceAccountName", String, "The target serviceAccount")
			//Required("jwt", "namespace", "serviceAccount" )
		})
		// Result describes the method result
		// Here the result is a simple integer value
		Result(StatusResult)
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

// Define the nested attribute for the status
var Status = Type("Status", func() {
	Description("Status with a token.")
	Attribute("token", String, "The status token")
	Required("token") // The token attribute is required
})

// Define the result type for the endpoint, including the nested 'status' attribute
var StatusResult = Type("StatusResult", func() {
	Description("The result type containing status information.")
	Attribute("status", Status, "The status information with a token") // Nested attribute
	Required("status")                                                 // The status attribute is required
})

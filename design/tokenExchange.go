package design

import (
	. "goa.design/goa/v3/dsl"
)

// Service describes a service
var _ = Service("tokenExchange", func() {
	Description("Exchange an incoming Kubernetes Token for Another Kubernetes Token")
	// Method describes a service method (endpoint)
	Method("exchangeToken", func() {
		Security(JWTAuth)
		// Payload describes the method payload
		// Here the payload is an object that consists of two fields
		Payload(func() {
			Token("Authorization", String, "The JWT Token from the impersonating service account", func() {
				Pattern("^Bearer [^ ]+") //Returning 400 not 401/403
				MaxLength(2000)
			})

			// Attribute describes an object field
			Field(0, "namespace", String, "The target namespace for impersonation", func() {
				Pattern(`[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*`)
				MaxLength(253)
			})
			Field(1, "serviceAccountName", String, "The target serviceAccount", func() {
				Pattern(`[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*`)
				MaxLength(253)
			})

			Required("Authorization") //Returning 400 not 401/403
			Required("namespace")
			Required("serviceAccountName")

		})

		// Result describes the method result
		// Here the result is a simple integer value
		Result(StatusResult)
		Error("internal_error", ErrorResult, "Internal Server Error")
		Error("forbidden", ErrorResult, "Forbidden")
		Error("not_found", ErrorResult, "Not Found")
		Error("not_acceptable", ErrorResult, "Not Acceptable")
		Error("unauthorized", ErrorResult, "Unauthorized")
		Error("bad_request_error", ErrorResult, "Bad Request")
		Error("unsupported_media_type", ErrorResult, "Unsupported Media Type")
		Error("too_many_requests", ErrorResult, "Too Many Requests")

		// HTTP describes the HTTP transport mapping
		HTTP(func() {
			// Requests to the service consist of HTTP GET requests
			// The payload fields are encoded as path parameters
			POST("/exchangeToken")
			// Responses use a "200 OK" HTTP status
			// The result is encoded in the response body

			Response(StatusOK)
			Response("internal_error", StatusInternalServerError, func() {
				ContentType("application/json")
			})
			Response("forbidden", StatusForbidden, func() {
				ContentType("application/json")
			})
			Response("not_found", StatusNotFound, func() {
				ContentType("application/json")
			})
			Response("not_acceptable", StatusNotAcceptable, func() {
				ContentType("application/json")
			})
			Response("unauthorized", StatusUnauthorized, func() {
				ContentType("application/json")
			})
			Response("bad_request_error", StatusBadRequest, func() {
				ContentType("application/json")
			})
			Response("unsupported_media_type", StatusUnsupportedMediaType, func() {
				ContentType("application/json")
			})
			Response("too_many_requests", StatusTooManyRequests, func() {
				ContentType("application/json")
			})
		})
	})
})

// Define the nested attribute for the status
var Status = Type("Status", func() {
	Description("Status with a token.")
	Attribute("token", String, "The status token", func() {
		Pattern(`^(?:[A-Za-z0-9+/]{4})*(?:[A-Za-z0-9+/]{2}==|[A-Za-z0-9+/]{3}=)?$`)
		MaxLength(2000)
	})
	Required("token") // The token attribute is required
})

// Define the result type for the endpoint, including the nested 'status' attribute
var StatusResult = Type("StatusResult", func() {
	Description("The result type containing status information.")
	Attribute("status", Status, "The status information with a token") // Nested attribute
	Required("status")                                                 // The status attribute is required
})

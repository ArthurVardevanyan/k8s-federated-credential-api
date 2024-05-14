package design

import (
	. "goa.design/goa/v3/dsl"
)

var _ = Service("readyz", func() {
	Method("readyz", func() {
		NoSecurity()
		Meta("openapi:generate", "false")
		Result(func() {
			Attribute("ready", Boolean)
			Required("ready")
		})
		// HTTP describes the HTTP transport mapping
		HTTP(func() {
			// Requests to the service consist of HTTP GET requests
			// The payload fields are encoded as path parameters
			GET("/readyz")
			// Responses use a "200 OK" HTTP status
			// The result is encoded in the response body
			Response(StatusOK)
		})
	})
})

package design

import (
	. "goa.design/goa/v3/dsl"
)

var _ = Service("livez", func() {
	Method("livez", func() {
		NoSecurity()
		Meta("openapi:generate", "false")
		Result(func() {
			Attribute("live", Boolean)
			Required("live")
		})
		// HTTP describes the HTTP transport mapping
		HTTP(func() {
			// Requests to the service consist of HTTP GET requests
			// The payload fields are encoded as path parameters
			GET("/livez")
			// Responses use a "200 OK" HTTP status
			// The result is encoded in the response body
			Response(StatusOK)
		})
	})
})

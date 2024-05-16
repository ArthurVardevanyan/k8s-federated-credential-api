package design

import (
	. "goa.design/goa/v3/dsl"
)

var JWTAuth = JWTSecurity("jwt", func() {
	Description("Use JWT to authenticate requests.")
})

// API describes the global properties of the API server.
var _ = API("kfca", func() {
	Title("")
	Description("The Kubernetes Federated Credential Api")
	Server("kfca", func() {
		Services("tokenExchange")
		Services("readyz")
		Services("livez")
		Host("localhost", func() { URI("http://0.0.0.0:8088") })
	})

	HTTP(func() {
		Produces("application/json")
		Consumes("application/json")
	})
})

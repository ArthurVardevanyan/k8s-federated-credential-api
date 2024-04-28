package design

import . "goa.design/goa/v3/dsl"

// API describes the global properties of the API server.
var _ = API("kfca", func() {
	Title("")
	Description("The Kubernetes Federated Credential Api")
	Server("kfca", func() {
		Services("tokenExchange")
		Host("localhost", func() { URI("http://0.0.0.0:8088") })
	})
})

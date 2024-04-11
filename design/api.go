package design

import . "goa.design/goa/v3/dsl"

// API describes the global properties of the API server.
var _ = API("kfcc", func() {
	Title("")
	Description("The Kubernetes Federated Credential Controller")
	Server("kfcc", func() {
		Services("tokenExchange")
		Host("localhost", func() { URI("http://localhost:8088") })
	})
})

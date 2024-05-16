package main

import (
	"context"
	livezsvr "k8s-federated-credential-api/gen/http/livez/server"
	readyzsvr "k8s-federated-credential-api/gen/http/readyz/server"
	tokenexchangesvr "k8s-federated-credential-api/gen/http/token_exchange/server"
	livez "k8s-federated-credential-api/gen/livez"
	readyz "k8s-federated-credential-api/gen/readyz"
	tokenexchange "k8s-federated-credential-api/gen/token_exchange"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	goahttp "goa.design/goa/v3/http"
	httpmdlwr "goa.design/goa/v3/http/middleware"
	"goa.design/goa/v3/middleware"
)

// handleHTTPServer starts configures and starts a HTTP server on the given
// URL. It shuts down the server if any error is received in the error channel.
func handleHTTPServer(ctx context.Context, u *url.URL, tokenExchangeEndpoints *tokenexchange.Endpoints, readyzEndpoints *readyz.Endpoints, livezEndpoints *livez.Endpoints, wg *sync.WaitGroup, errc chan error, logger *log.Logger, debug bool) {

	// Setup goa log adapter.
	var (
		adapter middleware.Logger
	)
	{
		adapter = middleware.NewLogger(logger)
	}

	// Provide the transport specific request decoder and response encoder.
	// The goa http package has built-in support for JSON, XML and gob.
	// Other encodings can be used by providing the corresponding functions,
	// see goa.design/implement/encoding.
	var (
		dec = goahttp.RequestDecoder
		enc = goahttp.ResponseEncoder
	)

	// Build the service HTTP request multiplexer and configure it to serve
	// HTTP requests to the service endpoints.
	var mux goahttp.Muxer
	{
		mux = goahttp.NewMuxer()
	}

	// Wrap the endpoints with the transport specific layers. The generated
	// server packages contains code generated from the design which maps
	// the service input and output data structures to HTTP requests and
	// responses.
	var (
		tokenExchangeServer *tokenexchangesvr.Server
		readyzServer        *readyzsvr.Server
		livezServer         *livezsvr.Server
	)
	{
		eh := errorHandler(logger)
		tokenExchangeServer = tokenexchangesvr.New(tokenExchangeEndpoints, mux, dec, enc, eh, nil)
		readyzServer = readyzsvr.New(readyzEndpoints, mux, dec, enc, eh, nil)
		livezServer = livezsvr.New(livezEndpoints, mux, dec, enc, eh, nil)
		if debug {
			servers := goahttp.Servers{
				tokenExchangeServer,
				readyzServer,
				livezServer,
			}
			servers.Use(httpmdlwr.Debug(mux, os.Stdout))
		}
	}
	// Configure the mux.
	tokenexchangesvr.Mount(mux, tokenExchangeServer)
	readyzsvr.Mount(mux, readyzServer)
	livezsvr.Mount(mux, livezServer)

	// Wrap the multiplexer with additional middlewares. Middlewares mounted
	// here apply to all the service endpoints.
	var handler http.Handler = mux
	{
		handler = httpmdlwr.Log(adapter)(handler)
		handler = httpmdlwr.RequestID()(handler)
		handler = contentTypeMiddleware(handler)
	}

	// Start HTTP server using default configuration, change the code to
	// configure the server as required by your service.
	srv := &http.Server{Addr: u.Host, Handler: handler, ReadHeaderTimeout: time.Second * 60}
	for _, m := range tokenExchangeServer.Mounts {
		logger.Printf("HTTP %q mounted on %s %s", m.Method, m.Verb, m.Pattern)
	}
	for _, m := range readyzServer.Mounts {
		logger.Printf("HTTP %q mounted on %s %s", m.Method, m.Verb, m.Pattern)
	}
	for _, m := range livezServer.Mounts {
		logger.Printf("HTTP %q mounted on %s %s", m.Method, m.Verb, m.Pattern)
	}

	(*wg).Add(1)
	go func() {
		defer (*wg).Done()

		// Start HTTP server in a separate goroutine.
		go func() {
			logger.Printf("HTTP server listening on %q", u.Host)
			errc <- srv.ListenAndServe()
		}()

		<-ctx.Done()
		logger.Printf("shutting down HTTP server at %q", u.Host)

		// Shutdown gracefully with a 30s timeout.
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		err := srv.Shutdown(ctx)
		if err != nil {
			logger.Printf("failed to shutdown: %v", err)
		}
	}()
}

// errorHandler returns a function that writes and logs the given error.
// The function also writes and logs the error unique ID so that it's possible
// to correlate.
func errorHandler(logger *log.Logger) func(context.Context, http.ResponseWriter, error) {
	return func(ctx context.Context, w http.ResponseWriter, err error) {
		id := ctx.Value(middleware.RequestIDKey).(string)
		_, _ = w.Write([]byte("[" + id + "] encoding: " + err.Error()))
		logger.Printf("[%s] ERROR: %s", id, err.Error())
	}
}

// ContentTypeMiddleware returns a function that pulls the request's Content-Type header
// and stores it in the context
func contentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "Content-Type", r.Header.Get("Content-Type"))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// // ContentTypeMiddleware returns a function that pulls the request's Content-Type header
// // and stores it in the context
// type contextKey string

// func contentTypeMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		ctx := context.WithValue(r.Context(), contextKey("contentType"), r.Header.Get("Content-Type"))
// 		next.ServeHTTP(w, r.WithContext(ctx))
// 	})
// }

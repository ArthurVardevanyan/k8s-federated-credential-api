package main

import (
	"context"
	"flag"
	"fmt"
	kfca "k8s-federated-credential-api/internal"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	var (
		httpPortF = flag.String("http-port", "8088", "HTTP port to listen on")
		httpAddrF = flag.String("http-addr", "0.0.0.0", "HTTP address to listen on")
	)
	flag.Parse()

	logger := log.New(os.Stderr, "[kfca] ", log.Ltime)

	mux := http.NewServeMux()
	mux.HandleFunc("/exchangeToken", kfca.NewTokenExchangeHandler(logger))
	mux.HandleFunc("/readyz", kfca.NewReadyzHandler(logger))
	mux.HandleFunc("/livez", kfca.NewLivezHandler(logger))

	var handler http.Handler = mux
	handler = loggingMiddleware(logger)(handler)
	handler = requestIDMiddleware(handler)

	addr := net.JoinHostPort(*httpAddrF, *httpPortF)
	srv := &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: 60 * time.Second,
	}

	errc := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-c)
	}()

	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())

	wg.Add(1)
	go func() {
		defer wg.Done()

		go func() {
			logger.Printf("HTTP server listening on %q", addr)
			logger.Printf("HTTP POST /exchangeToken")
			logger.Printf("HTTP GET  /readyz")
			logger.Printf("HTTP GET  /livez")
			errc <- srv.ListenAndServe()
		}()

		<-ctx.Done()
		logger.Printf("shutting down HTTP server at %q", addr)

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer shutdownCancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			logger.Printf("failed to shutdown: %v", err)
		}
	}()

	logger.Printf("exiting (%v)", <-errc)
	cancel()
	wg.Wait()
	logger.Println("exited")
}

type contextKey string

const requestIDKey contextKey = "request-id"

func requestIDMiddleware(next http.Handler) http.Handler {
	var counter uint64
	var mu sync.Mutex
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		counter++
		id := fmt.Sprintf("%d", counter)
		mu.Unlock()
		ctx := context.WithValue(r.Context(), requestIDKey, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func loggingMiddleware(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			logger.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
		})
	}
}

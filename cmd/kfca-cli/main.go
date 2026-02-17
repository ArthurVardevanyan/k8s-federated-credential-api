package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type ExchangeTokenRequest struct {
	Namespace          string `json:"namespace"`
	ServiceAccountName string `json:"serviceAccountName"`
}

func main() {
	var (
		urlF       = flag.String("url", "http://0.0.0.0:8088", "URL to service host")
		tokenF     = flag.String("token", "", "Bearer token for authentication")
		namespaceF = flag.String("namespace", "", "Target namespace for impersonation")
		saNameF    = flag.String("service-account", "", "Target service account name")
		timeoutF   = flag.Int("timeout", 30, "Maximum number of seconds to wait for response")
	)
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `kfca-cli is a command line client for the Kubernetes Federated Credential API.

Usage:
    %s [flags]

Flags:
`, os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, `
Example:
    %s -url http://localhost:8088 -token "$(kubectl create token default --duration=1h -n default)" -namespace smoke-tests -service-account default
`, os.Args[0])
	}
	flag.Parse()

	if *tokenF == "" {
		fmt.Fprintln(os.Stderr, "error: -token is required")
		flag.Usage()
		os.Exit(1)
	}
	if *namespaceF == "" {
		fmt.Fprintln(os.Stderr, "error: -namespace is required")
		flag.Usage()
		os.Exit(1)
	}
	if *saNameF == "" {
		fmt.Fprintln(os.Stderr, "error: -service-account is required")
		flag.Usage()
		os.Exit(1)
	}

	reqBody := ExchangeTokenRequest{
		Namespace:          *namespaceF,
		ServiceAccountName: *saNameF,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error marshalling request: %v\n", err)
		os.Exit(1)
	}

	client := &http.Client{Timeout: time.Duration(*timeoutF) * time.Second}

	req, err := http.NewRequest(http.MethodPost, *urlF+"/exchangeToken", bytes.NewReader(body))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating request: %v\n", err)
		os.Exit(1)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+*tokenF)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error sending request: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading response: %v\n", err)
		os.Exit(1)
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "error: server returned %s\n%s\n", resp.Status, string(respBody))
		os.Exit(1)
	}

	// Pretty-print the JSON response
	var out bytes.Buffer
	if err := json.Indent(&out, respBody, "", "    "); err != nil {
		fmt.Println(string(respBody))
	} else {
		fmt.Println(out.String())
	}
}

package kfca

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"

	authenticationV1 "k8s.io/api/authentication/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type ServiceAccountInfo struct {
	Issuer  string `json:"issuer"`
	Subject string `json:"subject"`
}

type ExchangeTokenRequest struct {
	Namespace          string `json:"namespace"`
	ServiceAccountName string `json:"serviceAccountName"`
}

type Status struct {
	Token string `json:"token"`
}

type StatusResult struct {
	Status *Status `json:"status"`
}

type ErrorResponse struct {
	Name    string `json:"name"`
	ID      string `json:"id"`
	Message string `json:"message"`
}

type tokenExchangesrvc struct {
	logger *log.Logger
	debug  bool
}

func NewTokenExchangeHandler(logger *log.Logger, debug bool) http.HandlerFunc {
	s := &tokenExchangesrvc{logger: logger, debug: debug}
	return s.ServeHTTP
}

func (s *tokenExchangesrvc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Method Not Allowed")
		return
	}

	// Check content-type of incoming request
	ct := r.Header.Get("Content-Type")
	if ct == "" {
		writeError(w, http.StatusBadRequest, "bad_request", "Content-Type header must be application/json")
		return
	}
	if ct != "application/json" {
		writeError(w, http.StatusUnsupportedMediaType, "unsupported_media_type", "Content-Type header must be application/json")
		return
	}

	// Check Authorization header
	authorization := r.Header.Get("Authorization")
	if authorization == "" {
		writeError(w, http.StatusForbidden, "forbidden", "missing token")
		return
	}

	// Validate Bearer prefix
	if !strings.HasPrefix(authorization, "Bearer ") {
		writeError(w, http.StatusBadRequest, "bad_request", "Authorization header must start with 'Bearer '")
		return
	}
	token := strings.TrimPrefix(authorization, "Bearer ")
	if len(authorization) > 2000 {
		writeError(w, http.StatusBadRequest, "bad_request", "Authorization token too long")
		return
	}

	// Decode request body
	var req ExchangeTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "bad_request", "invalid request body: "+err.Error())
		return
	}

	if req.Namespace == "" {
		writeError(w, http.StatusBadRequest, "bad_request", "namespace is required")
		return
	}
	if req.ServiceAccountName == "" {
		writeError(w, http.StatusBadRequest, "bad_request", "serviceAccountName is required")
		return
	}
	if len(req.Namespace) > 253 {
		writeError(w, http.StatusBadRequest, "bad_request", "namespace too long")
		return
	}
	if len(req.ServiceAccountName) > 253 {
		writeError(w, http.StatusBadRequest, "bad_request", "serviceAccountName too long")
		return
	}

	res, statusCode, err := s.exchangeToken(r.Context(), token, &req)
	if err != nil {
		writeError(w, statusCode, "error", err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (s *tokenExchangesrvc) exchangeToken(ctx context.Context, authorization string, req *ExchangeTokenRequest) (*StatusResult, int, error) {
	const hourSeconds = 3600

	s.logger.Print("tokenExchange.exchangeToken")

	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("internal server error")
	}
	// creates the clientSet
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("internal server error")
	}

	serviceAccount, err := clientSet.CoreV1().ServiceAccounts(req.Namespace).Get(ctx, req.ServiceAccountName, metav1.GetOptions{})
	if err != nil {
		return nil, http.StatusNotFound, fmt.Errorf("service account not found")
	}

	annotations := serviceAccount.GetAnnotations()
	for key, value := range annotations {
		if strings.Contains(key, "kfca") {

			jsonData := []byte(value)

			var serviceAccountInfo ServiceAccountInfo
			err := json.Unmarshal(jsonData, &serviceAccountInfo)
			if err != nil {
				return nil, http.StatusInternalServerError, fmt.Errorf("internal server error")
			}

			provider, err := oidc.NewProvider(ctx, serviceAccountInfo.Issuer)
			if err != nil {
				if s.debug {
					s.logger.Printf("debug: annotation %q: failed to create OIDC provider for issuer %q: %v", key, serviceAccountInfo.Issuer, err)
				}
				continue
			}

			oidcConfig := &oidc.Config{
				SkipClientIDCheck: true,
			}
			verifier := provider.Verifier(oidcConfig)

			idToken, err := verifier.Verify(ctx, authorization)
			if err != nil {
				if s.debug {
					s.logger.Printf("debug: annotation %q: token verification failed for issuer %q: %v", key, serviceAccountInfo.Issuer, err)
				}
				continue
			}

			// Extract claims from the ID token
			var claims struct {
				Aud        []string `json:"aud"`
				Exp        int      `json:"exp"`
				Iat        int      `json:"iat"`
				Iss        string   `json:"iss"`
				Kubernetes struct {
					Namespace      string `json:"namespace"`
					ServiceAccount struct {
						Name string `json:"name"`
						UID  string `json:"uid"`
					} `json:"serviceaccount"`
				} `json:"kubernetes.io"`
				Nbf int    `json:"nbf"`
				Sub string `json:"sub"`
			}
			if err := idToken.Claims(&claims); err != nil {
				if s.debug {
					s.logger.Printf("debug: annotation %q: failed to extract claims: %v", key, err)
				}
				continue
			}

			if (claims.Exp - claims.Iat) > hourSeconds {
				return nil, http.StatusBadRequest, fmt.Errorf("only tokens with a validity of one hour or less are accepted")
			}

			if claims.Iss == serviceAccountInfo.Issuer && claims.Sub == serviceAccountInfo.Subject {
				tokenRequest := kubernetesAuthToken(hourSeconds)
				token, err := clientSet.CoreV1().ServiceAccounts(req.Namespace).CreateToken(ctx, req.ServiceAccountName, tokenRequest, metav1.CreateOptions{})
				if err != nil {
					if s.debug {
						s.logger.Printf("debug: annotation %q: failed to create token for %s/%s: %v", key, req.Namespace, req.ServiceAccountName, err)
					}
					continue
				}

				return &StatusResult{
					Status: &Status{
						Token: token.Status.Token,
					},
				}, http.StatusOK, nil
			} else if s.debug {
				s.logger.Printf("debug: annotation %q: issuer/subject mismatch: got iss=%q sub=%q, want iss=%q sub=%q", key, claims.Iss, claims.Sub, serviceAccountInfo.Issuer, serviceAccountInfo.Subject)
			}
		}
	}
	return nil, http.StatusNotFound, fmt.Errorf("no matching binding found")
}

func kubernetesAuthToken(expirationSeconds int) *authenticationV1.TokenRequest {
	ExpirationSeconds := int64(expirationSeconds)

	tokenRequest := &authenticationV1.TokenRequest{
		Spec: authenticationV1.TokenRequestSpec{
			ExpirationSeconds: &ExpirationSeconds,
		},
	}

	return tokenRequest
}

func writeError(w http.ResponseWriter, statusCode int, name string, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{
		Name:    name,
		ID:      http.StatusText(statusCode),
		Message: message,
	})
}

package kfca

import (
	"context"
	"encoding/json"
	tokenexchange "k8s-federated-credential-api/gen/token_exchange"
	"log"
	"strings"

	authenticationV1 "k8s.io/api/authentication/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/coreos/go-oidc/v3/oidc"
)

type ServiceAccountInfo struct {
	Issuer             string `json:"issuer"`
	Namespace          string `json:"namespace"`
	ServiceAccountName string `json:"serviceAccountName"`
}

// tokenExchange service example implementation.
// The example methods log the requests and return zero values.
type tokenExchangesrvc struct {
	logger *log.Logger
}

// NewTokenExchange returns the tokenExchange service implementation.
func NewTokenExchange(logger *log.Logger) tokenexchange.Service {
	return &tokenExchangesrvc{logger}
}

// ExchangeToken implements exchangeToken.
func (s *tokenExchangesrvc) ExchangeToken(ctx context.Context, p *tokenexchange.ExchangeTokenPayload) (res string, err error) {
	s.logger.Print("tokenExchange.exchangeToken")

	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		return err.Error(), err
	}
	// creates the clientSet
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err.Error(), err
	}

	serviceAccount, err := clientSet.CoreV1().ServiceAccounts(*p.Namespace).Get(ctx, *p.ServiceAccountName, metav1.GetOptions{})
	if err != nil {
		return "service Account Not Found. Error: %v", err
	}

	annotations := serviceAccount.GetAnnotations()
	for key, value := range annotations {
		if strings.Contains(key, "kfcc") {

			jsonData := []byte(value)

			var serviceAccountInfo ServiceAccountInfo
			err := json.Unmarshal(jsonData, &serviceAccountInfo)
			if err != nil {
				return "Error un-marshalling JSON:", err
			}

			provider, err := oidc.NewProvider(ctx, serviceAccountInfo.Issuer)
			if err != nil {
				continue
				// return "Unable to Reach Bucket:", err
			}

			oidcConfig := &oidc.Config{
				SkipClientIDCheck: true,
			}
			verifier := provider.Verifier(oidcConfig)

			idToken, err := verifier.Verify(ctx, *p.JWT)
			if err != nil {
				continue
				// return "Failed to parse the JWT.\nError: %s", err
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
				continue
				// return "Failed to extract claims: %v", err
			}

			if claims.Iss == serviceAccountInfo.Issuer && claims.Kubernetes.Namespace == serviceAccountInfo.Namespace && claims.Kubernetes.ServiceAccount.Name == serviceAccountInfo.ServiceAccountName {
				const tokenExpirationSeconds = 3600
				tokenRequest := kubernetesAuthToken(tokenExpirationSeconds)
				token, err := clientSet.CoreV1().ServiceAccounts(*p.Namespace).CreateToken(ctx, *p.ServiceAccountName, tokenRequest, metav1.CreateOptions{})
				if err != nil {
					return "Failed to create token: %v", err
				}

				return token.Status.Token, nil
			}
		}
	}

	return "No Matching Binding Found", nil

}

func kubernetesAuthToken(expirationSeconds int) *authenticationV1.TokenRequest {
	ExpirationSeconds := int64(expirationSeconds)

	tokenRequest := &authenticationV1.TokenRequest{
		Spec: authenticationV1.TokenRequestSpec{
			Audiences:         []string{"openshift"},
			ExpirationSeconds: &ExpirationSeconds,
		},
	}

	return tokenRequest
}

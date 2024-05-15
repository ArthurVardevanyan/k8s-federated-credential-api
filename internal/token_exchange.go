package kfca

import (
	"context"
	"encoding/json"
	"fmt"
	tokenexchange "k8s-federated-credential-api/gen/token_exchange"
	"log"
	"strings"

	"goa.design/goa/v3/security"

	authenticationV1 "k8s.io/api/authentication/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/coreos/go-oidc/v3/oidc"
)

type ServiceAccountInfo struct {
	Issuer  string `json:"issuer"`
	Subject string `json:"subject"`
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

// JWTAuth implements the authorization logic for service "tokenExchange" for
// the "jwt" security scheme.
func (s *tokenExchangesrvc) JWTAuth(ctx context.Context, token string, scheme *security.JWTScheme) (context.Context, error) {

	// TBD: add authorization logic.
	//
	// In case of authorization failure this function should return
	// one of the generated error structs, e.g.:
	//
	//    return ctx, myservice.MakeUnauthorizedError("invalid token")
	//
	// Alternatively this function may return an instance of
	// goa.ServiceError with a Name field value that matches one of
	// the design error names, e.g:
	//
	//    return ctx, goa.PermanentError("unauthorized", "invalid token")
	//
	//return ctx, fmt.Errorf("not implemented")
	return ctx, nil

}

// ExchangeToken implements exchangeToken.
func (s *tokenExchangesrvc) ExchangeToken(ctx context.Context, p *tokenexchange.ExchangeTokenPayload) (res *tokenexchange.StatusResult, err error) {
	const hourSeconds = 3600
	res = &tokenexchange.StatusResult{
		Status: &tokenexchange.Status{
			Token: "", // Your desired token value
		},
	}

	s.logger.Print("tokenExchange.exchangeToken")

	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		return res, err
	}
	// creates the clientSet
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		//return err.Error(), err
		return res, err

	}

	serviceAccount, err := clientSet.CoreV1().ServiceAccounts(p.Namespace).Get(ctx, p.ServiceAccountName, metav1.GetOptions{})
	if err != nil {
		//return "service Account Not Found. Error: %v", err
		return res, err
	}

	annotations := serviceAccount.GetAnnotations()
	for key, value := range annotations {
		if strings.Contains(key, "kfca") {

			jsonData := []byte(value)

			var serviceAccountInfo ServiceAccountInfo
			err := json.Unmarshal(jsonData, &serviceAccountInfo)
			if err != nil {
				//return "Error un-marshalling JSON:", err
				return res, err
			}

			provider, err := oidc.NewProvider(ctx, serviceAccountInfo.Issuer)
			if err != nil {
				// If Log Level Debug
				//println(err.Error())
				continue
				// return "Unable to Reach Bucket:", err
			}

			oidcConfig := &oidc.Config{
				SkipClientIDCheck: true,
			}
			verifier := provider.Verifier(oidcConfig)

			idToken, err := verifier.Verify(ctx, *&p.Authorization)
			if err != nil {
				// If Log Level Debug
				//println(err.Error())
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

			if (claims.Exp - claims.Iat) > hourSeconds {
				return res, fmt.Errorf("only tokens with a validity of one hour or less or accepted. ")
			}

			if claims.Iss == serviceAccountInfo.Issuer && claims.Sub == serviceAccountInfo.Subject {
				tokenRequest := kubernetesAuthToken(hourSeconds)
				token, err := clientSet.CoreV1().ServiceAccounts(p.Namespace).CreateToken(ctx, p.ServiceAccountName, tokenRequest, metav1.CreateOptions{})
				if err != nil {
					//res.Token = "Failed to create token: %v" + err.Error()
					return res, err
				}

				res.Status.Token = token.Status.Token
				return res, nil
			} else {
				// If Log Level Debug
				//println(claims.Iss == serviceAccountInfo.Issuer && claims.Sub == serviceAccountInfo.Subject)
			}
		}
	}
	return res, fmt.Errorf("no matching binding found")

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

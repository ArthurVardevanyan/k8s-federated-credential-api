// Code generated by goa v3.16.1, DO NOT EDIT.
//
// tokenExchange HTTP client CLI support package
//
// Command:
// $ goa gen k8s-federated-credential-api/design

package client

import (
	"encoding/json"
	"fmt"
	tokenexchange "k8s-federated-credential-api/gen/token_exchange"
	"unicode/utf8"

	goa "goa.design/goa/v3/pkg"
)

// BuildExchangeTokenPayload builds the payload for the tokenExchange
// exchangeToken endpoint from CLI flags.
func BuildExchangeTokenPayload(tokenExchangeExchangeTokenBody string, tokenExchangeExchangeTokenAuthorization string) (*tokenexchange.ExchangeTokenPayload, error) {
	var err error
	var body ExchangeTokenRequestBody
	{
		err = json.Unmarshal([]byte(tokenExchangeExchangeTokenBody), &body)
		if err != nil {
			return nil, fmt.Errorf("invalid JSON for body, \nerror: %s, \nexample of valid JSON:\n%s", err, "'{\n      \"namespace\": \"q0t\",\n      \"serviceAccountName\": \"ki0\"\n   }'")
		}
		err = goa.MergeErrors(err, goa.ValidatePattern("body.namespace", body.Namespace, "[a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*"))
		if utf8.RuneCountInString(body.Namespace) > 253 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("body.namespace", body.Namespace, utf8.RuneCountInString(body.Namespace), 253, false))
		}
		err = goa.MergeErrors(err, goa.ValidatePattern("body.serviceAccountName", body.ServiceAccountName, "[a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*"))
		if utf8.RuneCountInString(body.ServiceAccountName) > 253 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("body.serviceAccountName", body.ServiceAccountName, utf8.RuneCountInString(body.ServiceAccountName), 253, false))
		}
		if err != nil {
			return nil, err
		}
	}
	var authorization string
	{
		authorization = tokenExchangeExchangeTokenAuthorization
		err = goa.MergeErrors(err, goa.ValidatePattern("Authorization", authorization, "^Bearer [^ ]+"))
		if utf8.RuneCountInString(authorization) > 2000 {
			err = goa.MergeErrors(err, goa.InvalidLengthError("Authorization", authorization, utf8.RuneCountInString(authorization), 2000, false))
		}
		if err != nil {
			return nil, err
		}
	}
	v := &tokenexchange.ExchangeTokenPayload{
		Namespace:          body.Namespace,
		ServiceAccountName: body.ServiceAccountName,
	}
	v.Authorization = authorization

	return v, nil
}

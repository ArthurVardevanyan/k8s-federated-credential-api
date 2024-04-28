// Code generated by goa v3.16.1, DO NOT EDIT.
//
// tokenExchange HTTP client encoders and decoders
//
// Command:
// $ goa gen k8s-federated-credential-api/design

package client

import (
	"bytes"
	"context"
	"io"
	tokenexchange "k8s-federated-credential-api/gen/token_exchange"
	"net/http"
	"net/url"

	goahttp "goa.design/goa/v3/http"
)

// BuildExchangeTokenRequest instantiates a HTTP request object with method and
// path set to call the "tokenExchange" service "exchangeToken" endpoint
func (c *Client) BuildExchangeTokenRequest(ctx context.Context, v any) (*http.Request, error) {
	u := &url.URL{Scheme: c.scheme, Host: c.host, Path: ExchangeTokenTokenExchangePath()}
	req, err := http.NewRequest("POST", u.String(), nil)
	if err != nil {
		return nil, goahttp.ErrInvalidURL("tokenExchange", "exchangeToken", u.String(), err)
	}
	if ctx != nil {
		req = req.WithContext(ctx)
	}

	return req, nil
}

// EncodeExchangeTokenRequest returns an encoder for requests sent to the
// tokenExchange exchangeToken server.
func EncodeExchangeTokenRequest(encoder func(*http.Request) goahttp.Encoder) func(*http.Request, any) error {
	return func(req *http.Request, v any) error {
		p, ok := v.(*tokenexchange.ExchangeTokenPayload)
		if !ok {
			return goahttp.ErrInvalidType("tokenExchange", "exchangeToken", "*tokenexchange.ExchangeTokenPayload", v)
		}
		body := NewExchangeTokenRequestBody(p)
		if err := encoder(req).Encode(&body); err != nil {
			return goahttp.ErrEncodingError("tokenExchange", "exchangeToken", err)
		}
		return nil
	}
}

// DecodeExchangeTokenResponse returns a decoder for responses returned by the
// tokenExchange exchangeToken endpoint. restoreBody controls whether the
// response body should be restored after having been read.
func DecodeExchangeTokenResponse(decoder func(*http.Response) goahttp.Decoder, restoreBody bool) func(*http.Response) (any, error) {
	return func(resp *http.Response) (any, error) {
		if restoreBody {
			b, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			resp.Body = io.NopCloser(bytes.NewBuffer(b))
			defer func() {
				resp.Body = io.NopCloser(bytes.NewBuffer(b))
			}()
		} else {
			defer resp.Body.Close()
		}
		switch resp.StatusCode {
		case http.StatusOK:
			var (
				body string
				err  error
			)
			err = decoder(resp).Decode(&body)
			if err != nil {
				return nil, goahttp.ErrDecodingError("tokenExchange", "exchangeToken", err)
			}
			return body, nil
		default:
			body, _ := io.ReadAll(resp.Body)
			return nil, goahttp.ErrInvalidResponse("tokenExchange", "exchangeToken", resp.StatusCode, string(body))
		}
	}
}

// Code generated by goa v3.16.1, DO NOT EDIT.
//
// kfca HTTP client CLI support package
//
// Command:
// $ goa gen k8s-federated-credential-api/design

package cli

import (
	"flag"
	"fmt"
	livezc "k8s-federated-credential-api/gen/http/livez/client"
	readyzc "k8s-federated-credential-api/gen/http/readyz/client"
	tokenexchangec "k8s-federated-credential-api/gen/http/token_exchange/client"
	"net/http"
	"os"

	goahttp "goa.design/goa/v3/http"
	goa "goa.design/goa/v3/pkg"
)

// UsageCommands returns the set of commands and sub-commands using the format
//
//	command (subcommand1|subcommand2|...)
func UsageCommands() string {
	return `token-exchange exchange-token
readyz readyz
livez livez
`
}

// UsageExamples produces an example of a valid invocation of the CLI tool.
func UsageExamples() string {
	return os.Args[0] + ` token-exchange exchange-token --body '{
      "namespace": "q0t",
      "serviceAccountName": "ki0"
   }' --authorization "y12"` + "\n" +
		os.Args[0] + ` readyz readyz` + "\n" +
		os.Args[0] + ` livez livez` + "\n" +
		""
}

// ParseEndpoint returns the endpoint and payload as specified on the command
// line.
func ParseEndpoint(
	scheme, host string,
	doer goahttp.Doer,
	enc func(*http.Request) goahttp.Encoder,
	dec func(*http.Response) goahttp.Decoder,
	restore bool,
) (goa.Endpoint, any, error) {
	var (
		tokenExchangeFlags = flag.NewFlagSet("token-exchange", flag.ContinueOnError)

		tokenExchangeExchangeTokenFlags             = flag.NewFlagSet("exchange-token", flag.ExitOnError)
		tokenExchangeExchangeTokenBodyFlag          = tokenExchangeExchangeTokenFlags.String("body", "REQUIRED", "")
		tokenExchangeExchangeTokenAuthorizationFlag = tokenExchangeExchangeTokenFlags.String("authorization", "REQUIRED", "")

		readyzFlags = flag.NewFlagSet("readyz", flag.ContinueOnError)

		readyzReadyzFlags = flag.NewFlagSet("readyz", flag.ExitOnError)

		livezFlags = flag.NewFlagSet("livez", flag.ContinueOnError)

		livezLivezFlags = flag.NewFlagSet("livez", flag.ExitOnError)
	)
	tokenExchangeFlags.Usage = tokenExchangeUsage
	tokenExchangeExchangeTokenFlags.Usage = tokenExchangeExchangeTokenUsage

	readyzFlags.Usage = readyzUsage
	readyzReadyzFlags.Usage = readyzReadyzUsage

	livezFlags.Usage = livezUsage
	livezLivezFlags.Usage = livezLivezUsage

	if err := flag.CommandLine.Parse(os.Args[1:]); err != nil {
		return nil, nil, err
	}

	if flag.NArg() < 2 { // two non flag args are required: SERVICE and ENDPOINT (aka COMMAND)
		return nil, nil, fmt.Errorf("not enough arguments")
	}

	var (
		svcn string
		svcf *flag.FlagSet
	)
	{
		svcn = flag.Arg(0)
		switch svcn {
		case "token-exchange":
			svcf = tokenExchangeFlags
		case "readyz":
			svcf = readyzFlags
		case "livez":
			svcf = livezFlags
		default:
			return nil, nil, fmt.Errorf("unknown service %q", svcn)
		}
	}
	if err := svcf.Parse(flag.Args()[1:]); err != nil {
		return nil, nil, err
	}

	var (
		epn string
		epf *flag.FlagSet
	)
	{
		epn = svcf.Arg(0)
		switch svcn {
		case "token-exchange":
			switch epn {
			case "exchange-token":
				epf = tokenExchangeExchangeTokenFlags

			}

		case "readyz":
			switch epn {
			case "readyz":
				epf = readyzReadyzFlags

			}

		case "livez":
			switch epn {
			case "livez":
				epf = livezLivezFlags

			}

		}
	}
	if epf == nil {
		return nil, nil, fmt.Errorf("unknown %q endpoint %q", svcn, epn)
	}

	// Parse endpoint flags if any
	if svcf.NArg() > 1 {
		if err := epf.Parse(svcf.Args()[1:]); err != nil {
			return nil, nil, err
		}
	}

	var (
		data     any
		endpoint goa.Endpoint
		err      error
	)
	{
		switch svcn {
		case "token-exchange":
			c := tokenexchangec.NewClient(scheme, host, doer, enc, dec, restore)
			switch epn {
			case "exchange-token":
				endpoint = c.ExchangeToken()
				data, err = tokenexchangec.BuildExchangeTokenPayload(*tokenExchangeExchangeTokenBodyFlag, *tokenExchangeExchangeTokenAuthorizationFlag)
			}
		case "readyz":
			c := readyzc.NewClient(scheme, host, doer, enc, dec, restore)
			switch epn {
			case "readyz":
				endpoint = c.Readyz()
				data = nil
			}
		case "livez":
			c := livezc.NewClient(scheme, host, doer, enc, dec, restore)
			switch epn {
			case "livez":
				endpoint = c.Livez()
				data = nil
			}
		}
	}
	if err != nil {
		return nil, nil, err
	}

	return endpoint, data, nil
}

// token-exchangeUsage displays the usage of the token-exchange command and its
// subcommands.
func tokenExchangeUsage() {
	fmt.Fprintf(os.Stderr, `Exchange an incoming Kubernetes Token for Another Kubernetes Token
Usage:
    %[1]s [globalflags] token-exchange COMMAND [flags]

COMMAND:
    exchange-token: ExchangeToken implements exchangeToken.

Additional help:
    %[1]s token-exchange COMMAND --help
`, os.Args[0])
}
func tokenExchangeExchangeTokenUsage() {
	fmt.Fprintf(os.Stderr, `%[1]s [flags] token-exchange exchange-token -body JSON -authorization STRING

ExchangeToken implements exchangeToken.
    -body JSON: 
    -authorization STRING: 

Example:
    %[1]s token-exchange exchange-token --body '{
      "namespace": "q0t",
      "serviceAccountName": "ki0"
   }' --authorization "y12"
`, os.Args[0])
}

// readyzUsage displays the usage of the readyz command and its subcommands.
func readyzUsage() {
	fmt.Fprintf(os.Stderr, `Service is the readyz service interface.
Usage:
    %[1]s [globalflags] readyz COMMAND [flags]

COMMAND:
    readyz: Readyz implements readyz.

Additional help:
    %[1]s readyz COMMAND --help
`, os.Args[0])
}
func readyzReadyzUsage() {
	fmt.Fprintf(os.Stderr, `%[1]s [flags] readyz readyz

Readyz implements readyz.

Example:
    %[1]s readyz readyz
`, os.Args[0])
}

// livezUsage displays the usage of the livez command and its subcommands.
func livezUsage() {
	fmt.Fprintf(os.Stderr, `Service is the livez service interface.
Usage:
    %[1]s [globalflags] livez COMMAND [flags]

COMMAND:
    livez: Livez implements livez.

Additional help:
    %[1]s livez COMMAND --help
`, os.Args[0])
}
func livezLivezUsage() {
	fmt.Fprintf(os.Stderr, `%[1]s [flags] livez livez

Livez implements livez.

Example:
    %[1]s livez livez
`, os.Args[0])
}

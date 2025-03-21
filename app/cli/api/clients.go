package api

import (
	"gpt4cli/auth"
	"gpt4cli/types"
	"net"
	"net/http"
	"os"
	"time"
)

const dialTimeout = 10 * time.Second
const fastReqTimeout = 30 * time.Second
const slowReqTimeout = 5 * time.Minute

type Api struct{}

var cloudApiHost string

var Client types.ApiClient = (*Api)(nil)

func init() {
	if os.Getenv("GPT4CLI_ENV") == "development" {
		cloudApiHost = os.Getenv("GPT4CLI_API_HOST")
		if cloudApiHost == "" {
			cloudApiHost = "http://localhost:8080"
		}
	} else {
		cloudApiHost = "https://api.gpt4cli.khulnasoft.com"
	}
}

func getApiHost() string {
	if auth.Current == nil {
		return ""
	} else if auth.Current.IsCloud {
		return cloudApiHost
	} else {
		return auth.Current.Host
	}
}

type authenticatedTransport struct {
	underlyingTransport http.RoundTripper
}

// RoundTrip executes a single HTTP transaction and adds a custom header
func (t *authenticatedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	err := auth.SetAuthHeader(req)
	if err != nil {
		return nil, err
	}
	return t.underlyingTransport.RoundTrip(req)
}

var netDialer = &net.Dialer{
	Timeout: dialTimeout,
}

var unauthenticatedClient = &http.Client{
	Transport: &http.Transport{
		Dial: netDialer.Dial,
	},
	Timeout: fastReqTimeout,
}

var authenticatedFastClient = &http.Client{
	Transport: &authenticatedTransport{
		underlyingTransport: &http.Transport{
			Dial: netDialer.Dial,
		},
	},
	Timeout: fastReqTimeout,
}

var authenticatedSlowClient = &http.Client{
	Transport: &authenticatedTransport{
		underlyingTransport: &http.Transport{
			Dial: netDialer.Dial,
		},
	},
	Timeout: slowReqTimeout,
}

var authenticatedStreamingClient = &http.Client{
	Transport: &authenticatedTransport{
		underlyingTransport: &http.Transport{
			Dial: netDialer.Dial,
		},
	},
	// No global timeout set for the streaming client
}

package api

import (
	"net"
	"net/http"
	"os"
	"gpt4cli/auth"
	"gpt4cli/types"
	"time"
)

const dialTimeout = 10 * time.Second
const fastReqTimeout = 30 * time.Second
const slowReqTimeout = 5 * time.Minute

type Api struct{}

var CloudApiHost string
var Client types.ApiClient = (*Api)(nil)

func init() {
	if os.Getenv("GPT4CLI_ENV") == "development" {
		CloudApiHost = os.Getenv("GPT4CLI_API_HOST")
		if CloudApiHost == "" {
			CloudApiHost = "http://localhost:8080"
		}
	} else {
		CloudApiHost = "https://api.gpt4cli.khulnasoft.com"
	}
}

func GetApiHost() string {
	if auth.Current == nil {
		return CloudApiHost
	} else if auth.Current.IsCloud {
		return CloudApiHost
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

type unauthenticatedTransport struct {
	underlyingTransport http.RoundTripper
}

func (t *unauthenticatedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return t.underlyingTransport.RoundTrip(req)
}

var netDialer = &net.Dialer{
	Timeout: dialTimeout,
}

var unauthenticatedClient = &http.Client{
	Transport: &unauthenticatedTransport{
		underlyingTransport: &http.Transport{
			Dial: netDialer.Dial,
		},
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

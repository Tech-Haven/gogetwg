package configs

import (
	"net/http"
	"time"
)

var NETMAKER_API_URL = "https://api.netmaker.techhaven.io"
var NETMAKER_INGRESS_NODE_ID = "0b119419-9ab5-46b8-ad45-2cb40088e877"

type Config struct {
	HttpClient            *http.Client
	NetmakerApiUrl        string
	NetmakerIngressNodeID string
}

func New() (*Config, error) {
	httpClient := newHttpClient()

	return &Config{
		HttpClient:            httpClient,
		NetmakerApiUrl:        NETMAKER_API_URL,
		NetmakerIngressNodeID: NETMAKER_INGRESS_NODE_ID,
	}, nil
}

func newHttpClient() *http.Client {
	// Define HTTP Client transport options
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 100
	t.MaxConnsPerHost = 100
	t.MaxIdleConnsPerHost = 100

	// Create HTTP client
	client := &http.Client{
		Timeout:   time.Second * 60,
		Transport: t,
	}

	return client
}

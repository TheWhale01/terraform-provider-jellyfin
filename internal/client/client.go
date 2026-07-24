package client

import (
	"net/http"
	"time"
)

type Client struct {
	Endpoint string
	APIKey string
	HTTPClient *http.Client
}

func NewClient(endpoint, apiKey string) *Client {
	return &Client {
		Endpoint: endpoint,
		APIKey: apiKey,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

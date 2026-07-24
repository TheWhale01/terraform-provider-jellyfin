package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	Endpoint string
	APIKey string
	HTTPClient *http.Client
}

type SystemInfo struct {
	Id string `json:"Id"`
	Version string `json:"Version"`
	OperatingSystem string `json:"OperatingSystem"`
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

func (client *Client) GetSystemInfo() (*SystemInfo, error) {
	req, err := http.NewRequest("GET", client.Endpoint + "/System/Info", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", `MediaBrowser Token="`+client.APIKey+`"`)
	resp, err := client.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected status code: %d", resp.StatusCode)
	}
	var info SystemInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, err
	}
	return &info, nil
}

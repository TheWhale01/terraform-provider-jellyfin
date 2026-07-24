package client

import (
	"fmt"
	"net/http"
	"encoding/json"
)

type SystemInfo struct {
	Id string `json:"Id"`
	Version string `json:"Version"`
}

func (client *Client) GetSystemInfo() (*SystemInfo, error) {
	req, err := http.NewRequest(http.MethodGet, client.Endpoint + "/System/Info", nil)
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

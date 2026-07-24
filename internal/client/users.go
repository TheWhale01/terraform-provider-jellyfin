package client

import (
	"fmt"
	"bytes"
	"errors"
	"net/http"
	"encoding/json"
)

type AuthByNameResult struct {
	AccessToken string `json:"AccessToken"`
}

func (client *Client) AuthenticateByName(username string, password string) (*AuthByNameResult, error) {
	if username == "" || password == "" {
		return nil, errors.New("Missing jellyfin username or password value(s)")
	}

	body := map[string]string {
		"Username": username,
		"Pw": password,
	}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("Marshaling auth request: %w", err)
	}
	req, err := http.NewRequest(http.MethodPost, client.Endpoint + "/Users/AuthenticateByName", bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("Creating auth request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", `MediaBrowser Client="Terraform", Device="Provider", DeviceId="terraform-provider-jellyfin", Version="1.0.0"`)
	resp, err := client.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Executing auth request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Authentication failed to user %s (status %d)", username, resp.StatusCode)
	}
	var result AuthByNameResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decoding auth response: %w", err)
	}
	return &result, nil
}

package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type VirtualFolder struct {
	Name string `json:"Name"`
	ItemId string `json:"ItemId"`
	CollectionType string `json:"CollectionType"`
	Locations []string `json:"Locations"`
}

func (client *Client) AddLibrary(options VirtualFolder) error {
	reqUrl := url.URL { Path: "/Library/VirtualFolders" }
	query := reqUrl.Query()
	query.Set("name", options.Name)
	query.Set("collectionType", options.CollectionType)
	for _, path := range options.Locations {
		query.Set("paths", path)
	}
	reqUrl.RawQuery = query.Encode()

	req, err := http.NewRequest(http.MethodPost, client.Endpoint + reqUrl.RequestURI(), nil)
	if err != nil {
		return fmt.Errorf("Creating request: %w", err)
	}
	req.Header.Set("Authorization", `MediaBrowser Token="`+client.APIKey+`"`)
	resp, err := client.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("Executing request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("Failed to create library (status: %d).", resp.StatusCode)
	}
	return nil
}

func (client *Client) RemoveLibrary(name string) error {
	reqUrl := url.URL { Path: "/Library/VirtualFolders" }
	query := reqUrl.Query()
	query.Set("name", name)
	reqUrl.RawQuery = query.Encode()
	req, err := http.NewRequest(http.MethodDelete, client.Endpoint + reqUrl.RequestURI(), nil)
	if err != nil {
		return fmt.Errorf("Creating request: %w", err)
	}
	req.Header.Set("Authorization", `MediaBrowser Token="`+client.APIKey+`"`)
	resp, err := client.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("Executing request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("Failed to delete library (status: %d).", resp.StatusCode)
	}
	return nil
}

func (client *Client) GetLibraries() ([]VirtualFolder, error) {
	req, err := http.NewRequest(http.MethodGet, client.Endpoint + "/Library/VirtualFolders", nil)
	if err != nil {
		return nil, fmt.Errorf("Creating request: %w", err)
	}
	req.Header.Set("Authorization", `MediaBrowser Token="`+client.APIKey+`"`)
	resp, err := client.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Executing request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Failed to get libraries (status: %d).", resp.StatusCode)
	}
	var folders []VirtualFolder
	if err := json.NewDecoder(resp.Body).Decode(&folders); err != nil {
		return nil, fmt.Errorf("Failed to parse API response: %w", err)
	}
	return folders, nil
}

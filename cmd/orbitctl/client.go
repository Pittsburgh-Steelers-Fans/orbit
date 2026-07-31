package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// Client calls the Orbit HTTP API.
type Client struct {
	baseURL string
	http    *http.Client
}

// NewClient builds a client from the ORBIT_API_URL environment variable.
func NewClient() (*Client, error) {
	baseURL := os.Getenv("ORBIT_API_URL")
	if baseURL == "" {
		return nil, fmt.Errorf("ORBIT_API_URL is required")
	}
	return &Client{baseURL: baseURL, http: http.DefaultClient}, nil
}

func (c *Client) get(path string, out any) error {
	response, err := c.http.Get(c.baseURL + path)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("GET %s returned %s", path, response.Status)
	}
	return json.NewDecoder(response.Body).Decode(out)
}

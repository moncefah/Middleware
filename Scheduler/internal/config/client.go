package config

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Agenda struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	UcaID string `json:"ucaId"`
}

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL: strings.TrimRight(baseURL, "/"),
		HTTPClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func (c *Client) FetchAgendas(ctx context.Context) ([]Agenda, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+"/agendas", nil)
	if err != nil {
		return nil, fmt.Errorf("create agendas request: %w", err)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch agendas: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var agendas []Agenda
	if err := json.NewDecoder(resp.Body).Decode(&agendas); err != nil {
		return nil, fmt.Errorf("decode agendas: %w", err)
	}

	return agendas, nil
}

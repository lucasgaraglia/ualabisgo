package ualabisgo

import (
	"sync"
	"time"
)

// Client is the Ualá Bis api client.
type Client struct {
	BaseURL        string
	ClientID       string
	ClientSecretId string
	Username       string

	token     string
	expiresAt time.Time
	mu        sync.Mutex
}

// NewProductionClient creates a new Ualá Bis api client for the production environment.
// Make sure to use the correct credentials for the production environment.
func NewProductionClient(clientId, clientSecretId, username string) *Client {
	return &Client{
		BaseURL:        PRODUCTION_URL,
		ClientID:       clientId,
		ClientSecretId: clientSecretId,
		Username:       username,
	}
}

// NewStageClient creates a new Ualá Bis api client for the stage environment.
// Make sure to use the correct credentials for the stage environment.
func NewStageClient(clientId, clientSecretId, username string) *Client {
	return &Client{
		BaseURL:        STAGE_URL,
		ClientID:       clientId,
		ClientSecretId: clientSecretId,
		Username:       username,
	}
}

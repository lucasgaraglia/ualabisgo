package ualabisgo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type tokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

// issueAccessToken issues a new access token to the Ualá bis api.
func (c *Client) issueAccessToken() (tokenResponse, error) {

	tokenURL := c.BaseURL + "/auth/token"

	type auth struct {
		Username       string `json:"username"`
		ClientId       string `json:"client_id"`
		ClientSecretId string `json:"client_secret_id"`
		GrantType      string `json:"grant_type"`
	}

	var _auth auth

	_auth = auth{
		Username:       c.Username,
		ClientId:       c.ClientID,
		ClientSecretId: c.ClientSecretId,
		GrantType:      "client_credentials",
	}

	data, err := json.Marshal(_auth)
	if err != nil {
		return tokenResponse{}, err
	}

	resp, err := http.Post(tokenURL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return tokenResponse{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		_, _ = io.ReadAll(resp.Body)

		return tokenResponse{}, fmt.Errorf("uala auth failed: %s", resp.Status)
	}

	var token tokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return tokenResponse{}, err
	}

	return token, nil
}

// getToken returns a valid access token, issuing a new one if necessary.
func (c *Client) getToken() (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.token != "" && time.Now().Before(c.expiresAt.Add(-1*time.Minute)) {
		return c.token, nil
	}
	resp, err := c.issueAccessToken()
	if err != nil {
		return "", err
	}

	c.token = resp.AccessToken
	c.expiresAt = time.Now().Add(time.Duration(resp.ExpiresIn) * time.Second)

	return c.token, nil
}

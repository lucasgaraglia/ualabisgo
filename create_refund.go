package ualabisgo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// RefundOrderParams represents the parameters for the RefundOrder endpoint.
type RefundOrderParams struct {
	OrderUUID       string // uuid of the order to be refunded
	Amount          string // amount to be refunded; int or float with 2 decimals
	NotificationURL string // url to be notified when the refund is completed
}

// RefundOrderResponse represents the response from the RefundOrder endpoint.
type RefundOrderResponse struct {
	// 200 OK
	Status string `json:"status"`

	// 400/401/403/404/500
	Code    string   `json:"code"`
	Message string   `json:"message"`
	Errors  []string `json:"errors"`
}

// RefundOrder creates a refund for an order.
func (c *Client) RefundOrder(request RefundOrderParams) (RefundOrderResponse, error) {
	url := fmt.Sprintf("%s/v2/api/orders/%s/refund", c.BaseURL, request.OrderUUID)

	accessToken, err := c.getToken()
	if err != nil {
		return RefundOrderResponse{}, fmt.Errorf("error getting access token: %v", err)
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return RefundOrderResponse{}, fmt.Errorf("error serializing payload: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return RefundOrderResponse{}, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	clientHTTP := &http.Client{}
	resp, err := clientHTTP.Do(req)
	if err != nil {
		return RefundOrderResponse{}, fmt.Errorf("error executing request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return RefundOrderResponse{}, fmt.Errorf("error reading response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return RefundOrderResponse{}, fmt.Errorf("error response: %s", body)
	}

	var response RefundOrderResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return RefundOrderResponse{}, fmt.Errorf("error deserializing response: %v", err)
	}

	return response, nil
}

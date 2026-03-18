package ualabisgo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// CreateOrderResponse is the response from the CreateOrder method.
type CreateOrderResponse struct {
	// 200 OK
	UUID              string `json:"uuid"`
	Status            string `json:"status"`
	ExternalReference string `json:"external_reference"`
	Links             struct {
		CheckoutLink string `json:"checkout_link"`
		Success      string `json:"success"`
		Failed       string `json:"failed"`
	} `json:"links"`
	Amount float64 `json:"amount"`

	// 400/401/403/500
	Code    string   `json:"code"`
	Message string   `json:"message"`
	Errors  []string `json:"errors"`
}

// CreateOrderParams is the parameters for the CreateOrder method.
type CreateOrderParams struct {
	Amount            string // amount to be charged; int or float with 2 decimals
	Description       string // description of the payment
	NotificationURL   string // url to be notified when the payment is completed
	CallbackFail      string // url to be redirected when the payment is failed
	CallbackSuccess   string // url to be redirected when the payment is completed
	ExternalReference string // external reference of the payment
}

// CreateOrder creates an order for the given parameters.
func (c *Client) CreateOrder(params CreateOrderParams) (CreateOrderResponse, error) {

	url := fmt.Sprintf("%s/%s", c.CheckoutBaseURL, "checkout")

	accessToken, err := c.getToken()
	if err != nil {
		return CreateOrderResponse{}, err
	}

	type checkout struct {
		Amount            string `json:"amount"`
		Description       string `json:"description"`
		NotificationURL   string `json:"notification_url"`
		CallbackFail      string `json:"callback_fail"`
		CallbackSuccess   string `json:"callback_success"`
		ExternalReference string `json:"external_reference"`
	}

	_checkout := checkout{
		Amount:            params.Amount,
		Description:       params.Description,
		NotificationURL:   params.NotificationURL,
		CallbackFail:      params.CallbackFail,
		CallbackSuccess:   params.CallbackSuccess,
		ExternalReference: params.ExternalReference,
	}

	data, err := json.Marshal(_checkout)
	if err != nil {
		return CreateOrderResponse{}, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return CreateOrderResponse{}, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return CreateOrderResponse{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return CreateOrderResponse{}, fmt.Errorf("uala checkout failed: %s", resp.Status)
	}

	var response CreateOrderResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return CreateOrderResponse{}, err
	}

	return response, nil

}

package ualabisgo

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type commission struct {
	Amount     string `json:"amount"`
	Percentage string `json:"percentage"`
	Type       string `json:"type"`
}

type tax struct {
	Percentage string `json:"percentage"`
	HeldTax    string `json:"held_tax"`
	Type       string `json:"type"`
}

type customerCard struct {
	HolderName   string `json:"holder_name"`
	Issuer       string `json:"issuer"`
	Pan          string `json:"pan"`
	Installments []struct {
		Number              int     `json:"number"`
		Total               float64 `json:"total"`
		FinancialCost       float64 `json:"financial_cost"`
		ValuePerInstallment float64 `json:"value_per_installment"`
	} `json:"installments"`
}

type customer struct {
	Name string       `json:"name"`
	Card customerCard `json:"card"`
}

type changelogEntry struct {
	NewStatus   string `json:"new_status"`
	OldStatus   string `json:"old_status"`
	UpdatedDate string `json:"updated_date"`
}

// Order represents the order object.
type Order struct {
	UUID              string           `json:"uuid"`
	Amount            float64          `json:"amount"`
	Status            string           `json:"status"`
	ExternalReference string           `json:"external_reference"`
	Commissions       []commission     `json:"commissions"`
	Taxes             []tax            `json:"taxes"`
	Customer          customer         `json:"customer"`
	Changelog         []changelogEntry `json:"changelog"`
	CreatedDate       string           `json:"created_date"`
	UpdatedDate       string           `json:"updated_date"`
}

// GetOrderResponse represents the response from the GetOrder endpoint.
type GetOrderResponse struct {
	// 200 OK
	Order

	// 400/401/403/500
	Code    string   `json:"code"`
	Message string   `json:"message"`
	Errors  []string `json:"errors"`
}

// GetOrder gets an order by its UUID.
func (c *Client) GetOrder(uuid string) (GetOrderResponse, error) {

	accessToken, err := c.getToken()
	if err != nil {
		return GetOrderResponse{}, err
	}

	url := fmt.Sprintf("%s/orders/%s", c.CheckoutBaseURL, uuid)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return GetOrderResponse{}, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return GetOrderResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return GetOrderResponse{}, fmt.Errorf("failed to fetch order, status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return GetOrderResponse{}, err
	}

	var result GetOrderResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return GetOrderResponse{}, err
	}

	return result, nil
}

package ualabisgo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// GetOrdersResponse represents the response from the GetOrders endpoint.
type GetOrdersResponse struct {
	// 200 OK
	LastSearchKey string  `json:"last_search_key"`
	HasMoreItems  bool    `json:"has_more_items"`
	Orders        []Order `json:"orders"`

	// 400/401/403/500
	Code    string   `json:"code"`
	Message string   `json:"message"`
	Errors  []string `json:"errors"`
}

// GetOrdersParams is the parameters for the GetOrders method.
type GetOrdersParams struct {
	Limit         int       // orders per page (default 10)
	FromDate      time.Time // order creation date from (default all)
	ToDate        time.Time // order creation date to (default all)
	Status        string    // order status (default all)
	LastSearchKey string    // last search key (it is used to do another request from the last search key) (default none)
}

// GetOrders returns the orders for the given parameters.
func (c *Client) GetOrders(params GetOrdersParams) (GetOrdersResponse, error) {

	accessToken, err := c.getToken()
	if err != nil {
		return GetOrdersResponse{}, err
	}

	orderURL := fmt.Sprintf("%s/%s", c.BaseURL, "orders")
	request, err := http.NewRequest("GET", orderURL, nil)
	if err != nil {
		return GetOrdersResponse{}, err
	}

	q := request.URL.Query()
	if params.Limit != 0 {
		q.Add("limit", fmt.Sprintf("%d", params.Limit))
	}
	if !params.FromDate.IsZero() {
		q.Add("fromDate", params.FromDate.Format("2006-01-02"))
	}
	if !params.ToDate.IsZero() {
		q.Add("toDate", params.ToDate.Format("2006-01-02"))
	}
	if params.Status != "" {
		q.Add("status", params.Status)
	}
	if params.LastSearchKey != "" {
		q.Add("last_search_key", params.LastSearchKey)
	}
	request.URL.RawQuery = q.Encode()

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+accessToken)
	cli := &http.Client{}
	resp, err := cli.Do(request)
	if err != nil {
		return GetOrdersResponse{}, err
	}
	defer resp.Body.Close()

	var response GetOrdersResponse
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return response, err
	}
	return response, nil
}

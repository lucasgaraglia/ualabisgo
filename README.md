# ualabisgo

`ualabisgo` is an unofficial SDK for Ualá Bis in Go.

## Installation

```bash
go get github.com/lucasgaraglia/ualabisgo
```

## Usage

```go
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/lucasgaraglia/ualabisgo"
)

func main() {
	// Initialize the client for the desired environment
	// client := ualabisgo.NewProductionClient("client_id", "client_secret", "username")
	client := ualabisgo.NewStageClient("client_id", "client_secret", "username")

	// 1. Create an Order
	order, err := client.CreateOrder(ualabisgo.CreateOrderParams{
		Amount:            "100.00",
		Description:       "Test order",
		NotificationURL:   "https://example.com/webhook",
		CallbackFail:      "https://example.com/fail",
		CallbackSuccess:   "https://example.com/success",
		ExternalReference: "ref_12345",
	})
	if err != nil {
		log.Fatalf("Failed to create order: %v", err)
	}
	fmt.Printf("Checkout Link: %s\n", order.Links.CheckoutLink)

	// 2. Get an Order by UUID
	getOrderResponse, err := client.GetOrder(order.UUID)
	if err != nil {
		log.Fatalf("Failed to get order: %v", err)
	}
	fmt.Printf("Order Status: %s\n", getOrderResponse.Status)

	// 3. Get Orders with filtering
	getOrdersResponse, err := client.GetOrders(ualabisgo.GetOrdersParams{
		Limit:         10,
		FromDate:      time.Now().AddDate(0, -1, 0),
		ToDate:        time.Now(),
		Status:        "APPROVED",
		LastSearchKey: "",
	})
	if err != nil {
		log.Fatalf("Failed to get orders: %v", err)
	}
	fmt.Printf("Fetched %d orders\n", len(getOrdersResponse.Orders))

	// 4. Refund an Order
	refundResponse, err := client.RefundOrder(ualabisgo.RefundOrderParams{
		OrderUUID:       order.UUID,
		Amount:          "100.00",
		NotificationURL: "https://example.com/webhook",
	})
	if err != nil {
		log.Fatalf("Failed to refund order: %v", err)
	}
	fmt.Printf("Refund Status: %s\n", refundResponse.Status)
}
```

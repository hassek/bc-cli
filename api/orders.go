package api

import (
	"fmt"
)

// OrderLineItem represents a single coffee preference within an order (for sending to API)
type OrderLineItem struct {
	Quantity      int    `json:"quantity"`
	GrindType     string `json:"grind_type"`
	BrewingMethod string `json:"brewing_method"`
	Notes         string `json:"notes,omitempty"`
}

// OrderLineItemResponse represents a line item received from API
type OrderLineItemResponse struct {
	ID            string `json:"id"`
	Quantity      int    `json:"quantity"`
	GrindType     string `json:"grind_type"`
	BrewingMethod string `json:"brewing_method"`
	Notes         string `json:"notes,omitempty"`
}

// GetQuantity returns the quantity as an int
func (o *OrderLineItemResponse) GetQuantity() int {
	return o.Quantity
}

// Order represents a coffee order
type Order struct {
	ID                   string                  `json:"id"`
	Tier                 string                  `json:"tier"`
	TotalQuantity        int                     `json:"total_quantity"`
	LineItems            []OrderLineItemResponse `json:"line_items"`
	Status               string                  `json:"status"`
	ExpectedShipmentDate *string                 `json:"expected_shipment_date"`
	CreatedOn            string                  `json:"created_on"`
}

// GetTotalQuantity returns the total quantity as an int
func (o *Order) GetTotalQuantity() int {
	return o.TotalQuantity
}

// CreateOrderRequest is the request body for creating an order
type CreateOrderRequest struct {
	Tier          string          `json:"tier,omitempty"`
	ProductID     string          `json:"product_id,omitempty"`
	TotalQuantity int             `json:"total_quantity"`
	LineItems     []OrderLineItem `json:"line_items"`
}

// CreateOrderResponse is the response from creating an order
type CreateOrderResponse struct {
	Meta struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"meta"`
	Data Order `json:"data"`
}

// CheckoutSession contains the checkout URL and session info
type CheckoutSession struct {
	CheckoutURL string `json:"checkout_url"`
	SessionID   string `json:"session_id"`
	OrderID     string `json:"order_id"`
}

// CheckoutResponse is the response from creating a checkout session
type CheckoutResponse struct {
	Meta struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"meta"`
	Data CheckoutSession `json:"data"`
}

// ListOrdersResponse is the response from listing orders
type ListOrdersResponse struct {
	Meta struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"meta"`
	Data []Order `json:"data"`
}

// CreateOrder creates a new draft order with coffee preferences
func (c *Client) CreateOrder(req CreateOrderRequest) (*Order, error) {
	resp, err := c.doRequest("POST", "/api/core/v1/orders/configure", req, true)
	if err != nil {
		return nil, err
	}

	var result CreateOrderResponse
	if err := c.handleResponse(resp, &result); err != nil {
		return nil, err
	}

	if err := validateOrder(&result.Data); err != nil {
		return nil, fmt.Errorf("invalid order response: %w", err)
	}

	return &result.Data, nil
}

// CreateCheckoutSession creates a Stripe checkout session for an order
func (c *Client) CreateCheckoutSession(orderID string) (*CheckoutSession, error) {
	url := fmt.Sprintf("/api/core/v1/orders/%s/checkout", orderID)
	resp, err := c.doRequest("POST", url, nil, true)
	if err != nil {
		return nil, err
	}

	var result CheckoutResponse
	if err := c.handleResponse(resp, &result); err != nil {
		return nil, err
	}

	if err := validateCheckoutSession(&result.Data); err != nil {
		return nil, fmt.Errorf("invalid checkout session: %w", err)
	}

	return &result.Data, nil
}

// GetOrder retrieves a specific order by ID
func (c *Client) GetOrder(orderID string) (*Order, error) {
	url := fmt.Sprintf("/api/core/v1/orders/%s", orderID)
	resp, err := c.doRequest("GET", url, nil, true)
	if err != nil {
		return nil, err
	}

	var result CreateOrderResponse
	if err := c.handleResponse(resp, &result); err != nil {
		return nil, err
	}

	if err := validateOrder(&result.Data); err != nil {
		return nil, fmt.Errorf("invalid order response: %w", err)
	}

	return &result.Data, nil
}

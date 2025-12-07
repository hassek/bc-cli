package api

import (
	"fmt"
	"strconv"
)

// OrderLineItem represents a single coffee preference within an order (for sending to API)
type OrderLineItem struct {
	QuantityKg    int    `json:"quantity_kg"`
	GrindType     string `json:"grind_type"`
	BrewingMethod string `json:"brewing_method"`
	Notes         string `json:"notes,omitempty"`
}

// OrderLineItemResponse represents a line item received from API (with string quantities)
type OrderLineItemResponse struct {
	ID            string `json:"id"`
	QuantityKg    string `json:"quantity_kg"` // Django DecimalField serializes as string
	GrindType     string `json:"grind_type"`
	BrewingMethod string `json:"brewing_method"`
	Notes         string `json:"notes,omitempty"`
}

// GetQuantity returns the quantity as a float64
func (o *OrderLineItemResponse) GetQuantity() float64 {
	val, _ := strconv.ParseFloat(o.QuantityKg, 64)
	return val
}

// Order represents a coffee order
type Order struct {
	ID                   string                      `json:"id"`
	Tier                 string                      `json:"tier"`
	TotalQuantityKg      string                      `json:"total_quantity_kg"` // Django DecimalField serializes as string
	LineItems            []OrderLineItemResponse     `json:"line_items"`
	Status               string                      `json:"status"`
	ExpectedShipmentDate *string                     `json:"expected_shipment_date"`
	CreatedOn            string                      `json:"created_on"`
}

// GetTotalQuantity returns the total quantity as a float64
func (o *Order) GetTotalQuantity() float64 {
	val, _ := strconv.ParseFloat(o.TotalQuantityKg, 64)
	return val
}

// CreateOrderRequest is the request body for creating an order
type CreateOrderRequest struct {
	Tier            string          `json:"tier,omitempty"`
	ProductID       string          `json:"product_id,omitempty"`
	TotalQuantityKg int             `json:"total_quantity_kg"`
	LineItems       []OrderLineItem `json:"line_items"`
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

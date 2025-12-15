package api

import "fmt"

type SubscriptionPaymentLinkRequest struct {
	Tier string `json:"tier"`
}

type SubscriptionPaymentLinkResponse struct {
	Meta struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"meta"`
	Data struct {
		Tier        string `json:"tier"`
		PaymentLink string `json:"payment_link"`
		Message     string `json:"message,omitempty"`
	} `json:"data"`
}

type SubscriptionPaymentLink struct {
	Tier        string `json:"tier"`
	PaymentLink string `json:"payment_link"`
	Message     string `json:"message,omitempty"`
}

type Subscription struct {
	ID                  string                    `json:"id"`
	Tier                string                    `json:"tier"`
	Status              string                    `json:"status"`
	StripePaymentLink   string                    `json:"stripe_payment_link,omitempty"`
	StartedAt           *string                   `json:"started_at"`
	ExpiresAt           *string                   `json:"expires_at"`
	CreatedOn           string                    `json:"created_on"`
	DefaultQuantity     string                    `json:"default_quantity,omitempty"`   // Django DecimalField as string
	DefaultPreferences  []SubscriptionPreference  `json:"default_preferences,omitempty"`   // From GetSubscription endpoint
}

// SubscriptionPreference represents a default coffee preference for a subscription
type SubscriptionPreference struct {
	ID            string `json:"id"`
	Quantity      string `json:"quantity"`    // Django DecimalField as string
	GrindType     string `json:"grind_type"`
	BrewingMethod string `json:"brewing_method"`
	Notes         string `json:"notes,omitempty"`
}

// GetTotalQuantity returns the total quantity as an int
func (s *Subscription) GetTotalQuantity() int {
	if s.DefaultQuantity == "" {
		return 0
	}
	// Parse and round to nearest int
	var qty float64
	if _, err := fmt.Sscanf(s.DefaultQuantity, "%f", &qty); err != nil {
		return 0
	}
	return int(qty + 0.5)
}

// GetQuantity returns the quantity for a preference as an int
func (p *SubscriptionPreference) GetQuantity() int {
	if p.Quantity == "" {
		return 0
	}
	var qty float64
	if _, err := fmt.Sscanf(p.Quantity, "%f", &qty); err != nil {
		return 0
	}
	return int(qty + 0.5)
}

// AvailablePlan represents both subscription tiers and one-time purchase products
type AvailablePlan struct {
	ID            string   `json:"id"`
	Tier          string   `json:"tier"`
	Name          string   `json:"name"`
	Price         string   `json:"price"`
	Currency      string   `json:"currency"`
	BillingPeriod string   `json:"billing_period"`
	Summary       string   `json:"summary"`
	Description   string   `json:"description"`
	Features      []string `json:"features"`
	IsSubscription bool    `json:"is_subscription"`
	IsActive       bool    `json:"is_active"`
}

// AvailableSubscription is an alias for backwards compatibility
// Deprecated: Use AvailablePlan instead
type AvailableSubscription = AvailablePlan

type ListSubscriptionsResponse struct {
	Meta struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"meta"`
	Data []Subscription `json:"data"`
}

func (c *Client) ListSubscriptions() ([]Subscription, error) {
	resp, err := c.doRequest("GET", "/api/core/v1/subscriptions", nil, true)
	if err != nil {
		return nil, err
	}

	var result ListSubscriptionsResponse
	if err := c.handleResponse(resp, &result); err != nil {
		return nil, err
	}

	// Validate each subscription
	for i := range result.Data {
		if err := validateSubscription(&result.Data[i]); err != nil {
			return nil, fmt.Errorf("invalid subscription at index %d: %w", i, err)
		}
	}

	return result.Data, nil
}

type AvailablePlansResponse struct {
	Meta struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"meta"`
	Data []AvailablePlan `json:"data"`
}

// AvailableSubscriptionsResponse is an alias for backwards compatibility
// Deprecated: Use AvailablePlansResponse instead
type AvailableSubscriptionsResponse = AvailablePlansResponse

func (c *Client) GetAvailableSubscriptions() ([]AvailablePlan, error) {
	resp, err := c.doRequest("GET", "/api/core/v1/subscriptions/available?is_subscription=true", nil, false)
	if err != nil {
		return nil, err
	}

	var result AvailablePlansResponse
	if err := c.handleResponse(resp, &result); err != nil {
		return nil, err
	}

	// Filter to only include active subscriptions
	var activeSubscriptions []AvailablePlan
	for _, subscription := range result.Data {
		if subscription.IsActive {
			activeSubscriptions = append(activeSubscriptions, subscription)
		}
	}

	return activeSubscriptions, nil
}

// GetAvailableProducts retrieves all available one-time purchase products
func (c *Client) GetAvailableProducts() ([]AvailablePlan, error) {
	resp, err := c.doRequest("GET", "/api/core/v1/subscriptions/available?is_subscription=false", nil, false)
	if err != nil {
		return nil, err
	}

	var result AvailablePlansResponse
	if err := c.handleResponse(resp, &result); err != nil {
		return nil, err
	}

	// Filter to only include active products
	var activeProducts []AvailablePlan
	for _, product := range result.Data {
		if product.IsActive {
			activeProducts = append(activeProducts, product)
		}
	}

	return activeProducts, nil
}

// GetSubscriptionPricing retrieves pricing information for a specific tier
func (c *Client) GetSubscriptionPricing(tier string) (*AvailablePlan, error) {
	plans, err := c.GetAvailableSubscriptions()
	if err != nil {
		return nil, err
	}

	for _, plan := range plans {
		if plan.Tier == tier {
			return &plan, nil
		}
	}

	return nil, fmt.Errorf("pricing not found for tier: %s", tier)
}

// GetSubscription retrieves a specific subscription with preferences by ID
func (c *Client) GetSubscription(subscriptionID string) (*Subscription, error) {
	url := fmt.Sprintf("/api/core/v1/subscriptions/%s/preferences", subscriptionID)
	resp, err := c.doRequest("GET", url, nil, true)
	if err != nil {
		return nil, err
	}

	var result SubscriptionActionResponse
	if err := c.handleResponse(resp, &result); err != nil {
		return nil, err
	}

	if err := validateSubscription(&result.Data); err != nil {
		return nil, fmt.Errorf("invalid subscription response: %w", err)
	}

	return &result.Data, nil
}

type SubscriptionActionResponse struct {
	Meta struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"meta"`
	Data Subscription `json:"data"`
}

// PauseSubscription pauses an active subscription
// Note: The backend does not currently support scheduled resume dates
func (c *Client) PauseSubscription(subscriptionID string) (*Subscription, error) {
	url := fmt.Sprintf("/api/core/v1/subscriptions/%s/pause", subscriptionID)
	resp, err := c.doRequest("POST", url, nil, true)
	if err != nil {
		return nil, err
	}

	var result SubscriptionActionResponse
	if err := c.handleResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result.Data, nil
}

// ResumeSubscription resumes a paused subscription
func (c *Client) ResumeSubscription(subscriptionID string) (*Subscription, error) {
	url := fmt.Sprintf("/api/core/v1/subscriptions/%s/resume", subscriptionID)
	resp, err := c.doRequest("POST", url, nil, true)
	if err != nil {
		return nil, err
	}

	var result SubscriptionActionResponse
	if err := c.handleResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result.Data, nil
}

// CancelSubscription cancels a subscription
func (c *Client) CancelSubscription(subscriptionID string) (*Subscription, error) {
	url := fmt.Sprintf("/api/core/v1/subscriptions/%s/cancel", subscriptionID)
	resp, err := c.doRequest("POST", url, nil, true)
	if err != nil {
		return nil, err
	}

	var result SubscriptionActionResponse
	if err := c.handleResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result.Data, nil
}

type UpdateSubscriptionRequest struct {
	TotalQuantity int             `json:"total_quantity,omitempty"`
	Preferences   []OrderLineItem `json:"preferences,omitempty"`
}

// UpdateSubscription updates subscription preferences (quantity, line items)
func (c *Client) UpdateSubscription(subscriptionID string, req UpdateSubscriptionRequest) (*Subscription, error) {
	url := fmt.Sprintf("/api/core/v1/subscriptions/%s/preferences", subscriptionID)
	resp, err := c.doRequest("PATCH", url, req, true)
	if err != nil {
		return nil, err
	}

	var result SubscriptionActionResponse
	if err := c.handleResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result.Data, nil
}

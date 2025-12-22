package api

import (
	"fmt"
	"net/url"
)

// validateURL checks if a URL string is valid and uses http/https
func validateURL(urlStr string) error {
	if urlStr == "" {
		return nil // Empty URLs are allowed for optional fields
	}

	parsed, err := url.Parse(urlStr)
	if err != nil {
		return fmt.Errorf("invalid URL format: %w", err)
	}

	// Only allow http and https schemes for security
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return fmt.Errorf("URL must use http or https scheme, got: %s", parsed.Scheme)
	}

	return nil
}

// validateStringLength checks if a string exceeds maximum length
func validateStringLength(s string, maxLen int, fieldName string) error {
	if len(s) > maxLen {
		return fmt.Errorf("%s exceeds maximum length of %d characters", fieldName, maxLen)
	}
	return nil
}

// validateSubscription validates a Subscription object
func validateSubscription(sub *Subscription) error {
	if sub == nil {
		return fmt.Errorf("subscription is nil")
	}

	// Validate string lengths
	if err := validateStringLength(sub.ID, 255, "subscription ID"); err != nil {
		return err
	}
	if err := validateStringLength(sub.Tier, 100, "tier"); err != nil {
		return err
	}
	if err := validateStringLength(sub.Status, 50, "status"); err != nil {
		return err
	}

	// Validate URL if present
	if err := validateURL(sub.StripePaymentLink); err != nil {
		return fmt.Errorf("invalid payment link: %w", err)
	}

	// Validate quantity
	qty := sub.GetTotalQuantity()
	if qty < 0 || qty > 1000 {
		return fmt.Errorf("invalid quantity: %d (must be between 0 and 1000)", qty)
	}

	return nil
}

// validateOrder validates an Order object
func validateOrder(order *Order) error {
	if order == nil {
		return fmt.Errorf("order is nil")
	}

	if err := validateStringLength(order.ID, 255, "order ID"); err != nil {
		return err
	}
	if err := validateStringLength(order.Tier, 100, "tier"); err != nil {
		return err
	}
	if err := validateStringLength(order.Status, 50, "status"); err != nil {
		return err
	}

	// Validate quantity
	qty := order.GetTotalQuantity()
	if qty < 0 || qty > 1000 {
		return fmt.Errorf("invalid order quantity: %d (must be between 0 and 1000)", qty)
	}

	// Validate line items count
	if len(order.LineItems) > 50 {
		return fmt.Errorf("too many line items: %d (maximum 50)", len(order.LineItems))
	}

	return nil
}

// validateCheckoutSession validates a CheckoutSession object
func validateCheckoutSession(session *CheckoutSession) error {
	if session == nil {
		return fmt.Errorf("checkout session is nil")
	}

	// Validate checkout URL
	if err := validateURL(session.CheckoutURL); err != nil {
		return fmt.Errorf("invalid checkout URL: %w", err)
	}

	if session.CheckoutURL == "" {
		return fmt.Errorf("checkout URL is required")
	}

	if err := validateStringLength(session.SessionID, 255, "session ID"); err != nil {
		return err
	}
	if err := validateStringLength(session.OrderID, 255, "order ID"); err != nil {
		return err
	}

	return nil
}

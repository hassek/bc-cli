package cmd

import "time"

const (
	// Payment and polling timeouts
	PaymentTimeoutSeconds     = 5 * 60 // 5 minutes for payment completion
	PaymentPollIntervalSeconds = 5      // Poll every 5 seconds

	// Subscription preferences
	DefaultPreferenceQuantityKg = 2 // Default quantity for new preferences

	// Token expiry safety margin
	TokenExpirySafetyMarginSeconds = 30 // Consider token expired 30 seconds before actual expiration
)

// Payment timeout as duration for convenience
var PaymentTimeout = time.Duration(PaymentTimeoutSeconds) * time.Second
var PaymentPollInterval = time.Duration(PaymentPollIntervalSeconds) * time.Second

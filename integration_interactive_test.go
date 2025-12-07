package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/hassek/bc-cli/api"
	"github.com/hassek/bc-cli/config"
)

// TestE2EInteractiveFullLifecycle is an interactive test that requires manual Stripe payment
// This test walks through the complete subscription lifecycle:
// 1. User registration
// 2. Order creation
// 3. Checkout session creation
// 4. Manual payment completion (opens browser, waits for user)
// 5. Subscription verification
// 6. Pause subscription
// 7. Resume subscription
// 8. Update subscription preferences
// 9. Cancel subscription
//
// Run with: BASE_HOSTNAME=http://localhost:8000 RUN_INTEGRATION_TESTS=true INTERACTIVE=true go test -v -run TestE2EInteractiveFullLifecycle -timeout 30m
func TestE2EInteractiveFullLifecycle(t *testing.T) {
	if os.Getenv("RUN_INTEGRATION_TESTS") != "true" {
		t.Skip("Skipping integration test. Set RUN_INTEGRATION_TESTS=true to run")
	}

	if os.Getenv("INTERACTIVE") != "true" {
		t.Skip("Skipping interactive test. Set INTERACTIVE=true to run")
	}

	apiURL := config.GetAPIURL()
	if apiURL == config.DefaultAPIURL {
		t.Log("Warning: Using production API. Set BASE_HOSTNAME=http://localhost:8000 for local testing")
	}
	t.Logf("Testing against: %s", apiURL)

	// Create test configuration
	testCfg := createTestConfig(t)
	client := api.NewClient(testCfg)

	// Generate unique test credentials
	timestamp := time.Now().Unix()
	randSuffix := rand.Intn(10000)
	username := fmt.Sprintf("test_user_%d_%d", timestamp, randSuffix)
	email := fmt.Sprintf("test_%d_%d@butler.test", timestamp, randSuffix)
	password := "TestPassword123!@#"

	separator := "=" + strings.Repeat("=", 70)
	t.Log(separator)
	t.Log("INTERACTIVE INTEGRATION TEST - FULL SUBSCRIPTION LIFECYCLE")
	t.Log(separator)
	t.Logf("")
	t.Logf("This test will:")
	t.Logf("  1. Create a new user account")
	t.Logf("  2. Create an order and checkout session")
	t.Logf("  3. Open Stripe checkout in your browser")
	t.Logf("  4. Wait for you to complete payment")
	t.Logf("  5. Test subscription management (pause, resume, update)")
	t.Logf("  6. Cancel the subscription")
	t.Logf("")
	t.Logf("Test credentials: %s / %s", username, password)
	t.Log(separator)
	t.Logf("")

	// Step 1: User Registration
	var userID string
	t.Run("1_UserRegistration", func(t *testing.T) {
		t.Logf("\n▶ Step 1: Creating user account...")
		userID = testInteractiveUserRegistration(t, client, username, email, password)
		t.Logf("✓ User created: %s\n", userID)
	})

	// Step 2: Get available subscriptions
	var selectedTier string
	t.Run("2_GetAvailableSubscriptions", func(t *testing.T) {
		t.Logf("\n▶ Step 2: Fetching available subscriptions...")
		selectedTier = testInteractiveGetSubscriptions(t, client)
		t.Logf("✓ Selected tier: %s\n", selectedTier)
	})

	// Step 3: Create order
	var orderID string
	t.Run("3_CreateOrder", func(t *testing.T) {
		t.Logf("\n▶ Step 3: Creating order with coffee preferences...")
		orderID = testInteractiveCreateOrder(t, client, selectedTier)
		t.Logf("✓ Order created: %s\n", orderID)
	})

	// Step 4: Create checkout session and wait for payment
	var subscriptionID string
	t.Run("4_CheckoutAndPayment", func(t *testing.T) {
		t.Logf("\n▶ Step 4: Creating checkout session and processing payment...")
		subscriptionID = testInteractiveCheckoutAndPay(t, client, orderID)
		t.Logf("✓ Payment completed! Subscription ID: %s\n", subscriptionID)
	})

	// Step 5: Verify subscription was created
	var originalSubscription *api.Subscription
	t.Run("5_VerifySubscription", func(t *testing.T) {
		t.Logf("\n▶ Step 5: Verifying subscription was created...")
		originalSubscription = testInteractiveVerifySubscription(t, client, subscriptionID)
		t.Logf("✓ Subscription verified\n")
	})

	// Step 6: Pause subscription
	t.Run("6_PauseSubscription", func(t *testing.T) {
		t.Logf("\n▶ Step 6: Pausing subscription...")
		testInteractivePauseSubscription(t, client, subscriptionID)
		t.Logf("✓ Subscription paused\n")
	})

	// Step 7: Resume subscription
	t.Run("7_ResumeSubscription", func(t *testing.T) {
		t.Logf("\n▶ Step 7: Resuming subscription...")
		testInteractiveResumeSubscription(t, client, subscriptionID)
		t.Logf("✓ Subscription resumed\n")
	})

	// Step 8: Update subscription preferences
	t.Run("8_UpdateSubscription", func(t *testing.T) {
		t.Logf("\n▶ Step 8: Updating subscription preferences...")
		testInteractiveUpdateSubscription(t, client, subscriptionID)
		t.Logf("✓ Subscription updated\n")
	})

	// Step 9: Restore original preferences
	t.Run("9_RestorePreferences", func(t *testing.T) {
		t.Logf("\n▶ Step 9: Restoring original preferences...")
		testInteractiveRestoreSubscription(t, client, subscriptionID, originalSubscription)
		t.Logf("✓ Preferences restored\n")
	})

	// Step 10: Cancel subscription
	t.Run("10_CancelSubscription", func(t *testing.T) {
		t.Logf("\n▶ Step 10: Cancelling subscription...")
		testInteractiveCancelSubscription(t, client, subscriptionID)
		t.Logf("✓ Subscription cancelled\n")
	})

	finalSeparator := "=" + strings.Repeat("=", 70)
	t.Logf("")
	t.Log(finalSeparator)
	t.Log("✓ ALL TESTS PASSED - COMPLETE LIFECYCLE VERIFIED")
	t.Log(finalSeparator)
}

// Helper functions for interactive testing

func testInteractiveUserRegistration(t *testing.T, client *api.Client, username, email, password string) string {
	req := api.RegisterRequest{
		Username: username,
		Email:    email,
		Password: password,
	}

	resp, err := client.Register(req)
	if err != nil {
		t.Fatalf("Registration failed: %v", err)
	}

	if resp.Data.ID == "" {
		t.Error("Expected user ID in response")
	}
	if resp.Data.AccessToken == "" {
		t.Error("Expected access token in response")
	}

	t.Logf("  User ID: %s", resp.Data.ID)
	t.Logf("  Username: %s", username)
	t.Logf("  Email: %s", email)

	return resp.Data.ID
}

func testInteractiveGetSubscriptions(t *testing.T, client *api.Client) string {
	subscriptions, err := client.GetAvailableSubscriptions()
	if err != nil {
		t.Fatalf("Failed to get available subscriptions: %v", err)
	}

	if len(subscriptions) == 0 {
		t.Fatal("Expected at least one available subscription")
	}

	t.Logf("  Available tiers:")
	for i, sub := range subscriptions {
		t.Logf("    [%d] %s - %s %s/%s", i+1, sub.Name, sub.Price, sub.Currency, sub.BillingPeriod)
	}

	// Use the first tier
	return subscriptions[0].Tier
}

func testInteractiveCreateOrder(t *testing.T, client *api.Client, tier string) string {
	req := api.CreateOrderRequest{
		Tier:            tier,
		TotalQuantityKg: 5,
		LineItems: []api.OrderLineItem{
			{
				QuantityKg:    3,
				GrindType:     "whole_bean",
				BrewingMethod: "espresso",
				Notes:         "Interactive test - whole beans",
			},
			{
				QuantityKg:    2,
				GrindType:     "ground",
				BrewingMethod: "v60",
				Notes:         "Interactive test - ground coffee",
			},
		},
	}

	order, err := client.CreateOrder(req)
	if err != nil {
		t.Fatalf("Failed to create order: %v", err)
	}

	t.Logf("  Order ID: %s", order.ID)
	t.Logf("  Tier: %s", order.Tier)
	t.Logf("  Total quantity: %.1f kg/month", order.GetTotalQuantity())
	t.Logf("  Line items:")
	for i, item := range order.LineItems {
		t.Logf("    [%d] %s → %.1f kg (%s)", i+1, item.BrewingMethod, item.GetQuantity(), item.GrindType)
	}

	return order.ID
}

func testInteractiveCheckoutAndPay(t *testing.T, client *api.Client, orderID string) string {
	// Create checkout session
	session, err := client.CreateCheckoutSession(orderID)
	if err != nil {
		t.Fatalf("Failed to create checkout session: %v", err)
	}

	t.Logf("  Checkout session created")
	t.Logf("  Session ID: %s", session.SessionID)
	t.Logf("")
	t.Logf("  ┌─────────────────────────────────────────────────────────────┐")
	t.Logf("  │ OPENING BROWSER FOR STRIPE CHECKOUT...                     │")
	t.Logf("  └─────────────────────────────────────────────────────────────┘")
	t.Logf("")

	// Open browser
	if err := openBrowser(session.CheckoutURL); err != nil {
		t.Logf("  ⚠ Could not open browser automatically: %v", err)
		t.Logf("  Please manually open: %s", session.CheckoutURL)
	} else {
		t.Logf("  ✓ Browser opened")
	}

	t.Logf("")
	t.Logf("  Use Stripe test card: 4242 4242 4242 4242")
	t.Logf("  Expiry: Any future date (e.g., 12/34)")
	t.Logf("  CVC: Any 3 digits (e.g., 123)")
	t.Logf("  ZIP: Any 5 digits (e.g., 12345)")
	t.Logf("")

	// Wait for payment with polling
	t.Logf("  Waiting for payment completion...")
	t.Logf("  (Polling every 3 seconds, timeout in 5 minutes)")
	t.Logf("")

	subscriptionID, err := pollForSubscription(t, client, 5*time.Minute, 3*time.Second)
	if err != nil {
		t.Fatalf("Failed to verify payment: %v", err)
	}

	t.Logf("")
	t.Logf("  ✓ Payment verified!")
	t.Logf("  ✓ Subscription created: %s", subscriptionID)

	return subscriptionID
}

func testInteractiveVerifySubscription(t *testing.T, client *api.Client, subscriptionID string) *api.Subscription {
	subscription, err := client.GetSubscription(subscriptionID)
	if err != nil {
		t.Fatalf("Failed to get subscription: %v", err)
	}

	if subscription.Status != "active" {
		t.Errorf("Expected status 'active', got '%s'", subscription.Status)
	}

	t.Logf("  ID: %s", subscription.ID)
	t.Logf("  Tier: %s", subscription.Tier)
	t.Logf("  Status: %s", subscription.Status)
	t.Logf("  Total quantity: %d kg/month", subscription.GetTotalQuantity())

	if len(subscription.DefaultPreferences) > 0 {
		t.Logf("  Preferences:")
		for i, pref := range subscription.DefaultPreferences {
			t.Logf("    [%d] %s → %d kg (%s)", i+1, pref.BrewingMethod, pref.GetQuantity(), pref.GrindType)
		}
	}

	return subscription
}

func testInteractivePauseSubscription(t *testing.T, client *api.Client, subscriptionID string) {
	subscription, err := client.PauseSubscription(subscriptionID)
	if err != nil {
		t.Fatalf("Failed to pause subscription: %v", err)
	}

	if subscription.Status != "paused" {
		t.Errorf("Expected status 'paused', got '%s'", subscription.Status)
	}

	t.Logf("  Status changed to: %s", subscription.Status)
}

func testInteractiveResumeSubscription(t *testing.T, client *api.Client, subscriptionID string) {
	subscription, err := client.ResumeSubscription(subscriptionID)
	if err != nil {
		t.Fatalf("Failed to resume subscription: %v", err)
	}

	if subscription.Status != "active" {
		t.Errorf("Expected status 'active', got '%s'", subscription.Status)
	}

	t.Logf("  Status changed to: %s", subscription.Status)
}

func testInteractiveUpdateSubscription(t *testing.T, client *api.Client, subscriptionID string) {
	// Update with different preferences
	updateReq := api.UpdateSubscriptionRequest{
		TotalQuantityKg: 8,
		Preferences: []api.OrderLineItem{
			{
				QuantityKg:    5,
				GrindType:     "ground",
				BrewingMethod: "french_press",
			},
			{
				QuantityKg:    3,
				GrindType:     "whole_bean",
				BrewingMethod: "aeropress",
			},
		},
	}

	subscription, err := client.UpdateSubscription(subscriptionID, updateReq)
	if err != nil {
		t.Fatalf("Failed to update subscription: %v", err)
	}

	if subscription.GetTotalQuantity() != 8 {
		t.Errorf("Expected total quantity 8, got %d", subscription.GetTotalQuantity())
	}

	t.Logf("  Updated total quantity: %d kg/month", subscription.GetTotalQuantity())

	// Verify the update
	updated, err := client.GetSubscription(subscriptionID)
	if err != nil {
		t.Fatalf("Failed to verify update: %v", err)
	}

	t.Logf("  New preferences:")
	for i, pref := range updated.DefaultPreferences {
		t.Logf("    [%d] %s → %d kg (%s)", i+1, pref.BrewingMethod, pref.GetQuantity(), pref.GrindType)
	}
}

func testInteractiveRestoreSubscription(t *testing.T, client *api.Client, subscriptionID string, original *api.Subscription) {
	var lineItems []api.OrderLineItem
	for _, pref := range original.DefaultPreferences {
		lineItems = append(lineItems, api.OrderLineItem{
			QuantityKg:    pref.GetQuantity(),
			GrindType:     pref.GrindType,
			BrewingMethod: pref.BrewingMethod,
		})
	}

	restoreReq := api.UpdateSubscriptionRequest{
		TotalQuantityKg: original.GetTotalQuantity(),
		Preferences:     lineItems,
	}

	subscription, err := client.UpdateSubscription(subscriptionID, restoreReq)
	if err != nil {
		t.Fatalf("Failed to restore subscription: %v", err)
	}

	t.Logf("  Restored to %d kg/month", subscription.GetTotalQuantity())
}

func testInteractiveCancelSubscription(t *testing.T, client *api.Client, subscriptionID string) {
	subscription, err := client.CancelSubscription(subscriptionID)
	if err != nil {
		t.Fatalf("Failed to cancel subscription: %v", err)
	}

	if subscription.Status != "cancelled" {
		t.Errorf("Expected status 'cancelled', got '%s'", subscription.Status)
	}

	t.Logf("  Status changed to: %s", subscription.Status)
	t.Logf("  Subscription is now permanently cancelled")
}

// pollForSubscription polls the API to check if a subscription was created
func pollForSubscription(t *testing.T, client *api.Client, timeout, interval time.Duration) (string, error) {
	deadline := time.Now().Add(timeout)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	attempt := 0
	for time.Now().Before(deadline) {
		attempt++
		<-ticker.C

		// List subscriptions
		subscriptions, err := client.ListSubscriptions()
		if err != nil {
			t.Logf("  [Attempt %d] Error checking subscriptions: %v", attempt, err)
			continue
		}

		// Check if we have an active subscription
		for _, sub := range subscriptions {
			if sub.Status == "active" {
				return sub.ID, nil
			}
		}

		remaining := time.Until(deadline)
		t.Logf("  [Attempt %d] No active subscription yet, retrying in %d seconds... (%s remaining)",
			attempt, int(interval.Seconds()), remaining.Round(time.Second))
	}

	return "", fmt.Errorf("timeout waiting for subscription (waited %s)", timeout)
}

// openBrowser opens a URL in the default browser
func openBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "darwin":
		cmd = "open"
		args = []string{url}
	case "linux":
		cmd = "xdg-open"
		args = []string{url}
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start", url}
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	return exec.Command(cmd, args...).Start()
}

package main

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/hassek/bc-cli/api"
	"github.com/hassek/bc-cli/config"
)

// TestE2EFullFlow runs a complete end-to-end test against the real backend
// This test requires BASE_HOSTNAME to be set to http://localhost:8000 for local testing
// or to be run against a staging/production environment
//
// Run with: BASE_HOSTNAME=http://localhost:8000 go test -v -run TestE2EFullFlow
func TestE2EFullFlow(t *testing.T) {
	// Skip if not explicitly enabled
	if os.Getenv("RUN_INTEGRATION_TESTS") != "true" {
		t.Skip("Skipping integration test. Set RUN_INTEGRATION_TESTS=true to run")
	}

	// Verify BASE_HOSTNAME is set
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

	t.Logf("Generated test credentials: username=%s, email=%s", username, email)

	// Step 1: Test user registration
	t.Run("UserRegistration", func(t *testing.T) {
		testUserRegistration(t, client, username, email, password)
	})

	// Step 2: Test login (verify we can login with the newly created user)
	t.Run("UserLogin", func(t *testing.T) {
		testUserLogin(t, client, username, password)
	})

	// Step 3: Get available subscriptions
	var selectedTier string
	t.Run("GetAvailableSubscriptions", func(t *testing.T) {
		selectedTier = testGetAvailableSubscriptions(t, client)
	})

	// Step 4: Create an order with preferences
	var orderID string
	t.Run("CreateOrder", func(t *testing.T) {
		orderID = testCreateOrder(t, client, selectedTier)
	})

	// Step 5: Create checkout session (note: we won't complete payment in automated tests)
	t.Run("CreateCheckoutSession", func(t *testing.T) {
		testCreateCheckoutSession(t, client, orderID)
	})

	// Step 6: List subscriptions (should be empty since we didn't complete payment)
	t.Run("ListSubscriptionsEmpty", func(t *testing.T) {
		testListSubscriptions(t, client, 0)
	})

	// For testing subscription management, we need to manually create a subscription
	// or use a test helper. Since we can't complete Stripe payments in automated tests,
	// we'll skip the subscription management tests here and handle them separately.
	t.Log("Note: Subscription management tests require completed payment and are tested separately")
}

// TestE2ESubscriptionManagement tests subscription management flows
// This test requires a pre-existing active subscription for the test user
//
// Run with: BASE_HOSTNAME=http://localhost:8000 go test -v -run TestE2ESubscriptionManagement
func TestE2ESubscriptionManagement(t *testing.T) {
	if os.Getenv("RUN_INTEGRATION_TESTS") != "true" {
		t.Skip("Skipping integration test. Set RUN_INTEGRATION_TESTS=true to run")
	}

	// This test requires TEST_SUBSCRIPTION_ID to be set
	subscriptionID := os.Getenv("TEST_SUBSCRIPTION_ID")
	if subscriptionID == "" {
		t.Skip("Skipping subscription management test. Set TEST_SUBSCRIPTION_ID to run")
	}

	apiURL := config.GetAPIURL()
	t.Logf("Testing against: %s", apiURL)

	// Load existing config (assumes user is already logged in)
	testCfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if !testCfg.IsAuthenticated() {
		t.Fatal("User must be authenticated. Run login first.")
	}

	client := api.NewClient(testCfg)

	// Step 1: Get subscription details
	var originalSubscription *api.Subscription
	t.Run("GetSubscription", func(t *testing.T) {
		originalSubscription = testGetSubscription(t, client, subscriptionID)
	})

	// Step 2: Pause subscription
	t.Run("PauseSubscription", func(t *testing.T) {
		testPauseSubscription(t, client, subscriptionID)
	})

	// Step 3: Resume subscription
	t.Run("ResumeSubscription", func(t *testing.T) {
		testResumeSubscription(t, client, subscriptionID)
	})

	// Step 4: Update subscription preferences
	t.Run("UpdateSubscription", func(t *testing.T) {
		testUpdateSubscription(t, client, subscriptionID, originalSubscription)
	})

	// Step 5: Restore original preferences
	t.Run("RestoreOriginalPreferences", func(t *testing.T) {
		testRestoreSubscription(t, client, subscriptionID, originalSubscription)
	})

	// Step 6: Cancel subscription (optional - leave commented to preserve test data)
	// t.Run("CancelSubscription", func(t *testing.T) {
	// 	testCancelSubscription(t, client, subscriptionID)
	// })
}

// Helper functions

func createTestConfig(t *testing.T) *config.Config {
	// Create a temporary config directory for testing
	tmpDir := t.TempDir()
	testConfigPath := filepath.Join(tmpDir, "test_config.json")

	// Override the config path for this test
	// Note: This requires modifying config package to support custom paths
	// For now, we'll create a config without saving to disk
	cfg := &config.Config{
		APIURL:      config.GetAPIURL(),
		MinQuantity: config.DefaultMinQuantity,
		MaxQuantity: config.DefaultMaxQuantity,
	}

	t.Logf("Created test config with API URL: %s (config path: %s)", cfg.APIURL, testConfigPath)
	return cfg
}

func testUserRegistration(t *testing.T, client *api.Client, username, email, password string) {
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
	if resp.Data.RefreshToken == "" {
		t.Error("Expected refresh token in response")
	}

	t.Logf("✓ User registered successfully: ID=%s", resp.Data.ID)
	t.Logf("✓ Received access token: %s...", resp.Data.AccessToken[:20])
	t.Logf("✓ Received refresh token: %s...", resp.Data.RefreshToken[:20])
}

func testUserLogin(t *testing.T, client *api.Client, username, password string) {
	req := api.LoginRequest{
		Username: username,
		Password: password,
	}

	resp, err := client.Login(req)
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}

	if resp.Data.AccessToken == "" {
		t.Error("Expected access token in response")
	}
	if resp.Data.RefreshToken == "" {
		t.Error("Expected refresh token in response")
	}
	if resp.Data.ExpiresAt == "" {
		t.Error("Expected expires_at in response")
	}
	if resp.Data.RefreshTokenExpiresAt == "" {
		t.Error("Expected refresh_token_expires_at in response")
	}
	if resp.Data.UserID == "" {
		t.Error("Expected user_id in response")
	}

	t.Logf("✓ User logged in successfully: UserID=%s", resp.Data.UserID)
	t.Logf("✓ Access token expires at: %s", resp.Data.ExpiresAt)
	t.Logf("✓ Refresh token expires at: %s", resp.Data.RefreshTokenExpiresAt)
}

func testGetAvailableSubscriptions(t *testing.T, client *api.Client) string {
	subscriptions, err := client.GetAvailableSubscriptions()
	if err != nil {
		t.Fatalf("Failed to get available subscriptions: %v", err)
	}

	if len(subscriptions) == 0 {
		t.Fatal("Expected at least one available subscription")
	}

	t.Logf("✓ Found %d available subscription(s)", len(subscriptions))

	for i, sub := range subscriptions {
		t.Logf("  [%d] Tier: %s, Name: %s, Price: %s %s/%s",
			i+1, sub.Tier, sub.Name, sub.Price, sub.Currency, sub.BillingPeriod)

		// Validate each subscription
		if sub.ID == "" {
			t.Errorf("Subscription %d missing ID", i)
		}
		if sub.Tier == "" {
			t.Errorf("Subscription %d missing tier", i)
		}
		if sub.Name == "" {
			t.Errorf("Subscription %d missing name", i)
		}
		if sub.Price == "" {
			t.Errorf("Subscription %d missing price", i)
		}
		if !sub.IsSubscription {
			t.Errorf("Subscription %d has IsSubscription=false", i)
		}
	}

	// Return the first tier for testing
	selectedTier := subscriptions[0].Tier
	t.Logf("✓ Selected tier for testing: %s", selectedTier)
	return selectedTier
}

func testCreateOrder(t *testing.T, client *api.Client, tier string) string {
	// Create an order with multiple preferences
	req := api.CreateOrderRequest{
		Tier:          tier,
		TotalQuantity: 5,
		LineItems: []api.OrderLineItem{
			{
				Quantity:      3,
				GrindType:     "whole_bean",
				BrewingMethod: "espresso",
				Notes:         "Test order - whole beans for espresso",
			},
			{
				Quantity:      2,
				GrindType:     "ground",
				BrewingMethod: "v60",
				Notes:         "Test order - ground for pour over",
			},
		},
	}

	order, err := client.CreateOrder(req)
	if err != nil {
		t.Fatalf("Failed to create order: %v", err)
	}

	// Validate order response
	if order.ID == "" {
		t.Error("Expected order ID")
	}
	if order.Tier != tier {
		t.Errorf("Expected tier %s, got %s", tier, order.Tier)
	}
	if order.Status == "" {
		t.Error("Expected order status")
	}
	if len(order.LineItems) != 2 {
		t.Errorf("Expected 2 line items, got %d", len(order.LineItems))
	}

	totalQty := order.GetTotalQuantity()
	if totalQty != 5.0 {
		t.Errorf("Expected total quantity 5.0, got %f", totalQty)
	}

	t.Logf("✓ Order created successfully: ID=%s", order.ID)
	t.Logf("✓ Order tier: %s, status: %s", order.Tier, order.Status)
	t.Logf("✓ Total quantity: %.1f kg", totalQty)

	for i, item := range order.LineItems {
		qty := item.GetQuantity()
		t.Logf("  [%d] %s → %.1f kg (%s)", i+1, item.BrewingMethod, qty, item.GrindType)
	}

	return order.ID
}

func testCreateCheckoutSession(t *testing.T, client *api.Client, orderID string) {
	session, err := client.CreateCheckoutSession(orderID)
	if err != nil {
		t.Fatalf("Failed to create checkout session: %v", err)
	}

	if session.CheckoutURL == "" {
		t.Error("Expected checkout URL")
	}
	if session.SessionID == "" {
		t.Error("Expected session ID")
	}
	if session.OrderID != orderID {
		t.Errorf("Expected order ID %s, got %s", orderID, session.OrderID)
	}

	t.Logf("✓ Checkout session created successfully")
	t.Logf("✓ Session ID: %s", session.SessionID)
	t.Logf("✓ Checkout URL: %s", session.CheckoutURL)
	t.Log("Note: Automated tests cannot complete Stripe payment flow")
}

func testListSubscriptions(t *testing.T, client *api.Client, expectedCount int) {
	subscriptions, err := client.ListSubscriptions()
	if err != nil {
		t.Fatalf("Failed to list subscriptions: %v", err)
	}

	if len(subscriptions) != expectedCount {
		t.Errorf("Expected %d subscription(s), got %d", expectedCount, len(subscriptions))
	}

	t.Logf("✓ Found %d subscription(s)", len(subscriptions))

	for i, sub := range subscriptions {
		t.Logf("  [%d] ID: %s, Tier: %s, Status: %s", i+1, sub.ID, sub.Tier, sub.Status)
	}
}

func testGetSubscription(t *testing.T, client *api.Client, subscriptionID string) *api.Subscription {
	subscription, err := client.GetSubscription(subscriptionID)
	if err != nil {
		t.Fatalf("Failed to get subscription: %v", err)
	}

	if subscription.ID != subscriptionID {
		t.Errorf("Expected subscription ID %s, got %s", subscriptionID, subscription.ID)
	}
	if subscription.Tier == "" {
		t.Error("Expected tier in subscription")
	}
	if subscription.Status == "" {
		t.Error("Expected status in subscription")
	}

	t.Logf("✓ Subscription retrieved: ID=%s", subscription.ID)
	t.Logf("✓ Tier: %s, Status: %s", subscription.Tier, subscription.Status)
	t.Logf("✓ Total quantity: %d kg/month", subscription.GetTotalQuantity())

	if len(subscription.DefaultPreferences) > 0 {
		t.Logf("✓ Default preferences:")
		for i, pref := range subscription.DefaultPreferences {
			t.Logf("  [%d] %s → %d kg (%s)",
				i+1, pref.BrewingMethod, pref.GetQuantity(), pref.GrindType)
		}
	}

	return subscription
}

func testPauseSubscription(t *testing.T, client *api.Client, subscriptionID string) {
	subscription, err := client.PauseSubscription(subscriptionID)
	if err != nil {
		t.Fatalf("Failed to pause subscription: %v", err)
	}

	if subscription.Status != "paused" {
		t.Errorf("Expected status 'paused', got '%s'", subscription.Status)
	}

	t.Logf("✓ Subscription paused successfully: ID=%s", subscription.ID)
	t.Logf("✓ Status: %s", subscription.Status)
}

func testResumeSubscription(t *testing.T, client *api.Client, subscriptionID string) {
	subscription, err := client.ResumeSubscription(subscriptionID)
	if err != nil {
		t.Fatalf("Failed to resume subscription: %v", err)
	}

	if subscription.Status != "active" {
		t.Errorf("Expected status 'active', got '%s'", subscription.Status)
	}

	t.Logf("✓ Subscription resumed successfully: ID=%s", subscription.ID)
	t.Logf("✓ Status: %s", subscription.Status)
}

func testUpdateSubscription(t *testing.T, client *api.Client, subscriptionID string, original *api.Subscription) {
	// Update with different preferences
	updateReq := api.UpdateSubscriptionRequest{
		TotalQuantity: 8,
		Preferences: []api.OrderLineItem{
			{
				Quantity:      5,
				GrindType:     "ground",
				BrewingMethod: "french_press",
			},
			{
				Quantity:      3,
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

	t.Logf("✓ Subscription updated successfully: ID=%s", subscription.ID)
	t.Logf("✓ New total quantity: %d kg/month", subscription.GetTotalQuantity())

	// Verify the updated preferences
	updated, err := client.GetSubscription(subscriptionID)
	if err != nil {
		t.Fatalf("Failed to verify updated subscription: %v", err)
	}

	if updated.GetTotalQuantity() != 8 {
		t.Errorf("Verification failed: Expected total quantity 8, got %d", updated.GetTotalQuantity())
	}

	t.Logf("✓ Verified updated preferences:")
	for i, pref := range updated.DefaultPreferences {
		t.Logf("  [%d] %s → %d kg (%s)",
			i+1, pref.BrewingMethod, pref.GetQuantity(), pref.GrindType)
	}
}

func testRestoreSubscription(t *testing.T, client *api.Client, subscriptionID string, original *api.Subscription) {
	// Build preferences from original subscription
	var lineItems []api.OrderLineItem
	for _, pref := range original.DefaultPreferences {
		lineItems = append(lineItems, api.OrderLineItem{
			Quantity:      pref.GetQuantity(),
			GrindType:     pref.GrindType,
			BrewingMethod: pref.BrewingMethod,
		})
	}

	restoreReq := api.UpdateSubscriptionRequest{
		TotalQuantity: original.GetTotalQuantity(),
		Preferences:   lineItems,
	}

	subscription, err := client.UpdateSubscription(subscriptionID, restoreReq)
	if err != nil {
		t.Fatalf("Failed to restore subscription: %v", err)
	}

	if subscription.GetTotalQuantity() != original.GetTotalQuantity() {
		t.Errorf("Expected total quantity %d, got %d",
			original.GetTotalQuantity(), subscription.GetTotalQuantity())
	}

	t.Logf("✓ Subscription restored to original preferences: ID=%s", subscription.ID)
	t.Logf("✓ Total quantity: %d kg/month", subscription.GetTotalQuantity())
}

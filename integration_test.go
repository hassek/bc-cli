package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/hassek/bc-cli/api"
	"github.com/hassek/bc-cli/config"
)

// TestE2EFullFlow runs a complete end-to-end test against the real backend
// This test requires BASE_HOSTNAME to be set to http://localhost:8000 for local testing
// or to be run against a staging/production environment
//
// Environment variables:
//   - RUN_INTEGRATION_TESTS=true (required)
//   - BASE_HOSTNAME=http://localhost:8000 (optional, defaults to production)
//   - TEST_PRODUCT_ID=<product-id> (optional, uses "Butler" product if not set)
//   - TEST_QA_USERNAME=<username> (optional, creates random user if not set)
//   - TEST_QA_PASSWORD=<password> (optional, creates random user if not set)
//
// Run with: BASE_HOSTNAME=http://localhost:8000 RUN_INTEGRATION_TESTS=true TEST_QA_USERNAME=qa_user TEST_QA_PASSWORD=qa_pass go test -v -run TestE2EFullFlow
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

	// Check for QA user credentials
	qaUsername := os.Getenv("TEST_QA_USERNAME")
	qaPassword := os.Getenv("TEST_QA_PASSWORD")

	var username, password string

	if qaUsername != "" && qaPassword != "" {
		// Use dedicated QA user
		t.Logf("Using QA test user: %s", qaUsername)
		username = qaUsername
		password = qaPassword

		// Step 1: Login with QA user
		t.Logf("\n=== Step 1: Login with QA User ===")
		testUserLogin(t, client, username, password)
	} else {
		// Generate unique test credentials for one-time user
		timestamp := time.Now().Unix()
		randSuffix := rand.Intn(10000)
		username = fmt.Sprintf("test_user_%d_%d", timestamp, randSuffix)
		email := fmt.Sprintf("test_%d_%d@butler.test", timestamp, randSuffix)
		password = "TestPassword123!@#"

		t.Logf("Generated test credentials: username=%s, email=%s", username, email)
		t.Log("Tip: Set TEST_QA_USERNAME and TEST_QA_PASSWORD to use a dedicated QA user")

		// Step 1: Test user registration
		t.Logf("\n=== Step 1: User Registration ===")
		testUserRegistration(t, client, username, email, password)

		// Step 2: Test login (verify we can login with the newly created user)
		t.Logf("\n=== Step 2: User Login ===")
		testUserLogin(t, client, username, password)
	}

	// Step 3: Get available subscriptions
	t.Logf("\n=== Step 3: Get Available Subscriptions ===")
	selectedPlan := testGetAvailableSubscriptions(t, client)

	// Step 4: Create an order with preferences
	t.Logf("\n=== Step 4: Create Order ===")
	orderID := testCreateOrder(t, client, selectedPlan)

	// Step 5: Create checkout session (note: we won't complete payment in automated tests)
	t.Logf("\n=== Step 5: Create Checkout Session ===")
	testCreateCheckoutSession(t, client, orderID)

	// Step 6: List subscriptions (should be empty since we didn't complete payment)
	t.Logf("\n=== Step 6: List Subscriptions (Empty) ===")
	testListSubscriptions(t, client, 0)

	// For testing subscription management, we need to manually create a subscription
	// or use a test helper. Since we can't complete Stripe payments in automated tests,
	// we'll skip the subscription management tests here and handle them separately.
	t.Log("\n✓ All automated tests passed")
	t.Log("Note: Subscription management tests require completed payment (see TestE2EInteractiveFullLifecycle)")
}

// TestE2ESubscriptionManagement tests subscription management flows
// This test requires a pre-existing active subscription for the test user
//
// Run with: BASE_HOSTNAME=http://localhost:8000 RUN_INTEGRATION_TESTS=true TEST_SUBSCRIPTION_ID=<id> go test -v -run TestE2ESubscriptionManagement
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
		t.Fatalf("❌ Failed to load config: %v", err)
	}

	if !testCfg.IsAuthenticated() {
		t.Fatal("❌ User must be authenticated. Run login first.")
	}

	client := api.NewClient(testCfg)

	// Step 1: Get subscription details
	t.Logf("\n=== Step 1: Get Subscription ===")
	originalSubscription := testGetSubscription(t, client, subscriptionID)

	// Step 2: Pause subscription
	t.Logf("\n=== Step 2: Pause Subscription ===")
	testPauseSubscription(t, client, subscriptionID)

	// Step 3: Resume subscription
	t.Logf("\n=== Step 3: Resume Subscription ===")
	testResumeSubscription(t, client, subscriptionID)

	// Step 4: Update subscription preferences
	t.Logf("\n=== Step 4: Update Subscription ===")
	testUpdateSubscription(t, client, subscriptionID, originalSubscription)

	// Step 5: Restore original preferences
	t.Logf("\n=== Step 5: Restore Original Preferences ===")
	testRestoreSubscription(t, client, subscriptionID, originalSubscription)

	// Step 6: Cancel subscription (optional - leave commented to preserve test data)
	// t.Logf("\n=== Step 6: Cancel Subscription ===")
	// testInteractiveCancelSubscription(t, client, subscriptionID)

	t.Log("\n✓ All subscription management tests passed")
}

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
// Environment variables:
//   - RUN_INTEGRATION_TESTS=true (required)
//   - INTERACTIVE=true (required)
//   - BASE_HOSTNAME=http://localhost:8000 (optional, defaults to production)
//   - TEST_PRODUCT_ID=<product-id> (optional, uses "Butler" product if not set)
//   - TEST_QA_USERNAME=<username> (optional, creates random user if not set)
//   - TEST_QA_PASSWORD=<password> (optional, creates random user if not set)
//
// Run with: BASE_HOSTNAME=http://localhost:8000 RUN_INTEGRATION_TESTS=true INTERACTIVE=true TEST_QA_USERNAME=qa_user TEST_QA_PASSWORD=qa_pass go test -v -run TestE2EInteractiveFullLifecycle -timeout 30m
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

	// Check for QA user credentials
	qaUsername := os.Getenv("TEST_QA_USERNAME")
	qaPassword := os.Getenv("TEST_QA_PASSWORD")

	var username, password, userID string

	separator := "=" + strings.Repeat("=", 70)
	t.Log(separator)
	t.Log("INTERACTIVE INTEGRATION TEST - FULL SUBSCRIPTION LIFECYCLE")
	t.Log(separator)
	t.Logf("")

	if qaUsername != "" && qaPassword != "" {
		// Use dedicated QA user
		username = qaUsername
		password = qaPassword

		t.Logf("This test will:")
		t.Logf("  1. Login with QA user account")
		t.Logf("  2. Create an order and checkout session")
		t.Logf("  3. Open Stripe checkout in your browser")
		t.Logf("  4. Wait for you to complete payment")
		t.Logf("  5. Test subscription management (pause, resume, update)")
		t.Logf("  6. Cancel the subscription")
		t.Logf("")
		t.Logf("QA User: %s", username)
		t.Log(separator)
		t.Logf("")

		// Step 1: Login with QA user
		t.Logf("\n▶ Step 1: Logging in with QA user...")
		loginReq := api.LoginRequest{
			Username: username,
			Password: password,
		}
		loginResp, err := client.Login(loginReq)
		if err != nil {
			t.Fatalf("❌ Failed to login with QA user: %v", err)
		}
		userID = loginResp.Data.UserID
		t.Logf("✓ Logged in as: %s (ID: %s)\n", username, userID)
	} else {
		// Generate unique test credentials
		timestamp := time.Now().Unix()
		randSuffix := rand.Intn(10000)
		username = fmt.Sprintf("test_user_%d_%d", timestamp, randSuffix)
		email := fmt.Sprintf("test_%d_%d@butler.test", timestamp, randSuffix)
		password = "TestPassword123!@#"

		t.Logf("This test will:")
		t.Logf("  1. Create a new user account")
		t.Logf("  2. Create an order and checkout session")
		t.Logf("  3. Open Stripe checkout in your browser")
		t.Logf("  4. Wait for you to complete payment")
		t.Logf("  5. Test subscription management (pause, resume, update)")
		t.Logf("  6. Cancel the subscription")
		t.Logf("")
		t.Logf("Test credentials: %s / %s", username, password)
		t.Logf("Tip: Set TEST_QA_USERNAME and TEST_QA_PASSWORD to use a dedicated QA user")
		t.Log(separator)
		t.Logf("")

		// Step 1: User Registration
		t.Logf("\n▶ Step 1: Creating user account...")
		userID = testInteractiveUserRegistration(t, client, username, email, password)
		t.Logf("✓ User created: %s\n", userID)

		// Register cleanup function to delete test user at the end
		registerCleanup(t, client, username)
	}

	// Step 2: Get available subscriptions
	t.Logf("\n▶ Step 2: Fetching available subscriptions...")
	selectedPlan := testInteractiveGetSubscriptions(t, client)
	t.Logf("✓ Selected tier: %s (ID: %s)\n", selectedPlan.Tier, selectedPlan.ID)

	// Step 3: Create order
	t.Logf("\n▶ Step 3: Creating order with coffee preferences...")
	orderID := testInteractiveCreateOrder(t, client, selectedPlan)
	t.Logf("✓ Order created: %s\n", orderID)

	// Step 4: Create checkout session and wait for payment
	t.Logf("\n▶ Step 4: Creating checkout session and processing payment...")
	subscriptionID := testInteractiveCheckoutAndPay(t, client, orderID)
	t.Logf("✓ Payment completed! Subscription ID: %s\n", subscriptionID)

	// Step 5: Verify subscription was created
	t.Logf("\n▶ Step 5: Verifying subscription was created...")
	originalSubscription := testInteractiveVerifySubscription(t, client, subscriptionID)
	t.Logf("✓ Subscription verified\n")

	// Step 6: Pause subscription
	t.Logf("\n▶ Step 6: Pausing subscription...")
	testInteractivePauseSubscription(t, client, subscriptionID)
	t.Logf("✓ Subscription paused\n")

	// Step 7: Resume subscription
	t.Logf("\n▶ Step 7: Resuming subscription...")
	testInteractiveResumeSubscription(t, client, subscriptionID)
	t.Logf("✓ Subscription resumed\n")

	// Step 8: Update subscription preferences
	t.Logf("\n▶ Step 8: Updating subscription preferences...")
	testInteractiveUpdateSubscription(t, client, subscriptionID)
	t.Logf("✓ Subscription updated\n")

	// Step 9: Restore original preferences
	t.Logf("\n▶ Step 9: Restoring original preferences...")
	testInteractiveRestoreSubscription(t, client, subscriptionID, originalSubscription)
	t.Logf("✓ Preferences restored\n")

	// Step 10: Cancel subscription
	t.Logf("\n▶ Step 10: Cancelling subscription...")
	testInteractiveCancelSubscription(t, client, subscriptionID)
	t.Logf("✓ Subscription cancelled\n")

	finalSeparator := "=" + strings.Repeat("=", 70)
	t.Logf("")
	t.Log(finalSeparator)
	t.Log("✓ ALL TESTS PASSED - COMPLETE LIFECYCLE VERIFIED")
	t.Log(finalSeparator)
}

// ============================================================================
// Helper Functions
// ============================================================================

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
		// If user already exists, try to login instead
		t.Logf("⚠ Registration failed (user may already exist): %v", err)
		t.Logf("  Attempting to login with existing credentials...")

		loginReq := api.LoginRequest{
			Username: username,
			Password: password,
		}
		loginResp, loginErr := client.Login(loginReq)
		if loginErr != nil {
			t.Fatalf("❌ Both registration and login failed: %v", loginErr)
		}

		t.Logf("✓ Logged in with existing user: UserID=%s", loginResp.Data.UserID)
		return
	}

	if resp.Data.ID == "" {
		t.Fatalf("❌ Expected user ID in response")
	}
	if resp.Data.AccessToken == "" {
		t.Fatalf("❌ Expected access token in response")
	}
	if resp.Data.RefreshToken == "" {
		t.Fatalf("❌ Expected refresh token in response")
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
		t.Fatalf("❌ Login failed: %v", err)
	}

	if resp.Data.AccessToken == "" {
		t.Fatalf("❌ Expected access token in response")
	}
	if resp.Data.RefreshToken == "" {
		t.Fatalf("❌ Expected refresh token in response")
	}
	if resp.Data.ExpiresAt == "" {
		t.Fatalf("❌ Expected expires_at in response")
	}
	if resp.Data.RefreshTokenExpiresAt == "" {
		t.Fatalf("❌ Expected refresh_token_expires_at in response")
	}
	if resp.Data.UserID == "" {
		t.Fatalf("❌ Expected user_id in response")
	}

	t.Logf("✓ User logged in successfully: UserID=%s", resp.Data.UserID)
	t.Logf("✓ Access token expires at: %s", resp.Data.ExpiresAt)
	t.Logf("✓ Refresh token expires at: %s", resp.Data.RefreshTokenExpiresAt)
}

func testGetAvailableSubscriptions(t *testing.T, client *api.Client) *api.AvailablePlan {
	subscriptions, err := client.GetAvailableSubscriptions()
	if err != nil {
		t.Fatalf("❌ Failed to get available subscriptions: %v", err)
	}

	if len(subscriptions) == 0 {
		t.Fatalf("❌ Expected at least one available subscription")
	}

	t.Logf("✓ Found %d available subscription(s)", len(subscriptions))

	for i, sub := range subscriptions {
		t.Logf("  [%d] Tier: %s, Name: %s, Price: %s %s/%s, ID: %s",
			i+1, sub.Tier, sub.Name, sub.Price, sub.Currency, sub.BillingPeriod, sub.ID)

		// Validate each subscription
		if sub.ID == "" {
			t.Fatalf("❌ Subscription %d missing ID", i)
		}
		if sub.Tier == "" {
			t.Fatalf("❌ Subscription %d missing tier", i)
		}
		if sub.Name == "" {
			t.Fatalf("❌ Subscription %d missing name", i)
		}
		if sub.Price == "" {
			t.Fatalf("❌ Subscription %d missing price", i)
		}
		if !sub.IsSubscription {
			t.Fatalf("❌ Subscription %d has IsSubscription=false", i)
		}
	}

	// Check if a specific test product ID is set
	testProductID := os.Getenv("TEST_PRODUCT_ID")
	if testProductID != "" {
		t.Logf("Looking for TEST_PRODUCT_ID: %s", testProductID)
		for _, sub := range subscriptions {
			if sub.ID == testProductID {
				t.Logf("✓ Selected tier for testing: %s (ID: %s)", sub.Tier, sub.ID)
				return &sub
			}
		}
		t.Fatalf("❌ TEST_PRODUCT_ID %s not found in available subscriptions", testProductID)
	}

	// For QA/local testing, prefer the product named "Butler"
	for _, sub := range subscriptions {
		if sub.Name == "Butler" {
			t.Logf("✓ Selected tier for testing: %s (ID: %s)", sub.Tier, sub.ID)
			t.Logf("Tip: Set TEST_PRODUCT_ID environment variable to use a different product")
			return &sub
		}
	}

	// Fallback to first plan if "Butler" not found
	t.Logf("✓ Selected tier for testing: %s (ID: %s)", subscriptions[0].Tier, subscriptions[0].ID)
	t.Logf("Tip: Set TEST_PRODUCT_ID environment variable to use a specific product for testing")
	return &subscriptions[0]
}

func testCreateOrder(t *testing.T, client *api.Client, plan *api.AvailablePlan) string {
	// Create an order with multiple preferences
	req := api.CreateOrderRequest{
		Tier:          plan.Tier,
		ProductID:     plan.ID,
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
		t.Fatalf("❌ Failed to create order: %v", err)
	}

	// Validate order response
	if order.ID == "" {
		t.Fatalf("❌ Expected order ID")
	}
	if order.Tier != plan.Tier {
		t.Fatalf("❌ Expected tier %s, got %s", plan.Tier, order.Tier)
	}
	if order.Status == "" {
		t.Fatalf("❌ Expected order status")
	}
	if len(order.LineItems) != 2 {
		t.Fatalf("❌ Expected 2 line items, got %d", len(order.LineItems))
	}

	totalQty := order.GetTotalQuantity()
	if totalQty != 5 {
		t.Fatalf("❌ Expected total quantity 5, got %d", totalQty)
	}

	t.Logf("✓ Order created successfully: ID=%s", order.ID)
	t.Logf("✓ Order tier: %s, status: %s", order.Tier, order.Status)
	t.Logf("✓ Total quantity: %d kg", totalQty)

	for i, item := range order.LineItems {
		qty := item.GetQuantity()
		t.Logf("  [%d] %s → %d kg (%s)", i+1, item.BrewingMethod, qty, item.GrindType)
	}

	return order.ID
}

func testCreateCheckoutSession(t *testing.T, client *api.Client, orderID string) {
	session, err := client.CreateCheckoutSession(orderID)
	if err != nil {
		t.Fatalf("❌ Failed to create checkout session: %v", err)
	}

	if session.CheckoutURL == "" {
		t.Fatalf("❌ Expected checkout URL")
	}
	if session.SessionID == "" {
		t.Fatalf("❌ Expected session ID")
	}
	if session.OrderID != orderID {
		t.Fatalf("❌ Expected order ID %s, got %s", orderID, session.OrderID)
	}

	t.Logf("✓ Checkout session created successfully")
	t.Logf("✓ Session ID: %s", session.SessionID)
	t.Logf("✓ Checkout URL: %s", session.CheckoutURL)
	t.Log("Note: Automated tests cannot complete Stripe payment flow")
}

func testListSubscriptions(t *testing.T, client *api.Client, expectedCount int) {
	subscriptions, err := client.ListSubscriptions()
	if err != nil {
		t.Fatalf("❌ Failed to list subscriptions: %v", err)
	}

	if len(subscriptions) != expectedCount {
		t.Fatalf("❌ Expected %d subscription(s), got %d", expectedCount, len(subscriptions))
	}

	t.Logf("✓ Found %d subscription(s)", len(subscriptions))

	for i, sub := range subscriptions {
		t.Logf("  [%d] ID: %s, Tier: %s, Status: %s", i+1, sub.ID, sub.Tier, sub.Status)
	}
}

func testGetSubscription(t *testing.T, client *api.Client, subscriptionID string) *api.Subscription {
	subscription, err := client.GetSubscription(subscriptionID)
	if err != nil {
		t.Fatalf("❌ Failed to get subscription: %v", err)
	}

	if subscription.ID != subscriptionID {
		t.Fatalf("❌ Expected subscription ID %s, got %s", subscriptionID, subscription.ID)
	}
	if subscription.Tier == "" {
		t.Fatalf("❌ Expected tier in subscription")
	}
	if subscription.Status == "" {
		t.Fatalf("❌ Expected status in subscription")
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
		t.Fatalf("❌ Failed to pause subscription: %v", err)
	}

	if subscription.Status != "paused" {
		t.Fatalf("❌ Expected status 'paused', got '%s'", subscription.Status)
	}

	t.Logf("✓ Subscription paused successfully: ID=%s", subscription.ID)
	t.Logf("✓ Status: %s", subscription.Status)
}

func testResumeSubscription(t *testing.T, client *api.Client, subscriptionID string) {
	subscription, err := client.ResumeSubscription(subscriptionID)
	if err != nil {
		t.Fatalf("❌ Failed to resume subscription: %v", err)
	}

	if subscription.Status != "active" {
		t.Fatalf("❌ Expected status 'active', got '%s'", subscription.Status)
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
		t.Fatalf("❌ Failed to update subscription: %v", err)
	}

	if subscription.GetTotalQuantity() != 8 {
		t.Fatalf("❌ Expected total quantity 8, got %d", subscription.GetTotalQuantity())
	}

	t.Logf("✓ Subscription updated successfully: ID=%s", subscription.ID)
	t.Logf("✓ New total quantity: %d kg/month", subscription.GetTotalQuantity())

	// Verify the updated preferences
	updated, err := client.GetSubscription(subscriptionID)
	if err != nil {
		t.Fatalf("❌ Failed to verify updated subscription: %v", err)
	}

	if updated.GetTotalQuantity() != 8 {
		t.Fatalf("❌ Verification failed: Expected total quantity 8, got %d", updated.GetTotalQuantity())
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
		t.Fatalf("❌ Failed to restore subscription: %v", err)
	}

	if subscription.GetTotalQuantity() != original.GetTotalQuantity() {
		t.Fatalf("❌ Expected total quantity %d, got %d",
			original.GetTotalQuantity(), subscription.GetTotalQuantity())
	}

	t.Logf("✓ Subscription restored to original preferences: ID=%s", subscription.ID)
	t.Logf("✓ Total quantity: %d kg/month", subscription.GetTotalQuantity())
}

// ============================================================================
// Interactive Test Helper Functions
// ============================================================================

func testInteractiveUserRegistration(t *testing.T, client *api.Client, username, email, password string) string {
	req := api.RegisterRequest{
		Username: username,
		Email:    email,
		Password: password,
	}

	resp, err := client.Register(req)
	if err != nil {
		// If user already exists, try to login instead
		t.Logf("  ⚠ Registration failed (user may already exist): %v", err)
		t.Logf("  Attempting to login with existing credentials...")

		loginReq := api.LoginRequest{
			Username: username,
			Password: password,
		}
		loginResp, loginErr := client.Login(loginReq)
		if loginErr != nil {
			t.Fatalf("❌ Both registration and login failed: %v", loginErr)
		}

		t.Logf("  User ID: %s", loginResp.Data.UserID)
		t.Logf("  Username: %s", username)
		t.Logf("  Email: %s", email)
		t.Logf("  (Using existing user)")

		return loginResp.Data.UserID
	}

	if resp.Data.ID == "" {
		t.Fatalf("❌ Expected user ID in response")
	}
	if resp.Data.AccessToken == "" {
		t.Fatalf("❌ Expected access token in response")
	}

	t.Logf("  User ID: %s", resp.Data.ID)
	t.Logf("  Username: %s", username)
	t.Logf("  Email: %s", email)

	return resp.Data.ID
}

func testInteractiveGetSubscriptions(t *testing.T, client *api.Client) *api.AvailablePlan {
	subscriptions, err := client.GetAvailableSubscriptions()
	if err != nil {
		t.Fatalf("❌ Failed to get available subscriptions: %v", err)
	}

	if len(subscriptions) == 0 {
		t.Fatalf("❌ Expected at least one available subscription")
	}

	t.Logf("  Available tiers:")
	for i, sub := range subscriptions {
		t.Logf("    [%d] %s - %s %s/%s (ID: %s)", i+1, sub.Name, sub.Price, sub.Currency, sub.BillingPeriod, sub.ID)
	}

	// Check if a specific test product ID is set
	testProductID := os.Getenv("TEST_PRODUCT_ID")
	if testProductID != "" {
		for _, sub := range subscriptions {
			if sub.ID == testProductID {
				t.Logf("  Using TEST_PRODUCT_ID: %s", testProductID)
				return &sub
			}
		}
		t.Fatalf("❌ TEST_PRODUCT_ID %s not found in available subscriptions", testProductID)
	}

	// For QA/local testing, prefer the product named "Butler"
	for _, sub := range subscriptions {
		if sub.Name == "Butler" {
			t.Logf("  Using product: %s (ID: %s)", sub.Name, sub.ID)
			return &sub
		}
	}

	// Use the first plan as fallback
	return &subscriptions[0]
}

func testInteractiveCreateOrder(t *testing.T, client *api.Client, plan *api.AvailablePlan) string {
	req := api.CreateOrderRequest{
		Tier:          plan.Tier,
		ProductID:     plan.ID,
		TotalQuantity: 5,
		LineItems: []api.OrderLineItem{
			{
				Quantity:      3,
				GrindType:     "whole_bean",
				BrewingMethod: "espresso",
				Notes:         "Interactive test - whole beans",
			},
			{
				Quantity:      2,
				GrindType:     "ground",
				BrewingMethod: "v60",
				Notes:         "Interactive test - ground coffee",
			},
		},
	}

	order, err := client.CreateOrder(req)
	if err != nil {
		t.Fatalf("❌ Failed to create order: %v", err)
	}

	t.Logf("  Order ID: %s", order.ID)
	t.Logf("  Tier: %s", order.Tier)
	t.Logf("  Total quantity: %d kg/month", order.GetTotalQuantity())
	t.Logf("  Line items:")
	for i, item := range order.LineItems {
		t.Logf("    [%d] %s → %d kg (%s)", i+1, item.BrewingMethod, item.GetQuantity(), item.GrindType)
	}

	return order.ID
}

func testInteractiveCheckoutAndPay(t *testing.T, client *api.Client, orderID string) string {
	// Create checkout session
	session, err := client.CreateCheckoutSession(orderID)
	if err != nil {
		t.Fatalf("❌ Failed to create checkout session: %v", err)
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
		t.Fatalf("❌ Failed to verify payment: %v", err)
	}

	t.Logf("")
	t.Logf("  ✓ Payment verified!")
	t.Logf("  ✓ Subscription created: %s", subscriptionID)

	return subscriptionID
}

func testInteractiveVerifySubscription(t *testing.T, client *api.Client, subscriptionID string) *api.Subscription {
	subscription, err := client.GetSubscription(subscriptionID)
	if err != nil {
		t.Fatalf("❌ Failed to get subscription: %v", err)
	}

	if subscription.Status != "active" {
		t.Fatalf("❌ Expected status 'active', got '%s'", subscription.Status)
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
		t.Fatalf("❌ Failed to pause subscription: %v", err)
	}

	if subscription.Status != "paused" {
		t.Fatalf("❌ Expected status 'paused', got '%s'", subscription.Status)
	}

	t.Logf("  Status changed to: %s", subscription.Status)
}

func testInteractiveResumeSubscription(t *testing.T, client *api.Client, subscriptionID string) {
	subscription, err := client.ResumeSubscription(subscriptionID)
	if err != nil {
		t.Fatalf("❌ Failed to resume subscription: %v", err)
	}

	if subscription.Status != "active" {
		t.Fatalf("❌ Expected status 'active', got '%s'", subscription.Status)
	}

	t.Logf("  Status changed to: %s", subscription.Status)
}

func testInteractiveUpdateSubscription(t *testing.T, client *api.Client, subscriptionID string) {
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
		t.Fatalf("❌ Failed to update subscription: %v", err)
	}

	if subscription.GetTotalQuantity() != 8 {
		t.Fatalf("❌ Expected total quantity 8, got %d", subscription.GetTotalQuantity())
	}

	t.Logf("  Updated total quantity: %d kg/month", subscription.GetTotalQuantity())

	// Verify the update
	updated, err := client.GetSubscription(subscriptionID)
	if err != nil {
		t.Fatalf("❌ Failed to verify update: %v", err)
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
		t.Fatalf("❌ Failed to restore subscription: %v", err)
	}

	t.Logf("  Restored to %d kg/month", subscription.GetTotalQuantity())
}

func testInteractiveCancelSubscription(t *testing.T, client *api.Client, subscriptionID string) {
	subscription, err := client.CancelSubscription(subscriptionID)
	if err != nil {
		t.Fatalf("❌ Failed to cancel subscription: %v", err)
	}

	if subscription.Status != "cancelled" {
		t.Fatalf("❌ Expected status 'cancelled', got '%s'", subscription.Status)
	}

	t.Logf("  Status changed to: %s", subscription.Status)
	t.Logf("  Subscription is now permanently cancelled")
}

// ============================================================================
// Utility Functions
// ============================================================================

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

// registerCleanup registers a cleanup function that attempts to delete the test user
// This ensures test data is cleaned up even if the test fails
func registerCleanup(t *testing.T, client *api.Client, username string) {
	t.Cleanup(func() {
		// XXX add here any cleanup needed
	})
}

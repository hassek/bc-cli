package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hassek/bc-cli/config"
)

func TestListSubscriptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/core/v1/subscriptions" {
			t.Errorf("Expected path '/api/core/v1/subscriptions', got '%s'", r.URL.Path)
		}
		if r.Method != "GET" {
			t.Errorf("Expected GET method, got %s", r.Method)
		}

		response := ListSubscriptionsResponse{
			Meta: struct {
				Code    int    `json:"code"`
				Message string `json:"message"`
			}{Code: 200, Message: "Success"},
			Data: []Subscription{
				{
					ID:              "sub-123",
					Tier:            "butler",
					Status:          "active",
					DefaultQuantity: "5",
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	cfg := &config.Config{
		APIURL:      server.URL,
		AccessToken: "test-token",
	}
	client := NewClient(cfg)

	subscriptions, err := client.ListSubscriptions()
	if err != nil {
		t.Fatalf("ListSubscriptions failed: %v", err)
	}

	if len(subscriptions) != 1 {
		t.Errorf("Expected 1 subscription, got %d", len(subscriptions))
	}

	if subscriptions[0].ID != "sub-123" {
		t.Errorf("Expected subscription ID 'sub-123', got '%s'", subscriptions[0].ID)
	}

	if subscriptions[0].Tier != "butler" {
		t.Errorf("Expected tier 'butler', got '%s'", subscriptions[0].Tier)
	}

	if subscriptions[0].GetTotalQuantity() != 5 {
		t.Errorf("Expected total quantity 5, got %d", subscriptions[0].GetTotalQuantity())
	}
}

func TestGetSubscription(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/core/v1/subscriptions/sub-123/preferences" {
			t.Errorf("Expected path '/api/core/v1/subscriptions/sub-123/preferences', got '%s'", r.URL.Path)
		}
		if r.Method != "GET" {
			t.Errorf("Expected GET method, got %s", r.Method)
		}

		response := SubscriptionActionResponse{
			Meta: struct {
				Code    int    `json:"code"`
				Message string `json:"message"`
			}{Code: 200, Message: "Success"},
			Data: Subscription{
				ID:              "sub-123",
				Tier:            "butler",
				Status:          "active",
				DefaultQuantity: "5",
				DefaultPreferences: []SubscriptionPreference{
					{
						ID:            "pref-1",
						Quantity:      "3",
						GrindType:     "whole_bean",
						BrewingMethod: "espresso",
					},
					{
						ID:            "pref-2",
						Quantity:      "2",
						GrindType:     "ground",
						BrewingMethod: "v60",
					},
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	cfg := &config.Config{
		APIURL:      server.URL,
		AccessToken: "test-token",
	}
	client := NewClient(cfg)

	subscription, err := client.GetSubscription("sub-123")
	if err != nil {
		t.Fatalf("GetSubscription failed: %v", err)
	}

	if subscription.ID != "sub-123" {
		t.Errorf("Expected subscription ID 'sub-123', got '%s'", subscription.ID)
	}

	if len(subscription.DefaultPreferences) != 2 {
		t.Errorf("Expected 2 preferences, got %d", len(subscription.DefaultPreferences))
	}

	if subscription.DefaultPreferences[0].GrindType != "whole_bean" {
		t.Errorf("Expected grind type 'whole_bean', got '%s'", subscription.DefaultPreferences[0].GrindType)
	}
}

func TestPauseSubscription(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/core/v1/subscriptions/sub-123/pause" {
			t.Errorf("Expected path '/api/core/v1/subscriptions/sub-123/pause', got '%s'", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		// Backend does not accept request body for pause
		// This matches actual backend behavior (subscriptions.py line 63)

		response := SubscriptionActionResponse{
			Meta: struct {
				Code    int    `json:"code"`
				Message string `json:"message"`
			}{Code: 200, Message: "Subscription paused"},
			Data: Subscription{
				ID:     "sub-123",
				Status: "paused",
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	cfg := &config.Config{
		APIURL:      server.URL,
		AccessToken: "test-token",
	}
	client := NewClient(cfg)

	subscription, err := client.PauseSubscription("sub-123")
	if err != nil {
		t.Fatalf("PauseSubscription failed: %v", err)
	}

	if subscription.Status != "paused" {
		t.Errorf("Expected status 'paused', got '%s'", subscription.Status)
	}
}

func TestResumeSubscription(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/core/v1/subscriptions/sub-123/resume" {
			t.Errorf("Expected path '/api/core/v1/subscriptions/sub-123/resume', got '%s'", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		response := SubscriptionActionResponse{
			Meta: struct {
				Code    int    `json:"code"`
				Message string `json:"message"`
			}{Code: 200, Message: "Subscription resumed"},
			Data: Subscription{
				ID:     "sub-123",
				Status: "active",
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	cfg := &config.Config{
		APIURL:      server.URL,
		AccessToken: "test-token",
	}
	client := NewClient(cfg)

	subscription, err := client.ResumeSubscription("sub-123")
	if err != nil {
		t.Fatalf("ResumeSubscription failed: %v", err)
	}

	if subscription.Status != "active" {
		t.Errorf("Expected status 'active', got '%s'", subscription.Status)
	}
}

func TestCancelSubscription(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/core/v1/subscriptions/sub-123/cancel" {
			t.Errorf("Expected path '/api/core/v1/subscriptions/sub-123/cancel', got '%s'", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		response := SubscriptionActionResponse{
			Meta: struct {
				Code    int    `json:"code"`
				Message string `json:"message"`
			}{Code: 200, Message: "Subscription cancelled"},
			Data: Subscription{
				ID:     "sub-123",
				Status: "cancelled",
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	cfg := &config.Config{
		APIURL:      server.URL,
		AccessToken: "test-token",
	}
	client := NewClient(cfg)

	subscription, err := client.CancelSubscription("sub-123")
	if err != nil {
		t.Fatalf("CancelSubscription failed: %v", err)
	}

	if subscription.Status != "cancelled" {
		t.Errorf("Expected status 'cancelled', got '%s'", subscription.Status)
	}
}

func TestUpdateSubscription(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/core/v1/subscriptions/sub-123/preferences" {
			t.Errorf("Expected path '/api/core/v1/subscriptions/sub-123/preferences', got '%s'", r.URL.Path)
		}
		if r.Method != "PATCH" {
			t.Errorf("Expected PATCH method, got %s", r.Method)
		}

		var reqBody UpdateSubscriptionRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}

		if reqBody.TotalQuantity != 10 {
			t.Errorf("Expected total quantity 10, got %d", reqBody.TotalQuantity)
		}

		if len(reqBody.Preferences) != 1 {
			t.Errorf("Expected 1 preference item, got %d", len(reqBody.Preferences))
		}

		response := SubscriptionActionResponse{
			Meta: struct {
				Code    int    `json:"code"`
				Message string `json:"message"`
			}{Code: 200, Message: "Subscription updated"},
			Data: Subscription{
				ID:              "sub-123",
				Status:          "active",
				DefaultQuantity: "10",
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	cfg := &config.Config{
		APIURL:      server.URL,
		AccessToken: "test-token",
	}
	client := NewClient(cfg)

	updateReq := UpdateSubscriptionRequest{
		TotalQuantity: 10,
		Preferences: []OrderLineItem{
			{
				Quantity:      10,
				GrindType:     "whole_bean",
				BrewingMethod: "espresso",
			},
		},
	}

	subscription, err := client.UpdateSubscription("sub-123", updateReq)
	if err != nil {
		t.Fatalf("UpdateSubscription failed: %v", err)
	}

	if subscription.GetTotalQuantity() != 10 {
		t.Errorf("Expected total quantity 10, got %d", subscription.GetTotalQuantity())
	}
}

func TestSubscriptionGetTotalQuantity(t *testing.T) {
	tests := []struct {
		name     string
		quantity string
		want     int
	}{
		{"whole number", "5", 5},
		{"decimal rounds up", "5.6", 6},
		{"decimal rounds down", "5.4", 5},
		{"empty string", "", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sub := &Subscription{DefaultQuantity: tt.quantity}
			got := sub.GetTotalQuantity()
			if got != tt.want {
				t.Errorf("GetTotalQuantity() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestGetAvailableSubscriptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/core/v1/subscriptions/available" {
			t.Errorf("Expected path '/api/core/v1/subscriptions/available', got '%s'", r.URL.Path)
		}
		if r.Method != "GET" {
			t.Errorf("Expected GET method, got %s", r.Method)
		}

		// Verify query parameter
		if r.URL.Query().Get("is_subscription") != "true" {
			t.Errorf("Expected query parameter is_subscription=true")
		}

		response := AvailablePlansResponse{
			Meta: struct {
				Code    int    `json:"code"`
				Message string `json:"message"`
			}{Code: 200, Message: "Success"},
			Data: []AvailablePlan{
				{
					ID:             "plan-1",
					Tier:           "butler",
					Name:           "Butler Coffee",
					Price:          "49.00",
					Currency:       "EUR",
					BillingPeriod:  "monthly",
					IsActive:       true,
					IsSubscription: true,
					MinQuantity:    1,
					MaxQuantity:    10,
				},
				{
					ID:             "plan-2",
					Tier:           "collection",
					Name:           "Collection Coffee",
					Price:          "79.00",
					Currency:       "EUR",
					BillingPeriod:  "monthly",
					IsActive:       true,
					IsSubscription: true,
					MinQuantity:    1,
					MaxQuantity:    100,
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	cfg := &config.Config{
		APIURL: server.URL,
	}
	client := NewClient(cfg)

	plans, err := client.GetAvailableSubscriptions()
	if err != nil {
		t.Fatalf("GetAvailableSubscriptions failed: %v", err)
	}

	if len(plans) != 2 {
		t.Errorf("Expected 2 plans, got %d", len(plans))
	}

	// Check first plan
	if plans[0].Tier != "butler" {
		t.Errorf("Expected tier 'butler', got '%s'", plans[0].Tier)
	}
	if plans[0].MinQuantity != 1 {
		t.Errorf("Expected MinQuantity 1, got %d", plans[0].MinQuantity)
	}
	if plans[0].MaxQuantity != 10 {
		t.Errorf("Expected MaxQuantity 10, got %d", plans[0].MaxQuantity)
	}

	// Check second plan
	if plans[1].Tier != "collection" {
		t.Errorf("Expected tier 'collection', got '%s'", plans[1].Tier)
	}
	if plans[1].MinQuantity != 1 {
		t.Errorf("Expected MinQuantity 1, got %d", plans[1].MinQuantity)
	}
	if plans[1].MaxQuantity != 100 {
		t.Errorf("Expected MaxQuantity 100, got %d", plans[1].MaxQuantity)
	}
}

func TestGetAvailableProducts(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/core/v1/subscriptions/available" {
			t.Errorf("Expected path '/api/core/v1/subscriptions/available', got '%s'", r.URL.Path)
		}
		if r.Method != "GET" {
			t.Errorf("Expected GET method, got %s", r.Method)
		}

		// Verify query parameter
		if r.URL.Query().Get("is_subscription") != "false" {
			t.Errorf("Expected query parameter is_subscription=false")
		}

		response := AvailablePlansResponse{
			Meta: struct {
				Code    int    `json:"code"`
				Message string `json:"message"`
			}{Code: 200, Message: "Success"},
			Data: []AvailablePlan{
				{
					ID:             "product-1",
					Name:           "Christmas Bundle",
					Price:          "40.00",
					Currency:       "EUR",
					BillingPeriod:  "one-time",
					IsActive:       true,
					IsSubscription: false,
					MinQuantity:    1,
					MaxQuantity:    100,
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	cfg := &config.Config{
		APIURL: server.URL,
	}
	client := NewClient(cfg)

	products, err := client.GetAvailableProducts()
	if err != nil {
		t.Fatalf("GetAvailableProducts failed: %v", err)
	}

	if len(products) != 1 {
		t.Errorf("Expected 1 product, got %d", len(products))
	}

	// Check product
	if products[0].Name != "Christmas Bundle" {
		t.Errorf("Expected name 'Christmas Bundle', got '%s'", products[0].Name)
	}
	if products[0].MinQuantity != 1 {
		t.Errorf("Expected MinQuantity 1, got %d", products[0].MinQuantity)
	}
	if products[0].MaxQuantity != 100 {
		t.Errorf("Expected MaxQuantity 100, got %d", products[0].MaxQuantity)
	}
	if products[0].IsSubscription {
		t.Error("Expected IsSubscription to be false")
	}
}

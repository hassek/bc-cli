package cmd

import (
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/hassek/bc-cli/api"
	"github.com/hassek/bc-cli/cmd/order"
	"github.com/hassek/bc-cli/config"
	"github.com/hassek/bc-cli/templates"
	"github.com/hassek/bc-cli/tui/models"
	"github.com/hassek/bc-cli/tui/prompts"
	"github.com/hassek/bc-cli/utils"
	"github.com/spf13/cobra"
)

var subscriptionsCmd = &cobra.Command{
	Use:   "subscriptions",
	Short: "Manage your Butler Coffee subscriptions",
	Long:  `View your active subscriptions and browse available tiers interactively.`,
	RunE:  runSubscriptions,
}

func init() {
	rootCmd.AddCommand(subscriptionsCmd)
}

func runSubscriptions(cmd *cobra.Command, args []string) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	client := api.NewClient(cfg)

	// Get available subscriptions
	available, err := client.GetAvailableSubscriptions()
	if err != nil {
		return fmt.Errorf("failed to get available subscriptions: %w", err)
	}

	if len(available) == 0 {
		fmt.Println("No subscription tiers available at this time.")
		return nil
	}

	// Use new subscription picker with duck animation
	selectedSub, err := models.PickSubscription(available)
	if err != nil {
		fmt.Println("\nExiting...")
		return nil
	}

	// User cancelled or selected exit
	if selectedSub == nil {
		return nil
	}

	displaySubscriptionDetails(*selectedSub, api.Subscription{}, cfg.IsAuthenticated())

	// Ask if user wants to subscribe (if authenticated)
	if cfg.IsAuthenticated() {
		fmt.Println()
		confirmed, err := prompts.PromptConfirm(fmt.Sprintf("Would you like to subscribe to %s now", selectedSub.Name))
		if err == nil && confirmed {
			// User wants to subscribe - start order configuration flow
			return createOrderAndSubscribe(cfg, client, *selectedSub)
		}
	} else if !cfg.IsAuthenticated() {
		fmt.Println("\nPlease login first to subscribe:")
		fmt.Println("  bc-cli login")
	}

	return nil
}

func displaySubscriptionDetails(sub api.AvailableSubscription, activeSub api.Subscription, isAuthenticated bool) {
	type activeSubData struct {
		ID        string
		Status    string
		StartedAt string
		ExpiresAt string
	}

	var activeData activeSubData
	if activeSub.ID != "" {
		activeData.ID = activeSub.ID
		activeData.Status = activeSub.Status
		if activeSub.StartedAt != nil {
			activeData.StartedAt = utils.FormatTimestamp(*activeSub.StartedAt)
		}
		if activeSub.ExpiresAt != nil {
			activeData.ExpiresAt = utils.FormatTimestamp(*activeSub.ExpiresAt)
		}
	}

	// Use viewport for scrollable display
	if err := templates.RenderInViewport(sub.Name, templates.SubscriptionDetailsTemplate, struct {
		Name          string
		Currency      string
		Price         string
		BillingPeriod string
		Description   string
		ActiveSub     activeSubData
	}{
		Name:          sub.Name,
		Currency:      sub.Currency,
		Price:         sub.Price,
		BillingPeriod: sub.BillingPeriod,
		Description:   sub.Description,
		ActiveSub:     activeData,
	}); err != nil {
		fmt.Printf("Error rendering template: %v\n", err)
	}
}

func createOrderAndSubscribe(cfg *config.Config, client *api.Client, tier api.AvailableSubscription) error {
	if !cfg.IsAuthenticated() {
		return fmt.Errorf("you must be logged in to subscribe. Please run 'bc-cli login' first")
	}

	// Use min/max quantity from backend, fallback to config defaults if not set
	minQty := tier.MinQuantity
	maxQty := tier.MaxQuantity
	if minQty == 0 {
		minQty = cfg.MinQuantity
	}
	if maxQty == 0 {
		maxQty = cfg.MaxQuantity
	}

	if err := templates.RenderToStdout(templates.OrderConfigIntroTemplate, struct {
		MinQuantity int
		MaxQuantity int
	}{
		MinQuantity: minQty,
		MaxQuantity: maxQty,
	}); err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}

	totalQuantity, err := prompts.PromptQuantityInt("Total quantity per month", minQty, maxQty, minQty)
	if err != nil {
		return err
	}

	fmt.Printf("\n✓ Total: %d per month\n", totalQuantity)

	var lineItems []api.OrderLineItem

	// Step 2: Ask if they want to split or keep uniform (skip if quantity is 1)
	if totalQuantity == 1 {
		// Only 1 unit - cannot split, go straight to uniform order
		lineItems, err = order.ConfigureUniformOrder(totalQuantity)
		if err != nil {
			return err
		}
	} else {
		// Multiple units - ask if they want to split
		if err := templates.RenderToStdout(templates.OrderSplitIntroTemplate, nil); err != nil {
			return fmt.Errorf("failed to render template: %w", err)
		}

		wantsSplit, err := prompts.PromptConfirm("Would you like different grind methods?")
		if err != nil {
			return err
		}

		if !wantsSplit {
			// Simple flow - all coffee the same way
			lineItems, err = order.ConfigureUniformOrder(totalQuantity)
			if err != nil {
				return err
			}
		} else {
			// Complex flow - split into multiple preferences
			lineItems, err = order.ConfigureLineItems(totalQuantity)
			if err != nil {
				return err
			}
		}
	}

	// Step 3: Show summary and confirm
	if err := showOrderSummary(tier, totalQuantity, lineItems); err != nil {
		return err
	}

	confirmed, err := prompts.PromptConfirm("Looks good! Proceed to checkout?")
	if err != nil {
		return err
	}
	if !confirmed {
		fmt.Println("\nOrder cancelled.")
		return nil
	}

	// Step 4: Create order via API
	fmt.Print("\nCreating order... ")
	order, err := client.CreateOrder(api.CreateOrderRequest{
		Tier:          tier.Tier,
		TotalQuantity: totalQuantity,
		LineItems:     lineItems,
	})
	if err != nil {
		fmt.Println("✗")
		return fmt.Errorf("failed to create order: %w", err)
	}
	fmt.Println("✓")

	// Step 5: Create checkout session
	fmt.Print("Opening checkout in your browser... ")
	checkout, err := client.CreateCheckoutSession(order.ID)
	if err != nil {
		fmt.Println("✗")
		return fmt.Errorf("failed to create checkout session: %w", err)
	}
	fmt.Println("✓")

	// Step 6: Open browser
	if err := openBrowser(checkout.CheckoutURL); err != nil {
		fmt.Printf("\nCouldn't open browser automatically. Please visit:\n%s\n", checkout.CheckoutURL)
	}

	fmt.Printf("\nOrder created successfully!\n")
	fmt.Printf("Order ID: %s\n\n", order.ID)

	// Step 7: Wait for payment completion
	fmt.Println("Waiting for payment confirmation...")
	fmt.Printf("(You have %d minutes to complete the payment)\n", PaymentTimeoutSeconds/60)

	subscription, completed := waitForSubscriptionActivation(client, order.ID, PaymentTimeoutSeconds)

	if completed && subscription != nil {
		// Payment successful!
		if err := templates.RenderToStdout(templates.SuccessArtTemplate, nil); err != nil {
			fmt.Printf("Error rendering template: %v\n", err)
		}
		if err := templates.RenderToStdout(templates.SuccessMessageTemplate, struct {
			TotalQuantity int
			TierName      string
		}{
			TotalQuantity: totalQuantity,
			TierName:      tier.Name,
		}); err != nil {
			fmt.Printf("Error rendering template: %v\n", err)
		}
	} else {
		// Timeout or user didn't complete payment
		fmt.Println("Complete your payment to activate your subscription.")
		fmt.Println("Your order will be processed once payment is received.")
	}

	return nil
}

// openBrowser opens the specified URL in the default browser
func openBrowser(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		return fmt.Errorf("unsupported platform")
	}

	return cmd.Start()
}

// Helper functions for order configuration

func showOrderSummary(tier api.AvailableSubscription, totalQuantity int, lineItems []api.OrderLineItem) error {
	// Calculate price based on quantity
	// tier.Price is the price per 1kg
	pricePerKg, _ := strconv.ParseFloat(tier.Price, 64)
	totalPrice := pricePerKg * float64(totalQuantity)

	// Format line items for display
	formattedItems := make([]string, len(lineItems))
	for i, item := range lineItems {
		if item.GrindType == "whole_bean" {
			formattedItems[i] = fmt.Sprintf("%d → Whole beans for %s",
				item.Quantity,
				order.BrewingMethodDisplay(item.BrewingMethod))
		} else {
			grindDesc := order.GetGrindDescription(item.BrewingMethod)
			formattedItems[i] = fmt.Sprintf("%d → Ground for %s (%s)",
				item.Quantity,
				order.BrewingMethodDisplay(item.BrewingMethod),
				grindDesc)
		}
	}

	fmt.Println(templates.RenderOrderSummary(
		tier.Name,
		totalQuantity,
		tier.Currency,
		totalPrice,
		tier.BillingPeriod,
		formattedItems,
	))
	return nil
}

// waitForSubscriptionActivation polls the API for subscription activation
func waitForSubscriptionActivation(client *api.Client, orderID string, timeoutSeconds int) (*api.Subscription, bool) {
	ticker := time.NewTicker(PaymentPollInterval)
	defer ticker.Stop()

	timeout := time.After(time.Duration(timeoutSeconds) * time.Second)
	dots := 0

	for {
		select {
		case <-ticker.C:
			// Poll for order status
			order, err := client.GetOrder(orderID)
			if err == nil && order.Status == "paid" {
				// Order is paid, fetch subscription
				subscriptions, err := client.ListSubscriptions()
				if err == nil && len(subscriptions) > 0 {
					// Find the active subscription
					for _, sub := range subscriptions {
						if sub.Status == "active" {
							return &sub, true
						}
					}
				}
			}

			// Show progress dots
			dots++
			if dots > 3 {
				dots = 1
			}
			fmt.Printf("\rChecking payment status%s   ", strings.Repeat(".", dots))

		case <-timeout:
			fmt.Println("\r" + strings.Repeat(" ", 50)) // Clear the line
			return nil, false
		}
	}
}

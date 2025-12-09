package cmd

import (
	"bufio"
	"fmt"
	"os"
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

var manageCmd = &cobra.Command{
	Use:   "manage",
	Short: "Manage your existing subscriptions",
	Long:  `Pause, resume, update, or cancel your Butler Coffee subscriptions.`,
	RunE:  runManage,
}

func init() {
	rootCmd.AddCommand(manageCmd)
}

func runManage(cmd *cobra.Command, args []string) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if !cfg.IsAuthenticated() {
		if err := templates.RenderToStdout(templates.ManageNotAuthenticatedTemplate, nil); err != nil {
			return err
		}
		return nil
	}

	client := api.NewClient(cfg)

	subscriptions, err := client.ListSubscriptions()
	if err != nil {
		return fmt.Errorf("failed to get subscriptions: %w", err)
	}

	if len(subscriptions) == 0 {
		if err := templates.RenderToStdout(templates.NoSubscriptionsTemplate, nil); err != nil {
			return err
		}
		return nil
	}

	subscription, err := selectSubscriptionToManage(subscriptions)
	if err != nil {
		return err
	}

	if subscription == nil {
		return nil
	}

	// Fetch full subscription details including order configuration
	fullSubscription, err := client.GetSubscription(subscription.ID)
	if err != nil {
		// If we can't fetch full details, continue with what we have
		fmt.Printf("Note: Could not fetch full subscription details: %v\n\n", err)
		fullSubscription = subscription
	}

	return showManagementMenu(cfg, client, fullSubscription)
}

func selectSubscriptionToManage(subscriptions []api.Subscription) (*api.Subscription, error) {
	// Filter out cancelled subscriptions
	var activeSubscriptions []api.Subscription
	for _, sub := range subscriptions {
		if sub.Status != "cancelled" {
			activeSubscriptions = append(activeSubscriptions, sub)
		}
	}

	// If no active subscriptions, return early
	if len(activeSubscriptions) == 0 {
		fmt.Println("You don't have any active subscriptions to manage.")
		fmt.Println("All your subscriptions have been cancelled.")
		return nil, nil
	}

	items := make([]models.ManageSubscriptionItem, len(activeSubscriptions)+1)
	for i, sub := range activeSubscriptions {
		statusIcon := getStatusIcon(sub.Status)
		display := fmt.Sprintf("%s %s (%s)", statusIcon, sub.Tier, sub.Status)

		item := models.ManageSubscriptionItem{
			Subscription:    sub,
			Display:         display,
			Status:          sub.Status,
			StartedAt:       "",
			ExpiresAt:       "",
			TotalQuantity:   0,
			HasOrderDetails: false,
			IsExit:          false,
		}

		if sub.StartedAt != nil {
			item.StartedAt = utils.FormatTimestamp(*sub.StartedAt)
		}
		if sub.ExpiresAt != nil {
			item.ExpiresAt = utils.FormatTimestamp(*sub.ExpiresAt)
		}
		if sub.DefaultQuantity != "" {
			item.TotalQuantity = sub.GetTotalQuantity()
			item.HasOrderDetails = true
		}

		items[i] = item
	}

	items[len(activeSubscriptions)] = models.ManageSubscriptionItem{
		Display: "← Exit",
		IsExit:  true,
	}

	return models.PickManageSubscription(items)
}

func showManagementMenu(cfg *config.Config, client *api.Client, subscription *api.Subscription) error {
	for {
		if err := displaySubscriptionInfo(client, subscription); err != nil {
			return err
		}

		actions := buildActionMenu(subscription.Status)

		if len(actions) == 0 {
			if err := templates.RenderToStdout(templates.NoActionsAvailableTemplate, nil); err != nil {
				return err
			}
			return nil
		}

		action, err := models.SelectAction(actions)
		if err != nil {
			return nil
		}

		if action == "exit" {
			return nil
		}

		updatedSub, err := executeAction(cfg, client, subscription, action)
		if err != nil {
			fmt.Printf("\nError: %v\n\n", err)
			fmt.Print("Press Enter to continue...")
			_, _ = bufio.NewReader(os.Stdin).ReadBytes('\n')
			continue
		}

		if updatedSub != nil {
			subscription = updatedSub
		}
	}
}

func buildActionMenu(status string) []models.ActionItem {
	var actions []models.ActionItem

	switch status {
	case "active":
		actions = append(actions, models.ActionItem{Action: "pause", Display: "⏸  Pause subscription"})
		actions = append(actions, models.ActionItem{Action: "update", Display: "✏  Update preferences"})
		actions = append(actions, models.ActionItem{Action: "cancel", Display: "✕ Cancel subscription"})
	case "paused":
		actions = append(actions, models.ActionItem{Action: "resume", Display: "▶  Resume subscription"})
		actions = append(actions, models.ActionItem{Action: "update", Display: "✏  Update preferences"})
		actions = append(actions, models.ActionItem{Action: "cancel", Display: "✕ Cancel subscription"})
	case "cancelled":
		return actions
	}

	actions = append(actions, models.ActionItem{Action: "exit", Display: "← Exit"})
	return actions
}

func executeAction(cfg *config.Config, client *api.Client, subscription *api.Subscription, action string) (*api.Subscription, error) {
	switch action {
	case "pause":
		return handlePause(client, subscription)
	case "resume":
		return handleResume(client, subscription)
	case "update":
		return handleUpdate(cfg, client, subscription)
	case "cancel":
		return handleCancel(client, subscription)
	}
	return nil, fmt.Errorf("unknown action: %s", action)
}

func handlePause(client *api.Client, subscription *api.Subscription) (*api.Subscription, error) {
	if err := templates.RenderToStdout(templates.PauseWarningTemplate, nil); err != nil {
		return nil, err
	}

	confirmed, err := prompts.PromptConfirm("Pause subscription? (You can resume it anytime)")
	if err != nil || !confirmed {
		if err := templates.RenderToStdout(templates.ActionCancelledTemplate, struct{ Action string }{Action: "Pause"}); err != nil {
			fmt.Println("Pause cancelled.")
		}
		return nil, nil
	}

	fmt.Print("\nPausing subscription... ")
	updatedSub, err := client.PauseSubscription(subscription.ID)
	if err != nil {
		fmt.Println("✗")
		return nil, err
	}
	fmt.Println("✓")

	if err := templates.RenderToStdout(templates.SubscriptionPausedTemplate, struct {
		HasResumeDate bool
		ResumeDate    string
	}{
		HasResumeDate: false,
		ResumeDate:    "",
	}); err != nil {
		return nil, err
	}

	return updatedSub, nil
}

func handleResume(client *api.Client, subscription *api.Subscription) (*api.Subscription, error) {
	if err := templates.RenderToStdout(templates.ResumeInfoTemplate, nil); err != nil {
		return nil, err
	}

	confirmed, err := prompts.PromptConfirm("Resume subscription?")
	if err != nil || !confirmed {
		if err := templates.RenderToStdout(templates.ActionCancelledTemplate, struct{ Action string }{Action: "Resume"}); err != nil {
			fmt.Println("Resume cancelled.")
		}
		return nil, nil
	}

	fmt.Print("\nResuming subscription... ")
	updatedSub, err := client.ResumeSubscription(subscription.ID)
	if err != nil {
		fmt.Println("✗")
		return nil, err
	}
	fmt.Println("✓")

	if err := templates.RenderToStdout(templates.SubscriptionResumedTemplate, nil); err != nil {
		return nil, err
	}

	return updatedSub, nil
}

func handleUpdate(cfg *config.Config, client *api.Client, subscription *api.Subscription) (*api.Subscription, error) {
	if err := templates.RenderToStdout(templates.UpdateSubscriptionHeaderTemplate, nil); err != nil {
		return nil, err
	}

	availableSubs, err := client.GetAvailableSubscriptions()
	if err != nil {
		return nil, fmt.Errorf("failed to get available subscriptions: %w", err)
	}

	var currentTier *api.AvailableSubscription
	for _, tier := range availableSubs {
		if tier.Tier == subscription.Tier {
			currentTier = &tier
			break
		}
	}

	if currentTier == nil {
		return nil, fmt.Errorf("could not find tier information")
	}

	// Use default preference quantity or the minimum quantity, whichever is larger
	defaultQty := max(order.DefaultPreferenceQuantity, cfg.MinQuantity)

	totalQuantity, err := prompts.PromptQuantityInt("New total quantity per month", cfg.MinQuantity, cfg.MaxQuantity, defaultQty)
	if err != nil {
		if err := templates.RenderToStdout(templates.ActionCancelledTemplate, struct{ Action string }{Action: "Update"}); err != nil {
			fmt.Println("Update cancelled.")
		}
		return nil, nil
	}

	fmt.Printf("\n✓ New quantity: %d per month\n\n", totalQuantity)

	wantsSplit, err := prompts.PromptConfirm("Would you like different grind methods?")
	if err != nil {
		if err := templates.RenderToStdout(templates.ActionCancelledTemplate, struct{ Action string }{Action: "Update"}); err != nil {
			fmt.Println("Update cancelled.")
		}
		return nil, nil
	}

	var lineItems []api.OrderLineItem

	if !wantsSplit {
		lineItems, err = order.ConfigureUniformOrder(totalQuantity)
		if err != nil {
			return nil, err
		}
	} else {
		lineItems, err = order.ConfigureLineItems(totalQuantity)
		if err != nil {
			return nil, err
		}
	}

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

	if err := templates.RenderToStdout(templates.UpdatePreferencesSummaryTemplate, struct {
		TotalQuantity int
		LineItems     []string
	}{
		TotalQuantity: totalQuantity,
		LineItems:     formattedItems,
	}); err != nil {
		return nil, err
	}

	confirmed, err := prompts.PromptConfirm("Update subscription with these preferences?")
	if err != nil || !confirmed {
		if err := templates.RenderToStdout(templates.ActionCancelledTemplate, struct{ Action string }{Action: "Update"}); err != nil {
			fmt.Println("Update cancelled.")
		}
		return nil, nil
	}

	fmt.Print("\nUpdating subscription... ")
	updatedSub, err := client.UpdateSubscription(subscription.ID, api.UpdateSubscriptionRequest{
		TotalQuantity: totalQuantity,
		Preferences:   lineItems,
	})
	if err != nil {
		fmt.Println("✗")
		return nil, err
	}
	fmt.Println("✓")

	if err := templates.RenderToStdout(templates.SubscriptionUpdatedTemplate, nil); err != nil {
		return nil, err
	}

	return updatedSub, nil
}

func handleCancel(client *api.Client, subscription *api.Subscription) (*api.Subscription, error) {
	if err := templates.RenderToStdout(templates.CancelWarningTemplate, nil); err != nil {
		return nil, err
	}

	// Offer pause as an alternative
	options := []models.ActionItem{
		{Action: "pause", Display: "⏸  Pause subscription instead (you can resume anytime)"},
		{Action: "cancel", Display: "✕ Cancel permanently"},
		{Action: "back", Display: "← Go back"},
	}

	action, err := models.SelectAction(options)
	if err != nil || action == "back" || action == "" {
		return nil, nil
	}

	if action == "pause" {
		return handlePause(client, subscription)
	}

	// Proceed with cancellation
	if err := templates.RenderToStdout(templates.CancelDoubleConfirmTemplate, nil); err != nil {
		return nil, err
	}

	confirmed, err := prompts.PromptConfirm("Type 'y' to permanently cancel")
	if err != nil || !confirmed {
		if err := templates.RenderToStdout(templates.ActionCancelledTemplate, struct{ Action string }{Action: "Cancellation"}); err != nil {
			fmt.Println("Cancellation aborted.")
		}
		return nil, nil
	}

	fmt.Print("\nCancelling subscription... ")
	updatedSub, err := client.CancelSubscription(subscription.ID)
	if err != nil {
		fmt.Println("✗")
		return nil, err
	}
	fmt.Println("✓")

	if err := templates.RenderToStdout(templates.SubscriptionCancelledTemplate, nil); err != nil {
		return nil, err
	}

	return updatedSub, nil
}

func displaySubscriptionInfo(client *api.Client, subscription *api.Subscription) error {
	statusIcon := getStatusIcon(subscription.Status)

	data := struct {
		Tier            string
		Status          string
		StatusIcon      string
		StartedAt       string
		ExpiresAt       string
		NextShipment    string
		HasNextShipment bool
		HasOrderDetails bool
		TotalQuantity   int
		LineItems       []string
		HasPricing      bool
		Price           string
		Currency        string
		BillingPeriod   string
	}{
		Tier:            subscription.Tier,
		Status:          subscription.Status,
		StatusIcon:      statusIcon,
		StartedAt:       "",
		ExpiresAt:       "",
		NextShipment:    "",
		HasNextShipment: false,
		HasOrderDetails: false,
		TotalQuantity:   0,
		LineItems:       []string{},
		HasPricing:      false,
		Price:           "",
		Currency:        "",
		BillingPeriod:   "",
	}

	if subscription.StartedAt != nil {
		data.StartedAt = utils.FormatTimestamp(*subscription.StartedAt)
	}
	if subscription.ExpiresAt != nil {
		data.ExpiresAt = utils.FormatTimestamp(*subscription.ExpiresAt)
	}

	// Calculate next shipment date for active subscriptions
	if subscription.Status == "active" && subscription.StartedAt != nil {
		if nextShipment := calculateNextShipment(*subscription.StartedAt); nextShipment != "" {
			data.NextShipment = nextShipment
			data.HasNextShipment = true
		}
	}

	// Include order details if available
	if subscription.DefaultQuantity != "" && len(subscription.DefaultPreferences) > 0 {
		data.HasOrderDetails = true
		data.TotalQuantity = subscription.GetTotalQuantity()

		// Fetch pricing information and calculate actual price based on quantity
		if pricing, err := client.GetSubscriptionPricing(subscription.Tier); err == nil {
			// Parse base price and multiply by quantity
			var basePrice float64
			if _, err := fmt.Sscanf(pricing.Price, "%f", &basePrice); err == nil {
				actualPrice := basePrice * float64(data.TotalQuantity)

				data.HasPricing = true
				data.Price = fmt.Sprintf("%.2f", actualPrice)
				data.Currency = pricing.Currency
				data.BillingPeriod = pricing.BillingPeriod
			}
		}

		// Format preferences
		data.LineItems = make([]string, len(subscription.DefaultPreferences))
		for i, pref := range subscription.DefaultPreferences {
			qty := pref.GetQuantity()
			if pref.GrindType == "whole_bean" {
				data.LineItems[i] = fmt.Sprintf("%d → Whole beans for %s",
					qty,
					order.BrewingMethodDisplay(pref.BrewingMethod))
			} else {
				grindDesc := order.GetGrindDescription(pref.BrewingMethod)
				data.LineItems[i] = fmt.Sprintf("%d → Ground for %s (%s)",
					qty,
					order.BrewingMethodDisplay(pref.BrewingMethod),
					grindDesc)
			}
		}
	}

	return templates.RenderToStdout(templates.ManageSubscriptionHeaderTemplate, data)
}

func calculateNextShipment(startedAtStr string) string {
	// Parse the started_at timestamp
	startedAt, err := utils.ParseTimestamp(startedAtStr)
	if err != nil {
		return ""
	}

	now := time.Now()

	// Calculate the next shipment date by adding months until we're in the future
	nextShipment := startedAt
	for nextShipment.Before(now) || nextShipment.Equal(now) {
		nextShipment = nextShipment.AddDate(0, 1, 0) // Add 1 month
	}

	return utils.FormatDate(nextShipment)
}

func getStatusIcon(status string) string {
	switch status {
	case "active":
		return "✓"
	case "paused":
		return "⏸"
	case "cancelled":
		return "✕"
	default:
		return "•"
	}
}

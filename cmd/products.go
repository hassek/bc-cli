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
	"github.com/spf13/cobra"
)

var productsCmd = &cobra.Command{
	Use:   "products",
	Short: "Browse and purchase Butler Coffee products",
	Long:  `Browse our product catalog and purchase one-time coffee deals directly from the terminal.`,
	RunE:  runProducts,
}

func init() {
	rootCmd.AddCommand(productsCmd)
}

func runProducts(cmd *cobra.Command, args []string) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	client := api.NewClient(cfg)

	// Get available products
	available, err := client.GetAvailableProducts()
	if err != nil {
		return fmt.Errorf("failed to get available products: %w", err)
	}

	if len(available) == 0 {
		fmt.Println("No products available at this time.")
		return nil
	}

	// Use new product picker with duck animation
	selectedProduct, err := models.PickProduct(available)
	if err != nil {
		fmt.Println("\nExiting...")
		return nil
	}

	// User cancelled or selected exit
	if selectedProduct == nil {
		return nil
	}

	displayProductDetails(*selectedProduct)

	// Ask if user wants to purchase (if authenticated)
	if cfg.IsAuthenticated() {
		fmt.Println()
		confirmed, err := prompts.PromptConfirm(fmt.Sprintf("Would you like to purchase %s now", selectedProduct.Name))
		if err == nil && confirmed {
			// User wants to purchase - start order configuration flow
			return createProductOrder(cfg, client, *selectedProduct)
		}
	} else if !cfg.IsAuthenticated() {
		fmt.Println("\nPlease login first to purchase:")
		fmt.Println("  bc-cli login")
	}

	return nil
}

func displayProductDetails(product api.AvailableSubscription) {
	// Pre-render description to support template syntax (highlights, etc.)
	renderedDescription := templates.RenderDescription(product.Description)

	// Use viewport for scrollable display with template
	if err := templates.RenderInViewport(product.Name, templates.ProductDetailsTemplate, struct {
		Name        string
		Currency    string
		Price       string
		Description string
	}{
		Name:        product.Name,
		Currency:    product.Currency,
		Price:       product.Price,
		Description: renderedDescription,
	}); err != nil {
		fmt.Printf("Error displaying product details: %v\n", err)
	}
}

func createProductOrder(cfg *config.Config, client *api.Client, product api.AvailableSubscription) error {
	if !cfg.IsAuthenticated() {
		return fmt.Errorf("you must be logged in to purchase. Please run 'bc-cli login' first")
	}

	fmt.Println("\n" + strings.Repeat("═", 60))
	fmt.Println("\n  Configure Your Order")
	fmt.Println("\n" + strings.Repeat("─", 60) + "\n")

	// Use min/max quantity from backend, fallback to defaults if not set
	minQty := product.MinQuantity
	maxQty := product.MaxQuantity
	if minQty == 0 {
		minQty = 1
	}
	if maxQty == 0 {
		maxQty = 10
	}

	// Step 1: Ask for quantity (number of items, not kg)
	quantity, err := prompts.PromptQuantityInt("How many would you like to purchase?", minQty, maxQty, minQty)
	if err != nil {
		return err
	}

	fmt.Printf("\n✓ Quantity: %d\n", quantity)

	// Step 2: Ask for grind type
	fmt.Println("\n  How would you like your coffee prepared?")
	fmt.Println()
	grindType, err := order.SelectGrindType()
	if err != nil {
		return err
	}

	// Show grind type confirmation
	fmt.Println()
	if grindType == "ground" {
		fmt.Println("✓ We'll grind these beans for you!")
	} else {
		fmt.Println("✓ You'll grind these beans yourself")
	}

	// Step 3: Ask for brewing method (ALWAYS)
	brewResult, err := order.SelectBrewingMethod(grindType)
	if err != nil {
		return err
	}

	// Confirmation message
	fmt.Printf("\n✓ Perfect! Your order: %d x %s - ", quantity, product.Name)
	if grindType == "whole_bean" {
		fmt.Printf("whole beans for %s.\n", order.BrewingMethodDisplay(brewResult.Method))
	} else {
		grindDesc := order.GetGrindDescription(brewResult.Method)
		fmt.Printf("ground for %s (%s).\n", order.BrewingMethodDisplay(brewResult.Method), grindDesc)
	}
	if brewResult.Notes != "" {
		fmt.Printf("  Notes: %s\n", brewResult.Notes)
	}
	fmt.Println()

	// Step 4: Show order summary
	if err := showProductOrderSummary(product, quantity, grindType, brewResult.Method, brewResult.Notes); err != nil {
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

	// Step 5: Create order via API
	// For products, we use ProductID instead of Tier
	fmt.Print("\nCreating order... ")
	order, err := client.CreateOrder(api.CreateOrderRequest{
		ProductID:     product.ID,
		TotalQuantity: quantity,
		LineItems: []api.OrderLineItem{
			{
				Quantity:      quantity,
				GrindType:     grindType,
				BrewingMethod: brewResult.Method,
				Notes:         brewResult.Notes,
			},
		},
	})
	if err != nil {
		fmt.Println("✗")
		return fmt.Errorf("failed to create order: %w", err)
	}
	fmt.Println("✓")

	// Step 6: Create checkout session
	fmt.Print("Opening checkout in your browser... ")
	checkout, err := client.CreateCheckoutSession(order.ID)
	if err != nil {
		fmt.Println("✗")
		return fmt.Errorf("failed to create checkout session: %w", err)
	}
	fmt.Println("✓")

	// Step 7: Open browser
	if err := openProductBrowser(checkout.CheckoutURL); err != nil {
		fmt.Printf("\nCouldn't open browser automatically. Please visit:\n%s\n", checkout.CheckoutURL)
	}

	fmt.Printf("\nOrder created successfully!\n")
	fmt.Printf("Order ID: %s\n\n", order.ID)

	// Step 8: Wait for payment completion
	fmt.Println("Waiting for payment confirmation...")
	fmt.Println("(You have 5 minutes to complete the payment)")

	completed := waitForProductPayment(client, order.ID, 5*60) // 5 minutes

	if completed {
		// Payment successful!
		fmt.Println("\n" + strings.Repeat("═", 60))
		fmt.Println("\n  🎉 Payment Successful!")
		fmt.Println("\n" + strings.Repeat("─", 60) + "\n")
		fmt.Printf("  Your order for %d x %s has been confirmed!\n", quantity, product.Name)
		fmt.Println("  We'll start preparing your coffee right away.")
		fmt.Println("\n" + strings.Repeat("═", 60) + "\n")
	} else {
		// Timeout or user didn't complete payment
		fmt.Println("\nComplete your payment to confirm your order.")
		fmt.Println("Your order will be processed once payment is received.")
	}

	return nil
}

func showProductOrderSummary(product api.AvailableSubscription, quantity int, grindType, brewingMethod, notes string) error {
	// Calculate total price
	pricePerUnit, _ := strconv.ParseFloat(product.Price, 64)
	totalPrice := pricePerUnit * float64(quantity)

	fmt.Println("\n" + strings.Repeat("═", 60))
	fmt.Println("\n  Order Summary")
	fmt.Println("\n" + strings.Repeat("─", 60) + "\n")

	fmt.Printf("  Product:  %s\n", product.Name)
	fmt.Printf("  Quantity: %d\n", quantity)
	fmt.Printf("  Price:    %s %.2f\n\n", product.Currency, totalPrice)

	fmt.Println("  Preparation:")
	if grindType == "whole_bean" {
		fmt.Printf("    • Whole beans for %s\n", order.BrewingMethodDisplay(brewingMethod))
	} else {
		grindDesc := order.GetGrindDescription(brewingMethod)
		fmt.Printf("    • Ground for %s (%s)\n", order.BrewingMethodDisplay(brewingMethod), grindDesc)
	}
	if notes != "" {
		fmt.Printf("    • Notes: %s\n", notes)
	}

	fmt.Println("\n" + strings.Repeat("═", 60) + "\n")

	return nil
}

// openProductBrowser opens the specified URL in the default browser
func openProductBrowser(url string) error {
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

// waitForProductPayment polls the API for payment completion
func waitForProductPayment(client *api.Client, orderID string, timeoutSeconds int) bool {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	timeout := time.After(time.Duration(timeoutSeconds) * time.Second)
	dots := 0

	for {
		select {
		case <-ticker.C:
			// Poll for order status
			order, err := client.GetOrder(orderID)
			if err == nil && order.Status == "paid" {
				return true
			}

			// Show progress dots
			dots++
			if dots > 3 {
				dots = 1
			}
			fmt.Printf("\rChecking payment status%s   ", strings.Repeat(".", dots))

		case <-timeout:
			fmt.Println("\r" + strings.Repeat(" ", 50)) // Clear the line
			return false
		}
	}
}

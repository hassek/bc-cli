package order

import (
	"fmt"
	"strings"

	"github.com/hassek/bc-cli/api"
	"github.com/hassek/bc-cli/templates"
	"github.com/hassek/bc-cli/tui/models"
	"github.com/hassek/bc-cli/tui/prompts"
)

const (
	DefaultPreferenceQuantity = 2 // Default quantity for new preferences
)

// ConfigureUniformOrder guides the user through configuring a uniform order (all same grind/brew method)
func ConfigureUniformOrder(totalQuantity int) ([]api.OrderLineItem, error) {
	if err := templates.RenderToStdout(templates.UniformOrderIntroTemplate, struct{ TotalQuantity int }{TotalQuantity: totalQuantity}); err != nil {
		return nil, fmt.Errorf("failed to render template: %w", err)
	}

	// Prompt for grind type
	grindType, err := SelectGrindType()
	if err != nil {
		return nil, err
	}

	// Show confirmation based on choice
	fmt.Println()
	if grindType == "whole_bean" {
		fmt.Println("✓ You'll grind these beans yourself")
	} else {
		fmt.Println("✓ We'll grind these beans for you!")
	}

	// ALWAYS prompt for brewing method
	brewingMethod, err := SelectBrewingMethod(grindType)
	if err != nil {
		return nil, err
	}

	// Confirmation message
	fmt.Printf("\n✓ Perfect! All %d will be ", totalQuantity)
	if grindType == "whole_bean" {
		fmt.Printf("whole beans, roasted for %s.\n", BrewingMethodDisplay(brewingMethod))
	} else {
		fmt.Printf("ground for %s.\n", BrewingMethodDisplay(brewingMethod))
	}
	fmt.Println()

	// Create single line item with full quantity
	lineItems := []api.OrderLineItem{
		{
			Quantity:      totalQuantity,
			GrindType:     grindType,
			BrewingMethod: brewingMethod,
		},
	}

	return lineItems, nil
}

// ConfigureLineItems guides the user through configuring split line items (different grind/brew methods)
func ConfigureLineItems(totalQuantity int) ([]api.OrderLineItem, error) {
	var lineItems []api.OrderLineItem
	remaining := totalQuantity
	preferenceNum := 1

	// Introduction
	if err := templates.RenderToStdout(templates.SplitOrderIntroTemplate, struct{ TotalQuantity int }{TotalQuantity: totalQuantity}); err != nil {
		return nil, fmt.Errorf("failed to render template: %w", err)
	}

	for remaining > 0 {
		// Show preference header with remaining amount
		lowRemaining := float64(remaining) < float64(totalQuantity)*0.3
		fmt.Println() // Add newline before the box
		fmt.Println(templates.RenderPreferenceHeader(preferenceNum, totalQuantity, remaining, lowRemaining))

		// Prompt for quantity with smart defaults
		maxQty := remaining
		defaultQty := remaining
		if remaining > DefaultPreferenceQuantity {
			defaultQty = min(DefaultPreferenceQuantity, remaining)
		}

		quantity, err := prompts.PromptQuantityInt("  How much for this preference?", 1, maxQty, defaultQty)
		if err != nil {
			return nil, err
		}

		// Show allocation confirmation
		if quantity >= remaining {
			fmt.Printf("\n✓ Allocating %d (this will complete your order!)\n\n", quantity)
		} else {
			fmt.Printf("\n✓ Allocating %d\n\n", quantity)
		}

		// Prompt for grind type with explanation
		fmt.Printf("  How would you like these %d prepared?\n\n", quantity)
		grindType, err := SelectGrindType()
		if err != nil {
			return nil, err
		}

		// Show grind type confirmation
		fmt.Println()
		if grindType == "ground" {
			fmt.Println("✓ We'll grind these beans for you!")
		} else {
			fmt.Println("✓ You'll grind these beans yourself")
		}

		// Prompt for brewing method (ALWAYS, regardless of grind type)
		brewingMethod, err := SelectBrewingMethod(grindType)
		if err != nil {
			return nil, err
		}

		lineItems = append(lineItems, api.OrderLineItem{
			Quantity:      quantity,
			GrindType:     grindType,
			BrewingMethod: brewingMethod,
		})

		// Updated confirmation message
		fmt.Printf("\n✓ Added: %d ", quantity)
		if grindType == "whole_bean" {
			fmt.Printf("whole beans for %s", BrewingMethodDisplay(brewingMethod))
		} else {
			fmt.Printf("ground for %s", BrewingMethodDisplay(brewingMethod))
		}
		fmt.Println()

		remaining -= quantity

		// Show progress bar
		ShowProgressBar(totalQuantity-remaining, totalQuantity)

		preferenceNum++

		// Check if we're done
		if remaining <= 0 {
			break
		}
	}

	// Success message
	fmt.Println("\n" + strings.Repeat("─", 60) + "\n")
	fmt.Printf("🎉 Perfect! You've allocated all %d!\n\n", totalQuantity)

	return lineItems, nil
}

// SelectGrindType prompts the user to select a grind type
func SelectGrindType() (string, error) {
	return models.SelectGrindType()
}

// SelectBrewingMethod prompts the user to select a brewing method
func SelectBrewingMethod(grindType string) (string, error) {
	// Show helpful message first
	fmt.Println("  What is your preferred brewing method?")
	fmt.Println("  This helps us understand the best profiles to ensure the best tasting experience!")
	fmt.Println()

	return models.SelectBrewingMethod(grindType)
}

// BrewingMethodDisplay returns the display name for a brewing method
func BrewingMethodDisplay(method string) string {
	displays := map[string]string{
		"espresso":     "Espresso",
		"moka":         "Moka Pot",
		"v60":          "V60 Pour Over",
		"french_press": "French Press",
		"pour_over":    "Pour Over",
		"drip":         "Drip Coffee",
		"cold_brew":    "Cold Brew",
	}
	if display, ok := displays[method]; ok {
		return display
	}
	return method
}

// GetGrindDescription returns the grind description for a brewing method
func GetGrindDescription(method string) string {
	descriptions := map[string]string{
		"espresso":     "very fine",
		"moka":         "fine-medium",
		"v60":          "medium",
		"french_press": "coarse",
		"pour_over":    "medium",
		"drip":         "medium",
		"cold_brew":    "extra coarse",
	}
	if desc, ok := descriptions[method]; ok {
		return desc
	}
	return ""
}

// ShowProgressBar displays a progress bar
func ShowProgressBar(current, total int) {
	fmt.Println(templates.RenderProgressBar(current, total))
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

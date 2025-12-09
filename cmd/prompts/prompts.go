package prompts

import (
	"fmt"
	"strings"

	"github.com/manifoldco/promptui"
)

// PromptQuantityInt prompts the user to enter a quantity as an integer
func PromptQuantityInt(label string, min, max, defaultVal int) (int, error) {
	validate := func(input string) error {
		var val int
		_, err := fmt.Sscanf(input, "%d", &val)
		if err != nil {
			return fmt.Errorf("please enter a valid whole number")
		}
		if val < min || val > max {
			return fmt.Errorf("quantity must be between %d and %d", min, max)
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    label,
		Validate: validate,
		Default:  fmt.Sprintf("%d", defaultVal),
	}

	result, err := prompt.Run()
	if err != nil {
		return 0, err
	}

	var quantity int
	if _, err := fmt.Sscanf(result, "%d", &quantity); err != nil {
		return 0, fmt.Errorf("invalid quantity: %w", err)
	}
	return quantity, nil
}

// PromptConfirm prompts the user for a yes/no confirmation
func PromptConfirm(label string) (bool, error) {
	prompt := promptui.Prompt{
		Label:     label,
		IsConfirm: true,
		Default:   "y",
	}

	result, err := prompt.Run()
	if err != nil {
		if err == promptui.ErrAbort {
			return false, nil
		}
		return false, err
	}

	// If empty (user just pressed Enter), default to yes
	result = strings.TrimSpace(result)
	if result == "" {
		return true, nil
	}

	return strings.ToLower(result) == "y" || strings.ToLower(result) == "yes", nil
}

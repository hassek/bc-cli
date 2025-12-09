package templates

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	// Box style for preference headers and progress bars
	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(0, 1).
			Width(60)

	// Style for headers
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("14"))

	// Style for info text
	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("7"))

	// Style for warnings
	warningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("11"))
)

// RenderPreferenceHeader renders a preference header box using Lipgloss
func RenderPreferenceHeader(preferenceNum, totalQuantity, remaining int, lowRemaining bool) string {
	var lines []string

	// Header line
	header := headerStyle.Render(fmt.Sprintf("Preference #%d", preferenceNum))
	lines = append(lines, header)

	// Allocating from line
	allocLine := infoStyle.Render(fmt.Sprintf("Allocating from: %d total", totalQuantity))
	lines = append(lines, allocLine)

	// Remaining line
	remainingText := fmt.Sprintf("Remaining: %d", remaining)
	if lowRemaining {
		remainingText += " ⚠️  (almost done!)"
		lines = append(lines, warningStyle.Render(remainingText))
	} else {
		lines = append(lines, infoStyle.Render(remainingText))
	}

	content := strings.Join(lines, "\n")
	return boxStyle.Render(content)
}

// RenderProgressBar renders a progress bar box using Lipgloss
func RenderProgressBar(current, total int) string {
	// Calculate percentage and build progress bar
	percentage := float64(current) / float64(total)
	barWidth := 30
	filled := int(percentage * float64(barWidth))
	bar := strings.Repeat("█", filled) + strings.Repeat("░", barWidth-filled)

	// Build the progress text
	progressText := fmt.Sprintf("Progress: %s %d/%d", bar, current, total)
	if current >= total {
		progressText += " ✓"
	}

	return boxStyle.Render(infoStyle.Render(progressText))
}

// RenderOrderSummary renders the order summary box using Lipgloss
func RenderOrderSummary(tierName string, totalQuantity int, currency string, totalPrice float64, billingPeriod string, lineItems []string) string {
	var lines []string

	// Title
	lines = append(lines, headerStyle.Render("Your Order Summary"))
	lines = append(lines, "")

	// Order details
	lines = append(lines, infoStyle.Render(fmt.Sprintf("Tier: %s", tierName)))
	lines = append(lines, infoStyle.Render(fmt.Sprintf("Total: %d/month", totalQuantity)))
	lines = append(lines, infoStyle.Render(fmt.Sprintf("Price: %s %.2f/%s", currency, totalPrice, billingPeriod)))
	lines = append(lines, "")
	lines = append(lines, infoStyle.Render("How your coffee will be prepared:"))

	// Line items
	for i, item := range lineItems {
		lines = append(lines, infoStyle.Render(fmt.Sprintf("   %d. %s", i+1, item)))
	}

	content := strings.Join(lines, "\n")

	summaryBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		Padding(1, 2).
		Width(65)

	return summaryBox.Render(content)
}

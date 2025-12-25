package templates

import (
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"

	"github.com/charmbracelet/lipgloss"
	"github.com/hassek/bc-cli/tui/components"
	"github.com/hassek/bc-cli/tui/styles"
	"github.com/hassek/bc-cli/utils"
)

// Common template functions
var funcMap = template.FuncMap{
	"upper":  strings.ToUpper,
	"lower":  strings.ToLower,
	"repeat": strings.Repeat,
	"printf": fmt.Sprintf,
	"add": func(a, b int) int {
		return a + b
	},
	"percentage": func(current, total int) float64 {
		if total == 0 {
			return 0
		}
		return float64(current) / float64(total) * 100
	},
	"progressBar": func(current, total, width int) string {
		if total == 0 {
			return strings.Repeat("░", width)
		}
		percentage := float64(current) / float64(total)
		filled := int(percentage * float64(width))
		bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
		return bar
	},
	"wrap": func(text string, width int) string {
		return utils.WrapText(text, width)
	},
	"wrapAuto": func(text string) string {
		// Auto-detect terminal width and use most of it with margins
		termWidth := utils.GetTerminalWidth()
		// Use 90% of terminal width with reasonable min/max bounds
		maxWidth := int(float64(termWidth) * 0.9)
		if maxWidth < 60 {
			maxWidth = 60 // Minimum width for readability
		}
		if maxWidth > 120 {
			maxWidth = 120 // Maximum width for readability
		}
		return utils.WrapText(text, maxWidth)
	},
	// Style functions for inline text formatting
	"highlight": func(text string) string {
		return styles.ActiveStyle.Render(text) // Bold cyan
	},
	"emphasis": func(text string) string {
		return lipgloss.NewStyle().
			Italic(true).
			Foreground(lipgloss.Color("247")).
			Render(text)
	},
	"section": func(text string) string {
		return lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("214")). // Orange/yellow
			Render(text)
	},
	"bold": func(text string) string {
		return lipgloss.NewStyle().Bold(true).Render(text)
	},
	"faint": func(text string) string {
		return styles.FaintStyle.Render(text)
	},
	"faintNoWrap": func(text string) string {
		// Faint text without wrapping - useful for dividers/lines
		style := lipgloss.NewStyle().
			Foreground(styles.Faint).
			Width(0) // Don't wrap
		return style.Render(text)
	},
	"cyan": func(text string) string {
		return lipgloss.NewStyle().Foreground(styles.Cyan).Render(text)
	},
	"green": func(text string) string {
		return lipgloss.NewStyle().Foreground(styles.Green).Render(text)
	},
	"yellow": func(text string) string {
		return lipgloss.NewStyle().Foreground(styles.Yellow).Render(text)
	},
	"red": func(text string) string {
		return lipgloss.NewStyle().Foreground(styles.Red).Render(text)
	},
	// Paragraph styling with wrapping
	"paragraph": func(text string, width int) string {
		style := lipgloss.NewStyle().
			Width(width).
			MarginBottom(1)
		return style.Render(text)
	},
	"paragraphAuto": func(text string) string {
		// Auto-detect terminal width and use most of it with margins
		termWidth := utils.GetTerminalWidth()
		// Use 90% of terminal width with reasonable min/max bounds
		maxWidth := int(float64(termWidth) * 0.9)
		if maxWidth < 60 {
			maxWidth = 60 // Minimum width for readability
		}
		if maxWidth > 120 {
			maxWidth = 120 // Maximum width for readability
		}
		style := lipgloss.NewStyle().
			Width(maxWidth).
			MarginBottom(1)
		return style.Render(text)
	},
}

// Render renders a template with the given data
func Render(w io.Writer, tmpl string, data any) error {
	t, err := template.New("output").Funcs(funcMap).Parse(tmpl)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}
	return t.Execute(w, data)
}

// RenderToStdout renders a template to stdout
func RenderToStdout(tmpl string, data any) error {
	return Render(os.Stdout, tmpl, data)
}

// RenderToString renders a template to a string
func RenderToString(tmpl string, data any) (string, error) {
	var buf strings.Builder
	if err := Render(&buf, tmpl, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// RenderInViewport renders a template in a scrollable viewport
func RenderInViewport(title, tmpl string, data any) error {
	content, err := RenderToString(tmpl, data)
	if err != nil {
		return err
	}
	return components.ShowTextViewer(title, content)
}

// RenderDescription renders a description string that may contain template syntax.
// If the description contains template markers ({{ }}), it renders them.
// Otherwise, it returns the description as-is.
// This allows backend CMS to send descriptions with highlighting and formatting.
func RenderDescription(description string) string {
	// Check if description contains template syntax
	if !strings.Contains(description, "{{") {
		return description
	}

	// Render as template
	rendered, err := RenderToString(description, nil)
	if err != nil {
		// If template rendering fails, show error for debugging
		return fmt.Sprintf("Template rendering error: %v\n\nOriginal content:\n%s", err, description)
	}
	return rendered
}

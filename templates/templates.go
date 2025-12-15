package templates

import (
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"

	"github.com/hassek/bc-cli/tui/components"
	"github.com/hassek/bc-cli/utils"
)

// Common template functions
var funcMap = template.FuncMap{
	"upper":  strings.ToUpper,
	"lower":  strings.ToLower,
	"repeat": strings.Repeat,
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
		// Auto-detect terminal width and use 60 as max for better readability
		termWidth := utils.GetTerminalWidth()
		maxWidth := 60
		if termWidth < maxWidth {
			maxWidth = termWidth - 4 // Leave some margin
		}
		return utils.WrapText(text, maxWidth)
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

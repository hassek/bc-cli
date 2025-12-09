package styles

import "github.com/charmbracelet/lipgloss"

// Color palette matching promptui design
var (
	// Primary colors
	Cyan   = lipgloss.Color("#00ffff")
	Green  = lipgloss.Color("#00ff00")
	Yellow = lipgloss.Color("#ffff00")
	Red    = lipgloss.Color("#ff0000")
	Faint  = lipgloss.Color("#808080")

	// Text styles
	ActiveStyle   = lipgloss.NewStyle().Foreground(Cyan).Bold(true)
	SelectedStyle = lipgloss.NewStyle().Foreground(Green).Bold(true)
	InactiveStyle = lipgloss.NewStyle()
	FaintStyle    = lipgloss.NewStyle().Foreground(Faint)
	ErrorStyle    = lipgloss.NewStyle().Foreground(Red).Bold(true)

	// Duck styles
	DuckStyle       = lipgloss.NewStyle().Foreground(Cyan)
	DuckAccentStyle = lipgloss.NewStyle().Foreground(Yellow)

	// Cursor
	CursorStyle = lipgloss.NewStyle().Foreground(Cyan)
	Cursor      = "▸"

	// Spacing
	Indent = "  "
)

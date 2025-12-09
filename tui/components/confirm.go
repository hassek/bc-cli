package components

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hassek/bc-cli/tui/styles"
)

// ConfirmComponent handles yes/no confirmation prompts
type ConfirmComponent struct {
	label     string
	cursor    int
	options   []string
	submitted bool
	cancelled bool
	result    bool
}

func NewConfirmComponent(label string) *ConfirmComponent {
	return &ConfirmComponent{
		label:   label,
		cursor:  0,
		options: []string{"Yes", "No"},
	}
}

func (c *ConfirmComponent) Init() tea.Cmd {
	return nil
}

func (c *ConfirmComponent) Update(msg tea.Msg) (*ConfirmComponent, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			c.cancelled = true
			return c, tea.Quit

		case "enter":
			c.result = (c.cursor == 0) // Yes is at index 0
			c.submitted = true
			return c, tea.Quit

		case "left", "h", "up", "k":
			c.cursor = 0

		case "right", "l", "down", "j":
			c.cursor = 1

		case "y", "Y":
			c.result = true
			c.submitted = true
			return c, tea.Quit

		case "n", "N":
			c.result = false
			c.submitted = true
			return c, tea.Quit
		}
	}
	return c, nil
}

func (c *ConfirmComponent) View() string {
	var b strings.Builder

	// Label
	b.WriteString(styles.ActiveStyle.Render(c.label))
	b.WriteString("\n\n")

	// Options
	for i, option := range c.options {
		cursor := "  "
		style := styles.InactiveStyle

		if i == c.cursor {
			cursor = styles.CursorStyle.Render(styles.Cursor) + " "
			style = styles.ActiveStyle
		}

		b.WriteString(cursor)
		b.WriteString(style.Render(option))
		b.WriteString("  ")
	}
	b.WriteString("\n\n")

	// Help text
	b.WriteString(styles.FaintStyle.Render("Use ←→ or y/n to select, Enter to confirm, Esc to cancel"))

	return b.String()
}

func (c *ConfirmComponent) Submitted() bool {
	return c.submitted
}

func (c *ConfirmComponent) Cancelled() bool {
	return c.cancelled
}

func (c *ConfirmComponent) Result() bool {
	return c.result
}

package components

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hassek/bc-cli/tui/styles"
)

const maxVisibleItems = 10

// SelectItem represents an item that can be selected
type SelectItem interface {
	// Label returns the text to display for this item
	Label() string
	// Details returns optional details text shown below the item
	Details() string
}

// SimpleItem is a basic implementation of SelectItem
type SimpleItem struct {
	LabelText   string
	DetailsText string
}

func (s SimpleItem) Label() string   { return s.LabelText }
func (s SimpleItem) Details() string { return s.DetailsText }

// SelectComponent handles list selection with cursor navigation
type SelectComponent struct {
	items       []SelectItem
	cursor      int
	selected    bool
	cancelled   bool
	title       string
	scrollOffset int
}

func NewSelectComponent(title string, items []SelectItem) *SelectComponent {
	return &SelectComponent{
		title: title,
		items: items,
		cursor: 0,
	}
}

func (s *SelectComponent) Init() tea.Cmd {
	return nil
}

func (s *SelectComponent) Update(msg tea.Msg) (*SelectComponent, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			s.cancelled = true
			return s, tea.Quit

		case "enter":
			s.selected = true
			return s, tea.Quit

		case "up", "k":
			if s.cursor > 0 {
				s.cursor--
				s.adjustScroll()
			}

		case "down", "j":
			if s.cursor < len(s.items)-1 {
				s.cursor++
				s.adjustScroll()
			}

		case "home", "g":
			s.cursor = 0
			s.scrollOffset = 0

		case "end", "G":
			s.cursor = len(s.items) - 1
			s.adjustScroll()
		}
	}
	return s, nil
}

func (s *SelectComponent) adjustScroll() {
	if s.cursor < s.scrollOffset {
		s.scrollOffset = s.cursor
	} else if s.cursor >= s.scrollOffset+maxVisibleItems {
		s.scrollOffset = s.cursor - maxVisibleItems + 1
	}
}

func (s *SelectComponent) View() string {
	var b strings.Builder

	// Title
	if s.title != "" {
		b.WriteString(styles.ActiveStyle.Render(s.title))
		b.WriteString("\n\n")
	}

	// Calculate visible range
	start := s.scrollOffset
	end := s.scrollOffset + maxVisibleItems
	if end > len(s.items) {
		end = len(s.items)
	}

	// Show scroll indicator at top if needed
	if s.scrollOffset > 0 {
		b.WriteString(styles.FaintStyle.Render(fmt.Sprintf("  ↑ %d more above\n", s.scrollOffset)))
	}

	// Render visible items (no inline details)
	for i := start; i < end; i++ {
		item := s.items[i]
		cursor := "  "
		style := styles.InactiveStyle

		if i == s.cursor {
			cursor = styles.CursorStyle.Render(styles.Cursor) + " "
			style = styles.ActiveStyle
		}

		b.WriteString(cursor)
		b.WriteString(style.Render(item.Label()))
		b.WriteString("\n")
	}

	// Show scroll indicator at bottom if needed
	remaining := len(s.items) - end
	if remaining > 0 {
		b.WriteString(styles.FaintStyle.Render(fmt.Sprintf("  ↓ %d more below\n", remaining)))
	}

	// Fixed details panel below menu
	b.WriteString("\n")
	b.WriteString(strings.Repeat("━", 60))
	b.WriteString("\n")

	if s.cursor >= 0 && s.cursor < len(s.items) {
		details := s.items[s.cursor].Details()
		if details != "" {
			// Render details with better visibility (not faint)
			b.WriteString(details)
			b.WriteString("\n")
		} else {
			// Placeholder if no details
			b.WriteString(styles.FaintStyle.Render("No details available"))
			b.WriteString("\n")
		}
	}

	b.WriteString(strings.Repeat("━", 60))

	// Instructions
	b.WriteString("\n\n")
	b.WriteString(styles.FaintStyle.Render("Use ↑↓ or j/k to navigate, Enter to select, Esc to cancel"))

	return b.String()
}

func (s *SelectComponent) Selected() bool {
	return s.selected
}

func (s *SelectComponent) Cancelled() bool {
	return s.cancelled
}

func (s *SelectComponent) SelectedItem() SelectItem {
	if s.selected && s.cursor >= 0 && s.cursor < len(s.items) {
		return s.items[s.cursor]
	}
	return nil
}

func (s *SelectComponent) SelectedIndex() int {
	if s.selected {
		return s.cursor
	}
	return -1
}

package components

import (
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// TextViewerComponent provides a scrollable text viewer with viewport
type TextViewerComponent struct {
	viewport viewport.Model
	ready    bool
	title    string
}

// NewTextViewerComponent creates a new text viewer with the given content and title
func NewTextViewerComponent(title, content string) *TextViewerComponent {
	// Start with reasonable defaults - will be adjusted on WindowSizeMsg
	vp := viewport.New(80, 20)
	vp.SetContent(content)

	return &TextViewerComponent{
		viewport: vp,
		ready:    false,
		title:    title,
	}
}

func (c *TextViewerComponent) Init() tea.Cmd {
	return nil
}

func (c *TextViewerComponent) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Reserve space for title (1 line) and help text (2 lines)
		headerHeight := 1
		footerHeight := 2
		verticalMargin := headerHeight + footerHeight

		c.viewport.Width = msg.Width
		c.viewport.Height = msg.Height - verticalMargin
		c.ready = true

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "enter":
			return c, tea.Quit
		}
	}

	c.viewport, cmd = c.viewport.Update(msg)
	return c, cmd
}

func (c *TextViewerComponent) View() string {
	if !c.ready {
		return "Loading..."
	}

	var b strings.Builder

	// Title
	if c.title != "" {
		titleStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("86")) // Cyan
		b.WriteString(titleStyle.Render(c.title))
		b.WriteString("\n")
	}

	// Viewport content
	b.WriteString(c.viewport.View())
	b.WriteString("\n")

	// Help text
	helpStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241")) // Faint gray
	help := helpStyle.Render("↑/↓/PgUp/PgDn: scroll • q/esc/enter: exit")
	b.WriteString(help)

	return b.String()
}

// SetContent updates the viewport content
func (c *TextViewerComponent) SetContent(content string) {
	c.viewport.SetContent(content)
}

// ShowTextViewer displays a scrollable text viewer and waits for user to exit
func ShowTextViewer(title, content string) error {
	p := tea.NewProgram(NewTextViewerComponent(title, content))
	_, err := p.Run()
	return err
}

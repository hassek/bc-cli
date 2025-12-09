package components

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hassek/bc-cli/tui/styles"
)

// InputComponent handles validated number input
type InputComponent struct {
	textInput   textinput.Model
	label       string
	min         int
	max         int
	defaultVal  int
	value       int
	err         error
	submitted   bool
	cancelled   bool
}

func NewInputComponent(label string, min, max, defaultVal int) *InputComponent {
	ti := textinput.New()
	ti.Placeholder = fmt.Sprintf("%d", defaultVal)
	ti.Focus()
	ti.CharLimit = 10
	ti.Width = 20

	return &InputComponent{
		textInput:  ti,
		label:      label,
		min:        min,
		max:        max,
		defaultVal: defaultVal,
	}
}

func (i *InputComponent) Init() tea.Cmd {
	return textinput.Blink
}

func (i *InputComponent) Update(msg tea.Msg) (*InputComponent, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			i.cancelled = true
			return i, tea.Quit

		case "enter":
			// Validate and submit
			value := strings.TrimSpace(i.textInput.Value())
			if value == "" {
				i.value = i.defaultVal
				i.submitted = true
				return i, tea.Quit
			}

			num, err := strconv.Atoi(value)
			if err != nil {
				i.err = fmt.Errorf("please enter a valid number")
				return i, nil
			}

			if num < i.min || num > i.max {
				i.err = fmt.Errorf("please enter a number between %d and %d", i.min, i.max)
				return i, nil
			}

			i.value = num
			i.submitted = true
			i.err = nil
			return i, tea.Quit
		}
	}

	i.textInput, cmd = i.textInput.Update(msg)
	return i, cmd
}

func (i *InputComponent) View() string {
	var b strings.Builder

	// Label
	b.WriteString(styles.ActiveStyle.Render(i.label))
	b.WriteString("\n\n")

	// Input field
	b.WriteString(i.textInput.View())
	b.WriteString("\n")

	// Error message
	if i.err != nil {
		b.WriteString("\n")
		b.WriteString(styles.ErrorStyle.Render("✗ " + i.err.Error()))
		b.WriteString("\n")
	}

	// Help text
	b.WriteString("\n")
	helpText := fmt.Sprintf("Enter a number between %d and %d (default: %d)", i.min, i.max, i.defaultVal)
	b.WriteString(styles.FaintStyle.Render(helpText))
	b.WriteString("\n")
	b.WriteString(styles.FaintStyle.Render("Press Enter to confirm, Esc to cancel"))

	return b.String()
}

func (i *InputComponent) Submitted() bool {
	return i.submitted
}

func (i *InputComponent) Cancelled() bool {
	return i.cancelled
}

func (i *InputComponent) Value() int {
	return i.value
}

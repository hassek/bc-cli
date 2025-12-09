package prompts

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hassek/bc-cli/tui/components"
)

// Model that composes duck with input/confirm components
type promptModel struct {
	duck      *components.DuckComponent
	input     *components.InputComponent
	confirm   *components.ConfirmComponent
	inputMode bool
}

func (m promptModel) Init() tea.Cmd {
	var cmds []tea.Cmd
	cmds = append(cmds, m.duck.Init())
	if m.inputMode {
		cmds = append(cmds, m.input.Init())
	} else {
		cmds = append(cmds, m.confirm.Init())
	}
	return tea.Batch(cmds...)
}

func (m promptModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	// Update duck (handles tick messages)
	var duckCmd tea.Cmd
	m.duck, duckCmd = m.duck.Update(msg)
	if duckCmd != nil {
		cmds = append(cmds, duckCmd)
	}

	// Update input or confirm component
	if m.inputMode {
		var inputCmd tea.Cmd
		m.input, inputCmd = m.input.Update(msg)
		if inputCmd != nil {
			cmds = append(cmds, inputCmd)
		}

		// Trigger duck action on submit
		if m.input.Submitted() {
			m.duck.TriggerAction()
		}
	} else {
		var confirmCmd tea.Cmd
		m.confirm, confirmCmd = m.confirm.Update(msg)
		if confirmCmd != nil {
			cmds = append(cmds, confirmCmd)
		}

		// Trigger duck action on submit
		if m.confirm.Submitted() {
			m.duck.TriggerAction()
		}
	}

	return m, tea.Batch(cmds...)
}

func (m promptModel) View() string {
	var view string
	view += m.duck.View()
	if m.inputMode {
		view += m.input.View()
	} else {
		view += m.confirm.View()
	}
	return view
}

// PromptQuantityInt prompts the user to enter a quantity as an integer
// Matches the old API: func PromptQuantityInt(label string, min, max, defaultVal int) (int, error)
func PromptQuantityInt(label string, min, max, defaultVal int) (int, error) {
	duck := components.NewDuckComponent()
	input := components.NewInputComponent(label, min, max, defaultVal)

	m := promptModel{
		duck:      duck,
		input:     input,
		inputMode: true,
	}

	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		return 0, err
	}

	result := finalModel.(promptModel)
	if result.input.Cancelled() {
		return 0, ErrUserCancelled
	}

	return result.input.Value(), nil
}

// PromptConfirm prompts the user for a yes/no confirmation
// Matches the old API: func PromptConfirm(label string) (bool, error)
func PromptConfirm(label string) (bool, error) {
	duck := components.NewDuckComponent()
	confirm := components.NewConfirmComponent(label)

	m := promptModel{
		duck:      duck,
		confirm:   confirm,
		inputMode: false,
	}

	p := tea.NewProgram(m)
	finalModel, err := p.Run()
	if err != nil {
		return false, err
	}

	result := finalModel.(promptModel)
	if result.confirm.Cancelled() {
		return false, nil
	}

	return result.confirm.Result(), nil
}

// ErrUserCancelled is returned when the user cancels the prompt
type userCancelledError struct{}

func (userCancelledError) Error() string {
	return "user cancelled"
}

var ErrUserCancelled = userCancelledError{}

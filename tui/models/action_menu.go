package models

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hassek/bc-cli/tui/components"
)

// ActionItem represents a management action
type ActionItem struct {
	Action  string
	Display string
}

func (a ActionItem) Label() string {
	return a.Display
}

func (a ActionItem) Details() string {
	return ""
}

// ActionMenuModel composes duck + select for management actions
type ActionMenuModel struct {
	duck     *components.DuckComponent
	selector *components.SelectComponent
}

func NewActionMenuModel(actions []ActionItem) ActionMenuModel {
	items := make([]components.SelectItem, len(actions))
	for i, action := range actions {
		items[i] = action
	}

	return ActionMenuModel{
		duck:     components.NewDuckComponent(),
		selector: components.NewSelectComponent("What would you like to do?", items),
	}
}

func (m ActionMenuModel) Init() tea.Cmd {
	var cmds []tea.Cmd
	cmds = append(cmds, m.duck.Init())
	cmds = append(cmds, m.selector.Init())
	return tea.Batch(cmds...)
}

func (m ActionMenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	// Update duck (handles tick messages)
	var duckCmd tea.Cmd
	m.duck, duckCmd = m.duck.Update(msg)
	if duckCmd != nil {
		cmds = append(cmds, duckCmd)
	}

	// Update selector (handles key messages)
	var selectCmd tea.Cmd
	m.selector, selectCmd = m.selector.Update(msg)
	if selectCmd != nil {
		cmds = append(cmds, selectCmd)
	}

	// Trigger duck action on selection
	if m.selector.Selected() {
		m.duck.TriggerAction()
	}

	return m, tea.Batch(cmds...)
}

func (m ActionMenuModel) View() string {
	return m.duck.View() + m.selector.View()
}

// SelectAction shows the action menu and returns the selected action
func SelectAction(actions []ActionItem) (string, error) {
	p := tea.NewProgram(NewActionMenuModel(actions))
	model, err := p.Run()
	if err != nil {
		return "", err
	}

	m := model.(ActionMenuModel)
	if m.selector.Cancelled() {
		return "", nil
	}

	selectedItem := m.selector.SelectedItem()
	if selectedItem == nil {
		return "", nil
	}

	actionItem := selectedItem.(ActionItem)
	return actionItem.Action, nil
}

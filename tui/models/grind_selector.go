package models

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hassek/bc-cli/tui/components"
)

// GrindOption represents a grind type option
type GrindOption struct {
	Value   string
	Display string
}

func (g GrindOption) Label() string {
	return g.Display
}

func (g GrindOption) Details() string {
	return ""
}

// GrindSelectorModel composes duck + select for grind type selection
type GrindSelectorModel struct {
	duck     *components.DuckComponent
	selector *components.SelectComponent
}

func NewGrindSelectorModel() GrindSelectorModel {
	options := []GrindOption{
		{"whole_bean", "Whole Bean (I'll grind it myself)"},
		{"ground", "Ground (We'll grind it for you)"},
	}

	items := make([]components.SelectItem, len(options))
	for i, opt := range options {
		items[i] = opt
	}

	return GrindSelectorModel{
		duck:     components.NewDuckComponent(),
		selector: components.NewSelectComponent("Grind type", items),
	}
}

func (m GrindSelectorModel) Init() tea.Cmd {
	var cmds []tea.Cmd
	cmds = append(cmds, m.duck.Init())
	cmds = append(cmds, m.selector.Init())
	return tea.Batch(cmds...)
}

func (m GrindSelectorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m GrindSelectorModel) View() string {
	return m.duck.View() + m.selector.View()
}

// SelectGrindType shows the grind selector and returns the selected grind type
func SelectGrindType() (string, error) {
	p := tea.NewProgram(NewGrindSelectorModel())
	model, err := p.Run()
	if err != nil {
		return "", err
	}

	m := model.(GrindSelectorModel)
	if m.selector.Cancelled() {
		return "", nil
	}

	selectedItem := m.selector.SelectedItem()
	if selectedItem == nil {
		return "", nil
	}

	grindOpt := selectedItem.(GrindOption)
	return grindOpt.Value, nil
}

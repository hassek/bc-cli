package models

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hassek/bc-cli/api"
	"github.com/hassek/bc-cli/tui/components"
)

// SectionItem wraps an api.Section for use with SelectComponent
type SectionItem struct {
	Section api.Section
	IsBack  bool
}

func (s SectionItem) Label() string {
	if s.IsBack {
		return "← Back"
	}
	return s.Section.Name
}

func (s SectionItem) Details() string {
	if s.IsBack {
		return "Return to categories"
	}
	return s.Section.Description
}

// SectionPickerModel composes duck + select for section browsing
type SectionPickerModel struct {
	duck     *components.DuckComponent
	selector *components.SelectComponent
}

func NewSectionPickerModel(sections []api.Section) SectionPickerModel {
	// Convert sections to SelectItems
	items := make([]components.SelectItem, len(sections)+1)
	for i, section := range sections {
		items[i] = SectionItem{Section: section}
	}
	// Add back option
	items[len(sections)] = SectionItem{IsBack: true}

	return SectionPickerModel{
		duck:     components.NewDuckComponent(),
		selector: components.NewSelectComponent("Select a section", items),
	}
}

func (m SectionPickerModel) Init() tea.Cmd {
	var cmds []tea.Cmd
	cmds = append(cmds, m.duck.Init())
	cmds = append(cmds, m.selector.Init())
	return tea.Batch(cmds...)
}

func (m SectionPickerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	var duckCmd tea.Cmd
	m.duck, duckCmd = m.duck.Update(msg)
	if duckCmd != nil {
		cmds = append(cmds, duckCmd)
	}

	var selectCmd tea.Cmd
	m.selector, selectCmd = m.selector.Update(msg)
	if selectCmd != nil {
		cmds = append(cmds, selectCmd)
	}

	if m.selector.Selected() {
		m.duck.TriggerAction()
	}

	return m, tea.Batch(cmds...)
}

func (m SectionPickerModel) View() string {
	return m.duck.View() + m.selector.View()
}

// PickSection returns selected section or nil if back/cancelled
func PickSection(sections []api.Section) (*api.Section, error) {
	p := tea.NewProgram(NewSectionPickerModel(sections))
	model, err := p.Run()
	if err != nil {
		return nil, err
	}

	m := model.(SectionPickerModel)
	if m.selector.Cancelled() {
		return nil, nil
	}

	selectedItem := m.selector.SelectedItem()
	if selectedItem == nil {
		return nil, nil
	}

	sectionItem := selectedItem.(SectionItem)
	if sectionItem.IsBack {
		return nil, nil
	}

	return &sectionItem.Section, nil
}

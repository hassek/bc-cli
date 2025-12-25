package models

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hassek/bc-cli/api"
	"github.com/hassek/bc-cli/tui/components"
)

// CategoryItem wraps an api.Category for use with SelectComponent
type CategoryItem struct {
	Category          api.Category
	IsBookmarksOption bool // Special option to view bookmarks
	IsExit            bool
}

func (c CategoryItem) Label() string {
	if c.IsExit {
		return "← Exit"
	}
	if c.IsBookmarksOption {
		return "★ My Bookmarks"
	}
	return c.Category.Name
}

func (c CategoryItem) Details() string {
	if c.IsExit {
		return "Return to main menu"
	}
	if c.IsBookmarksOption {
		return "View your saved articles"
	}
	return c.Category.Description
}

// CategoryPickerModel composes duck + select for category browsing
type CategoryPickerModel struct {
	duck     *components.DuckComponent
	selector *components.SelectComponent
}

func NewCategoryPickerModel(categories []api.Category, showBookmarks bool) CategoryPickerModel {
	itemsCount := len(categories)
	if showBookmarks {
		itemsCount++ // Add bookmarks option
	}
	itemsCount++ // Add exit option

	items := make([]components.SelectItem, 0, itemsCount)

	// Add bookmarks option first if authenticated
	if showBookmarks {
		items = append(items, CategoryItem{IsBookmarksOption: true})
	}

	// Add categories
	for _, cat := range categories {
		items = append(items, CategoryItem{Category: cat})
	}

	// Add exit
	items = append(items, CategoryItem{IsExit: true})

	return CategoryPickerModel{
		duck:     components.NewDuckComponent(),
		selector: components.NewSelectComponent("Explore Coffee Knowledge", items),
	}
}

func (m CategoryPickerModel) Init() tea.Cmd {
	var cmds []tea.Cmd
	cmds = append(cmds, m.duck.Init())
	cmds = append(cmds, m.selector.Init())
	return tea.Batch(cmds...)
}

func (m CategoryPickerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m CategoryPickerModel) View() string {
	return m.duck.View() + m.selector.View()
}

// PickCategory returns selected category, special flag for bookmarks, or error
// Second return value is true if user selected bookmarks option
func PickCategory(categories []api.Category, showBookmarks bool) (*api.Category, bool, error) {
	p := tea.NewProgram(NewCategoryPickerModel(categories, showBookmarks))
	model, err := p.Run()
	if err != nil {
		return nil, false, err
	}

	m := model.(CategoryPickerModel)
	if m.selector.Cancelled() {
		return nil, false, nil
	}

	selectedItem := m.selector.SelectedItem()
	if selectedItem == nil {
		return nil, false, nil
	}

	catItem := selectedItem.(CategoryItem)
	if catItem.IsExit {
		return nil, false, nil
	}
	if catItem.IsBookmarksOption {
		return nil, true, nil // Signal to show bookmarks
	}

	return &catItem.Category, false, nil
}

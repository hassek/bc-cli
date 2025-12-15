package models

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hassek/bc-cli/api"
	"github.com/hassek/bc-cli/tui/components"
	"github.com/hassek/bc-cli/utils"
)

// ProductItem wraps an api.AvailableSubscription for use with SelectComponent
type ProductItem struct {
	Product api.AvailableSubscription
	IsExit  bool
}

func (p ProductItem) Label() string {
	if p.IsExit {
		return "← Exit"
	}
	return p.Product.Name
}

func (p ProductItem) Details() string {
	if p.IsExit {
		return "Return to main menu"
	}

	// Get terminal width and calculate available space for details
	termWidth := utils.GetTerminalWidth()
	detailsWidth := termWidth - 10 // Account for margins and duck
	if detailsWidth > 100 {
		detailsWidth = 100 // Cap at 100 for readability
	}

	details := fmt.Sprintf("Name:    %s\n", p.Product.Name)
	details += fmt.Sprintf("Price:   %s %s\n", p.Product.Currency, p.Product.Price)

	// Wrap summary to fit available space
	indentWidth := 9 // "Summary: " is 9 chars
	wrapWidth := detailsWidth - indentWidth
	if wrapWidth < 20 {
		wrapWidth = 20 // Minimum width
	}
	wrappedSummary := utils.WrapTextWithIndent(p.Product.Summary, wrapWidth, "         ")
	details += fmt.Sprintf("Summary: %s", wrappedSummary)
	return details
}

// ProductPickerModel composes duck + select for product browsing
type ProductPickerModel struct {
	duck     *components.DuckComponent
	selector *components.SelectComponent
}

func NewProductPickerModel(products []api.AvailableSubscription) ProductPickerModel {
	// Convert products to SelectItems
	items := make([]components.SelectItem, len(products)+1)
	for i, product := range products {
		items[i] = ProductItem{Product: product}
	}
	// Add exit option
	items[len(products)] = ProductItem{IsExit: true}

	return ProductPickerModel{
		duck:     components.NewDuckComponent(),
		selector: components.NewSelectComponent("Select a product to learn more", items),
	}
}

func (m ProductPickerModel) Init() tea.Cmd {
	var cmds []tea.Cmd
	cmds = append(cmds, m.duck.Init())
	cmds = append(cmds, m.selector.Init())
	return tea.Batch(cmds...)
}

func (m ProductPickerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m ProductPickerModel) View() string {
	return m.duck.View() + m.selector.View()
}

// PickProduct shows the product picker and returns the selected product
// Returns nil if user cancelled or selected exit
func PickProduct(products []api.AvailableSubscription) (*api.AvailableSubscription, error) {
	p := tea.NewProgram(NewProductPickerModel(products))
	model, err := p.Run()
	if err != nil {
		return nil, err
	}

	m := model.(ProductPickerModel)
	if m.selector.Cancelled() {
		return nil, nil
	}

	selectedItem := m.selector.SelectedItem()
	if selectedItem == nil {
		return nil, nil
	}

	productItem := selectedItem.(ProductItem)
	if productItem.IsExit {
		return nil, nil
	}

	return &productItem.Product, nil
}

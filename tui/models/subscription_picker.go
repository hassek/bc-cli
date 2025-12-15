package models

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hassek/bc-cli/api"
	"github.com/hassek/bc-cli/tui/components"
	"github.com/hassek/bc-cli/utils"
)

// SubscriptionItem wraps an api.AvailableSubscription for use with SelectComponent
type SubscriptionItem struct {
	Subscription api.AvailableSubscription
	IsExit       bool
}

func (s SubscriptionItem) Label() string {
	if s.IsExit {
		return "← Exit"
	}
	return s.Subscription.Name
}

func (s SubscriptionItem) Details() string {
	if s.IsExit {
		return "Return to main menu"
	}
	details := fmt.Sprintf("Name:    %s\n", s.Subscription.Name)
	details += fmt.Sprintf("Price:   %s %s/%s\n", s.Subscription.Currency, s.Subscription.Price, s.Subscription.BillingPeriod)

	// Wrap summary to fit in details panel (60 chars wide, minus label)
	wrappedSummary := utils.WrapTextWithIndent(s.Subscription.Summary, 55, "         ")
	details += fmt.Sprintf("Summary: %s\n", wrappedSummary)

	if len(s.Subscription.Features) > 0 {
		details += "\nFeatures:\n"
		for _, feature := range s.Subscription.Features {
			// Wrap each feature to fit (60 chars wide, minus bullet indent)
			wrappedFeature := utils.WrapTextWithIndent(feature, 56, "    ")
			details += fmt.Sprintf("  • %s\n", wrappedFeature)
		}
	}

	return details
}

// SubscriptionPickerModel composes duck + select for subscription browsing
type SubscriptionPickerModel struct {
	duck     *components.DuckComponent
	selector *components.SelectComponent
}

func NewSubscriptionPickerModel(subscriptions []api.AvailableSubscription) SubscriptionPickerModel {
	// Convert subscriptions to SelectItems
	items := make([]components.SelectItem, len(subscriptions)+1)
	for i, subscription := range subscriptions {
		items[i] = SubscriptionItem{Subscription: subscription}
	}
	// Add exit option
	items[len(subscriptions)] = SubscriptionItem{IsExit: true}

	return SubscriptionPickerModel{
		duck:     components.NewDuckComponent(),
		selector: components.NewSelectComponent("Select a subscription tier to learn more", items),
	}
}

func (m SubscriptionPickerModel) Init() tea.Cmd {
	var cmds []tea.Cmd
	cmds = append(cmds, m.duck.Init())
	cmds = append(cmds, m.selector.Init())
	return tea.Batch(cmds...)
}

func (m SubscriptionPickerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m SubscriptionPickerModel) View() string {
	return m.duck.View() + m.selector.View()
}

// PickSubscription shows the subscription picker and returns the selected subscription
// Returns nil if user cancelled or selected exit
func PickSubscription(subscriptions []api.AvailableSubscription) (*api.AvailableSubscription, error) {
	p := tea.NewProgram(NewSubscriptionPickerModel(subscriptions))
	model, err := p.Run()
	if err != nil {
		return nil, err
	}

	m := model.(SubscriptionPickerModel)
	if m.selector.Cancelled() {
		return nil, nil
	}

	selectedItem := m.selector.SelectedItem()
	if selectedItem == nil {
		return nil, nil
	}

	subItem := selectedItem.(SubscriptionItem)
	if subItem.IsExit {
		return nil, nil
	}

	return &subItem.Subscription, nil
}

package models

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hassek/bc-cli/api"
	"github.com/hassek/bc-cli/tui/components"
)

// ManageSubscriptionItem wraps a subscription for management
type ManageSubscriptionItem struct {
	Subscription    api.Subscription
	Display         string
	Status          string
	StartedAt       string
	ExpiresAt       string
	TotalQuantity   int
	HasOrderDetails bool
	IsExit          bool
}

func (m ManageSubscriptionItem) Label() string {
	return m.Display
}

func (m ManageSubscriptionItem) Details() string {
	if m.IsExit {
		return "Return to main menu"
	}

	details := ""
	if m.Subscription.ID != "" {
		details += fmt.Sprintf("Tier:     %s\n", m.Subscription.Tier)
		details += fmt.Sprintf("Status:   %s\n", m.Status)
		if m.StartedAt != "" {
			details += fmt.Sprintf("Started:  %s\n", m.StartedAt)
		}
		if m.HasOrderDetails {
			details += fmt.Sprintf("Quantity: %d/month", m.TotalQuantity)
		}
	}
	return details
}

// ManageSubscriptionPickerModel composes duck + select for subscription management
type ManageSubscriptionPickerModel struct {
	duck     *components.DuckComponent
	selector *components.SelectComponent
}

func NewManageSubscriptionPickerModel(subscriptions []ManageSubscriptionItem) ManageSubscriptionPickerModel {
	items := make([]components.SelectItem, len(subscriptions))
	for i, sub := range subscriptions {
		items[i] = sub
	}

	return ManageSubscriptionPickerModel{
		duck:     components.NewDuckComponent(),
		selector: components.NewSelectComponent("Select a subscription to manage", items),
	}
}

func (m ManageSubscriptionPickerModel) Init() tea.Cmd {
	var cmds []tea.Cmd
	cmds = append(cmds, m.duck.Init())
	cmds = append(cmds, m.selector.Init())
	return tea.Batch(cmds...)
}

func (m ManageSubscriptionPickerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m ManageSubscriptionPickerModel) View() string {
	return m.duck.View() + m.selector.View()
}

// PickManageSubscription shows the subscription picker for management
func PickManageSubscription(subscriptions []ManageSubscriptionItem) (*api.Subscription, error) {
	p := tea.NewProgram(NewManageSubscriptionPickerModel(subscriptions))
	model, err := p.Run()
	if err != nil {
		return nil, err
	}

	m := model.(ManageSubscriptionPickerModel)
	if m.selector.Cancelled() {
		return nil, nil
	}

	selectedItem := m.selector.SelectedItem()
	if selectedItem == nil {
		return nil, nil
	}

	subItem := selectedItem.(ManageSubscriptionItem)
	if subItem.IsExit {
		return nil, nil
	}

	return &subItem.Subscription, nil
}

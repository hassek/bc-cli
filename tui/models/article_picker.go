package models

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hassek/bc-cli/api"
	"github.com/hassek/bc-cli/tui/components"
	"github.com/hassek/bc-cli/utils"
)

// ArticleItem wraps an api.Article for use with SelectComponent
type ArticleItem struct {
	Article api.Article
	IsBack  bool
}

func (a ArticleItem) Label() string {
	if a.IsBack {
		return "← Back"
	}
	label := a.Article.Title
	if a.Article.IsBookmarked {
		label = "★ " + label // Show bookmark indicator
	}
	return label
}

func (a ArticleItem) Details() string {
	if a.IsBack {
		return "Return to previous menu"
	}

	// Get terminal width and calculate available space for details
	termWidth := utils.GetTerminalWidth()
	detailsWidth := termWidth - 10 // Account for margins and duck
	if detailsWidth > 100 {
		detailsWidth = 100 // Cap at 100 for readability
	}

	var details strings.Builder

	// Wrap summary to fit available space
	if a.Article.Summary != "" {
		indentWidth := 9 // "Summary: " is 9 chars
		wrapWidth := detailsWidth - indentWidth
		if wrapWidth < 20 {
			wrapWidth = 20 // Minimum width
		}
		wrappedSummary := utils.WrapTextWithIndent(a.Article.Summary, wrapWidth, "         ")
		details.WriteString(fmt.Sprintf("Summary: %s\n", wrappedSummary))
	}

	// Reading time
	if a.Article.ReadTime > 0 {
		details.WriteString(fmt.Sprintf("Reading time: %d min\n", a.Article.ReadTime))
	}

	// Tags
	if a.Article.Tags != "" {
		details.WriteString(fmt.Sprintf("Tags: %s\n", a.Article.Tags))
	}

	// Bookmark status
	if a.Article.IsBookmarked {
		details.WriteString("★ Bookmarked")
	}

	return strings.TrimSpace(details.String())
}

// ArticlePickerModel composes duck + select for article browsing
type ArticlePickerModel struct {
	duck     *components.DuckComponent
	selector *components.SelectComponent
}

func NewArticlePickerModel(articles []api.Article) ArticlePickerModel {
	// Convert articles to SelectItems
	items := make([]components.SelectItem, len(articles)+1)
	for i, article := range articles {
		items[i] = ArticleItem{Article: article}
	}
	// Add back option
	items[len(articles)] = ArticleItem{IsBack: true}

	return ArticlePickerModel{
		duck:     components.NewDuckComponent(),
		selector: components.NewSelectComponent("Select an article to read", items),
	}
}

func (m ArticlePickerModel) Init() tea.Cmd {
	var cmds []tea.Cmd
	cmds = append(cmds, m.duck.Init())
	cmds = append(cmds, m.selector.Init())
	return tea.Batch(cmds...)
}

func (m ArticlePickerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m ArticlePickerModel) View() string {
	return m.duck.View() + m.selector.View()
}

// PickArticle returns selected article or nil if back/cancelled
func PickArticle(articles []api.Article) (*api.Article, error) {
	p := tea.NewProgram(NewArticlePickerModel(articles))
	model, err := p.Run()
	if err != nil {
		return nil, err
	}

	m := model.(ArticlePickerModel)
	if m.selector.Cancelled() {
		return nil, nil
	}

	selectedItem := m.selector.SelectedItem()
	if selectedItem == nil {
		return nil, nil
	}

	articleItem := selectedItem.(ArticleItem)
	if articleItem.IsBack {
		return nil, nil
	}

	return &articleItem.Article, nil
}

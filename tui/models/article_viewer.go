package models

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hassek/bc-cli/api"
	"github.com/hassek/bc-cli/templates"
	"github.com/hassek/bc-cli/tui/styles"
)

// ArticleAction represents what action was taken in the article viewer
type ArticleAction int

const (
	ArticleActionNone ArticleAction = iota
	ArticleActionToggleBookmark
	ArticleActionShowRelated
	ArticleActionBack
	ArticleActionQuit // Ctrl+C - quit entire program
)

// ArticleViewerModel displays an article with scrollable content and actions
type ArticleViewerModel struct {
	viewport    viewport.Model
	article     *api.Article
	ready       bool
	canBookmark bool
	lastAction  ArticleAction
}

func NewArticleViewerModel(article *api.Article, canBookmark bool) *ArticleViewerModel {
	vp := viewport.New(80, 20)

	// Render article content using template
	content, err := templates.RenderArticleContent(article)
	if err != nil {
		content = fmt.Sprintf("Error rendering article: %v", err)
	}
	vp.SetContent(content)

	return &ArticleViewerModel{
		viewport:    vp,
		article:     article,
		canBookmark: canBookmark,
		lastAction:  ArticleActionNone,
	}
}

func (m *ArticleViewerModel) Init() tea.Cmd {
	return nil
}

func (m *ArticleViewerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		headerHeight := 3
		footerHeight := 2
		verticalMargins := headerHeight + footerHeight

		if !m.ready {
			m.viewport = viewport.New(msg.Width, msg.Height-verticalMargins)
			m.viewport.YPosition = headerHeight

			// Re-render content with current article
			content, err := templates.RenderArticleContent(m.article)
			if err != nil {
				content = fmt.Sprintf("Error rendering article: %v", err)
			}
			m.viewport.SetContent(content)
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMargins
		}

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			// Force quit to terminal
			m.lastAction = ArticleActionQuit
			return m, tea.Quit

		case "b":
			// Toggle bookmark
			if m.canBookmark {
				m.lastAction = ArticleActionToggleBookmark
				return m, tea.Quit
			}
		case "r":
			// Show related articles
			m.lastAction = ArticleActionShowRelated
			return m, tea.Quit
		case "q", "esc":
			// Back
			m.lastAction = ArticleActionBack
			return m, tea.Quit
		// Vim-like navigation
		case "g":
			// Go to top (gg in vim, but we use single 'g' for simplicity)
			m.viewport.GotoTop()
			return m, nil
		case "G":
			// Go to bottom (Shift+G)
			m.viewport.GotoBottom()
			return m, nil
		case "d", "ctrl+d":
			// Half page down
			m.viewport.HalfPageDown()
			return m, nil
		case "u", "ctrl+u":
			// Half page up
			m.viewport.HalfPageUp()
			return m, nil
		case "f", "ctrl+f", "pgdown":
			// Full page down
			m.viewport.PageDown()
			return m, nil
		case "ctrl+b", "pgup":
			// Full page up
			m.viewport.PageUp()
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m *ArticleViewerModel) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}

	var b strings.Builder

	// Header with title and bookmark status
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(styles.Cyan)
	bookmarkIndicator := ""
	if m.article.IsBookmarked {
		bookmarkIndicator = " " + styles.DuckAccentStyle.Render("★")
	}
	b.WriteString(titleStyle.Render(m.article.Title) + bookmarkIndicator + "\n\n")

	// Viewport content
	b.WriteString(m.viewport.View())
	b.WriteString("\n\n")

	// Help text footer
	helpParts := []string{"↑/↓/j/k: scroll", "g/G: top/bottom", "d/u: half page"}
	if m.canBookmark {
		if m.article.IsBookmarked {
			helpParts = append(helpParts, "b: remove bookmark")
		} else {
			helpParts = append(helpParts, "b: bookmark")
		}
	}
	helpParts = append(helpParts, "r: related", "q: back")
	helpText := strings.Join(helpParts, " • ")
	b.WriteString(styles.FaintStyle.Render(helpText))

	return b.String()
}

// ViewArticle displays an article and returns the action taken
func ViewArticle(article *api.Article, canBookmark bool) (ArticleAction, error) {
	p := tea.NewProgram(NewArticleViewerModel(article, canBookmark))
	model, err := p.Run()
	if err != nil {
		return ArticleActionNone, err
	}

	m := model.(*ArticleViewerModel)
	return m.lastAction, nil
}

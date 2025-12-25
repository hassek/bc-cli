package templates

import (
	"fmt"

	"github.com/hassek/bc-cli/api"
	"github.com/hassek/bc-cli/utils"
)

const ArticleContentTemplate = `{{ if .Author }}{{ faint "By" }} {{ cyan .Author }}{{ if .ReadTime }} {{ faint "•" }} {{ .ReadTime }} {{ faint "min read" }}{{ end }}{{ if .PublishedAt }} {{ faint "•" }} {{ .PublishedAt }}{{ end }}
{{ else }}{{ if .ReadTime }}{{ .ReadTime }} {{ faint "min read" }}{{ if .PublishedAt }} {{ faint "•" }} {{ .PublishedAt }}{{ end }}{{ end }}{{ end }}

{{ if .Tags }}{{ faint "Tags:" }} {{ cyan .Tags }}

{{ end }}{{ .Content }}
`

// RenderArticleContent renders an article with formatting for terminal display
func RenderArticleContent(article *api.Article) (string, error) {
	// Pre-render markdown content to support template syntax
	renderedContent := RenderDescription(article.Content)

	// Format published date if available
	var publishedAt string
	if article.PublishedAt != nil && *article.PublishedAt != "" {
		publishedAt = utils.FormatTimestamp(*article.PublishedAt)
	}

	data := struct {
		Title       string
		Author      string
		ReadTime    string
		PublishedAt string
		Tags        string
		Content     string
	}{
		Title:       article.Title,
		Author:      article.Author,
		ReadTime:    fmt.Sprintf("%d", article.ReadTime),
		PublishedAt: publishedAt,
		Tags:        article.Tags,
		Content:     renderedContent,
	}

	return RenderToString(ArticleContentTemplate, data)
}

package api

import (
	"fmt"
)

// Category represents a content category
type Category struct {
	ID          string  `json:"id"`
	Slug        string  `json:"slug"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Order       int     `json:"order"`
	PublishedAt *string `json:"published_at"`
}

// Section represents a section within a category
type Section struct {
	ID          string  `json:"id"`
	CategoryID  string  `json:"category_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Order       int     `json:"order"`
	PublishedAt *string `json:"published_at"`
}

// Article represents a knowledge article
type Article struct {
	ID           string  `json:"id"`
	CategoryID   string  `json:"category_id"`
	SectionID    *string `json:"section_id"` // null if in default section
	Title        string  `json:"title"`
	Summary      string  `json:"summary"`
	Content      string  `json:"content"` // Full markdown content
	Author       string  `json:"author"`
	ReadTime     int     `json:"read_time"` // Minutes
	Tags         string  `json:"tags"`      // Comma-separated tags
	PublishedAt  *string `json:"published_at"`
	IsBookmarked bool    `json:"is_bookmarked"` // Only present for authenticated users
}

// Bookmark represents a user's bookmarked article
type Bookmark struct {
	ID        string  `json:"id"`
	ArticleID string  `json:"article_id"`
	Article   Article `json:"article"`
	CreatedAt string  `json:"created_at"`
}

// Response wrappers
type categoriesResponse struct {
	Meta struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"meta"`
	Data struct {
		Count    int        `json:"count"`
		Next     *string    `json:"next"`
		Previous *string    `json:"previous"`
		Results  []Category `json:"results"`
	} `json:"data"`
}

type categoryResponse struct {
	Meta struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"meta"`
	Data Category `json:"data"`
}

type sectionsResponse struct {
	Meta struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"meta"`
	Data []Section `json:"data"`
}

type articlesResponse struct {
	Meta struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"meta"`
	Data []Article `json:"data"`
}

type articleResponse struct {
	Meta struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"meta"`
	Data Article `json:"data"`
}

type bookmarksResponse struct {
	Meta struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"meta"`
	Data []Bookmark `json:"data"`
}

type bookmarkResponse struct {
	Meta struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"meta"`
	Data Bookmark `json:"data"`
}

// ListCategories retrieves all published categories
func (c *Client) ListCategories() ([]Category, error) {
	resp, err := c.doRequest("GET", "/api/core/v1/content/categories/", nil, false)
	if err != nil {
		return nil, err
	}

	var result categoriesResponse
	if err := c.handleResponse(resp, &result); err != nil {
		return nil, err
	}

	return result.Data.Results, nil
}

// GetCategory retrieves a category by slug
func (c *Client) GetCategory(slug string) (*Category, error) {
	url := fmt.Sprintf("/api/core/v1/content/categories/%s/", slug)
	resp, err := c.doRequest("GET", url, nil, false)
	if err != nil {
		return nil, err
	}

	var result categoryResponse
	if err := c.handleResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result.Data, nil
}

// ListCategorySections retrieves sections for a category
func (c *Client) ListCategorySections(categorySlug string) ([]Section, error) {
	url := fmt.Sprintf("/api/core/v1/content/categories/%s/sections/", categorySlug)
	resp, err := c.doRequest("GET", url, nil, false)
	if err != nil {
		return nil, err
	}

	var result sectionsResponse
	if err := c.handleResponse(resp, &result); err != nil {
		return nil, err
	}

	return result.Data, nil
}

// ListCategoryArticles retrieves articles in category's default section
func (c *Client) ListCategoryArticles(categorySlug string) ([]Article, error) {
	url := fmt.Sprintf("/api/core/v1/content/categories/%s/articles/", categorySlug)
	resp, err := c.doRequest("GET", url, nil, false)
	if err != nil {
		return nil, err
	}

	var result articlesResponse
	if err := c.handleResponse(resp, &result); err != nil {
		return nil, err
	}

	return result.Data, nil
}

// ListSectionArticles retrieves articles in a specific section
func (c *Client) ListSectionArticles(sectionID string) ([]Article, error) {
	url := fmt.Sprintf("/api/core/v1/content/sections/%s/articles/", sectionID)
	resp, err := c.doRequest("GET", url, nil, false)
	if err != nil {
		return nil, err
	}

	var result articlesResponse
	if err := c.handleResponse(resp, &result); err != nil {
		return nil, err
	}

	return result.Data, nil
}

// GetArticle retrieves full article with content
func (c *Client) GetArticle(articleID string) (*Article, error) {
	url := fmt.Sprintf("/api/core/v1/content/articles/%s/", articleID)
	resp, err := c.doRequest("GET", url, nil, false)
	if err != nil {
		return nil, err
	}

	var result articleResponse
	if err := c.handleResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result.Data, nil
}

// ListBookmarks retrieves user's bookmarked articles (requires auth)
func (c *Client) ListBookmarks() ([]Bookmark, error) {
	resp, err := c.doRequest("GET", "/api/core/v1/content/bookmarks/", nil, true)
	if err != nil {
		return nil, err
	}

	var result bookmarksResponse
	if err := c.handleResponse(resp, &result); err != nil {
		return nil, err
	}

	return result.Data, nil
}

// CreateBookmark bookmarks an article (requires auth)
func (c *Client) CreateBookmark(articleID string) (*Bookmark, error) {
	body := map[string]string{"article_id": articleID}
	resp, err := c.doRequest("POST", "/api/core/v1/content/bookmarks/", body, true)
	if err != nil {
		return nil, err
	}

	var result bookmarkResponse
	if err := c.handleResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result.Data, nil
}

// DeleteBookmark removes a bookmark (requires auth)
func (c *Client) DeleteBookmark(bookmarkID string) error {
	url := fmt.Sprintf("/api/core/v1/content/bookmarks/%s/", bookmarkID)
	resp, err := c.doRequest("DELETE", url, nil, true)
	if err != nil {
		return err
	}

	// DELETE returns 204 No Content on success
	if resp.StatusCode != 204 {
		return c.handleResponse(resp, nil)
	}

	return nil
}

// CategoryHasSections checks if a category has sections (helper for smart navigation)
func (c *Client) CategoryHasSections(categorySlug string) (bool, error) {
	sections, err := c.ListCategorySections(categorySlug)
	if err != nil {
		return false, err
	}
	return len(sections) > 0, nil
}

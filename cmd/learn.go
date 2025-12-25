package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"

	"github.com/hassek/bc-cli/api"
	"github.com/hassek/bc-cli/config"
	"github.com/hassek/bc-cli/tui/models"
	"github.com/spf13/cobra"
)

// ErrUserQuit is returned when user presses Ctrl+C to exit
var ErrUserQuit = errors.New("user quit")

var learnCmd = &cobra.Command{
	Use:   "learn",
	Short: "Explore coffee knowledge and articles",
	Long:  `Browse categories, sections, and articles about coffee. Save your favorites with bookmarks.`,
	RunE:  runLearn,
}

func init() {
	rootCmd.AddCommand(learnCmd)
}

func runLearn(cmd *cobra.Command, args []string) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	client := api.NewClient(cfg)

	// Main navigation loop
	for {
		// Fetch categories
		categories, err := client.ListCategories()
		if err != nil {
			return fmt.Errorf("failed to fetch categories: %w", err)
		}

		if len(categories) == 0 {
			fmt.Println("No content available at this time.")
			return nil
		}

		// Show category picker (with bookmarks option if authenticated)
		category, showBookmarks, err := models.PickCategory(categories, cfg.IsAuthenticated())
		if err != nil {
			return err
		}

		// User cancelled or selected exit
		if category == nil && !showBookmarks {
			return nil
		}

		// Handle bookmarks view
		if showBookmarks {
			if err := showBookmarksView(cfg, client); err != nil {
				fmt.Printf("\nError showing bookmarks: %v\n", err)
				fmt.Print("Press Enter to continue...")
				_, _ = bufio.NewReader(os.Stdin).ReadBytes('\n')
				continue
			}
			continue
		}

		// Navigate into category
		if err := navigateCategory(cfg, client, category); err != nil {
			if errors.Is(err, ErrUserQuit) {
				return nil // Exit cleanly
			}
			fmt.Printf("\nError: %v\n", err)
			fmt.Print("Press Enter to continue...")
			_, _ = bufio.NewReader(os.Stdin).ReadBytes('\n')
			continue
		}
	}
}

func navigateCategory(cfg *config.Config, client *api.Client, category *api.Category) error {
	// Smart detection: check if category has sections
	hasSections, err := client.CategoryHasSections(category.Slug)
	if err != nil {
		return err
	}

	if hasSections {
		// Show sections first
		return navigateSections(cfg, client, category)
	}
	// Show articles directly
	return navigateArticles(cfg, client, category.Slug, nil)
}

func navigateSections(cfg *config.Config, client *api.Client, category *api.Category) error {
	for {
		sections, err := client.ListCategorySections(category.Slug)
		if err != nil {
			return err
		}

		if len(sections) == 0 {
			fmt.Println("No sections available in this category.")
			return nil
		}

		section, err := models.PickSection(sections)
		if err != nil || section == nil {
			return err // Back or exit
		}

		// Navigate into section's articles
		if err := navigateArticles(cfg, client, category.Slug, &section.ID); err != nil {
			if errors.Is(err, ErrUserQuit) {
				return err // Propagate quit signal
			}
			fmt.Printf("\nError: %v\n", err)
			fmt.Print("Press Enter to continue...")
			_, _ = bufio.NewReader(os.Stdin).ReadBytes('\n')
			continue
		}
	}
}

func navigateArticles(cfg *config.Config, client *api.Client, categorySlug string, sectionID *string) error {
	for {
		var articles []api.Article
		var err error

		if sectionID == nil {
			// Category's default articles
			articles, err = client.ListCategoryArticles(categorySlug)
		} else {
			// Section's articles
			articles, err = client.ListSectionArticles(*sectionID)
		}

		if err != nil {
			return err
		}

		if len(articles) == 0 {
			fmt.Println("No articles available in this section.")
			return nil
		}

		article, err := models.PickArticle(articles)
		if err != nil || article == nil {
			return err // Back or exit
		}

		// Fetch full article with content
		fullArticle, err := client.GetArticle(article.ID)
		if err != nil {
			fmt.Printf("\nError loading article: %v\n", err)
			fmt.Print("Press Enter to continue...")
			_, _ = bufio.NewReader(os.Stdin).ReadBytes('\n')
			continue
		}

		// View article with actions
		action, err := viewArticleWithActions(cfg, client, fullArticle)
		if err != nil {
			fmt.Printf("\nError: %v\n", err)
			fmt.Print("Press Enter to continue...")
			_, _ = bufio.NewReader(os.Stdin).ReadBytes('\n')
			continue
		}

		// Handle post-view actions
		if action == models.ArticleActionQuit {
			// User pressed Ctrl+C - exit completely
			return ErrUserQuit
		}
		if action == models.ArticleActionShowRelated {
			// Loop continues to show articles list again (related articles from same section)
			continue
		}
	}
}

func viewArticleWithActions(cfg *config.Config, client *api.Client, article *api.Article) (models.ArticleAction, error) {
	canBookmark := cfg.IsAuthenticated()

	for {
		action, err := models.ViewArticle(article, canBookmark)
		if err != nil {
			return models.ArticleActionNone, err
		}

		switch action {
		case models.ArticleActionToggleBookmark:
			// Handle bookmark toggle
			if err := toggleBookmark(client, article); err != nil {
				fmt.Printf("\nError toggling bookmark: %v\n", err)
				fmt.Print("Press Enter to continue...")
				_, _ = bufio.NewReader(os.Stdin).ReadBytes('\n')
			} else {
				// Refresh article to get updated bookmark status
				updated, err := client.GetArticle(article.ID)
				if err == nil {
					article = updated
					if article.IsBookmarked {
						fmt.Println("\n✓ Article bookmarked!")
					} else {
						fmt.Println("\n✓ Bookmark removed")
					}
					fmt.Print("Press Enter to continue reading...")
					_, _ = bufio.NewReader(os.Stdin).ReadBytes('\n')
				}
			}
			// Continue viewing article
			continue

		case models.ArticleActionShowRelated:
			return action, nil

		default:
			// Back/quit
			return action, nil
		}
	}
}

func toggleBookmark(client *api.Client, article *api.Article) error {
	// Check if user is authenticated
	if !client.Config.IsAuthenticated() {
		return fmt.Errorf("you need to login to bookmark articles. Run 'bc-cli login' to authenticate")
	}

	if article.IsBookmarked {
		// Need to find bookmark ID to delete
		bookmarks, err := client.ListBookmarks()
		if err != nil {
			return err
		}
		for _, bm := range bookmarks {
			if bm.ArticleID == article.ID {
				return client.DeleteBookmark(bm.ID)
			}
		}
		return fmt.Errorf("bookmark not found")
	}
	_, err := client.CreateBookmark(article.ID)
	return err
}

func showBookmarksView(cfg *config.Config, client *api.Client) error {
	if !cfg.IsAuthenticated() {
		fmt.Println("\nPlease login to view bookmarks.")
		fmt.Println("Run 'bc-cli login' to authenticate.")
		return nil
	}

	bookmarks, err := client.ListBookmarks()
	if err != nil {
		return err
	}

	if len(bookmarks) == 0 {
		fmt.Println("\nYou don't have any bookmarks yet.")
		fmt.Println("Press 'b' while reading an article to bookmark it!")
		fmt.Print("\nPress Enter to continue...")
		_, _ = bufio.NewReader(os.Stdin).ReadBytes('\n')
		return nil
	}

	// Extract articles from bookmarks
	articles := make([]api.Article, len(bookmarks))
	for i, bm := range bookmarks {
		articles[i] = bm.Article
	}

	// Use article picker
	for {
		article, err := models.PickArticle(articles)
		if err != nil || article == nil {
			return err
		}

		// View article
		fullArticle, err := client.GetArticle(article.ID)
		if err != nil {
			fmt.Printf("\nError loading article: %v\n", err)
			fmt.Print("Press Enter to continue...")
			_, _ = bufio.NewReader(os.Stdin).ReadBytes('\n')
			continue
		}

		action, err := viewArticleWithActions(cfg, client, fullArticle)
		if err != nil {
			return err
		}

		// If user removed bookmark, refresh the bookmarks list
		if action == models.ArticleActionToggleBookmark && !fullArticle.IsBookmarked {
			// Refresh bookmarks list
			bookmarks, err = client.ListBookmarks()
			if err != nil {
				return err
			}

			if len(bookmarks) == 0 {
				fmt.Println("\nYou don't have any bookmarks left.")
				fmt.Print("Press Enter to continue...")
				_, _ = bufio.NewReader(os.Stdin).ReadBytes('\n')
				return nil
			}

			// Update articles list
			articles = make([]api.Article, len(bookmarks))
			for i, bm := range bookmarks {
				articles[i] = bm.Article
			}
		}
	}
}

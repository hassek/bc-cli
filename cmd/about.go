package cmd

import (
	"strings"

	"github.com/hassek/bc-cli/tui/components"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var aboutCmd = &cobra.Command{
	Use:   "about",
	Short: "Learn about Butler Coffee",
	Long:  `Discover our story, principles, and what makes Butler Coffee special.`,
	RunE:  runAbout,
}

func init() {
	rootCmd.AddCommand(aboutCmd)
}

func runAbout(cmd *cobra.Command, args []string) error {
	content := formatAboutContent()
	return components.ShowTextViewer("☕ About Butler Coffee", content)
}

func formatAboutContent() string {
	var b strings.Builder

	// Styles
	// titleStyle := lipgloss.NewStyle().
	// 	Bold(true).
	// 	Foreground(lipgloss.Color("86")). // Cyan
	// 	MarginBottom(1)

	sectionStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("214")). // Orange/yellow
		MarginTop(1).
		MarginBottom(1)

	paragraphStyle := lipgloss.NewStyle().
		Width(70).
		MarginBottom(1)

	principleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("86")). // Cyan
		Bold(true)

	emphasisStyle := lipgloss.NewStyle().
		Italic(true).
		Foreground(lipgloss.Color("247")) // Light gray

	// // Title
	// b.WriteString(titleStyle.Render("☕ About Butler Coffee"))
	// b.WriteString("\n")

	// Introduction
	intro := "\nAt Butler Coffee we keep things simple: we only share what we genuinely love, and we have fun doing it."
	b.WriteString(paragraphStyle.Render(intro))
	b.WriteString("\n")

	// First Principle
	b.WriteString(sectionStyle.Render("Our First Principle"))
	b.WriteString("\n")

	principle1 := principleStyle.Render("We only offer what we like") + " — nothing goes on our stock unless we'd " +
		"happily drink it ourselves. Every coffee, machine, and product we offer " +
		"has been tested, tasted and enjoyed by us first. If it doesn't meet " +
		"our own standards, it never makes it to yours."
	b.WriteString(paragraphStyle.Render(principle1))
	b.WriteString("\n")

	// Second Principle
	b.WriteString(sectionStyle.Render("Our Second Principle"))
	b.WriteString("\n")

	principle2 := "Just as important: " + principleStyle.Render("we should enjoy the ride") + ".\nThat means " +
		"we sometimes try ideas that doesn't have to make much sense on paper, simply " +
		"because they make us smile. Coffee is meant to be enjoyed, and we want " +
		"that spirit to show through everything we do."
	b.WriteString(paragraphStyle.Render(principle2))
	b.WriteString("\n")

	// What We Do
	b.WriteString(sectionStyle.Render("What We Do"))
	b.WriteString("\n")

	whatWeDo := "Today, we focus on bringing high-quality specialty coffee to both homes " +
		"and workplaces through curated subscriptions and office setups " +
		"offering beans from all around the world. Whether it's a single " +
		"bag or a full machine solution, our goal is to deliver products that make your " +
		"daily coffee something worth looking forward to."
	b.WriteString(paragraphStyle.Render(whatWeDo))
	b.WriteString("\n")

	// End Goal
	b.WriteString("\n")
	endGoal := emphasisStyle.Render("Our end goal is to make you make good coffee.")
	b.WriteString(paragraphStyle.Render(endGoal))
	b.WriteString("\n\n")

	// Footer decoration
	footerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")). // Faint
		Align(lipgloss.Center).
		Width(70)

	footer := "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n" +
		"                Made with ☕ and love\n" +
		"                butler.coffee"
	b.WriteString(footerStyle.Render(footer))

	return b.String()
}

package cmd

import (
	"github.com/hassek/bc-cli/templates"
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
	// In the future, this content would come from the API
	// For now, we demonstrate the template capabilities
	tmpl := getAboutTemplate()
	return templates.RenderInViewport("☕ About Butler Coffee", tmpl, nil)
}

// getAboutTemplate returns the About Us content as a template string.
// In production, this would be fetched from the backend CMS API.
// The backend would store this template and return it via an endpoint.
func getAboutTemplate() string {
	return `
{{paragraphAuto "\nAt Butler Coffee we keep things simple: we only share what we genuinely love, and we have fun doing it."}}

{{section "Our First Principle"}}

{{paragraphAuto (printf "%s — nothing goes on our stock unless we'd happily drink it ourselves. Every coffee, machine, and product we offer has been tested, tasted and enjoyed by us first. If it doesn't meet our own standards, it never makes it to yours." (highlight "We only offer what we like"))}}

{{section "Our Second Principle"}}

{{paragraphAuto (printf "Just as important: %s.\nThat means we sometimes try ideas that don't make much sense on paper, simply because they make us smile. Coffee is meant to be enjoyed, and we want that spirit to show through everything we do." (highlight "we should enjoy the ride"))}}

{{section "What We Do"}}

{{paragraphAuto "Today, we focus on bringing high-quality specialty coffee to both homes and workplaces through curated subscriptions and office setups offering beans from all around the world. Whether it's a single bag or a full machine solution, our goal is to deliver products that make your daily coffee something worth looking forward to."}}


{{paragraphAuto (emphasis "Our end goal is to make you make good coffee.")}}


{{faint "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n                Made with ☕ and love\n                butler.coffee"}}
`
}

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Version information - set via build flags if needed
var (
	Version   = "dev"
	GitCommit = "unknown"
	BuildDate = "unknown"
)

var rootCmd = &cobra.Command{
	Use:   "bc-cli",
	Short: "Butler Coffee CLI tool",
	Long: `Butler Coffee CLI tool - A command line interface for managing your coffee operations.

Complete documentation is available at https://github.com/hassek/bc-cli`,
	Run: func(cmd *cobra.Command, args []string) {
		// Check if version flag is set
		version, _ := cmd.Flags().GetBool("version")
		if version {
			fmt.Printf("bc-cli version %s\n", Version)
			fmt.Printf("  git commit: %s\n", GitCommit)
			fmt.Printf("  built: %s\n", BuildDate)
			return
		}
		_ = cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("version", "v", false, "Print version information")
}

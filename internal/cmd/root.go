// internal/cmd/root.go
package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "intervals",
	Short: "CLI for the intervals.icu API",
	Long:  "intervals is a command-line interface for the intervals.icu training analytics platform.\nDesigned for AI agents and human operators. 100% API coverage.",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().String("format", "json", "Output format: json, ndjson, text")
	rootCmd.PersistentFlags().String("fields", "", "Comma-separated fields to include in output")
	rootCmd.PersistentFlags().Bool("dry-run", false, "Show request without executing")
	rootCmd.PersistentFlags().Bool("yes", false, "Skip confirmation prompts")
	rootCmd.PersistentFlags().String("api-key", "", "API key (overrides config)")
	rootCmd.PersistentFlags().String("athlete-id", "", "Athlete ID (overrides config)")
	rootCmd.PersistentFlags().String("base-url", "", "API base URL (overrides config)")
}

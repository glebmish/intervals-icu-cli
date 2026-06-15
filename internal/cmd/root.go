package cmd

import (
	"fmt"
	"os"

	"github.com/glebmish/intervals-icu-cli/internal/api"
	"github.com/glebmish/intervals-icu-cli/internal/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:          "intervals",
	Short:        "CLI for the intervals.icu API",
	Long:         "intervals is a command-line interface for the intervals.icu training analytics platform.\nDesigned for AI agents and human operators. 100% API coverage.",
	SilenceUsage: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Offline commands skip config loading.
		switch cmd.Name() {
		case "schema", "skills", "config", "help", "init":
			return nil
		}
		// Walk up to find an offline-group ancestor (e.g. `config init`, `skills list`).
		for p := cmd.Parent(); p != nil; p = p.Parent() {
			switch p.Name() {
			case "schema", "skills", "config":
				return nil
			}
		}

		cfgPath := os.Getenv("INTERVALS_CONFIG")
		if cfgPath == "" {
			cfgPath = config.DefaultPath()
		}

		cfg, err := config.Load(cfgPath)
		if err != nil {
			return err
		}
		cfg.ApplyEnv()

		apiKey, _ := cmd.Flags().GetString("api-key")
		athleteID, _ := cmd.Flags().GetString("athlete-id")
		baseURL, _ := cmd.Flags().GetString("base-url")
		cfg.ApplyFlags(apiKey, athleteID, baseURL)

		if err := cfg.Validate(); err != nil {
			return err
		}

		client := api.NewClient(cfg.BaseURL, cfg.APIKey, cfg.AthleteID)
		cmd.SetContext(api.WithContext(cmd.Context(), client))
		return nil
	},
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
	rootCmd.PersistentFlags().String("json", "", "Raw JSON request body for write ops; see 'intervals schema <op>' for the shape")
	rootCmd.PersistentFlags().String("params", "", "Raw JSON object overlaying query/path params, e.g. '{\"id\":\"abc\"}'")
}

// confirmDelete prompts for confirmation on delete operations.
func confirmDelete(cmd *cobra.Command, resource, id string) error {
	yes, _ := cmd.Flags().GetBool("yes")
	if yes {
		return nil
	}

	fi, err := os.Stdin.Stat()
	if err != nil || (fi.Mode()&os.ModeCharDevice) == 0 {
		return fmt.Errorf("delete %s %s requires --yes flag in non-interactive mode", resource, id)
	}

	fmt.Fprintf(os.Stderr, "Delete %s %s? [y/N] ", resource, id)
	var response string
	fmt.Scanln(&response)
	if response != "y" && response != "Y" {
		return fmt.Errorf("cancelled")
	}
	return nil
}

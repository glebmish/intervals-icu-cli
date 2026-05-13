package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/glebmish/intervals-icu-cli/internal/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage CLI configuration",
}

var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize configuration file",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfgPath := configFilePath()

		// Preserve existing values if the file is already there.
		existing, _ := config.Load(cfgPath)
		if existing == nil {
			existing = &config.Config{}
		}
		existing.ApplyEnv()

		reader := bufio.NewReader(os.Stdin)

		fmt.Fprint(os.Stderr, "API key")
		if existing.APIKey != "" {
			fmt.Fprintf(os.Stderr, " [%s]", maskToken(existing.APIKey))
		}
		fmt.Fprint(os.Stderr, ": ")
		apiKey, _ := reader.ReadString('\n')
		apiKey = strings.TrimSpace(apiKey)
		if apiKey == "" {
			apiKey = existing.APIKey
		}
		if apiKey == "" {
			return fmt.Errorf("API key is required\n  Get one at https://intervals.icu/settings")
		}

		fmt.Fprint(os.Stderr, "Athlete ID")
		if existing.AthleteID != "" {
			fmt.Fprintf(os.Stderr, " [%s]", existing.AthleteID)
		}
		fmt.Fprint(os.Stderr, ": ")
		athleteID, _ := reader.ReadString('\n')
		athleteID = strings.TrimSpace(athleteID)
		if athleteID == "" {
			athleteID = existing.AthleteID
		}
		if athleteID == "" {
			return fmt.Errorf("athlete ID is required")
		}

		content := fmt.Sprintf("api_key: %q\nathlete_id: %q\nbase_url: \"https://intervals.icu\"\n", apiKey, athleteID)

		dir := filepath.Dir(cfgPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("creating config directory: %w", err)
		}

		if err := os.WriteFile(cfgPath, []byte(content), 0600); err != nil {
			return fmt.Errorf("writing config: %w", err)
		}

		fmt.Fprintf(os.Stderr, "Config written to %s\n", cfgPath)
		return nil
	},
}

var configPathCmd = &cobra.Command{
	Use:   "path",
	Short: "Print the effective config file path",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println(configFilePath())
		return nil
	},
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Print the resolved config (API key masked unless --unmasked)",
	RunE: func(cmd *cobra.Command, args []string) error {
		path := configFilePath()
		cfg, err := config.Load(path)
		if err != nil {
			return err
		}
		cfg.ApplyEnv()
		unmasked, _ := cmd.Flags().GetBool("unmasked")
		key := maskToken(cfg.APIKey)
		if unmasked {
			key = cfg.APIKey
		}
		fmt.Printf("path:       %s\n", path)
		fmt.Printf("api_key:    %s\n", key)
		fmt.Printf("athlete_id: %s\n", cfg.AthleteID)
		fmt.Printf("base_url:   %s\n", cfg.BaseURL)
		return nil
	},
}

func configFilePath() string {
	if v := os.Getenv("INTERVALS_CONFIG"); v != "" {
		return v
	}
	return config.DefaultPath()
}

func maskToken(t string) string {
	if t == "" {
		return "(unset)"
	}
	if len(t) <= 8 {
		return strings.Repeat("*", len(t))
	}
	return t[:4] + strings.Repeat("*", len(t)-8) + t[len(t)-4:]
}

func init() {
	configShowCmd.Flags().Bool("unmasked", false, "Print the API key in clear text (for headless bootstrap)")
	configCmd.AddCommand(configInitCmd, configPathCmd, configShowCmd)
	rootCmd.AddCommand(configCmd)
}

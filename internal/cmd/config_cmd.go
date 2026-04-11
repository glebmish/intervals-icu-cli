// internal/cmd/config_cmd.go
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
		cfgPath := os.Getenv("INTERVALS_CONFIG")
		if cfgPath == "" {
			cfgPath = config.DefaultPath()
		}

		reader := bufio.NewReader(os.Stdin)

		fmt.Print("API key: ")
		apiKey, _ := reader.ReadString('\n')
		apiKey = strings.TrimSpace(apiKey)

		fmt.Print("Athlete ID: ")
		athleteID, _ := reader.ReadString('\n')
		athleteID = strings.TrimSpace(athleteID)

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

func init() {
	configCmd.AddCommand(configInitCmd)
	rootCmd.AddCommand(configCmd)
}

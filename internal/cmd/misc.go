// internal/cmd/misc.go
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var miscCmd = &cobra.Command{
	Use:   "misc",
	Short: "Miscellaneous API operations",
}

func init() {
	miscCmd.AddCommand(
		newMiscPaceDistancesCmd(),
		newMiscDownloadWorkoutCmd(),
		newMiscUpdateAthletePlansCmd(),
		newMiscDisconnectAppCmd(),
	)
	rootCmd.AddCommand(miscCmd)
}

func newMiscPaceDistancesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pace-distances",
		Short: "Get pace distances",
		RunE: func(cmd *cobra.Command, args []string) error {
			return doGet(cmd, "/api/v1/pace_distances", nil)
		},
	}
	return cmd
}

func newMiscDownloadWorkoutCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "download-workout",
		Short: "Download a workout file",
		RunE: func(cmd *cobra.Command, args []string) error {
			ext, _ := cmd.Flags().GetString("ext")
			jsonBody, _ := cmd.Flags().GetString("json")
			if jsonBody == "" {
				return fmt.Errorf("--json is required with workout payload")
			}
			params := map[string]string{}
			if ext != "" {
				params["ext"] = ext
			}
			path := "/api/v1/download-workout" + ext
			return doMutate(cmd, "POST", path, params, jsonBody)
		},
	}
	cmd.Flags().String("ext", "", "File extension (e.g. .zwo)")
	cmd.Flags().String("json", "", "Workout JSON payload (required)")
	return cmd
}

func newMiscUpdateAthletePlansCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-athlete-plans",
		Short: "Update athlete plans",
		RunE: func(cmd *cobra.Command, args []string) error {
			jsonBody, _ := cmd.Flags().GetString("json")
			if jsonBody == "" {
				return fmt.Errorf("--json is required with athlete plans payload")
			}
			return doMutate(cmd, "PUT", "/api/v1/athlete-plans", nil, jsonBody)
		},
	}
	cmd.Flags().String("json", "", "Athlete plans JSON payload (required)")
	return cmd
}

func newMiscDisconnectAppCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "disconnect-app",
		Short: "Disconnect an integrated app",
		RunE: func(cmd *cobra.Command, args []string) error {
			return doDelete(cmd, "/api/v1/disconnect-app", nil, "app", "integration")
		},
	}
	return cmd
}

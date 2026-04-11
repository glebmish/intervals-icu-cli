// internal/cmd/sport_settings.go
package cmd

import (
	"fmt"

	"github.com/glebmish/intervals-icu-cli/internal/validate"
	"github.com/spf13/cobra"
)

var sportSettingsCmd = &cobra.Command{
	Use:   "sport-settings",
	Short: "Sport-specific settings and zones",
}

func init() {
	sportSettingsCmd.AddCommand(
		newSportSettingsListCmd(),
		newSportSettingsGetCmd(),
		newSportSettingsGetDeviceCmd(),
		newSportSettingsCreateCmd(),
		newSportSettingsUpdateCmd(),
		newSportSettingsUpdateMultiCmd(),
		newSportSettingsDeleteCmd(),
		newSportSettingsApplyCmd(),
		newSportSettingsPaceDistancesCmd(),
		newSportSettingsMatchingActivitiesCmd(),
	)
	rootCmd.AddCommand(sportSettingsCmd)
}

// 1. GET /api/v1/athlete/{athleteId}/sport-settings
func newSportSettingsListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List sport settings",
		RunE: func(cmd *cobra.Command, args []string) error {
			return doGet(cmd, "/api/v1/athlete/{id}/sport-settings", nil)
		},
	}
	return cmd
}

// 2. GET /api/v1/athlete/{athleteId}/sport-settings/{id}
func newSportSettingsGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get a sport setting",
		RunE: func(cmd *cobra.Command, args []string) error {
			settingID, _ := cmd.Flags().GetString("setting-id")
			if settingID == "" {
				return fmt.Errorf("--setting-id is required")
			}
			if err := validate.PathParam("setting-id", settingID); err != nil {
				return err
			}
			params := map[string]string{"id": settingID}
			return doGet(cmd, "/api/v1/athlete/{athleteId}/sport-settings/{id}", params)
		},
	}
	cmd.Flags().String("setting-id", "", "Sport setting ID (required)")
	return cmd
}

// 3. GET /api/v1/athlete/{id}/settings/{deviceClass}
func newSportSettingsGetDeviceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-device",
		Short: "Get device settings",
		RunE: func(cmd *cobra.Command, args []string) error {
			deviceClass, _ := cmd.Flags().GetString("device-class")
			if deviceClass == "" {
				return fmt.Errorf("--device-class is required")
			}
			if err := validate.PathParam("device-class", deviceClass); err != nil {
				return err
			}
			params := map[string]string{"deviceClass": deviceClass}
			return doGet(cmd, "/api/v1/athlete/{id}/settings/{deviceClass}", params)
		},
	}
	cmd.Flags().String("device-class", "", "Device class (required)")
	return cmd
}

// 4. POST /api/v1/athlete/{athleteId}/sport-settings
func newSportSettingsCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create sport settings",
		RunE: func(cmd *cobra.Command, args []string) error {
			jsonBody, _ := cmd.Flags().GetString("json")
			if jsonBody == "" {
				return fmt.Errorf("--json is required with sport settings payload")
			}
			return doMutate(cmd, "POST", "/api/v1/athlete/{id}/sport-settings", nil, jsonBody)
		},
	}
	cmd.Flags().String("json", "", "Sport settings JSON payload (required)")
	return cmd
}

// 5. PUT /api/v1/athlete/{athleteId}/sport-settings/{id}
func newSportSettingsUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update sport settings",
		RunE: func(cmd *cobra.Command, args []string) error {
			settingID, _ := cmd.Flags().GetString("setting-id")
			if settingID == "" {
				return fmt.Errorf("--setting-id is required")
			}
			jsonBody, _ := cmd.Flags().GetString("json")
			if jsonBody == "" {
				return fmt.Errorf("--json is required with sport settings payload")
			}
			params := map[string]string{"id": settingID}
			if v, _ := cmd.Flags().GetBool("recalc-hr-zones"); v {
				params["recalcHrZones"] = "true"
			}
			return doMutate(cmd, "PUT", "/api/v1/athlete/{athleteId}/sport-settings/{id}", params, jsonBody)
		},
	}
	cmd.Flags().String("setting-id", "", "Sport setting ID (required)")
	cmd.Flags().String("json", "", "Sport settings JSON payload (required)")
	cmd.Flags().Bool("recalc-hr-zones", false, "Recalculate HR zones")
	return cmd
}

// 6. PUT /api/v1/athlete/{athleteId}/sport-settings
func newSportSettingsUpdateMultiCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-multi",
		Short: "Update multiple sport settings",
		RunE: func(cmd *cobra.Command, args []string) error {
			jsonBody, _ := cmd.Flags().GetString("json")
			if jsonBody == "" {
				return fmt.Errorf("--json is required with sport settings payload")
			}
			params := map[string]string{}
			if v, _ := cmd.Flags().GetBool("recalc-hr-zones"); v {
				params["recalcHrZones"] = "true"
			}
			return doMutate(cmd, "PUT", "/api/v1/athlete/{id}/sport-settings", params, jsonBody)
		},
	}
	cmd.Flags().String("json", "", "Sport settings JSON payload (required)")
	cmd.Flags().Bool("recalc-hr-zones", false, "Recalculate HR zones")
	return cmd
}

// 7. DELETE /api/v1/athlete/{athleteId}/sport-settings/{id}
func newSportSettingsDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete sport settings",
		RunE: func(cmd *cobra.Command, args []string) error {
			settingID, _ := cmd.Flags().GetString("setting-id")
			if settingID == "" {
				return fmt.Errorf("--setting-id is required")
			}
			params := map[string]string{"id": settingID}
			return doDelete(cmd, "/api/v1/athlete/{athleteId}/sport-settings/{id}", params, "sport-settings", settingID)
		},
	}
	cmd.Flags().String("setting-id", "", "Sport setting ID (required)")
	return cmd
}

// 8. PUT /api/v1/athlete/{athleteId}/sport-settings/{id}/apply
func newSportSettingsApplyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "apply",
		Short: "Apply sport settings",
		RunE: func(cmd *cobra.Command, args []string) error {
			settingID, _ := cmd.Flags().GetString("setting-id")
			if settingID == "" {
				return fmt.Errorf("--setting-id is required")
			}
			params := map[string]string{"id": settingID}
			return doMutate(cmd, "PUT", "/api/v1/athlete/{athleteId}/sport-settings/{id}/apply", params, "")
		},
	}
	cmd.Flags().String("setting-id", "", "Sport setting ID (required)")
	return cmd
}

// 9. GET /api/v1/athlete/{athleteId}/sport-settings/{id}/pace_distances
func newSportSettingsPaceDistancesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pace-distances",
		Short: "Get pace distances for sport settings",
		RunE: func(cmd *cobra.Command, args []string) error {
			settingID, _ := cmd.Flags().GetString("setting-id")
			if settingID == "" {
				return fmt.Errorf("--setting-id is required")
			}
			params := map[string]string{"id": settingID}
			return doGet(cmd, "/api/v1/athlete/{athleteId}/sport-settings/{id}/pace_distances", params)
		},
	}
	cmd.Flags().String("setting-id", "", "Sport setting ID (required)")
	return cmd
}

// 10. GET /api/v1/athlete/{athleteId}/sport-settings/{id}/matching-activities
func newSportSettingsMatchingActivitiesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "matching-activities",
		Short: "Get activities matching sport settings",
		RunE: func(cmd *cobra.Command, args []string) error {
			settingID, _ := cmd.Flags().GetString("setting-id")
			if settingID == "" {
				return fmt.Errorf("--setting-id is required")
			}
			params := map[string]string{"id": settingID}
			return doGet(cmd, "/api/v1/athlete/{athleteId}/sport-settings/{id}/matching-activities", params)
		},
	}
	cmd.Flags().String("setting-id", "", "Sport setting ID (required)")
	return cmd
}

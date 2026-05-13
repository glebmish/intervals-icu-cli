// internal/cmd/gear.go
package cmd

import (
	"fmt"

	"github.com/glebmish/intervals-icu-cli/internal/validate"
	"github.com/spf13/cobra"
)

var gearCmd = &cobra.Command{
	Use:   "gear",
	Short: "Equipment tracking and reminders",
}

func init() {
	gearCmd.AddCommand(
		newGearListCmd(),
		newGearCreateCmd(),
		newGearUpdateCmd(),
		newGearDeleteCmd(),
		newGearReplaceCmd(),
		newGearCalcCmd(),
		newGearCreateReminderCmd(),
		newGearUpdateReminderCmd(),
		newGearDeleteReminderCmd(),
	)
	rootCmd.AddCommand(gearCmd)
}

// 1. GET /api/v1/athlete/{id}/gear{ext}
func newGearListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List athlete gear",
		RunE: func(cmd *cobra.Command, args []string) error {
			ext, _ := cmd.Flags().GetString("ext")
			params := map[string]string{"ext": ext}
			return doGet(cmd, "/api/v1/athlete/{id}/gear{ext}", params)
		},
	}
	cmd.Flags().String("ext", "", "File extension (e.g. .csv)")
	return cmd
}

// 2. POST /api/v1/athlete/{id}/gear
func newGearCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create gear item",
		RunE: func(cmd *cobra.Command, args []string) error {
			jsonBody, err := requireJSON(cmd)
			if err != nil {
				return err
			}
			return doMutate(cmd, "POST", "/api/v1/athlete/{id}/gear", nil, jsonBody)
		},
	}
	return cmd
}

// 3. PUT /api/v1/athlete/{id}/gear/{gearId}
func newGearUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update gear item",
		RunE: func(cmd *cobra.Command, args []string) error {
			gearID, err := requireString(cmd, "gear-id")
			if err != nil {
				return err
			}
			if err := validate.PathParam("gear-id", gearID); err != nil {
				return err
			}
			jsonBody, err := requireJSON(cmd)
			if err != nil {
				return err
			}
			params := map[string]string{"gearId": gearID}
			return doMutate(cmd, "PUT", "/api/v1/athlete/{id}/gear/{gearId}", params, jsonBody)
		},
	}
	cmd.Flags().String("gear-id", "", "Gear ID (required)")
	return cmd
}

// 4. DELETE /api/v1/athlete/{id}/gear/{gearId}
func newGearDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete gear item",
		RunE: func(cmd *cobra.Command, args []string) error {
			gearID, err := requireString(cmd, "gear-id")
			if err != nil {
				return err
			}
			params := map[string]string{"gearId": gearID}
			return doDelete(cmd, "/api/v1/athlete/{id}/gear/{gearId}", params, "gear", gearID)
		},
	}
	cmd.Flags().String("gear-id", "", "Gear ID (required)")
	return cmd
}

// 5. POST /api/v1/athlete/{id}/gear/{gearId}/replace
func newGearReplaceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "replace",
		Short: "Replace gear item",
		RunE: func(cmd *cobra.Command, args []string) error {
			gearID, err := requireString(cmd, "gear-id")
			if err != nil {
				return err
			}
			jsonBody, err := requireJSON(cmd)
			if err != nil {
				return err
			}
			params := map[string]string{"gearId": gearID}
			return doMutate(cmd, "POST", "/api/v1/athlete/{id}/gear/{gearId}/replace", params, jsonBody)
		},
	}
	cmd.Flags().String("gear-id", "", "Gear ID (required)")
	return cmd
}

// 6. GET /api/v1/athlete/{id}/gear/{gearId}/calc
func newGearCalcCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "calc",
		Short: "Calculate gear statistics",
		RunE: func(cmd *cobra.Command, args []string) error {
			gearID, err := requireString(cmd, "gear-id")
			if err != nil {
				return err
			}
			params := map[string]string{"gearId": gearID}
			return doGet(cmd, "/api/v1/athlete/{id}/gear/{gearId}/calc", params)
		},
	}
	cmd.Flags().String("gear-id", "", "Gear ID (required)")
	return cmd
}

// 7. POST /api/v1/athlete/{id}/gear/{gearId}/reminder
func newGearCreateReminderCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-reminder",
		Short: "Create gear reminder",
		RunE: func(cmd *cobra.Command, args []string) error {
			gearID, err := requireString(cmd, "gear-id")
			if err != nil {
				return err
			}
			jsonBody, err := requireJSON(cmd)
			if err != nil {
				return err
			}
			params := map[string]string{"gearId": gearID}
			return doMutate(cmd, "POST", "/api/v1/athlete/{id}/gear/{gearId}/reminder", params, jsonBody)
		},
	}
	cmd.Flags().String("gear-id", "", "Gear ID (required)")
	return cmd
}

// 8. PUT /api/v1/athlete/{id}/gear/{gearId}/reminder/{reminderId}
func newGearUpdateReminderCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-reminder",
		Short: "Update gear reminder",
		RunE: func(cmd *cobra.Command, args []string) error {
			gearID, err := requireString(cmd, "gear-id")
			if err != nil {
				return err
			}
			reminderID, err := requireString(cmd, "reminder-id")
			if err != nil {
				return err
			}
			jsonBody, err := requireJSON(cmd)
			if err != nil {
				return err
			}
			params := map[string]string{"gearId": gearID, "reminderId": reminderID}
			if v, _ := cmd.Flags().GetBool("reset"); v {
				params["reset"] = "true"
			}
			if v, _ := cmd.Flags().GetInt("snooze-days"); v > 0 {
				params["snoozeDays"] = fmt.Sprintf("%d", v)
			}
			return doMutate(cmd, "PUT", "/api/v1/athlete/{id}/gear/{gearId}/reminder/{reminderId}", params, jsonBody)
		},
	}
	cmd.Flags().String("gear-id", "", "Gear ID (required)")
	cmd.Flags().String("reminder-id", "", "Reminder ID (required)")
	cmd.Flags().Bool("reset", false, "Reset reminder counter")
	cmd.Flags().Int("snooze-days", 0, "Snooze reminder by N days")
	return cmd
}

// 9. DELETE /api/v1/athlete/{id}/gear/{gearId}/reminder/{reminderId}
func newGearDeleteReminderCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-reminder",
		Short: "Delete gear reminder",
		RunE: func(cmd *cobra.Command, args []string) error {
			gearID, err := requireString(cmd, "gear-id")
			if err != nil {
				return err
			}
			reminderID, err := requireString(cmd, "reminder-id")
			if err != nil {
				return err
			}
			params := map[string]string{"gearId": gearID, "reminderId": reminderID}
			return doDelete(cmd, "/api/v1/athlete/{id}/gear/{gearId}/reminder/{reminderId}", params, "gear reminder", reminderID)
		},
	}
	cmd.Flags().String("gear-id", "", "Gear ID (required)")
	cmd.Flags().String("reminder-id", "", "Reminder ID (required)")
	return cmd
}

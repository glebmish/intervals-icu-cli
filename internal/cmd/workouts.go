// internal/cmd/workouts.go
package cmd

import (
	"fmt"

	"github.com/glebmish/intervals-icu-cli/internal/validate"
	"github.com/spf13/cobra"
)

var workoutsCmd = &cobra.Command{
	Use:   "workouts",
	Short: "Workout library — CRUD, bulk, duplicate",
}

func init() {
	workoutsCmd.AddCommand(
		newWorkoutsListCmd(),
		newWorkoutsGetCmd(),
		newWorkoutsCreateCmd(),
		newWorkoutsUpdateCmd(),
		newWorkoutsDeleteCmd(),
		newWorkoutsCreateBulkCmd(),
		newWorkoutsDuplicateCmd(),
		newWorkoutsDownloadCmd(),
		newWorkoutsTagsCmd(),
	)
	rootCmd.AddCommand(workoutsCmd)
}

func newWorkoutsListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List workouts",
		RunE: func(cmd *cobra.Command, args []string) error {
			return doGet(cmd, "/api/v1/athlete/{id}/workouts", nil)
		},
	}
	return cmd
}

func newWorkoutsGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get a single workout by ID",
		RunE: func(cmd *cobra.Command, args []string) error {
			workoutID, _ := cmd.Flags().GetString("workout-id")
			if workoutID == "" {
				return fmt.Errorf("--workout-id is required")
			}
			if err := validate.PathParam("workout-id", workoutID); err != nil {
				return err
			}
			params := map[string]string{"workoutId": workoutID}
			return doGet(cmd, "/api/v1/athlete/{id}/workouts/{workoutId}", params)
		},
	}
	cmd.Flags().String("workout-id", "", "Workout ID (required)")
	return cmd
}

func newWorkoutsCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new workout",
		RunE: func(cmd *cobra.Command, args []string) error {
			jsonBody, _ := cmd.Flags().GetString("json")
			if jsonBody == "" {
				return fmt.Errorf("--json is required with Workout payload")
			}
			return doMutate(cmd, "POST", "/api/v1/athlete/{id}/workouts", nil, jsonBody)
		},
	}
	cmd.Flags().String("json", "", "Workout JSON payload (required)")
	return cmd
}

func newWorkoutsUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update a workout",
		RunE: func(cmd *cobra.Command, args []string) error {
			workoutID, _ := cmd.Flags().GetString("workout-id")
			if workoutID == "" {
				return fmt.Errorf("--workout-id is required")
			}
			if err := validate.PathParam("workout-id", workoutID); err != nil {
				return err
			}
			jsonBody, _ := cmd.Flags().GetString("json")
			if jsonBody == "" {
				return fmt.Errorf("--json is required with Workout payload")
			}
			params := map[string]string{"workoutId": workoutID}
			return doMutate(cmd, "PUT", "/api/v1/athlete/{id}/workouts/{workoutId}", params, jsonBody)
		},
	}
	cmd.Flags().String("workout-id", "", "Workout ID (required)")
	cmd.Flags().String("json", "", "Workout JSON payload (required)")
	return cmd
}

func newWorkoutsDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a workout",
		RunE: func(cmd *cobra.Command, args []string) error {
			workoutID, _ := cmd.Flags().GetString("workout-id")
			if workoutID == "" {
				return fmt.Errorf("--workout-id is required")
			}
			if err := validate.PathParam("workout-id", workoutID); err != nil {
				return err
			}
			params := map[string]string{"workoutId": workoutID}
			if v, _ := cmd.Flags().GetBool("others"); v {
				params["others"] = "true"
			}
			return doDelete(cmd, "/api/v1/athlete/{id}/workouts/{workoutId}", params, "workout", workoutID)
		},
	}
	cmd.Flags().String("workout-id", "", "Workout ID (required)")
	cmd.Flags().Bool("others", false, "Delete other related workouts")
	return cmd
}

func newWorkoutsCreateBulkCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-bulk",
		Short: "Bulk create workouts",
		RunE: func(cmd *cobra.Command, args []string) error {
			jsonBody, _ := cmd.Flags().GetString("json")
			if jsonBody == "" {
				return fmt.Errorf("--json is required with array of Workout payloads")
			}
			return doMutate(cmd, "POST", "/api/v1/athlete/{id}/workouts/bulk", nil, jsonBody)
		},
	}
	cmd.Flags().String("json", "", "JSON array of Workout payloads (required)")
	return cmd
}

func newWorkoutsDuplicateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "duplicate",
		Short: "Duplicate workouts",
		RunE: func(cmd *cobra.Command, args []string) error {
			jsonBody, _ := cmd.Flags().GetString("json")
			if jsonBody == "" {
				return fmt.Errorf("--json is required with duplicate payload")
			}
			return doMutate(cmd, "POST", "/api/v1/athlete/{id}/duplicate-workouts", nil, jsonBody)
		},
	}
	cmd.Flags().String("json", "", "Duplicate JSON payload (required)")
	return cmd
}

func newWorkoutsDownloadCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "download",
		Short: "Download workouts as zip",
		RunE: func(cmd *cobra.Command, args []string) error {
			params := map[string]string{}
			if v, _ := cmd.Flags().GetString("oldest"); v != "" {
				params["oldest"] = v
			}
			if v, _ := cmd.Flags().GetString("newest"); v != "" {
				params["newest"] = v
			}
			output, _ := cmd.Flags().GetString("output")
			ext, _ := cmd.Flags().GetString("ext")
			path := "/api/v1/athlete/{id}/workouts" + ext
			return doDownload(cmd, path, params, output)
		},
	}
	cmd.Flags().String("ext", ".zip", "File extension (default: .zip)")
	cmd.Flags().String("oldest", "", "Start date (YYYY-MM-DD)")
	cmd.Flags().String("newest", "", "End date (YYYY-MM-DD)")
	cmd.Flags().String("output", "", "Output file path (default: stdout)")
	return cmd
}

func newWorkoutsTagsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tags",
		Short: "List workout tags",
		RunE: func(cmd *cobra.Command, args []string) error {
			return doGet(cmd, "/api/v1/athlete/{id}/workout-tags", nil)
		},
	}
	return cmd
}

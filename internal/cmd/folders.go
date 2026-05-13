// internal/cmd/folders.go
package cmd

import (
	"github.com/glebmish/intervals-icu-cli/internal/validate"
	"github.com/spf13/cobra"
)

var foldersCmd = &cobra.Command{
	Use:   "folders",
	Short: "Workout folder management",
}

func init() {
	foldersCmd.AddCommand(
		newFoldersListCmd(),
		newFoldersCreateCmd(),
		newFoldersUpdateCmd(),
		newFoldersDeleteCmd(),
		newFoldersUpdateWorkoutsCmd(),
		newFoldersSharedWithCmd(),
		newFoldersUpdateSharedWithCmd(),
		newFoldersImportWorkoutCmd(),
	)
	rootCmd.AddCommand(foldersCmd)
}

// 1. GET /api/v1/athlete/{id}/folders
func newFoldersListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List workout folders",
		RunE: func(cmd *cobra.Command, args []string) error {
			return doGet(cmd, "/api/v1/athlete/{id}/folders", nil)
		},
	}
	return cmd
}

// 2. POST /api/v1/athlete/{id}/folders
func newFoldersCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a workout folder",
		RunE: func(cmd *cobra.Command, args []string) error {
			jsonBody, err := requireJSON(cmd)
			if err != nil {
				return err
			}
			return doMutate(cmd, "POST", "/api/v1/athlete/{id}/folders", nil, jsonBody)
		},
	}
	return cmd
}

// 3. PUT /api/v1/athlete/{id}/folders/{folderId}
func newFoldersUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update a workout folder",
		RunE: func(cmd *cobra.Command, args []string) error {
			folderID, err := requireString(cmd, "folder-id")
			if err != nil {
				return err
			}
			if err := validate.PathParam("folder-id", folderID); err != nil {
				return err
			}
			jsonBody, err := requireJSON(cmd)
			if err != nil {
				return err
			}
			params := map[string]string{"folderId": folderID}
			return doMutate(cmd, "PUT", "/api/v1/athlete/{id}/folders/{folderId}", params, jsonBody)
		},
	}
	cmd.Flags().String("folder-id", "", "Folder ID (required)")
	return cmd
}

// 4. DELETE /api/v1/athlete/{id}/folders/{folderId}
func newFoldersDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a workout folder",
		RunE: func(cmd *cobra.Command, args []string) error {
			folderID, err := requireString(cmd, "folder-id")
			if err != nil {
				return err
			}
			params := map[string]string{"folderId": folderID}
			return doDelete(cmd, "/api/v1/athlete/{id}/folders/{folderId}", params, "folder", folderID)
		},
	}
	cmd.Flags().String("folder-id", "", "Folder ID (required)")
	return cmd
}

// 5. PUT /api/v1/athlete/{id}/folders/{folderId}/workouts
func newFoldersUpdateWorkoutsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-workouts",
		Short: "Update workouts in a folder",
		RunE: func(cmd *cobra.Command, args []string) error {
			folderID, err := requireString(cmd, "folder-id")
			if err != nil {
				return err
			}
			jsonBody, err := requireJSON(cmd)
			if err != nil {
				return err
			}
			params := map[string]string{"folderId": folderID}
			if v, _ := cmd.Flags().GetString("oldest"); v != "" {
				params["oldest"] = v
			}
			if v, _ := cmd.Flags().GetString("newest"); v != "" {
				params["newest"] = v
			}
			return doMutate(cmd, "PUT", "/api/v1/athlete/{id}/folders/{folderId}/workouts", params, jsonBody)
		},
	}
	cmd.Flags().String("folder-id", "", "Folder ID (required)")
	cmd.Flags().String("oldest", "", "Start date filter")
	cmd.Flags().String("newest", "", "End date filter")
	return cmd
}

// 6. GET /api/v1/athlete/{id}/folders/{folderId}/shared-with
func newFoldersSharedWithCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "shared-with",
		Short: "Get folder sharing settings",
		RunE: func(cmd *cobra.Command, args []string) error {
			folderID, err := requireString(cmd, "folder-id")
			if err != nil {
				return err
			}
			params := map[string]string{"folderId": folderID}
			return doGet(cmd, "/api/v1/athlete/{id}/folders/{folderId}/shared-with", params)
		},
	}
	cmd.Flags().String("folder-id", "", "Folder ID (required)")
	return cmd
}

// 7. PUT /api/v1/athlete/{id}/folders/{folderId}/shared-with
func newFoldersUpdateSharedWithCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-shared-with",
		Short: "Update folder sharing settings",
		RunE: func(cmd *cobra.Command, args []string) error {
			folderID, err := requireString(cmd, "folder-id")
			if err != nil {
				return err
			}
			jsonBody, err := requireJSON(cmd)
			if err != nil {
				return err
			}
			params := map[string]string{"folderId": folderID}
			return doMutate(cmd, "PUT", "/api/v1/athlete/{id}/folders/{folderId}/shared-with", params, jsonBody)
		},
	}
	cmd.Flags().String("folder-id", "", "Folder ID (required)")
	return cmd
}

// 8. POST /api/v1/athlete/{id}/folders/{folderId}/import-workout
func newFoldersImportWorkoutCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "import-workout",
		Short: "Import a workout into a folder",
		RunE: func(cmd *cobra.Command, args []string) error {
			folderID, err := requireString(cmd, "folder-id")
			if err != nil {
				return err
			}
			jsonBody, err := requireJSON(cmd)
			if err != nil {
				return err
			}
			params := map[string]string{"folderId": folderID}
			if v, _ := cmd.Flags().GetString("type"); v != "" {
				params["type"] = v
			}
			return doMutate(cmd, "POST", "/api/v1/athlete/{id}/folders/{folderId}/import-workout", params, jsonBody)
		},
	}
	cmd.Flags().String("folder-id", "", "Folder ID (required)")
	cmd.Flags().String("type", "", "Workout type")
	return cmd
}

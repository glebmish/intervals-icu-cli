// internal/cmd/events.go
package cmd

import (
	"fmt"

	"github.com/glebmish/intervals-icu-cli/internal/validate"
	"github.com/spf13/cobra"
)

var eventsCmd = &cobra.Command{
	Use:   "events",
	Short: "Training calendar events — CRUD, bulk, plans",
}

func init() {
	eventsCmd.AddCommand(
		newEventsListCmd(),
		newEventsGetCmd(),
		newEventsCreateCmd(),
		newEventsUpdateCmd(),
		newEventsDeleteCmd(),
		newEventsUpdateBulkCmd(),
		newEventsDeleteRangeCmd(),
		newEventsDeleteBulkCmd(),
		newEventsCreateBulkCmd(),
		newEventsMarkDoneCmd(),
		newEventsApplyPlanCmd(),
		newEventsDownloadWorkoutCmd(),
		newEventsTagsCmd(),
	)
	rootCmd.AddCommand(eventsCmd)
}

func newEventsListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List events in date range",
		RunE: func(cmd *cobra.Command, args []string) error {
			params := map[string]string{}
			if v, _ := cmd.Flags().GetString("oldest"); v != "" {
				if err := validate.DateParam("oldest", v); err != nil {
					return err
				}
				params["oldest"] = v
			}
			if v, _ := cmd.Flags().GetString("newest"); v != "" {
				if err := validate.DateParam("newest", v); err != nil {
					return err
				}
				params["newest"] = v
			}
			if v, _ := cmd.Flags().GetString("category"); v != "" {
				params["category"] = v
			}
			if v, _ := cmd.Flags().GetInt("limit"); v > 0 {
				params["limit"] = fmt.Sprintf("%d", v)
			}
			ext, _ := cmd.Flags().GetString("format-ext")
			path := "/api/v1/athlete/{id}/events" + ext
			return doGet(cmd, path, params)
		},
	}
	cmd.Flags().String("format-ext", "", "Format extension for path (e.g. .ics)")
	cmd.Flags().String("oldest", "", "Start date (YYYY-MM-DD)")
	cmd.Flags().String("newest", "", "End date (YYYY-MM-DD)")
	cmd.Flags().String("category", "", "Category filter")
	cmd.Flags().Int("limit", 0, "Max number of events")
	return cmd
}

func newEventsGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get a single event by ID",
		RunE: func(cmd *cobra.Command, args []string) error {
			eventID, err := requireString(cmd, "event-id")
			if err != nil {
				return err
			}
			if err := validate.PathParam("event-id", eventID); err != nil {
				return err
			}
			params := map[string]string{"eventId": eventID}
			return doGet(cmd, "/api/v1/athlete/{id}/events/{eventId}", params)
		},
	}
	cmd.Flags().String("event-id", "", "Event ID (required)")
	return cmd
}

func newEventsCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new event",
		RunE: func(cmd *cobra.Command, args []string) error {
			jsonBody, err := requireJSON(cmd)
			if err != nil {
				return err
			}
			params := map[string]string{}
			if v, _ := cmd.Flags().GetString("upsert-on-uid"); v != "" {
				params["upsertOnUID"] = v
			}
			return doMutate(cmd, "POST", "/api/v1/athlete/{id}/events", params, jsonBody)
		},
	}
	cmd.Flags().String("upsert-on-uid", "", "Upsert on UID query param")
	return cmd
}

func newEventsUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update an event",
		RunE: func(cmd *cobra.Command, args []string) error {
			eventID, err := requireString(cmd, "event-id")
			if err != nil {
				return err
			}
			if err := validate.PathParam("event-id", eventID); err != nil {
				return err
			}
			jsonBody, err := requireJSON(cmd)
			if err != nil {
				return err
			}
			params := map[string]string{"eventId": eventID}
			return doMutate(cmd, "PUT", "/api/v1/athlete/{id}/events/{eventId}", params, jsonBody)
		},
	}
	cmd.Flags().String("event-id", "", "Event ID (required)")
	return cmd
}

func newEventsDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete an event",
		RunE: func(cmd *cobra.Command, args []string) error {
			eventID, err := requireString(cmd, "event-id")
			if err != nil {
				return err
			}
			if err := validate.PathParam("event-id", eventID); err != nil {
				return err
			}
			params := map[string]string{"eventId": eventID}
			if v, _ := cmd.Flags().GetBool("others"); v {
				params["others"] = "true"
			}
			if v, _ := cmd.Flags().GetString("not-before"); v != "" {
				if err := validate.DateParam("not-before", v); err != nil {
					return err
				}
				params["notBefore"] = v
			}
			return doDelete(cmd, "/api/v1/athlete/{id}/events/{eventId}", params, "event", eventID)
		},
	}
	cmd.Flags().String("event-id", "", "Event ID (required)")
	cmd.Flags().Bool("others", false, "Delete other occurrences")
	cmd.Flags().String("not-before", "", "Do not delete events before this date")
	return cmd
}

func newEventsUpdateBulkCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-bulk",
		Short: "Bulk update events",
		RunE: func(cmd *cobra.Command, args []string) error {
			jsonBody, err := requireJSON(cmd)
			if err != nil {
				return err
			}
			params := map[string]string{}
			if v, _ := cmd.Flags().GetString("oldest"); v != "" {
				if err := validate.DateParam("oldest", v); err != nil {
					return err
				}
				params["oldest"] = v
			}
			if v, _ := cmd.Flags().GetString("newest"); v != "" {
				if err := validate.DateParam("newest", v); err != nil {
					return err
				}
				params["newest"] = v
			}
			return doMutate(cmd, "PUT", "/api/v1/athlete/{id}/events", params, jsonBody)
		},
	}
	cmd.Flags().String("oldest", "", "Start date (YYYY-MM-DD)")
	cmd.Flags().String("newest", "", "End date (YYYY-MM-DD)")
	return cmd
}

func newEventsDeleteRangeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-range",
		Short: "Delete events in date range",
		RunE: func(cmd *cobra.Command, args []string) error {
			oldest, err := requireString(cmd, "oldest")
			if err != nil {
				return err
			}
			newest, err := requireString(cmd, "newest")
			if err != nil {
				return err
			}
			if err := validate.DateParam("oldest", oldest); err != nil {
				return err
			}
			if err := validate.DateParam("newest", newest); err != nil {
				return err
			}
			params := map[string]string{"oldest": oldest, "newest": newest}
			if v, _ := cmd.Flags().GetString("category"); v != "" {
				params["category"] = v
			}
			return doDelete(cmd, "/api/v1/athlete/{id}/events", params, "events in range", oldest+" to "+newest)
		},
	}
	cmd.Flags().String("oldest", "", "Start date (required)")
	cmd.Flags().String("newest", "", "End date (required)")
	cmd.Flags().String("category", "", "Category filter")
	return cmd
}

func newEventsDeleteBulkCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-bulk",
		Short: "Bulk delete events by IDs",
		RunE: func(cmd *cobra.Command, args []string) error {
			jsonBody, err := requireJSON(cmd)
			if err != nil {
				return err
			}
			return doMutate(cmd, "PUT", "/api/v1/athlete/{id}/events/bulk-delete", nil, jsonBody)
		},
	}
	return cmd
}

func newEventsCreateBulkCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-bulk",
		Short: "Bulk create events",
		RunE: func(cmd *cobra.Command, args []string) error {
			jsonBody, err := requireJSON(cmd)
			if err != nil {
				return err
			}
			params := map[string]string{}
			if v, _ := cmd.Flags().GetBool("upsert"); v {
				params["upsert"] = "true"
			}
			if v, _ := cmd.Flags().GetBool("upsert-on-uid"); v {
				params["upsertOnUID"] = "true"
			}
			return doMutate(cmd, "POST", "/api/v1/athlete/{id}/events/bulk", params, jsonBody)
		},
	}
	cmd.Flags().Bool("upsert", false, "Upsert events")
	cmd.Flags().Bool("upsert-on-uid", false, "Upsert on UID")
	return cmd
}

func newEventsMarkDoneCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mark-done",
		Short: "Mark an event as done",
		RunE: func(cmd *cobra.Command, args []string) error {
			eventID, err := requireString(cmd, "event-id")
			if err != nil {
				return err
			}
			if err := validate.PathParam("event-id", eventID); err != nil {
				return err
			}
			params := map[string]string{"eventId": eventID}
			return doMutate(cmd, "POST", "/api/v1/athlete/{id}/events/{eventId}/mark-done", params, "{}")
		},
	}
	cmd.Flags().String("event-id", "", "Event ID (required)")
	return cmd
}

func newEventsApplyPlanCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "apply-plan",
		Short: "Apply a training plan",
		RunE: func(cmd *cobra.Command, args []string) error {
			jsonBody, err := requireJSON(cmd)
			if err != nil {
				return err
			}
			return doMutate(cmd, "POST", "/api/v1/athlete/{id}/events/apply-plan", nil, jsonBody)
		},
	}
	return cmd
}

func newEventsDownloadWorkoutCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "download-workout",
		Short: "Download workout file for an event",
		RunE: func(cmd *cobra.Command, args []string) error {
			eventID, err := requireString(cmd, "event-id")
			if err != nil {
				return err
			}
			if err := validate.PathParam("event-id", eventID); err != nil {
				return err
			}
			ext, _ := cmd.Flags().GetString("ext")
			if ext != "" {
				if err := validate.PathParam("ext", ext); err != nil {
					return err
				}
			}
			output, _ := cmd.Flags().GetString("output")
			params := map[string]string{"eventId": eventID}
			path := "/api/v1/athlete/{id}/events/{eventId}/download" + ext
			return doDownload(cmd, path, params, output)
		},
	}
	cmd.Flags().String("event-id", "", "Event ID (required)")
	cmd.Flags().String("ext", "", "File extension (e.g. .zwo)")
	cmd.Flags().String("output", "", "Output file path (default: stdout)")
	return cmd
}

func newEventsTagsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tags",
		Short: "List event tags",
		RunE: func(cmd *cobra.Command, args []string) error {
			return doGet(cmd, "/api/v1/athlete/{id}/event-tags", nil)
		},
	}
	return cmd
}

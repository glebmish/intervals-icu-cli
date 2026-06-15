package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/glebmish/intervals-icu-cli/internal/api"
	"github.com/glebmish/intervals-icu-cli/internal/format"
	"github.com/glebmish/intervals-icu-cli/internal/validate"
	"github.com/spf13/cobra"
)

var activitiesCmd = &cobra.Command{
	Use:   "activities",
	Short: "Athlete activities — list, search, create, upload",
}

func init() {
	activitiesCmd.AddCommand(
		newActivitiesListCmd(),
		newActivitiesSearchCmd(),
		newActivitiesSearchFullCmd(),
		newActivitiesIntervalSearchCmd(),
		newActivitiesCreateManualCmd(),
		newActivitiesCreateManualBulkCmd(),
		newActivitiesUploadCmd(),
		newActivitiesDownloadCSVCmd(),
		newActivitiesListAroundCmd(),
		newActivitiesGetMultipleCmd(),
	)
	rootCmd.AddCommand(activitiesCmd)
}

func newActivitiesListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List athlete activities",
		RunE: func(cmd *cobra.Command, args []string) error {
			oldest, err := requireString(cmd, "oldest")
			if err != nil {
				return err
			}
			if err := validate.DateParam("oldest", oldest); err != nil {
				return err
			}
			params := map[string]string{"oldest": oldest}
			if v, _ := cmd.Flags().GetString("newest"); v != "" {
				if err := validate.DateParam("newest", v); err != nil {
					return err
				}
				params["newest"] = v
			}
			if v, _ := cmd.Flags().GetInt("limit"); v > 0 {
				params["limit"] = fmt.Sprintf("%d", v)
			}
			return doGet(cmd, "/api/v1/athlete/{id}/activities", params)
		},
	}
	cmd.Flags().String("oldest", "", "Start date (YYYY-MM-DD, required)")
	cmd.Flags().String("newest", "", "End date (YYYY-MM-DD)")
	cmd.Flags().Int("limit", 0, "Max number of activities")
	return cmd
}

func newActivitiesSearchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search",
		Short: "Search activities by text",
		RunE: func(cmd *cobra.Command, args []string) error {
			q, err := requireString(cmd, "query")
			if err != nil {
				return err
			}
			params := map[string]string{"q": q}
			if v, _ := cmd.Flags().GetInt("limit"); v > 0 {
				params["limit"] = fmt.Sprintf("%d", v)
			}
			return doGet(cmd, "/api/v1/athlete/{id}/activities/search", params)
		},
	}
	cmd.Flags().String("query", "", "Search query (required)")
	cmd.Flags().Int("limit", 0, "Max results")
	return cmd
}

func newActivitiesSearchFullCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search-full",
		Short: "Search activities with full details",
		RunE: func(cmd *cobra.Command, args []string) error {
			q, err := requireString(cmd, "query")
			if err != nil {
				return err
			}
			params := map[string]string{"q": q}
			if v, _ := cmd.Flags().GetInt("limit"); v > 0 {
				params["limit"] = fmt.Sprintf("%d", v)
			}
			return doGet(cmd, "/api/v1/athlete/{id}/activities/search-full", params)
		},
	}
	cmd.Flags().String("query", "", "Search query (required)")
	cmd.Flags().Int("limit", 0, "Max results")
	return cmd
}

func newActivitiesIntervalSearchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "interval-search",
		Short: "Search activities by interval criteria",
		RunE: func(cmd *cobra.Command, args []string) error {
			params := map[string]string{}
			for _, flag := range []string{"min-secs", "max-secs", "min-intensity", "max-intensity"} {
				if !cmd.Flags().Changed(flag) {
					return validationErr(fmt.Errorf("--%s is required", flag))
				}
				v, _ := cmd.Flags().GetInt(flag)
				apiParam := map[string]string{
					"min-secs":      "minSecs",
					"max-secs":      "maxSecs",
					"min-intensity": "minIntensity",
					"max-intensity": "maxIntensity",
				}[flag]
				params[apiParam] = fmt.Sprintf("%d", v)
			}
			if v, _ := cmd.Flags().GetString("type"); v != "" {
				params["type"] = v
			}
			if v, _ := cmd.Flags().GetInt("limit"); v > 0 {
				params["limit"] = fmt.Sprintf("%d", v)
			}
			return doGet(cmd, "/api/v1/athlete/{id}/activities/interval-search", params)
		},
	}
	cmd.Flags().Int("min-secs", 0, "Minimum interval duration in seconds (required)")
	cmd.Flags().Int("max-secs", 0, "Maximum interval duration in seconds (required)")
	cmd.Flags().Int("min-intensity", 0, "Minimum intensity (required)")
	cmd.Flags().Int("max-intensity", 0, "Maximum intensity (required)")
	cmd.Flags().String("type", "", "Activity type filter")
	cmd.Flags().Int("limit", 0, "Max results")
	return cmd
}

func newActivitiesCreateManualCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-manual",
		Short: "Create a manual activity",
		RunE: func(cmd *cobra.Command, args []string) error {
			jsonBody, err := requireJSON(cmd)
			if err != nil {
				return err
			}
			return doMutate(cmd, "POST", "/api/v1/athlete/{id}/activities/manual", nil, jsonBody)
		},
	}
	return cmd
}

func newActivitiesCreateManualBulkCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-manual-bulk",
		Short: "Create multiple manual activities",
		RunE: func(cmd *cobra.Command, args []string) error {
			jsonBody, err := requireJSON(cmd)
			if err != nil {
				return err
			}
			return doMutate(cmd, "POST", "/api/v1/athlete/{id}/activities/manual/bulk", nil, jsonBody)
		},
	}
	return cmd
}

func newActivitiesUploadCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upload",
		Short: "Upload an activity file (FIT, TCX, GPX)",
		RunE: func(cmd *cobra.Command, args []string) error {
			file, data, err := requireUploadFile(cmd)
			if err != nil {
				return err
			}

			params := map[string]string{}
			if v, _ := cmd.Flags().GetString("name"); v != "" {
				params["name"] = v
			}

			c := api.FromContext(cmd.Context())

			dryRun, _ := cmd.Flags().GetBool("dry-run")
			if dryRun {
				out, err := c.DryRun("POST", "/api/v1/athlete/{id}/activities", params, []byte(fmt.Sprintf("<binary file: %s, %d bytes>", file, len(data))))
				if err != nil {
					return err
				}
				return format.DryRunOutput(os.Stdout, out)
			}

			resp, err := c.Do("POST", "/api/v1/athlete/{id}/activities", params, data)
			if err != nil {
				return err
			}
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			return format.Write(os.Stdout, body, fmtOpts(cmd))
		},
	}
	addUploadFlag(cmd, "Path to activity file FIT/TCX/GPX (required)")
	cmd.Flags().String("name", "", "Activity name")
	return cmd
}

func newActivitiesDownloadCSVCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "download-csv",
		Short: "Download all activities as CSV",
		RunE: func(cmd *cobra.Command, args []string) error {
			output, _ := cmd.Flags().GetString("output")
			return doDownload(cmd, "/api/v1/athlete/{id}/activities.csv", nil, output)
		},
	}
	cmd.Flags().String("output", "", "Output file path (default: stdout)")
	return cmd
}

func newActivitiesListAroundCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-around",
		Short: "List activities around a specific activity",
		RunE: func(cmd *cobra.Command, args []string) error {
			activityID, err := requireString(cmd, "activity-id")
			if err != nil {
				return err
			}
			if err := validate.PathParam("activity-id", activityID); err != nil {
				return err
			}
			params := map[string]string{"activity_id": activityID}
			if v, _ := cmd.Flags().GetInt("limit"); v > 0 {
				params["limit"] = fmt.Sprintf("%d", v)
			}
			return doGet(cmd, "/api/v1/athlete/{id}/activities-around", params)
		},
	}
	cmd.Flags().String("activity-id", "", "Activity ID to search around (required)")
	cmd.Flags().Int("limit", 0, "Max results")
	return cmd
}

func newActivitiesGetMultipleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-multiple",
		Short: "Get multiple activities by IDs",
		RunE: func(cmd *cobra.Command, args []string) error {
			ids, err := requireString(cmd, "ids")
			if err != nil {
				return err
			}
			if err := validate.PathParam("ids", ids); err != nil {
				return err
			}
			params := map[string]string{"ids": ids}
			if v, _ := cmd.Flags().GetBool("intervals"); v {
				params["intervals"] = "true"
			}
			return doGet(cmd, "/api/v1/athlete/{athleteId}/activities/{ids}", params)
		},
	}
	cmd.Flags().String("ids", "", "Comma-separated activity IDs (required)")
	cmd.Flags().Bool("intervals", false, "Include interval data")
	return cmd
}

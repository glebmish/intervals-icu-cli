// internal/cmd/activities.go
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

// Helper to get format options from command flags.
func fmtOpts(cmd *cobra.Command) format.Options {
	f, _ := cmd.Flags().GetString("format")
	fields, _ := cmd.Flags().GetString("fields")
	return format.FormatFromFlag(f, fields)
}

// Helper to handle the common GET pattern.
func doGet(cmd *cobra.Command, path string, params map[string]string) error {
	c := api.FromContext(cmd.Context())

	dryRun, _ := cmd.Flags().GetBool("dry-run")
	if dryRun {
		fmt.Print(c.DryRun("GET", path, params, nil))
		return nil
	}

	resp, err := c.Do("GET", path, params, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response: %w", err)
	}
	return format.Write(os.Stdout, body, fmtOpts(cmd))
}

// Helper to handle the common POST/PUT pattern.
func doMutate(cmd *cobra.Command, method, path string, params map[string]string, jsonFlag string) error {
	c := api.FromContext(cmd.Context())

	var body []byte
	if jsonFlag != "" {
		if err := validate.JSONBody(jsonFlag); err != nil {
			return err
		}
		body = []byte(jsonFlag)
	}

	dryRun, _ := cmd.Flags().GetBool("dry-run")
	if dryRun {
		fmt.Print(c.DryRun(method, path, params, body))
		return nil
	}

	resp, err := c.Do(method, path, params, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response: %w", err)
	}
	return format.Write(os.Stdout, respBody, fmtOpts(cmd))
}

// Helper to handle DELETE with confirmation.
func doDelete(cmd *cobra.Command, path string, params map[string]string, resource, id string) error {
	if err := confirmDelete(cmd, resource, id); err != nil {
		return err
	}

	c := api.FromContext(cmd.Context())

	dryRun, _ := cmd.Flags().GetBool("dry-run")
	if dryRun {
		fmt.Print(c.DryRun("DELETE", path, params, nil))
		return nil
	}

	resp, err := c.Do("DELETE", path, params, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response: %w", err)
	}
	if len(body) > 0 {
		return format.Write(os.Stdout, body, fmtOpts(cmd))
	}
	return nil
}

// Helper to handle binary file downloads.
func doDownload(cmd *cobra.Command, path string, params map[string]string, outputPath string) error {
	c := api.FromContext(cmd.Context())

	dryRun, _ := cmd.Flags().GetBool("dry-run")
	if dryRun {
		fmt.Print(c.DryRun("GET", path, params, nil))
		return nil
	}

	resp, err := c.Do("GET", path, params, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response: %w", err)
	}

	if outputPath != "" {
		if err := os.WriteFile(outputPath, body, 0644); err != nil {
			return fmt.Errorf("writing file: %w", err)
		}
		fmt.Fprintf(os.Stderr, "Downloaded to %s\n", outputPath)
		return nil
	}

	return format.WriteRaw(os.Stdout, body)
}

func newActivitiesListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List athlete activities",
		RunE: func(cmd *cobra.Command, args []string) error {
			oldest, _ := cmd.Flags().GetString("oldest")
			if oldest == "" {
				return fmt.Errorf("--oldest is required (YYYY-MM-DD)")
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
			q, _ := cmd.Flags().GetString("query")
			if q == "" {
				return fmt.Errorf("--query is required")
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
			q, _ := cmd.Flags().GetString("query")
			if q == "" {
				return fmt.Errorf("--query is required")
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
				v, _ := cmd.Flags().GetInt(flag)
				if v == 0 {
					return fmt.Errorf("--%s is required", flag)
				}
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
			jsonBody, _ := cmd.Flags().GetString("json")
			if jsonBody == "" {
				return fmt.Errorf("--json is required with Activity payload")
			}
			return doMutate(cmd, "POST", "/api/v1/athlete/{id}/activities/manual", nil, jsonBody)
		},
	}
	cmd.Flags().String("json", "", "Activity JSON payload")
	return cmd
}

func newActivitiesCreateManualBulkCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-manual-bulk",
		Short: "Create multiple manual activities",
		RunE: func(cmd *cobra.Command, args []string) error {
			jsonBody, _ := cmd.Flags().GetString("json")
			if jsonBody == "" {
				return fmt.Errorf("--json is required with array of Activity payloads")
			}
			return doMutate(cmd, "POST", "/api/v1/athlete/{id}/activities/manual/bulk", nil, jsonBody)
		},
	}
	cmd.Flags().String("json", "", "JSON array of Activity payloads")
	return cmd
}

func newActivitiesUploadCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upload",
		Short: "Upload an activity file (FIT, TCX, GPX)",
		RunE: func(cmd *cobra.Command, args []string) error {
			file, _ := cmd.Flags().GetString("file")
			if file == "" {
				return fmt.Errorf("--file is required")
			}

			data, err := os.ReadFile(file)
			if err != nil {
				return fmt.Errorf("reading file %s: %w", file, err)
			}

			params := map[string]string{}
			if v, _ := cmd.Flags().GetString("name"); v != "" {
				params["name"] = v
			}

			c := api.FromContext(cmd.Context())

			dryRun, _ := cmd.Flags().GetBool("dry-run")
			if dryRun {
				fmt.Print(c.DryRun("POST", "/api/v1/athlete/{id}/activities", params, []byte(fmt.Sprintf("<binary file: %s, %d bytes>", file, len(data)))))
				return nil
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
	cmd.Flags().String("file", "", "Path to activity file (required)")
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
			activityID, _ := cmd.Flags().GetString("activity-id")
			if activityID == "" {
				return fmt.Errorf("--activity-id is required")
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
			ids, _ := cmd.Flags().GetString("ids")
			if ids == "" {
				return fmt.Errorf("--ids is required (comma-separated)")
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

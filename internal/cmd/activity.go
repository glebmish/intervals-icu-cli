// internal/cmd/activity.go
package cmd

import (
	"fmt"

	"github.com/glebmish/intervals-icu-cli/internal/validate"
	"github.com/spf13/cobra"
)

var activityCmd = &cobra.Command{
	Use:   "activity",
	Short: "Activity analysis — power curves, HR, pace, streams, intervals",
}

func init() {
	activityCmd.AddCommand(
		newActivityGetCmd(),
		newActivityUpdateCmd(),
		newActivityDeleteCmd(),
		newActivityIntervalsCmd(),
		newActivityUpdateIntervalsCmd(),
		newActivityUpdateIntervalCmd(),
		newActivityDeleteIntervalsCmd(),
		newActivitySplitIntervalCmd(),
		newActivityStreamsCmd(),
		newActivityUpdateStreamsCmd(),
		newActivityUploadStreamsCSVCmd(),
		newActivityPowerCurveCmd(),
		newActivityPowerCurvesCmd(),
		newActivityPowerHistogramCmd(),
		newActivityPowerSpikeModelCmd(),
		newActivityPowerVsHRCmd(),
		newActivityHRCurveCmd(),
		newActivityHRHistogramCmd(),
		newActivityHRLoadModelCmd(),
		newActivityPaceCurveCmd(),
		newActivityPaceHistogramCmd(),
		newActivityGapHistogramCmd(),
		newActivityIntervalStatsCmd(),
		newActivitySegmentsCmd(),
		newActivityMapCmd(),
		newActivityWeatherSummaryCmd(),
		newActivityTimeAtHRCmd(),
		newActivityBestEffortsCmd(),
		newActivityMessagesCmd(),
		newActivitySendMessageCmd(),
		newActivityFitFileCmd(),
		newActivityGpxFileCmd(),
		newActivityFileCmd(),
	)
	rootCmd.AddCommand(activityCmd)
}

func requireActivityID(cmd *cobra.Command) (string, error) {
	id, err := requireString(cmd, "activity-id")
	if err != nil {
		return "", err
	}
	if err := validate.PathParam("activity-id", id); err != nil {
		return "", err
	}
	return id, nil
}

// 1. GET /api/v1/activity/{activityId}
func newActivityGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get an activity by ID",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := requireActivityID(cmd)
			if err != nil {
				return err
			}
			params := map[string]string{"activityId": id}
			if v, _ := cmd.Flags().GetBool("intervals"); v {
				params["intervals"] = "true"
			}
			return doGet(cmd, "/api/v1/activity/{activityId}", params)
		},
	}
	cmd.Flags().String("activity-id", "", "Activity ID (required)")
	cmd.Flags().Bool("intervals", false, "Include interval data")
	return cmd
}

// 2. PUT /api/v1/activity/{activityId}
func newActivityUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update an activity",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := requireActivityID(cmd)
			if err != nil {
				return err
			}
			jsonBody, err := requireJSON(cmd)
			if err != nil {
				return err
			}
			params := map[string]string{"activityId": id}
			return doMutate(cmd, "PUT", "/api/v1/activity/{activityId}", params, jsonBody)
		},
	}
	cmd.Flags().String("activity-id", "", "Activity ID (required)")
	return cmd
}

// 3. DELETE /api/v1/activity/{activityId}
func newActivityDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete an activity",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := requireActivityID(cmd)
			if err != nil {
				return err
			}
			params := map[string]string{"activityId": id}
			return doDelete(cmd, "/api/v1/activity/{activityId}", params, "activity", id)
		},
	}
	cmd.Flags().String("activity-id", "", "Activity ID (required)")
	return cmd
}

// 4. GET /api/v1/activity/{activityId}/intervals
func newActivityIntervalsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "intervals",
		Short: "Get intervals for an activity",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := requireActivityID(cmd)
			if err != nil {
				return err
			}
			params := map[string]string{"activityId": id}
			return doGet(cmd, "/api/v1/activity/{activityId}/intervals", params)
		},
	}
	cmd.Flags().String("activity-id", "", "Activity ID (required)")
	return cmd
}

// 5. PUT /api/v1/activity/{activityId}/intervals
func newActivityUpdateIntervalsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-intervals",
		Short: "Update intervals for an activity",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := requireActivityID(cmd)
			if err != nil {
				return err
			}
			jsonBody, err := requireJSON(cmd)
			if err != nil {
				return err
			}
			params := map[string]string{"activityId": id}
			if v, _ := cmd.Flags().GetBool("all"); v {
				params["all"] = "true"
			}
			return doMutate(cmd, "PUT", "/api/v1/activity/{activityId}/intervals", params, jsonBody)
		},
	}
	cmd.Flags().String("activity-id", "", "Activity ID (required)")
	cmd.Flags().Bool("all", false, "Update all intervals")
	return cmd
}

// 6. PUT /api/v1/activity/{activityId}/intervals/{intervalId}
func newActivityUpdateIntervalCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-interval",
		Short: "Update a single interval",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := requireActivityID(cmd)
			if err != nil {
				return err
			}
			intervalID, err := requireString(cmd, "interval-id")
			if err != nil {
				return err
			}
			if err := validate.PathParam("interval-id", intervalID); err != nil {
				return err
			}
			jsonBody, err := requireJSON(cmd)
			if err != nil {
				return err
			}
			params := map[string]string{"activityId": id, "intervalId": intervalID}
			return doMutate(cmd, "PUT", "/api/v1/activity/{activityId}/intervals/{intervalId}", params, jsonBody)
		},
	}
	cmd.Flags().String("activity-id", "", "Activity ID (required)")
	cmd.Flags().String("interval-id", "", "Interval ID (required)")
	return cmd
}

// 7. PUT /api/v1/activity/{activityId}/delete-intervals
func newActivityDeleteIntervalsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-intervals",
		Short: "Delete specific intervals from an activity",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := requireActivityID(cmd)
			if err != nil {
				return err
			}
			jsonBody, err := requireJSON(cmd)
			if err != nil {
				return err
			}
			params := map[string]string{"activityId": id}
			return doMutate(cmd, "PUT", "/api/v1/activity/{activityId}/delete-intervals", params, jsonBody)
		},
	}
	cmd.Flags().String("activity-id", "", "Activity ID (required)")
	return cmd
}

// 8. PUT /api/v1/activity/{activityId}/split-interval
func newActivitySplitIntervalCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "split-interval",
		Short: "Split an interval at a specific point",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := requireActivityID(cmd)
			if err != nil {
				return err
			}
			splitAt, _ := cmd.Flags().GetInt("split-at")
			if splitAt == 0 {
				return validationErr(fmt.Errorf("--split-at is required"))
			}
			params := map[string]string{
				"activityId": id,
				"splitAt":    fmt.Sprintf("%d", splitAt),
			}
			return doMutate(cmd, "PUT", "/api/v1/activity/{activityId}/split-interval", params, "")
		},
	}
	cmd.Flags().String("activity-id", "", "Activity ID (required)")
	cmd.Flags().Int("split-at", 0, "Split point index (required)")
	return cmd
}

// 9. GET /api/v1/activity/{activityId}/streams{ext}
func newActivityStreamsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "streams",
		Short: "Get activity streams",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := requireActivityID(cmd)
			if err != nil {
				return err
			}
			ext, _ := cmd.Flags().GetString("ext")
			params := map[string]string{"activityId": id, "ext": ext}
			if v, _ := cmd.Flags().GetString("types"); v != "" {
				params["types"] = v
			}
			return doGet(cmd, "/api/v1/activity/{activityId}/streams{ext}", params)
		},
	}
	cmd.Flags().String("activity-id", "", "Activity ID (required)")
	cmd.Flags().String("ext", "", "File extension (e.g. .csv)")
	cmd.Flags().String("types", "", "Comma-separated stream types")
	return cmd
}

// 10. PUT /api/v1/activity/{activityId}/streams
func newActivityUpdateStreamsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-streams",
		Short: "Update activity streams",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := requireActivityID(cmd)
			if err != nil {
				return err
			}
			jsonBody, err := requireJSON(cmd)
			if err != nil {
				return err
			}
			params := map[string]string{"activityId": id}
			return doMutate(cmd, "PUT", "/api/v1/activity/{activityId}/streams", params, jsonBody)
		},
	}
	cmd.Flags().String("activity-id", "", "Activity ID (required)")
	return cmd
}

// 11. PUT /api/v1/activity/{activityId}/streams.csv
func newActivityUploadStreamsCSVCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upload-streams-csv",
		Short: "Upload activity streams as CSV",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := requireActivityID(cmd)
			if err != nil {
				return err
			}
			jsonBody, err := requireJSON(cmd)
			if err != nil {
				return err
			}
			params := map[string]string{"activityId": id}
			return doMutate(cmd, "PUT", "/api/v1/activity/{activityId}/streams.csv", params, jsonBody)
		},
	}
	cmd.Flags().String("activity-id", "", "Activity ID (required)")
	return cmd
}

// 12. GET /api/v1/activity/{activityId}/power-curve{ext}
func newActivityPowerCurveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "power-curve",
		Short: "Get activity power curve",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := requireActivityID(cmd)
			if err != nil {
				return err
			}
			ext, _ := cmd.Flags().GetString("ext")
			params := map[string]string{"activityId": id, "ext": ext}
			if v, _ := cmd.Flags().GetBool("fatigue"); v {
				params["fatigue"] = "true"
			}
			return doGet(cmd, "/api/v1/activity/{activityId}/power-curve{ext}", params)
		},
	}
	cmd.Flags().String("activity-id", "", "Activity ID (required)")
	cmd.Flags().String("ext", "", "File extension (e.g. .csv)")
	cmd.Flags().Bool("fatigue", false, "Include fatigue data")
	return cmd
}

// 13. GET /api/v1/activity/{activityId}/power-curves{ext}
func newActivityPowerCurvesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "power-curves",
		Short: "Get activity power curves",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := requireActivityID(cmd)
			if err != nil {
				return err
			}
			ext, _ := cmd.Flags().GetString("ext")
			params := map[string]string{"activityId": id, "ext": ext}
			if v, _ := cmd.Flags().GetString("types"); v != "" {
				params["types"] = v
			}
			return doGet(cmd, "/api/v1/activity/{activityId}/power-curves{ext}", params)
		},
	}
	cmd.Flags().String("activity-id", "", "Activity ID (required)")
	cmd.Flags().String("ext", "", "File extension (e.g. .csv)")
	cmd.Flags().String("types", "", "Comma-separated curve types")
	return cmd
}

// 14. GET /api/v1/activity/{activityId}/power-histogram
func newActivityPowerHistogramCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "power-histogram",
		Short: "Get activity power histogram",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := requireActivityID(cmd)
			if err != nil {
				return err
			}
			params := map[string]string{"activityId": id}
			if v, _ := cmd.Flags().GetInt("bucket-size"); v > 0 {
				params["bucketSize"] = fmt.Sprintf("%d", v)
			}
			return doGet(cmd, "/api/v1/activity/{activityId}/power-histogram", params)
		},
	}
	cmd.Flags().String("activity-id", "", "Activity ID (required)")
	cmd.Flags().Int("bucket-size", 0, "Histogram bucket size in watts")
	return cmd
}

// 15. GET /api/v1/activity/{activityId}/power-spike-model
func newActivityPowerSpikeModelCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "power-spike-model",
		Short: "Get activity power spike model",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := requireActivityID(cmd)
			if err != nil {
				return err
			}
			params := map[string]string{"activityId": id}
			return doGet(cmd, "/api/v1/activity/{activityId}/power-spike-model", params)
		},
	}
	cmd.Flags().String("activity-id", "", "Activity ID (required)")
	return cmd
}

// 16. GET /api/v1/activity/{activityId}/power-vs-hr{ext}
func newActivityPowerVsHRCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "power-vs-hr",
		Short: "Get activity power vs heart rate data",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := requireActivityID(cmd)
			if err != nil {
				return err
			}
			ext, _ := cmd.Flags().GetString("ext")
			params := map[string]string{"activityId": id, "ext": ext}
			return doGet(cmd, "/api/v1/activity/{activityId}/power-vs-hr{ext}", params)
		},
	}
	cmd.Flags().String("activity-id", "", "Activity ID (required)")
	cmd.Flags().String("ext", "", "File extension (e.g. .csv)")
	return cmd
}

// 17. GET /api/v1/activity/{activityId}/hr-curve{ext}
func newActivityHRCurveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hr-curve",
		Short: "Get activity heart rate curve",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := requireActivityID(cmd)
			if err != nil {
				return err
			}
			ext, _ := cmd.Flags().GetString("ext")
			params := map[string]string{"activityId": id, "ext": ext}
			return doGet(cmd, "/api/v1/activity/{activityId}/hr-curve{ext}", params)
		},
	}
	cmd.Flags().String("activity-id", "", "Activity ID (required)")
	cmd.Flags().String("ext", "", "File extension (e.g. .csv)")
	return cmd
}

// 18. GET /api/v1/activity/{activityId}/hr-histogram
func newActivityHRHistogramCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hr-histogram",
		Short: "Get activity heart rate histogram",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := requireActivityID(cmd)
			if err != nil {
				return err
			}
			params := map[string]string{"activityId": id}
			if v, _ := cmd.Flags().GetInt("bucket-size"); v > 0 {
				params["bucketSize"] = fmt.Sprintf("%d", v)
			}
			return doGet(cmd, "/api/v1/activity/{activityId}/hr-histogram", params)
		},
	}
	cmd.Flags().String("activity-id", "", "Activity ID (required)")
	cmd.Flags().Int("bucket-size", 0, "Histogram bucket size in bpm")
	return cmd
}

// 19. GET /api/v1/activity/{activityId}/hr-load-model
func newActivityHRLoadModelCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hr-load-model",
		Short: "Get activity HR load model",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := requireActivityID(cmd)
			if err != nil {
				return err
			}
			params := map[string]string{"activityId": id}
			return doGet(cmd, "/api/v1/activity/{activityId}/hr-load-model", params)
		},
	}
	cmd.Flags().String("activity-id", "", "Activity ID (required)")
	return cmd
}

// 20. GET /api/v1/activity/{activityId}/pace-curve{ext}
func newActivityPaceCurveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pace-curve",
		Short: "Get activity pace curve",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := requireActivityID(cmd)
			if err != nil {
				return err
			}
			ext, _ := cmd.Flags().GetString("ext")
			params := map[string]string{"activityId": id, "ext": ext}
			if v, _ := cmd.Flags().GetBool("gap"); v {
				params["gap"] = "true"
			}
			return doGet(cmd, "/api/v1/activity/{activityId}/pace-curve{ext}", params)
		},
	}
	cmd.Flags().String("activity-id", "", "Activity ID (required)")
	cmd.Flags().String("ext", "", "File extension (e.g. .csv)")
	cmd.Flags().Bool("gap", false, "Use grade-adjusted pace")
	return cmd
}

// 21. GET /api/v1/activity/{activityId}/pace-histogram
func newActivityPaceHistogramCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pace-histogram",
		Short: "Get activity pace histogram",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := requireActivityID(cmd)
			if err != nil {
				return err
			}
			params := map[string]string{"activityId": id}
			return doGet(cmd, "/api/v1/activity/{activityId}/pace-histogram", params)
		},
	}
	cmd.Flags().String("activity-id", "", "Activity ID (required)")
	return cmd
}

// 22. GET /api/v1/activity/{activityId}/gap-histogram
func newActivityGapHistogramCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gap-histogram",
		Short: "Get activity grade-adjusted pace histogram",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := requireActivityID(cmd)
			if err != nil {
				return err
			}
			params := map[string]string{"activityId": id}
			return doGet(cmd, "/api/v1/activity/{activityId}/gap-histogram", params)
		},
	}
	cmd.Flags().String("activity-id", "", "Activity ID (required)")
	return cmd
}

// 23. GET /api/v1/activity/{activityId}/interval-stats
func newActivityIntervalStatsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "interval-stats",
		Short: "Get interval statistics for an activity",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := requireActivityID(cmd)
			if err != nil {
				return err
			}
			params := map[string]string{"activityId": id}
			if v, _ := cmd.Flags().GetInt("start-index"); v > 0 {
				params["startIndex"] = fmt.Sprintf("%d", v)
			}
			if v, _ := cmd.Flags().GetInt("end-index"); v > 0 {
				params["endIndex"] = fmt.Sprintf("%d", v)
			}
			return doGet(cmd, "/api/v1/activity/{activityId}/interval-stats", params)
		},
	}
	cmd.Flags().String("activity-id", "", "Activity ID (required)")
	cmd.Flags().Int("start-index", 0, "Start stream index")
	cmd.Flags().Int("end-index", 0, "End stream index")
	return cmd
}

// 24. GET /api/v1/activity/{activityId}/segments
func newActivitySegmentsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "segments",
		Short: "Get segments for an activity",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := requireActivityID(cmd)
			if err != nil {
				return err
			}
			params := map[string]string{"activityId": id}
			return doGet(cmd, "/api/v1/activity/{activityId}/segments", params)
		},
	}
	cmd.Flags().String("activity-id", "", "Activity ID (required)")
	return cmd
}

// 25. GET /api/v1/activity/{activityId}/map
func newActivityMapCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "map",
		Short: "Get activity map data",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := requireActivityID(cmd)
			if err != nil {
				return err
			}
			params := map[string]string{"activityId": id}
			if v, _ := cmd.Flags().GetBool("bounds-only"); v {
				params["boundsOnly"] = "true"
			}
			if v, _ := cmd.Flags().GetBool("weather"); v {
				params["weather"] = "true"
			}
			return doGet(cmd, "/api/v1/activity/{activityId}/map", params)
		},
	}
	cmd.Flags().String("activity-id", "", "Activity ID (required)")
	cmd.Flags().Bool("bounds-only", false, "Return only map bounds")
	cmd.Flags().Bool("weather", false, "Include weather data")
	return cmd
}

// 26. GET /api/v1/activity/{activityId}/weather-summary
func newActivityWeatherSummaryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "weather-summary",
		Short: "Get activity weather summary",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := requireActivityID(cmd)
			if err != nil {
				return err
			}
			params := map[string]string{"activityId": id}
			if v, _ := cmd.Flags().GetInt("start-index"); v > 0 {
				params["startIndex"] = fmt.Sprintf("%d", v)
			}
			if v, _ := cmd.Flags().GetInt("end-index"); v > 0 {
				params["endIndex"] = fmt.Sprintf("%d", v)
			}
			return doGet(cmd, "/api/v1/activity/{activityId}/weather-summary", params)
		},
	}
	cmd.Flags().String("activity-id", "", "Activity ID (required)")
	cmd.Flags().Int("start-index", 0, "Start stream index")
	cmd.Flags().Int("end-index", 0, "End stream index")
	return cmd
}

// 27. GET /api/v1/activity/{activityId}/time-at-hr
func newActivityTimeAtHRCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "time-at-hr",
		Short: "Get time at heart rate zones for an activity",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := requireActivityID(cmd)
			if err != nil {
				return err
			}
			params := map[string]string{"activityId": id}
			return doGet(cmd, "/api/v1/activity/{activityId}/time-at-hr", params)
		},
	}
	cmd.Flags().String("activity-id", "", "Activity ID (required)")
	return cmd
}

// 28. GET /api/v1/activity/{activityId}/best-efforts
func newActivityBestEffortsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "best-efforts",
		Short: "Get best efforts for an activity",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := requireActivityID(cmd)
			if err != nil {
				return err
			}
			stream, err := requireString(cmd, "stream")
			if err != nil {
				return err
			}
			params := map[string]string{"activityId": id, "stream": stream}
			if v, _ := cmd.Flags().GetInt("duration"); v > 0 {
				params["duration"] = fmt.Sprintf("%d", v)
			}
			if v, _ := cmd.Flags().GetFloat64("distance"); v > 0 {
				params["distance"] = fmt.Sprintf("%g", v)
			}
			if v, _ := cmd.Flags().GetInt("count"); v > 0 {
				params["count"] = fmt.Sprintf("%d", v)
			}
			return doGet(cmd, "/api/v1/activity/{activityId}/best-efforts", params)
		},
	}
	cmd.Flags().String("activity-id", "", "Activity ID (required)")
	cmd.Flags().String("stream", "", "Stream type for best efforts (required)")
	cmd.Flags().Int("duration", 0, "Duration in seconds")
	cmd.Flags().Float64("distance", 0, "Distance in meters")
	cmd.Flags().Int("count", 0, "Number of results")
	return cmd
}

// 29. GET /api/v1/activity/{activityId}/messages
func newActivityMessagesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "messages",
		Short: "Get messages for an activity",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := requireActivityID(cmd)
			if err != nil {
				return err
			}
			params := map[string]string{"activityId": id}
			if v, _ := cmd.Flags().GetInt("limit"); v > 0 {
				params["limit"] = fmt.Sprintf("%d", v)
			}
			return doGet(cmd, "/api/v1/activity/{activityId}/messages", params)
		},
	}
	cmd.Flags().String("activity-id", "", "Activity ID (required)")
	cmd.Flags().Int("limit", 0, "Max number of messages")
	return cmd
}

// 30. POST /api/v1/activity/{activityId}/messages
func newActivitySendMessageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send-message",
		Short: "Send a message on an activity",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := requireActivityID(cmd)
			if err != nil {
				return err
			}
			jsonBody, err := requireJSON(cmd)
			if err != nil {
				return err
			}
			params := map[string]string{"activityId": id}
			return doMutate(cmd, "POST", "/api/v1/activity/{activityId}/messages", params, jsonBody)
		},
	}
	cmd.Flags().String("activity-id", "", "Activity ID (required)")
	return cmd
}

// 31. GET /api/v1/activity/{activityId}/fit-file (download)
func newActivityFitFileCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fit-file",
		Short: "Download the FIT file for an activity",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := requireActivityID(cmd)
			if err != nil {
				return err
			}
			output, _ := cmd.Flags().GetString("output")
			params := map[string]string{"activityId": id}
			return doDownload(cmd, "/api/v1/activity/{activityId}/fit-file", params, output)
		},
	}
	cmd.Flags().String("activity-id", "", "Activity ID (required)")
	cmd.Flags().String("output", "", "Output file path (default: stdout)")
	return cmd
}

// 32. GET /api/v1/activity/{activityId}/gpx-file (download)
func newActivityGpxFileCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gpx-file",
		Short: "Download the GPX file for an activity",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := requireActivityID(cmd)
			if err != nil {
				return err
			}
			output, _ := cmd.Flags().GetString("output")
			params := map[string]string{"activityId": id}
			if v, _ := cmd.Flags().GetBool("power"); v {
				params["power"] = "true"
			}
			if v, _ := cmd.Flags().GetBool("hr"); v {
				params["hr"] = "true"
			}
			return doDownload(cmd, "/api/v1/activity/{activityId}/gpx-file", params, output)
		},
	}
	cmd.Flags().String("activity-id", "", "Activity ID (required)")
	cmd.Flags().String("output", "", "Output file path (default: stdout)")
	cmd.Flags().Bool("power", false, "Include power data in GPX")
	cmd.Flags().Bool("hr", false, "Include heart rate data in GPX")
	return cmd
}

// 33. GET /api/v1/activity/{activityId}/file (download)
func newActivityFileCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "file",
		Short: "Download the original file for an activity",
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := requireActivityID(cmd)
			if err != nil {
				return err
			}
			output, _ := cmd.Flags().GetString("output")
			params := map[string]string{"activityId": id}
			return doDownload(cmd, "/api/v1/activity/{activityId}/file", params, output)
		},
	}
	cmd.Flags().String("activity-id", "", "Activity ID (required)")
	cmd.Flags().String("output", "", "Output file path (default: stdout)")
	return cmd
}

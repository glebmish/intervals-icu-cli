// internal/cmd/athlete.go
package cmd

import (
	"github.com/glebmish/intervals-icu-cli/internal/validate"
	"github.com/spf13/cobra"
)

var athleteCmd = &cobra.Command{
	Use:   "athlete",
	Short: "Athlete profile, settings, curves, routes, fitness model",
}

func init() {
	athleteCmd.AddCommand(
		newAthleteGetCmd(),
		newAthleteUpdateCmd(),
		newAthleteProfileCmd(),
		newAthleteSummaryCmd(),
		newAthleteTrainingPlanCmd(),
		newAthleteUpdateTrainingPlanCmd(),
		newAthleteApplyPlanChangesCmd(),
		newAthleteWeatherConfigCmd(),
		newAthleteUpdateWeatherConfigCmd(),
		newAthleteWeatherForecastCmd(),
		newAthleteRoutesCmd(),
		newAthleteRouteGetCmd(),
		newAthleteRouteUpdateCmd(),
		newAthleteRouteSimilarityCmd(),
		newAthleteChatsCmd(),
		newAthleteFitnessModelEventsCmd(),
		newAthletePowerHRCurveCmd(),
		newAthletePowerCurvesCmd(),
		newAthletePaceCurvesCmd(),
		newAthleteHRCurvesCmd(),
		newAthleteMmpModelCmd(),
		newAthleteActivityPowerCurvesCmd(),
		newAthleteActivityPaceCurvesCmd(),
		newAthleteActivityHRCurvesCmd(),
		newAthleteDuplicateEventsCmd(),
		newAthleteDownloadWorkoutCmd(),
		newAthleteDownloadFitFilesCmd(),
		newAthleteTagsCmd(),
	)
	rootCmd.AddCommand(athleteCmd)
}

func newAthleteGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get athlete profile",
		RunE: func(cmd *cobra.Command, args []string) error {
			return doGet(cmd, "/api/v1/athlete/{id}", nil)
		},
	}
	return cmd
}

func newAthleteUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update athlete profile",
		RunE: func(cmd *cobra.Command, args []string) error {
			jsonBody, err := requireJSON(cmd)
			if err != nil {
				return err
			}
			return doMutate(cmd, "PUT", "/api/v1/athlete/{id}", nil, jsonBody)
		},
	}
	return cmd
}

func newAthleteProfileCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "profile",
		Short: "Get athlete profile details",
		RunE: func(cmd *cobra.Command, args []string) error {
			return doGet(cmd, "/api/v1/athlete/{id}/profile", nil)
		},
	}
	return cmd
}

func newAthleteSummaryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "summary",
		Short: "Get athlete summary",
		RunE: func(cmd *cobra.Command, args []string) error {
			ext, _ := cmd.Flags().GetString("ext")
			params := map[string]string{}
			if ext != "" {
				params["ext"] = ext
			}
			if v, _ := cmd.Flags().GetString("start"); v != "" {
				if err := validate.DateParam("start", v); err != nil {
					return err
				}
				params["start"] = v
			}
			if v, _ := cmd.Flags().GetString("end"); v != "" {
				if err := validate.DateParam("end", v); err != nil {
					return err
				}
				params["end"] = v
			}
			if v, _ := cmd.Flags().GetString("tags"); v != "" {
				params["tags"] = v
			}
			path := "/api/v1/athlete/{id}/athlete-summary" + ext
			return doGet(cmd, path, params)
		},
	}
	cmd.Flags().String("ext", "", "Format extension (e.g. .csv)")
	cmd.Flags().String("start", "", "Start date (YYYY-MM-DD)")
	cmd.Flags().String("end", "", "End date (YYYY-MM-DD)")
	cmd.Flags().String("tags", "", "Tags filter")
	return cmd
}

func newAthleteTrainingPlanCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "training-plan",
		Short: "Get athlete training plan",
		RunE: func(cmd *cobra.Command, args []string) error {
			return doGet(cmd, "/api/v1/athlete/{id}/training-plan", nil)
		},
	}
	return cmd
}

func newAthleteUpdateTrainingPlanCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-training-plan",
		Short: "Update athlete training plan",
		RunE: func(cmd *cobra.Command, args []string) error {
			jsonBody, err := requireJSON(cmd)
			if err != nil {
				return err
			}
			return doMutate(cmd, "PUT", "/api/v1/athlete/{id}/training-plan", nil, jsonBody)
		},
	}
	return cmd
}

func newAthleteApplyPlanChangesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "apply-plan-changes",
		Short: "Apply pending plan changes",
		RunE: func(cmd *cobra.Command, args []string) error {
			return doMutate(cmd, "PUT", "/api/v1/athlete/{id}/apply-plan-changes", nil, "")
		},
	}
	return cmd
}

func newAthleteWeatherConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "weather-config",
		Short: "Get athlete weather configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			return doGet(cmd, "/api/v1/athlete/{id}/weather-config", nil)
		},
	}
	return cmd
}

func newAthleteUpdateWeatherConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-weather-config",
		Short: "Update athlete weather configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			jsonBody, err := requireJSON(cmd)
			if err != nil {
				return err
			}
			return doMutate(cmd, "PUT", "/api/v1/athlete/{id}/weather-config", nil, jsonBody)
		},
	}
	return cmd
}

func newAthleteWeatherForecastCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "weather-forecast",
		Short: "Get athlete weather forecast",
		RunE: func(cmd *cobra.Command, args []string) error {
			return doGet(cmd, "/api/v1/athlete/{id}/weather-forecast", nil)
		},
	}
	return cmd
}

func newAthleteRoutesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "routes",
		Short: "List athlete routes",
		RunE: func(cmd *cobra.Command, args []string) error {
			return doGet(cmd, "/api/v1/athlete/{id}/routes", nil)
		},
	}
	return cmd
}

func newAthleteRouteGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "route-get",
		Short: "Get a specific route",
		RunE: func(cmd *cobra.Command, args []string) error {
			routeID, err := requireString(cmd, "route-id")
			if err != nil {
				return err
			}
			if err := validate.PathParam("route-id", routeID); err != nil {
				return err
			}
			params := map[string]string{"route_id": routeID}
			if v, _ := cmd.Flags().GetBool("include-path"); v {
				params["includePath"] = "true"
			}
			return doGet(cmd, "/api/v1/athlete/{id}/routes/{route_id}", params)
		},
	}
	cmd.Flags().String("route-id", "", "Route ID (required)")
	cmd.Flags().Bool("include-path", false, "Include path coordinates")
	return cmd
}

func newAthleteRouteUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "route-update",
		Short: "Update a route",
		RunE: func(cmd *cobra.Command, args []string) error {
			routeID, err := requireString(cmd, "route-id")
			if err != nil {
				return err
			}
			if err := validate.PathParam("route-id", routeID); err != nil {
				return err
			}
			jsonBody, err := requireJSON(cmd)
			if err != nil {
				return err
			}
			params := map[string]string{"route_id": routeID}
			return doMutate(cmd, "PUT", "/api/v1/athlete/{id}/routes/{route_id}", params, jsonBody)
		},
	}
	cmd.Flags().String("route-id", "", "Route ID (required)")
	return cmd
}

func newAthleteRouteSimilarityCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "route-similarity",
		Short: "Get similarity between two routes",
		RunE: func(cmd *cobra.Command, args []string) error {
			routeID, err := requireString(cmd, "route-id")
			if err != nil {
				return err
			}
			if err := validate.PathParam("route-id", routeID); err != nil {
				return err
			}
			otherID, err := requireString(cmd, "other-id")
			if err != nil {
				return err
			}
			if err := validate.PathParam("other-id", otherID); err != nil {
				return err
			}
			params := map[string]string{"route_id": routeID, "other_id": otherID}
			return doGet(cmd, "/api/v1/athlete/{id}/routes/{route_id}/similarity/{other_id}", params)
		},
	}
	cmd.Flags().String("route-id", "", "Route ID (required)")
	cmd.Flags().String("other-id", "", "Other route ID (required)")
	return cmd
}

func newAthleteChatsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "chats",
		Short: "List athlete chats",
		RunE: func(cmd *cobra.Command, args []string) error {
			return doGet(cmd, "/api/v1/athlete/{id}/chats", nil)
		},
	}
	return cmd
}

func newAthleteFitnessModelEventsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fitness-model-events",
		Short: "Get fitness model events",
		RunE: func(cmd *cobra.Command, args []string) error {
			return doGet(cmd, "/api/v1/athlete/{id}/fitness-model-events", nil)
		},
	}
	return cmd
}

func newAthletePowerHRCurveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "power-hr-curve",
		Short: "Get power/HR curve",
		RunE: func(cmd *cobra.Command, args []string) error {
			params := map[string]string{}
			if v, _ := cmd.Flags().GetString("start"); v != "" {
				if err := validate.DateParam("start", v); err != nil {
					return err
				}
				params["start"] = v
			}
			if v, _ := cmd.Flags().GetString("end"); v != "" {
				if err := validate.DateParam("end", v); err != nil {
					return err
				}
				params["end"] = v
			}
			return doGet(cmd, "/api/v1/athlete/{id}/power-hr-curve", params)
		},
	}
	cmd.Flags().String("start", "", "Start date (YYYY-MM-DD)")
	cmd.Flags().String("end", "", "End date (YYYY-MM-DD)")
	return cmd
}

func newAthletePowerCurvesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "power-curves",
		Short: "Get power curves",
		RunE: func(cmd *cobra.Command, args []string) error {
			ext, _ := cmd.Flags().GetString("ext")
			params := map[string]string{}
			if ext != "" {
				params["ext"] = ext
			}
			for _, flag := range []string{"oldest", "newest", "type", "fatigue", "filters", "secs"} {
				if v, _ := cmd.Flags().GetString(flag); v != "" {
					params[flag] = v
				}
			}
			path := "/api/v1/athlete/{id}/power-curves" + ext
			return doGet(cmd, path, params)
		},
	}
	cmd.Flags().String("ext", "", "Format extension (e.g. .csv)")
	cmd.Flags().String("oldest", "", "Oldest date filter")
	cmd.Flags().String("newest", "", "Newest date filter")
	cmd.Flags().String("type", "", "Activity type filter")
	cmd.Flags().String("fatigue", "", "Fatigue filter")
	cmd.Flags().String("filters", "", "Additional filters")
	cmd.Flags().String("secs", "", "Seconds filter")
	return cmd
}

func newAthletePaceCurvesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pace-curves",
		Short: "Get pace curves",
		RunE: func(cmd *cobra.Command, args []string) error {
			ext, _ := cmd.Flags().GetString("ext")
			params := map[string]string{}
			if ext != "" {
				params["ext"] = ext
			}
			for _, flag := range []string{"oldest", "newest", "type", "filters", "distances"} {
				if v, _ := cmd.Flags().GetString(flag); v != "" {
					params[flag] = v
				}
			}
			if v, _ := cmd.Flags().GetBool("gap"); v {
				params["gap"] = "true"
			}
			path := "/api/v1/athlete/{id}/pace-curves" + ext
			return doGet(cmd, path, params)
		},
	}
	cmd.Flags().String("ext", "", "Format extension (e.g. .csv)")
	cmd.Flags().String("oldest", "", "Oldest date filter")
	cmd.Flags().String("newest", "", "Newest date filter")
	cmd.Flags().String("type", "", "Activity type filter")
	cmd.Flags().String("filters", "", "Additional filters")
	cmd.Flags().String("distances", "", "Distances filter")
	cmd.Flags().Bool("gap", false, "Use GAP (grade adjusted pace)")
	return cmd
}

func newAthleteHRCurvesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hr-curves",
		Short: "Get HR curves",
		RunE: func(cmd *cobra.Command, args []string) error {
			ext, _ := cmd.Flags().GetString("ext")
			params := map[string]string{}
			if ext != "" {
				params["ext"] = ext
			}
			for _, flag := range []string{"oldest", "newest", "type", "filters", "secs"} {
				if v, _ := cmd.Flags().GetString(flag); v != "" {
					params[flag] = v
				}
			}
			path := "/api/v1/athlete/{id}/hr-curves" + ext
			return doGet(cmd, path, params)
		},
	}
	cmd.Flags().String("ext", "", "Format extension (e.g. .csv)")
	cmd.Flags().String("oldest", "", "Oldest date filter")
	cmd.Flags().String("newest", "", "Newest date filter")
	cmd.Flags().String("type", "", "Activity type filter")
	cmd.Flags().String("filters", "", "Additional filters")
	cmd.Flags().String("secs", "", "Seconds filter")
	return cmd
}

func newAthleteMmpModelCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mmp-model",
		Short: "Get MMP model",
		RunE: func(cmd *cobra.Command, args []string) error {
			params := map[string]string{}
			if v, _ := cmd.Flags().GetString("type"); v != "" {
				params["type"] = v
			}
			return doGet(cmd, "/api/v1/athlete/{id}/mmp-model", params)
		},
	}
	cmd.Flags().String("type", "", "Model type")
	return cmd
}

func newAthleteActivityPowerCurvesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "activity-power-curves",
		Short: "Get activity power curves",
		RunE: func(cmd *cobra.Command, args []string) error {
			ext, _ := cmd.Flags().GetString("ext")
			params := map[string]string{}
			if ext != "" {
				params["ext"] = ext
			}
			for _, flag := range []string{"oldest", "newest", "filters", "secs", "type", "fatigue"} {
				if v, _ := cmd.Flags().GetString(flag); v != "" {
					params[flag] = v
				}
			}
			path := "/api/v1/athlete/{id}/activity-power-curves" + ext
			return doGet(cmd, path, params)
		},
	}
	cmd.Flags().String("ext", "", "Format extension (e.g. .csv)")
	cmd.Flags().String("oldest", "", "Oldest date filter")
	cmd.Flags().String("newest", "", "Newest date filter")
	cmd.Flags().String("filters", "", "Additional filters")
	cmd.Flags().String("secs", "", "Seconds filter")
	cmd.Flags().String("type", "", "Activity type filter")
	cmd.Flags().String("fatigue", "", "Fatigue filter")
	return cmd
}

func newAthleteActivityPaceCurvesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "activity-pace-curves",
		Short: "Get activity pace curves",
		RunE: func(cmd *cobra.Command, args []string) error {
			ext, _ := cmd.Flags().GetString("ext")
			params := map[string]string{}
			if ext != "" {
				params["ext"] = ext
			}
			for _, flag := range []string{"oldest", "newest", "type", "filters", "distances"} {
				if v, _ := cmd.Flags().GetString(flag); v != "" {
					params[flag] = v
				}
			}
			if v, _ := cmd.Flags().GetBool("gap"); v {
				params["gap"] = "true"
			}
			path := "/api/v1/athlete/{id}/activity-pace-curves" + ext
			return doGet(cmd, path, params)
		},
	}
	cmd.Flags().String("ext", "", "Format extension (e.g. .csv)")
	cmd.Flags().String("oldest", "", "Oldest date filter")
	cmd.Flags().String("newest", "", "Newest date filter")
	cmd.Flags().String("type", "", "Activity type filter")
	cmd.Flags().String("filters", "", "Additional filters")
	cmd.Flags().String("distances", "", "Distances filter")
	cmd.Flags().Bool("gap", false, "Use GAP (grade adjusted pace)")
	return cmd
}

func newAthleteActivityHRCurvesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "activity-hr-curves",
		Short: "Get activity HR curves",
		RunE: func(cmd *cobra.Command, args []string) error {
			ext, _ := cmd.Flags().GetString("ext")
			params := map[string]string{}
			if ext != "" {
				params["ext"] = ext
			}
			for _, flag := range []string{"oldest", "newest", "filters", "secs", "type"} {
				if v, _ := cmd.Flags().GetString(flag); v != "" {
					params[flag] = v
				}
			}
			path := "/api/v1/athlete/{id}/activity-hr-curves" + ext
			return doGet(cmd, path, params)
		},
	}
	cmd.Flags().String("ext", "", "Format extension (e.g. .csv)")
	cmd.Flags().String("oldest", "", "Oldest date filter")
	cmd.Flags().String("newest", "", "Newest date filter")
	cmd.Flags().String("filters", "", "Additional filters")
	cmd.Flags().String("secs", "", "Seconds filter")
	cmd.Flags().String("type", "", "Activity type filter")
	return cmd
}

func newAthleteDuplicateEventsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "duplicate-events",
		Short: "Duplicate events for the athlete",
		RunE: func(cmd *cobra.Command, args []string) error {
			jsonBody, err := requireJSON(cmd)
			if err != nil {
				return err
			}
			return doMutate(cmd, "POST", "/api/v1/athlete/{id}/duplicate-events", nil, jsonBody)
		},
	}
	return cmd
}

func newAthleteDownloadWorkoutCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "download-workout",
		Short: "Download a workout file",
		RunE: func(cmd *cobra.Command, args []string) error {
			ext, _ := cmd.Flags().GetString("ext")
			jsonBody, err := requireJSON(cmd)
			if err != nil {
				return err
			}
			params := map[string]string{}
			if ext != "" {
				params["ext"] = ext
			}
			path := "/api/v1/athlete/{id}/download-workout" + ext
			return doMutate(cmd, "POST", path, params, jsonBody)
		},
	}
	cmd.Flags().String("ext", "", "File extension (e.g. .zwo)")
	return cmd
}

func newAthleteDownloadFitFilesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "download-fit-files",
		Short: "Download FIT files for activities",
		RunE: func(cmd *cobra.Command, args []string) error {
			params := map[string]string{}
			if v, _ := cmd.Flags().GetBool("power"); v {
				params["power"] = "true"
			}
			if v, _ := cmd.Flags().GetBool("hr"); v {
				params["hr"] = "true"
			}
			if v, _ := cmd.Flags().GetString("ids"); v != "" {
				params["ids"] = v
			}
			jsonBody, _ := cmd.Flags().GetString("json")
			return doMutate(cmd, "POST", "/api/v1/athlete/{id}/download-fit-files", params, jsonBody)
		},
	}
	cmd.Flags().Bool("power", false, "Include power data")
	cmd.Flags().Bool("hr", false, "Include heart rate data")
	cmd.Flags().String("ids", "", "Comma-separated activity IDs")
	return cmd
}

func newAthleteTagsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tags",
		Short: "List athlete activity tags",
		RunE: func(cmd *cobra.Command, args []string) error {
			return doGet(cmd, "/api/v1/athlete/{id}/activity-tags", nil)
		},
	}
	return cmd
}

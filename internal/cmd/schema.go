// internal/cmd/schema.go
package cmd

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/glebmish/intervals-icu-cli/internal/cliexit"
	"github.com/spf13/cobra"
)

//go:embed openapi-spec.json
var specData []byte

func discoveryErr(err error) error {
	if err == nil {
		return nil
	}
	return &cliexit.DiscoveryError{Err: err}
}

type operationEntry struct {
	CommandName string
	OperationID string
	Method      string
	Path        string
	PathItem    map[string]interface{}
}

// operationIDToCommand maps spec operationIds → CLI command names.
// The bijection tests in schema_test.go enforce that this map is in sync with
// both the embedded spec (every operationId mapped) and the cobra tree
// (every mapped name resolves to a registered command).
var operationIDToCommand = map[string]string{
	"applyCurrentPlanChanges":        "athlete.apply-plan-changes",
	"applyPlan":                      "events.apply-plan",
	"applyToActivities":              "sport-settings.apply",
	"calcDistanceEtc":                "gear.calc",
	"checkMerge":                     "athlete.route-similarity",
	"createCustomItem":               "custom-items.create",
	"createEvent":                    "events.create",
	"createFolder":                   "folders.create",
	"createGear":                     "gear.create",
	"createManualActivity":           "activities.create-manual",
	"createMultipleEvents":           "events.create-bulk",
	"createMultipleManualActivities": "activities.create-manual-bulk",
	"createMultipleWorkouts":         "workouts.create-bulk",
	"createReminder":                 "gear.create-reminder",
	"createSettings":                 "sport-settings.create",
	"createSharedEvent":              "shared-events.create",
	"createWorkout":                  "workouts.create",
	"duplicateEvents":                "athlete.duplicate-events",
	"deleteActivity":                 "activity.delete",
	"deleteCustomItem":               "custom-items.delete",
	"deleteEvent":                    "events.delete",
	"deleteEvents":                   "events.delete-range",
	"deleteEventsBulk":               "events.delete-bulk",
	"deleteFolder":                   "folders.delete",
	"deleteGear":                     "gear.delete",
	"deleteIntervals":                "activity.delete-intervals",
	"deleteMessage":                  "chats.delete-message",
	"deleteReminder":                 "gear.delete-reminder",
	"deleteSettings":                 "sport-settings.delete",
	"deleteSharedEvent":              "shared-events.delete",
	"deleteWorkout":                  "workouts.delete",
	"disconnectApp":                  "misc.disconnect-app",
	"downloadActivitiesAsCSV":        "activities.download-csv",
	"downloadActivityFile":           "activity.file",
	"downloadActivityFitFile":        "activity.fit-file",
	"downloadActivityFitFiles":       "athlete.download-fit-files",
	"downloadActivityGpxFile":        "activity.gpx-file",
	"downloadWorkout":                "misc.download-workout",
	"downloadWorkoutForAthlete":      "athlete.download-workout",
	"downloadWorkout_1":              "events.download-workout",
	"downloadWorkouts":               "workouts.download",
	"duplicateWorkouts":              "workouts.duplicate",
	"findBestEfforts":                "activity.best-efforts",
	"getActivities":                  "activities.get-multiple",
	"getActivity":                    "activity.get",
	"getActivityHRCurve":             "activity.hr-curve",
	"getActivityMap":                 "activity.map",
	"getActivityPaceCurve":           "activity.pace-curve",
	"getActivityPowerCurve":          "activity.power-curve",
	"getActivityPowerSpikeModel":     "activity.power-spike-model",
	"getActivitySegments":            "activity.segments",
	"getActivityStreams":             "activity.streams",
	"getActivityWeatherSummary":      "activity.weather-summary",
	"getAthlete":                     "athlete.get",
	"getAthleteMMPModel":             "athlete.mmp-model",
	"getAthleteProfile":              "athlete.profile",
	"getAthleteRoute":                "athlete.route-get",
	"getAthleteSummary":              "athlete.summary",
	"getAthleteTrainingPlan":         "athlete.training-plan",
	"getCustomItem":                  "custom-items.get",
	"getForecast":                    "athlete.weather-forecast",
	"getGapHistogram":                "activity.gap-histogram",
	"getHRHistogram":                 "activity.hr-histogram",
	"getHRTrainingLoadModel":         "activity.hr-load-model",
	"getIntervalStats":               "activity.interval-stats",
	"getIntervals":                   "activity.intervals",
	"getPaceHistogram":               "activity.pace-histogram",
	"getPowerHRCurve":                "athlete.power-hr-curve",
	"getPowerHistogram":              "activity.power-histogram",
	"getPowerVsHR":                   "activity.power-vs-hr",
	"getRecord":                      "wellness.get",
	"getSettings":                    "sport-settings.get-device",
	"getSettings_1":                  "sport-settings.get",
	"getSharedEvent":                 "shared-events.get-by-slug",
	"getSharedEvent_1":               "shared-events.get",
	"getTimeAtHR":                    "activity.time-at-hr",
	"getWeatherConfig":               "athlete.weather-config",
	"importWorkoutFile":              "folders.import-workout",
	"listActivities":                 "activities.list",
	"listActivitiesAround":           "activities.list-around",
	"listActivityHRCurves":           "athlete.activity-hr-curves",
	"listActivityMessages":           "activity.messages",
	"listActivityPaceCurves":         "athlete.activity-pace-curves",
	"listActivityPowerCurves":        "athlete.activity-power-curves",
	"listActivityPowerCurves_1":      "activity.power-curves",
	"listAthleteHRCurves":            "athlete.hr-curves",
	"listAthletePaceCurves":          "athlete.pace-curves",
	"listAthletePowerCurves":         "athlete.power-curves",
	"listAthleteRoutes":              "athlete.routes",
	"listChats":                      "athlete.chats",
	"listCustomItems":                "custom-items.list",
	"listEvents":                     "events.list",
	"listFitnessModelEvents":         "athlete.fitness-model-events",
	"listFolderSharedWith":           "folders.shared-with",
	"listFolders":                    "folders.list",
	"listGear":                       "gear.list",
	"listMatchingActivities":         "sport-settings.matching-activities",
	"listMessages":                   "chats.messages",
	"listPaceDistances":              "misc.pace-distances",
	"listPaceDistancesForSport":      "sport-settings.pace-distances",
	"listSettings":                   "sport-settings.list",
	"listTags":                       "workouts.tags",
	"listTags_1":                     "events.tags",
	"listTags_2":                     "athlete.tags",
	"listWellnessRecords":            "wellness.list",
	"listWorkouts":                   "workouts.list",
	"markEventAsDone":                "events.mark-done",
	"replaceGear":                    "gear.replace",
	"searchForActivities":            "activities.search",
	"searchForActivitiesFull":        "activities.search-full",
	"searchForIntervals":             "activities.interval-search",
	"sendActivityMessage":            "activity.send-message",
	"sendMessage":                    "chats.send",
	"showChat":                       "chats.get",
	"showEvent":                      "events.get",
	"showWorkout":                    "workouts.get",
	"splitInterval":                  "activity.split-interval",
	"updateActivity":                 "activity.update",
	"updateActivityStreams":          "activity.update-streams",
	"updateAthlete":                  "athlete.update",
	"updateAthletePlan":              "athlete.update-training-plan",
	"updateAthletePlans":             "misc.update-athlete-plans",
	"updateAthleteRoute":             "athlete.route-update",
	"updateCustomItem":               "custom-items.update",
	"updateCustomItemImage":          "custom-items.upload-image",
	"updateCustomItemIndexes":        "custom-items.update-indexes",
	"updateEvent":                    "events.update",
	"updateEvents":                   "events.update-bulk",
	"updateFolder":                   "folders.update",
	"updateFolderSharedWith":         "folders.update-shared-with",
	"updateGear":                     "gear.update",
	"updateInterval":                 "activity.update-interval",
	"updateIntervals":                "activity.update-intervals",
	"updateLastSeenMessageId":        "chats.mark-seen",
	"updatePlanWorkouts":             "folders.update-workouts",
	"updateReminder":                 "gear.update-reminder",
	"updateSettings":                 "sport-settings.update",
	"updateSettingsMulti":            "sport-settings.update-multi",
	"updateSharedEvent":              "shared-events.update",
	"updateSharedEventImage":         "shared-events.upload-image",
	"updateWeatherConfig":            "athlete.update-weather-config",
	"updateWellness":                 "wellness.update",
	"updateWellnessBulk":             "wellness.update-bulk",
	"updateWellness_1":               "wellness.update-current",
	"updateWorkout":                  "workouts.update",
	"uploadActivity":                 "activities.upload",
	"uploadActivityStreamsCSV":       "activity.upload-streams-csv",
	"uploadWellness":                 "wellness.upload",
}

var schemaCmd = &cobra.Command{
	Use:   "schema [operation|type]",
	Short: "Inspect API operation schemas or type definitions",
	Long: `Inspect the intervals.icu API schema.

  intervals schema activities.list       # Show operation schema
  intervals schema Activity              # Show type definition
  intervals schema --list                # List all operations`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		listFlag, _ := cmd.Flags().GetBool("list")
		resolveRefs, _ := cmd.Flags().GetBool("resolve-refs")

		spec, err := parseSpec()
		if err != nil {
			return err
		}

		if listFlag {
			return listOperations(spec)
		}

		if len(args) == 0 {
			return discoveryErr(fmt.Errorf("provide an operation path (e.g., activities.list) or type name (e.g., Activity)\n  Use --list to see all operations"))
		}

		path := args[0]
		if !strings.Contains(path, ".") {
			return showType(spec, path, resolveRefs)
		}
		return showOperation(spec, path, resolveRefs)
	},
}

func init() {
	schemaCmd.Flags().Bool("list", false, "List all available operations")
	schemaCmd.Flags().Bool("resolve-refs", false, "Inline $ref references in output")
	rootCmd.AddCommand(schemaCmd)
}

func parseSpec() (map[string]interface{}, error) {
	var spec map[string]interface{}
	if err := json.Unmarshal(specData, &spec); err != nil {
		return nil, discoveryErr(fmt.Errorf("parsing embedded OpenAPI spec: %w", err))
	}
	return spec, nil
}

func parseSchemaPath(path string) (string, string, error) {
	parts := strings.SplitN(path, ".", 2)
	if len(parts) != 2 {
		return "", "", discoveryErr(fmt.Errorf("schema path must be 'resource.method' (e.g., activities.list), got %q", path))
	}
	return parts[0], parts[1], nil
}


func buildOperationIndex(spec map[string]interface{}) map[string]operationEntry {
	index := make(map[string]operationEntry)
	paths, ok := spec["paths"].(map[string]interface{})
	if !ok {
		return index
	}
	for apiPath, methods := range paths {
		methodMap, ok := methods.(map[string]interface{})
		if !ok {
			continue
		}
		for method, opRaw := range methodMap {
			if method != "get" && method != "post" && method != "put" && method != "delete" {
				continue
			}
			op, ok := opRaw.(map[string]interface{})
			if !ok {
				continue
			}
			operationID, _ := op["operationId"].(string)
			cmdName, ok := operationIDToCommand[operationID]
			if !ok {
				continue
			}
			index[cmdName] = operationEntry{
				CommandName: cmdName,
				OperationID: operationID,
				Method:      strings.ToUpper(method),
				Path:        apiPath,
				PathItem:    op,
			}
		}
	}
	return index
}

func buildSchemaOutput(spec map[string]interface{}, entry operationEntry) map[string]interface{} {
	output := map[string]interface{}{
		"command":     entry.CommandName,
		"operationId": entry.OperationID,
		"method":      entry.Method,
		"path":        entry.Path,
	}
	if desc, ok := entry.PathItem["description"].(string); ok {
		output["description"] = desc
	}
	if summary, ok := entry.PathItem["summary"].(string); ok {
		output["summary"] = summary
	}
	if params, ok := entry.PathItem["parameters"].([]interface{}); ok {
		paramList := []map[string]interface{}{}
		for _, p := range params {
			param, ok := p.(map[string]interface{})
			if !ok {
				continue
			}
			paramEntry := map[string]interface{}{
				"name": param["name"], "in": param["in"], "required": param["required"],
			}
			if schema, ok := param["schema"].(map[string]interface{}); ok {
				paramEntry["type"] = schema["type"]
				if format, ok := schema["format"]; ok {
					paramEntry["format"] = format
				}
			}
			if desc, ok := param["description"].(string); ok {
				paramEntry["description"] = desc
			}
			paramList = append(paramList, paramEntry)
		}
		output["parameters"] = paramList
	}
	if reqBody, ok := entry.PathItem["requestBody"].(map[string]interface{}); ok {
		if content, ok := reqBody["content"].(map[string]interface{}); ok {
			for contentType, schemaRaw := range content {
				schema, ok := schemaRaw.(map[string]interface{})
				if !ok {
					continue
				}
				bodyInfo := map[string]interface{}{"contentType": contentType}
				if s, ok := schema["schema"].(map[string]interface{}); ok {
					bodyInfo["schema"] = s
				}
				output["requestBody"] = bodyInfo
				break
			}
		}
	}
	if responses, ok := entry.PathItem["responses"].(map[string]interface{}); ok {
		respInfo := map[string]interface{}{}
		for code, respRaw := range responses {
			resp, ok := respRaw.(map[string]interface{})
			if !ok {
				continue
			}
			entry := map[string]interface{}{}
			if desc, ok := resp["description"].(string); ok {
				entry["description"] = desc
			}
			if content, ok := resp["content"].(map[string]interface{}); ok {
				for ct, schemaRaw := range content {
					s, ok := schemaRaw.(map[string]interface{})
					if !ok {
						continue
					}
					entry["contentType"] = ct
					if schema, ok := s["schema"].(map[string]interface{}); ok {
						entry["schema"] = schema
					}
					break
				}
			}
			respInfo[code] = entry
		}
		output["responses"] = respInfo
	}
	return output
}

func showOperation(spec map[string]interface{}, path string, resolveRefs bool) error {
	index := buildOperationIndex(spec)
	entry, ok := index[path]
	if !ok {
		var suggestions []string
		resource, _, _ := parseSchemaPath(path)
		for name := range index {
			if strings.HasPrefix(name, resource+".") {
				suggestions = append(suggestions, name)
			}
		}
		sort.Strings(suggestions)
		msg := fmt.Sprintf("operation %q not found", path)
		if len(suggestions) > 0 {
			msg += fmt.Sprintf("\n  Available %s operations: %s", resource, strings.Join(suggestions, ", "))
		}
		return discoveryErr(fmt.Errorf("%s", msg))
	}
	output := buildSchemaOutput(spec, entry)
	if resolveRefs {
		components, _ := spec["components"].(map[string]interface{})
		schemas, _ := components["schemas"].(map[string]interface{})
		if schemas != nil {
			resolveAllRefs(output, schemas, map[string]bool{})
		}
	}
	data, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}

func showType(spec map[string]interface{}, name string, resolveRefs bool) error {
	components, ok := spec["components"].(map[string]interface{})
	if !ok {
		return discoveryErr(fmt.Errorf("no components in spec"))
	}
	schemas, ok := components["schemas"].(map[string]interface{})
	if !ok {
		return discoveryErr(fmt.Errorf("no schemas in spec"))
	}
	schema, ok := schemas[name]
	if !ok {
		var available []string
		for k := range schemas {
			available = append(available, k)
		}
		sort.Strings(available)
		return discoveryErr(fmt.Errorf("type %q not found\n  Available types: %s", name, strings.Join(available, ", ")))
	}
	output, ok := schema.(map[string]interface{})
	if !ok {
		return discoveryErr(fmt.Errorf("invalid schema for type %q", name))
	}
	if resolveRefs {
		resolveAllRefs(output, schemas, map[string]bool{name: true})
	}
	data, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}

func listOperations(spec map[string]interface{}) error {
	index := buildOperationIndex(spec)
	var names []string
	for name := range index {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		entry := index[name]
		fmt.Printf("%-40s  %s %s\n", name, entry.Method, entry.Path)
	}
	return nil
}

func resolveAllRefs(v interface{}, schemas map[string]interface{}, seen map[string]bool) {
	switch val := v.(type) {
	case map[string]interface{}:
		if ref, ok := val["$ref"].(string); ok {
			refName := strings.TrimPrefix(ref, "#/components/schemas/")
			if seen[refName] {
				return // cycle: leave $ref to terminate expansion
			}
			if schema, ok := schemas[refName].(map[string]interface{}); ok {
				seen[refName] = true
				delete(val, "$ref")
				for k, sv := range schema {
					if k != "$ref" {
						val[k] = deepCopy(sv)
					}
				}
				resolveAllRefs(val, schemas, seen)
				delete(seen, refName)
			}
			return
		}
		for _, child := range val {
			resolveAllRefs(child, schemas, seen)
		}
	case []interface{}:
		for _, item := range val {
			resolveAllRefs(item, schemas, seen)
		}
	}
}

func deepCopy(v interface{}) interface{} {
	switch val := v.(type) {
	case map[string]interface{}:
		cp := make(map[string]interface{}, len(val))
		for k, v := range val {
			cp[k] = deepCopy(v)
		}
		return cp
	case []interface{}:
		cp := make([]interface{}, len(val))
		for i, v := range val {
			cp[i] = deepCopy(v)
		}
		return cp
	default:
		return v
	}
}

// internal/cmd/schema.go
package cmd

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

//go:embed openapi-spec.json
var specData []byte

type operationEntry struct {
	CommandName string
	OperationID string
	Method      string
	Path        string
	PathItem    map[string]interface{}
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
			return fmt.Errorf("provide an operation path (e.g., activities.list) or type name (e.g., Activity)\n  Use --list to see all operations")
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
		return nil, fmt.Errorf("parsing embedded OpenAPI spec: %w", err)
	}
	return spec, nil
}

func parseSchemaPath(path string) (string, string, error) {
	parts := strings.SplitN(path, ".", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("schema path must be 'resource.method' (e.g., activities.list), got %q", path)
	}
	return parts[0], parts[1], nil
}

func commandNameForOperation(httpMethod, apiPath, operationID string) string {
	clean := strings.TrimPrefix(apiPath, "/api/v1/")

	if strings.HasPrefix(clean, "activity/") {
		return activityCommandName(httpMethod, clean, operationID)
	}
	if strings.HasPrefix(clean, "athlete/") {
		return athleteCommandName(httpMethod, clean, operationID)
	}
	if strings.HasPrefix(clean, "chats/") {
		return chatsCommandName(httpMethod, clean, operationID)
	}
	if strings.HasPrefix(clean, "shared-event") {
		return sharedEventsCommandName(httpMethod, clean, operationID)
	}
	return miscCommandName(httpMethod, clean, operationID)
}

func activityCommandName(httpMethod, clean, operationID string) string {
	parts := strings.Split(clean, "/")
	if len(parts) <= 2 {
		switch httpMethod {
		case "get":
			return "activity.get"
		case "put":
			return "activity.update"
		case "delete":
			return "activity.delete"
		}
	}
	sub := parts[2]
	sub = strings.TrimSuffix(sub, "{ext}")

	if len(parts) > 3 {
		switch {
		case sub == "intervals" && httpMethod == "put" && strings.Contains(clean, "{intervalId}"):
			return "activity.update-interval"
		case sub == "delete-intervals":
			return "activity.delete-intervals"
		}
	}

	switch {
	case sub == "streams.csv":
		return "activity.upload-streams-csv"
	case sub == "streams" && httpMethod == "put":
		return "activity.update-streams"
	case sub == "streams":
		return "activity.streams"
	case sub == "split-interval":
		return "activity.split-interval"
	case sub == "intervals" && httpMethod == "get":
		return "activity.intervals"
	case sub == "intervals" && httpMethod == "put":
		return "activity.update-intervals"
	case sub == "messages" && httpMethod == "get":
		return "activity.messages"
	case sub == "messages" && httpMethod == "post":
		return "activity.send-message"
	default:
		return "activity." + sub
	}
}

func athleteCommandName(httpMethod, clean, operationID string) string {
	parts := strings.Split(clean, "/")
	if len(parts) <= 2 {
		switch httpMethod {
		case "get":
			return "athlete.get"
		case "put":
			return "athlete.update"
		}
	}
	sub := parts[2]

	switch {
	case sub == "activities" || sub == "activities.csv" || strings.HasPrefix(sub, "activities-"):
		return activitiesCommandName(httpMethod, clean, operationID)
	case sub == "wellness" || sub == "wellness-bulk" || strings.HasSuffix(sub, "wellness"):
		return wellnessCommandName(httpMethod, clean, operationID)
	case sub == "events" || strings.HasPrefix(sub, "events"):
		return eventsCommandName(httpMethod, clean, operationID)
	case sub == "event-tags":
		return "tags.event"
	case sub == "workouts" || sub == "workouts.zip" || sub == "workout-tags" || sub == "duplicate-workouts":
		return workoutsCommandName(httpMethod, clean, operationID)
	case sub == "gear" || strings.HasSuffix(sub, "gear"):
		return gearCommandName(httpMethod, clean, operationID)
	case sub == "folders":
		return foldersCommandName(httpMethod, clean, operationID)
	case sub == "sport-settings" || sub == "settings":
		return sportSettingsCommandName(httpMethod, clean, operationID)
	case sub == "custom-item" || sub == "custom-item-indexes":
		return customItemsCommandName(httpMethod, clean, operationID)
	case sub == "chats":
		return "athlete.chats"
	case sub == "duplicate-events":
		return "events.duplicate"
	}

	subClean := strings.TrimSuffix(sub, "{ext}")
	switch {
	case sub == "profile":
		return "athlete.profile"
	case sub == "training-plan" && httpMethod == "get":
		return "athlete.training-plan"
	case sub == "training-plan" && httpMethod == "put":
		return "athlete.update-training-plan"
	case sub == "weather-config" && httpMethod == "get":
		return "athlete.weather-config"
	case sub == "weather-config" && httpMethod == "put":
		return "athlete.update-weather-config"
	case sub == "weather-forecast":
		return "athlete.weather-forecast"
	case sub == "apply-plan-changes":
		return "athlete.apply-plan-changes"
	case sub == "routes":
		return athleteRoutesCommandName(httpMethod, clean, operationID)
	case sub == "download-fit-files":
		return "athlete.download-fit-files"
	case strings.HasPrefix(sub, "download-workout"):
		return "athlete.download-workout"
	case sub == "fitness-model-events":
		return "athlete.fitness-model-events"
	case sub == "activity-tags":
		return "tags.activity"
	default:
		return "athlete." + subClean
	}
}

func activitiesCommandName(httpMethod, clean, operationID string) string {
	switch operationID {
	case "listActivities":
		return "activities.list"
	case "uploadActivity":
		return "activities.upload"
	case "createManualActivity":
		return "activities.create-manual"
	case "createMultipleManualActivities":
		return "activities.create-manual-bulk"
	case "searchForActivities":
		return "activities.search"
	case "searchForActivitiesFull":
		return "activities.search-full"
	case "searchForIntervals":
		return "activities.interval-search"
	case "downloadActivitiesAsCSV":
		return "activities.download-csv"
	case "listActivitiesAround":
		return "activities.list-around"
	case "getActivities":
		return "activities.get-multiple"
	default:
		return "activities." + operationID
	}
}

func wellnessCommandName(httpMethod, clean, operationID string) string {
	switch operationID {
	case "getRecord":
		return "wellness.get"
	case "updateWellness":
		return "wellness.update"
	case "updateWellness_1":
		return "wellness.update-current"
	case "updateWellnessBulk":
		return "wellness.update-bulk"
	case "uploadWellness":
		return "wellness.upload"
	case "listWellnessRecords":
		return "wellness.list"
	default:
		return "wellness." + operationID
	}
}

func eventsCommandName(httpMethod, clean, operationID string) string {
	switch operationID {
	case "listEvents":
		return "events.list"
	case "createEvent":
		return "events.create"
	case "showEvent":
		return "events.get"
	case "updateEvent":
		return "events.update"
	case "deleteEvent":
		return "events.delete"
	case "updateEvents":
		return "events.update-bulk"
	case "deleteEvents":
		return "events.delete-range"
	case "deleteEventsBulk":
		return "events.delete-bulk"
	case "createMultipleEvents":
		return "events.create-bulk"
	case "markEventAsDone":
		return "events.mark-done"
	case "applyPlan":
		return "events.apply-plan"
	case "downloadWorkout_1":
		return "events.download-workout"
	case "listTags_1":
		return "tags.event"
	default:
		return "events." + operationID
	}
}

func workoutsCommandName(httpMethod, clean, operationID string) string {
	switch operationID {
	case "listWorkouts":
		return "workouts.list"
	case "showWorkout":
		return "workouts.get"
	case "createWorkout":
		return "workouts.create"
	case "updateWorkout":
		return "workouts.update"
	case "deleteWorkout":
		return "workouts.delete"
	case "createMultipleWorkouts":
		return "workouts.create-bulk"
	case "duplicateWorkouts":
		return "workouts.duplicate"
	case "downloadWorkouts":
		return "workouts.download"
	case "listTags":
		return "tags.workout"
	default:
		return "workouts." + operationID
	}
}

func gearCommandName(httpMethod, clean, operationID string) string {
	switch operationID {
	case "listGear":
		return "gear.list"
	case "createGear":
		return "gear.create"
	case "updateGear":
		return "gear.update"
	case "deleteGear":
		return "gear.delete"
	case "replaceGear":
		return "gear.replace"
	case "calcDistanceEtc":
		return "gear.calc"
	case "createReminder":
		return "gear.create-reminder"
	case "updateReminder":
		return "gear.update-reminder"
	case "deleteReminder":
		return "gear.delete-reminder"
	default:
		return "gear." + operationID
	}
}

func foldersCommandName(httpMethod, clean, operationID string) string {
	switch operationID {
	case "listFolders":
		return "folders.list"
	case "createFolder":
		return "folders.create"
	case "updateFolder":
		return "folders.update"
	case "deleteFolder":
		return "folders.delete"
	case "updatePlanWorkouts":
		return "folders.update-workouts"
	case "listFolderSharedWith":
		return "folders.shared-with"
	case "updateFolderSharedWith":
		return "folders.update-shared-with"
	case "importWorkoutFile":
		return "folders.import-workout"
	default:
		return "folders." + operationID
	}
}

func sportSettingsCommandName(httpMethod, clean, operationID string) string {
	switch operationID {
	case "listSettings":
		return "sport-settings.list"
	case "getSettings_1":
		return "sport-settings.get"
	case "getSettings":
		return "sport-settings.get-device"
	case "createSettings":
		return "sport-settings.create"
	case "updateSettings":
		return "sport-settings.update"
	case "updateSettingsMulti":
		return "sport-settings.update-multi"
	case "deleteSettings":
		return "sport-settings.delete"
	case "applyToActivities":
		return "sport-settings.apply"
	case "listPaceDistancesForSport":
		return "sport-settings.pace-distances"
	case "listMatchingActivities":
		return "sport-settings.matching-activities"
	default:
		return "sport-settings." + operationID
	}
}

func customItemsCommandName(httpMethod, clean, operationID string) string {
	switch operationID {
	case "listCustomItems":
		return "custom-items.list"
	case "getCustomItem":
		return "custom-items.get"
	case "createCustomItem":
		return "custom-items.create"
	case "updateCustomItem":
		return "custom-items.update"
	case "deleteCustomItem":
		return "custom-items.delete"
	case "updateCustomItemImage":
		return "custom-items.upload-image"
	case "updateCustomItemIndexes":
		return "custom-items.update-indexes"
	default:
		return "custom-items." + operationID
	}
}

func athleteRoutesCommandName(httpMethod, clean, operationID string) string {
	switch operationID {
	case "listAthleteRoutes":
		return "athlete.routes"
	case "getAthleteRoute":
		return "athlete.route-get"
	case "updateAthleteRoute":
		return "athlete.route-update"
	case "checkMerge":
		return "athlete.route-similarity"
	default:
		return "athlete." + operationID
	}
}

func chatsCommandName(httpMethod, clean, operationID string) string {
	switch operationID {
	case "showChat":
		return "chats.get"
	case "listMessages":
		return "chats.messages"
	case "sendMessage":
		return "chats.send"
	case "deleteMessage":
		return "chats.delete-message"
	case "updateLastSeenMessageId":
		return "chats.mark-seen"
	default:
		return "chats." + operationID
	}
}

func sharedEventsCommandName(httpMethod, clean, operationID string) string {
	switch operationID {
	case "getSharedEvent":
		return "shared-events.get-by-slug"
	case "getSharedEvent_1":
		return "shared-events.get"
	case "createSharedEvent":
		return "shared-events.create"
	case "updateSharedEvent":
		return "shared-events.update"
	case "deleteSharedEvent":
		return "shared-events.delete"
	case "updateSharedEventImage":
		return "shared-events.upload-image"
	default:
		return "shared-events." + operationID
	}
}

func miscCommandName(httpMethod, clean, operationID string) string {
	switch operationID {
	case "updateAthletePlans":
		return "misc.update-athlete-plans"
	case "downloadWorkout":
		return "misc.download-workout"
	case "listPaceDistances":
		return "misc.pace-distances"
	case "disconnectApp":
		return "misc.disconnect-app"
	default:
		return "misc." + operationID
	}
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
			cmdName := commandNameForOperation(method, apiPath, operationID)
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
		return fmt.Errorf("%s", msg)
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
		return fmt.Errorf("no components in spec")
	}
	schemas, ok := components["schemas"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("no schemas in spec")
	}
	schema, ok := schemas[name]
	if !ok {
		var available []string
		for k := range schemas {
			available = append(available, k)
		}
		sort.Strings(available)
		return fmt.Errorf("type %q not found\n  Available types: %s", name, strings.Join(available, ", "))
	}
	output, ok := schema.(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid schema for type %q", name)
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
			if !seen[refName] {
				if schema, ok := schemas[refName].(map[string]interface{}); ok {
					seen[refName] = true
					for k, sv := range schema {
						if k != "$ref" {
							val[k] = deepCopy(sv)
						}
					}
					resolveAllRefs(val, schemas, seen)
					delete(seen, refName)
				}
			}
		}
		for _, v := range val {
			resolveAllRefs(v, schemas, seen)
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

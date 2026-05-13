package cmd

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

// TestAllOperationIdsMapped: every operationId in the spec must appear in
// operationIDToCommand. Catches drift — a new endpoint added to the spec
// without a CLI mapping.
func TestAllOperationIdsMapped(t *testing.T) {
	spec := loadEmbeddedSpec(t)
	paths, _ := spec["paths"].(map[string]interface{})
	var missing []string
	for _, item := range paths {
		pi, _ := item.(map[string]interface{})
		for _, m := range []string{"get", "post", "put", "delete", "patch"} {
			op, ok := pi[m].(map[string]interface{})
			if !ok {
				continue
			}
			opID, _ := op["operationId"].(string)
			if opID == "" {
				continue
			}
			if _, ok := operationIDToCommand[opID]; !ok {
				missing = append(missing, opID)
			}
		}
	}
	if len(missing) > 0 {
		t.Fatalf("%d operationIds missing CLI mapping: %v", len(missing), missing)
	}
}

// TestEveryMappedCLINameExists: every value in operationIDToCommand must
// resolve to a real cobra leaf command.
func TestEveryMappedCLINameExists(t *testing.T) {
	var missing []string
	for opID, cliName := range operationIDToCommand {
		parts := strings.SplitN(cliName, ".", 2)
		if len(parts) != 2 {
			missing = append(missing, cliName+" (malformed name)")
			continue
		}
		var group *cobra.Command
		for _, sub := range rootCmd.Commands() {
			if sub.Name() == parts[0] {
				group = sub
				break
			}
		}
		if group == nil {
			missing = append(missing, cliName+" (group not found, from "+opID+")")
			continue
		}
		found := false
		for _, sub := range group.Commands() {
			if sub.Name() == parts[1] {
				found = true
				break
			}
		}
		if !found {
			missing = append(missing, cliName+" (from "+opID+")")
		}
	}
	if len(missing) > 0 {
		t.Fatalf("%d mapped CLI names not registered: %v", len(missing), missing)
	}
}

func TestParseOperationPath(t *testing.T) {
	tests := []struct {
		input    string
		resource string
		method   string
	}{
		{"activities.list", "activities", "list"},
		{"activity.power-curve", "activity", "power-curve"},
		{"events.create", "events", "create"},
	}
	for _, tt := range tests {
		resource, method, err := parseSchemaPath(tt.input)
		if err != nil {
			t.Errorf("parseSchemaPath(%q) error: %v", tt.input, err)
			continue
		}
		if resource != tt.resource {
			t.Errorf("parseSchemaPath(%q) resource = %q, want %q", tt.input, resource, tt.resource)
		}
		if method != tt.method {
			t.Errorf("parseSchemaPath(%q) method = %q, want %q", tt.input, method, tt.method)
		}
	}
}

func TestParseSchemaPathSinglePart(t *testing.T) {
	_, _, err := parseSchemaPath("Activity")
	if err == nil {
		t.Error("expected error for single-part path")
	}
}

func TestBuildOperationIndex(t *testing.T) {
	spec := loadEmbeddedSpec(t)
	index := buildOperationIndex(spec)
	if _, ok := index["activities.list"]; !ok {
		t.Error("missing activities.list from index")
	}
	if _, ok := index["activity.get"]; !ok {
		t.Error("missing activity.get from index")
	}
	if len(index) < 100 {
		t.Errorf("index has %d entries, expected at least 100", len(index))
	}
}

func TestSchemaOutputIsValidJSON(t *testing.T) {
	spec := loadEmbeddedSpec(t)
	index := buildOperationIndex(spec)
	entry, ok := index["activities.list"]
	if !ok {
		t.Fatal("missing activities.list")
	}
	output := buildSchemaOutput(spec, entry)
	data, err := json.Marshal(output)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	if len(data) < 10 {
		t.Error("schema output too small")
	}
}

func loadEmbeddedSpec(t *testing.T) map[string]interface{} {
	t.Helper()
	var spec map[string]interface{}
	if err := json.Unmarshal(specData, &spec); err != nil {
		t.Fatalf("failed to parse embedded spec: %v", err)
	}
	return spec
}

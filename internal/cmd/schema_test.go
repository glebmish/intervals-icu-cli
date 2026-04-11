package cmd

import (
	"encoding/json"
	"testing"
)

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

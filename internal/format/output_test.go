package format_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/glebmish/intervals-icu-cli/internal/format"
)

func TestJSONFormat(t *testing.T) {
	input := []byte(`{"id":1,"name":"test"}`)
	var buf bytes.Buffer
	opts := format.Options{Format: "json"}
	if err := format.Write(&buf, input, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if len(out) == 0 {
		t.Fatal("expected non-empty output")
	}
	// should be valid JSON
	var v interface{}
	if err := json.Unmarshal([]byte(out), &v); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	// should be pretty-printed (contains newlines/indentation)
	if !strings.Contains(out, "\n") {
		t.Error("expected pretty-printed JSON with newlines")
	}
}

func TestJSONFieldsFilter(t *testing.T) {
	input := []byte(`{"id":1,"name":"test","extra":"should-be-removed"}`)
	var buf bytes.Buffer
	opts := format.Options{Format: "json", Fields: []string{"id", "name"}}
	if err := format.Write(&buf, input, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if strings.Contains(out, "extra") {
		t.Error("expected 'extra' field to be filtered out")
	}
	if !strings.Contains(out, "id") {
		t.Error("expected 'id' field to be present")
	}
	if !strings.Contains(out, "name") {
		t.Error("expected 'name' field to be present")
	}
}

func TestNDJSONFormatArray(t *testing.T) {
	input := []byte(`[{"id":1},{"id":2},{"id":3}]`)
	var buf bytes.Buffer
	opts := format.Options{Format: "ndjson"}
	if err := format.Write(&buf, input, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
	if len(lines) != 3 {
		t.Errorf("expected 3 lines, got %d: %q", len(lines), out)
	}
	for i, line := range lines {
		var v interface{}
		if err := json.Unmarshal([]byte(line), &v); err != nil {
			t.Errorf("line %d is not valid JSON: %v", i, err)
		}
	}
}

func TestNDJSONFormatObject(t *testing.T) {
	input := []byte(`{"id":1,"name":"test"}`)
	var buf bytes.Buffer
	opts := format.Options{Format: "ndjson"}
	if err := format.Write(&buf, input, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
	if len(lines) != 1 {
		t.Errorf("expected 1 line, got %d: %q", len(lines), out)
	}
	var v interface{}
	if err := json.Unmarshal([]byte(lines[0]), &v); err != nil {
		t.Errorf("output is not valid JSON: %v", err)
	}
}

func TestNDJSONWithFieldsFilter(t *testing.T) {
	input := []byte(`[{"id":1,"extra":"x"},{"id":2,"extra":"y"}]`)
	var buf bytes.Buffer
	opts := format.Options{Format: "ndjson", Fields: []string{"id"}}
	if err := format.Write(&buf, input, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if strings.Contains(out, "extra") {
		t.Error("expected 'extra' field to be filtered out")
	}
	lines := strings.Split(strings.TrimRight(out, "\n"), "\n")
	if len(lines) != 2 {
		t.Errorf("expected 2 lines, got %d", len(lines))
	}
}

func TestFieldsFilterOnArray(t *testing.T) {
	input := []byte(`[{"id":1,"name":"a","extra":"x"},{"id":2,"name":"b","extra":"y"}]`)
	var buf bytes.Buffer
	opts := format.Options{Format: "json", Fields: []string{"id", "name"}}
	if err := format.Write(&buf, input, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if strings.Contains(out, "extra") {
		t.Error("expected 'extra' field to be filtered out")
	}
	if !strings.Contains(out, "id") {
		t.Error("expected 'id' field to be present")
	}
}

func TestTextFormatFallsBackToJSON(t *testing.T) {
	input := []byte(`{"id":1,"name":"test"}`)
	var buf bytes.Buffer
	opts := format.Options{Format: "text"}
	if err := format.Write(&buf, input, opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	var v interface{}
	if err := json.Unmarshal([]byte(out), &v); err != nil {
		t.Fatalf("text format output is not valid JSON: %v", err)
	}
	if !strings.Contains(out, "\n") {
		t.Error("expected pretty-printed JSON with newlines for text fallback")
	}
}

func TestRawBinaryOutput(t *testing.T) {
	data := []byte{0x01, 0x02, 0x03, 0xFF, 0xFE}
	var buf bytes.Buffer
	if err := format.WriteRaw(&buf, data); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !bytes.Equal(buf.Bytes(), data) {
		t.Errorf("expected raw bytes %v, got %v", data, buf.Bytes())
	}
}

package format

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// Options controls how output is formatted.
type Options struct {
	Format string
	Fields []string
}

// FormatFromFlag parses CLI flags into Options.
func FormatFromFlag(formatFlag, fieldsFlag string) Options {
	opts := Options{Format: formatFlag}
	if fieldsFlag != "" {
		for _, f := range strings.Split(fieldsFlag, ",") {
			f = strings.TrimSpace(f)
			if f != "" {
				opts.Fields = append(opts.Fields, f)
			}
		}
	}
	return opts
}

// Write parses data as JSON, applies field filtering, and writes formatted output to w.
func Write(w io.Writer, data []byte, opts Options) error {
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return fmt.Errorf("parse JSON: %w", err)
	}

	if len(opts.Fields) > 0 {
		v = filterFields(v, opts.Fields)
	}

	switch opts.Format {
	case "ndjson":
		return writeNDJSON(w, v)
	case "json", "text", "":
		return writePrettyJSON(w, v)
	default:
		return writePrettyJSON(w, v)
	}
}

// WriteRaw passes bytes through to w unchanged.
func WriteRaw(w io.Writer, data []byte) error {
	_, err := w.Write(data)
	return err
}

// DryRunOutput prints dry-run text to w.
func DryRunOutput(w io.Writer, dryRunText string) error {
	_, err := fmt.Fprintln(w, dryRunText)
	return err
}

func writePrettyJSON(w io.Writer, v interface{}) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	return enc.Encode(v)
}

func writeNDJSON(w io.Writer, v interface{}) error {
	arr, ok := v.([]interface{})
	if !ok {
		// single object — write one line
		return writeJSONLine(w, v)
	}
	for _, elem := range arr {
		if err := writeJSONLine(w, elem); err != nil {
			return err
		}
	}
	return nil
}

func writeJSONLine(w io.Writer, v interface{}) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	buf.Write(b)
	buf.WriteByte('\n')
	_, err = w.Write(buf.Bytes())
	return err
}

// filterFields recursively removes keys not in the fields list.
// Works on maps and slices.
func filterFields(v interface{}, fields []string) interface{} {
	allowed := make(map[string]struct{}, len(fields))
	for _, f := range fields {
		allowed[f] = struct{}{}
	}
	return filterValue(v, allowed)
}

func filterValue(v interface{}, allowed map[string]struct{}) interface{} {
	switch val := v.(type) {
	case map[string]interface{}:
		result := make(map[string]interface{})
		for k, child := range val {
			if _, ok := allowed[k]; ok {
				result[k] = child
			}
		}
		return result
	case []interface{}:
		result := make([]interface{}, len(val))
		for i, elem := range val {
			result[i] = filterValue(elem, allowed)
		}
		return result
	default:
		return v
	}
}

package format

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"regexp"
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

// Write parses data as JSON, sanitizes, applies field filtering, and writes
// formatted output to w. Non-JSON bodies (plain strings, raw integers, etc.)
// pass through unchanged so callers don't have to handle that case.
func Write(w io.Writer, data []byte, opts Options) error {
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		// Non-JSON bodies pass through.
		_, err := w.Write(append(data, '\n'))
		return err
	}

	v = sanitize(v)

	if len(opts.Fields) > 0 {
		v = filterFields(v, opts.Fields)
	}

	switch opts.Format {
	case "ndjson":
		return writeNDJSON(w, v)
	default:
		return writePrettyJSON(w, v)
	}
}

// injectionTagPattern matches XML-ish tag wrappers used in some prompt-injection
// attempts, e.g. <system>...</system>, <assistant>...</assistant>. We strip the
// wrapper and keep the inner text — opening/closing tags only, case-insensitive.
var injectionTagPattern = regexp.MustCompile(`(?i)</?(system|assistant|tool_use|tool_result)[^>]*>`)

// sanitize walks parsed JSON and cleans string values to defend against
// prompt injection embedded in API responses. Strips control characters and
// known injection tag wrappers. Defensive minimum per design-cli §13.
func sanitize(v interface{}) interface{} {
	switch val := v.(type) {
	case string:
		return sanitizeString(val)
	case map[string]interface{}:
		out := make(map[string]interface{}, len(val))
		for k, child := range val {
			out[k] = sanitize(child)
		}
		return out
	case []interface{}:
		out := make([]interface{}, len(val))
		for i, child := range val {
			out[i] = sanitize(child)
		}
		return out
	default:
		return v
	}
}

func sanitizeString(s string) string {
	if s == "" {
		return s
	}
	s = injectionTagPattern.ReplaceAllString(s, "")
	if !needsControlStrip(s) {
		return s
	}
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		if r < 0x20 && r != '\t' && r != '\n' && r != '\r' {
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}

func needsControlStrip(s string) bool {
	for _, r := range s {
		if r < 0x20 && r != '\t' && r != '\n' && r != '\r' {
			return true
		}
	}
	return false
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

// filterFields keeps only the requested paths. Dotted paths descend through
// objects and arrays: "rows.id" keeps {"rows":[{"id":1}]} → {"rows":[{"id":1}]}.
func filterFields(v interface{}, fields []string) interface{} {
	type node struct {
		leaf     bool
		children map[string]*node
	}
	root := &node{children: map[string]*node{}}
	for _, f := range fields {
		segs := strings.Split(f, ".")
		cur := root
		for _, s := range segs {
			child, ok := cur.children[s]
			if !ok {
				child = &node{children: map[string]*node{}}
				cur.children[s] = child
			}
			cur = child
		}
		cur.leaf = true
	}
	var walk func(v interface{}, n *node) interface{}
	walk = func(v interface{}, n *node) interface{} {
		if n == nil {
			return v
		}
		switch val := v.(type) {
		case map[string]interface{}:
			result := make(map[string]interface{})
			for k, child := range val {
				cn, ok := n.children[k]
				if !ok {
					continue
				}
				if cn.leaf && len(cn.children) == 0 {
					result[k] = child
				} else {
					result[k] = walk(child, cn)
				}
			}
			return result
		case []interface{}:
			out := make([]interface{}, len(val))
			for i, elem := range val {
				out[i] = walk(elem, n)
			}
			return out
		default:
			return v
		}
	}
	return walk(v, root)
}

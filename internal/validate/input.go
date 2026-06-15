package validate

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/glebmish/intervals-icu-cli/internal/cliexit"
)

func vErr(format string, args ...any) error {
	return &cliexit.ValidationError{Err: fmt.Errorf(format, args...)}
}

// PathParam validates that a user-supplied value is safe to embed in a URL path.
func PathParam(name, value string) error {
	if value == "" {
		return vErr("field %q: must not be empty", name)
	}
	if strings.Contains(value, "..") {
		return vErr("field %q: contains path traversal characters", name)
	}
	if strings.ContainsAny(value, "?#&") {
		return vErr("field %q: contains query injection characters", name)
	}
	if strings.Contains(value, "%") {
		return vErr("field %q: contains percent-encoded characters (provide raw values)", name)
	}
	for _, r := range value {
		if r < 0x20 {
			return vErr("field %q: contains control characters", name)
		}
	}
	return nil
}

// FilePath validates a user-supplied local file path before it is opened for
// upload. Treats the path as adversarial (agents hallucinate): rejects control
// characters and NUL bytes. Returns the cleaned path. Callers read the file;
// this does not enforce a CWD sandbox so legitimate absolute upload paths still
// work, but it blocks the obvious injection shapes.
func FilePath(name, value string) (string, error) {
	if value == "" {
		return "", vErr("field %q: must not be empty", name)
	}
	for _, r := range value {
		if r < 0x20 || r == 0x7f {
			return "", vErr("field %q: contains control characters", name)
		}
	}
	if strings.ContainsRune(value, 0) {
		return "", vErr("field %q: contains NUL byte", name)
	}
	return filepath.Clean(value), nil
}

// DateParam accepts ISO YYYY-MM-DD only.
func DateParam(name, value string) error {
	if _, err := time.Parse("2006-01-02", value); err != nil {
		return vErr("field %q: invalid date %q (expected YYYY-MM-DD)", name, value)
	}
	return nil
}

// JSONBody ensures the body is syntactically valid JSON with no rogue control chars.
func JSONBody(body string) error {
	if body == "" {
		return vErr("empty JSON body")
	}
	for i, r := range body {
		if r < 0x20 && r != '\t' && r != '\n' && r != '\r' {
			return vErr("JSON body contains control character at position %d", i)
		}
	}
	if !json.Valid([]byte(body)) {
		return vErr("invalid JSON body")
	}
	return nil
}

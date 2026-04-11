package validate

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

func PathParam(name, value string) error {
	if value == "" {
		return fmt.Errorf("field %q: must not be empty", name)
	}
	if strings.Contains(value, "..") {
		return fmt.Errorf("field %q: contains path traversal characters", name)
	}
	if strings.ContainsAny(value, "?#&") {
		return fmt.Errorf("field %q: contains query injection characters", name)
	}
	if strings.Contains(value, "%") {
		return fmt.Errorf("field %q: contains percent-encoded characters (provide raw values)", name)
	}
	for _, r := range value {
		if r < 0x20 {
			return fmt.Errorf("field %q: contains control characters", name)
		}
	}
	return nil
}

func DateParam(name, value string) error {
	_, err := time.Parse("2006-01-02", value)
	if err != nil {
		return fmt.Errorf("field %q: invalid date %q (expected YYYY-MM-DD)", name, value)
	}
	return nil
}

func JSONBody(body string) error {
	if body == "" {
		return fmt.Errorf("empty JSON body")
	}
	for i, r := range body {
		if r < 0x20 && r != '\t' && r != '\n' && r != '\r' {
			return fmt.Errorf("JSON body contains control character at position %d", i)
		}
	}
	if !json.Valid([]byte(body)) {
		return fmt.Errorf("invalid JSON body")
	}
	return nil
}

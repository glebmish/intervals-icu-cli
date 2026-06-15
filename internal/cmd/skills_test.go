package cmd

import (
	"strings"
	"testing"
)

func TestSplitFrontmatterClosingDelimAtEOF(t *testing.T) {
	// Closing --- at EOF with no trailing newline.
	content := []byte("---\nname: foo\ndescription: bar\n---")
	fm, body := splitFrontmatter(content)
	if len(fm) == 0 {
		t.Fatal("expected frontmatter, got none")
	}
	if !strings.Contains(string(fm), "name: foo") {
		t.Errorf("frontmatter = %q, want it to contain name: foo", fm)
	}
	if len(body) != 0 {
		t.Errorf("body = %q, want empty", body)
	}
}

func TestSplitFrontmatterClosingDelimWithNewline(t *testing.T) {
	content := []byte("---\nname: foo\n---\nbody text\n")
	fm, body := splitFrontmatter(content)
	if !strings.Contains(string(fm), "name: foo") {
		t.Errorf("frontmatter = %q, want name: foo", fm)
	}
	if strings.TrimSpace(string(body)) != "body text" {
		t.Errorf("body = %q, want 'body text'", body)
	}
}

func TestSplitFrontmatterNoFrontmatter(t *testing.T) {
	content := []byte("just body, no frontmatter")
	fm, body := splitFrontmatter(content)
	if fm != nil {
		t.Errorf("frontmatter = %q, want nil", fm)
	}
	if string(body) != string(content) {
		t.Errorf("body = %q, want full content", body)
	}
}

package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestMergeParamsKeepsIntegerLiteral(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String("params", `{"oldest":1700000000}`, "")

	merged, err := mergeParams(cmd, nil)
	if err != nil {
		t.Fatalf("mergeParams error = %v", err)
	}
	if got := merged["oldest"]; got != "1700000000" {
		t.Errorf("merged[oldest] = %q, want %q (no scientific notation)", got, "1700000000")
	}
}

func TestMergeParamsBaseWins(t *testing.T) {
	cmd := &cobra.Command{}
	cmd.Flags().String("params", `{"id":"fromparams","extra":"x"}`, "")

	merged, err := mergeParams(cmd, map[string]string{"id": "base"})
	if err != nil {
		t.Fatalf("mergeParams error = %v", err)
	}
	if merged["id"] != "base" {
		t.Errorf("merged[id] = %q, want base (caller-set entries win)", merged["id"])
	}
	if merged["extra"] != "x" {
		t.Errorf("merged[extra] = %q, want x", merged["extra"])
	}
}

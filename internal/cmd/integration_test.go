// internal/cmd/integration_test.go
package cmd

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func setupTestEnv(t *testing.T, handler http.HandlerFunc) (server *httptest.Server, cleanup func()) {
	t.Helper()
	server = httptest.NewServer(handler)

	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.yaml")
	cfg := []byte("api_key: test-key\nathlete_id: i123\nbase_url: " + server.URL + "\n")
	os.WriteFile(cfgPath, cfg, 0644)
	t.Setenv("INTERVALS_CONFIG", cfgPath)

	return server, func() {
		server.Close()
	}
}

func TestActivitiesListE2E(t *testing.T) {
	server, cleanup := setupTestEnv(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/athlete/i123/activities" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("oldest") != "2026-01-01" {
			t.Errorf("missing oldest param")
		}
		user, pass, _ := r.BasicAuth()
		if user != "API_KEY" || pass != "test-key" {
			t.Errorf("bad auth: %s/%s", user, pass)
		}
		w.WriteHeader(200)
		json.NewEncoder(w).Encode([]map[string]interface{}{
			{"id": "a1", "name": "Morning Run"},
		})
	})
	defer cleanup()
	_ = server

	rootCmd.SetArgs([]string{"activities", "list", "--oldest", "2026-01-01"})
	var out bytes.Buffer
	rootCmd.SetOut(&out)
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("execute failed: %v", err)
	}
}

func TestDryRunE2E(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.yaml")
	cfg := []byte("api_key: test-key\nathlete_id: i123\n")
	os.WriteFile(cfgPath, cfg, 0644)
	t.Setenv("INTERVALS_CONFIG", cfgPath)

	rootCmd.SetArgs([]string{"activities", "list", "--oldest", "2026-01-01", "--dry-run"})
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("dry-run failed: %v", err)
	}
}

func TestSchemaListE2E(t *testing.T) {
	rootCmd.SetArgs([]string{"schema", "--list"})
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("schema --list failed: %v", err)
	}
}

func TestSchemaOperationE2E(t *testing.T) {
	rootCmd.SetArgs([]string{"schema", "activities.list"})
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("schema activities.list failed: %v", err)
	}
}

func TestSchemaTypeE2E(t *testing.T) {
	rootCmd.SetArgs([]string{"schema", "Activity"})
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("schema Activity failed: %v", err)
	}
}

func TestMissingConfigError(t *testing.T) {
	t.Setenv("INTERVALS_CONFIG", "/nonexistent/config.yaml")
	t.Setenv("INTERVALS_API_KEY", "")
	t.Setenv("INTERVALS_ATHLETE_ID", "")
	rootCmd.SetArgs([]string{"activities", "list", "--oldest", "2026-01-01"})
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error for missing config")
	}
}

func TestSkillsListTextE2E(t *testing.T) {
	rootCmd.SetArgs([]string{"skills", "list"})
	var out bytes.Buffer
	rootCmd.SetOut(&out)
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("skills list failed: %v", err)
	}
	got := out.String()
	for _, want := range []string{"intervals-shared", "intervals-training", "intervals-settings"} {
		if !bytes.Contains([]byte(got), []byte(want)) {
			t.Errorf("skills list output missing %q\noutput:\n%s", want, got)
		}
	}
}

func TestSkillsListJSONE2E(t *testing.T) {
	rootCmd.SetArgs([]string{"skills", "list", "--format", "json"})
	var out bytes.Buffer
	rootCmd.SetOut(&out)
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("skills list --format json failed: %v", err)
	}
	var entries []map[string]any
	if err := json.Unmarshal(out.Bytes(), &entries); err != nil {
		t.Fatalf("output not valid JSON: %v\noutput:\n%s", err, out.String())
	}
	if len(entries) < 3 {
		t.Fatalf("expected at least 3 skill entries, got %d", len(entries))
	}
	for _, e := range entries {
		if _, ok := e["name"].(string); !ok {
			t.Errorf("entry missing string 'name': %v", e)
		}
		if _, ok := e["description"].(string); !ok {
			t.Errorf("entry missing string 'description': %v", e)
		}
	}
}

func TestSkillsGetRawE2E(t *testing.T) {
	// Reset persistent --format to default before running, since prior tests
	// may have marked it changed.
	rootCmd.PersistentFlags().Set("format", "json")
	rootCmd.PersistentFlags().Lookup("format").Changed = false

	rootCmd.SetArgs([]string{"skills", "get", "intervals-shared"})
	var out bytes.Buffer
	rootCmd.SetOut(&out)
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("skills get failed: %v", err)
	}
	// Body should not start with frontmatter delimiter.
	first := bytes.SplitN(out.Bytes(), []byte("\n"), 2)[0]
	if bytes.Equal(first, []byte("---")) {
		t.Errorf("skills get raw output starts with frontmatter delimiter; want body only.\noutput head:\n%s", string(first))
	}
	if out.Len() == 0 {
		t.Errorf("skills get raw output is empty")
	}
}

func TestSkillsGetUnknownE2E(t *testing.T) {
	rootCmd.SetArgs([]string{"skills", "get", "this-skill-does-not-exist"})
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error for unknown skill")
	}
}

func TestFieldsFilterE2E(t *testing.T) {
	server, cleanup := setupTestEnv(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id": "a1", "name": "Run", "extra": "data",
		})
	})
	defer cleanup()
	_ = server

	rootCmd.SetArgs([]string{"activity", "get", "--activity-id", "a1", "--fields", "id,name"})
	var out bytes.Buffer
	rootCmd.SetOut(&out)
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("execute failed: %v", err)
	}
}

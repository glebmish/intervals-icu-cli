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

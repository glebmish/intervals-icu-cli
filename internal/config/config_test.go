package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadFromFile(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")
	content := `api_key: myapikey
athlete_id: athlete123
`
	if err := os.WriteFile(configPath, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.APIKey != "myapikey" {
		t.Errorf("APIKey = %q, want %q", cfg.APIKey, "myapikey")
	}
	if cfg.AthleteID != "athlete123" {
		t.Errorf("AthleteID = %q, want %q", cfg.AthleteID, "athlete123")
	}
	if cfg.BaseURL != "https://intervals.icu" {
		t.Errorf("BaseURL = %q, want %q", cfg.BaseURL, "https://intervals.icu")
	}
}

func TestLoadMissingFileReturnsDefaults(t *testing.T) {
	cfg, err := Load("/nonexistent/path/config.yaml")
	if err != nil {
		t.Fatalf("Load() error = %v, want nil for missing file", err)
	}

	if cfg.APIKey != "" {
		t.Errorf("APIKey = %q, want empty", cfg.APIKey)
	}
	if cfg.AthleteID != "" {
		t.Errorf("AthleteID = %q, want empty", cfg.AthleteID)
	}
	if cfg.BaseURL != "https://intervals.icu" {
		t.Errorf("BaseURL = %q, want %q", cfg.BaseURL, "https://intervals.icu")
	}
}

func TestEnvVarOverridesFile(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")
	content := `api_key: fileapikey
athlete_id: fileathlete
`
	if err := os.WriteFile(configPath, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	t.Setenv("INTERVALS_API_KEY", "envapikey")
	t.Setenv("INTERVALS_ATHLETE_ID", "envathlete")

	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	cfg.ApplyEnv()

	if cfg.APIKey != "envapikey" {
		t.Errorf("APIKey = %q, want %q", cfg.APIKey, "envapikey")
	}
	if cfg.AthleteID != "envathlete" {
		t.Errorf("AthleteID = %q, want %q", cfg.AthleteID, "envathlete")
	}
}

func TestApplyFlags(t *testing.T) {
	cfg := &Config{
		APIKey:    "original_key",
		AthleteID: "original_athlete",
		BaseURL:   "https://original.example.com",
	}

	cfg.ApplyFlags("flag_key", "flag_athlete", "https://flag.example.com")

	if cfg.APIKey != "flag_key" {
		t.Errorf("APIKey = %q, want %q", cfg.APIKey, "flag_key")
	}
	if cfg.AthleteID != "flag_athlete" {
		t.Errorf("AthleteID = %q, want %q", cfg.AthleteID, "flag_athlete")
	}
	if cfg.BaseURL != "https://flag.example.com" {
		t.Errorf("BaseURL = %q, want %q", cfg.BaseURL, "https://flag.example.com")
	}
}

func TestApplyFlagsEmptyDoesNotOverride(t *testing.T) {
	cfg := &Config{
		APIKey:    "original_key",
		AthleteID: "original_athlete",
		BaseURL:   "https://original.example.com",
	}

	cfg.ApplyFlags("", "", "")

	if cfg.APIKey != "original_key" {
		t.Errorf("APIKey = %q, want %q", cfg.APIKey, "original_key")
	}
	if cfg.AthleteID != "original_athlete" {
		t.Errorf("AthleteID = %q, want %q", cfg.AthleteID, "original_athlete")
	}
	if cfg.BaseURL != "https://original.example.com" {
		t.Errorf("BaseURL = %q, want %q", cfg.BaseURL, "https://original.example.com")
	}
}

func TestValidateMissingAPIKey(t *testing.T) {
	cfg := &Config{
		AthleteID: "athlete123",
		BaseURL:   "https://intervals.icu",
	}

	err := cfg.Validate()
	if err == nil {
		t.Error("Validate() = nil, want error for missing APIKey")
	}
}

func TestValidateMissingAthleteID(t *testing.T) {
	cfg := &Config{
		APIKey:  "myapikey",
		BaseURL: "https://intervals.icu",
	}

	err := cfg.Validate()
	if err == nil {
		t.Error("Validate() = nil, want error for missing AthleteID")
	}
}

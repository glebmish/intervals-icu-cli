package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/glebmish/intervals-icu-cli/internal/cliexit"
	"gopkg.in/yaml.v3"
)

type Config struct {
	APIKey    string `yaml:"api_key"`
	AthleteID string `yaml:"athlete_id"`
	BaseURL   string `yaml:"base_url"`
}

func DefaultPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "intervals", "config.yaml")
}

func Load(path string) (*Config, error) {
	cfg := &Config{BaseURL: "https://intervals.icu"}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, fmt.Errorf("reading config %s: %w", path, err)
	}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parsing config %s: %w", path, err)
	}
	if cfg.BaseURL == "" {
		cfg.BaseURL = "https://intervals.icu"
	}
	return cfg, nil
}

func (c *Config) ApplyEnv() {
	if v := os.Getenv("INTERVALS_API_KEY"); v != "" {
		c.APIKey = v
	}
	if v := os.Getenv("INTERVALS_ATHLETE_ID"); v != "" {
		c.AthleteID = v
	}
	if v := os.Getenv("INTERVALS_BASE_URL"); v != "" {
		c.BaseURL = v
	}
}

func (c *Config) ApplyFlags(apiKey, athleteID, baseURL string) {
	if apiKey != "" {
		c.APIKey = apiKey
	}
	if athleteID != "" {
		c.AthleteID = athleteID
	}
	if baseURL != "" {
		c.BaseURL = baseURL
	}
}

func (c *Config) Validate() error {
	if c.APIKey == "" {
		return &cliexit.AuthError{Err: fmt.Errorf(
			"api_key not configured\n  Set it in %s or INTERVALS_API_KEY env var\n  Run: intervals config init\n  Get an API key at https://intervals.icu/settings",
			DefaultPath())}
	}
	if c.AthleteID == "" {
		return &cliexit.AuthError{Err: fmt.Errorf(
			"athlete_id not configured\n  Set it in %s or INTERVALS_ATHLETE_ID env var\n  Run: intervals config init",
			DefaultPath())}
	}
	return nil
}

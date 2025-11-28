package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

// Profile represents a saved profile.
type Profile struct {
	Name    string            `json:"name"`
	BaseURL string            `json:"base_url"`
	Headers map[string]string `json:"headers,omitempty"`

	AuthType string `json:"auth_type,omitempty"` // none|basic|bearer
	User     string `json:"user,omitempty"`
	Pass     string `json:"pass,omitempty"`
	Token    string `json:"token,omitempty"`
}

// Config is the root config file structure.
type Config struct {
	Profiles map[string]Profile `json:"profiles"`
}

func defaultConfig() *Config {
	return &Config{
		Profiles: make(map[string]Profile),
	}
}

func configPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		home, err2 := os.UserHomeDir()
		if err2 != nil {
			return "", err
		}
		dir = filepath.Join(home, ".go-rest-api-cli")
	} else {
		dir = filepath.Join(dir, "go-rest-api-cli")
	}
	return filepath.Join(dir, "config.json"), nil
}

// Load loads config from disk, or returns an empty config if file is missing.
func Load() (*Config, error) {
	path, err := configPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return defaultConfig(), nil
		}
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	if cfg.Profiles == nil {
		cfg.Profiles = make(map[string]Profile)
	}
	return &cfg, nil
}

// Save writes config to disk.
func Save(cfg *Config) error {
	path, err := configPath()
	if err != nil {
		return err
	}
	if cfg.Profiles == nil {
		cfg.Profiles = make(map[string]Profile)
	}
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

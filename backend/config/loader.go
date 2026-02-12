package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

var globalConfig *Config

// Load reads and parses the configuration file
// Following rule #9: Zero Trust Input - validates all configuration values
func Load(configPath string) (*Config, error) {
	// Read config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Apply defaults
	applyDefaults(&cfg)

	// Apply environment variable overrides
	applyEnvOverrides(&cfg)

	// Validate configuration
	if err := validate(&cfg); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	// Store globally for access throughout application
	globalConfig = &cfg

	return &cfg, nil
}

// Get returns the global configuration
// Must call Load() first
func Get() *Config {
	if globalConfig == nil {
		panic("configuration not loaded - call config.Load() first")
	}
	return globalConfig
}

// applyDefaults sets default values for missing configuration
func applyDefaults(cfg *Config) {
	// Server defaults
	if cfg.Server.Host == "" {
		cfg.Server.Host = "0.0.0.0"
	}
	if cfg.Server.Port == 0 {
		cfg.Server.Port = 8080
	}
	if len(cfg.Server.CORSOrigins) == 0 {
		cfg.Server.CORSOrigins = []string{"http://localhost:5173", "http://localhost:8080"}
	}

	// Download defaults
	if cfg.Download.Path == "" {
		// Use ~/Music as default
		homeDir, err := os.UserHomeDir()
		if err == nil {
			cfg.Download.Path = filepath.Join(homeDir, "Music")
		} else {
			cfg.Download.Path = "./downloads"
		}
	}
	if cfg.Download.FilenameFormat == "" {
		cfg.Download.FilenameFormat = "title-artist"
	}
	if cfg.Download.AudioFormat == "" {
		cfg.Download.AudioFormat = "LOSSLESS"
	}

	// Services defaults
	if cfg.Services.DefaultService == "" {
		cfg.Services.DefaultService = "tidal"
	}

	// UI defaults
	if cfg.UI.Theme == "" {
		cfg.UI.Theme = "default"
	}
	if cfg.UI.ThemeMode == "" {
		cfg.UI.ThemeMode = "auto"
	}
	if cfg.UI.FontFamily == "" {
		cfg.UI.FontFamily = "Inter"
	}

	// Database defaults
	if cfg.Database.Path == "" {
		cfg.Database.Path = "SpotiFLAC"
	}
}

// applyEnvOverrides applies environment variable overrides
// Environment variables use SPOTIFLAC_ prefix and underscore separators
// Example: SPOTIFLAC_SERVER_PORT=9000
func applyEnvOverrides(cfg *Config) {
	// Server overrides
	if port := os.Getenv("SPOTIFLAC_SERVER_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			cfg.Server.Port = p
		}
	}
	if host := os.Getenv("SPOTIFLAC_SERVER_HOST"); host != "" {
		cfg.Server.Host = host
	}

	// Download path override
	if path := os.Getenv("SPOTIFLAC_DOWNLOAD_PATH"); path != "" {
		cfg.Download.Path = path
	}

	// Service overrides
	if service := os.Getenv("SPOTIFLAC_DEFAULT_SERVICE"); service != "" {
		cfg.Services.DefaultService = service
	}
}

// validate checks configuration values for correctness
// Following rule #9: Validate input types and sanitize dangerous values
func validate(cfg *Config) error {
	// Validate server port range (rule #10: Least Privilege - don't allow privileged ports)
	if cfg.Server.Port < 1024 || cfg.Server.Port > 65535 {
		return fmt.Errorf("server port must be between 1024-65535 (non-privileged ports)")
	}

	// Validate download path (rule #9: prevent path traversal)
	if strings.Contains(cfg.Download.Path, "..") {
		return fmt.Errorf("download path cannot contain '..' (path traversal attempt)")
	}

	// Validate audio format
	validFormats := map[string]bool{
		"LOSSLESS": true,
		"6":        true, // Qobuz quality levels
		"7":        true,
		"27":       true,
	}
	if !validFormats[cfg.Download.AudioFormat] {
		return fmt.Errorf("invalid audio format: %s (must be LOSSLESS, 6, 7, or 27)", cfg.Download.AudioFormat)
	}

	// Validate default service
	validServices := map[string]bool{
		"tidal":  true,
		"qobuz":  true,
		"amazon": true,
	}
	if !validServices[cfg.Services.DefaultService] {
		return fmt.Errorf("invalid default service: %s (must be tidal, qobuz, or amazon)", cfg.Services.DefaultService)
	}

	// Validate theme mode
	validThemeModes := map[string]bool{
		"light": true,
		"dark":  true,
		"auto":  true,
	}
	if !validThemeModes[cfg.UI.ThemeMode] {
		return fmt.Errorf("invalid theme mode: %s (must be light, dark, or auto)", cfg.UI.ThemeMode)
	}

	return nil
}

// Save writes the current configuration to a file
// Used when updating settings via API
func Save(cfg *Config, configPath string) error {
	// Validate before saving
	if err := validate(cfg); err != nil {
		return fmt.Errorf("cannot save invalid configuration: %w", err)
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Write with restricted permissions (rule #14: Secure by Default)
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	// Update global config
	globalConfig = cfg

	return nil
}

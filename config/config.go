package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	DefaultAPIURL      = "https://api.butler.coffee"
	ConfigDir          = ".butler-coffee"
	ConfigFile         = "config.json"
	DefaultMinQuantity = 1  // Minimum kg per month
	DefaultMaxQuantity = 10 // Maximum kg per month

	// Token expiry safety margin - consider token expired this many seconds before actual expiration
	TokenExpirySafetyMarginSeconds = 30
)

type Config struct {
	APIURL                string `json:"api_url"`
	AccessToken           string `json:"access_token,omitempty"`
	RefreshToken          string `json:"refresh_token,omitempty"`
	ExpiresAt             string `json:"expires_at,omitempty"`
	RefreshTokenExpiresAt string `json:"refresh_token_expires_at,omitempty"`
	MinQuantityKg         int    `json:"min_quantity_kg,omitempty"`
	MaxQuantityKg         int    `json:"max_quantity_kg,omitempty"`
}

func GetAPIURL() string {
	// Check for BASE_HOSTNAME environment variable
	if hostname := os.Getenv("BASE_HOSTNAME"); hostname != "" {
		return hostname
	}
	return DefaultAPIURL
}

func GetConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ConfigDir, ConfigFile), nil
}

func LoadConfig() (*Config, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	apiURL := GetAPIURL()

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return &Config{
			APIURL:        apiURL,
			MinQuantityKg: DefaultMinQuantity,
			MaxQuantityKg: DefaultMaxQuantity,
		}, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	// Override with environment variable if set, otherwise use config or default
	if envURL := os.Getenv("BASE_HOSTNAME"); envURL != "" {
		cfg.APIURL = envURL
	} else if cfg.APIURL == "" {
		cfg.APIURL = apiURL
	}

	// Set default quantity limits if not configured
	if cfg.MinQuantityKg == 0 {
		cfg.MinQuantityKg = DefaultMinQuantity
	}
	if cfg.MaxQuantityKg == 0 {
		cfg.MaxQuantityKg = DefaultMaxQuantity
	}

	return &cfg, nil
}

func (c *Config) Save() error {
	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return err
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0600)
}

func (c *Config) IsAuthenticated() bool {
	return c.AccessToken != ""
}

func (c *Config) IsTokenExpired() bool {
	if c.ExpiresAt == "" {
		return false
	}

	expiresAt, err := parseTimestamp(c.ExpiresAt)
	if err != nil {
		return true
	}

	// Consider token expired before actual expiration for safety margin
	safetyMargin := time.Duration(TokenExpirySafetyMarginSeconds) * time.Second
	return time.Now().Add(safetyMargin).After(expiresAt)
}

func (c *Config) IsRefreshTokenExpired() bool {
	if c.RefreshTokenExpiresAt == "" {
		return false
	}

	expiresAt, err := parseTimestamp(c.RefreshTokenExpiresAt)
	if err != nil {
		return true
	}

	return time.Now().After(expiresAt)
}

// parseTimestamp handles both Unix timestamp (milliseconds) and RFC3339 formats
func parseTimestamp(timestamp string) (time.Time, error) {
	// Try parsing as Unix timestamp in milliseconds (string format)
	var unixMs int64
	if _, err := fmt.Sscanf(timestamp, "%d", &unixMs); err == nil {
		return time.Unix(unixMs/1000, (unixMs%1000)*1000000), nil
	}

	// Fallback to RFC3339 format
	return time.Parse(time.RFC3339, timestamp)
}

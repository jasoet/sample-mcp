package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"

	"sample-mcp/pkg/db"
)

const (
	// EnvMCPServerConfig is the environment variable name for the MCP server configuration file path
	EnvMCPServerConfig = "MCP_SERVER_CONFIG"
)

// DefaultConnectionConfig returns the default database connection configuration
func DefaultConnectionConfig() *db.ConnectionConfig {
	return &db.ConnectionConfig{
		DbType:       db.Postgresql,
		Host:         "localhost",
		Port:         5432,
		Username:     "jasoet",
		Password:     "localhost",
		DbName:       "mcp_db",
		Timeout:      3 * time.Second,
		MaxIdleConns: 5,
		MaxOpenConns: 10,
	}
}

// Config represents the application configuration
type Config struct {
	Database *db.ConnectionConfig `yaml:"database"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		Database: DefaultConnectionConfig(),
	}
}

// LoadConfig loads the configuration from the config file
// It looks for the config file in the following order:
// 1. Path specified in MCP_SERVER_CONFIG environment variable
// 2. config.yml in the same directory as the binary
// If no config file is found, it returns the default configuration
func LoadConfig() (*Config, error) {
	config := DefaultConfig()

	// Try to load from environment variable
	configPath := os.Getenv(EnvMCPServerConfig)
	if configPath == "" {
		// Try to load from the same directory as the binary
		execPath, err := os.Executable()
		if err != nil {
			return config, fmt.Errorf("failed to get executable path: %w", err)
		}
		configPath = filepath.Join(filepath.Dir(execPath), "config.yml")
	}

	// Check if the config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Config file doesn't exist, return default config
		return config, nil
	}

	// Read the config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return config, fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse the config file
	if err := yaml.Unmarshal(data, config); err != nil {
		return config, fmt.Errorf("failed to parse config file: %w", err)
	}

	return config, nil
}

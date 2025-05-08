package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"

	"sample-mcp/pkg/db"
)

func TestDefaultConnectionConfig(t *testing.T) {
	config := DefaultConnectionConfig()

	assert.Equal(t, db.Postgresql, config.DbType)
	assert.Equal(t, "localhost", config.Host)
	assert.Equal(t, 5432, config.Port)
	assert.Equal(t, "jasoet", config.Username)
	assert.Equal(t, "localhost", config.Password)
	assert.Equal(t, "mcp_db", config.DbName)
	assert.Equal(t, 3*time.Second, config.Timeout)
	assert.Equal(t, 5, config.MaxIdleConns)
	assert.Equal(t, 10, config.MaxOpenConns)
}

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	assert.NotNil(t, config.Database)
	assert.Equal(t, db.Postgresql, config.Database.DbType)
	assert.Equal(t, "localhost", config.Database.Host)
	assert.Equal(t, 5432, config.Database.Port)
	assert.Equal(t, "jasoet", config.Database.Username)
	assert.Equal(t, "localhost", config.Database.Password)
	assert.Equal(t, "mcp_db", config.Database.DbName)
	assert.Equal(t, 3*time.Second, config.Database.Timeout)
	assert.Equal(t, 5, config.Database.MaxIdleConns)
	assert.Equal(t, 10, config.Database.MaxOpenConns)
}

func TestLoadConfig_NoConfigFile(t *testing.T) {
	// Ensure environment variable is not set
	os.Unsetenv(EnvMCPServerConfig)

	config, err := LoadConfig()
	assert.NoError(t, err)
	assert.NotNil(t, config)
	assert.NotNil(t, config.Database)
	assert.Equal(t, db.Postgresql, config.Database.DbType)
}

func TestLoadConfig_FromEnvVar(t *testing.T) {
	// Create a temporary config file
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yml")

	// Create a custom config
	customConfig := &Config{
		Database: &db.ConnectionConfig{
			DbType:       db.Mysql,
			Host:         "custom-host",
			Port:         3306,
			Username:     "custom-user",
			Password:     "custom-password",
			DbName:       "custom-db",
			Timeout:      5 * time.Second,
			MaxIdleConns: 10,
			MaxOpenConns: 20,
		},
	}

	// Write the config to the file
	data, err := yaml.Marshal(customConfig)
	assert.NoError(t, err)
	err = os.WriteFile(configPath, data, 0644)
	assert.NoError(t, err)

	// Set the environment variable
	os.Setenv(EnvMCPServerConfig, configPath)
	defer os.Unsetenv(EnvMCPServerConfig)

	// Load the config
	config, err := LoadConfig()
	assert.NoError(t, err)
	assert.NotNil(t, config)
	assert.NotNil(t, config.Database)
	assert.Equal(t, db.Mysql, config.Database.DbType)
	assert.Equal(t, "custom-host", config.Database.Host)
	assert.Equal(t, 3306, config.Database.Port)
	assert.Equal(t, "custom-user", config.Database.Username)
	assert.Equal(t, "custom-password", config.Database.Password)
	assert.Equal(t, "custom-db", config.Database.DbName)
	assert.Equal(t, 5*time.Second, config.Database.Timeout)
	assert.Equal(t, 10, config.Database.MaxIdleConns)
	assert.Equal(t, 20, config.Database.MaxOpenConns)
}

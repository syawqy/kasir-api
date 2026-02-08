package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

type ServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

func LoadConfig() (*Config, error) {
	// Enable automatic environment variable reading
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Try to read .env file if it exists (for local development)
	// On server environments like Leapcell, env vars are set directly
	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig() // Ignore error, env vars will be used as fallback
	}

	// Build config directly from viper (works with both .env file and env vars)
	config := Config{
		Server: ServerConfig{
			Port: viper.GetString("SERVER_PORT"),
		},
		Database: DatabaseConfig{
			Host:     viper.GetString("DB_HOST"),
			Port:     viper.GetString("DB_PORT"),
			User:     viper.GetString("DB_USER"),
			Password: viper.GetString("DB_PASSWORD"),
			Name:     viper.GetString("DB_NAME"),
		},
	}

	// Set defaults
	if config.Server.Port == "" {
		config.Server.Port = "8080" // Default server port
	}

	// Validate required fields
	if config.Database.Host == "" {
		return nil, fmt.Errorf("DB_HOST is required")
	}
	if config.Database.Port == "" {
		config.Database.Port = "5432" // Default PostgreSQL port
	}

	return &config, nil
}

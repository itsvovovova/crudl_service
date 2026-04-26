package config

import (
	"fmt"
	"os"
)

type Config struct {
	Server   *ServerConfig
	Database *DatabaseConfig
	JWT      *JWTConfig
}

type ServerConfig struct {
	Port     string
	LogLevel string
}

type DatabaseConfig struct {
	Username      string
	Name          string
	Password      string
	Host          string
	Port          string
	SSLMode       string
	PathMigration string
}

type JWTConfig struct {
	SecretKey string
}

func InitConfig() (*Config, error) {
	cfg := newConfig()
	return cfg, cfg.Validate()
}

func newConfig() *Config {
	sslMode := os.Getenv("DB_SSL_MODE")
	if sslMode == "" {
		sslMode = "disable"
	}
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}
	return &Config{
		Server: &ServerConfig{
			Port:     port,
			LogLevel: os.Getenv("LOG_LEVEL"),
		},
		Database: &DatabaseConfig{
			Username:      os.Getenv("DB_USER"),
			Password:      os.Getenv("DB_PASSWORD"),
			Host:          os.Getenv("DB_HOST"),
			Port:          os.Getenv("DB_PORT"),
			Name:          os.Getenv("DB_NAME"),
			SSLMode:       sslMode,
			PathMigration: os.Getenv("DB_PATH_MIGRATION"),
		},
		JWT: &JWTConfig{
			SecretKey: os.Getenv("JWT_SECRET_KEY"),
		},
	}
}

func (c *Config) Validate() error {
	required := map[string]string{
		"DB_USER":             c.Database.Username,
		"DB_PASSWORD":         c.Database.Password,
		"DB_HOST":             c.Database.Host,
		"DB_PORT":             c.Database.Port,
		"DB_NAME":             c.Database.Name,
		"JWT_SECRET_KEY":      c.JWT.SecretKey,
		"DB_PATH_MIGRATION":   c.Database.PathMigration,
	}
	for key, val := range required {
		if val == "" {
			return fmt.Errorf("required environment variable %s is not set", key)
		}
	}
	return nil
}

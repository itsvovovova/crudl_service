package config

import (
	"os"
)

var CurrentConfig = NewConfig()

type Config struct {
	Server   *ServerConfig
	Database *DatabaseConfig
}

type ServerConfig struct {
	Port     string
	Host     string
	HostPort string
}

type DatabaseConfig struct {
	Username      string
	Name          string
	Password      string
	Host          string
	Port          string
	DatabaseName  string
	SSLMode       string
	PathMigration string
}

func NewConfig() *Config {
	var databaseConfig = DatabaseConfig{
		Username:      os.Getenv("DB_USER"),
		Password:      os.Getenv("DB_PASSWORD"),
		Host:          os.Getenv("DB_HOST"),
		Port:          os.Getenv("DB_PORT"),
		Name:          os.Getenv("DB_NAME"),
		DatabaseName:  os.Getenv("DB_NAME"),
		SSLMode:       os.Getenv("DB_SSL_MODE"),
		PathMigration: os.Getenv("DB_PATH_MIGRATION"),
	}

	var serverConfig = ServerConfig{
		Port:     os.Getenv("SERVER_PORT"),
		Host:     os.Getenv("SERVER_HOST"),
		HostPort: os.Getenv("HOST_PORT"),
	}
	var Config = &Config{
		Server:   &serverConfig,
		Database: &databaseConfig,
	}
	return Config
}

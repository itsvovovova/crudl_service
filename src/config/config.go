package config

import (
	"os"
)

var CurrentConfig = NewConfig()

type Config struct {
	Server   *ServerConfig
	Database *DatabaseConfig
	Logger   *LoggerConfig
}

type ServerConfig struct {
	Port     string
	Host     string
	HostPort string
}

type DatabaseConfig struct {
	Username     string
	Name         string
	Password     string
	Host         string
	Port         string
	DatabaseName string
	SSLMode      string
}

type LoggerConfig struct {
	Level  string
	Format string
}

func NewConfig() *Config {
	var databaseConfig = DatabaseConfig{
		Username:     os.Getenv("DB_USER"),
		Password:     os.Getenv("DB_PASSWORD"),
		Host:         os.Getenv("DB_HOST"),
		Port:         os.Getenv("DB_PORT"),
		Name:         os.Getenv("DB_NAME"),
		DatabaseName: os.Getenv("DB_NAME"),
		SSLMode:      os.Getenv("DB_SSL_MODE"),
	}

	var loggerConfig = LoggerConfig{
		Level:  os.Getenv("LOG_LEVEL"),
		Format: os.Getenv("LOG_FORMAT"),
	}

	var serverConfig = ServerConfig{
		Port:     os.Getenv("SERVER_PORT"),
		Host:     os.Getenv("SERVER_HOST"),
		HostPort: os.Getenv("HOST_PORT"),
	}
	var Config = &Config{
		Server:   &serverConfig,
		Database: &databaseConfig,
		Logger:   &loggerConfig,
	}
	return Config
}

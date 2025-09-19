package config

import (
	"os"
)

var CurrentConfig *Config

func InitConfig() {
	CurrentConfig = NewConfig()
}

func ShutdownConfig() {
	CurrentConfig = nil
}

type Config struct {
	Server   *ServerConfig
	Database *DatabaseConfig
	JWT      *JWTConfig
}

type ServerConfig struct {
	Port     string
	Host     string
	HostPort string
	LogLevel string
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

type JWTConfig struct {
	SecretKey string
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
		LogLevel: os.Getenv("LOG_LEVEL"),
	}
	var jwtConfig = JWTConfig{
		SecretKey: os.Getenv("JWT_SECRET_KEY"),
	}

	var Config = &Config{
		Server:   &serverConfig,
		Database: &databaseConfig,
		JWT:      &jwtConfig,
	}
	return Config
}

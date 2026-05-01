package db

import (
	"crudl_service/src/config"
	"database/sql"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

type NotFoundError struct{}

func (e *NotFoundError) Error() string {
	return "not found"
}

func buildConnURL(cfg *config.DatabaseConfig) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Name, cfg.SSLMode,
	)
}

func InitDB(cfg *config.DatabaseConfig) (*sql.DB, error) {
	log.Info("Initializing database connection")
	urlConnection := buildConnURL(cfg)

	conn, err := sql.Open("postgres", urlConnection)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	const maxRetries = 10
	const retryDelay = 3 * time.Second

	for i := 1; i <= maxRetries; i++ {
		if err = conn.Ping(); err == nil {
			break
		}
		log.WithError(err).Warnf("Database not ready, retry %d/%d", i, maxRetries)
		if i == maxRetries {
			_ = conn.Close()
			return nil, fmt.Errorf("database connection failed after %d retries", maxRetries)
		}
		time.Sleep(retryDelay)
	}
	log.Info("Database connection established successfully")

	m, err := migrate.New(cfg.PathMigration, urlConnection)
	if err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("failed to create migration instance: %w", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		_ = conn.Close()
		return nil, fmt.Errorf("failed to apply migrations: %w", err)
	}

	sourceErr, dbErr := m.Close()
	if sourceErr != nil || dbErr != nil {
		log.WithFields(log.Fields{"source": sourceErr, "db": dbErr}).Warn("Error closing migrate instance")
	}

	log.Info("Database migrations completed successfully")

	return conn, nil
}

func CloseDB(db *sql.DB) {
	if db == nil {
		return
	}
	log.Info("Closing database connection")
	if err := db.Close(); err != nil {
		log.WithError(err).Error("Error closing database connection")
	} else {
		log.Info("Database connection closed successfully")
	}
}

package db

import (
	"crudl_service/src/config"
	"database/sql"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

var db *sql.DB

func init() {
	db = InitDB()
}

func InitDB() *sql.DB {
	urlConnection := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.CurrentConfig.Database.Username,
		config.CurrentConfig.Database.Password,
		config.CurrentConfig.Database.Host,
		config.CurrentConfig.Database.Port,
		config.CurrentConfig.Database.Name)

	connStr := urlConnection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to open database connection:", err)
	}
	if err := db.Ping(); err != nil {
		_ = db.Close()
		log.Fatal("Database connection error:", err)
	}

	m, err := migrate.New(
		config.CurrentConfig.Database.PathMigration,
		urlConnection)
	if err != nil {
		log.Fatal("Failed to create migration instance:", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("Failed to apply migrations:", err)
	}
	return db
}

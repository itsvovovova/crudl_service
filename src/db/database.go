package db

import (
	"crudl_service/src/config"
	"database/sql"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var db *sql.DB

func init() {
	db = InitDB()
}

func InitDB() *sql.DB {
	urlConnection := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.CurrentConfig.Database.Host,
		config.CurrentConfig.Database.Port,
		config.CurrentConfig.Database.Username,
		config.CurrentConfig.Database.Password,
		config.CurrentConfig.Database.Name)

	connStr := urlConnection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to open database connection")
	}
	if err := db.Ping(); err != nil {
		_ = db.Close()
		log.Fatal("Database connection error")
	}
	m, err := migrate.New(
		config.CurrentConfig.Database.PathMigration,
		urlConnection)
	if err != nil {
		log.Fatal("Failed to create migration instance")
	}
	if err := m.Up(); err != nil {
		log.Fatal("Failed to apply migrations")
	}
	return db
}

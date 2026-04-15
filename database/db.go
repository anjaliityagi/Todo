package database

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var (
	Todo *sqlx.DB
)

const (
	SSLModeDisable SSLMode = "disable"
)

type SSLMode string

func ConnectandMigrate(host, port, databaseName, user, password string, sslMode SSLMode) error {
	connectionStr := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s", host, port, databaseName, user, password, sslMode)

	DB, err := sqlx.Open("postgres", connectionStr)
	if err != nil {
		return err
	}
	err = DB.Ping()
	if err != nil {
		return fmt.Errorf("database ping failed %w", err)
	}
	fmt.Println("Database connected successfully")
	Todo = DB
	return migrateUp(DB)
}

func migrateUp(db *sqlx.DB) error {
	fmt.Println("Starting database migrations...")
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://database/migrations",
		"postgres", driver)

	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("No new migrations to apply")
			return nil
		}
		return fmt.Errorf("migration failed: %w", err)
	}
	fmt.Println("migrations applied successfully")
	return nil
}

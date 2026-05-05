package migrator

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

type Migrator struct {
	srcDriver source.Driver
}

func NewMigrator(sqlFiles embed.FS, dirname string) (*Migrator, error) {
	d, err := iofs.New(sqlFiles, dirname)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize migration driver: %w", err)
	}
	return &Migrator{srcDriver: d}, nil
}

func (m *Migrator) ApplyMigrations(db *sql.DB) (err error) {
	driver, err := sqlite.WithInstance(db, &sqlite.Config{})
	if err != nil {
		return fmt.Errorf("unable to create db instance: %w", err)
	}

	migrator, err := migrate.NewWithInstance("migration_embeded_sql_files", m.srcDriver, "sqlite", driver)
	if err != nil {
		return fmt.Errorf("unable to create migration: %w", err)
	}

	if err = migrator.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		migrator.Close()
		return fmt.Errorf("unable to apply migrations: %w", err)
	}

	return nil
}

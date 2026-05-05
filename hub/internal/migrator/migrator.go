package migrator

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
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
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("unable to create db instance: %w", err)
	}

	migrator, err := migrate.NewWithInstance("migration_embeded_sql_files", m.srcDriver, "psql_db", driver)
	if err != nil {
		return fmt.Errorf("unable to create migration: %w", err)
	}

	defer func() {
		closeErr, _ := migrator.Close()
		if err == nil {
			err = closeErr
		}
	}()

	if err = migrator.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("unable to apply migrations: %w", err)
	}

	return nil
}

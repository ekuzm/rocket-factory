package migrator

import (
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"
)

type migrator struct {
	db           *sql.DB
	migrationDir string
}

func New(db *sql.DB, migrationDir string) *migrator {
	return &migrator{
		db:           db,
		migrationDir: migrationDir,
	}
}

func (m *migrator) Up() error {
	if err := goose.Up(m.db, m.migrationDir); err != nil {
		return fmt.Errorf("goose.Up: %w", err)
	}

	return nil
}

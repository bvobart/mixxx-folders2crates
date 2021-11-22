package mixxxdb

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type MixxxDB interface {
	// TODO: define an interface.
	// e.g. InsertCrate(crate Crate) error
	// e.g. WipeCrate(name string) error
	// etc.
}

type mixxxDB struct {
	db *sql.DB
}

func OpenDefault() (MixxxDB, error) {
	return Open(DefaultMixxxDBPath)
}

// Open opens a Mixxx DB SQLite file
func Open(mixxxdbPath string) (MixxxDB, error) {
	db, err := sql.Open("sqlite3", mixxxdbPath)
	if err != nil {
		return nil, fmt.Errorf("cannot open Mixxx DB file %s: %w", mixxxdbPath, err)
	}

	return &mixxxDB{db}, db.Ping()
}

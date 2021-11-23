package mixxxdb

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

// Opens the default MixxxDB SQLite file, as defined by the platform-specific `DefaultMixxxDBPath` variable.
func OpenDefault() (MixxxDB, error) {
	return Open(DefaultMixxxDBPath)
}

// Open opens a Mixxx DB SQLite file
func Open(mixxxdbPath string) (MixxxDB, error) {
	db, err := sql.Open("sqlite3", mixxxdbPath)
	if err != nil {
		return nil, fmt.Errorf("cannot open Mixxx DB file %s: %w", mixxxdbPath, err)
	}

	crates := &cratesDB{db}
	tracks := &tracksDB{db}
	return &mixxxDB{crates, tracks}, db.Ping()
}

//-------------------------------------------------------------------------------------------------------------

type MixxxDB interface {
	Crates() CratesDB
	Tracks() TracksDB
}

type CratesDB interface {
	FindByName(name string) (*Crate, error)
	Insert(crate Crate) error
	List() ([]Crate, error)
	WipeTracks(crateid uint) error
}

type TracksDB interface {
	FindByPath(filepath string) (*Track, error)
}

//-------------------------------------------------------------------------------------------------------------

type mixxxDB struct {
	crates CratesDB
	tracks TracksDB
}

func (m mixxxDB) Crates() CratesDB {
	return m.crates
}

func (m mixxxDB) Tracks() TracksDB {
	return m.tracks
}

//-------------------------------------------------------------------------------------------------------------

type cratesDB struct {
	db *sql.DB
}

func (c *cratesDB) FindByName(name string) (*Crate, error) {
	var id uint
	err := c.db.QueryRow("select id from crates where name = ?", name).Scan(&id)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &Crate{ID: id, Name: name}, err
}

func (c *cratesDB) Insert(crate Crate) error {
	return errors.New("TODO")
}

func (c *cratesDB) List() ([]Crate, error) {
	rows, err := c.db.Query("select id, name from crates")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	crates := []Crate{}
	for rows.Next() {
		var id uint
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			return nil, err
		}

		crates = append(crates, Crate{ID: id, Name: name})
	}

	// TODO: retrieve each crate's tracks.

	return crates, rows.Err()
}

func (c *cratesDB) WipeTracks(crateid uint) error {
	return errors.New("TODO")
}

type tracksDB struct {
	db *sql.DB
}

func (t *tracksDB) FindByPath(filepath string) (*Track, error) {
	return nil, errors.New("TODO")
}

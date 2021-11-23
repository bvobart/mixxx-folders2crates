package mixxxdb

import (
	"database/sql"
	"fmt"

	folders2crates "github.com/bvobart/mixxx-folders2crates"
	_ "github.com/mattn/go-sqlite3"
)

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
	FindByName(name string) (*folders2crates.Crate, error)
	Insert(crate folders2crates.Crate) error
	InsertMany(crates []folders2crates.Crate) error
	List() ([]folders2crates.Crate, error)
	// TODO: Wipe?
}

type TracksDB interface {
	FindByPath(filepath string) (*folders2crates.Track, error)
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

func (c *cratesDB) FindByName(name string) (*folders2crates.Crate, error) {
	var id uint
	err := c.db.QueryRow("select id from crates where name = ?", name).Scan(&id)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &folders2crates.Crate{Id: id, Name: name}, err
}

func (c *cratesDB) Insert(crate folders2crates.Crate) error {
	return nil
}

func (c *cratesDB) InsertMany(crates []folders2crates.Crate) error {
	return nil
}

func (c *cratesDB) List() ([]folders2crates.Crate, error) {
	rows, err := c.db.Query("select id, name from crates")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	crates := []folders2crates.Crate{}
	for rows.Next() {
		var id uint
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			return nil, err
		}

		crates = append(crates, folders2crates.Crate{Id: id, Name: name})
	}

	return crates, rows.Err()
}

type tracksDB struct {
	db *sql.DB
}

func (t *tracksDB) FindByPath(filepath string) (*folders2crates.Track, error) {
	return nil, nil
}

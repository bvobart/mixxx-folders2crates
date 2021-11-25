package mixxxdb

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

var ErrTrackNoID = errors.New("track has no ID")

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
	Insert(crate Crate) (int64, error)
	InsertTracks(crateT CrateTracks) (int64, error)
	List() ([]Crate, error)
	WipeTracks(crateid int64) error
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
	var id int64
	err := c.db.QueryRow("select id from crates where name = ?", name).Scan(&id)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &Crate{ID: id, Name: name}, err
}

func (c *cratesDB) Insert(crate Crate) (int64, error) {
	result, err := c.db.Exec("insert into crates(name) values(?)", crate.Name)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func (c *cratesDB) InsertTracks(crate CrateTracks) (crateID int64, err error) {
	// if crate already has a valid ID, then we assume it is already in the DB under this ID.
	crateID = crate.ID
	if crateID == 0 {
		crateID, err = c.Insert(crate.Crate)
		if err != nil {
			return 0, err
		}
	}

	inserts := []string{}
	args := []interface{}{}
	for _, track := range crate.Tracks {
		// check that all tracks have IDs
		if track.ID == 0 {
			return 0, fmt.Errorf("%w: %s", ErrTrackNoID, track.Path)
		}

		inserts = append(inserts, "(?, ?)")
		args = append(args, crateID, track.ID)
	}

	query := fmt.Sprint("insert into crate_tracks(crate_id, track_id) values", strings.Join(inserts, ","))
	_, err = c.db.Exec(query, args...)
	return crateID, err
}

func (c *cratesDB) List() ([]Crate, error) {
	rows, err := c.db.Query("select id, name from crates")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	crates := []Crate{}
	for rows.Next() {
		var id int64
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			return nil, err
		}

		crates = append(crates, Crate{ID: id, Name: name})
	}

	return crates, rows.Err()
}

func (c *cratesDB) WipeTracks(crateid int64) error {
	_, err := c.db.Exec("delete from crate_tracks where crate_id = ?", crateid)
	return err
}

type tracksDB struct {
	db *sql.DB
}

func (t *tracksDB) FindByPath(filepath string) (*Track, error) {
	var id uint
	err := t.db.QueryRow("select id from track_locations where location = ? and fs_deleted = 0", filepath).Scan(&id)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &Track{ID: id, Path: filepath}, err
}

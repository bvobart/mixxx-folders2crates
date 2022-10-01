package mixxxdb

import (
	"fmt"
	"path"
	"runtime"

	"github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/sqlite"

	"github.com/bvobart/mixxx-folders2crates/utils"
)

// See https://manual.mixxx.org/2.3/en/chapters/appendix/settings_directory.html?#location
// If this is an empty string, then the current OS is not supported.
var DefaultMixxxDBPath = func() string {
	switch runtime.GOOS {

	// Unix-based OSes
	case "linux":
		fallthrough
	case "freebsd":
		fallthrough
	case "netbsd":
		fallthrough
	case "openbsd":
		fallthrough
	case "plan9":
		fallthrough
	case "solaris":
		return path.Join(utils.HomeDir(), ".mixxx", "mixxxdb.sqlite")

	// Windows
	case "windows":
		return path.Join(utils.HomeDir(), "Local Settings", "Application Data", "Mixxx", "mixxxdb.sqlite")

	// MacOS
	case "darwin":
		return path.Join(utils.HomeDir(), "Library", "Containers", "org.mixxx.mixxx", "Data", "Library", "Application Support", "Mixxx", "mixxxdb.sqlite")

	// Unsupported OSes
	default:
		return ""

	}
}()

// Opens the default MixxxDB SQLite file, as defined by the platform-specific `DefaultMixxxDBPath` variable.
func OpenDefault() (MixxxDB, error) {
	return Open(DefaultMixxxDBPath)
}

// Open opens a Mixxx DB SQLite file
func Open(mixxxdbPath string) (MixxxDB, error) {
	db, err := OpenSession(mixxxdbPath)
	if err != nil {
		return nil, err
	}

	return NewMixxxDB(db), nil
}

// OpenSession opens a db.Session to the Mixxx DB SQLite file at the given location.
// Use with `NewMixxxDB(session)`, or open a `session.Tx(func (tx db.Session) error)` first, then use the
// `tx` session in its callback to call `NewMixxxDB(tx)` such that all updates are bundled in one transaction.
func OpenSession(mixxxdbPath string) (db.Session, error) {
	settings := sqlite.ConnectionURL{Database: mixxxdbPath}
	db, err := sqlite.Open(settings)
	if err != nil {
		return nil, fmt.Errorf("cannot open Mixxx DB file %s: %w", mixxxdbPath, err)
	}

	return db, db.Ping()
}

// NewMixxxDB creates a new instance of a struct fulfilling the MixxxDB interface.
// No fancy trickery here, just an object instantiation.
func NewMixxxDB(session db.Session) MixxxDB {
	crates := NewCratesDB(session)
	crateTracks := NewCrateTracksDB(session)
	tracks := NewTracksDB(session)
	trackLocations := NewTrackLocationsDB(session)
	return &mixxxDB{session, crates, crateTracks, tracks, trackLocations}
}

//-------------------------------------------------------------------------------------------------------------

type MixxxDB interface {
	Session() db.Session
	Crates() CratesDB
	CrateTracks() CrateTracksDB
	Tracks() TracksDB
	TrackLocations() TrackLocationsDB
}

type CratesDB interface {
	db.Collection
	FindByName(name string) (*Crate, error)
	InsertCrate(crate Crate) (*Crate, error)
}

type CrateTracksDB interface {
	db.Collection
	InsertTracks(crateID int, trackIDs []int) error
	WipeCrate(crateID int) error
}

type TracksDB interface {
	db.Collection
	FindByLocationID(id int) (*Track, error)
}

type TrackLocationsDB interface {
	db.Collection
	FindByLocation(filepath string) (*TrackLocation, error)
}

//-------------------------------------------------------------------------------------------------------------

type mixxxDB struct {
	session        db.Session
	crates         CratesDB
	crateTracks    CrateTracksDB
	tracks         TracksDB
	trackLocations TrackLocationsDB
}

func (m mixxxDB) Session() db.Session {
	return m.session
}

func (m mixxxDB) Crates() CratesDB {
	return m.crates
}

func (m mixxxDB) CrateTracks() CrateTracksDB {
	return m.crateTracks
}

func (m mixxxDB) Tracks() TracksDB {
	return m.tracks
}

func (m mixxxDB) TrackLocations() TrackLocationsDB {
	return m.trackLocations
}

//-------------------------------------------------------------------------------------------------------------

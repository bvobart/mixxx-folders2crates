package mixxxdb

import (
	"errors"

	"github.com/upper/db/v4"
)

// Compile-time check on conformance to interface.
var _ = db.Record(&Track{})
var _ = db.Store(&tracksDB{})
var _ = TracksDB(&tracksDB{})

// Track represents a music file that can be played by Mixxx.
type Track struct {
	// The track's ID in Mixxx DB. Can be found in `library` table. Note: do not confuse with the ID in `track_locations` table.
	ID int `db:"id,omitempty"`

	// ID of the current TrackLocation of this music track.
	Location int `db:"location"`
}

func (_ *Track) Store(sess db.Session) db.Store {
	return NewTracksDB(sess)
}

//-------------------------------------------------------------------------------------------------------------

func NewTracksDB(sess db.Session) TracksDB {
	return &tracksDB{sess.Collection("library")}
}

type tracksDB struct {
	db.Collection
}

func (tracks *tracksDB) FindByLocationID(id int) (*Track, error) {
	var track Track
	err := tracks.Find(db.Cond{"location": id}).One(&track)
	if errors.Is(err, db.ErrNoMoreRows) {
		return nil, nil
	}
	return &track, err
}

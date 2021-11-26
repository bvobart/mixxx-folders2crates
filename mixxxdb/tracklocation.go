package mixxxdb

import (
	"errors"

	"github.com/upper/db/v4"
)

// Compile-time check on conformance to interface.
var _ = db.Record(&TrackLocation{})
var _ = db.Store(&trackLocationsDB{})
var _ = TrackLocationsDB(&trackLocationsDB{})

type TrackLocation struct {
	// ID of this track location entity. NOTE: this is NOT the track ID.
	ID int `db:"id,omitempty"`

	// The full filepath to the track on disk.
	Location string `db:"location"`

	// The filename of the track, i.e. only the last section of the path. `this.Directory + this.Filename == this.Location`
	Filename string `db:"filename"`

	// The directory the music file is located in. `this.Directory + this.Filename == this.Location`
	Directory string `db:"directory"`

	// Size of the music file on disk
	FileSize int `db:"filesize"`

	// Whether the file is deleted from disk (no longer available on this.Location)
	IsDeleted bool `db:"fs_deleted"`

	// Whether the track needs verification. No idea what kind of verification though, check the Mixxx docs for this.
	NeedsVerification bool `db:"needs_verification"`
}

func (_ *TrackLocation) Store(sess db.Session) db.Store {
	return NewTrackLocationsDB(sess)
}

//-------------------------------------------------------------------------------------------------------------

func NewTrackLocationsDB(sess db.Session) TrackLocationsDB {
	return &trackLocationsDB{sess.Collection("track_locations")}
}

type trackLocationsDB struct {
	db.Collection
}

func (tracklocs *trackLocationsDB) FindByLocation(filepath string) (*TrackLocation, error) {
	var loc TrackLocation
	err := tracklocs.Find(db.Cond{"location": filepath}).One(&loc)
	if errors.Is(err, db.ErrNoMoreRows) {
		return nil, nil
	}
	return &loc, err
}

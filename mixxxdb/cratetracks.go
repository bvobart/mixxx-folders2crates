package mixxxdb

import (
	"fmt"

	"github.com/upper/db/v4"
)

// Compile-time check on conformance to interface.
var _ = db.Record(&CrateTrack{})
var _ = db.Store(&crateTracksDB{})
var _ = CrateTracksDB(&crateTracksDB{})

// CrateTrack is the couple table enabling the many-to-many relationship between crates and tracks
type CrateTrack struct {
	// ID of the crate in CratesDB
	CrateID int `db:"crate_id"`

	// ID of the track in TracksDB
	TrackID int `db:"track_id"`
}

func (_ *CrateTrack) Store(sess db.Session) db.Store {
	return NewCrateTracksDB(sess)
}

//-------------------------------------------------------------------------------------------------------------

func NewCrateTracksDB(sess db.Session) CrateTracksDB {
	return &crateTracksDB{sess.Collection("crate_tracks")}
}

type crateTracksDB struct {
	db.Collection
}

func (cratetracks *crateTracksDB) InsertTracks(crateID int, trackIDs []int) error {
	// use a transaction for batch insert.
	return cratetracks.Session().Tx(func(sess db.Session) error {
		cratetracks := NewCrateTracksDB(sess)
		for _, trackID := range trackIDs {
			ct := CrateTrack{CrateID: crateID, TrackID: trackID}

			if _, err := cratetracks.Insert(ct); err != nil {
				return fmt.Errorf("failed to insert %+v: %w", ct, err)
			}
		}
		return nil
	})
}

func (cratetracks *crateTracksDB) WipeCrate(crateID int) error {
	return cratetracks.Find(db.Cond{"crate_id": crateID}).Delete()
}

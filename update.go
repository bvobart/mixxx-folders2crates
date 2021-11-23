package folders2crates

import (
	"errors"
	"fmt"

	"github.com/bvobart/mixxx-folders2crates/mixxxdb"
	"github.com/hashicorp/go-multierror"
)

var ErrTrackNotFound = errors.New("track not found in MixxxDB")

// TODO: insert into Mixxx's SQLite DB
// first focus on starting with an empty crates table (but songs already analysed so library and track_locations are populated)
// then insert all crates and their tracks.
// - for every crate, get crate id either by searching for crate with that name or creating a new crate.
// - for each track, get track id from searching DB table `track_locations` with track path,
// - add entry to `crate_tracks` with crate id and track id
//
// Next TODO: what if library is already populated from previous use?
// - Don't delete all crates, to protect personal custom crates.
// - If a crate with the same name already exists:
//   - Wipe the crate.
//   - Add all tracks that are in the generated crate.

// UpdateCratesDB ...
func UpdateCratesDB(db mixxxdb.MixxxDB, crates []mixxxdb.CrateTracks) error {
	var errors *multierror.Error
	for _, crate := range crates {
		errors = multierror.Append(errors, UpdateCrateInDB(db, crate))
	}
	return errors.ErrorOrNil()
}

// UpdateCrateInDB ...
func UpdateCrateInDB(db mixxxdb.MixxxDB, crate mixxxdb.CrateTracks) error {
	// check if all the crate's tracks are already inserted in Mixxx's DB.
	// Simultaneously, get their IDs
	var errors *multierror.Error
	for _, track := range crate.Tracks {
		multierror.Append(errors, FindAndSetTrackID(db, &track))
	}
	if err := errors.ErrorOrNil(); err != nil {
		return err
	}

	// check whether the crate already exists
	dbcrate, err := db.Crates().FindByName(crate.Name)
	if err != nil {
		return fmt.Errorf("error finding crate %s: %w", dbcrate.Name, err)
	}

	// if so, then wipe the crate
	if dbcrate != nil {
		err := db.Crates().WipeTracks(dbcrate.ID)
		if err != nil {
			return fmt.Errorf("error wiping crate %d (%s): %w", dbcrate.ID, dbcrate.Name, err)
		}
	}

	// otherwise just insert the crate and its tracks.
	return db.Crates().InsertTracks(crate)
}

// FindAndSetTrackID searches the MixxxDB for the given track by its path and, if found,
// sets track.ID to the ID the track has in the DB.
// Returns an error wrapping `ErrTrackNotFound`
func FindAndSetTrackID(db mixxxdb.MixxxDB, track *mixxxdb.Track) error {
	dbtrack, err := db.Tracks().FindByPath(track.Path)
	if err != nil {
		return fmt.Errorf("error finding track with path %s in Mixxx DB: %w", track.Path, err)
	}

	if dbtrack == nil {
		return fmt.Errorf("%s: %w", track.Path, ErrTrackNotFound)
	}

	track.ID = dbtrack.ID
	return nil
}

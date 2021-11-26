package folders2crates

import (
	"errors"
	"fmt"
	"log"

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
func UpdateCratesDB(db mixxxdb.MixxxDB, crates []CrateFolder) error {
	var errors *multierror.Error
	for _, crate := range crates {
		errors = multierror.Append(errors, UpdateCrateInDB(db, crate))
		log.Println("Inserted crate: ", crate.Name)
	}

	return errors.ErrorOrNil()
}

// UpdateCrateInDB ...
func UpdateCrateInDB(db mixxxdb.MixxxDB, crate CrateFolder) error {
	// find all the crate's tracks in Mixxx's DB.
	trackIDs, err := FindTrackIDs(db, crate.Tracks)
	if err != nil {
		return err
	}

	// find the crate in Mixxx's DB
	dbcrate, err := db.Crates().FindByName(crate.Name)
	if err != nil {
		return fmt.Errorf("error finding crate '%s': %w", dbcrate.Name, err)
	}

	// if it doesn't exist yet, create it
	if dbcrate == nil {
		dbcrate, err = db.Crates().InsertCrate(mixxxdb.Crate{Name: crate.Name})
		if err != nil {
			return fmt.Errorf("error inserting crate '%s': %w", crate.Name, err)
		}
	}

	// if it does, then remove all tracks in it
	if dbcrate != nil {
		err := db.CrateTracks().WipeCrate(dbcrate.ID)
		if err != nil {
			return fmt.Errorf("error wiping crate %d (%s): %w", dbcrate.ID, dbcrate.Name, err)
		}
	}

	// Now with an empty crate, we insert all tracks we want to have in there.
	err = db.CrateTracks().InsertTracks(dbcrate.ID, trackIDs)
	if err != nil {
		return fmt.Errorf("error inserting tracks from crate '%s': %w", crate.Name, err)
	}
	return nil
}

// FindTrackIDs searches the DB to find the track IDs corresponding to the given track files.
// Returns an error wrapping one or multiple `ErrTrackNotFound`s if there are tracks that cannot be found in Mixxx's DB.
func FindTrackIDs(db mixxxdb.MixxxDB, tracks []TrackFile) ([]int, error) {
	var errors *multierror.Error
	trackIDs := make([]int, 0, len(tracks))

	for _, fpath := range tracks {
		trackLoc, err := db.TrackLocations().FindByLocation(string(fpath))
		if trackLoc == nil {
			err = fmt.Errorf("%s: %w", fpath, ErrTrackNotFound)
		}
		if err != nil {
			errors = multierror.Append(errors, err)
			continue
		}

		track, err := db.Tracks().FindByLocationID(trackLoc.ID)
		if track == nil {
			err = fmt.Errorf("%s: %w", fpath, ErrTrackNotFound)
		}
		if err != nil {
			errors = multierror.Append(errors, err)
			continue
		}

		trackIDs = append(trackIDs, track.ID)
	}

	return trackIDs, errors.ErrorOrNil()
}

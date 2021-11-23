package mixxxdb

type Crate struct {
	// The crate's ID in Mixxx's DB. Can be found in `crates` table
	ID uint
	// The name of the crate.
	Name string

	Tracks []Track
}

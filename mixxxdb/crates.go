package mixxxdb

type Crate struct {
	// The crate's ID in Mixxx's DB. Can be found in `crates` table
	ID int64
	// The name of the crate.
	Name string
}

type CrateTracks struct {
	Crate
	Tracks []Track
}

func NewCrate(id int64, name string) Crate {
	return Crate{ID: id, Name: name}
}

func NewCrateTracks(id int64, name string, tracks []Track) CrateTracks {
	return CrateTracks{Crate: Crate{ID: id, Name: name}, Tracks: tracks}
}

package mixxxdb

// Track represents a music file that can be played by Mixxx.
type Track struct {
	// The track's ID in Mixxx DB. Can be found in `library` table or in `track_locations` table.
	ID uint

	// Path to the music file. Stored in `track_locations` table.
	Path string
}

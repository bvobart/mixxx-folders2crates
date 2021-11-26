package folders2crates

import (
	"os"
	"path"

	ignore "github.com/sabhiram/go-gitignore"
)

// TrackFile is a filepath to a music file that can be played by Mixxx.
// See `IsTrackFile` to see what types of files qualify.
type TrackFile string

// IsTrackFile detects if a file at the given path is indeed a track that can be played by Mixxx.
// Returns true iff the filename has an extension that is supported by Mixxx.
// See https://manual.mixxx.org/2.0/en/chapters/configuration.html#importing-your-audio-files for the list of supported files
func IsTrackFile(filename string) bool {
	switch path.Ext(filename) {
	case ".wav", ".aiff", ".aif", ".mp3", ".ogg", ".flac", ".opus":
		return true
	}
	return false
}

// FindTrackFiles finds all tracks that can be played by Mixxx in a given folder.
// Respects the ignore patterns specified with the github.com/sabhiram/go-gitignore library.
// Note: Does *not* search recursively.
func FindTrackFiles(folder string, ignore *ignore.GitIgnore) ([]TrackFile, error) {
	files, err := os.ReadDir(folder)
	if err != nil {
		return nil, err
	}

	tracks := []TrackFile{}
	for _, file := range files {
		// if file is directory or if it's not a music file, skip.
		if file.IsDir() || !IsTrackFile(file.Name()) {
			continue
		}

		// ignore .crateignore'd paths
		fullpath := path.Join(folder, file.Name())
		if ignore != nil && ignore.MatchesPath(fullpath) {
			continue
		}

		tracks = append(tracks, TrackFile(fullpath))
	}

	return tracks, nil
}

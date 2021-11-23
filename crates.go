package folders2crates

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/bvobart/mixxx-folders2crates/mixxxdb"
	ignore "github.com/sabhiram/go-gitignore"
)

// FindCrates searches the given folder for music tracks that Mixxx can play.
// When it finds a folder with at least track in it, it will make a crate with the name of the path to that folder
// and the tracks that are directly inside the folder.
// Respects the ignore patterns specified with the github.com/sabhiram/go-gitignore library.
func FindCrates(libfolder string, ignore *ignore.GitIgnore) ([]mixxxdb.Crate, error) {
	crates := []mixxxdb.Crate{}

	err := filepath.WalkDir(libfolder, func(fpath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// ignore the library root directory and anything that's not a directory.
		if !d.IsDir() || fpath == libfolder {
			return nil
		}

		// ignore folders from .crateignore
		if ignore.MatchesPath(fpath) {
			return nil
		}

		// when encountering a non-ignored directory, find all tracks in it and make a crate
		relpath := strings.TrimPrefix(fpath, libfolder+"/")
		tracks, err := FindTracks(fpath, ignore)
		if err != nil {
			return fmt.Errorf("failed to find tracks in %s: %w", relpath, err)
		}

		// skip folders that don't contain any tracks.
		if len(tracks) == 0 {
			return nil
		}

		name := NameCrate(relpath)
		crates = append(crates, mixxxdb.Crate{
			ID:     0, // unknown at this point
			Name:   name,
			Tracks: tracks,
		})
		return nil
	})

	return crates, err
}

// NameCrate creates a name for a crate based on its path relative to the music library.
// Replaces folder separators with -
// E.g. House/90's becomes House - 90's
func NameCrate(relpath string) string {
	return strings.Replace(relpath, string(os.PathSeparator), " - ", -1)
}

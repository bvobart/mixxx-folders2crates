package folders2crates

import (
	"io/fs"
	"path"
	"path/filepath"
	"strings"

	"github.com/bvobart/mixxx-folders2crates/utils"
	ignore "github.com/sabhiram/go-gitignore"
)

type Crate struct {
	Name   string
	Tracks []Track
}

type Track struct {
	Path string
	Size int64
}

func DetectCrates(libfolder string, ignore *ignore.GitIgnore) ([]Crate, error) {
	crates := []Crate{}
	cratesByPath := map[string]*Crate{}

	err := filepath.WalkDir(libfolder, func(fpath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// ignore any files and folders from .crateignore
		if ignore.MatchesPath(fpath) {
			return nil
		}

		// when encountering a directory, make a crate
		// TODO; only create crate if there are more than 0 music files in the folder.
		if d.IsDir() {
			// skip the library root.
			if fpath == libfolder {
				return nil
			}

			relpath := strings.TrimPrefix(fpath, libfolder+"/")
			crates = append(crates, Crate{Name: relpath, Tracks: []Track{}})
			cratesByPath[fpath] = &crates[len(crates)-1]
			return nil
		}

		// skip any non-music files
		if !utils.IsMusicFile(fpath) {
			return nil
		}

		// when encountering a music file, add it to the crate it belongs to (i.e. its parent folder)
		fInfo, err := d.Info()
		if err != nil {
			return err
		}
		track := Track{Path: fpath, Size: fInfo.Size()}
		crate := cratesByPath[path.Dir(fpath)]
		crate.Tracks = append(crate.Tracks, track)

		return nil
	})

	return crates, err
}

package folders2crates_test

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/bvobart/mixxx-folders2crates/folders2crates"
	"github.com/stretchr/testify/require"
)

func TestFindCrateFolders(t *testing.T) {
	t.Run("Finds some crates in fabricated folder", func(t *testing.T) {
		// Given:
		dir, err := os.MkdirTemp("", "test-library")
		require.NoError(t, err)

		expectedCrates := []folders2crates.CrateFolder{
			{Name: "House", Tracks: []folders2crates.TrackFile{}},
			{Name: "Techno", Tracks: []folders2crates.TrackFile{}},
			// TODO: add subdirectories to test recursive searching behaviour
		}

		// ... a House folder with a load of .mp3 files in it
		houseDir := path.Join(dir, "House")
		require.NoError(t, os.Mkdir(houseDir, 0755))
		for i := 0; i < 10; i++ {
			fpath := path.Join(houseDir, fmt.Sprint("song-", i, ".mp3"))
			f, err := os.Create(fpath)
			require.NoError(t, err)
			defer f.Close()

			expectedCrates[0].Tracks = append(expectedCrates[0].Tracks, folders2crates.TrackFile(fpath))
		}

		// ... a Techno folder with a load of .flac files in it and a load of .cue files (which should be ignored)
		technoDir := path.Join(dir, "Techno")
		require.NoError(t, os.Mkdir(technoDir, 0755))
		for i := 0; i < 10; i++ {
			fpath := path.Join(technoDir, fmt.Sprint("song-", i, ".flac"))
			f, err := os.Create(fpath)
			require.NoError(t, err)
			defer f.Close()

			expectedCrates[1].Tracks = append(expectedCrates[1].Tracks, folders2crates.TrackFile(fpath))

			f, err = os.Create(path.Join(houseDir, fmt.Sprint("cue-", i, ".cue")))
			require.NoError(t, err)
			defer f.Close()
		}

		// When:
		crates, err := folders2crates.FindCrateFolders(dir, nil)

		// Then:
		require.NoError(t, err)
		require.Equal(t, expectedCrates, crates)
	})
}

package folders2crates_test

import (
	"os"
	"path"
	"runtime"
	"testing"

	folders2crates "github.com/bvobart/mixxx-folders2crates"
	"github.com/stretchr/testify/require"
)

func TestFindTrackFiles(t *testing.T) {
	t.Run("Finds none in source code directory", func(t *testing.T) {
		_, filename, _, ok := runtime.Caller(0)
		require.True(t, ok)

		dir := path.Join(path.Dir(filename), "mixxxdb") // mixxxdb folder in this repo.
		tracks, err := folders2crates.FindTrackFiles(dir, nil)

		require.NoError(t, err)
		require.Equal(t, []folders2crates.TrackFile{}, tracks)
	})

	t.Run("Finds all Mixxx's supported files in fabricated dir", func(t *testing.T) {
		// Given:
		trackNames := []string{"song.aif", "song.aiff", "song.flac", "song.mp3", "song.ogg", "song.opus", "song.wav"}
		dir, err := os.MkdirTemp("", "test-library")
		require.NoError(t, err)
		defer os.RemoveAll(dir) // cleanup

		expectedTracks := []folders2crates.TrackFile{}
		for _, name := range trackNames {
			fpath := path.Join(dir, name)
			expectedTracks = append(expectedTracks, folders2crates.TrackFile(fpath))

			f, err := os.Create(fpath)
			require.NoError(t, err)
			defer f.Close()
		}

		// When:
		tracks, err := folders2crates.FindTrackFiles(dir, nil)

		// Then:
		require.NoError(t, err)
		require.Equal(t, expectedTracks, tracks)
	})
}

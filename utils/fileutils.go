package utils

import (
	"errors"
	"io"
	"os"
	"path"
)

// FileExists checks if a file exists and is not a directory
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	return !os.IsNotExist(err) && info != nil && !info.IsDir()
}

// FolderExists checks if a folder exists and is indeed a folder.
func FolderExists(filename string) bool {
	info, err := os.Stat(filename)
	return !os.IsNotExist(err) && info != nil && info.IsDir()
}

// FolderIsEmpty checks if a folder is empty
func FolderIsEmpty(filename string) (bool, error) {
	file, err := os.Open(filename)
	if err != nil {
		return false, err
	}
	defer file.Close()

	_, err = file.Readdirnames(1)
	if errors.Is(err, io.EOF) {
		return true, nil
	}
	return false, err
}

// Returns the user's home directory. On Linux and MacOS, this is equivalent to resolving `~`
func HomeDir() string {
	if dir := os.Getenv("HOME"); dir != "" {
		return dir
	}
	return os.Getenv("USERPROFILE") // windows
}

// Returns true iff the filename has an extension that is supported by Mixxx
// See https://manual.mixxx.org/2.0/en/chapters/configuration.html#importing-your-audio-files for the list of supported files
func IsMusicFile(filename string) bool {
	switch path.Ext(filename) {
	case ".wav", ".aiff", ".aif", ".mp3", ".ogg", ".flac", ".opus":
		return true
	}
	return false
}

package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/fatih/color"
	ignore "github.com/sabhiram/go-gitignore"

	folders2crates "github.com/bvobart/mixxx-folders2crates"
	"github.com/bvobart/mixxx-folders2crates/utils"
)

var mixxxDB = path.Join(utils.HomeDir(), ".mixxx/mixxxdb.sqlite") // default Mixxx SQLite DB location under Linux

func main() {
	startTime := time.Now()
	libfolder := parseArgs(os.Args)

	color.Green("Mixxx DB:      %s", color.HiWhiteString(mixxxDB))
	color.Green("Music Library: %s", color.HiWhiteString(libfolder))
	fmt.Println()

	// parse .crateignore
	ignoreFile, err := ignore.CompileIgnoreFile(path.Join(libfolder, ".crateignore"))
	if err != nil {
		ignoreFile = &ignore.GitIgnore{}
	}

	// detect which folders in the music library should become crates and what tracks should be in them according to the folder layout.
	crates, err := folders2crates.FindCrates(libfolder, ignoreFile)
	if err != nil {
		color.Red("cannot detect crates from your music library: %s", color.YellowString(err.Error()))
		os.Exit(5)
	}

	// temporary: print all crates
	bold := color.New(color.Bold)
	faint := color.New(color.Faint)
	for _, crate := range crates {
		if len(crate.Tracks) == 0 {
			continue
		}

		bold.Print(crate.Name, " (", len(crate.Tracks), " tracks)\n")
		for _, track := range crate.Tracks {
			fmt.Print("- ")
			faint.Print(path.Dir(track.Path), "/")
			fmt.Println(path.Base(track.Path))
		}
		fmt.Println()
	}

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

	fmt.Println("took:", time.Since(startTime))
}

// parseArgs parses the arguments passed to folders2crates, deals with invalid arguments and returns the one valid argument: the path to a music library folder
func parseArgs(args []string) string {
	if len(args) < 2 {
		color.Red("expecting a music library folder as argument, but nothing was provided")
		os.Exit(1)
	}
	if len(args) > 2 {
		color.Yellow("WARNING: provided multiple arguments, only the first one will be used: %s. Ignored: %s", args[1], args[2:])
	}

	libfolder := args[1]
	if !utils.FolderExists(libfolder) {
		if utils.FileExists(libfolder) {
			color.Red("expecting a music library folder, but got a file: %s", libfolder)
			os.Exit(3)
		}

		color.Red("music library folder doesn't exist: %s", libfolder)
		os.Exit(2)
	}

	if !utils.FileExists(mixxxDB) {
		color.Red("cannot open your Mixxx DB, because it does not exist.")
		color.Yellow("Try starting Mixxx, then closing it, then running this program again.")
		os.Exit(4)
	}

	libfolder, _ = filepath.Abs(libfolder)
	return libfolder
}

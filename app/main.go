package main

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/fatih/color"
	ignore "github.com/sabhiram/go-gitignore"

	folders2crates "github.com/bvobart/mixxx-folders2crates"
	"github.com/bvobart/mixxx-folders2crates/mixxxdb"
	"github.com/bvobart/mixxx-folders2crates/utils"
)

func main() {
	startTime := time.Now()
	libfolder := parseArgs(os.Args)

	color.Green("Mixxx DB:      %s", color.HiWhiteString(mixxxdb.DefaultMixxxDBPath))
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
		color.Red("Error detecting crates from your music library:")
		color.Red("  %s", color.YellowString(err.Error()))
		os.Exit(5)
	}

	// open Mixxx's database
	db, err := mixxxdb.OpenDefault()
	if err != nil {
		color.Red("Error opening Mixxx's DB:")
		color.Red("  %s", color.YellowString(err.Error()))
		os.Exit(6)
	}

	// temporary: print all crates
	bold := color.New(color.Bold)
	faint := color.New(color.Faint)
	green := color.New(color.FgGreen)
	for _, crate := range crates {
		if len(crate.Tracks) == 0 {
			continue
		}

		bold.Print(crate.Name, " (", len(crate.Tracks), " tracks)")
		dbCrate, err := db.Crates().FindByName(crate.Name)
		if err != nil {
			panic(err)
		}

		if dbCrate != nil {
			green.Println(" - exists!")
		} else {
			fmt.Println()
		}

		for _, track := range crate.Tracks {
			fmt.Print("- ")
			faint.Print(path.Dir(track.Path), "/")
			fmt.Println(path.Base(track.Path))
		}
		fmt.Println()
	}

	if utils.IsInteractive() {
		pauseTime := time.Now()
		if err := utils.PromptConfirm("Are you sure you want these crates and tracks in Mixxx's DB?"); err != nil {
			return
		}
		startTime = startTime.Add(time.Since(pauseTime)) // ignores the time taken to confirm
		fmt.Println()
	}

	err = folders2crates.UpdateCratesDB(db, crates)
	if errors.Is(err, folders2crates.ErrTrackNotFound) || errors.Is(err, mixxxdb.ErrTrackNoID) {
		color.Yellow("Warning!")
		color.Yellow("  There were one or more tracks that couldn't be found in Mixxx's library.")
		color.Yellow("  The easiest way to fix this problem is to start Mixxx, re-scan your library, then close Mixxx and run this program again.")
		fmt.Println()
	}
	if err != nil {
		color.Red("Error:")
		color.Red("  %s", color.YellowString(err.Error()))
		os.Exit(7)
	}

	color.Green("Done!")
	faint.Println("took:", time.Since(startTime))
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

	if !utils.FileExists(mixxxdb.DefaultMixxxDBPath) {
		color.Red("cannot open your Mixxx DB, because it does not exist.")
		color.Yellow("Try starting Mixxx, then closing it, then running this program again.")
		os.Exit(4)
	}

	libfolder, _ = filepath.Abs(libfolder)
	return libfolder
}

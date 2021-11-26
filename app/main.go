package main

import (
	"errors"
	"fmt"
	"math"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/hashicorp/go-multierror"
	ignore "github.com/sabhiram/go-gitignore"

	folders2crates "github.com/bvobart/mixxx-folders2crates"
	"github.com/bvobart/mixxx-folders2crates/mixxxdb"
	"github.com/bvobart/mixxx-folders2crates/utils"
)

var bold = color.New(color.Bold)
var faint = color.New(color.Faint)
var green = color.New(color.FgGreen)
var red = color.New(color.FgRed)
var yellow = color.New(color.FgYellow)

func main() {
	startTime := time.Now()
	defer func() { faint.Println("took:", time.Since(startTime)) }()
	libfolder := parseArgs(os.Args)

	green.Println("Mixxx DB:     ", color.HiWhiteString(mixxxdb.DefaultMixxxDBPath))
	green.Println("Music Library:", color.HiWhiteString(libfolder))
	fmt.Println()

	// parse .crateignore
	ignoreFile, err := ignore.CompileIgnoreFile(path.Join(libfolder, ".crateignore"))
	if err != nil {
		ignoreFile = &ignore.GitIgnore{}
	}

	// detect which folders in the music library should become crates and what tracks should be in them according to the folder layout.
	crates, err := folders2crates.FindCrateFolders(libfolder, ignoreFile)
	if err != nil {
		red.Println("Error detecting crates from your music library:")
		red.Println("  ", yellow.Sprint(err.Error()))
		os.Exit(5)
	}

	if len(crates) == 0 {
		green.Println("No crates found, nothing to do!")
		return
	}

	// open Mixxx's database
	db, err := mixxxdb.OpenDefault()
	if err != nil {
		red.Println("Error opening Mixxx's DB:")
		red.Println("  ", yellow.Sprint(err.Error()))
		os.Exit(6)
	}

	// temporary: print all crates
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

		for i, track := range crate.Tracks {
			if i < len(crate.Tracks)-1 {
				fmt.Print("├── ")
			} else {
				fmt.Print("└── ")
			}
			faint.Print(".", strings.TrimPrefix(path.Dir(string(track)), libfolder), "/")
			fmt.Println(path.Base(string(track)))
		}
		fmt.Println()
	}

	if utils.IsInteractive() {
		pauseTime := time.Now()
		if err := utils.PromptConfirm("Are you sure you want these crates and tracks in Mixxx's DB?"); err != nil {
			fmt.Println()
			yellow.Println("Alright, no problem, just let me know when you need those crates inserted! 😊")
			return
		}
		startTime = startTime.Add(time.Since(pauseTime)) // ignores the time taken to confirm
		fmt.Println()
	}

	yellow.Println("Inserting ", bold.Sprint(len(crates)), yellow.Sprint(" crates into Mixxx's DB..."))
	fmt.Println()

	// update the crates in Mixxx' DB
	var multierr *multierror.Error
	for _, crate := range crates {
		t := time.Now()
		err := folders2crates.UpdateCrateInDB(db, crate)
		multierr = multierror.Append(multierr, err)
		if err == nil {
			green.Print("✅ ", crate.Name, strings.Repeat(" ", int(math.Max(0, float64(48-len(crate.Name))))))
		} else {
			red.Print("❌ ", crate.Name, strings.Repeat(" ", int(math.Max(0, float64(48-len(crate.Name))))))
		}
		faint.Println("\t", time.Since(t))
	}

	fmt.Println()

	err = multierr.ErrorOrNil()
	if errors.Is(err, folders2crates.ErrTrackNotFound) {
		yellow.Println("Warning!")
		yellow.Println("  There were one or more tracks that couldn't be found in Mixxx's library.")
		yellow.Println("  The easiest way to fix this problem is to start Mixxx, let it re-scan your library, then close Mixxx and run this program again.")
		fmt.Println()
	}
	if err != nil {
		red.Println("Error:")
		red.Println("  ", yellow.Sprint(err.Error()))
		os.Exit(7)
	}

	green.Println("Done!")
}

// parseArgs parses the arguments passed to folders2crates, deals with invalid arguments and returns the one valid argument: the path to a music library folder
func parseArgs(args []string) string {
	if len(args) < 2 {
		red.Println("expecting a music library folder as argument, but nothing was provided")
		os.Exit(1)
	}
	if len(args) > 2 {
		yellow.Println("WARNING: provided multiple arguments, only the first one will be used:", args[1])
		yellow.Println("WARNING: arguments ignored:", args[2:])
	}

	libfolder := args[1]
	if !utils.FolderExists(libfolder) {
		if utils.FileExists(libfolder) {
			red.Println("expecting a music library folder, but got a file: ", libfolder)
			os.Exit(3)
		}

		red.Println("music library folder doesn't exist: ", libfolder)
		os.Exit(2)
	}

	if !utils.FileExists(mixxxdb.DefaultMixxxDBPath) {
		red.Println("cannot open your Mixxx DB, because it does not exist.")
		yellow.Println("Try starting Mixxx, then closing it, then running this program again.")
		os.Exit(4)
	}

	libfolder, _ = filepath.Abs(libfolder)
	return libfolder
}

# `mixxx-folders2crates`

<p align="center">
  <img alt="GitHub Workflow Status" src="https://img.shields.io/github/workflow/status/bvobart/mixxx-folders2crates/Release">
  <img alt="GitHub go.mod Go version" src="https://img.shields.io/github/go-mod/go-version/bvobart/mixxx-folders2crates">
  <a href="https://pkg.go.dev/github.com/bvobart/mixxx-folders2crates"><img src="https://pkg.go.dev/badge/github.com/bvobart/mixxx-folders2crates.svg" alt="Go Reference"></a>
</p>

This is the repository for a tool I developed to help keep my crates up to date with my music library. My library consists of loads of songs that I manually sorted into a bunch of folders per genre and subgenre of that music, e.g. something like:

```
├── House
│   ├── 90s
│   │   ├── Funky
│   │   │   ├── Fatboy Slim & Riva Starr - Eat, Sleep, Rave, Repeat.mp3
│   │   │   ├── Felix - Don't U Want Me.flac
│   │   │   ├── Mr. Matey - Acid Party (Aciieed Groove Mix).flac
│   │   │   ├── ...
│   │   │   └── Xpansions - Move Your Body (Elevation) [Club Mix].mp3
│   │   └── Mr Oizo - Flat Beat.mp3
│   ├── Acid
│   │   ├── Acid Arab - Le Gaz Qui Fait Rire.mp3
│   │   └── Harder
│   │       ├── Hot Hanas Hula - 70th & King Drive.flac
│   │       ├── Hot Hanas Hula - Hot Hands.flac
│   │       ├── Jack Frost & The Circle Jerks - Acid Rout.flac
│   │       └── Mix Machine - Bikini.flac
│   ├── Bangers
│   │   ├── BICEP - VISION OF LOVE.mp3
│   │   ├── David Harness - Al Greenz (Dance) (Original Mix).mp3
│   │   ├── Disclosure - Ecstasy.mp3
│   │   ├── DMX Krew - Come to Me.flac
...
```

My problem is that I want to import these songs into Mixxx crates based on this folder structure.

Doing this manually requires opening Mixxx, then:
- for every one of those (sub-)folders in my library: 
  - navigate to this folder
  - create a crate for the folder if it doesn't already exist, e.g. `House - 90s - Acid`
  - select every song in it and right click -> Add to crate -> select the crate for this folder -> tick the checkbox.

That's quite some manual effort which compounds greatly when adding a few songs to my library every now and then, or moving songs around within my library (making them incorrectly be shown in the crate while they no longer exist).

So I decided to spend my manual effort into building this tool to automate the process. 
It's written in Go and takes about 1 second to update all 46 crates and 832 tracks in my music library in Mixxx's SQLite DB file.

This tool is tested and used by me on Linux (x86 and ARM), but you may also be able to get it working on Windows or MacOS. I included the default Mixxx DB file paths for Linux, Windows and MacOS, so if you can get it to compile it on your machine, this tool should theoretically work. Please let me know if it actually does :)

## Usage

With the latest version of Go (1.17 or higher) installed, `git clone` this repo, then open a terminal in this folder and run:

```sh
go run ./app FOLDER
# e.g.
go run ./app ~/Music
```

replacing `FOLDER` with the location of your music library.

> Note: this tool depends on [`go-sqlite3`](https://pkg.go.dev/github.com/mattn/go-sqlite3) which is a `CGO_ENABLED=1` package and thus requires a C compiler to be installed for it to compile.

### `.crateignore`

Place a `.crateignore` file at the root of your music library folder to specify files and folders that `mixxx-folders2crates` should ignore.
This file uses the same syntax as a `.gitignore` file (thanks to [`go-gitignore`](https://pkg.go.dev/github.com/sabhiram/go-gitignore)).

For example:

```sh
.stfolder
Phone
# this is just a comment
MixxxRecordings
Beatport Top 100*
```

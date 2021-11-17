# Mixxx - folders2crates

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
│   ├── Afro Exotic
│   │   ├── Afriquoi - Acid Attack.mp3
│   │   ├── Afriquoi - Bayeke.mp3
│   │   ├── Afriquoi - Kudaushe.mp3
│   │   ├── Afriquoi - Ndeko Solo.mp3
│   │   ├── Afriquoi - Ndeko Solo (Voilaaa Remix).mp3
│   │   ├── Afriquoi - Sam Sam.mp3
│   │   ├── Disclosure & Eko Roosevelt - Tondo.aif
│   │   ├── Disclosure, FATOUMATA DIAWARA - Douha (Mali Mali).flac
│   │   ├── Nickodemus - Ndini (feat ismael kouyate).mp3
│   │   └── Raoul K, Manoo feat. Ahmed Sosso - Toukan (Dixon Rework) (Dixon Rework).mp3
│   ├── Bangers
│   │   ├── BICEP - VISION OF LOVE.mp3
│   │   ├── David Harness - Al Greenz (Dance) (Original Mix).mp3
│   │   ├── Disclosure - Ecstasy.mp3
│   │   ├── DMX Krew - Come to Me.flac
...
```

My problem is that I want to import these songs into Mixxx crates based on this folder structure. Doing this manually requires opening Mixxx, then for every leaf folder in my library: navigate to this folder, select every song in it and right click to add to the crate (creating the crate first if it doesn't exist yet). That's quite some manual effort which compounds greatly when adding a few songs to your library every now and then, or moving songs around within my library (making them incorrectly be shown in the crate while they no longer exist).

This tool aims to help overcome that problem by:
- [ ] importing all the songs in the leaf folders in my library into crates
  - [ ] creating new crates if they do not yet exist
  - [ ] if a song is already in your Mixxx library, but in a different folder than it is now, move the song while trying to keep the analysis.
- [ ] Read from a `.crateignore` file (with syntax like `.gitignore` files) in the library folder to determine which files / folders to ignore for adding to crates.
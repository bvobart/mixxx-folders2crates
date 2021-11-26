package mixxxdb

import (
	"errors"

	"github.com/upper/db/v4"
)

// Compile-time checks on conformance to interfaces.
var _ = db.Record(&Crate{})
var _ = db.Store(&cratesDB{})
var _ = CratesDB(&cratesDB{})

type Crate struct {
	// The crate's ID in Mixxx's DB. Can be found in `crates` table
	ID int `db:"id,omitempty"`

	// The name of the crate.
	Name string `db:"name"`
}

func (_ *Crate) Store(sess db.Session) db.Store {
	return NewCratesDB(sess)
}

//-------------------------------------------------------------------------------------------------------------

func NewCratesDB(sess db.Session) CratesDB {
	return &cratesDB{sess.Collection("crates")}
}

type cratesDB struct {
	db.Collection
}

func (crates *cratesDB) FindByName(name string) (*Crate, error) {
	var crate Crate
	err := crates.Find(db.Cond{"name": name}).One(&crate)
	if errors.Is(err, db.ErrNoMoreRows) {
		return nil, nil
	}
	return &crate, err
}

func (crates *cratesDB) InsertCrate(crate Crate) (*Crate, error) {
	return &crate, crates.InsertReturning(&crate)
}

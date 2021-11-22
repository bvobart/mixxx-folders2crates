package mixxxdb

import (
	"path"

	"github.com/bvobart/mixxx-folders2crates/utils"
)

// Default Mixxx SQLite DB location under Linux / Unix systems.
var DefaultMixxxDBPath = path.Join(utils.HomeDir(), ".mixxx", "mixxxdb.sqlite")

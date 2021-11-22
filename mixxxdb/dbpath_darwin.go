package mixxxdb

import (
	"path"

	"github.com/bvobart/mixxx-folders2crates/utils"
)

// Default Mixxx SQLite DB location under MacOS
var DefaultMixxxDBPath = path.Join(utils.HomeDir(), "Library", "Application Support", "Mixxx", "mixxxdb.sqlite")

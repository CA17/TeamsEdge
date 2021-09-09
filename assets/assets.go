package assets

import (
	_ "embed"
)

//go:embed build.txt
var BuildInfo string

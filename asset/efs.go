package asset

import (
	"embed"
)

//go:embed migrations config.yaml
var EmbeddedFiles embed.FS

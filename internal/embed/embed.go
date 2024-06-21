package embed

import (
	"embed"
	"io/fs"
)

//go:embed static
var static embed.FS

// GetStatic returns the static files
func GetStatic() fs.FS {
	sfs, _ := fs.Sub(static, "static")
	return sfs
}

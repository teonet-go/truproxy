//go:build !dev

package main

import (
	"embed"
	"io/fs"
)

//go:embed frontend/dist
var embedFrontend embed.FS

func getFrontendAssets() fs.FS {
	f, err := fs.Sub(embedFrontend, "frontend/dist")
	if err != nil {
		panic(err)
	}

	return f
}

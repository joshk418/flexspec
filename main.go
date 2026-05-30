/*
Copyright © 2026 Josh Kyte
*/
package main

import (
	"embed"
	"io/fs"
	"log"

	"github.com/joshk418/flexspec/cmd"
)

// templatesFS embeds the template files shipped with the binary so that
// `flexspec init` can scaffold them into a project's .flexspec/templates dir.
//
//go:embed all:templates
var templatesFS embed.FS

// uiFS embeds the production web UI (ui/dist). Run `make build-ui` before release builds.
//
//go:embed all:ui/dist
var uiFS embed.FS

func main() {
	cmd.TemplatesFS = templatesFS

	// The embed paths are rooted at ui/dist; serve the UI from that subtree so
	// index.html and assets resolve at the server root.
	uiRoot, err := fs.Sub(uiFS, "ui/dist")
	if err != nil {
		log.Fatalf("mount embedded web UI: %v", err)
	}
	cmd.UIFS = uiRoot

	cmd.Execute()
}

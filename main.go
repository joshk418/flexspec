package main

import (
	"embed"
	"io/fs"
	"log"

	"github.com/joshk418/flexspec/cmd"
)

// templatesFS embeds templates for `flexspec init`.
//
//go:embed all:templates
var templatesFS embed.FS

// skillsFS embeds agent skills for `flexspec update --skills`.
//
//go:embed all:skills
var skillsFS embed.FS

// uiFS embeds the production web UI (ui/dist). Run `make build-ui` before release builds.
//
//go:embed all:ui/dist
var uiFS embed.FS

func main() {
	cmd.TemplatesFS = templatesFS

	// Mount skills/ so the installer walks from ".".
	skillsRoot, err := fs.Sub(skillsFS, "skills")
	if err != nil {
		log.Fatalf("mount embedded skills: %v", err)
	}
	cmd.SkillsFS = skillsRoot

	// Mount ui/dist so index.html and assets resolve at the server root.
	uiRoot, err := fs.Sub(uiFS, "ui/dist")
	if err != nil {
		log.Fatalf("mount embedded web UI: %v", err)
	}
	cmd.UIFS = uiRoot

	cmd.Execute()
}

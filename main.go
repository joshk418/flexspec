/*
Copyright © 2026 Josh Kyte
*/
package main

import (
	"embed"

	"github.com/joshk418/flexspec/cmd"
)

// templatesFS embeds the template files shipped with the binary so that
// `flexspec init` can scaffold them into a project's .flexspec/templates dir.
//
//go:embed all:templates
var templatesFS embed.FS

func main() {
	cmd.TemplatesFS = templatesFS
	cmd.Execute()
}

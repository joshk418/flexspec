package cmd

import (
	"os"
	"testing"

	"github.com/joshk418/flexspec/internal/ui"
)

func TestRunUI_missingConfig(t *testing.T) {
	root := t.TempDir()
	oldWd, _ := os.Getwd()
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(oldWd) }()

	UIFS = ui.StubStaticFS()
	uiNoOpen = true
	defer func() { uiNoOpen = false }()

	err := runUI(uiCmd, nil)
	if err == nil {
		t.Fatal("expected error when .flexspec/config.yaml is missing")
	}
}

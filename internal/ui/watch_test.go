package ui

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/joshk418/flexspec/internal/spec"
)

func TestWatcher_broadcastsOnSpecChange(t *testing.T) {
	root := t.TempDir()
	writeUIProject(t, root)

	srv, err := NewServer(root, "127.0.0.1", 0, StubStaticFS())
	if err != nil {
		t.Fatal(err)
	}
	if err := srv.startWatcher(); err != nil {
		t.Fatal(err)
	}
	defer func() {
		close(srv.done)
		if srv.watcher != nil {
			_ = srv.watcher.Close()
		}
	}()

	ch := srv.hub.Subscribe()
	defer srv.hub.Unsubscribe(ch)

	readme := filepath.Join(root, "specs", "001-test", "README.md")
	if err := spec.SetFileStatus(readme, "in_progress"); err != nil {
		t.Fatal(err)
	}

	select {
	case ev := <-ch:
		if ev != "specs-changed" {
			t.Fatalf("event = %q", ev)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for specs-changed after file write")
	}
}

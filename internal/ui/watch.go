package ui

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

func (s *Server) startWatcher() error {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	s.watcher = w

	specsPath := filepath.Join(s.root, s.cfg.SpecsDir)
	_ = filepath.Walk(specsPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			return w.Add(path)
		}
		return nil
	})
	if err := w.Add(specsPath); err != nil && !os.IsNotExist(err) {
		_ = w.Close()
		return err
	}
	configPath := filepath.Join(s.root, ".flexspec", "config.yaml")
	if err := w.Add(filepath.Dir(configPath)); err != nil {
		_ = w.Close()
		return err
	}

	go s.runWatcher(w)
	return nil
}

func (s *Server) runWatcher(w *fsnotify.Watcher) {
	var debounce *time.Timer
	for {
		select {
		case <-s.done:
			return
		case ev, ok := <-w.Events:
			if !ok {
				return
			}
			if !s.watchPathInteresting(ev.Name) {
				continue
			}
			if debounce != nil {
				debounce.Stop()
			}
			debounce = time.AfterFunc(500*time.Millisecond, func() {
				s.hub.Broadcast("specs-changed")
			})
		case _, ok := <-w.Errors:
			if !ok {
				return
			}
		}
	}
}

func (s *Server) watchPathInteresting(path string) bool {
	path = filepath.Clean(path)
	if strings.HasSuffix(path, "config.yaml") {
		return true
	}
	if !strings.HasSuffix(path, ".md") {
		return false
	}
	specsRoot := filepath.Join(s.root, s.cfg.SpecsDir)
	rel, err := filepath.Rel(specsRoot, path)
	if err != nil || strings.HasPrefix(rel, "..") {
		return false
	}
	return true
}

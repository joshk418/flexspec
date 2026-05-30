package ui

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"

	"github.com/joshk418/flexspec/internal/config"
)

// Server serves the FlexSpec management UI and JSON API.
type Server struct {
	root    string
	cfg     config.Config
	addr    string
	static  fs.FS
	hub     *EventHub
	watcher *fsnotify.Watcher
	done    chan struct{}
	http    *http.Server
}

// NewServer creates a UI server for a FlexSpec project root.
func NewServer(root, host string, port int, static fs.FS) (*Server, error) {
	cfg, err := config.Load(root)
	if err != nil {
		return nil, err
	}
	addr := fmt.Sprintf("%s:%d", host, port)
	s := &Server{
		root:   root,
		cfg:    cfg,
		addr:   addr,
		static: static,
		hub:    NewEventHub(),
		done:   make(chan struct{}),
	}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", s.handleHealth)
	mux.HandleFunc("GET /api/specs", s.handleSpecs)
	mux.HandleFunc("GET /api/specs/{dir}", s.handleSpecDetail)
	mux.HandleFunc("GET /api/config", s.handleConfigGet)
	mux.HandleFunc("GET /api/config/raw", s.handleConfigRaw)
	mux.HandleFunc("PUT /api/config", s.handleConfigPut)
	mux.HandleFunc("PATCH /api/specs/{dir}/status", s.handleSpecStatus)
	mux.HandleFunc("PATCH /api/specs/{dir}/tasks/{file}/status", s.handleTaskStatus)
	mux.HandleFunc("GET /api/events", s.handleEvents)
	mux.Handle("/", s.handleStatic())
	s.http = &http.Server{Addr: addr, Handler: withCORS(mux)}
	return s, nil
}

// Addr returns the listen address.
func (s *Server) Addr() string {
	return s.addr
}

// Run starts the server and blocks until ctx is canceled.
func (s *Server) Run(ctx context.Context) error {
	if err := s.startWatcher(); err != nil {
		return fmt.Errorf("start filesystem watcher: %w", err)
	}
	ln, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	errCh := make(chan error, 1)
	go func() {
		errCh <- s.http.Serve(ln)
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = s.http.Shutdown(shutdownCtx)
		close(s.done)
		if s.watcher != nil {
			_ = s.watcher.Close()
		}
		return nil
	case err := <-errCh:
		if err == http.ErrServerClosed {
			return nil
		}
		return err
	}
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, PATCH, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s *Server) specDir(dir string) (string, error) {
	if dir == "" || dir == "." || dir == ".." || strings.Contains(dir, "..") {
		return "", errNotFound
	}
	clean := filepath.Clean(dir)
	if clean != dir {
		return "", errNotFound
	}
	path := filepath.Join(s.root, s.cfg.SpecsDir, dir)
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return "", errNotFound
		}
		return "", err
	}
	return path, nil
}

var errNotFound = errors.New("not found")

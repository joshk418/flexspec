package ui

import (
	"io"
	"io/fs"
	"net/http"
	"strings"
	"time"
)

func (s *Server) handleStatic() http.Handler {
	fileServer := http.FileServer(http.FS(s.static))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/") {
			http.NotFound(w, r)
			return
		}
		path := strings.TrimPrefix(r.URL.Path, "/")
		if path == "" {
			s.serveIndex(w, r)
			return
		}
		info, err := fs.Stat(s.static, path)
		if err != nil || info.IsDir() {
			// Unknown path or directory: serve the SPA entry point so client-side
			// routing can take over. Never emit a directory listing.
			s.serveIndex(w, r)
			return
		}
		r.URL.Path = "/" + path
		fileServer.ServeHTTP(w, r)
	})
}

func (s *Server) serveIndex(w http.ResponseWriter, _ *http.Request) {
	data, err := fs.ReadFile(s.static, "index.html")
	if err != nil {
		http.Error(w, "management UI not built; run `make build-ui`", http.StatusServiceUnavailable)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write(data)
}

// StubStaticFS returns a minimal FS when ui/dist is not built.
func StubStaticFS() fs.FS {
	return stubFS{}
}

type stubFS struct{}

func (stubFS) Open(name string) (fs.File, error) {
	if name == "." || name == "index.html" {
		return stubFile{}, nil
	}
	return nil, fs.ErrNotExist
}

type stubFile struct{}

func (stubFile) Stat() (fs.FileInfo, error) { return stubInfo{}, nil }
func (stubFile) Read([]byte) (int, error)   { return 0, io.EOF }
func (stubFile) Close() error               { return nil }

type stubInfo struct{}

func (stubInfo) Name() string       { return "index.html" }
func (stubInfo) Size() int64        { return 0 }
func (stubInfo) Mode() fs.FileMode  { return 0 }
func (stubInfo) ModTime() time.Time { return time.Time{} }
func (stubInfo) IsDir() bool        { return false }
func (stubInfo) Sys() any           { return nil }

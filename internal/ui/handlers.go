package ui

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/joshk418/flexspec/internal/config"
	"github.com/joshk418/flexspec/internal/spec"
)

func (s *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (s *Server) handleSpecs(w http.ResponseWriter, _ *http.Request) {
	entries, err := spec.List(s.root, s.cfg)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, encodeSpecs(entries))
}

func (s *Server) handleSpecDetail(w http.ResponseWriter, r *http.Request) {
	dir := r.PathValue("dir")
	specPath, err := s.specDir(dir)
	if err != nil {
		if errors.Is(err, errNotFound) {
			writeError(w, http.StatusNotFound, "spec not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	readme := filepath.Join(specPath, "README.md")
	parts, err := spec.ReadFileParts(readme)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	meta, err := spec.ParseSpecMeta(readme)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	detail := SpecDetailJSON{
		SpecJSON: SpecJSON{
			ID:          specIDFromDir(dir),
			Dir:         dir,
			Name:        meta.Name,
			Description: meta.Description,
			Status:      meta.Status,
			SpecType:    meta.SpecType,
		},
		Markdown: parts.Body,
	}
	tasksDir := filepath.Join(specPath, "tasks")
	entries, err := os.ReadDir(tasksDir)
	if err != nil && !os.IsNotExist(err) {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	for _, e := range entries {
		if e.IsDir() || !strings.HasPrefix(e.Name(), "T-") || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		taskPath := filepath.Join(tasksDir, e.Name())
		tparts, err := spec.ReadFileParts(taskPath)
		if err != nil {
			continue
		}
		tmeta, err := spec.ParseTaskMeta(taskPath)
		if err != nil {
			continue
		}
		detail.Tasks = append(detail.Tasks, TaskDetailJSON{
			TaskJSON: TaskJSON{
				ID:     tmeta.ID,
				File:   e.Name(),
				Name:   tmeta.Name,
				Status: tmeta.Status,
			},
			Markdown: tparts.Body,
		})
	}
	writeJSON(w, http.StatusOK, detail)
}

func specIDFromDir(dir string) string {
	i := strings.IndexByte(dir, '-')
	if i <= 0 {
		return dir
	}
	return dir[:i]
}

func (s *Server) handleConfigGet(w http.ResponseWriter, _ *http.Request) {
	cfg, err := config.Load(s.root)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, cfg)
}

func (s *Server) handleConfigRaw(w http.ResponseWriter, _ *http.Request) {
	data, err := os.ReadFile(config.ConfigPath(s.root))
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}

func (s *Server) handleConfigPut(w http.ResponseWriter, r *http.Request) {
	var cfg config.Config
	if err := json.NewDecoder(r.Body).Decode(&cfg); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if err := config.Save(s.root, cfg); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	s.cfg = cfg
	s.hub.Broadcast("specs-changed")
	writeJSON(w, http.StatusOK, cfg)
}

func (s *Server) handleSpecStatus(w http.ResponseWriter, r *http.Request) {
	dir := r.PathValue("dir")
	specPath, err := s.specDir(dir)
	if err != nil {
		if errors.Is(err, errNotFound) {
			writeError(w, http.StatusNotFound, "spec not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	var req StatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if strings.TrimSpace(req.Status) == "" {
		writeError(w, http.StatusBadRequest, "status is required")
		return
	}
	readme := filepath.Join(specPath, "README.md")
	if err := spec.SetFileStatus(readme, req.Status); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	s.hub.Broadcast("specs-changed")
	writeJSON(w, http.StatusOK, map[string]string{"status": req.Status})
}

func (s *Server) handleTaskStatus(w http.ResponseWriter, r *http.Request) {
	dir := r.PathValue("dir")
	file := r.PathValue("file")
	specPath, err := s.specDir(dir)
	if err != nil {
		if errors.Is(err, errNotFound) {
			writeError(w, http.StatusNotFound, "spec not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if file == "" || strings.Contains(file, "..") || strings.ContainsRune(file, filepath.Separator) {
		writeError(w, http.StatusBadRequest, "invalid task file")
		return
	}
	var req StatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	taskPath := filepath.Join(specPath, "tasks", file)
	if _, err := os.Stat(taskPath); err != nil {
		writeError(w, http.StatusNotFound, "task not found")
		return
	}
	if err := spec.SetFileStatus(taskPath, req.Status); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	s.hub.Broadcast("specs-changed")
	writeJSON(w, http.StatusOK, map[string]string{"status": req.Status})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

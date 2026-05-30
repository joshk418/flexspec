package spec

import (
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"sync"
	"time"
)

type memFS struct {
	mu    sync.Mutex
	dirs  map[string]struct{}
	files map[string][]byte
}

func newMemFS() *memFS {
	return &memFS{
		dirs:  map[string]struct{}{},
		files: map[string][]byte{},
	}
}

func (m *memFS) clean(path string) string {
	return filepath.Clean(path)
}

func (m *memFS) ensureParent(path string) {
	current := m.clean(path)
	for {
		parent := filepath.Dir(current)
		if parent == current || parent == "." {
			break
		}
		m.dirs[parent] = struct{}{}
		current = parent
	}
}

func (m *memFS) ReadDir(name string) ([]os.DirEntry, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	name = m.clean(name)
	children := m.childNames(name)
	if len(children) == 0 {
		if _, ok := m.dirs[name]; !ok {
			return nil, &os.PathError{Op: "open", Path: name, Err: fs.ErrNotExist}
		}
		return nil, nil
	}

	names := make([]string, 0, len(children))
	for child := range children {
		names = append(names, child)
	}
	slices.Sort(names)

	entries := make([]os.DirEntry, 0, len(names))
	for _, child := range names {
		childPath := filepath.Join(name, child)
		if _, ok := m.dirs[childPath]; ok {
			entries = append(entries, memDirEntry{name: child, isDir: true})
			continue
		}
		entries = append(entries, memDirEntry{name: child, isDir: false})
	}
	return entries, nil
}

func (m *memFS) childNames(dir string) map[string]struct{} {
	dir = m.clean(dir)
	out := map[string]struct{}{}

	for path := range m.dirs {
		if path == dir {
			continue
		}
		if filepath.Dir(path) == dir {
			out[filepath.Base(path)] = struct{}{}
		}
	}
	for path := range m.files {
		if filepath.Dir(path) == dir {
			out[filepath.Base(path)] = struct{}{}
		}
	}
	return out
}

func (m *memFS) MkdirAll(path string, perm os.FileMode) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	path = m.clean(path)
	m.dirs[path] = struct{}{}
	m.ensureParent(path)
	return nil
}

func (m *memFS) Stat(name string) (os.FileInfo, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	name = m.clean(name)
	if data, ok := m.files[name]; ok {
		return memFileInfo{name: filepath.Base(name), size: int64(len(data))}, nil
	}
	if _, ok := m.dirs[name]; ok {
		return memFileInfo{name: filepath.Base(name), dir: true}, nil
	}
	return nil, &os.PathError{Op: "stat", Path: name, Err: fs.ErrNotExist}
}

func (m *memFS) ReadFile(name string) ([]byte, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	name = m.clean(name)
	data, ok := m.files[name]
	if !ok {
		return nil, &os.PathError{Op: "open", Path: name, Err: fs.ErrNotExist}
	}
	out := make([]byte, len(data))
	copy(out, data)
	return out, nil
}

func (m *memFS) WriteFile(name string, data []byte, perm os.FileMode) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	name = m.clean(name)
	m.files[name] = append([]byte(nil), data...)
	m.ensureParent(name)
	return nil
}

type memDirEntry struct {
	name  string
	isDir bool
}

func (e memDirEntry) Name() string               { return e.name }
func (e memDirEntry) IsDir() bool                { return e.isDir }
func (e memDirEntry) Type() fs.FileMode          { return fs.ModePerm }
func (e memDirEntry) Info() (fs.FileInfo, error) { return memFileInfo{name: e.name, dir: e.isDir}, nil }

type memFileInfo struct {
	name string
	size int64
	dir  bool
}

func (i memFileInfo) Name() string       { return i.name }
func (i memFileInfo) Size() int64        { return i.size }
func (i memFileInfo) Mode() fs.FileMode  { return fs.ModePerm }
func (i memFileInfo) ModTime() time.Time { return time.Time{} }
func (i memFileInfo) IsDir() bool        { return i.dir }
func (i memFileInfo) Sys() any           { return nil }

func setupProjectMem(fsys *memFS, root string) {
	templates := filepath.Join(root, ".flexspec", "templates")
	simple := []byte("---\nspec_type: simple\n---\n# simple\n")
	expanded := []byte("---\nspec_type: expanded\n---\n# expanded\n")
	_ = fsys.MkdirAll(filepath.Join(templates, "expanded"), dirPerm)
	_ = fsys.WriteFile(filepath.Join(templates, "flexspec-simple.md"), simple, filePerm)
	_ = fsys.WriteFile(filepath.Join(templates, "expanded", "flexspec-expanded.md"), expanded, filePerm)
}

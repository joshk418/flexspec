package spec

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/joshk418/flexspec/internal/config"
)

// SpecMeta is YAML frontmatter from a spec README.md.
type SpecMeta struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Status      string `yaml:"status"`
	SpecType    string `yaml:"spec_type"`
}

// TaskMeta is YAML frontmatter from a task file under tasks/.
type TaskMeta struct {
	ID     string `yaml:"id"`
	Name   string `yaml:"name"`
	Status string `yaml:"status"`
}

// SpecEntry is one spec directory with optional tasks.
type SpecEntry struct {
	ID    string
	Dir   string
	Meta  SpecMeta
	Tasks []TaskEntry
}

// TaskEntry is one task file under an expanded spec.
type TaskEntry struct {
	File string
	Meta TaskMeta
}

// List discovers specs under cfg.SpecsDir relative to root.
func List(root string, cfg config.Config) ([]SpecEntry, error) {
	specsPath := filepath.Join(root, cfg.SpecsDir)
	entries, err := os.ReadDir(specsPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read specs directory %s: %w", specsPath, err)
	}

	var out []SpecEntry
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		dir := filepath.Join(specsPath, e.Name())
		readme := filepath.Join(dir, "README.md")
		if _, err := os.Stat(readme); err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, fmt.Errorf("stat %s: %w", readme, err)
		}

		meta, err := ParseSpecMeta(readme)
		if err != nil {
			return nil, err
		}

		entry := SpecEntry{
			ID:   specID(e.Name()),
			Dir:  e.Name(),
			Meta: meta,
		}

		if strings.EqualFold(meta.SpecType, "expanded") {
			tasks, err := loadTasks(filepath.Join(dir, "tasks"))
			if err != nil {
				return nil, err
			}
			entry.Tasks = tasks
		}

		out = append(out, entry)
	}

	sort.Slice(out, func(i, j int) bool {
		ni, ei := specSortKey(out[i].Dir)
		nj, ej := specSortKey(out[j].Dir)
		if ni != nj {
			return ni < nj
		}
		return ei < ej
	})

	return out, nil
}

func loadTasks(tasksDir string) ([]TaskEntry, error) {
	entries, err := os.ReadDir(tasksDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read tasks directory %s: %w", tasksDir, err)
	}

	var out []TaskEntry
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if name == "README.md" {
			continue
		}
		if !strings.HasPrefix(name, "T-") || !strings.HasSuffix(name, ".md") {
			continue
		}

		path := filepath.Join(tasksDir, name)
		meta, err := ParseTaskMeta(path)
		if err != nil {
			return nil, err
		}
		out = append(out, TaskEntry{File: name, Meta: meta})
	}

	sort.Slice(out, func(i, j int) bool {
		return taskSortKey(out[i].Meta.ID, out[i].File) < taskSortKey(out[j].Meta.ID, out[j].File)
	})

	return out, nil
}

// ParseSpecMeta reads YAML frontmatter from a spec README.md.
func ParseSpecMeta(path string) (SpecMeta, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return SpecMeta{}, fmt.Errorf("read %s: %w", path, err)
	}
	fm, err := splitFrontmatter(string(data))
	if err != nil {
		return SpecMeta{}, fmt.Errorf("%s: %w", path, err)
	}
	var meta SpecMeta
	if err := yaml.Unmarshal([]byte(fm), &meta); err != nil {
		return SpecMeta{}, fmt.Errorf("parse frontmatter in %s: %w", path, err)
	}
	return meta, nil
}

// ParseTaskMeta reads YAML frontmatter from a task markdown file.
func ParseTaskMeta(path string) (TaskMeta, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return TaskMeta{}, fmt.Errorf("read %s: %w", path, err)
	}
	fm, err := splitFrontmatter(string(data))
	if err != nil {
		return TaskMeta{}, fmt.Errorf("%s: %w", path, err)
	}
	var meta TaskMeta
	if err := yaml.Unmarshal([]byte(fm), &meta); err != nil {
		return TaskMeta{}, fmt.Errorf("parse frontmatter in %s: %w", path, err)
	}
	return meta, nil
}

func splitFrontmatter(content string) (string, error) {
	content = strings.TrimPrefix(content, "\ufeff")
	if !strings.HasPrefix(content, "---") {
		return "", fmt.Errorf("missing opening ---")
	}
	rest := content[3:]
	if strings.HasPrefix(rest, "\n") {
		rest = rest[1:]
	} else if strings.HasPrefix(rest, "\r\n") {
		rest = rest[2:]
	}
	end := strings.Index(rest, "\n---")
	if end < 0 {
		return "", fmt.Errorf("missing closing ---")
	}
	return rest[:end], nil
}

func specID(dirName string) string {
	i := strings.IndexByte(dirName, '-')
	if i <= 0 {
		return dirName
	}
	return dirName[:i]
}

func specSortKey(dirName string) (int, string) {
	id := specID(dirName)
	n, err := strconv.Atoi(id)
	if err != nil {
		return 1 << 30, dirName
	}
	return n, dirName
}

func taskSortKey(id, file string) string {
	if id != "" {
		return id
	}
	return file
}

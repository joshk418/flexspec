package spec

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/joshk418/flexspec/internal/config"
)

func isNotExist(err error) bool {
	return errors.Is(err, fs.ErrNotExist) || os.IsNotExist(err)
}

const (
	flexspecDir      = ".flexspec"
	templatesDir     = "templates"
	simpleTemplate   = "flexspec-simple.md"
	expandedTemplate = "expanded/flexspec-expanded.md"
	dirPerm          = 0o755
	filePerm         = 0o644
)

var slugSanitizer = regexp.MustCompile(`[^a-z0-9]+`)

// CreateResult is the outcome of scaffolding a new spec directory.
type CreateResult struct {
	DirName  string
	SpecPath string
}

// Slugify normalizes a spec name into a lowercase hyphenated slug.
func Slugify(name string) (string, error) {
	s := strings.ToLower(strings.TrimSpace(name))
	s = slugSanitizer.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	if s == "" {
		return "", fmt.Errorf("name %q does not produce a valid slug", name)
	}
	return s, nil
}

// NextSequence returns the next zero-padded spec sequence number for specsDir.
func NextSequence(specsDir string) (int, error) {
	return nextSequenceWithFS(defaultFS, specsDir)
}

func nextSequenceWithFS(fsys fileSystem, specsDir string) (int, error) {
	entries, err := fsys.ReadDir(specsDir)
	if err != nil {
		if isNotExist(err) {
			return 1, nil
		}
		return 0, fmt.Errorf("read specs directory %s: %w", specsDir, err)
	}

	max := 0
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		n, err := strconv.Atoi(specID(e.Name()))
		if err != nil {
			continue
		}
		if n > max {
			max = n
		}
	}
	return max + 1, nil
}

// Create scaffolds a new spec directory under cfg.SpecsDir relative to root.
func Create(root string, cfg config.Config, slug string, template string) (CreateResult, error) {
	return createWithFS(defaultFS, root, cfg, slug, template)
}

func createWithFS(fsys fileSystem, root string, cfg config.Config, slug string, template string) (CreateResult, error) {
	if template != "simple" && template != "expanded" {
		return CreateResult{}, fmt.Errorf("invalid template %q; must be simple or expanded", template)
	}

	specsPath := filepath.Join(root, cfg.SpecsDir)
	if err := fsys.MkdirAll(specsPath, dirPerm); err != nil {
		return CreateResult{}, fmt.Errorf("create specs directory %s: %w", specsPath, err)
	}

	seq, err := nextSequenceWithFS(fsys, specsPath)
	if err != nil {
		return CreateResult{}, err
	}

	dirName := fmt.Sprintf("%03d-%s", seq, slug)
	specDir := filepath.Join(specsPath, dirName)

	if _, err := fsys.Stat(specDir); err == nil {
		return CreateResult{}, fmt.Errorf("spec directory %s already exists", dirName)
	} else if !isNotExist(err) {
		return CreateResult{}, fmt.Errorf("stat %s: %w", specDir, err)
	}

	templatePath, err := templatePathFor(fsys, root, template)
	if err != nil {
		return CreateResult{}, err
	}

	data, err := fsys.ReadFile(templatePath)
	if err != nil {
		return CreateResult{}, fmt.Errorf("read template %s: %w", templatePath, err)
	}

	if err := fsys.MkdirAll(specDir, dirPerm); err != nil {
		return CreateResult{}, fmt.Errorf("create spec directory %s: %w", specDir, err)
	}

	readmePath := filepath.Join(specDir, "README.md")
	if err := fsys.WriteFile(readmePath, data, filePerm); err != nil {
		return CreateResult{}, fmt.Errorf("write %s: %w", readmePath, err)
	}

	if template == "expanded" {
		tasksDir := filepath.Join(specDir, "tasks")
		if err := fsys.MkdirAll(tasksDir, dirPerm); err != nil {
			return CreateResult{}, fmt.Errorf("create tasks directory %s: %w", tasksDir, err)
		}
	}

	return CreateResult{
		DirName:  dirName,
		SpecPath: specDir,
	}, nil
}

func templatePathFor(fsys fileSystem, root, template string) (string, error) {
	var rel string
	switch template {
	case "simple":
		rel = simpleTemplate
	case "expanded":
		rel = expandedTemplate
	default:
		return "", fmt.Errorf("invalid template %q", template)
	}

	path := filepath.Join(root, flexspecDir, templatesDir, rel)
	if _, err := fsys.Stat(path); err != nil {
		if isNotExist(err) {
			return "", fmt.Errorf("template %s not found; run `flexspec init` first", path)
		}
		return "", fmt.Errorf("stat %s: %w", path, err)
	}
	return path, nil
}

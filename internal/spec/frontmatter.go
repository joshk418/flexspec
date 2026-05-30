package spec

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// FileParts is a markdown file split into YAML frontmatter and body.
type FileParts struct {
	Frontmatter string
	Body        string
}

// SplitFrontmatter extracts YAML between opening and closing --- delimiters.
func SplitFrontmatter(content string) (string, error) {
	return splitFrontmatter(content)
}

// ReadFileParts reads a markdown file and splits frontmatter from body.
func ReadFileParts(path string) (FileParts, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return FileParts{}, fmt.Errorf("read %s: %w", path, err)
	}
	content := string(data)
	fm, err := splitFrontmatter(content)
	if err != nil {
		return FileParts{}, fmt.Errorf("%s: %w", path, err)
	}
	body, err := bodyAfterFrontmatter(content)
	if err != nil {
		return FileParts{}, fmt.Errorf("%s: %w", path, err)
	}
	return FileParts{Frontmatter: fm, Body: body}, nil
}

// SetFileStatus updates the status field in frontmatter without changing the body.
func SetFileStatus(path, status string) error {
	parts, err := ReadFileParts(path)
	if err != nil {
		return err
	}
	var fields map[string]any
	if err := yaml.Unmarshal([]byte(parts.Frontmatter), &fields); err != nil {
		return fmt.Errorf("parse frontmatter in %s: %w", path, err)
	}
	if fields == nil {
		fields = map[string]any{}
	}
	fields["status"] = status
	updated, err := yaml.Marshal(fields)
	if err != nil {
		return fmt.Errorf("marshal frontmatter in %s: %w", path, err)
	}
	content := "---\n" + string(updated) + "---\n" + parts.Body
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}

func bodyAfterFrontmatter(content string) (string, error) {
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
	after := rest[end+4:]
	if strings.HasPrefix(after, "\n") {
		after = after[1:]
	} else if strings.HasPrefix(after, "\r\n") {
		after = after[2:]
	}
	return after, nil
}

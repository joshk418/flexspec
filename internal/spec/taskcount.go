package spec

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

func countExpandedTaskFiles(tasksDir string) (int, error) {
	entries, err := os.ReadDir(tasksDir)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, nil
		}
		return 0, fmt.Errorf("read tasks directory %s: %w", tasksDir, err)
	}
	n := 0
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if name == "README.md" {
			continue
		}
		if strings.HasPrefix(name, "T-") && strings.HasSuffix(name, ".md") {
			n++
		}
	}
	return n, nil
}

var simpleTaskBullet = regexp.MustCompile(`(?m)^\s*-\s+\*\*T-\d{3}\*\*`)
var simpleTaskTableRow = regexp.MustCompile(`(?m)^\s*\|\s*(?:\*\*)?T-\d{3}(?:\*\*)?\s*\|`)
var metadataTasksSegment = regexp.MustCompile(` · \*\*Tasks\*\*: \d+`)

// CountTasks returns the number of implementation tasks for a spec README.
func CountTasks(readmePath string, meta SpecMeta) (int, error) {
	if strings.EqualFold(meta.SpecType, "expanded") {
		specDir := filepath.Dir(readmePath)
		return countExpandedTaskFiles(filepath.Join(specDir, "tasks"))
	}

	parts, err := ReadFileParts(readmePath)
	if err != nil {
		return 0, err
	}
	return countSimpleTasks(parts.Body), nil
}

func countSimpleTasks(body string) int {
	return len(simpleTaskBullet.FindAllStringIndex(body, -1)) + len(simpleTaskTableRow.FindAllStringIndex(body, -1))
}

// EffectiveTaskCount returns frontmatter task_count when set, otherwise CountTasks.
func EffectiveTaskCount(readmePath string, meta SpecMeta) (int, error) {
	if meta.TaskCount != nil {
		return *meta.TaskCount, nil
	}
	return CountTasks(readmePath, meta)
}

// SyncTaskCount writes task_count to frontmatter and updates the metadata blockquote line.
func SyncTaskCount(readmePath string) error {
	parts, err := ReadFileParts(readmePath)
	if err != nil {
		return err
	}
	var fields map[string]any
	if err := yaml.Unmarshal([]byte(parts.Frontmatter), &fields); err != nil {
		return fmt.Errorf("parse frontmatter in %s: %w", readmePath, err)
	}
	if fields == nil {
		fields = map[string]any{}
	}

	meta, err := ParseSpecMeta(readmePath)
	if err != nil {
		return err
	}
	count, err := CountTasks(readmePath, meta)
	if err != nil {
		return err
	}

	fields["task_count"] = count
	updated, err := yaml.Marshal(fields)
	if err != nil {
		return fmt.Errorf("marshal frontmatter in %s: %w", readmePath, err)
	}

	body := upsertMetadataTasksLine(parts.Body, count)
	content := "---\n" + string(updated) + "---\n" + body
	if err := os.WriteFile(readmePath, []byte(content), 0o644); err != nil {
		return fmt.Errorf("write %s: %w", readmePath, err)
	}
	return nil
}

func upsertMetadataTasksLine(body string, count int) string {
	segment := fmt.Sprintf(" · **Tasks**: %d", count)
	if metadataTasksSegment.MatchString(body) {
		return metadataTasksSegment.ReplaceAllString(body, segment)
	}
	lines := strings.Split(body, "\n")
	for i, line := range lines {
		if strings.HasPrefix(line, "> **Status**:") {
			lines[i] = strings.TrimRight(line, " \t\r") + segment
			return strings.Join(lines, "\n")
		}
	}
	return body
}

// NeedsTaskCountSync reports whether SyncTaskCount would change the README.
func NeedsTaskCountSync(readmePath string) (bool, error) {
	meta, err := ParseSpecMeta(readmePath)
	if err != nil {
		return false, err
	}
	computed, err := CountTasks(readmePath, meta)
	if err != nil {
		return false, err
	}
	if meta.TaskCount == nil || *meta.TaskCount != computed {
		return true, nil
	}
	parts, err := ReadFileParts(readmePath)
	if err != nil {
		return false, err
	}
	want := fmt.Sprintf("**Tasks**: %d", computed)
	return !strings.Contains(parts.Body, want), nil
}

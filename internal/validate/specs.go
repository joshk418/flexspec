package validate

import (
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/joshk418/flexspec/internal/config"
	"github.com/joshk418/flexspec/internal/spec"
)

var specDirPattern = regexp.MustCompile(`^\d{3}-`)

// CheckSpecs validates directories and frontmatter under specs_dir.
func CheckSpecs(root string, cfg config.Config, _ Options) []Finding {
	specsRel := cfg.SpecsDir
	specsPath := filepath.Join(root, specsRel)

	entries, err := os.ReadDir(specsPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []Finding{{
				Severity: SeverityWarning,
				Path:     specsRel,
				Rule:     "specs.dir_missing",
				Message:  "specs directory does not exist yet",
			}}
		}
		return []Finding{{
			Severity: SeverityError,
			Path:     specsRel,
			Rule:     "specs.read",
			Message:  err.Error(),
		}}
	}

	var findings []Finding
	seenSeq := make(map[int][]string)

	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		name := e.Name()
		dirRel := filepath.Join(specsRel, name)

		if !specDirPattern.MatchString(name) {
			findings = append(findings, Finding{
				Severity: SeverityWarning,
				Path:     dirRel,
				Rule:     "specs.orphan_dir",
				Message:  "directory name does not match NNN-slug pattern",
			})
			continue
		}

		if n, err := strconv.Atoi(name[:3]); err == nil {
			seenSeq[n] = append(seenSeq[n], name)
		}

		readmeRel := filepath.Join(dirRel, "README.md")
		readmePath := filepath.Join(specsPath, name, "README.md")
		if _, err := os.Stat(readmePath); err != nil {
			if os.IsNotExist(err) {
				findings = append(findings, Finding{
					Severity: SeverityWarning,
					Path:     dirRel,
					Rule:     "specs.orphan_dir",
					Message:  "spec directory has no README.md",
				})
				continue
			}
			findings = append(findings, Finding{
				Severity: SeverityError,
				Path:     readmeRel,
				Rule:     "specs.readme_stat",
				Message:  err.Error(),
			})
			continue
		}

		meta, err := spec.ParseSpecMeta(readmePath)
		if err != nil {
			findings = append(findings, Finding{
				Severity: SeverityError,
				Path:     readmeRel,
				Rule:     "specs.frontmatter",
				Message:  err.Error(),
			})
			continue
		}

		computed, err := spec.CountTasks(readmePath, meta)
		if err != nil {
			findings = append(findings, Finding{
				Severity: SeverityError,
				Path:     readmeRel,
				Rule:     "specs.task_count",
				Message:  err.Error(),
			})
			continue
		}
		if meta.TaskCount != nil && *meta.TaskCount != computed {
			findings = append(findings, Finding{
				Severity: SeverityWarning,
				Path:     readmeRel,
				Rule:     "specs.task_count_mismatch",
				Message:  "frontmatter task_count " + strconv.Itoa(*meta.TaskCount) + " does not match computed " + strconv.Itoa(computed),
			})
		}

		if !strings.EqualFold(meta.SpecType, "expanded") {
			continue
		}

		tasksRel := filepath.Join(dirRel, "tasks")
		tasksPath := filepath.Join(specsPath, name, "tasks")
		taskEntries, err := os.ReadDir(tasksPath)
		if err != nil {
			if os.IsNotExist(err) {
				findings = append(findings, Finding{
					Severity: SeverityWarning,
					Path:     tasksRel,
					Rule:     "specs.tasks_missing",
					Message:  "expanded spec has no tasks directory",
				})
				continue
			}
			findings = append(findings, Finding{
				Severity: SeverityError,
				Path:     tasksRel,
				Rule:     "specs.tasks_read",
				Message:  err.Error(),
			})
			continue
		}

		for _, te := range taskEntries {
			if te.IsDir() {
				continue
			}
			fileName := te.Name()
			if !strings.HasPrefix(fileName, "T-") || !strings.HasSuffix(fileName, ".md") {
				continue
			}
			taskRel := filepath.Join(tasksRel, fileName)
			taskPath := filepath.Join(tasksPath, fileName)
			if _, err := spec.ParseTaskMeta(taskPath); err != nil {
				findings = append(findings, Finding{
					Severity: SeverityError,
					Path:     taskRel,
					Rule:     "specs.task_frontmatter",
					Message:  err.Error(),
				})
			}
		}
	}

	for n, dirs := range seenSeq {
		if len(dirs) < 2 {
			continue
		}
		findings = append(findings, Finding{
			Severity: SeverityWarning,
			Path:     specsRel,
			Rule:     "specs.duplicate_sequence",
			Message:  "duplicate spec sequence " + strconv.Itoa(n) + ": " + strings.Join(dirs, ", "),
		})
	}

	return findings
}

package glossary

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

const glossaryFile = "glossary.yaml"

// Document is the top-level glossary file shape.
type Document struct {
	Version string  `yaml:"version" json:"version"`
	Updated string  `yaml:"updated" json:"updated"`
	Terms   []Entry `yaml:"terms" json:"terms"`
}

// Entry is one glossary term definition.
type Entry struct {
	Term       string   `yaml:"term" json:"term"`
	Definition string   `yaml:"definition" json:"definition"`
	Category   string   `yaml:"category,omitempty" json:"category,omitempty"`
	Aliases    []string `yaml:"aliases,omitempty" json:"aliases,omitempty"`
	Sources    []string `yaml:"sources,omitempty" json:"sources,omitempty"`
	Created    string   `yaml:"created" json:"created"`
	Updated    string   `yaml:"updated" json:"updated"`
}

// DefaultDocument returns an empty glossary document.
func DefaultDocument() Document {
	return Document{
		Version: "1.0",
		Updated: time.Now().Format(time.RFC3339),
		Terms:   []Entry{},
	}
}

// Load reads .flexspec/glossary.yaml under root. A missing file returns an empty document.
func Load(root string) (Document, error) {
	path := filepath.Join(root, ".flexspec", glossaryFile)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return DefaultDocument(), nil
		}
		return Document{}, fmt.Errorf("read %s: %w", path, err)
	}

	var doc Document
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return Document{}, fmt.Errorf("parse %s: %w", path, err)
	}
	if err := validate(doc, path); err != nil {
		return Document{}, err
	}
	return doc, nil
}

// Save writes doc to .flexspec/glossary.yaml under root, sorting terms by name.
func Save(root string, doc Document) error {
	sortEntries(doc.Terms)
	path := filepath.Join(root, ".flexspec", glossaryFile)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create %s: %w", filepath.Dir(path), err)
	}
	data, err := yaml.Marshal(&doc)
	if err != nil {
		return fmt.Errorf("marshal glossary: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}

// Upsert inserts or replaces an entry by term, preserving unprovided fields unless redefined.
func Upsert(doc Document, entry Entry) (Document, error) {
	if strings.TrimSpace(entry.Term) == "" {
		return doc, fmt.Errorf("term is required")
	}
	if strings.TrimSpace(entry.Definition) == "" {
		return doc, fmt.Errorf("definition is required")
	}

	now := time.Now().Format(time.RFC3339)
	entry.Updated = now

	found := false
	for i, e := range doc.Terms {
		if strings.EqualFold(e.Term, entry.Term) {
			entry = mergeEntry(e, entry)
			if entry.Created == "" {
				entry.Created = e.Created
			}
			if entry.Created == "" {
				entry.Created = now
			}
			doc.Terms[i] = entry
			found = true
			break
		}
	}
	if !found {
		entry.Created = now
		doc.Terms = append(doc.Terms, entry)
	}

	doc.Updated = now
	sortEntries(doc.Terms)
	return doc, nil
}

// Query returns matching entries ranked by relevance: exact term/alias first, then substring.
func Query(doc Document, text string) []Entry {
	text = strings.TrimSpace(text)
	if text == "" {
		return nil
	}
	lower := strings.ToLower(text)

	var exact []Entry
	var alias []Entry
	var substring []Entry

	for _, e := range doc.Terms {
		termLower := strings.ToLower(e.Term)
		if termLower == lower {
			exact = append(exact, e)
			continue
		}
		matchedAlias := false
		for _, a := range e.Aliases {
			if strings.ToLower(a) == lower {
				alias = append(alias, e)
				matchedAlias = true
				break
			}
		}
		if matchedAlias {
			continue
		}
		if strings.Contains(termLower, lower) ||
			containsFold(e.Aliases, lower) ||
			strings.Contains(strings.ToLower(e.Definition), lower) ||
			strings.Contains(strings.ToLower(e.Category), lower) {
			substring = append(substring, e)
		}
	}

	out := make([]Entry, 0, len(exact)+len(alias)+len(substring))
	out = append(out, exact...)
	out = append(out, alias...)
	out = append(out, substring...)
	return out
}

func mergeEntry(existing, next Entry) Entry {
	if next.Category == "" {
		next.Category = existing.Category
	}
	next.Aliases = appendUnique(existing.Aliases, next.Aliases)
	next.Sources = appendUnique(existing.Sources, next.Sources)
	return next
}

func appendUnique(existing, added []string) []string {
	out := make([]string, 0, len(existing)+len(added))
	seen := make(map[string]struct{}, len(existing)+len(added))
	for _, value := range append(existing, added...) {
		trimmed := strings.TrimSpace(value)
		if trimmed == "" {
			continue
		}
		key := strings.ToLower(trimmed)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		out = append(out, trimmed)
	}
	return out
}

func containsFold(values []string, lower string) bool {
	for _, value := range values {
		if strings.Contains(strings.ToLower(value), lower) {
			return true
		}
	}
	return false
}

func sortEntries(entries []Entry) {
	sort.Slice(entries, func(i, j int) bool {
		return strings.ToLower(entries[i].Term) < strings.ToLower(entries[j].Term)
	})
}

func validate(doc Document, path string) error {
	for i, e := range doc.Terms {
		if strings.TrimSpace(e.Term) == "" {
			return fmt.Errorf("%s: terms[%d].term is required", path, i)
		}
		if strings.TrimSpace(e.Definition) == "" {
			return fmt.Errorf("%s: terms[%d].definition is required", path, i)
		}
	}
	return nil
}

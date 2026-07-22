// Package corpus provides read-only access to the frozen prompt corpus.
// It parses the flat YAML frontmatter and returns verbatim bodies; nothing in
// this package ever writes under corpus/.
package corpus

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Artifact is one corpus file: frontmatter metadata plus the verbatim prompt
// body below it. BodyStart is the byte offset of Body within the raw file, so
// segment offsets (relative to Body) can always be mapped back to the file.
type Artifact struct {
	Path      string
	Rel       string // e.g. "corpus/tier-a-official/anthropic--claude-ai--fable-5--2026-06-09.md"
	Meta      map[string]string
	Body      string
	BodyStart int
}

// Load reads every artifact under dir recursively, skipping the corpus
// documentation files. Results are sorted by Rel for determinism.
func Load(dir string) ([]Artifact, error) {
	skip := map[string]bool{"README.md": true, "SOURCES.md": true, "ETHICS-AND-PUBLISHABILITY.md": true}
	root := filepath.Clean(dir)
	var arts []Artifact
	err := filepath.WalkDir(root, func(p string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() || !strings.HasSuffix(p, ".md") || skip[filepath.Base(p)] {
			return nil
		}
		raw, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		meta, body, off, err := ParseFrontmatter(string(raw))
		if err != nil {
			return fmt.Errorf("%s: %w", p, err)
		}
		rel, err := filepath.Rel(filepath.Dir(root), p)
		if err != nil {
			rel = p
		}
		arts = append(arts, Artifact{Path: p, Rel: filepath.ToSlash(rel), Meta: meta, Body: body, BodyStart: off})
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Slice(arts, func(i, j int) bool { return arts[i].Rel < arts[j].Rel })
	return arts, nil
}

// ParseFrontmatter splits a corpus file into flat key→value frontmatter and
// the verbatim body. The frontmatter must start at byte 0 with a "---" line
// and end at the next "---" line. A value is the remainder of its line after
// the first ":"; indented continuation lines are appended to the previous
// key. Returns the metadata, the body, and the byte offset of the body.
func ParseFrontmatter(s string) (map[string]string, string, int, error) {
	const fence = "---"
	if !strings.HasPrefix(s, fence+"\n") {
		return nil, "", 0, fmt.Errorf("missing frontmatter opening fence")
	}
	meta := map[string]string{}
	lastKey := ""
	off := len(fence) + 1
	rest := s[off:]
	for {
		nl := strings.IndexByte(rest, '\n')
		line := rest
		if nl >= 0 {
			line = rest[:nl]
		}
		consumed := len(line)
		if nl >= 0 {
			consumed++
		}
		off += consumed
		rest = rest[consumed:]
		if strings.TrimRight(line, " \t") == fence {
			break
		}
		if nl < 0 {
			return nil, "", 0, fmt.Errorf("missing frontmatter closing fence")
		}
		if strings.HasPrefix(line, " ") || strings.HasPrefix(line, "\t") {
			if lastKey != "" {
				meta[lastKey] += " " + strings.TrimSpace(line)
			}
			continue
		}
		if i := strings.IndexByte(line, ':'); i > 0 {
			k := strings.TrimSpace(line[:i])
			v := strings.Trim(strings.TrimSpace(line[i+1:]), `"`)
			meta[k] = v
			lastKey = k
		}
	}
	body := rest
	if strings.HasPrefix(body, "\n") {
		body = body[1:]
		off++
	}
	return meta, body, off, nil
}

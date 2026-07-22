// Command segment is step 1 of the decomposition pipeline (see CLAUDE.md
// "Decomposition methodology").
//
// It reads the frozen corpus (read-only), strips each file's YAML
// frontmatter, and splits the verbatim prompt body into analysis segments
// along its own structure: markdown headers, whole-line XML tags, `====`
// separators, ChatML markers, fenced code blocks, and blank lines. Each
// segment is tagged with a heuristic layer (behavioral | tool-definition-
// schema | few-shot | template-var) and carries its source path plus byte
// offsets into the body, so every downstream claim can be traced back to the
// exact corpus bytes. Segments shorter than 40 characters are dropped.
//
// Input:  corpus/  (never modified)
// Output: artefacts/segments.jsonl — one segment per line:
//
//	{"id","source","product","layer","start","end","text"}
//
// No network access; fully deterministic.
package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

const minChars = 40

type seg struct {
	ID      string `json:"id"`
	Source  string `json:"source"`
	Product string `json:"product"`
	Layer   string `json:"layer"`
	Start   int    `json:"start"`
	End     int    `json:"end"`
	Text    string `json:"text"`
}

var (
	reHeader      = regexp.MustCompile(`^(#{1,6})\s+.+`)
	reXMLLine     = regexp.MustCompile(`^</?[A-Za-z][A-Za-z0-9_.-]*(\s[^<>]*)?>\s*$`)
	reSeparator   = regexp.MustCompile(`^={4,}\s*$`)
	reChatML      = regexp.MustCompile(`^<\|[a-z_]+\|>`)
	reFence       = regexp.MustCompile("^(```|~~~)")
	rePlaceholder = regexp.MustCompile(`^\s*(\{\{[^}]+\}\}|\[[A-Z0-9_ -]+\])\s*$`)
	reJSONHint    = regexp.MustCompile(`"(type|properties|parameters|description|required|name)"\s*:`)
	reExampleSec  = regexp.MustCompile(`(?i)^#{1,6}\s+.*(example|alignment|demonstration)`)
)

func main() {
	log.SetFlags(0)
	corpusDir := flag.String("corpus", "corpus", "corpus root (read-only)")
	out := flag.String("out", "artefacts/segments.jsonl", "output JSONL path")
	flag.Parse()

	skip := map[string]bool{"README.md": true, "SOURCES.md": true, "ETHICS-AND-PUBLISHABILITY.md": true}
	var files []string
	err := filepath.WalkDir(*corpusDir, func(p string, d os.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if !d.IsDir() && strings.HasSuffix(p, ".md") && !skip[filepath.Base(p)] {
			files = append(files, p)
		}
		return nil
	})
	if err != nil {
		log.Fatalf("corpus: %v", err)
	}
	sort.Strings(files)
	if len(files) == 0 {
		log.Fatalf("no artifacts under %s", *corpusDir)
	}

	var all []seg
	layerCount := map[string]int{}
	for _, p := range files {
		raw, err := os.ReadFile(p)
		if err != nil {
			log.Fatal(err)
		}
		meta, body, err := parseFrontmatter(string(raw))
		if err != nil {
			log.Fatalf("%s: %v", p, err)
		}
		rel, _ := filepath.Rel(filepath.Dir(filepath.Clean(*corpusDir)), p)
		segs := split(filepath.ToSlash(rel), meta["product"], body)
		for _, s := range segs {
			layerCount[s.Layer]++
		}
		all = append(all, segs...)
	}

	if err := os.MkdirAll(filepath.Dir(*out), 0o755); err != nil {
		log.Fatal(err)
	}
	f, err := os.Create(*out)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	enc := json.NewEncoder(w)
	for _, s := range all {
		if err := enc.Encode(s); err != nil {
			log.Fatal(err)
		}
	}
	if err := w.Flush(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%d artifacts → %d segments (behavioral=%d, tool-schema=%d, few-shot=%d, template-var=%d)\n",
		len(files), len(all), layerCount["behavioral"], layerCount["tool-definition-schema"],
		layerCount["few-shot"], layerCount["template-var"])
	fmt.Printf("wrote %s\n", *out)
}

// parseFrontmatter splits "---\nkey: value...\n---\n" from the verbatim body.
// Values take the remainder of the line after the first ":"; indented lines
// continue the previous key.
func parseFrontmatter(s string) (map[string]string, string, error) {
	if !strings.HasPrefix(s, "---\n") {
		return nil, "", fmt.Errorf("missing frontmatter fence")
	}
	meta := map[string]string{}
	lastKey := ""
	rest := s[4:]
	for {
		nl := strings.IndexByte(rest, '\n')
		if nl < 0 {
			return nil, "", fmt.Errorf("missing closing fence")
		}
		line := rest[:nl]
		rest = rest[nl+1:]
		if strings.TrimRight(line, " \t") == "---" {
			break
		}
		if strings.HasPrefix(line, " ") || strings.HasPrefix(line, "\t") {
			if lastKey != "" {
				meta[lastKey] += " " + strings.TrimSpace(line)
			}
			continue
		}
		if i := strings.IndexByte(line, ':'); i > 0 {
			k := strings.TrimSpace(line[:i])
			meta[k] = strings.Trim(strings.TrimSpace(line[i+1:]), `"`)
			lastKey = k
		}
	}
	return meta, strings.TrimPrefix(rest, "\n"), nil
}

// split walks the body line by line, flushing a segment at every structural
// boundary. Offsets are byte positions within the body.
func split(source, product, body string) []seg {
	var out []seg
	var lines []string
	start, off, n := -1, 0, 0
	inFence, inExample := false, false

	flush := func(end int, fenced bool) {
		if start < 0 {
			return
		}
		text := strings.Join(lines, "\n")
		s := start
		lines, start = nil, -1
		trimmed := strings.TrimSpace(text)
		if trimmed == "" {
			return
		}
		layer := classify(trimmed, fenced, inExample)
		if len(trimmed) < minChars && layer != "tool-definition-schema" {
			return
		}
		out = append(out, seg{
			ID:     fmt.Sprintf("%s#%04d", strings.TrimSuffix(filepath.Base(source), ".md"), n),
			Source: source, Product: product, Layer: layer, Start: s, End: end, Text: text,
		})
		n++
	}

	for off <= len(body) {
		nl := strings.IndexByte(body[off:], '\n')
		line := body[off:]
		lineEnd := len(body)
		if nl >= 0 {
			line = body[off : off+nl]
			lineEnd = off + nl
		}
		switch {
		case inFence:
			if start < 0 {
				start = off
			}
			lines = append(lines, line)
			if reFence.MatchString(line) {
				inFence = false
				flush(lineEnd, true)
			}
		case reFence.MatchString(line):
			flush(off, false)
			inFence = true
			start = off
			lines = append(lines, line)
		case strings.TrimSpace(line) == "":
			flush(off, false)
		case reHeader.MatchString(line):
			flush(off, false)
			inExample = reExampleSec.MatchString(line)
		case reXMLLine.MatchString(line), reSeparator.MatchString(line), reChatML.MatchString(line):
			flush(off, false)
		default:
			if start < 0 {
				start = off
			}
			lines = append(lines, line)
		}
		if nl < 0 {
			break
		}
		off = lineEnd + 1
	}
	flush(len(body), false)
	return out
}

func classify(trimmed string, fenced, inExample bool) string {
	hints := len(reJSONHint.FindAllString(trimmed, 4))
	switch {
	case rePlaceholder.MatchString(trimmed):
		return "template-var"
	case (strings.HasPrefix(trimmed, "{") || strings.HasPrefix(trimmed, "[") || (fenced && hints >= 1)) && hints >= 1 && len(trimmed) >= 120:
		return "tool-definition-schema"
	case hints >= 3 && len(trimmed) >= 400:
		return "tool-definition-schema"
	case inExample || fenced:
		return "few-shot"
	default:
		return "behavioral"
	}
}

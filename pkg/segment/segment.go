// Package segment turns corpus artifacts into the analysis units consumed by
// the embedding pipeline (roadmap step 2). Corpus bodies are read-only input;
// every segment carries byte offsets back into the body so each boundary is
// auditable. Layer tags are v1 heuristics — the paper reports them as such.
package segment

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/mikolajsemeniuk/prism/pkg/corpus"
)

// Layer tags which prompt layer a segment belongs to (findings F1/F11).
type Layer string

const (
	Behavioral  Layer = "behavioral"
	ToolSchema  Layer = "tool-definition-schema"
	FewShot     Layer = "few-shot"
	TemplateVar Layer = "template-var"
	Wrapper     Layer = "wrapper"
)

// Segment is one analysis unit. Start/End are byte offsets into the artifact
// BODY (add Artifact.BodyStart for file offsets).
type Segment struct {
	ID           string   `json:"id"`
	Source       string   `json:"source"`
	Product      string   `json:"product"`
	Vendor       string   `json:"vendor"`
	Tier         string   `json:"tier"`
	ProductClass string   `json:"product_class"`
	Layer        Layer    `json:"layer"`
	SectionPath  []string `json:"section_path"`
	Start        int      `json:"start"`
	End          int      `json:"end"`
	Flags        []string `json:"flags,omitempty"`
	Text         string   `json:"text"`
}

type Options struct {
	MinChars int // segments shorter than this (after trimming) are dropped
}

var (
	reHeader      = regexp.MustCompile(`^(#{1,6})\s+(.+?)\s*$`)
	reXMLOpen     = regexp.MustCompile(`^<([A-Za-z][A-Za-z0-9_.-]*)(\s[^<>]*)?>\s*$`)
	reXMLClose    = regexp.MustCompile(`^</([A-Za-z][A-Za-z0-9_.-]*)>\s*$`)
	reSeparator   = regexp.MustCompile(`^={4,}\s*$`)
	reChatML      = regexp.MustCompile(`^<\|[a-z_]+\|>`)
	reFence       = regexp.MustCompile("^(```|~~~)")
	rePlaceholder = regexp.MustCompile(`^\s*(\{\{[^}]+\}\}|\[[A-Z0-9_ -]+\])\s*$`)
	reFewShotPath = regexp.MustCompile(`(?i)example|alignment|demonstration`)
	reJSONHint    = regexp.MustCompile(`"(type|properties|parameters|description|required|name)"\s*:`)
)

type secEntry struct {
	name    string
	mdLevel int // 0 for XML-tag sections
}

// Split segments one artifact. Deterministic: same input bytes and options →
// same output.
func Split(a corpus.Artifact, opt Options) []Segment {
	base := baseName(a.Rel)
	class := Class(base)
	flags := FileFlags(base)
	if strings.EqualFold(a.Meta["provenance_tier"], "B") {
		flags = append(flags, "source-code-carrier")
	}
	hasWrapper := contains(flags, "wrapper")

	var (
		out           []Segment
		stack         []secEntry
		blockLines    []string
		blockStart    = -1
		inFence       bool
		fenceMark     string
		blockIsFence  bool
		seenStructure bool
		blockIndex    int
		partCounter   int
	)

	flush := func(end int) {
		if blockStart < 0 {
			return
		}
		text := strings.Join(blockLines, "\n")
		start := blockStart
		blockLines, blockStart = nil, -1
		fenced := blockIsFence
		blockIsFence = false
		trimmed := strings.TrimSpace(text)
		idx := blockIndex
		blockIndex++
		if trimmed == "" {
			return
		}
		layer := classify(trimmed, fenced, stack, hasWrapper && !seenStructure && idx < 3)
		if len(trimmed) < opt.MinChars && layer != ToolSchema {
			return
		}
		out = append(out, Segment{
			Source:       a.Rel,
			Product:      a.Meta["product"],
			Vendor:       a.Meta["vendor"],
			Tier:         a.Meta["provenance_tier"],
			ProductClass: class,
			Layer:        layer,
			SectionPath:  pathNames(stack),
			Start:        start,
			End:          end,
			Flags:        flags,
			Text:         text,
		})
	}

	off := 0
	body := a.Body
	for off <= len(body) {
		nl := strings.IndexByte(body[off:], '\n')
		var line string
		lineEnd := len(body)
		if nl >= 0 {
			line = body[off : off+nl]
			lineEnd = off + nl
		} else {
			line = body[off:]
		}

		switch {
		case inFence:
			appendLine(&blockLines, &blockStart, line, off)
			if reFence.MatchString(line) && strings.HasPrefix(line, fenceMark) {
				inFence = false
				flush(lineEnd)
			}
		case reFence.MatchString(line):
			flush(off)
			inFence = true
			blockIsFence = true
			fenceMark = line[:3]
			appendLine(&blockLines, &blockStart, line, off)
		case strings.TrimSpace(line) == "":
			flush(off)
		case reSeparator.MatchString(line) || reChatML.MatchString(line):
			flush(off)
			seenStructure = true
			partCounter++
			stack = []secEntry{{name: fmt.Sprintf("part-%02d", partCounter)}}
		case reHeader.MatchString(line):
			flush(off)
			seenStructure = true
			m := reHeader.FindStringSubmatch(line)
			level := len(m[1])
			for len(stack) > 0 && (stack[len(stack)-1].mdLevel == 0 || stack[len(stack)-1].mdLevel >= level) {
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, secEntry{name: m[2], mdLevel: level})
		case reXMLClose.MatchString(line):
			flush(off)
			name := reXMLClose.FindStringSubmatch(line)[1]
			for i := len(stack) - 1; i >= 0; i-- {
				if stack[i].mdLevel == 0 && stack[i].name == name {
					stack = stack[:i]
					break
				}
			}
		case reXMLOpen.MatchString(line):
			flush(off)
			seenStructure = true
			stack = append(stack, secEntry{name: reXMLOpen.FindStringSubmatch(line)[1]})
		default:
			appendLine(&blockLines, &blockStart, line, off)
		}

		if nl < 0 {
			break
		}
		off = lineEnd + 1
	}
	flush(len(body))

	for i := range out {
		out[i].ID = fmt.Sprintf("%s#%04d", strings.TrimSuffix(base, ".md"), i)
	}
	return out
}

func appendLine(lines *[]string, start *int, line string, off int) {
	if *start < 0 {
		*start = off
	}
	*lines = append(*lines, line)
}

func classify(trimmed string, fenced bool, stack []secEntry, wrapperZone bool) Layer {
	inner := trimmed
	if fenced {
		inner = stripFence(trimmed)
	}
	switch {
	case wrapperZone && strings.Contains(trimmed, "http"):
		return Wrapper
	case rePlaceholder.MatchString(trimmed):
		return TemplateVar
	case looksJSON(inner):
		return ToolSchema
	case pathMatches(stack):
		return FewShot
	case fenced:
		return FewShot
	default:
		return Behavioral
	}
}

func looksJSON(t string) bool {
	t = strings.TrimSpace(t)
	hints := len(reJSONHint.FindAllString(t, 4))
	if (strings.HasPrefix(t, "{") || strings.HasPrefix(t, "[")) && hints >= 1 && len(t) >= 120 {
		return true
	}
	return hints >= 3 && len(t) >= 400
}

func stripFence(t string) string {
	lines := strings.Split(t, "\n")
	if len(lines) >= 2 && reFence.MatchString(lines[0]) {
		lines = lines[1:]
		if len(lines) > 0 && reFence.MatchString(lines[len(lines)-1]) {
			lines = lines[:len(lines)-1]
		}
	}
	return strings.Join(lines, "\n")
}

func pathMatches(stack []secEntry) bool {
	for _, e := range stack {
		if reFewShotPath.MatchString(e.name) {
			return true
		}
	}
	return false
}

func pathNames(stack []secEntry) []string {
	names := make([]string, len(stack))
	for i, e := range stack {
		names[i] = e.name
	}
	return names
}

func baseName(rel string) string {
	if i := strings.LastIndexByte(rel, '/'); i >= 0 {
		return rel[i+1:]
	}
	return rel
}

func contains(xs []string, s string) bool {
	for _, x := range xs {
		if x == s {
			return true
		}
	}
	return false
}

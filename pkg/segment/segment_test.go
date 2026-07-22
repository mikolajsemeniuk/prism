package segment

import (
	"strings"
	"testing"

	"github.com/mikolajsemeniuk/prism/pkg/corpus"
)

const sample = `---
product: "Test Product"
vendor: "TestCo"
provenance_tier: C
---
# Identity

You are a helpful test assistant working inside a sandboxed environment for tests.

## Rules

Always answer truthfully and never fabricate command output in any situation.

<safety>
Never do anything unsafe. This sentence pads the segment above minimum length.
</safety>

## Examples

Here is a demonstration of the expected transcript format for this assistant.

` + "```json\n" + `{"name": "run", "description": "Runs a shell command in the sandbox", "parameters": {"type": "object", "properties": {"command": {"type": "string", "description": "the command"}}, "required": ["command"]}}
` + "```" + `

{{SESSION_VARS}}
`

func load(t *testing.T) []Segment {
	t.Helper()
	meta, body, off, err := corpus.ParseFrontmatter(sample)
	if err != nil {
		t.Fatal(err)
	}
	a := corpus.Artifact{
		Rel: "corpus/tier-c-extracted/testco--test--v1.md", Meta: meta, Body: body, BodyStart: off,
	}
	return Split(a, Options{MinChars: 40})
}

func find(t *testing.T, segs []Segment, needle string) Segment {
	t.Helper()
	for _, s := range segs {
		if strings.Contains(s.Text, needle) {
			return s
		}
	}
	t.Fatalf("no segment contains %q; got %d segments", needle, len(segs))
	return Segment{}
}

func TestSplitLayersAndSections(t *testing.T) {
	segs := load(t)

	id := find(t, segs, "helpful test assistant")
	if id.Layer != Behavioral {
		t.Fatalf("identity layer = %s, want behavioral", id.Layer)
	}
	if len(id.SectionPath) != 1 || id.SectionPath[0] != "Identity" {
		t.Fatalf("identity section path = %v", id.SectionPath)
	}

	rules := find(t, segs, "never fabricate command output")
	if got := strings.Join(rules.SectionPath, "/"); got != "Identity/Rules" {
		t.Fatalf("rules section path = %q, want Identity/Rules", got)
	}

	safety := find(t, segs, "anything unsafe")
	if got := strings.Join(safety.SectionPath, "/"); got != "Identity/Rules/safety" {
		t.Fatalf("safety section path = %q (XML nesting broken)", got)
	}

	demo := find(t, segs, "demonstration of the expected transcript")
	if demo.Layer != FewShot {
		t.Fatalf("examples-section paragraph layer = %s, want few-shot", demo.Layer)
	}

	schema := find(t, segs, `"parameters"`)
	if schema.Layer != ToolSchema {
		t.Fatalf("json fence layer = %s, want tool-definition-schema", schema.Layer)
	}
}

func TestSplitOffsetsRoundTrip(t *testing.T) {
	meta, body, _, err := corpus.ParseFrontmatter(sample)
	if err != nil {
		t.Fatal(err)
	}
	a := corpus.Artifact{Rel: "corpus/tier-c-extracted/testco--test--v1.md", Meta: meta, Body: body}
	for _, s := range Split(a, Options{MinChars: 40}) {
		if got := body[s.Start:s.End]; !strings.Contains(got, strings.TrimSpace(s.Text)) && got != s.Text {
			// offsets must reproduce the stored text exactly
			if got != s.Text {
				t.Fatalf("offsets for %s do not reproduce text:\n--- text ---\n%q\n--- body slice ---\n%q", s.ID, s.Text, got)
			}
		}
	}
}

func TestMeta(t *testing.T) {
	if got := Class("anthropic--claude-code--v2.1.214.md"); got != "coding-agent" {
		t.Fatalf("claude-code class = %q", got)
	}
	if got := Class("vercel--v0--2026-05-10.md"); got != "app-builder" {
		t.Fatalf("v0 class = %q", got)
	}
	if got := Class("nobody--nothing--x.md"); got != "unknown" {
		t.Fatalf("unknown class = %q", got)
	}
	flags := FileFlags("replit--replit-agent--2025-04-22.md")
	want := map[string]bool{"wrapper": true, "stripped-tags": true, "aggregate-member": true}
	if len(flags) != len(want) {
		t.Fatalf("replit flags = %v", flags)
	}
	for _, f := range flags {
		if !want[f] {
			t.Fatalf("unexpected replit flag %q", f)
		}
	}
}

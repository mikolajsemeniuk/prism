// Command fragments is step 7 of the decomposition pipeline (see CLAUDE.md
// "Decomposition methodology") — it exports the VERBATIM experiment stimuli
// for the ablation study.
//
// For every component in the taxonomy it emits the medoid (the cluster's most
// typical real production sentence) and the remaining exemplars as alternate
// stimuli, each with full provenance (corpus source file + byte offsets into
// the prompt body) and size measures (words/chars) for building the
// length-matched placebo condition. Before emitting, every fragment is
// re-verified BYTE-FOR-BYTE against the frozen corpus: the recorded offsets
// must reproduce the recorded text, so a stimulus can never silently drift
// from the artifact it claims to quote.
//
// Deliberately mechanical: stimuli are selected by the medoid rule, never by
// hand, and component names (the interpretive layer) are absent — the
// ablation runner joins them only for reporting. Components the clustering
// did not consolidate (e.g. failure-escalation) are NOT covered here; their
// corpus-attested quotes will be a separate hand-curated input, clearly
// marked as such.
//
// Inputs: artefacts/taxonomy.json, artefacts/segments.jsonl, corpus/
// Output: artefacts/fragments.jsonl — one stimulus per line:
//
//	{"component","role":"medoid"|"exemplar","segment_id","source",
//	 "start","end","words","chars","text"}
package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type seg struct {
	ID     string `json:"id"`
	Source string `json:"source"`
	Start  int    `json:"start"`
	End    int    `json:"end"`
	Text   string `json:"text"`
}

type fragment struct {
	Component string `json:"component"`
	Role      string `json:"role"`
	SegmentID string `json:"segment_id"`
	Source    string `json:"source"`
	Start     int    `json:"start"`
	End       int    `json:"end"`
	Words     int    `json:"words"`
	Chars     int    `json:"chars"`
	Text      string `json:"text"`
}

func main() {
	log.SetFlags(0)
	dir := flag.String("artefacts", "artefacts", "artefacts directory")
	corpusDir := flag.String("corpus", "corpus", "corpus root (read-only, for byte verification)")
	flag.Parse()

	var tx struct {
		Entries []struct {
			ID          int      `json:"id"`
			Category    string   `json:"category"`
			MedoidID    string   `json:"medoid_id"`
			ExemplarIDs []string `json:"exemplar_ids"`
		} `json:"entries"`
	}
	raw, err := os.ReadFile(filepath.Join(*dir, "taxonomy.json"))
	if err != nil {
		log.Fatalf("%v (run cmd/taxonomy first)", err)
	}
	if err := json.Unmarshal(raw, &tx); err != nil {
		log.Fatal(err)
	}

	byID := map[string]seg{}
	f, err := os.Open(filepath.Join(*dir, "segments.jsonl"))
	if err != nil {
		log.Fatalf("%v (run cmd/segment first)", err)
	}
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 0, 1<<20), 64<<20)
	for sc.Scan() {
		var s seg
		if err := json.Unmarshal(sc.Bytes(), &s); err != nil {
			log.Fatal(err)
		}
		byID[s.ID] = s
	}
	f.Close()

	bodies := map[string]string{} // corpus source → prompt body (frontmatter stripped)

	var frags []fragment
	for _, e := range tx.Entries {
		if e.Category != "component" || e.MedoidID == "" {
			continue
		}
		emit := func(segID, role string) {
			s, ok := byID[segID]
			if !ok {
				log.Fatalf("component C%d references unknown segment %s", e.ID, segID)
			}
			verify(*corpusDir, bodies, s)
			frags = append(frags, fragment{
				Component: fmt.Sprintf("C%d", e.ID), Role: role, SegmentID: s.ID,
				Source: s.Source, Start: s.Start, End: s.End,
				Words: len(strings.Fields(s.Text)), Chars: len(s.Text), Text: s.Text,
			})
		}
		emit(e.MedoidID, "medoid")
		for _, ex := range e.ExemplarIDs {
			if ex != e.MedoidID {
				emit(ex, "exemplar")
			}
		}
	}
	if len(frags) == 0 {
		log.Fatal("no components in taxonomy.json — nothing to export")
	}

	outPath := filepath.Join(*dir, "fragments.jsonl")
	out, err := os.Create(outPath)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	w := bufio.NewWriter(out)
	enc := json.NewEncoder(w)
	for _, fr := range frags {
		if err := enc.Encode(fr); err != nil {
			log.Fatal(err)
		}
	}
	if err := w.Flush(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%-6s %-9s %-46s %6s %6s\n", "comp", "role", "source", "words", "chars")
	for _, fr := range frags {
		src := fr.Source
		if len(src) > 46 {
			src = "…" + src[len(src)-45:]
		}
		fmt.Printf("%-6s %-9s %-46s %6d %6d\n", fr.Component, fr.Role, src, fr.Words, fr.Chars)
	}
	fmt.Printf("\nwrote %s: %d stimuli (all byte-verified against corpus)\n", outPath, len(frags))
}

// verify re-derives the fragment from the frozen corpus bytes: the segment's
// [start,end) slice of the prompt body must reproduce its recorded text
// (modulo the trailing newline consumed by the segment boundary).
func verify(corpusDir string, bodies map[string]string, s seg) {
	body, ok := bodies[s.Source]
	if !ok {
		rel := strings.TrimPrefix(s.Source, "corpus/")
		raw, err := os.ReadFile(filepath.Join(corpusDir, rel))
		if err != nil {
			log.Fatalf("verify %s: %v", s.ID, err)
		}
		body = stripFrontmatter(string(raw))
		bodies[s.Source] = body
	}
	if s.Start < 0 || s.End > len(body) || s.Start >= s.End {
		log.Fatalf("verify %s: offsets [%d,%d) out of range (body %d bytes)", s.ID, s.Start, s.End, len(body))
	}
	got := strings.TrimRight(body[s.Start:s.End], "\n")
	want := strings.TrimRight(s.Text, "\n")
	if got != want {
		log.Fatalf("verify %s: corpus bytes at [%d,%d) do not match segment text — corpus and segments.jsonl are out of sync, re-run cmd/segment", s.ID, s.Start, s.End)
	}
}

func stripFrontmatter(s string) string {
	if !strings.HasPrefix(s, "---\n") {
		return s
	}
	rest := s[4:]
	for {
		nl := strings.IndexByte(rest, '\n')
		if nl < 0 {
			return s
		}
		line := rest[:nl]
		rest = rest[nl+1:]
		if strings.TrimRight(line, " \t") == "---" {
			return strings.TrimPrefix(rest, "\n")
		}
	}
}

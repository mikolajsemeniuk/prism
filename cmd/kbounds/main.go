// Command kbounds is step 4 of the decomposition pipeline (see CLAUDE.md
// "Decomposition methodology") — the pre-registered controls exam.
//
// The exam asks, at every cut k of every model's clustering, four questions
// whose answers we know in advance from the manual corpus analysis:
//
//	same-cluster:      the byte-identical Anthropic child-safety passage
//	                   (3 files) must land in ONE cluster;
//	same-cluster:      the near-verbatim Cursor↔Windsurf UX checklist;
//	same-cluster:      the byte-identical Codex rg-preference line (2 files);
//	different-cluster: the child-safety passage and the rg line must NOT
//	                   share a cluster.
//
// A cut k is ADMISSIBLE only if EVERY model passes ALL four questions at k.
// Coarse cuts fail the negative control (unrelated topics get mixed), fine
// cuts fail the positive controls (near-verbatim duplicates get torn apart),
// so the admissible set is bounded from below AND above — the k window. If no
// cut is admissible the command exits non-zero: the pipeline must not produce
// results in that state.
//
// Inputs: artefacts/segments.jsonl + every artefacts/clusters-*.json
// Output: artefacts/kbounds.json —
//
//	{"lo","hi","per_k":[{"k","pass":{model:bool},"pass_all"}]}
package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Pre-registered control needles (verbatim substrings of corpus bodies).
type control struct {
	name    string
	expect  string // "same" | "different"
	needleA string
	needleB string // "different" only
}

var controlSet = []control{
	{"anthropic-child-safety", "same", "Claude cares deeply about child safety and exercises special caution", ""},
	{"cursor-windsurf-ux", "same", "imbued with best UX practices", ""},
	{"codex-rg-preference", "same", "prefer using `rg` or `rg --files` respectively", ""},
	{"child-safety-vs-rg", "different",
		"Claude cares deeply about child safety and exercises special caution",
		"prefer using `rg` or `rg --files` respectively"},
}

type seg struct {
	ID     string `json:"id"`
	Source string `json:"source"`
	Text   string `json:"text"`
}

type clusters struct {
	Model string   `json:"model"`
	IDs   []string `json:"ids"`
	Grid  []struct {
		K      int   `json:"k"`
		Labels []int `json:"labels"`
	} `json:"grid"`
}

type perK struct {
	K       int             `json:"k"`
	Pass    map[string]bool `json:"pass"`
	PassAll bool            `json:"pass_all"`
}

func main() {
	log.SetFlags(0)
	segsPath := flag.String("segments", "artefacts/segments.jsonl", "input segments")
	inDir := flag.String("indir", "artefacts", "directory with clusters-*.json")
	out := flag.String("out", "artefacts/kbounds.json", "output path")
	flag.Parse()

	byID := map[string]seg{}
	f, err := os.Open(*segsPath)
	if err != nil {
		log.Fatalf("segments: %v (run cmd/segment first)", err)
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

	paths, _ := filepath.Glob(filepath.Join(*inDir, "clusters-*.json"))
	sort.Strings(paths)
	if len(paths) < 2 {
		log.Fatalf("need >=2 clusters-*.json in %s (run cmd/cluster per model), found %d", *inDir, len(paths))
	}
	var all []clusters
	for _, p := range paths {
		raw, err := os.ReadFile(p)
		if err != nil {
			log.Fatal(err)
		}
		var c clusters
		if err := json.Unmarshal(raw, &c); err != nil {
			log.Fatalf("%s: %v", p, err)
		}
		all = append(all, c)
	}

	// shared k set across models
	kCount := map[int]int{}
	for _, c := range all {
		for _, g := range c.Grid {
			kCount[g.K]++
		}
	}
	var ks []int
	for k, cnt := range kCount {
		if cnt == len(all) {
			ks = append(ks, k)
		}
	}
	sort.Ints(ks)

	var rows []perK
	lo, hi := 0, 0
	for _, k := range ks {
		row := perK{K: k, Pass: map[string]bool{}, PassAll: true}
		for _, c := range all {
			p := passControls(c, k, byID)
			row.Pass[c.Model] = p
			row.PassAll = row.PassAll && p
		}
		if row.PassAll {
			if lo == 0 {
				lo = k
			}
			hi = k
		}
		rows = append(rows, row)
	}

	models := make([]string, len(all))
	for i, c := range all {
		models[i] = c.Model
	}
	fmt.Printf("%-6s", "k")
	for _, m := range models {
		fmt.Printf("%-14s", m)
	}
	fmt.Println()
	for _, r := range rows {
		fmt.Printf("%-6d", r.K)
		for _, m := range models {
			s := "pass"
			if !r.Pass[m] {
				s = "FAIL"
			}
			fmt.Printf("%-14s", s)
		}
		fmt.Println()
	}
	if lo == 0 {
		log.Fatal("no admissible cut: every k fails at least one control — do not proceed")
	}
	fmt.Printf("\nadmissible window: k in [%d, %d]\n", lo, hi)

	raw, _ := json.MarshalIndent(map[string]any{"lo": lo, "hi": hi, "per_k": rows}, "", "  ")
	if err := os.WriteFile(*out, append(raw, '\n'), 0o644); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("wrote %s\n", *out)
}

// passControls answers all four exam questions for one model at one cut.
func passControls(c clusters, k int, byID map[string]seg) bool {
	var labels []int
	for _, g := range c.Grid {
		if g.K == k {
			labels = g.Labels
		}
	}
	if labels == nil {
		return false
	}
	find := func(needle string) (cl map[int]bool, files map[string]bool) {
		cl, files = map[int]bool{}, map[string]bool{}
		for i, id := range c.IDs {
			s, ok := byID[id]
			if ok && strings.Contains(s.Text, needle) {
				cl[labels[i]] = true
				files[s.Source] = true
			}
		}
		return
	}
	for _, ctl := range controlSet {
		switch ctl.expect {
		case "same":
			cl, files := find(ctl.needleA)
			if len(files) < 2 || len(cl) != 1 {
				return false
			}
		case "different":
			ca, _ := find(ctl.needleA)
			cb, _ := find(ctl.needleB)
			if len(ca) == 0 || len(cb) == 0 {
				return false
			}
			for x := range ca {
				if cb[x] {
					return false
				}
			}
		}
	}
	return true
}

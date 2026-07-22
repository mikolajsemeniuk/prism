// Command gentex walks the derived pipeline outputs and emits the paper's
// generated LaTeX fragments (roadmap step 6). The article never contains a
// hand-copied number: it inputs these files, and these files are never edited
// by hand.
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/mikolajsemeniuk/prism/pkg/cluster"
	"github.com/mikolajsemeniuk/prism/pkg/config"
	"github.com/mikolajsemeniuk/prism/pkg/segment"
	"github.com/mikolajsemeniuk/prism/pkg/tex"
)

func main() {
	log.SetFlags(0)
	cfgPath := flag.String("config", "pipeline.json", "pre-registered pipeline config")
	segsPath := flag.String("segments", "derived/segments.jsonl", "input segments")
	clustersDir := flag.String("clustersdir", "derived/clusters", "clustering results directory")
	derivedDir := flag.String("derived", "derived", "directory with stability/controls/matrix json")
	outDir := flag.String("outdir", "paper", "output directory for *.gen.tex")
	flag.Parse()

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		log.Fatalf("config: %v", err)
	}
	segs, err := segment.LoadJSONL(*segsPath)
	if err != nil {
		log.Fatalf("segments: %v", err)
	}
	results, err := cluster.LoadResultsDir(*clustersDir)
	if err != nil {
		log.Fatalf("clusters: %v", err)
	}
	inputs := map[string]string{
		"segments": fileHash(*segsPath),
	}
	for _, name := range []string{"stability.json", "controls.json", "matrix.json"} {
		inputs[name] = fileHash(filepath.Join(*derivedDir, name))
	}
	if err := os.MkdirAll(*outDir, 0o755); err != nil {
		log.Fatal(err)
	}

	stab := loadStability(*derivedDir)
	writeCorpusStats(*outDir, cfg, inputs, segs)
	writeClusterGrid(*outDir, cfg, inputs, results, stab)
	writeStability(*outDir, cfg, inputs, results, stab)
	writeControls(*outDir, cfg, inputs, *derivedDir)
	writeMatrix(*outDir, cfg, inputs, *derivedDir)
	n := 5
	if writeTaxonomy(*outDir, cfg, inputs, *derivedDir) {
		n++
	}
	fmt.Printf("wrote %d .gen.tex files to %s\n", n, *outDir)
}

func writeTaxonomy(dir string, cfg config.Config, inputs map[string]string, derivedDir string) bool {
	path := filepath.Join(derivedDir, "taxonomy.json")
	if _, err := os.Stat(path); err != nil {
		return false
	}
	var tx struct {
		ComponentK int `json:"component_k"`
		Counts     map[string]int
		Entries    []struct {
			ID        int      `json:"id"`
			Category  string   `json:"category"`
			Size      int      `json:"size"`
			Stability float64  `json:"stability"`
			Products  []string `json:"products"`
			Terms     []string `json:"terms"`
		} `json:"entries"`
	}
	loadJSON(path, &tx)
	inputs["taxonomy.json"] = fileHash(path)
	var rows [][]string
	for _, e := range tx.Entries {
		if e.Category != "component" {
			continue
		}
		terms := e.Terms
		if len(terms) > 4 {
			terms = terms[:4]
		}
		rows = append(rows, []string{
			fmt.Sprintf("C%d", e.ID), fmt.Sprintf("%d", e.Size),
			fmt.Sprintf("%d", len(e.Products)), fmt.Sprintf("%.2f", e.Stability),
			tex.Escape(strings.Join(terms, ", ")),
		})
	}
	writeTex(dir, "taxonomy.gen.tex", cfg, inputs, func(f *os.File) {
		tex.Table(f,
			fmt.Sprintf("Data-driven components at the consensus cut ($k=%d$): clusters reproduced by both embedding models (containment $\\geq$ %.2f) spanning $\\geq$ %d products. Distinctive terms are tf-idf over clusters; names are assigned in a separate interpretive pass.",
				tx.ComponentK, cfg.Taxonomy.StabilityMin, cfg.Taxonomy.MinProducts),
			"tab:taxonomy", "lrrrl",
			[]string{"Cluster", "Size", "Products", "Stability", "Distinctive terms"}, rows)
	})
	return true
}

type stability struct {
	Pairs []struct {
		A       string  `json:"a"`
		B       string  `json:"b"`
		NShared int     `json:"n_shared"`
		ARI     float64 `json:"ari"`
		NMI     float64 `json:"nmi"`
	} `json:"pairs"`
	Consensus *struct {
		K    int `json:"k"`
		PerK []struct {
			K            int     `json:"k"`
			MeanARI      float64 `json:"mean_ari"`
			ControlsPass bool    `json:"controls_pass"`
		} `json:"per_k"`
	} `json:"consensus"`
}

func loadStability(derivedDir string) stability {
	var s stability
	loadJSON(filepath.Join(derivedDir, "stability.json"), &s)
	return s
}

func writeCorpusStats(dir string, cfg config.Config, inputs map[string]string, segs []segment.Segment) {
	type agg struct {
		product, class, tier string
		count                int
		chars                map[segment.Layer]int
	}
	bySource := map[string]*agg{}
	var sources []string
	for _, s := range segs {
		a := bySource[s.Source]
		if a == nil {
			a = &agg{product: s.Product, class: s.ProductClass, tier: s.Tier, chars: map[segment.Layer]int{}}
			bySource[s.Source] = a
			sources = append(sources, s.Source)
		}
		a.count++
		a.chars[s.Layer] += len(s.Text)
	}
	sort.Strings(sources)
	var rows [][]string
	for _, src := range sources {
		a := bySource[src]
		total := 0
		for _, c := range a.chars {
			total += c
		}
		pct := func(l segment.Layer) string {
			if total == 0 {
				return "--"
			}
			return fmt.Sprintf("%.0f\\%%", 100*float64(a.chars[l])/float64(total))
		}
		rows = append(rows, []string{
			tex.Escape(a.product), tex.Escape(a.class), tex.Escape(a.tier),
			fmt.Sprintf("%d", a.count),
			pct(segment.Behavioral), pct(segment.ToolSchema), pct(segment.FewShot),
		})
	}
	writeTex(dir, "corpus-stats.gen.tex", cfg, inputs, func(f *os.File) {
		tex.Table(f,
			"Segmented corpus composition per artifact: segment count and character share of the behavioral, tool-schema, and few-shot layers (v1 heuristic layer tags).",
			"tab:corpus-stats", "lllrrrr",
			[]string{"Product", "Class", "Tier", "Segs", "Behav.", "Schema", "Few-shot"}, rows)
	})
}

func writeClusterGrid(dir string, cfg config.Config, inputs map[string]string, results []*cluster.Result, stab stability) {
	ks := map[int]bool{}
	for _, r := range results {
		for _, g := range r.Grid {
			ks[g.K] = true
		}
	}
	var sortedK []int
	for k := range ks {
		sortedK = append(sortedK, k)
	}
	sort.Ints(sortedK)
	consensusARI := map[int]string{}
	componentK := -1
	if stab.Consensus != nil {
		componentK = stab.Consensus.K
		for _, p := range stab.Consensus.PerK {
			cell := fmt.Sprintf("%.4f", p.MeanARI)
			if !p.ControlsPass {
				cell += "$^\\dagger$"
			}
			if p.K == componentK {
				cell = "\\textbf{" + cell + "}"
			}
			consensusARI[p.K] = cell
		}
	}
	header := []string{"$k$"}
	for _, r := range results {
		header = append(header, tex.Escape(r.Model))
	}
	if stab.Consensus != nil {
		header = append(header, "cross-model ARI")
	}
	var rows [][]string
	for _, k := range sortedK {
		row := []string{fmt.Sprintf("%d", k)}
		for _, r := range results {
			cell := "--"
			for _, g := range r.Grid {
				if g.K == k {
					cell = fmt.Sprintf("%.4f", g.Silhouette)
					if !g.ControlsPass {
						cell += "$^\\dagger$"
					}
					if k == r.ChosenK {
						cell = "\\textbf{" + cell + "}"
					}
				}
			}
			row = append(row, cell)
		}
		if stab.Consensus != nil {
			cell, ok := consensusARI[k]
			if !ok {
				cell = "--"
			}
			row = append(row, cell)
		}
		rows = append(rows, row)
	}
	colspec := "r" + repeat("r", len(header)-1)
	writeTex(dir, "cluster-grid.gen.tex", cfg, inputs, func(f *os.File) {
		tex.Table(f,
			"Two-level readout of the dendrogram across the pre-registered $k$ grid: per-model mean silhouette (bold: fine cut, max silhouette subject to controls) and cross-model agreement (bold: component cut, max mean pairwise ARI subject to controls). $\\dagger$ marks cuts violating a control assertion.",
			"tab:cluster-grid", colspec, header, rows)
	})
}

func writeStability(dir string, cfg config.Config, inputs map[string]string, results []*cluster.Result, stab stability) {
	var rows [][]string
	for _, p := range stab.Pairs {
		rows = append(rows, []string{
			tex.Escape(p.A), tex.Escape(p.B), fmt.Sprintf("%d", p.NShared),
			fmt.Sprintf("%.3f", p.ARI), fmt.Sprintf("%.3f", p.NMI),
		})
	}
	for _, r := range results {
		if r.Bootstrap == nil {
			continue
		}
		rows = append(rows, []string{
			tex.Escape(r.Model), "(bootstrap)", fmt.Sprintf("%d", r.Bootstrap.Rounds),
			fmt.Sprintf("%.3f", r.Bootstrap.MeanARI),
			fmt.Sprintf("[%.3f, %.3f]", r.Bootstrap.P2_5, r.Bootstrap.P97_5),
		})
	}
	if stab.Consensus != nil {
		best := ""
		for _, p := range stab.Consensus.PerK {
			if p.K == stab.Consensus.K {
				best = fmt.Sprintf("%.3f", p.MeanARI)
			}
		}
		rows = append(rows, []string{
			"consensus", "(component cut)", fmt.Sprintf("%d", stab.Consensus.K), best, "--",
		})
	}
	writeTex(dir, "stability.gen.tex", cfg, inputs, func(f *os.File) {
		tex.Table(f,
			"Cluster stability: pairwise cross-model agreement at the fine cuts (ARI/NMI), per-model bootstrap resampling at the fine $k$ (mean ARI, 95\\% percentile interval), and the consensus component cut ($k$, cross-model ARI).",
			"tab:stability", "llrrr",
			[]string{"A", "B", "$n$ / $k$", "ARI", "NMI / CI"}, rows)
	})
}

func writeControls(dir string, cfg config.Config, inputs map[string]string, derivedDir string) {
	var ctl struct {
		Pass     bool `json:"pass"`
		Controls []struct {
			Name     string `json:"name"`
			Expect   string `json:"expect"`
			Pass     bool   `json:"pass"`
			PerModel []struct {
				Model string `json:"model"`
				Pass  bool   `json:"pass"`
			} `json:"per_model"`
		} `json:"controls"`
	}
	loadJSON(filepath.Join(derivedDir, "controls.json"), &ctl)
	var rows [][]string
	for _, c := range ctl.Controls {
		nPass := 0
		for _, m := range c.PerModel {
			if m.Pass {
				nPass++
			}
		}
		status := "\\textbf{FAIL}"
		if c.Pass {
			status = "pass"
		}
		rows = append(rows, []string{
			tex.Escape(c.Name), tex.Escape(c.Expect),
			fmt.Sprintf("%d/%d", nPass, len(c.PerModel)), status,
		})
	}
	writeTex(dir, "controls.gen.tex", cfg, inputs, func(f *os.File) {
		tex.Table(f,
			"Pre-registered clustering control assertions (positive controls: known duplicated passages must co-cluster; negative control: unrelated categories must not).",
			"tab:controls", "llrr",
			[]string{"Control", "Expect", "Models pass", "Status"}, rows)
	})
}

func writeMatrix(dir string, cfg config.Config, inputs map[string]string, derivedDir string) {
	var mx struct {
		Model      string `json:"model"`
		FineK      int    `json:"fine_k"`
		ComponentK int    `json:"component_k"`
		Clusters   []struct {
			Cluster int            `json:"cluster"`
			Size    int            `json:"size"`
			Classes map[string]int `json:"classes"`
		} `json:"clusters"`
	}
	loadJSON(filepath.Join(derivedDir, "matrix.json"), &mx)
	classes := []string{"chat", "coding-agent", "app-builder"}
	var rows [][]string
	limit := 15
	for i, c := range mx.Clusters {
		if i >= limit {
			break
		}
		row := []string{fmt.Sprintf("C%d", c.Cluster), fmt.Sprintf("%d", c.Size)}
		for _, cl := range classes {
			row = append(row, fmt.Sprintf("%d", c.Classes[cl]))
		}
		rows = append(rows, row)
	}
	writeTex(dir, "component-matrix.gen.tex", cfg, inputs, func(f *os.File) {
		tex.Table(f,
			fmt.Sprintf("Largest clusters at the consensus component cut ($k=%d$, model %s; fine cut $k=%d$), stratified by product class (findings F4). Cluster naming happens in the codebook-validation step and never feeds back into clustering.", mx.ComponentK, tex.Escape(mx.Model), mx.FineK),
			"tab:component-matrix", "lrrrr",
			[]string{"Cluster", "Size", "Chat", "Coding", "App-builder"}, rows)
	})
}

func writeTex(dir, name string, cfg config.Config, inputs map[string]string, body func(*os.File)) {
	f, err := os.Create(filepath.Join(dir, name))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	tex.Header(f, "cmd/gentex", cfg.Hash, inputs)
	body(f)
}

func loadJSON(path string, v any) {
	raw, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("%s: %v (run cmd/matrix first)", path, err)
	}
	if err := json.Unmarshal(raw, v); err != nil {
		log.Fatalf("%s: %v", path, err)
	}
}

func fileHash(path string) string {
	raw, err := os.ReadFile(path)
	if err != nil {
		return "missing"
	}
	sum := sha256.Sum256(raw)
	return hex.EncodeToString(sum[:])
}

func repeat(s string, n int) string {
	out := ""
	for i := 0; i < n; i++ {
		out += s
	}
	return out
}

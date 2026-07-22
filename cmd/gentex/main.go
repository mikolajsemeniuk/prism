// Command gentex is step 8 of the decomposition pipeline (see CLAUDE.md
// "Decomposition methodology") — the ONLY step that writes under paper/.
//
// It is pure formatting: every number comes verbatim from artefacts/*.json,
// nothing is computed here, and each emitted file starts with a provenance
// header (generator + sha256 of every input) and a do-not-edit notice, so the
// article can never contain a hand-copied or hand-tweaked value.
//
// Emitted tables:
//
//	paper/corpus-stats.gen.tex — per-product segment counts and layer shares
//	                             (from segments.jsonl);
//	paper/kwindow.gen.tex      — the controls exam per k per model merged
//	                             with cross-model ARI: admissible window and
//	                             the component cut k* (from kbounds.json +
//	                             admit.json);
//	paper/taxonomy.gen.tex     — components at k* (from taxonomy.json), with
//	                             names merged from the hand-edited
//	                             analysis/component-names.json (interpretive
//	                             step E; missing names render as "—"). Only
//	                             emitted when taxonomy.json exists.
package main

import (
	"bufio"
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
)

func main() {
	log.SetFlags(0)
	dir := flag.String("artefacts", "artefacts", "artefacts directory")
	namesPath := flag.String("names", "analysis/component-names.json", "hand-edited component names (step E)")
	outDir := flag.String("outdir", "paper", "output directory")
	flag.Parse()

	if err := os.MkdirAll(*outDir, 0o755); err != nil {
		log.Fatal(err)
	}
	n := 0
	writeCorpusStats(*dir, *outDir)
	n++
	writeKWindow(*dir, *outDir)
	n++
	if writeTaxonomy(*dir, *namesPath, *outDir) {
		n++
	}
	fmt.Printf("wrote %d .gen.tex files to %s\n", n, *outDir)
}

// ----- corpus-stats.gen.tex -----

func writeCorpusStats(dir, outDir string) {
	segsPath := filepath.Join(dir, "segments.jsonl")
	f, err := os.Open(segsPath)
	if err != nil {
		log.Fatalf("%v (run cmd/segment first)", err)
	}
	defer f.Close()
	type agg struct {
		count int
		chars map[string]int
	}
	byProduct := map[string]*agg{}
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 0, 1<<20), 64<<20)
	for sc.Scan() {
		var s struct {
			Product string `json:"product"`
			Layer   string `json:"layer"`
			Text    string `json:"text"`
		}
		if err := json.Unmarshal(sc.Bytes(), &s); err != nil {
			log.Fatal(err)
		}
		a := byProduct[s.Product]
		if a == nil {
			a = &agg{chars: map[string]int{}}
			byProduct[s.Product] = a
		}
		a.count++
		a.chars[s.Layer] += len(s.Text)
	}
	var products []string
	for p := range byProduct {
		products = append(products, p)
	}
	sort.Strings(products)

	var rows [][]string
	for _, p := range products {
		a := byProduct[p]
		total := 0
		for _, c := range a.chars {
			total += c
		}
		pct := func(layer string) string {
			if total == 0 {
				return "--"
			}
			return fmt.Sprintf("%.0f\\%%", 100*float64(a.chars[layer])/float64(total))
		}
		rows = append(rows, []string{
			escape(p), fmt.Sprintf("%d", a.count),
			pct("behavioral"), pct("tool-definition-schema"), pct("few-shot"),
		})
	}
	writeTable(filepath.Join(outDir, "corpus-stats.gen.tex"),
		[]string{segsPath},
		"Segmented corpus composition per product: segment count and character share of the behavioral, tool-schema, and few-shot layers (heuristic layer tags; only the behavioral layer is embedded).",
		"tab:corpus-stats", "lrrrr",
		[]string{"Product", "Segs", "Behav.", "Schema", "Few-shot"}, rows)
}

// ----- kwindow.gen.tex -----

func writeKWindow(dir, outDir string) {
	kbPath := filepath.Join(dir, "kbounds.json")
	adPath := filepath.Join(dir, "admit.json")
	var kb struct {
		Lo   int `json:"lo"`
		Hi   int `json:"hi"`
		PerK []struct {
			K    int             `json:"k"`
			Pass map[string]bool `json:"pass"`
		} `json:"per_k"`
	}
	loadJSON(kbPath, &kb)
	var ad struct {
		KStar int     `json:"k_star"`
		ARI   float64 `json:"ari"`
		PerK  []struct {
			K       int     `json:"k"`
			MeanARI float64 `json:"mean_ari"`
		} `json:"per_k"`
	}
	loadJSON(adPath, &ad)

	ariAt := map[int]float64{}
	for _, r := range ad.PerK {
		ariAt[r.K] = r.MeanARI
	}
	var models []string
	if len(kb.PerK) > 0 {
		for m := range kb.PerK[0].Pass {
			models = append(models, m)
		}
		sort.Strings(models)
	}

	header := []string{"$k$"}
	for _, m := range models {
		header = append(header, escape(m))
	}
	header = append(header, "cross-model ARI")
	var rows [][]string
	for _, r := range kb.PerK {
		row := []string{fmt.Sprintf("%d", r.K)}
		for _, m := range models {
			if r.Pass[m] {
				row = append(row, "pass")
			} else {
				row = append(row, "$\\dagger$")
			}
		}
		cell := "--"
		if a, ok := ariAt[r.K]; ok {
			cell = fmt.Sprintf("%.4f", a)
			if r.K == ad.KStar {
				cell = "\\textbf{" + cell + "}"
			}
		}
		row = append(row, cell)
		rows = append(rows, row)
	}
	writeTable(filepath.Join(outDir, "kwindow.gen.tex"),
		[]string{kbPath, adPath},
		fmt.Sprintf("Pre-registered controls exam and cross-model agreement across the $k$ grid. $\\dagger$ marks cuts failing a control assertion; the admissible window is $k \\in [%d, %d]$ and the component cut $k^*=%d$ maximizes mean cross-model ARI within it (bold, %.3f).",
			kb.Lo, kb.Hi, ad.KStar, ad.ARI),
		"tab:kwindow", "r"+strings.Repeat("c", len(models))+"r", header, rows)
}

// ----- taxonomy.gen.tex -----

func writeTaxonomy(dir, namesPath, outDir string) bool {
	txPath := filepath.Join(dir, "taxonomy.json")
	if _, err := os.Stat(txPath); err != nil {
		return false
	}
	var tx struct {
		KStar   int `json:"k_star"`
		Entries []struct {
			ID        int      `json:"id"`
			Category  string   `json:"category"`
			Size      int      `json:"size"`
			Stability float64  `json:"stability"`
			Products  []string `json:"products"`
			Terms     []string `json:"terms"`
		} `json:"entries"`
	}
	loadJSON(txPath, &tx)

	// step E: hand-edited names, merged at format time only. Non-object
	// values (e.g. a top-level "_comment" string) are ignored.
	names := map[string]struct {
		Name string `json:"name"`
	}{}
	inputs := []string{txPath}
	if raw, err := os.ReadFile(namesPath); err == nil {
		var loose map[string]json.RawMessage
		if err := json.Unmarshal(raw, &loose); err != nil {
			log.Fatalf("%s: %v", namesPath, err)
		}
		for k, v := range loose {
			var nm struct {
				Name string `json:"name"`
			}
			if json.Unmarshal(v, &nm) == nil && nm.Name != "" {
				names[k] = nm
			}
		}
		inputs = append(inputs, namesPath)
	}

	var rows [][]string
	for _, e := range tx.Entries {
		if e.Category != "component" {
			continue
		}
		name := "---"
		if nm, ok := names[fmt.Sprintf("C%d", e.ID)]; ok && nm.Name != "" {
			name = escape(nm.Name)
		}
		terms := e.Terms
		if len(terms) > 4 {
			terms = terms[:4]
		}
		rows = append(rows, []string{
			name, fmt.Sprintf("C%d", e.ID), fmt.Sprintf("%d", e.Size),
			fmt.Sprintf("%d", len(e.Products)), fmt.Sprintf("%.2f", e.Stability),
			escape(strings.Join(terms, ", ")),
		})
	}
	writeTable(filepath.Join(outDir, "taxonomy.gen.tex"), inputs,
		fmt.Sprintf("Data-driven components at the consensus cut ($k^*=%d$): clusters reproduced by every embedding model (containment $\\geq 0.5$) spanning $\\geq 3$ products. Distinctive terms are tf-idf over clusters. Names are an interpretive layer maintained separately (analysis/component-names.json) and never influence the mechanical pipeline.", tx.KStar),
		"tab:taxonomy", "llrrrl",
		[]string{"Name", "Cluster", "Size", "Products", "Stability", "Distinctive terms"}, rows)
	return true
}

// ----- shared helpers -----

func writeTable(path string, inputs []string, caption, label, colspec string, header []string, rows [][]string) {
	var b strings.Builder
	b.WriteString("% GENERATED FILE — do not edit by hand (CLAUDE.md rule: *.gen.tex are never hand-edited).\n")
	b.WriteString("% generator: cmd/gentex\n")
	for _, in := range inputs {
		fmt.Fprintf(&b, "%% input %s sha256: %s\n", in, fileHash(in))
	}
	b.WriteString("\\begin{table}[t]\n\\centering\n\\small\n")
	fmt.Fprintf(&b, "\\begin{tabular}{%s}\n\\toprule\n", colspec)
	fmt.Fprintf(&b, "%s \\\\\n\\midrule\n", strings.Join(header, " & "))
	for _, r := range rows {
		fmt.Fprintf(&b, "%s \\\\\n", strings.Join(r, " & "))
	}
	b.WriteString("\\bottomrule\n\\end{tabular}\n")
	fmt.Fprintf(&b, "\\caption{%s}\n\\label{%s}\n\\end{table}\n", caption, label)
	if err := os.WriteFile(path, []byte(b.String()), 0o644); err != nil {
		log.Fatal(err)
	}
}

func loadJSON(path string, v any) {
	raw, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("%v (run the producing step first)", err)
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

var escaper = strings.NewReplacer(
	`\`, `\textbackslash{}`, `&`, `\&`, `%`, `\%`, `$`, `\$`, `#`, `\#`,
	`_`, `\_`, `{`, `\{`, `}`, `\}`, `~`, `\textasciitilde{}`, `^`, `\textasciicircum{}`,
)

func escape(s string) string { return escaper.Replace(s) }

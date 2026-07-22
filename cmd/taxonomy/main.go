// Command taxonomy is step 6 of the decomposition pipeline (see CLAUDE.md
// "Decomposition methodology") — it derives the data-driven component
// taxonomy from the component cut, entirely mechanically (steps A–D):
//
//	A. stability filter — keep clusters whose members the OTHER embedding
//	   model also groups together: best-match containment |A∩B|/|A| >= 0.5.
//	   Containment, not Jaccard: Jaccard conflates instability with
//	   granularity mismatch (a cluster the other model merely splits in two
//	   is co-grouped, not unstable).
//	B. componentness filter — a component must span >= 3 distinct products
//	   (fewer → "product-specific": a feature, not an industry component)
//	   and hold >= 5 segments (fewer → "micro": duplicate groups).
//	C. characterization — medoid (member with max mean cosine similarity to
//	   the rest, i.e. the most typical real sentence) plus distinctive terms
//	   (tf-idf with clusters as documents: frequent here, rare elsewhere).
//	D. emission — machine-readable taxonomy.json + human-readable taxonomy.md.
//
// Step E (naming/describing components) is deliberately OUTSIDE this command:
// names live in the hand-edited analysis/component-names.json and are merged
// only at paper-generation time (cmd/gentex), so the interpretive layer can
// never influence these numbers.
//
// Inputs: artefacts/segments.jsonl, artefacts/clusters-*.json,
//
//	artefacts/embeddings-<primary>.jsonl, artefacts/admit.json
//
// Output: artefacts/taxonomy.json + artefacts/taxonomy.md
package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Pre-registered thresholds (see CLAUDE.md; calibrated on the pilot).
const (
	stabilityMin = 0.5
	minProducts  = 3
	minSize      = 5
	topTerms     = 8
	nExemplars   = 3
)

type seg struct {
	ID      string `json:"id"`
	Source  string `json:"source"`
	Product string `json:"product"`
	Text    string `json:"text"`
}

type clusters struct {
	Model string   `json:"model"`
	IDs   []string `json:"ids"`
	Grid  []struct {
		K      int   `json:"k"`
		Labels []int `json:"labels"`
	} `json:"grid"`
}

type entry struct {
	ID          int      `json:"id"`
	Category    string   `json:"category"` // component | product-specific | unstable | micro
	Size        int      `json:"size"`
	Stability   float64  `json:"stability"` // min over other models of best-match containment
	Products    []string `json:"products"`
	Terms       []string `json:"terms,omitempty"`
	MedoidID    string   `json:"medoid_id,omitempty"`
	MedoidText  string   `json:"medoid_text,omitempty"`
	ExemplarIDs []string `json:"exemplar_ids,omitempty"`
}

func main() {
	log.SetFlags(0)
	dir := flag.String("artefacts", "artefacts", "artefacts directory")
	flag.Parse()

	byID := loadSegments(filepath.Join(*dir, "segments.jsonl"))
	all := loadClusters(*dir)
	kStar := loadKStar(filepath.Join(*dir, "admit.json"))
	primary := all[0]
	labels := labelsAt(primary, kStar)
	if labels == nil {
		log.Fatalf("model %s has no labels at k*=%d", primary.Model, kStar)
	}

	// cluster members (primary model, component cut)
	members := map[int][]string{}
	for i, id := range primary.IDs {
		members[labels[i]] = append(members[labels[i]], id)
	}
	for _, ids := range members {
		sort.Strings(ids)
	}

	// A: stability — min over other models of best-match containment
	stability := map[int]float64{}
	for cl := range members {
		stability[cl] = 1
	}
	for _, other := range all[1:] {
		oLabels := labelsAt(other, kStar)
		if oLabels == nil {
			log.Fatalf("model %s has no labels at k*=%d", other.Model, kStar)
		}
		oMembers := map[int]map[string]bool{}
		for i, id := range other.IDs {
			if oMembers[oLabels[i]] == nil {
				oMembers[oLabels[i]] = map[string]bool{}
			}
			oMembers[oLabels[i]][id] = true
		}
		for cl, ids := range members {
			best := 0.0
			for _, set := range oMembers {
				inter := 0
				for _, id := range ids {
					if set[id] {
						inter++
					}
				}
				if c := float64(inter) / float64(len(ids)); c > best {
					best = c
				}
			}
			if best < stability[cl] {
				stability[cl] = best
			}
		}
	}

	// C prerequisites: embeddings of the primary model (medoids) and
	// per-cluster term frequencies vs cluster document frequency (tf-idf).
	vec := loadEmbeddings(filepath.Join(*dir, "embeddings-"+sanitize(primary.Model)+".jsonl"))
	df := map[string]int{}
	tf := map[int]map[string]int{}
	for cl, ids := range members {
		tf[cl] = map[string]int{}
		seen := map[string]bool{}
		for _, id := range ids {
			for _, tok := range tokenize(byID[id].Text) {
				tf[cl][tok]++
				if !seen[tok] {
					seen[tok] = true
					df[tok]++
				}
			}
		}
	}

	// B + C + D
	var entries []entry
	counts := map[string]int{}
	for cl, ids := range members {
		e := entry{ID: cl, Size: len(ids), Stability: stability[cl]}
		prodSet := map[string]bool{}
		for _, id := range ids {
			prodSet[byID[id].Product] = true
		}
		for p := range prodSet {
			e.Products = append(e.Products, p)
		}
		sort.Strings(e.Products)
		switch {
		case e.Size < minSize:
			e.Category = "micro"
		case e.Stability < stabilityMin:
			e.Category = "unstable"
		case len(e.Products) < minProducts:
			e.Category = "product-specific"
		default:
			e.Category = "component"
		}
		counts[e.Category]++
		if e.Category == "component" || e.Category == "product-specific" {
			e.Terms = distinctiveTerms(tf[cl], df, float64(len(members)))
			e.MedoidID, e.ExemplarIDs = medoid(ids, vec)
			if e.MedoidID != "" {
				e.MedoidText = clip(byID[e.MedoidID].Text, 300)
			}
		}
		entries = append(entries, e)
	}
	sort.Slice(entries, func(i, j int) bool {
		if len(entries[i].Products) != len(entries[j].Products) {
			return len(entries[i].Products) > len(entries[j].Products)
		}
		if entries[i].Size != entries[j].Size {
			return entries[i].Size > entries[j].Size
		}
		return entries[i].ID < entries[j].ID
	})

	out := map[string]any{
		"model":  primary.Model,
		"k_star": kStar,
		"params": map[string]any{
			"stability_min": stabilityMin, "min_products": minProducts,
			"min_size": minSize, "top_terms": topTerms, "exemplars": nExemplars,
		},
		"counts":  counts,
		"entries": entries,
	}
	raw, _ := json.MarshalIndent(out, "", "  ")
	jsonPath := filepath.Join(*dir, "taxonomy.json")
	if err := os.WriteFile(jsonPath, append(raw, '\n'), 0o644); err != nil {
		log.Fatal(err)
	}
	if err := writeMarkdown(filepath.Join(*dir, "taxonomy.md"), primary.Model, kStar, counts, entries, byID); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("taxonomy at k*=%d (model %s): %d components, %d product-specific, %d unstable, %d micro\n",
		kStar, primary.Model, counts["component"], counts["product-specific"], counts["unstable"], counts["micro"])
	fmt.Printf("wrote %s and taxonomy.md\n", jsonPath)
}

func loadSegments(path string) map[string]seg {
	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("%v (run cmd/segment first)", err)
	}
	defer f.Close()
	byID := map[string]seg{}
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 0, 1<<20), 64<<20)
	for sc.Scan() {
		var s seg
		if err := json.Unmarshal(sc.Bytes(), &s); err != nil {
			log.Fatal(err)
		}
		byID[s.ID] = s
	}
	return byID
}

func loadClusters(dir string) []clusters {
	paths, _ := filepath.Glob(filepath.Join(dir, "clusters-*.json"))
	sort.Strings(paths)
	if len(paths) < 2 {
		log.Fatalf("need >=2 clusters-*.json in %s (run cmd/cluster per model)", dir)
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
	return all
}

func loadKStar(path string) int {
	raw, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("%v (run cmd/admit first)", err)
	}
	var a struct {
		KStar int `json:"k_star"`
	}
	if err := json.Unmarshal(raw, &a); err != nil {
		log.Fatal(err)
	}
	if a.KStar == 0 {
		log.Fatal("admit.json has no k_star")
	}
	return a.KStar
}

func loadEmbeddings(path string) map[string][]float32 {
	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("%v (run cmd/embed first)", err)
	}
	defer f.Close()
	vec := map[string][]float32{}
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 0, 1<<20), 64<<20)
	for sc.Scan() {
		var r struct {
			ID  string    `json:"id"`
			Vec []float32 `json:"vec"`
		}
		if json.Unmarshal(sc.Bytes(), &r) == nil && r.ID != "" {
			vec[r.ID] = r.Vec
		}
	}
	return vec
}

func labelsAt(c clusters, k int) []int {
	for _, g := range c.Grid {
		if g.K == k {
			return g.Labels
		}
	}
	return nil
}

// medoid returns the member with the highest mean cosine similarity to the
// rest of its cluster, plus the nExemplars most central members overall.
func medoid(ids []string, vec map[string][]float32) (string, []string) {
	type scored struct {
		id string
		s  float64
	}
	var xs []scored
	for _, a := range ids {
		va, ok := vec[a]
		if !ok {
			continue
		}
		var sum float64
		for _, b := range ids {
			if a == b {
				continue
			}
			if vb, ok := vec[b]; ok {
				var dot float64
				for d := range va {
					dot += float64(va[d]) * float64(vb[d])
				}
				sum += dot
			}
		}
		xs = append(xs, scored{a, sum})
	}
	if len(xs) == 0 {
		return "", nil
	}
	sort.Slice(xs, func(i, j int) bool {
		if xs[i].s != xs[j].s {
			return xs[i].s > xs[j].s
		}
		return xs[i].id < xs[j].id
	})
	var ex []string
	for i := 0; i < len(xs) && i < nExemplars; i++ {
		ex = append(ex, xs[i].id)
	}
	return xs[0].id, ex
}

// distinctiveTerms scores tokens by tf * ln(nClusters / (1 + df)): frequent
// in this cluster, present in few others.
func distinctiveTerms(tf map[string]int, df map[string]int, nClusters float64) []string {
	type ts struct {
		t string
		s float64
	}
	var xs []ts
	for t, f := range tf {
		if f < 2 {
			continue
		}
		xs = append(xs, ts{t, float64(f) * math.Log(nClusters/(1+float64(df[t])))})
	}
	sort.Slice(xs, func(i, j int) bool {
		if xs[i].s != xs[j].s {
			return xs[i].s > xs[j].s
		}
		return xs[i].t < xs[j].t
	})
	var out []string
	for i := 0; i < len(xs) && i < topTerms; i++ {
		out = append(out, xs[i].t)
	}
	return out
}

var stopwords = map[string]bool{
	"the": true, "and": true, "for": true, "with": true, "that": true,
	"this": true, "you": true, "your": true, "are": true, "not": true,
	"can": true, "will": true, "must": true, "should": true, "when": true,
	"any": true, "all": true, "use": true, "only": true, "may": true,
	"never": true, "always": true, "than": true, "from": true, "into": true,
	"has": true, "have": true, "was": true, "were": true, "been": true,
	"its": true, "it's": true, "does": true, "els": true, "each": true,
}

func tokenize(t string) []string {
	raw := strings.FieldsFunc(strings.ToLower(t), func(r rune) bool {
		return !('a' <= r && r <= 'z' || '0' <= r && r <= '9' || r == '-' || r == '_')
	})
	out := make([]string, 0, len(raw))
	for _, w := range raw {
		if len(w) >= 3 && !stopwords[w] {
			out = append(out, w)
		}
	}
	return out
}

func writeMarkdown(path, model string, kStar int, counts map[string]int, entries []entry, byID map[string]seg) error {
	var b strings.Builder
	fmt.Fprintf(&b, "# Data-driven component taxonomy (auto-generated by cmd/taxonomy — do not edit)\n\n")
	fmt.Fprintf(&b, "Component cut k*=%d, primary model %s. Thresholds: containment>=%.2f, products>=%d, size>=%d.\n",
		kStar, model, stabilityMin, minProducts, minSize)
	fmt.Fprintf(&b, "Names/definitions live in analysis/component-names.json (interpretive step E) and never feed back here.\n\n")
	fmt.Fprintf(&b, "Counts: %d components, %d product-specific, %d unstable, %d micro.\n\n",
		counts["component"], counts["product-specific"], counts["unstable"], counts["micro"])
	fmt.Fprintf(&b, "## Components (stable in every embedding model, >=%d products)\n\n", minProducts)
	for _, e := range entries {
		if e.Category != "component" {
			continue
		}
		fmt.Fprintf(&b, "### C%d\n\n", e.ID)
		fmt.Fprintf(&b, "- size: %d | products: %d (%s) | stability: %.2f\n", e.Size, len(e.Products), strings.Join(e.Products, ", "), e.Stability)
		fmt.Fprintf(&b, "- terms: %s\n", strings.Join(e.Terms, ", "))
		fmt.Fprintf(&b, "- medoid [%s]: %q\n", e.MedoidID, e.MedoidText)
		for _, ex := range e.ExemplarIDs {
			if ex != e.MedoidID {
				fmt.Fprintf(&b, "- exemplar [%s]: %q\n", ex, clip(byID[ex].Text, 160))
			}
		}
		b.WriteString("\n")
	}
	fmt.Fprintf(&b, "## Product-specific clusters (stable, <%d products)\n\n", minProducts)
	for _, e := range entries {
		if e.Category != "product-specific" {
			continue
		}
		fmt.Fprintf(&b, "- C%d: size %d, %s | terms: %s\n", e.ID, e.Size, strings.Join(e.Products, ", "), strings.Join(e.Terms, ", "))
	}
	return os.WriteFile(path, []byte(b.String()), 0o644)
}

func clip(s string, n int) string {
	s = strings.Join(strings.Fields(s), " ")
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}

func sanitize(s string) string {
	return strings.Map(func(r rune) rune {
		if r == ':' || r == '/' || r == ' ' {
			return '-'
		}
		return r
	}, s)
}

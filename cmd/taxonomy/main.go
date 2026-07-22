// Command taxonomy derives the data-driven component taxonomy from the
// component-cut clusters (steps A–D of the taxonomy methodology):
//
//	A. stability filter    — keep clusters whose members the second embedding
//	                         model also groups together: best-match containment
//	                         |A ∩ B|/|A| >= taxonomy.stability_min. Containment
//	                         (not Jaccard) because Jaccard punishes granularity
//	                         mismatches between models — a cluster split in two
//	                         by the other model is co-grouped, not unstable.
//	B. componentness filter — a component must span >= min_products products;
//	                         smaller spans are tagged product-specific
//	C. auto-characterization — medoid + exemplars (nearest to centroid) and
//	                         distinctive terms (tf-idf across clusters)
//	D. emission            — derived/taxonomy.json + derived/taxonomy.md
//
// Step E (naming/describing entries) is interpretive and happens OUTSIDE this
// command; nothing here depends on any human-authored codebook.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/mikolajsemeniuk/prism/pkg/cluster"
	"github.com/mikolajsemeniuk/prism/pkg/config"
	"github.com/mikolajsemeniuk/prism/pkg/embed"
	"github.com/mikolajsemeniuk/prism/pkg/segment"
)

type entry struct {
	ID          int            `json:"id"`
	Category    string         `json:"category"` // component | product-specific | unstable | micro
	Size        int            `json:"size"`
	Stability   float64        `json:"stability"` // best cross-model containment |A∩B|/|A|
	Products    []string       `json:"products"`
	Classes     map[string]int `json:"classes"`
	Terms       []string       `json:"terms,omitempty"`
	MedoidID    string         `json:"medoid_id,omitempty"`
	MedoidText  string         `json:"medoid_text,omitempty"`
	ExemplarIDs []string       `json:"exemplar_ids,omitempty"`
}

func main() {
	log.SetFlags(0)
	cfgPath := flag.String("config", "pipeline.json", "pre-registered pipeline config")
	segsPath := flag.String("segments", "derived/segments.jsonl", "input segments")
	clustersDir := flag.String("clustersdir", "derived/clusters", "clustering results directory")
	embDir := flag.String("embeddingsdir", "derived/embeddings", "embeddings directory (primary model, for medoids)")
	outDir := flag.String("outdir", "derived", "output directory")
	flag.Parse()

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		log.Fatalf("config: %v", err)
	}
	segs, err := segment.LoadJSONL(*segsPath)
	if err != nil {
		log.Fatalf("segments: %v", err)
	}
	byID := map[string]segment.Segment{}
	for _, s := range segs {
		byID[s.ID] = s
	}
	results, err := cluster.LoadResultsDir(*clustersDir)
	if err != nil {
		log.Fatalf("clusters: %v", err)
	}

	componentK := results[0].ChosenK
	if len(results) >= 2 {
		k, _, err := cluster.Consensus(results)
		if err != nil {
			log.Fatalf("consensus: %v", err)
		}
		componentK = k
	} else {
		log.Printf("WARNING: single model — using its fine cut k=%d", componentK)
	}
	primary := results[0]
	labA, ok := primary.LabelsAt(componentK)
	if !ok {
		log.Fatalf("no stored labels at k=%d for %s", componentK, primary.Model)
	}

	// members per primary cluster
	members := map[int][]string{}
	for id, cl := range labA {
		members[cl] = append(members[cl], id)
	}
	for _, ids := range members {
		sort.Strings(ids)
	}

	// A: stability — best containment vs secondary model's clusters at same k
	stability := map[int]float64{}
	if len(results) >= 2 {
		labB, ok := results[1].LabelsAt(componentK)
		if !ok {
			log.Fatalf("no stored labels at k=%d for %s", componentK, results[1].Model)
		}
		bMembers := map[int]map[string]bool{}
		for id, cl := range labB {
			if bMembers[cl] == nil {
				bMembers[cl] = map[string]bool{}
			}
			bMembers[cl][id] = true
		}
		for cl, ids := range members {
			best := 0.0
			for _, bset := range bMembers {
				inter := 0
				for _, id := range ids {
					if bset[id] {
						inter++
					}
				}
				if c := float64(inter) / float64(len(ids)); c > best {
					best = c
				}
			}
			stability[cl] = best
		}
	} else {
		for cl := range members {
			stability[cl] = 1
		}
	}

	// embeddings of the primary model, for medoids
	embPath := filepath.Join(*embDir, primary.Model+".jsonl")
	recs, err := embed.LoadRecords(embPath)
	if err != nil {
		log.Fatalf("embeddings %s: %v", embPath, err)
	}
	vec := map[string][]float32{}
	for _, r := range recs {
		vec[r.ID] = r.Vec
	}

	// C: distinctive terms — tf-idf with clusters as documents
	df := map[string]int{}
	tf := map[int]map[string]int{}
	for cl, ids := range members {
		tf[cl] = map[string]int{}
		seen := map[string]bool{}
		for _, id := range ids {
			for _, tok := range embed.Tokenize(byID[id].Text) {
				tf[cl][tok]++
				if !seen[tok] {
					seen[tok] = true
					df[tok]++
				}
			}
		}
	}
	nClusters := float64(len(members))

	var entries []entry
	counts := map[string]int{}
	for cl, ids := range members {
		e := entry{ID: cl, Size: len(ids), Stability: stability[cl], Classes: map[string]int{}}
		prodSet := map[string]bool{}
		for _, id := range ids {
			s := byID[id]
			prodSet[s.Product] = true
			e.Classes[s.ProductClass]++
		}
		for p := range prodSet {
			e.Products = append(e.Products, p)
		}
		sort.Strings(e.Products)

		switch {
		case e.Size < cfg.Taxonomy.MinSize:
			e.Category = "micro"
		case e.Stability < cfg.Taxonomy.StabilityMin:
			e.Category = "unstable"
		case len(e.Products) < cfg.Taxonomy.MinProducts:
			e.Category = "product-specific"
		default:
			e.Category = "component"
		}
		counts[e.Category]++

		if e.Category == "component" || e.Category == "product-specific" {
			e.Terms = topTerms(tf[cl], df, nClusters, cfg.Taxonomy.TopTerms)
			e.MedoidID, e.ExemplarIDs = medoidAndExemplars(ids, vec, cfg.Taxonomy.Exemplars)
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
		"model":       primary.Model,
		"component_k": componentK,
		"params":      cfg.Taxonomy,
		"config_hash": cfg.Hash,
		"counts":      counts,
		"entries":     entries,
	}
	raw, _ := json.MarshalIndent(out, "", "  ")
	if err := os.WriteFile(filepath.Join(*outDir, "taxonomy.json"), append(raw, '\n'), 0o644); err != nil {
		log.Fatal(err)
	}
	if err := writeMarkdown(filepath.Join(*outDir, "taxonomy.md"), componentK, primary.Model, cfg, counts, entries, byID); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("taxonomy at k=%d: %d components, %d product-specific, %d unstable, %d micro → %s\n",
		componentK, counts["component"], counts["product-specific"], counts["unstable"], counts["micro"],
		filepath.Join(*outDir, "taxonomy.md"))
}

func topTerms(tf map[string]int, df map[string]int, nClusters float64, n int) []string {
	type ts struct {
		t string
		s float64
	}
	var scored []ts
	for t, f := range tf {
		if f < 2 {
			continue
		}
		scored = append(scored, ts{t, float64(f) * math.Log(nClusters/(1+float64(df[t])))})
	}
	sort.Slice(scored, func(i, j int) bool {
		if scored[i].s != scored[j].s {
			return scored[i].s > scored[j].s
		}
		return scored[i].t < scored[j].t
	})
	var out []string
	for i := 0; i < len(scored) && i < n; i++ {
		out = append(out, scored[i].t)
	}
	return out
}

// medoidAndExemplars returns the member with the highest mean cosine
// similarity to the rest, plus the next most central members.
func medoidAndExemplars(ids []string, vec map[string][]float32, nEx int) (string, []string) {
	type ms struct {
		id string
		s  float64
	}
	var scored []ms
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
				sum += cluster.CosineSim(va, vb)
			}
		}
		scored = append(scored, ms{a, sum})
	}
	if len(scored) == 0 {
		return "", nil
	}
	sort.Slice(scored, func(i, j int) bool {
		if scored[i].s != scored[j].s {
			return scored[i].s > scored[j].s
		}
		return scored[i].id < scored[j].id
	})
	var ex []string
	for i := 0; i < len(scored) && i < nEx; i++ {
		ex = append(ex, scored[i].id)
	}
	return scored[0].id, ex
}

func writeMarkdown(path string, k int, model string, cfg config.Config, counts map[string]int, entries []entry, byID map[string]segment.Segment) error {
	var b strings.Builder
	fmt.Fprintf(&b, "# Data-driven component taxonomy (auto-generated)\n\n")
	fmt.Fprintf(&b, "Generated by cmd/taxonomy — do not edit. Component cut k=%d, primary model %s, config %s.\n\n", k, model, cfg.Hash[:12])
	fmt.Fprintf(&b, "Steps A–D are mechanical (thresholds pre-registered in pipeline.json: containment>=%.2f, products>=%d, size>=%d). Names/definitions (step E) are added in a separate interpretive pass and never modify this file.\n\n",
		cfg.Taxonomy.StabilityMin, cfg.Taxonomy.MinProducts, cfg.Taxonomy.MinSize)
	fmt.Fprintf(&b, "Counts: %d components, %d product-specific, %d unstable, %d micro.\n\n",
		counts["component"], counts["product-specific"], counts["unstable"], counts["micro"])

	fmt.Fprintf(&b, "## Components (cross-product, stable in both embedding models)\n\n")
	for _, e := range entries {
		if e.Category != "component" {
			continue
		}
		fmt.Fprintf(&b, "### component C%d\n\n", e.ID)
		fmt.Fprintf(&b, "- size: %d segments | products: %d (%s) | classes: %s | stability: %.2f\n",
			e.Size, len(e.Products), strings.Join(shorten(e.Products), ", "), classStr(e.Classes), e.Stability)
		fmt.Fprintf(&b, "- distinctive terms: %s\n", strings.Join(e.Terms, ", "))
		fmt.Fprintf(&b, "- medoid [%s]: %q\n", e.MedoidID, e.MedoidText)
		for _, ex := range e.ExemplarIDs {
			if ex == e.MedoidID {
				continue
			}
			fmt.Fprintf(&b, "- exemplar [%s]: %q\n", ex, clip(byID[ex].Text, 160))
		}
		fmt.Fprintf(&b, "\n")
	}

	fmt.Fprintf(&b, "## Product-specific clusters (stable but <%d products — features, not components)\n\n", cfg.Taxonomy.MinProducts)
	for _, e := range entries {
		if e.Category != "product-specific" {
			continue
		}
		fmt.Fprintf(&b, "- C%d: size %d, %s, terms: %s\n", e.ID, e.Size, strings.Join(shorten(e.Products), ", "), strings.Join(e.Terms, ", "))
	}
	return os.WriteFile(path, []byte(b.String()), 0o644)
}

func classStr(m map[string]int) string {
	var ks []string
	for k := range m {
		ks = append(ks, k)
	}
	sort.Strings(ks)
	parts := make([]string, len(ks))
	for i, k := range ks {
		parts[i] = fmt.Sprintf("%s %d", k, m[k])
	}
	return strings.Join(parts, ", ")
}

func shorten(products []string) []string {
	out := make([]string, len(products))
	for i, p := range products {
		if len(p) > 24 {
			p = p[:24] + "…"
		}
		out[i] = p
	}
	return out
}

func clip(s string, n int) string {
	s = strings.Join(strings.Fields(s), " ")
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}

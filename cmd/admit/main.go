// Command admit is step 5 of the decomposition pipeline (see CLAUDE.md
// "Decomposition methodology") — it selects the COMPONENT CUT k*.
//
// Rationale: silhouette-style quality scores monotonically reward
// near-duplicate granularity on this corpus (vendors copy passages from each
// other), so instead the component granularity is defined as the cut at which
// INDEPENDENT embedding models agree most about the partition. For every
// admissible k (from cmd/kbounds) it computes the mean pairwise Adjusted Rand
// Index between the models' labelings and picks k* = argmax (ties → smaller
// k). ARI is agreement corrected for chance: 1.0 = identical partitions,
// ~0 = what two random sorters would achieve.
//
// Inputs: every artefacts/clusters-*.json + artefacts/kbounds.json
// Output: artefacts/admit.json —
//
//	{"k_star","ari","per_k":[{"k","mean_ari","admissible"}]}
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
)

type clusters struct {
	Model string   `json:"model"`
	IDs   []string `json:"ids"`
	Grid  []struct {
		K      int   `json:"k"`
		Labels []int `json:"labels"`
	} `json:"grid"`
}

type kbounds struct {
	Lo   int `json:"lo"`
	Hi   int `json:"hi"`
	PerK []struct {
		K       int  `json:"k"`
		PassAll bool `json:"pass_all"`
	} `json:"per_k"`
}

type perK struct {
	K          int     `json:"k"`
	MeanARI    float64 `json:"mean_ari"`
	Admissible bool    `json:"admissible"`
}

func main() {
	log.SetFlags(0)
	inDir := flag.String("indir", "artefacts", "directory with clusters-*.json and kbounds.json")
	out := flag.String("out", "artefacts/admit.json", "output path")
	flag.Parse()

	paths, _ := filepath.Glob(filepath.Join(*inDir, "clusters-*.json"))
	sort.Strings(paths)
	if len(paths) < 2 {
		log.Fatalf("need >=2 clusters-*.json, found %d", len(paths))
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

	var kb kbounds
	raw, err := os.ReadFile(filepath.Join(*inDir, "kbounds.json"))
	if err != nil {
		log.Fatalf("%v (run cmd/kbounds first)", err)
	}
	if err := json.Unmarshal(raw, &kb); err != nil {
		log.Fatal(err)
	}
	admissible := map[int]bool{}
	for _, r := range kb.PerK {
		admissible[r.K] = r.PassAll
	}

	// All models must share the segment id order for label vectors to align;
	// cmd/cluster sorts ids, so verify and map by intersection if needed.
	base := all[0].IDs
	for _, c := range all[1:] {
		if len(c.IDs) != len(base) {
			log.Fatalf("id sets differ between models (%d vs %d) — re-run cmd/embed+cluster from the same segments file", len(c.IDs), len(base))
		}
		for i := range base {
			if c.IDs[i] != base[i] {
				log.Fatalf("id order differs between models at %d — re-run cmd/cluster", i)
			}
		}
	}

	labelsAt := func(c clusters, k int) []int {
		for _, g := range c.Grid {
			if g.K == k {
				return g.Labels
			}
		}
		return nil
	}

	var ks []int
	for _, g := range all[0].Grid {
		ks = append(ks, g.K)
	}
	sort.Ints(ks)

	var rows []perK
	bestK, bestARI := 0, -2.0
	for _, k := range ks {
		var sum float64
		var pairs int
		ok := true
		for i := 0; i < len(all); i++ {
			for j := i + 1; j < len(all); j++ {
				la, lb := labelsAt(all[i], k), labelsAt(all[j], k)
				if la == nil || lb == nil {
					ok = false
					continue
				}
				sum += ari(la, lb)
				pairs++
			}
		}
		if !ok || pairs == 0 {
			continue
		}
		row := perK{K: k, MeanARI: sum / float64(pairs), Admissible: admissible[k]}
		rows = append(rows, row)
		if row.Admissible && row.MeanARI > bestARI {
			bestARI, bestK = row.MeanARI, k
		}
	}
	if bestK == 0 {
		log.Fatal("no admissible cut to choose from — check kbounds.json")
	}

	for _, r := range rows {
		mark := ""
		if !r.Admissible {
			mark = "  (outside window)"
		}
		if r.K == bestK {
			mark = "  <== component cut k*"
		}
		fmt.Printf("k=%-5d meanARI=%.4f%s\n", r.K, r.MeanARI, mark)
	}
	fmt.Printf("\ncomponent cut: k*=%d (mean cross-model ARI %.3f)\n", bestK, bestARI)

	outRaw, _ := json.MarshalIndent(map[string]any{"k_star": bestK, "ari": bestARI, "per_k": rows}, "", "  ")
	if err := os.WriteFile(*out, append(outRaw, '\n'), 0o644); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("wrote %s\n", *out)
}

// ari computes the Adjusted Rand Index between two labelings of the same
// points (agreement corrected for chance).
func ari(x, y []int) float64 {
	n := len(x)
	cont := map[[2]int]int{}
	rows, cols := map[int]int{}, map[int]int{}
	for i := range x {
		cont[[2]int{x[i], y[i]}]++
		rows[x[i]]++
		cols[y[i]]++
	}
	comb2 := func(v int) float64 { return float64(v) * float64(v-1) / 2 }
	var cells, sa, sb float64
	for _, v := range cont {
		cells += comb2(v)
	}
	for _, v := range rows {
		sa += comb2(v)
	}
	for _, v := range cols {
		sb += comb2(v)
	}
	expected := sa * sb / comb2(n)
	max := (sa + sb) / 2
	if max == expected {
		return 1
	}
	return (cells - expected) / (max - expected)
}

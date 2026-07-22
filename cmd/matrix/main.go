// Command matrix consumes all clustering results and emits the descriptive
// outputs: cross-model stability (ARI/NMI), the pre-registered control
// assertions (exit code 1 if any fails — this is the pipeline's smoke gate),
// and the cluster×product matrix.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/mikolajsemeniuk/prism/pkg/cluster"
	"github.com/mikolajsemeniuk/prism/pkg/config"
	"github.com/mikolajsemeniuk/prism/pkg/controls"
	"github.com/mikolajsemeniuk/prism/pkg/segment"
)

type stabilityPair struct {
	A       string  `json:"a"`
	B       string  `json:"b"`
	NShared int     `json:"n_shared"`
	ARI     float64 `json:"ari"`
	NMI     float64 `json:"nmi"`
}

type controlModelResult struct {
	Model string `json:"model"`
	controls.Result
}

type controlResult struct {
	Name     string               `json:"name"`
	Expect   string               `json:"expect"`
	PerModel []controlModelResult `json:"per_model"`
	Pass     bool                 `json:"pass"`
}

type clusterRow struct {
	Cluster  int            `json:"cluster"`
	Size     int            `json:"size"`
	Products map[string]int `json:"products"`
	Classes  map[string]int `json:"classes"`
	Examples []string       `json:"examples"`
}

func main() {
	log.SetFlags(0)
	cfgPath := flag.String("config", "pipeline.json", "pre-registered pipeline config")
	segsPath := flag.String("segments", "derived/segments.jsonl", "input segments")
	clustersDir := flag.String("clustersdir", "derived/clusters", "directory of clustering results")
	outDir := flag.String("outdir", "derived", "output directory")
	strict := flag.Bool("strict", true, "exit 1 when a control assertion fails")
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

	// --- stability across models (fine cuts) ---
	var pairs []stabilityPair
	for i := 0; i < len(results); i++ {
		for j := i + 1; j < len(results); j++ {
			ids := sharedIDs(results[i].Assignments, results[j].Assignments)
			x := make([]int, len(ids))
			y := make([]int, len(ids))
			for t, id := range ids {
				x[t] = results[i].Assignments[id]
				y[t] = results[j].Assignments[id]
			}
			pairs = append(pairs, stabilityPair{
				A: results[i].Model, B: results[j].Model, NShared: len(ids),
				ARI: cluster.ARI(x, y), NMI: cluster.NMI(x, y),
			})
		}
	}

	// --- component cut: cross-model consensus (two-level readout) ---
	componentK := results[0].ChosenK
	var consensus map[string]any
	if len(results) >= 2 {
		k, points, err := cluster.Consensus(results)
		if err != nil {
			log.Fatalf("consensus: %v", err)
		}
		componentK = k
		consensus = map[string]any{"k": k, "per_k": points}
	} else {
		log.Printf("WARNING: single model — component cut falls back to its fine cut k=%d", componentK)
	}
	writeJSON(filepath.Join(*outDir, "stability.json"), map[string]any{"pairs": pairs, "consensus": consensus})

	// --- pre-registered controls ---
	allPass := true
	var ctlResults []controlResult
	for _, c := range cfg.Controls {
		cr := controlResult{Name: c.Name, Expect: c.Expect, Pass: true}
		for _, res := range results {
			mr := controlModelResult{Model: res.Model, Result: controls.Evaluate(c, res.Assignments, byID)}
			cr.PerModel = append(cr.PerModel, mr)
			if !mr.Pass {
				cr.Pass = false
			}
		}
		if !cr.Pass {
			allPass = false
		}
		ctlResults = append(ctlResults, cr)
	}
	writeJSON(filepath.Join(*outDir, "controls.json"), map[string]any{"pass": allPass, "controls": ctlResults})

	// --- cluster × product matrix at the COMPONENT cut (primary = first
	// model alphabetically) ---
	primary := results[0]
	componentAssignments, ok := primary.LabelsAt(componentK)
	if !ok {
		log.Fatalf("model %s has no stored labels at component k=%d", primary.Model, componentK)
	}
	rowsByCluster := map[int]*clusterRow{}
	for id, cl := range componentAssignments {
		row := rowsByCluster[cl]
		if row == nil {
			row = &clusterRow{Cluster: cl, Products: map[string]int{}, Classes: map[string]int{}}
			rowsByCluster[cl] = row
		}
		row.Size++
		if s, ok := byID[id]; ok {
			row.Products[s.Product]++
			row.Classes[s.ProductClass]++
		}
		if len(row.Examples) < 3 {
			row.Examples = append(row.Examples, id)
		}
	}
	rows := make([]clusterRow, 0, len(rowsByCluster))
	for _, r := range rowsByCluster {
		sort.Strings(r.Examples)
		rows = append(rows, *r)
	}
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].Size != rows[j].Size {
			return rows[i].Size > rows[j].Size
		}
		return rows[i].Cluster < rows[j].Cluster
	})
	writeJSON(filepath.Join(*outDir, "matrix.json"), map[string]any{
		"model": primary.Model, "fine_k": primary.ChosenK, "component_k": componentK, "clusters": rows,
	})

	for _, p := range pairs {
		fmt.Printf("stability (fine) %s vs %s: ARI=%.3f NMI=%.3f (n=%d)\n", p.A, p.B, p.ARI, p.NMI, p.NShared)
	}
	if consensus != nil {
		fmt.Printf("component cut: k=%d (max cross-model ARI subject to controls)\n", componentK)
	}
	for _, c := range ctlResults {
		status := "PASS"
		if !c.Pass {
			status = "FAIL"
		}
		fmt.Printf("control %-28s expect=%-9s %s\n", c.Name, c.Expect, status)
	}
	if !allPass {
		if *strict {
			log.Fatal("control assertions FAILED — the pipeline must not be trusted in this state")
		}
		fmt.Println("WARNING: control assertions failed (non-strict mode)")
	}
}

func sharedIDs(a, b map[string]int) []string {
	var ids []string
	for id := range a {
		if _, ok := b[id]; ok {
			ids = append(ids, id)
		}
	}
	sort.Strings(ids)
	return ids
}

func writeJSON(path string, v any) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		log.Fatal(err)
	}
	raw, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile(path, append(raw, '\n'), 0o644); err != nil {
		log.Fatal(err)
	}
}

// Command cluster runs deterministic average-linkage agglomerative
// clustering on one model's embeddings and evaluates the pre-registered k
// grid. Selection rule (pre-registered): among grid cuts where ALL control
// assertions hold, pick max mean silhouette (ties → smaller k). If no cut
// satisfies the controls the command fails — a pipeline that cannot separate
// known-different components while keeping known-identical passages together
// must not produce results.
package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"sort"
	"strings"

	"github.com/mikolajsemeniuk/prism/pkg/cluster"
	"github.com/mikolajsemeniuk/prism/pkg/config"
	"github.com/mikolajsemeniuk/prism/pkg/controls"
	"github.com/mikolajsemeniuk/prism/pkg/embed"
	"github.com/mikolajsemeniuk/prism/pkg/segment"
)

func main() {
	log.SetFlags(0)
	cfgPath := flag.String("config", "pipeline.json", "pre-registered pipeline config")
	segsPath := flag.String("segments", "derived/segments.jsonl", "segments (for control-assertion text matching)")
	embPath := flag.String("embeddings", "", "input embeddings JSONL (required)")
	out := flag.String("out", "", "output JSON path (required)")
	bootstrapOverride := flag.Int("bootstrap", -1, "override bootstrap rounds (-1 = use config)")
	flag.Parse()
	if *embPath == "" || *out == "" {
		log.Fatal("-embeddings and -out are required")
	}

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
	recs, err := embed.LoadRecords(*embPath)
	if err != nil {
		log.Fatalf("embeddings: %v", err)
	}
	sort.Slice(recs, func(i, j int) bool { return recs[i].ID < recs[j].ID })
	n := len(recs)
	if n < 3 {
		log.Fatalf("need at least 3 vectors, got %d", n)
	}
	vecs := make([][]float32, n)
	ids := make([]string, n)
	for i, r := range recs {
		vecs[i], ids[i] = r.Vec, r.ID
	}

	m := cluster.CosineDistances(vecs)
	merges := cluster.Agglomerative(m)

	toAssignments := func(labels []int) map[string]int {
		a := make(map[string]int, n)
		for i, id := range ids {
			a[id] = labels[i]
		}
		return a
	}

	var grid []cluster.GridPoint
	var gridLabels []cluster.GridLabels
	bestK, bestSil := 0, -2.0
	labelsByK := map[int][]int{}
	for _, k := range cfg.Cluster.KGrid {
		if k >= n {
			continue
		}
		labels := cluster.Cut(merges, n, k)
		sil := cluster.Silhouette(m, labels)
		pass := controls.AllPass(cfg.Controls, toAssignments(labels), byID)
		grid = append(grid, cluster.GridPoint{K: k, Silhouette: sil, ControlsPass: pass})
		gridLabels = append(gridLabels, cluster.GridLabels{K: k, Labels: labels})
		labelsByK[k] = labels
		if pass && (sil > bestSil || (sil == bestSil && k < bestK)) {
			bestSil, bestK = sil, k
		}
	}
	if bestK == 0 {
		for _, g := range grid {
			log.Printf("  k=%-3d silhouette=%.4f controls_pass=%v", g.K, g.Silhouette, g.ControlsPass)
		}
		log.Fatal("no cut in the pre-registered k grid satisfies all control assertions — pipeline must not produce results in this state")
	}
	chosen := labelsByK[bestK]

	rounds := cfg.Cluster.Bootstrap.Rounds
	if *bootstrapOverride >= 0 {
		rounds = *bootstrapOverride
	}
	var bs *cluster.BootstrapSummary
	if rounds > 0 {
		aris := cluster.BootstrapARI(m, chosen, bestK, rounds, cfg.Cluster.Bootstrap.Seed)
		bs = &cluster.BootstrapSummary{
			Rounds:  rounds,
			Seed:    cfg.Cluster.Bootstrap.Seed,
			MeanARI: cluster.Mean(aris),
			P2_5:    cluster.Percentile(aris, 2.5),
			P97_5:   cluster.Percentile(aris, 97.5),
		}
	}

	res := &cluster.Result{
		Model:       modelName(*embPath),
		Embeddings:  *embPath,
		ConfigHash:  cfg.Hash,
		N:           n,
		Grid:        grid,
		ChosenK:     bestK,
		Bootstrap:   bs,
		Assignments: toAssignments(chosen),
		IDs:         ids,
		GridLabels:  gridLabels,
	}
	if err := res.Save(*out); err != nil {
		log.Fatalf("write %s: %v", *out, err)
	}
	fmt.Printf("wrote %s: n=%d chosen_k=%d silhouette=%.4f (controls-constrained)", *out, n, bestK, bestSil)
	if bs != nil {
		fmt.Printf(" bootstrap_ari=%.3f [%.3f, %.3f] (%d rounds)", bs.MeanARI, bs.P2_5, bs.P97_5, bs.Rounds)
	}
	fmt.Println()
}

func modelName(embPath string) string {
	return strings.TrimSuffix(filepath.Base(embPath), ".jsonl")
}

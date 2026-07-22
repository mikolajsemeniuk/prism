// Command cluster is step 3 of the decomposition pipeline (see CLAUDE.md
// "Decomposition methodology").
//
// Pure math, no network, no policy: it loads one model's embeddings, builds
// the pairwise cosine-distance matrix, runs average-linkage agglomerative
// clustering with the nearest-neighbor-chain algorithm (fully deterministic —
// no random seeds; ties break toward the smallest index), and stores the flat
// labeling for EVERY cut of the pre-registered k grid. Downstream steps
// (cmd/kbounds, cmd/consensus) decide which cuts are admissible and which one
// to use; this command deliberately knows nothing about controls or needles.
//
// Run once per model:
//
//	go run ./cmd/cluster -model bge-m3
//	go run ./cmd/cluster -model nomic-embed-text
//
// Input:  artefacts/embeddings-<model>.jsonl
// Output: artefacts/clusters-<model>.json —
//
//	{"model","n","ids":[...],"grid":[{"k":10,"labels":[...]}, ...]}
//
// labels[i] is the cluster of ids[i]; labels are renumbered 0..k-1 in order
// of first appearance.
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

// kGrid is the pre-registered grid of dendrogram cuts (CLAUDE.md).
var kGrid = []int{10, 15, 20, 25, 30, 35, 40, 50, 60, 70, 80, 100, 120, 150, 200, 250, 300, 400, 500}

type record struct {
	ID  string    `json:"id"`
	Vec []float32 `json:"vec"`
}

type gridCut struct {
	K      int   `json:"k"`
	Labels []int `json:"labels"`
}

func main() {
	log.SetFlags(0)
	model := flag.String("model", "", "model name whose embeddings to cluster (required)")
	inDir := flag.String("indir", "artefacts", "directory with embeddings-<model>.jsonl")
	outDir := flag.String("outdir", "artefacts", "output directory")
	flag.Parse()
	if *model == "" {
		log.Fatal("-model is required")
	}

	inPath := filepath.Join(*inDir, "embeddings-"+sanitize(*model)+".jsonl")
	f, err := os.Open(inPath)
	if err != nil {
		log.Fatalf("%v (run cmd/embed -model %s first)", err, *model)
	}
	var recs []record
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 0, 1<<20), 64<<20)
	for sc.Scan() {
		var r record
		if err := json.Unmarshal(sc.Bytes(), &r); err != nil {
			log.Fatalf("%s: %v", inPath, err)
		}
		recs = append(recs, r)
	}
	f.Close()
	if err := sc.Err(); err != nil {
		log.Fatal(err)
	}
	sort.Slice(recs, func(i, j int) bool { return recs[i].ID < recs[j].ID })
	n := len(recs)
	if n < 3 {
		log.Fatalf("need >=3 vectors, got %d", n)
	}

	// cosine distance; vectors are L2-normalized by cmd/embed, so 1 − dot.
	dist := make([]float64, n*n)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			var dot float64
			for d := range recs[i].Vec {
				dot += float64(recs[i].Vec[d]) * float64(recs[j].Vec[d])
			}
			v := 1 - dot
			if v < 0 {
				v = 0
			}
			dist[i*n+j], dist[j*n+i] = v, v
		}
	}
	merges := agglomerative(dist, n)

	ids := make([]string, n)
	for i, r := range recs {
		ids[i] = r.ID
	}
	var grid []gridCut
	for _, k := range kGrid {
		if k >= n {
			continue
		}
		grid = append(grid, gridCut{K: k, Labels: cut(merges, n, k)})
	}

	out := map[string]any{"model": *model, "n": n, "ids": ids, "grid": grid}
	raw, _ := json.Marshal(out)
	outPath := filepath.Join(*outDir, "clusters-"+sanitize(*model)+".json")
	if err := os.MkdirAll(*outDir, 0o755); err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile(outPath, append(raw, '\n'), 0o644); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("wrote %s: n=%d, %d grid cuts (k=%d..%d)\n", outPath, n, len(grid), grid[0].K, grid[len(grid)-1].K)
}

type merge struct {
	a, b, id int
	height   float64
}

// agglomerative implements average-linkage clustering via the
// nearest-neighbor-chain algorithm with Lance-Williams distance updates.
func agglomerative(dist []float64, n int) []merge {
	d := append([]float64(nil), dist...)
	at := func(i, j int) float64 { return d[i*n+j] }
	set := func(i, j int, v float64) { d[i*n+j], d[j*n+i] = v, v }
	size := make([]int, n)
	active := make([]bool, n)
	id := make([]int, n)
	for i := range size {
		size[i], active[i], id[i] = 1, true, i
	}
	next, remaining := n, n
	merges := make([]merge, 0, n-1)
	var chain []int
	for remaining > 1 {
		if len(chain) == 0 {
			for i := 0; i < n; i++ {
				if active[i] {
					chain = append(chain, i)
					break
				}
			}
		}
		for {
			a := chain[len(chain)-1]
			b, best := -1, math.Inf(1)
			for j := 0; j < n; j++ {
				if active[j] && j != a && at(a, j) < best {
					best, b = at(a, j), j
				}
			}
			if len(chain) >= 2 && b == chain[len(chain)-2] {
				chain = chain[:len(chain)-2]
				i, j := a, b
				if j < i {
					i, j = j, i
				}
				ni, nj := size[i], size[j]
				for k := 0; k < n; k++ {
					if active[k] && k != i && k != j {
						set(i, k, (float64(ni)*at(i, k)+float64(nj)*at(j, k))/float64(ni+nj))
					}
				}
				merges = append(merges, merge{id[i], id[j], next, best})
				active[j] = false
				size[i] = ni + nj
				id[i] = next
				next++
				remaining--
				break
			}
			chain = append(chain, b)
		}
	}
	return merges
}

// cut applies the n−k lowest merges (stable by height, then recorded order —
// valid because average linkage is reducible, so children precede parents)
// and renumbers cluster labels 0..k-1 in order of first appearance.
func cut(merges []merge, n, k int) []int {
	type ord struct {
		m   merge
		pos int
	}
	sorted := make([]ord, len(merges))
	for i, m := range merges {
		sorted[i] = ord{m, i}
	}
	sort.SliceStable(sorted, func(i, j int) bool {
		if sorted[i].m.height != sorted[j].m.height {
			return sorted[i].m.height < sorted[j].m.height
		}
		return sorted[i].pos < sorted[j].pos
	})
	parent := make([]int, n)
	for i := range parent {
		parent[i] = i
	}
	var find func(int) int
	find = func(x int) int {
		for parent[x] != x {
			parent[x] = parent[parent[x]]
			x = parent[x]
		}
		return x
	}
	rep := map[int]int{}
	for i := 0; i < n; i++ {
		rep[i] = i
	}
	for t := 0; t < n-k && t < len(sorted); t++ {
		m := sorted[t].m
		ra, rb := find(rep[m.a]), find(rep[m.b])
		parent[rb] = ra
		rep[m.id] = ra
	}
	labels := make([]int, n)
	seen := map[int]int{}
	for i := 0; i < n; i++ {
		r := find(i)
		if _, ok := seen[r]; !ok {
			seen[r] = len(seen)
		}
		labels[i] = seen[r]
	}
	return labels
}

func sanitize(s string) string {
	return strings.Map(func(r rune) rune {
		if r == ':' || r == '/' || r == ' ' {
			return '-'
		}
		return r
	}, s)
}

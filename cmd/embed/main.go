// Command embed turns behavioral segments into vectors for one embedding
// model. Vectors are cached by segment-text hash, so unchanged segments never
// re-embed and reruns are cheap and byte-stable.
package main

import (
	"context"
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

	"github.com/mikolajsemeniuk/prism/pkg/config"
	"github.com/mikolajsemeniuk/prism/pkg/embed"
	"github.com/mikolajsemeniuk/prism/pkg/segment"
)

func main() {
	log.SetFlags(0)
	cfgPath := flag.String("config", "pipeline.json", "pre-registered pipeline config")
	segsPath := flag.String("segments", "derived/segments.jsonl", "input segments")
	backend := flag.String("backend", "ollama", "embedding backend: ollama | fake")
	model := flag.String("model", "", "embedding model name (required)")
	outDir := flag.String("outdir", "derived/embeddings", "output directory")
	flag.Parse()
	if *model == "" {
		log.Fatal("-model is required")
	}

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		log.Fatalf("config: %v", err)
	}
	segs, err := segment.LoadJSONL(*segsPath)
	if err != nil {
		log.Fatalf("segments: %v", err)
	}

	include := map[string]bool{}
	for _, l := range cfg.Segment.IncludeLayers {
		include[l] = true
	}
	var picked []segment.Segment
	for _, s := range segs {
		if include[string(s.Layer)] {
			picked = append(picked, s)
		}
	}
	sort.Slice(picked, func(i, j int) bool { return picked[i].ID < picked[j].ID })
	if len(picked) == 0 {
		log.Fatalf("no segments matched include_layers=%v", cfg.Segment.IncludeLayers)
	}

	outPath := filepath.Join(*outDir, sanitize(*model)+".jsonl")
	cache := map[string][]float32{}
	if prev, err := embed.LoadRecords(outPath); err == nil {
		for _, r := range prev {
			cache[r.Hash] = r.Vec
		}
	}

	em, err := embed.New(*backend, *model, cfg.Embed.Endpoint, cfg.Embed.FakeDims)
	if err != nil {
		log.Fatal(err)
	}

	recs := make([]embed.Record, len(picked))
	var missIdx []int
	for i, s := range picked {
		h := sha256.Sum256([]byte(s.Text))
		recs[i] = embed.Record{ID: s.ID, Hash: hex.EncodeToString(h[:])}
		if v, ok := cache[recs[i].Hash]; ok {
			recs[i].Vec = v
		} else {
			missIdx = append(missIdx, i)
		}
	}
	ctx := context.Background()
	for start := 0; start < len(missIdx); start += cfg.Embed.BatchSize {
		end := start + cfg.Embed.BatchSize
		if end > len(missIdx) {
			end = len(missIdx)
		}
		batch := missIdx[start:end]
		texts := make([]string, len(batch))
		for i, bi := range batch {
			texts[i] = picked[bi].Text
		}
		vecs, err := em.Embed(ctx, texts)
		if err != nil {
			log.Fatalf("embed batch %d-%d via %s: %v", start, end, em.Name(), err)
		}
		for i, bi := range batch {
			recs[bi].Vec = vecs[i]
		}
	}

	if err := embed.SaveRecords(outPath, recs); err != nil {
		log.Fatalf("write %s: %v", outPath, err)
	}
	manifest := map[string]any{
		"backend":     *backend,
		"model":       *model,
		"embedder":    em.Name(),
		"config_hash": cfg.Hash,
		"segments":    *segsPath,
		"count":       len(recs),
		"embedded":    len(missIdx),
		"cached":      len(recs) - len(missIdx),
		"dims":        len(recs[0].Vec),
	}
	raw, _ := json.MarshalIndent(manifest, "", "  ")
	manifestPath := filepath.Join(*outDir, sanitize(*model)+".manifest.json")
	if err := os.WriteFile(manifestPath, append(raw, '\n'), 0o644); err != nil {
		log.Fatalf("write %s: %v", manifestPath, err)
	}
	fmt.Printf("wrote %s: %d vectors (%d embedded, %d cached, dims=%d)\n",
		outPath, len(recs), len(missIdx), len(recs)-len(missIdx), len(recs[0].Vec))
}

func sanitize(s string) string {
	return strings.Map(func(r rune) rune {
		switch r {
		case ':', '/', ' ':
			return '-'
		}
		return r
	}, s)
}

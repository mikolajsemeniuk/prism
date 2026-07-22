// Command embed is step 2 of the decomposition pipeline (see CLAUDE.md
// "Decomposition methodology").
//
// It reads the segments produced by cmd/segment, keeps only the BEHAVIORAL
// layer (tool schemas, few-shot examples, and template variables would
// dominate the embedding space with boilerplate — findings F1/F2), and embeds
// each segment with ONE Ollama embedding model per invocation. Vectors are
// L2-normalized so downstream cosine similarity is a plain dot product.
//
// The output file doubles as the cache: on re-runs, segments whose sha256
// text hash already appears in the existing output are not re-embedded, so
// only the first run per model touches the network. Run once per model, e.g.:
//
//	go run ./cmd/embed -model bge-m3
//	go run ./cmd/embed -model nomic-embed-text
//
// Input:  artefacts/segments.jsonl
// Output: artefacts/embeddings-<model>.jsonl — {"id","hash","vec"} per line,
//
//	sorted by id.
package main

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const batchSize = 16

type seg struct {
	ID    string `json:"id"`
	Layer string `json:"layer"`
	Text  string `json:"text"`
}

type record struct {
	ID   string    `json:"id"`
	Hash string    `json:"hash"`
	Vec  []float32 `json:"vec"`
}

func main() {
	log.SetFlags(0)
	segsPath := flag.String("segments", "artefacts/segments.jsonl", "input segments")
	model := flag.String("model", "", "Ollama embedding model (required)")
	endpoint := flag.String("endpoint", "http://localhost:11434", "Ollama endpoint")
	outDir := flag.String("outdir", "artefacts", "output directory")
	flag.Parse()
	if *model == "" {
		log.Fatal("-model is required")
	}

	var behavioral []seg
	f, err := os.Open(*segsPath)
	if err != nil {
		log.Fatalf("segments: %v (run cmd/segment first)", err)
	}
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 0, 1<<20), 64<<20)
	for sc.Scan() {
		var s seg
		if err := json.Unmarshal(sc.Bytes(), &s); err != nil {
			log.Fatalf("%s: %v", *segsPath, err)
		}
		if s.Layer == "behavioral" {
			behavioral = append(behavioral, s)
		}
	}
	f.Close()
	if err := sc.Err(); err != nil {
		log.Fatal(err)
	}
	sort.Slice(behavioral, func(i, j int) bool { return behavioral[i].ID < behavioral[j].ID })

	outPath := filepath.Join(*outDir, "embeddings-"+sanitize(*model)+".jsonl")
	cache := loadCache(outPath)

	recs := make([]record, len(behavioral))
	var missing []int
	for i, s := range behavioral {
		h := sha256.Sum256([]byte(s.Text))
		recs[i] = record{ID: s.ID, Hash: hex.EncodeToString(h[:])}
		if v, ok := cache[recs[i].Hash]; ok {
			recs[i].Vec = v
		} else {
			missing = append(missing, i)
		}
	}

	client := &http.Client{Timeout: 5 * time.Minute}
	for s := 0; s < len(missing); s += batchSize {
		e := min(s+batchSize, len(missing))
		texts := make([]string, 0, e-s)
		for _, i := range missing[s:e] {
			texts = append(texts, behavioral[i].Text)
		}
		vecs, err := ollamaEmbed(client, *endpoint, *model, texts)
		if err != nil {
			log.Fatal(err)
		}
		for j, i := range missing[s:e] {
			recs[i].Vec = vecs[j]
		}
		fmt.Printf("  %s: embedded %d/%d\r", *model, e, len(missing))
	}
	if len(missing) > 0 {
		fmt.Println()
	}

	if err := os.MkdirAll(*outDir, 0o755); err != nil {
		log.Fatal(err)
	}
	out, err := os.Create(outPath)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	w := bufio.NewWriter(out)
	enc := json.NewEncoder(w)
	for _, r := range recs {
		if err := enc.Encode(r); err != nil {
			log.Fatal(err)
		}
	}
	if err := w.Flush(); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("wrote %s: %d vectors (%d embedded, %d cached, dims=%d)\n",
		outPath, len(recs), len(missing), len(recs)-len(missing), len(recs[0].Vec))
}

// loadCache reads a previous output file (if any) and returns hash → vector.
func loadCache(path string) map[string][]float32 {
	cache := map[string][]float32{}
	f, err := os.Open(path)
	if err != nil {
		return cache
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 0, 1<<20), 64<<20)
	for sc.Scan() {
		var r record
		if json.Unmarshal(sc.Bytes(), &r) == nil && r.Hash != "" && len(r.Vec) > 0 {
			cache[r.Hash] = r.Vec
		}
	}
	return cache
}

func ollamaEmbed(client *http.Client, endpoint, model string, texts []string) ([][]float32, error) {
	body, _ := json.Marshal(map[string]any{"model": model, "input": texts})
	resp, err := client.Post(strings.TrimRight(endpoint, "/")+"/api/embed", "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("ollama (%s): %w", model, err)
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ollama (%s): %s: %.160s", model, resp.Status, string(raw))
	}
	var parsed struct {
		Embeddings [][]float64 `json:"embeddings"`
	}
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return nil, err
	}
	if len(parsed.Embeddings) != len(texts) {
		return nil, fmt.Errorf("ollama (%s): %d embeddings for %d inputs", model, len(parsed.Embeddings), len(texts))
	}
	out := make([][]float32, len(parsed.Embeddings))
	for i, e := range parsed.Embeddings {
		v := make([]float32, len(e))
		var norm float64
		for d, x := range e {
			v[d] = float32(x)
			norm += x * x
		}
		if norm > 0 {
			inv := float32(1 / math.Sqrt(norm))
			for d := range v {
				v[d] *= inv
			}
		}
		out[i] = v
	}
	return out, nil
}

func sanitize(s string) string {
	return strings.Map(func(r rune) rune {
		if r == ':' || r == '/' || r == ' ' {
			return '-'
		}
		return r
	}, s)
}

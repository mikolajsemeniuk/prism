// Package embed provides embedding backends behind one interface: a
// deterministic offline fake (for smoke runs and tests) and an Ollama HTTP
// client (for the real open-weights models). Vectors are L2-normalized so
// cosine similarity is a plain dot product downstream.
package embed

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"io"
	"math"
	"net/http"
	"strings"
	"time"
)

type Embedder interface {
	Name() string
	Embed(ctx context.Context, texts []string) ([][]float32, error)
}

// New returns an embedder. backend is "fake" or "ollama"; dims applies to the
// fake backend only.
func New(backend, model, endpoint string, dims int) (Embedder, error) {
	switch backend {
	case "fake":
		return &Fake{Model: model, Dims: dims}, nil
	case "ollama":
		return &Ollama{Model: model, Endpoint: strings.TrimRight(endpoint, "/")}, nil
	default:
		return nil, fmt.Errorf("unknown embed backend %q", backend)
	}
}

// Fake is a deterministic bag-of-ngrams hashing embedder. Identical texts map
// to identical vectors and near-identical texts to near-identical vectors, so
// the positive-control assertions are meaningful without any model server.
type Fake struct {
	Model string
	Dims  int
}

func (f *Fake) Name() string { return "fake:" + f.Model }

// stopwords keeps ubiquitous function words out of the fake vectors so
// unrelated texts land near-orthogonal (roughly what a real embedding model
// achieves semantically).
var stopwords = map[string]bool{
	"the": true, "a": true, "an": true, "and": true, "or": true, "of": true,
	"to": true, "in": true, "on": true, "for": true, "with": true, "is": true,
	"are": true, "be": true, "it": true, "this": true, "that": true, "as": true,
	"at": true, "by": true, "if": true, "you": true, "your": true, "not": true,
	"do": true, "does": true, "can": true, "will": true, "must": true,
	"should": true, "when": true, "any": true, "all": true, "no": true,
	"never": true, "always": true, "use": true, "only": true, "may": true,
}

// Tokenize lowercases, splits on non-alphanumeric runes, and drops stopwords
// and tokens shorter than 3 runes. Shared by the taxonomy term extraction.
func Tokenize(t string) []string {
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

func (f *Fake) Embed(_ context.Context, texts []string) ([][]float32, error) {
	seed := fnvHash(f.Model)
	out := make([][]float32, len(texts))
	for i, t := range texts {
		counts := map[string]float64{}
		words := strings.FieldsFunc(strings.ToLower(t), func(r rune) bool {
			return !('a' <= r && r <= 'z' || '0' <= r && r <= '9' || r == '-' || r == '_')
		})
		prev := ""
		for _, w := range words {
			if !stopwords[w] {
				counts[w]++
			}
			if prev != "" && !(stopwords[prev] && stopwords[w]) {
				counts[prev+" "+w]++
			}
			prev = w
		}
		v := make([]float32, f.Dims)
		for tok, c := range counts {
			h := fnvHash(tok) ^ seed
			idx := int(h % uint64(f.Dims))
			w := float32(math.Sqrt(c)) // sublinear tf
			if h&(1<<63) != 0 {
				v[idx] -= w
			} else {
				v[idx] += w
			}
		}
		normalize(v)
		out[i] = v
	}
	return out, nil
}

// Ollama talks to a local Ollama server's /api/embed endpoint.
type Ollama struct {
	Model    string
	Endpoint string
	Client   *http.Client
}

func (o *Ollama) Name() string { return "ollama:" + o.Model }

func (o *Ollama) Embed(ctx context.Context, texts []string) ([][]float32, error) {
	client := o.Client
	if client == nil {
		client = &http.Client{Timeout: 5 * time.Minute}
	}
	reqBody, err := json.Marshal(map[string]any{"model": o.Model, "input": texts})
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, o.Endpoint+"/api/embed", bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ollama %s: %s", resp.Status, truncate(string(raw), 200))
	}
	var parsed struct {
		Embeddings [][]float64 `json:"embeddings"`
	}
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return nil, err
	}
	if len(parsed.Embeddings) != len(texts) {
		return nil, fmt.Errorf("ollama returned %d embeddings for %d inputs", len(parsed.Embeddings), len(texts))
	}
	out := make([][]float32, len(parsed.Embeddings))
	for i, e := range parsed.Embeddings {
		v := make([]float32, len(e))
		for j, x := range e {
			v[j] = float32(x)
		}
		normalize(v)
		out[i] = v
	}
	return out, nil
}

func normalize(v []float32) {
	var sum float64
	for _, x := range v {
		sum += float64(x) * float64(x)
	}
	if sum == 0 {
		return
	}
	inv := 1 / math.Sqrt(sum)
	for i := range v {
		v[i] = float32(float64(v[i]) * inv)
	}
}

func fnvHash(s string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(s))
	return h.Sum64()
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}

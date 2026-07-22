package cluster

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// GridPoint is one evaluated dendrogram cut. ControlsPass records whether
// the pre-registered control assertions hold at this cut — the selection rule
// only considers cuts where they do.
type GridPoint struct {
	K            int     `json:"k"`
	Silhouette   float64 `json:"silhouette"`
	ControlsPass bool    `json:"controls_pass"`
}

// BootstrapSummary summarizes the bootstrap ARI distribution at ChosenK.
type BootstrapSummary struct {
	Rounds  int     `json:"rounds"`
	Seed    int64   `json:"seed"`
	MeanARI float64 `json:"mean_ari"`
	P2_5    float64 `json:"p2_5"`
	P97_5   float64 `json:"p97_5"`
}

// GridLabels stores the full labeling at one grid cut, aligned with IDs.
type GridLabels struct {
	K      int   `json:"k"`
	Labels []int `json:"labels"`
}

// Result is the persisted output of one clustering run. ChosenK/Assignments
// are the FINE cut (max silhouette subject to controls); GridLabels keeps
// every evaluated cut so cmd/matrix can select the cross-model consensus
// (component) cut without re-clustering.
type Result struct {
	Model       string            `json:"model"`
	Embeddings  string            `json:"embeddings"`
	ConfigHash  string            `json:"config_hash"`
	N           int               `json:"n"`
	Grid        []GridPoint       `json:"grid"`
	ChosenK     int               `json:"chosen_k"`
	Bootstrap   *BootstrapSummary `json:"bootstrap,omitempty"`
	Assignments map[string]int    `json:"assignments"`
	IDs         []string          `json:"ids"`
	GridLabels  []GridLabels      `json:"grid_labels"`
}

// LabelsAt returns the id→cluster map at grid cut k, or false if k was not
// evaluated.
func (r *Result) LabelsAt(k int) (map[string]int, bool) {
	for _, gl := range r.GridLabels {
		if gl.K == k && len(gl.Labels) == len(r.IDs) {
			m := make(map[string]int, len(r.IDs))
			for i, id := range r.IDs {
				m[id] = gl.Labels[i]
			}
			return m, true
		}
	}
	return nil, false
}

// ControlsPassAt reports whether the control assertions held at grid cut k.
func (r *Result) ControlsPassAt(k int) bool {
	for _, g := range r.Grid {
		if g.K == k {
			return g.ControlsPass
		}
	}
	return false
}

func (r *Result) Save(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	raw, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(raw, '\n'), 0o644)
}

func LoadResult(path string) (*Result, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var r Result
	if err := json.Unmarshal(raw, &r); err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	return &r, nil
}

// LoadResultsDir loads every *.json clustering result in dir, sorted by model
// name for determinism.
func LoadResultsDir(dir string) ([]*Result, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var out []*Result
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		r, err := LoadResult(filepath.Join(dir, e.Name()))
		if err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("no clustering results in %s", dir)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Model < out[j].Model })
	return out, nil
}

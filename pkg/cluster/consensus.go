package cluster

import (
	"fmt"
	"sort"
)

// ConsensusPoint is the cross-model agreement at one grid cut.
type ConsensusPoint struct {
	K            int     `json:"k"`
	MeanARI      float64 `json:"mean_ari"`
	ControlsPass bool    `json:"controls_pass"` // true iff controls hold for EVERY model at this k
}

// Consensus selects the COMPONENT cut: the grid k at which the models'
// partitions agree most (max mean pairwise ARI over shared segment ids),
// restricted to cuts where every model's control assertions hold. Ties break
// toward smaller k. This is the pre-registered rule
// "max-cross-model-ari-subject-to-controls": component granularity is the
// granularity at which independent embedding models agree most.
func Consensus(results []*Result) (int, []ConsensusPoint, error) {
	if len(results) < 2 {
		return 0, nil, fmt.Errorf("consensus needs >=2 models, got %d", len(results))
	}

	// ks evaluated by every model
	kCount := map[int]int{}
	for _, r := range results {
		for _, gl := range r.GridLabels {
			kCount[gl.K]++
		}
	}
	var ks []int
	for k, c := range kCount {
		if c == len(results) {
			ks = append(ks, k)
		}
	}
	sort.Ints(ks)
	if len(ks) == 0 {
		return 0, nil, fmt.Errorf("models share no grid cuts")
	}

	// shared ids across all models
	shared := map[string]int{}
	for _, r := range results {
		for _, id := range r.IDs {
			shared[id]++
		}
	}
	var ids []string
	for id, c := range shared {
		if c == len(results) {
			ids = append(ids, id)
		}
	}
	sort.Strings(ids)
	if len(ids) < 3 {
		return 0, nil, fmt.Errorf("models share only %d segment ids", len(ids))
	}

	bestK, bestARI := 0, -2.0
	var points []ConsensusPoint
	for _, k := range ks {
		pass := true
		labelVecs := make([][]int, len(results))
		for i, r := range results {
			if !r.ControlsPassAt(k) {
				pass = false
			}
			m, ok := r.LabelsAt(k)
			if !ok {
				pass = false
				continue
			}
			vec := make([]int, len(ids))
			for t, id := range ids {
				vec[t] = m[id]
			}
			labelVecs[i] = vec
		}
		var sum float64
		var nPairs int
		for i := 0; i < len(labelVecs); i++ {
			for j := i + 1; j < len(labelVecs); j++ {
				if labelVecs[i] == nil || labelVecs[j] == nil {
					continue
				}
				sum += ARI(labelVecs[i], labelVecs[j])
				nPairs++
			}
		}
		if nPairs == 0 {
			continue
		}
		mean := sum / float64(nPairs)
		points = append(points, ConsensusPoint{K: k, MeanARI: mean, ControlsPass: pass})
		if pass && mean > bestARI {
			bestARI, bestK = mean, k
		}
	}
	if bestK == 0 {
		return 0, points, fmt.Errorf("no shared cut satisfies all models' control assertions")
	}
	return bestK, points, nil
}

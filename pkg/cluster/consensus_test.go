package cluster

import "testing"

func mkResult(model string, ids []string, grid map[int][]int, controlsFailAt map[int]bool) *Result {
	r := &Result{Model: model, IDs: ids}
	for k, labels := range grid {
		r.GridLabels = append(r.GridLabels, GridLabels{K: k, Labels: labels})
		r.Grid = append(r.Grid, GridPoint{K: k, ControlsPass: !controlsFailAt[k]})
	}
	return r
}

func TestConsensusPicksMaxAgreement(t *testing.T) {
	ids := []string{"a", "b", "c", "d"}
	m1 := mkResult("m1", ids, map[int][]int{
		2: {0, 0, 1, 1},
		3: {0, 0, 1, 2},
	}, nil)
	m2 := mkResult("m2", ids, map[int][]int{
		2: {1, 1, 0, 0}, // identical partition to m1 up to relabeling → ARI 1
		3: {0, 1, 1, 2}, // disagrees with m1 at k=3
	}, nil)
	k, points, err := Consensus([]*Result{m1, m2})
	if err != nil {
		t.Fatal(err)
	}
	if k != 2 {
		t.Fatalf("consensus k = %d, want 2 (max agreement)", k)
	}
	if len(points) != 2 {
		t.Fatalf("want 2 consensus points, got %d", len(points))
	}
}

func TestConsensusRespectsControls(t *testing.T) {
	ids := []string{"a", "b", "c", "d"}
	m1 := mkResult("m1", ids, map[int][]int{
		2: {0, 0, 1, 1},
		3: {0, 0, 1, 2},
	}, map[int]bool{2: true}) // controls fail at the best-agreement cut
	m2 := mkResult("m2", ids, map[int][]int{
		2: {1, 1, 0, 0},
		3: {0, 0, 1, 2},
	}, nil)
	k, _, err := Consensus([]*Result{m1, m2})
	if err != nil {
		t.Fatal(err)
	}
	if k != 3 {
		t.Fatalf("consensus k = %d, want 3 (k=2 excluded by controls)", k)
	}
}

func TestConsensusNeedsTwoModels(t *testing.T) {
	m1 := mkResult("m1", []string{"a", "b", "c"}, map[int][]int{2: {0, 0, 1}}, nil)
	if _, _, err := Consensus([]*Result{m1}); err == nil {
		t.Fatal("want error for a single model")
	}
}

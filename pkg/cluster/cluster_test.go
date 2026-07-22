package cluster

import (
	"math"
	"testing"
)

func almostEqual(a, b, eps float64) bool { return math.Abs(a-b) <= eps }

func TestCosineSim(t *testing.T) {
	if got := CosineSim([]float32{1, 0}, []float32{1, 0}); !almostEqual(got, 1, 1e-9) {
		t.Fatalf("identical vectors: got %v, want 1", got)
	}
	if got := CosineSim([]float32{1, 0}, []float32{0, 1}); !almostEqual(got, 0, 1e-9) {
		t.Fatalf("orthogonal vectors: got %v, want 0", got)
	}
	if got := CosineSim([]float32{1, 1}, []float32{1, 1}); !almostEqual(got, 1, 1e-9) {
		t.Fatalf("parallel vectors: got %v, want 1", got)
	}
}

// two tight pairs far apart: cutting at k=2 must recover them.
func twoPairMatrix() *Matrix {
	m := NewMatrix(4)
	m.Set(0, 1, 0.10)
	m.Set(2, 3, 0.12)
	m.Set(0, 2, 1.00)
	m.Set(0, 3, 1.01)
	m.Set(1, 2, 1.02)
	m.Set(1, 3, 1.03)
	return m
}

func TestAgglomerativeCut(t *testing.T) {
	m := twoPairMatrix()
	merges := Agglomerative(m)
	if len(merges) != 3 {
		t.Fatalf("want 3 merges for n=4, got %d", len(merges))
	}
	labels := Cut(merges, 4, 2)
	if labels[0] != labels[1] || labels[2] != labels[3] || labels[0] == labels[2] {
		t.Fatalf("k=2 should recover the two pairs, got %v", labels)
	}
	if l := Cut(merges, 4, 1); l[0] != l[1] || l[1] != l[2] || l[2] != l[3] {
		t.Fatalf("k=1 should merge everything, got %v", l)
	}
	l4 := Cut(merges, 4, 4)
	seen := map[int]bool{}
	for _, x := range l4 {
		seen[x] = true
	}
	if len(seen) != 4 {
		t.Fatalf("k=4 should keep all points distinct, got %v", l4)
	}
}

func TestAgglomerativeDeterminism(t *testing.T) {
	a := Agglomerative(twoPairMatrix())
	b := Agglomerative(twoPairMatrix())
	for i := range a {
		if a[i] != b[i] {
			t.Fatalf("non-deterministic merges at %d: %+v vs %+v", i, a[i], b[i])
		}
	}
}

func TestSilhouette(t *testing.T) {
	// within-cluster distance 0.1, across 1.0 → s = (1.0-0.1)/1.0 = 0.9 for
	// every point.
	m := NewMatrix(4)
	m.Set(0, 1, 0.1)
	m.Set(2, 3, 0.1)
	for _, p := range [][2]int{{0, 2}, {0, 3}, {1, 2}, {1, 3}} {
		m.Set(p[0], p[1], 1.0)
	}
	got := Silhouette(m, []int{0, 0, 1, 1})
	if !almostEqual(got, 0.9, 1e-9) {
		t.Fatalf("silhouette: got %v, want 0.9", got)
	}
	if got := Silhouette(m, []int{0, 0, 0, 0}); got != 0 {
		t.Fatalf("single cluster silhouette must be 0, got %v", got)
	}
}

func TestARI(t *testing.T) {
	if got := ARI([]int{0, 0, 1, 1}, []int{1, 1, 0, 0}); !almostEqual(got, 1, 1e-9) {
		t.Fatalf("label-permuted identical partitions: got %v, want 1", got)
	}
	// Hand-computed: contingency all-ones 2x2 → sumCells=0, sumA=sumB=2,
	// expected=4/6, max=2 → ARI=(0-2/3)/(2-2/3) = -0.5.
	if got := ARI([]int{0, 0, 1, 1}, []int{0, 1, 0, 1}); !almostEqual(got, -0.5, 1e-9) {
		t.Fatalf("crossed partitions: got %v, want -0.5", got)
	}
}

func TestNMI(t *testing.T) {
	if got := NMI([]int{0, 0, 1, 1}, []int{1, 1, 0, 0}); !almostEqual(got, 1, 1e-9) {
		t.Fatalf("identical partitions: got %v, want 1", got)
	}
	if got := NMI([]int{0, 0, 1, 1}, []int{0, 1, 0, 1}); !almostEqual(got, 0, 1e-9) {
		t.Fatalf("independent partitions: got %v, want 0", got)
	}
}

func TestBootstrapARIDeterminism(t *testing.T) {
	m := twoPairMatrix()
	full := Cut(Agglomerative(m), 4, 2)
	a := BootstrapARI(m, full, 2, 20, 42)
	b := BootstrapARI(m, full, 2, 20, 42)
	if len(a) != len(b) {
		t.Fatalf("round counts differ: %d vs %d", len(a), len(b))
	}
	for i := range a {
		if a[i] != b[i] {
			t.Fatalf("bootstrap not deterministic at round %d", i)
		}
	}
}

func TestPercentile(t *testing.T) {
	xs := []float64{1, 2, 3, 4, 5}
	if got := Percentile(xs, 50); !almostEqual(got, 3, 1e-9) {
		t.Fatalf("median: got %v, want 3", got)
	}
	if got := Percentile(xs, 100); !almostEqual(got, 5, 1e-9) {
		t.Fatalf("p100: got %v, want 5", got)
	}
}

// Package cluster implements the deterministic clustering core of the
// pipeline: cosine distance, average-linkage agglomerative clustering
// (nearest-neighbor-chain algorithm — no random seed anywhere), dendrogram
// cutting, silhouette scoring, and the ARI/NMI/bootstrap stability metrics.
// Everything is pure Go and covered by fixture tests so a reviewer can audit
// the math directly.
package cluster

import (
	"math"
	"sort"
)

// Matrix is a dense symmetric distance matrix.
type Matrix struct {
	N int
	D []float64 // row-major N×N
}

func NewMatrix(n int) *Matrix { return &Matrix{N: n, D: make([]float64, n*n)} }

func (m *Matrix) At(i, j int) float64 { return m.D[i*m.N+j] }

func (m *Matrix) Set(i, j int, v float64) {
	m.D[i*m.N+j] = v
	m.D[j*m.N+i] = v
}

// Sub extracts the submatrix over idx (indices into the original matrix).
func (m *Matrix) Sub(idx []int) *Matrix {
	s := NewMatrix(len(idx))
	for a, i := range idx {
		for b, j := range idx {
			s.D[a*s.N+b] = m.At(i, j)
		}
	}
	return s
}

// CosineSim returns the cosine similarity of two vectors.
func CosineSim(a, b []float32) float64 {
	var dot, na, nb float64
	for i := range a {
		dot += float64(a[i]) * float64(b[i])
		na += float64(a[i]) * float64(a[i])
		nb += float64(b[i]) * float64(b[i])
	}
	if na == 0 || nb == 0 {
		return 0
	}
	return dot / math.Sqrt(na*nb)
}

// CosineDistances builds the pairwise cosine-distance matrix (1 − cosine
// similarity, clamped at 0).
func CosineDistances(vecs [][]float32) *Matrix {
	n := len(vecs)
	m := NewMatrix(n)
	for i := range n {
		for j := i + 1; j < n; j++ {
			d := 1 - CosineSim(vecs[i], vecs[j])
			if d < 0 {
				d = 0
			}
			m.Set(i, j, d)
		}
	}
	return m
}

// Merge is one dendrogram node. A and B are cluster ids (0..n-1 are original
// points; each merge t in recorded order creates cluster id n+t = ID).
type Merge struct {
	A, B   int
	ID     int
	Height float64
	Size   int
}

// Agglomerative runs average-linkage clustering with the nearest-neighbor
// chain algorithm. Deterministic: ties break toward the smallest slot index.
func Agglomerative(m *Matrix) []Merge {
	n := m.N
	if n < 2 {
		return nil
	}
	d := make([]float64, len(m.D))
	copy(d, m.D)
	at := func(i, j int) float64 { return d[i*n+j] }
	set := func(i, j int, v float64) { d[i*n+j], d[j*n+i] = v, v }

	size := make([]int, n)
	active := make([]bool, n)
	id := make([]int, n)
	for i := 0; i < n; i++ {
		size[i], active[i], id[i] = 1, true, i
	}
	next := n
	merges := make([]Merge, 0, n-1)
	chain := make([]int, 0, n)
	remaining := n

	nearest := func(a int) (int, float64) {
		b, best := -1, math.Inf(1)
		for j := 0; j < n; j++ {
			if !active[j] || j == a {
				continue
			}
			if v := at(a, j); v < best {
				best, b = v, j
			}
		}
		return b, best
	}

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
			b, best := nearest(a)
			if len(chain) >= 2 && b == chain[len(chain)-2] {
				chain = chain[:len(chain)-2]
				i, j := a, b
				if j < i {
					i, j = j, i
				}
				ni, nj := size[i], size[j]
				for k := 0; k < n; k++ {
					if !active[k] || k == i || k == j {
						continue
					}
					set(i, k, (float64(ni)*at(i, k)+float64(nj)*at(j, k))/float64(ni+nj))
				}
				merges = append(merges, Merge{A: id[i], B: id[j], ID: next, Height: best, Size: ni + nj})
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

// Cut assigns each of the n original points to one of k clusters by applying
// the n−k lowest merges (stable order: height, then recorded order — valid
// because average linkage is reducible, so children precede parents). Cluster
// labels are renumbered 0..k-1 in order of first appearance by point index.
func Cut(merges []Merge, n, k int) []int {
	if k < 1 {
		k = 1
	}
	if k > n {
		k = n
	}
	type ord struct {
		m   Merge
		pos int
	}
	sorted := make([]ord, len(merges))
	for i, m := range merges {
		sorted[i] = ord{m, i}
	}
	sort.SliceStable(sorted, func(i, j int) bool {
		if sorted[i].m.Height != sorted[j].m.Height {
			return sorted[i].m.Height < sorted[j].m.Height
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
	rep := make(map[int]int, n+len(merges)) // cluster id → a member point
	for i := 0; i < n; i++ {
		rep[i] = i
	}
	for t := 0; t < n-k && t < len(sorted); t++ {
		m := sorted[t].m
		ra, rb := find(rep[m.A]), find(rep[m.B])
		parent[rb] = ra
		rep[m.ID] = ra
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

// Silhouette computes the mean silhouette coefficient over all points on a
// precomputed distance matrix. Points in singleton clusters score 0 (the
// scikit-learn convention).
func Silhouette(m *Matrix, labels []int) float64 {
	n := m.N
	if n == 0 {
		return 0
	}
	k := 0
	for _, l := range labels {
		if l+1 > k {
			k = l + 1
		}
	}
	if k < 2 {
		return 0
	}
	sizes := make([]int, k)
	for _, l := range labels {
		sizes[l]++
	}
	var total float64
	for i := 0; i < n; i++ {
		li := labels[i]
		if sizes[li] <= 1 {
			continue // s(i) = 0
		}
		sums := make([]float64, k)
		for j := 0; j < n; j++ {
			if j != i {
				sums[labels[j]] += m.At(i, j)
			}
		}
		a := sums[li] / float64(sizes[li]-1)
		b := math.Inf(1)
		for c := 0; c < k; c++ {
			if c == li || sizes[c] == 0 {
				continue
			}
			if v := sums[c] / float64(sizes[c]); v < b {
				b = v
			}
		}
		if max := math.Max(a, b); max > 0 {
			total += (b - a) / max
		}
	}
	return total / float64(n)
}

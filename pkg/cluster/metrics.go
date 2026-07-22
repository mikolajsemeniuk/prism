package cluster

import (
	"math"
	"math/rand"
	"sort"
)

// ARI computes the Adjusted Rand Index between two labelings of the same
// points. 1 = identical partitions, ~0 = chance agreement.
func ARI(x, y []int) float64 {
	n := len(x)
	if n == 0 || n != len(y) {
		return math.NaN()
	}
	cont, rows, cols := contingency(x, y)
	var sumCells, sumA, sumB float64
	for _, row := range cont {
		for _, v := range row {
			sumCells += comb2(v)
		}
	}
	for _, v := range rows {
		sumA += comb2(v)
	}
	for _, v := range cols {
		sumB += comb2(v)
	}
	expected := sumA * sumB / comb2(n)
	max := (sumA + sumB) / 2
	if max == expected {
		return 1
	}
	return (sumCells - expected) / (max - expected)
}

// NMI computes normalized mutual information with arithmetic-mean
// normalization (the scikit-learn default).
func NMI(x, y []int) float64 {
	n := len(x)
	if n == 0 || n != len(y) {
		return math.NaN()
	}
	cont, rows, cols := contingency(x, y)
	fn := float64(n)
	var mi float64
	for i, row := range cont {
		for j, v := range row {
			if v == 0 {
				continue
			}
			pij := float64(v) / fn
			mi += pij * math.Log(pij/(float64(rows[i])/fn*float64(cols[j])/fn))
		}
	}
	hx := entropy(rows, n)
	hy := entropy(cols, n)
	if hx == 0 && hy == 0 {
		return 1
	}
	if hx == 0 || hy == 0 {
		return 0
	}
	return mi / ((hx + hy) / 2)
}

// BootstrapARI reclusters bootstrap resamples of the points at fixed k and
// returns the ARI of each resample's labeling against the full labeling,
// restricted to the sampled points. Deterministic given seed.
func BootstrapARI(m *Matrix, full []int, k, rounds int, seed int64) []float64 {
	rng := rand.New(rand.NewSource(seed))
	out := make([]float64, 0, rounds)
	for r := 0; r < rounds; r++ {
		picked := map[int]bool{}
		for i := 0; i < m.N; i++ {
			picked[rng.Intn(m.N)] = true
		}
		idx := make([]int, 0, len(picked))
		for i := range picked {
			idx = append(idx, i)
		}
		sort.Ints(idx)
		if len(idx) < 3 {
			continue
		}
		sub := m.Sub(idx)
		kk := k
		if kk > len(idx) {
			kk = len(idx)
		}
		labels := Cut(Agglomerative(sub), len(idx), kk)
		ref := make([]int, len(idx))
		for a, i := range idx {
			ref[a] = full[i]
		}
		out = append(out, ARI(labels, ref))
	}
	return out
}

// Percentile returns the p-th percentile (0..100) of xs by linear
// interpolation; xs need not be sorted.
func Percentile(xs []float64, p float64) float64 {
	if len(xs) == 0 {
		return math.NaN()
	}
	s := append([]float64(nil), xs...)
	sort.Float64s(s)
	if len(s) == 1 {
		return s[0]
	}
	pos := p / 100 * float64(len(s)-1)
	lo := int(math.Floor(pos))
	hi := int(math.Ceil(pos))
	if lo == hi {
		return s[lo]
	}
	return s[lo] + (pos-float64(lo))*(s[hi]-s[lo])
}

func Mean(xs []float64) float64 {
	if len(xs) == 0 {
		return math.NaN()
	}
	var sum float64
	for _, x := range xs {
		sum += x
	}
	return sum / float64(len(xs))
}

func contingency(x, y []int) (map[int]map[int]int, map[int]int, map[int]int) {
	cont := map[int]map[int]int{}
	rows := map[int]int{}
	cols := map[int]int{}
	for i := range x {
		if cont[x[i]] == nil {
			cont[x[i]] = map[int]int{}
		}
		cont[x[i]][y[i]]++
		rows[x[i]]++
		cols[y[i]]++
	}
	return cont, rows, cols
}

func entropy(counts map[int]int, n int) float64 {
	var h float64
	for _, c := range counts {
		if c == 0 {
			continue
		}
		p := float64(c) / float64(n)
		h -= p * math.Log(p)
	}
	return h
}

func comb2(v int) float64 { return float64(v) * float64(v-1) / 2 }

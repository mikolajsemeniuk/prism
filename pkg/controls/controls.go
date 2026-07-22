// Package controls evaluates the pre-registered clustering assertions from
// pipeline.json against a clustering assignment. Shared by cmd/cluster
// (where the assertions constrain the pre-registered k selection) and
// cmd/matrix (where they are re-checked and reported).
package controls

import (
	"fmt"
	"sort"
	"strings"

	"github.com/mikolajsemeniuk/prism/pkg/config"
	"github.com/mikolajsemeniuk/prism/pkg/segment"
)

// Result is one control evaluated against one assignment.
type Result struct {
	Name    string              `json:"name"`
	Expect  string              `json:"expect"`
	Matches map[string][]string `json:"matches"`
	Pass    bool                `json:"pass"`
	Detail  string              `json:"detail,omitempty"`
}

// Evaluate checks one control against a clustering assignment (segment id →
// cluster). byID must cover every assigned id.
func Evaluate(c config.Control, assignments map[string]int, byID map[string]segment.Segment) Result {
	r := Result{Name: c.Name, Expect: c.Expect, Matches: map[string][]string{}}
	find := func(needle string) (ids []string, clusters map[int]bool, files map[string]bool) {
		clusters = map[int]bool{}
		files = map[string]bool{}
		for id, cl := range assignments {
			s, ok := byID[id]
			if !ok || !strings.Contains(s.Text, needle) {
				continue
			}
			ids = append(ids, id)
			clusters[cl] = true
			files[s.Source] = true
		}
		sort.Strings(ids)
		return
	}
	switch c.Expect {
	case "same":
		ids, clusters, files := find(c.Needle)
		r.Matches["needle"] = ids
		switch {
		case len(files) < 2:
			r.Detail = fmt.Sprintf("needle matched %d file(s), need >=2", len(files))
		case len(clusters) != 1:
			r.Detail = fmt.Sprintf("matched segments span %d clusters", len(clusters))
		default:
			r.Pass = true
		}
	case "different":
		idsA, clustersA, _ := find(c.NeedleA)
		idsB, clustersB, _ := find(c.NeedleB)
		r.Matches["needle_a"] = idsA
		r.Matches["needle_b"] = idsB
		if len(idsA) == 0 || len(idsB) == 0 {
			r.Detail = "a needle matched no segments"
			return r
		}
		r.Pass = true
		for cl := range clustersA {
			if clustersB[cl] {
				r.Pass = false
				r.Detail = fmt.Sprintf("cluster %d contains both needles", cl)
				break
			}
		}
	default:
		r.Detail = fmt.Sprintf("unknown expect %q", c.Expect)
	}
	return r
}

// AllPass evaluates every control and reports whether all pass.
func AllPass(cs []config.Control, assignments map[string]int, byID map[string]segment.Segment) bool {
	for _, c := range cs {
		if !Evaluate(c, assignments, byID).Pass {
			return false
		}
	}
	return true
}

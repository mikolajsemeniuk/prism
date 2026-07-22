package embed

import (
	"context"
	"testing"

	"github.com/mikolajsemeniuk/prism/pkg/cluster"
)

func TestFakeDeterministicAndDiscriminative(t *testing.T) {
	f := &Fake{Model: "smoke-a", Dims: 256}
	texts := []string{
		"Claude cares deeply about child safety and exercises special caution with minors.",
		"Claude cares deeply about child safety and exercises special caution with minors.",
		"Claude cares deeply about child safety and exercises extra caution with minors.",
		"prefer using rg or rg --files because rg is much faster than alternatives like grep",
	}
	vecs, err := f.Embed(context.Background(), texts)
	if err != nil {
		t.Fatal(err)
	}
	if sim := cluster.CosineSim(vecs[0], vecs[1]); sim < 0.9999 {
		t.Fatalf("identical texts similarity = %v, want ~1", sim)
	}
	near := cluster.CosineSim(vecs[0], vecs[2])
	far := cluster.CosineSim(vecs[0], vecs[3])
	if near <= far {
		t.Fatalf("near-duplicate (%v) must beat unrelated (%v)", near, far)
	}
	if near < 0.7 {
		t.Fatalf("near-duplicate similarity too low: %v", near)
	}
	if far > 0.5 {
		t.Fatalf("unrelated similarity too high: %v", far)
	}
}

func TestFakeModelSeedsDiffer(t *testing.T) {
	a := &Fake{Model: "smoke-a", Dims: 256}
	b := &Fake{Model: "smoke-b", Dims: 256}
	va, _ := a.Embed(context.Background(), []string{"hello world example"})
	vb, _ := b.Embed(context.Background(), []string{"hello world example"})
	if sim := cluster.CosineSim(va[0], vb[0]); sim > 0.999 {
		t.Fatalf("different fake models should produce different spaces, sim=%v", sim)
	}
}

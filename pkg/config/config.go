// Package config loads pipeline.json, the pre-registered parameter set for
// the whole pipeline. Commands hash the raw bytes into their output manifests
// so every derived artifact is traceable to the exact configuration.
package config

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Segment struct {
		MinChars      int      `json:"min_chars"`
		IncludeLayers []string `json:"include_layers"`
	} `json:"segment"`
	Embed struct {
		Endpoint  string `json:"endpoint"`
		BatchSize int    `json:"batch_size"`
		FakeDims  int    `json:"fake_dims"`
	} `json:"embed"`
	Cluster struct {
		Linkage   string `json:"linkage"`
		Metric    string `json:"metric"`
		KGrid     []int  `json:"k_grid"`
		Selection string `json:"selection"`
		Bootstrap struct {
			Rounds int   `json:"rounds"`
			Seed   int64 `json:"seed"`
		} `json:"bootstrap"`
	} `json:"cluster"`
	Taxonomy struct {
		StabilityMin  float64 `json:"stability_min"`
		MinProducts int     `json:"min_products"`
		MinSize     int     `json:"min_size"`
		TopTerms    int     `json:"top_terms"`
		Exemplars   int     `json:"exemplars"`
	} `json:"taxonomy"`
	Controls []Control `json:"controls"`

	// Hash is the sha256 of the raw config file, for manifests.
	Hash string `json:"-"`
}

// Control is a pre-registered clustering assertion: segments matching Needle
// (expect "same") must land in one cluster; segments matching NeedleA vs
// NeedleB (expect "different") must not share a cluster.
type Control struct {
	Name    string `json:"name"`
	Expect  string `json:"expect"`
	Needle  string `json:"needle,omitempty"`
	NeedleA string `json:"needle_a,omitempty"`
	NeedleB string `json:"needle_b,omitempty"`
	Note    string `json:"note,omitempty"`
}

func Load(path string) (Config, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}
	var c Config
	if err := json.Unmarshal(raw, &c); err != nil {
		return Config{}, fmt.Errorf("%s: %w", path, err)
	}
	sum := sha256.Sum256(raw)
	c.Hash = hex.EncodeToString(sum[:])
	if c.Segment.MinChars <= 0 {
		c.Segment.MinChars = 40
	}
	if c.Embed.BatchSize <= 0 {
		c.Embed.BatchSize = 16
	}
	if c.Embed.FakeDims <= 0 {
		c.Embed.FakeDims = 256
	}
	if len(c.Cluster.KGrid) == 0 {
		return Config{}, fmt.Errorf("%s: cluster.k_grid must not be empty", path)
	}
	if c.Taxonomy.StabilityMin <= 0 {
		c.Taxonomy.StabilityMin = 0.5
	}
	if c.Taxonomy.MinProducts <= 0 {
		c.Taxonomy.MinProducts = 3
	}
	if c.Taxonomy.MinSize <= 0 {
		c.Taxonomy.MinSize = 5
	}
	if c.Taxonomy.TopTerms <= 0 {
		c.Taxonomy.TopTerms = 8
	}
	if c.Taxonomy.Exemplars <= 0 {
		c.Taxonomy.Exemplars = 3
	}
	return c, nil
}

// Command segment parses the frozen corpus into derived/segments.jsonl
// (roadmap step 2). It never writes under corpus/.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/mikolajsemeniuk/prism/pkg/config"
	"github.com/mikolajsemeniuk/prism/pkg/corpus"
	"github.com/mikolajsemeniuk/prism/pkg/segment"
)

func main() {
	log.SetFlags(0)
	cfgPath := flag.String("config", "pipeline.json", "pre-registered pipeline config")
	corpusDir := flag.String("corpus", "corpus", "corpus root (read-only)")
	out := flag.String("out", "derived/segments.jsonl", "output JSONL path")
	flag.Parse()

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		log.Fatalf("config: %v", err)
	}
	arts, err := corpus.Load(*corpusDir)
	if err != nil {
		log.Fatalf("corpus: %v", err)
	}
	if len(arts) == 0 {
		log.Fatalf("no artifacts found under %s", *corpusDir)
	}

	var all []segment.Segment
	w := tabwriter.NewWriter(os.Stdout, 2, 4, 2, ' ', 0)
	fmt.Fprintln(w, "ARTIFACT\tCLASS\tSEGS\tBEHAV\tSCHEMA\tFEWSHOT\tTMPL\tWRAP")
	layerTotals := map[segment.Layer]int{}
	for _, a := range arts {
		segs := segment.Split(a, segment.Options{MinChars: cfg.Segment.MinChars})
		counts := map[segment.Layer]int{}
		for _, s := range segs {
			counts[s.Layer]++
			layerTotals[s.Layer]++
		}
		class := "?"
		if len(segs) > 0 {
			class = segs[0].ProductClass
			if class == "unknown" {
				log.Printf("WARNING: %s has unknown product class — extend pkg/segment/meta.go", a.Rel)
			}
		}
		fmt.Fprintf(w, "%s\t%s\t%d\t%d\t%d\t%d\t%d\t%d\n", a.Rel, class, len(segs),
			counts[segment.Behavioral], counts[segment.ToolSchema], counts[segment.FewShot],
			counts[segment.TemplateVar], counts[segment.Wrapper])
		all = append(all, segs...)
	}
	w.Flush()
	fmt.Printf("\nTOTAL %d segments  behavioral=%d tool-schema=%d few-shot=%d template-var=%d wrapper=%d\n",
		len(all), layerTotals[segment.Behavioral], layerTotals[segment.ToolSchema],
		layerTotals[segment.FewShot], layerTotals[segment.TemplateVar], layerTotals[segment.Wrapper])

	if err := segment.WriteJSONL(*out, all); err != nil {
		log.Fatalf("write %s: %v", *out, err)
	}
	fmt.Printf("wrote %s (config %s)\n", *out, cfg.Hash[:12])
}

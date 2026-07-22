# Prism pipeline. Each target is one small cmd doing exactly one job;
# `make paper` chains corpus → segments → embeddings → clusters → matrix →
# paper/*.gen.tex. `make smoke` runs the whole chain offline on the fake
# backend and fails if any pre-registered control assertion fails.

GO        ?= go
CONFIG    ?= pipeline.json
BACKEND   ?= ollama
MODELS    ?= bge-m3,nomic-embed-text
BOOTSTRAP ?= -1            # -1 = use the pre-registered config value
DERIVED   ?= derived
PAPER     ?= paper

.PHONY: segment embed cluster matrix gentex paper smoke clean

segment:
	$(GO) run ./cmd/segment -config $(CONFIG) -corpus corpus -out $(DERIVED)/segments.jsonl

embed: segment
	@for m in $$(echo $(MODELS) | tr ',' ' '); do \
		$(GO) run ./cmd/embed -config $(CONFIG) -segments $(DERIVED)/segments.jsonl \
			-backend $(BACKEND) -model $$m -outdir $(DERIVED)/embeddings || exit 1; \
	done

cluster: embed
	@for m in $$(echo $(MODELS) | tr ',' ' '); do \
		$(GO) run ./cmd/cluster -config $(CONFIG) -segments $(DERIVED)/segments.jsonl \
			-embeddings $(DERIVED)/embeddings/$$m.jsonl -bootstrap $(BOOTSTRAP) \
			-out $(DERIVED)/clusters/$$m.json || exit 1; \
	done

matrix: cluster
	$(GO) run ./cmd/matrix -config $(CONFIG) -segments $(DERIVED)/segments.jsonl \
		-clustersdir $(DERIVED)/clusters -outdir $(DERIVED)

gentex:
	$(GO) run ./cmd/gentex -config $(CONFIG) -segments $(DERIVED)/segments.jsonl \
		-clustersdir $(DERIVED)/clusters -derived $(DERIVED) -outdir $(PAPER)

paper: matrix gentex

# Offline end-to-end run: fake deterministic embedder, two simulated models,
# reduced bootstrap. Exercises every cmd and the control assertions.
smoke:
	$(MAKE) BACKEND=fake MODELS=smoke-a,smoke-b BOOTSTRAP=25 \
		DERIVED=derived/smoke PAPER=paper/smoke paper
	$(GO) test ./...

clean:
	rm -rf $(DERIVED) paper/smoke

# Prism decomposition pipeline (see CLAUDE.md "Decomposition methodology").
# One target per step, in dependency order; each step consumes the previous
# step's output in artefacts/:
#
#   segment → embed (×MODELS) → cluster (×MODELS) → kbounds → admit
#
# Only `embed` touches the network (local Ollama); its output file doubles as
# a cache, so re-runs are cheap. Everything else is deterministic offline
# math — changing a threshold re-runs only its suffix of the chain.

MODELS ?= bge-m3,nomic-embed-text

.PHONY: all build vet segment embed cluster kbounds admit clean clean-all

# step 1: corpus → artefacts/segments.jsonl
segment:
	go run ./cmd/segment

# step 2: segments → artefacts/embeddings-<model>.jsonl (one run per model)
embed:
	@for m in $$(echo $(MODELS) | tr ',' ' '); do \
		go run ./cmd/embed -model $$m || exit 1; \
	done

# step 3: embeddings → artefacts/clusters-<model>.json (one run per model)
cluster:
	@for m in $$(echo $(MODELS) | tr ',' ' '); do \
		go run ./cmd/cluster -model $$m || exit 1; \
	done

# step 4: controls exam → artefacts/kbounds.json (admissible k window)
kbounds:
	go run ./cmd/kbounds

# step 5: cross-model consensus → artefacts/admit.json (component cut k*)
admit:
	go run ./cmd/admit

clean:
	rm -rf artefacts

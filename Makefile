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

.PHONY: all build vet segment embed cluster kbounds admit taxonomy fragments gentex clean \
        vllm-up vllm-down vllm-logs smoke-runner pilot

# --- experiment (roadmap step 5): scenario runner + vLLM ---------------------
# One vLLM profile runs at a time on port 12000. PROFILE selects the model
# service; MODEL is its --served-model-name (must match). Examples:
#   make vllm-up PROFILE=qwen3  &&  make pilot MODEL=qwen3:30b-a3b
PROFILE  ?= qwen3
MODEL    ?= qwen3:30b-a3b
SCENARIO ?= scenarios/selfrepair-01
K        ?= 4

# start/stop the selected vLLM profile (see docker-compose.yaml)
vllm-up:
	docker compose --profile $(PROFILE) up -d

vllm-down:
	docker compose --profile $(PROFILE) down

vllm-logs:
	docker compose logs -f

# phase-0 smoke: replay a scripted "good" run — validates the harness (loop,
# real tool execution, success + manipulation-check detection) with zero API.
smoke-runner:
	go run ./cmd/runner -scenario $(SCENARIO) -scripted $(SCENARIO)/scripted-good.json -label smoke-good -k 1

# pilot: the four core conditions of one scenario × one model × K episodes.
# baseline (no prompt), grounding (env context), fragment (the C30 medoid),
# placebo (length-matched neutral text). Reveals direction + variance fast.
pilot:
	go run ./cmd/runner -scenario $(SCENARIO) -model $(MODEL) -label baseline -k $(K)
	go run ./cmd/runner -scenario $(SCENARIO) -model $(MODEL) -label grounding -grounding -k $(K)
	go run ./cmd/runner -scenario $(SCENARIO) -model $(MODEL) -label fragment \
		-system-file $(SCENARIO)/stimuli/fragment.txt -k $(K)
	go run ./cmd/runner -scenario $(SCENARIO) -model $(MODEL) -label placebo \
		-system-file $(SCENARIO)/stimuli/placebo.txt -k $(K)


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

# step 6: component taxonomy (steps A–D) → artefacts/taxonomy.json + .md
taxonomy:
	go run ./cmd/taxonomy

# step 8: paper tables (pure formatting) → paper/*.gen.tex
gentex:
	go run ./cmd/gentex

clean:
	rm -rf artefacts

# step 7: verbatim ablation stimuli → artefacts/fragments.jsonl
fragments:
	go run ./cmd/fragments

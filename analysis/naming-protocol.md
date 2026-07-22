# Component naming protocol (taxonomy step E)

How entries of `analysis/component-names.json` are produced from
`artefacts/taxonomy.json` / `taxonomy.md`. This is the ONLY interpretive step
of the decomposition methodology; everything upstream is mechanical. The
architecture guarantees names cannot influence any number: the pipeline never
reads the names file — `cmd/gentex` merges it at paper-format time only, and
both file hashes appear in the provenance header of `paper/taxonomy.gen.tex`.

## Input (per component entry — nothing outside the entry may be used)

1. **Medoid** — the cluster's most typical verbatim production sentence
   (argmax mean cosine similarity to the other members). Primary evidence.
2. **Exemplars** — the next most central members; confirm the medoid is
   representative, adjudicate when signals conflict.
3. **Distinctive terms** — tf-idf tokens (frequent in this cluster, rare in
   the others); a topical sanity check on the medoid.

## Rules

1. The name describes the BEHAVIORAL FUNCTION of the directives, never a
   vendor, product, or tool brand. Vendor tokens appearing in the terms list
   (e.g. "replit") are ignored — the name must fit every product the
   component spans.
2. Only the entry's own material may be used. Prior knowledge (the manual
   codebook, the ablation-candidates list) must NOT drive the name; that the
   resulting names often match the manual codebook is a finding, not the
   method.
3. Evidence precedence: medoid > exemplars > terms. If the medoid and terms
   disagree, exemplars adjudicate. If the entry remains unclear, it stays
   UNNAMED (rendered as "---" in the paper table) — naming under uncertainty
   is forbidden.
4. Form: kebab-case, at most three words. The definition is ONE sentence that
   paraphrases the medoid/exemplars without extending their scope.

## Worked examples (from the 2026-07-22 run)

- C31 — terms `debugging, problem, tests, error`; medoid literally contains
  "Address the root cause" → `root-cause-fix`.
- C30 — medoid "don't stop if the application is not installed. Instead,
  please install..." → `environment-self-repair`.
- C29 — terms contain vendor token `replit` (ignored, rule 1); medoid
  "Always speak in simple, everyday language. User is non-technical..." →
  `audience-adaptation`.

## Validation (when a reviewer asks for numbers)

- Reverse-matching test: give a third party the names and the shuffled
  medoids; a high re-pairing rate shows the names are faithful to the data.
- Stronger: two annotators name all entries independently under this
  protocol; report agreement and adjudicate differences.

## Maintenance

Cluster ids (C31, ...) are only meaningful for the current component cut:
after any re-run that changes the segmentation or k*, re-check every key
against the new `taxonomy.md` before trusting the names file.

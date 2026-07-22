# Cross-corpus findings — manual analysis (roadmap step 1)

Date: 2026-07-21. Method: qualitative inventory of all 17 corpus artifacts by four
parallel readers sharing a 27-slug seed codebook; per-artifact inventories with
line-referenced verbatim quotes live in `analysis/artifacts/` (one file per corpus
artifact — line refs below point into those). Category slugs are defined in
`codebook-draft.md`. Per the CLAUDE.md guardrail, everything here is hypothesis
material for validating the unsupervised clustering, not a component-selection
mechanism.

## F1. Layer composition varies by an order of magnitude

Behavioral share of artifact bodies ranges from ~4% (Claude Code) to ~62%
(Windsurf) and ~100% (Anthropic chat prompts, whose captures contain no tools).
Approximate byte/line shares:

| Artifact | Behavioral | Tool defs | Few-shot | Template/session | Other |
|---|---|---|---|---|---|
| Anthropic chat (×3) | ~100% | 0% | 0% | 0% | — |
| Codex CLI base | ~40% | ~15% | ~17% | — | ~29% presentation spec |
| Aider | ~12% | ~48% (edit grammar) | ~29% | — | ~11% |
| ChatGPT 5.6 | 9% | 79% | 10% | 1% | 1% |
| Gemini 3.5 | 47% | 41% | 5% | 5% | 2% |
| Grok 4.3 | 8% | 85% | 1% | 2% | 4% |
| Claude Code | ~4% | ~92% | 0% | ~2.5% | ~1.5% |
| Cursor 2.0 | ~33% | ~65% | (35% embedded) | ~1% | <1% |
| Windsurf | ~62% | 0% (external) | ~33% | ~5% | <1% |
| Devin | ~38% | ~50% | ~4% | ~7% | ~1% |
| Replit Agent | ~30% | ~56% | ~9% | ~2% | ~2% |
| v0 | ~52% | 0% (external) | ~22% | ~9% | ~13% domain-knowledge |

Implication (steps 2–3): embedding whole prompts would cluster by capture scope,
not by design. Layer tagging before embedding is mandatory, and the
component×product matrix must be computed per layer.

## F2. Behavioral policy hides inside tool schemas

Replit's strongest conduct policy (✓/→ progress format, no emojis, 30-word
summaries, one question at a time) lives in JSON `description` fields; Devin's
browser/editor doctrine sits inside command docs; ChatGPT's visibility-channel
policy is bound to tool namespaces. Implication: `tool-definition` ≠
non-behavioral. The step-2 schema needs a sublayer split
(`tool-definition-schema` vs `tool-definition-prose`), and schema-prose must be
embedded alongside behavioral segments.

## F3. Three similarity regimes — and free positive controls

The corpus exhibits (a) byte reuse within a vendor, (b) near-verbatim cross-vendor
copying, (c) convergent independent solutions. Manually verified duplicate or
near-duplicate passages:

| Passage | Where | Regime |
|---|---|---|
| child-safety core + four wellbeing passages | Fable 5 ↔ Opus 4.8 ↔ Opus 4.7 | within-vendor, byte-identical |
| rg-preference + file-reference grammar | Codex base ↔ GPT-5-Codex | within-vendor, byte-identical |
| `making_code_changes` checklist (incl. "beautiful and modern UI, imbued with best UX practices") | Cursor ↔ Windsurf | cross-vendor, near-byte |
| GDPR-Art.9-shaped sensitive-data lists | Gemini ↔ ChatGPT bio tool | cross-vendor, parallel boilerplate |
| Codex preamble guidance ("8–12 words", "friendly and curious tone") | OpenHands ← Codex base | cross-vendor, near-byte |
| Anthropic parallel-tool-call / "use that value EXACTLY" boilerplate | v0 (← Anthropic harness) | cross-vendor |
| str_replace "match EXACTLY ... Be mindful of whitespaces!" | Devin ↔ Replit (← Anthropic editor tool) | cross-vendor |
| SKILL.md progressive-disclosure convention | ChatGPT, Grok, Claude Code | copied convention |

These become the step-3 smoke-test assertions. **Caveat:** the Anthropic
child-safety control holds at *passage* granularity only — block membership churns
between versions (4.7 lacks the slang bullet; Fable 5 swaps two bullets). Assert
co-clustering of passages, not blocks.

## F4. Component distributions are product-class-conditional

Chat prompts: mental-health/wellbeing up to 31% of body (Fable 5), zero code
content. Coding agents: zero wellbeing text, and three of five Tier-B artifacts
contain no refusal text at all. A flat corpus-wide clustering will split along
product class before it splits along function. Implication: stratify the
component×product matrix by product class (chat / coding-agent / app-builder) and
report within-class and cross-class clusters separately.

## F5. The Anthropic trio is a natural ablation series

Opus 4.7 → 4.8 was condense-and-harden (adds `default_stance` anti-over-refusal,
CSAM-slang bullet, conventional-weapons=CBRN paragraph, tone micro-bans,
`tool_discovery`); 4.8 → Fable 5 was revert-and-deepen (drops most 4.8
experiments, grows wellbeing 17%→25%→31%, adds `end_conversation` single-warning
protocol). The vendor effectively ran add-then-revert experiments for us —
per-version component diffs are documented in `analysis/artifacts/anthropic--*`.
Side-finding: the weapons text weakens exactly when the prompt begins claiming
model-level safety measures — safety migrating from prompt into model.

## F6. Products ship contradictory directives for the same component

| Component | Opposing poles |
|---|---|
| memory-context-mgmt | Windsurf "create memories liberally" ↔ Cursor "DO NOT ... unless the user explicitly asks" |
| tool-parallelism | Cline "batch all independent calls simultaneously" ↔ OpenHands "combine multiple bash commands into one" (serialize) |
| model identity | Windsurf masks ("respond with `GPT 4.1`") ↔ Cursor asserts substrate ↔ Claude Code elaborates identity |
| persistence vs approval | Cline "Don't ask for permission!" ↔ Aider human-gated file access |
| audience-adaptation | Replit "user is non-technical" ↔ Devin note "usually an expert programmer" |
| math markup | Gemini "LaTeX only for formal/complex math" ↔ Grok "Always use KaTeX" |

Prime ablation targets: where the industry demonstrably disagrees, per-component
deltas answer a live design question rather than confirming folklore.

## F7. Prompts are accreted patch layers, not designed documents

Shipped typos/grammar errors in 5 of 8 Tier-A/B artifacts; ChatGPT contradicts
itself on consecutive lines (447–448) and carries stale internal dates; Opus 4.8
ends with an orphan tag outside its root element; Codex base references a
nonexistent section; Grok's anti-founder-deference bullet and GPT-5-Codex's git
protections read as post-incident patches. Paper-worthy descriptive finding, and a
known noise source for embeddings.

## F8. A model-adaptation layer exists

Codex ships per-model prompt files; OpenHands ships per-family `<IMPORTANT>`
patches pointing in opposite directions (Claude: "don't make extra or fewer
actions"; Gemini: "Avoid being too proactive"). Implication for step 5: "same
prompt, many models" ablations must strip or declare model-adaptation segments —
part of any prompt's effect is correction of one model's failure modes.

## F9. Tier C bodies are per-session serializations with capture damage

All three chatbot artifacts carry one extractor's session state (Iceland
location/timezone/name); Windsurf retains an end user's Windows path; Claude Code
retains the extractor's git identity; Replit is a three-document aggregate with
stripped XML tag names; v0 is truncated mid-sentence. Implications: step-2
corruption flags (`truncated | stripped-tags | aggregate-member | session-fill`),
a threats-to-validity paragraph ("serializations, not vendor masters"), and
cross-validation notes in the SOURCES.md TODO.

## F10. Formatting/verbosity is the most universal, most measurable component

Present in all 17 artifacts; the largest single policy area in Codex base (~29%)
and Opus 4.7 (~25%); frequently quantified (8–12-word preambles, 10-line answers,
3–5 colors, "oververbosity: 4", 30-word summaries, 15-second update cadence).
Best first target for programmatic scoring in the scenario suite, and the natural
material for length-matched placebo construction.

## F11. Requirements pushed into later pipeline steps

- **Step 2 (segmentation)** — segment record fields: id, source_file, char range,
  section_path, layer (`behavioral | tool-definition-schema | tool-definition-prose
  | few-shot | template-var | domain-knowledge | wrapper`), corruption flags,
  product_class, tier. Parser must handle: Markdown headers, XML sections, ChatML
  markers (Cursor), `====` separators (v0), TS/Python string literals (Tier B),
  single-line 20 KB JSON blobs (Replit), republisher wrappers (Devin/Replit).
- **Step 3 (clustering) smoke tests** — assert co-clustering of the F3 passage
  pairs (at minimum the four byte/near-byte ones); assert separation of unrelated
  categories; stability across ≥2 embedding models. Codex base vs GPT-5-Codex
  (same function, 3× compression) tests function-over-surface clustering.
- **Step 4 (scenarios)** — start from F10 (programmatic scoring) and the F6
  oppositions (memory, parallelism, persistence, audience).
- **Step 5 (ablation)** — strip/declare model-adaptation segments (F8); stratify
  by product class (F4); keep Tier C claims artifact-relative (F9).

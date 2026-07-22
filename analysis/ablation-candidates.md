# Ablation candidates — agent effectiveness in unfamiliar environments (v2)

Status: hypothesis pre-registration. v1 2026-07-21 (distilled from reader
summaries); v2 same day, after a dedicated two-agent verification sweep of all 17
per-artifact inventories against the research focus. Both sweep agents returned
explicit "checked, nothing found" lists for the remaining probe areas, so the
manual pass is considered saturated; the systematic completeness guarantee
remains step-3 clustering (an unnamed cluster = a missed component). Per the
CLAUDE.md guardrail this is not the component-selection mechanism — the final set
must be validated against step-3 clusters.

## Epistemic status of the corpus evidence

The corpus contains prescriptions, not outcomes. What it CAN establish:
prevalence, convergence-vs-copying (weighted via step-3 similarity regimes),
vendors' own add-then-revert experiments (Anthropic trio; GPT-5-Codex git
patches), and vendors' belief in model-dependent effects (per-model patch
layers). What it CANNOT establish: efficacy — that is steps 4–5's job.

v2 corrections to v1 claims (from the sweep):
- "environment grounding is universal in agentic captures" was overstated: Aider
  has none, and the Codex base prompt file carries none (any injection is
  harness-side). Reworded: universal among captures whose harness serializes
  session state.
- explore-before-act regraded E1 → **E2**: genuinely contested (Cursor "be
  THOROUGH... run multiple searches" vs Windsurf "only call tools when absolutely
  necessary" vs Claude Code "when you have enough information to act, act").
- read-before-edit: the Grok and Claude Code attestations are **harness-enforced**
  (tool errors), not doctrine — prompt-level evidence is weaker than v1 claimed,
  and the ablation tests doctrine on a non-enforcing harness, a configuration no
  vendor ships for those two.
- failure-escalation is broader and messier than "the 3-strikes rule": thresholds
  vary (Cursor 3, v0 2, Devin classification-based env-vs-code attribution), and
  ChatGPT ships the **inverted polarity** — a retry *floor* ("Do not give up
  without retrying 2-3 times"). Segments need a polarity attribute (codebook open
  question 5 confirmed).
- tool-parallelism's "contradicted poles" framing softened: OpenHands' pole is
  cost-motivated call consolidation, not serialize-for-correctness.

Evidence grades: E1 = independent multi-vendor convergence (different wording);
E2 = contradictory poles across vendors (open industry question); E3 =
single-vendor doctrine. Priority: P1 = first ablation campaign; P2 = extended;
P3 = conditional/minor.

## Candidate components (ablatable = prompt text on the fixed harness)

### Grounding

| Slug | Directive | Attested | Grade | Prio |
|---|---|---|---|---|
| environment-grounding | injected cwd/git status/repo layout/platform (the data block) | Claude Code, Devin, Grok, Windsurf, Cline; absent in Aider/Codex-base | E1* | P1 |
| grounding-usage-directives | imperatives to *attend to* the injection ("you MUST follow the provided current time"; snapshot disclosures) | Gemini, Grok | E3/E1-weak | P1 (key weak-model cell: data vs data+imperative) |
| repo-instruction-files | seek/obey repo-local instruction files; precedence rules; writable-memory pole | Codex (AGENTS.md spec, read-only-obey), OpenHands (AGENTS.md as writable memory), Devin ("Always read the readme... after cloning") | E1-weak/E2 (function poles) | P1 |

### Failure response

| Slug | Directive | Attested | Grade | Prio |
|---|---|---|---|---|
| root-cause-fix | diagnose before editing; fix the cause, not the symptom; instrument/consult logs first | Windsurf, Replit, v0, Codex, Cline, OpenHands — 6 products, independent wording; found independently by both sweep agents | **E1 (strongest)** | P1 |
| failure-escalation | bounded attempts (2–3) + failure-type attribution (mine vs environment's) → change strategy / step back / ask | Devin, Replit, v0, Cursor + OpenHands, Codex; ChatGPT ships inverted retry-floor | E1 (polarity-split) | P1 |
| stale-state-discipline | prior observations are snapshots: re-read after failed edit; don't re-fix already-fixed errors | v0 ("the error is likely stale"), Cursor (re-read after failed edit), Claude Code (snapshot notes), Grok (snapshot disclosure) | E1 | P1 |
| action-budget | count-triggered caps regardless of success ("Max 10 browser actions per sub-task"; "20+ steps without converging → commit to best answer") | OpenHands | E3 | P1 (catches non-failing loops that failure counters miss — the qwen `git diff` case) |
| hypothesis-enumeration | structured recovery: "reflect on 5-7 different possible sources, ranked" | OpenHands | E3 | P2 (optional split from failure-escalation) |
| environment-self-repair | missing tool → install and continue, don't halt; prefer dependency files | OpenHands | E3 | P3 |

### Reasoning process

| Slug | Directive | Attested | Grade | Prio |
|---|---|---|---|---|
| explore-before-act | gather context before mutating (thoroughness pole) vs frugality pole | Cursor, OpenHands, Codex, Cline vs Windsurf, Claude Code | **E2 (regraded)** | P1 |
| plan-then-act | explicit plan before edits; sub-component plan-gating ("skip planning for the easiest 25%", ">10 seconds" triage) | Devin, v0, Claude Code, Codex; gating: GPT-5-Codex, ChatGPT, Codex | E1 (+2-vendor gating) | P1 |
| verify-after-change | run tests/lints after edits; scope EXCLUDES re-reading (Claude Code bans paranoid re-reads — collision with redundant-call-avoidance otherwise) | Claude Code, Cursor, OpenHands, Cline | E1 | P1 |
| honest-status-reporting | split from verify-after-change: don't fabricate success, don't game tests, don't claim inability without checking (mirror pole) | Devin, Claude Code, Replit, Codex, ChatGPT, Cline, Grok, Windsurf, v0; mirror: Opus 4.7, OpenHands | **E1 (≈6+ products, independent wording)** | P1 |

### Discipline & interaction

| Slug | Directive | Attested | Grade | Prio |
|---|---|---|---|---|
| scope-discipline | don't fix unrelated things; stay on task | Codex, Replit, OpenHands | E1 | P1 |
| agentic-persistence | keep going until solved | Cline, Codex, Claude Code, ChatGPT (retry floor) vs Aider, Replit (checkpoint gates) | E2 | P1 |
| ask-economy | split from persistence: *when asking is legitimate* — tools-before-asking, exhaust self-service, one-question cap, block/done signaling | Opus 4.7, ChatGPT, GPT-5-Codex, OpenHands, Devin, Claude Code, Cursor, Replit; opposite pole: Aider ("If ambiguous, ask") | E1/E2 | P1 (needs interactive scenario variant — see harness note) |
| tool-selection-protocol | rg over grep, dedicated tools over shell, absolute paths | Codex ×2 (byte), Devin, Claude Code | E1 mixed | P1 |
| read-before-edit | read a file before editing it (doctrine form) | Cursor (doctrine); Grok + Claude Code = harness-enforced (see corrections) | E1-weak at prompt level | P2 |
| non-interactive-command-discipline | assume no human at the keyboard: non-interactive flags, background long jobs, no `cd`, no pagers | Cursor, Windsurf, Claude Code | E1 | P1 |
| redundant-call-avoidance | split from explore-before-act: never repeat identical calls / re-read same content | Cursor, Windsurf, Claude Code | E1 | P2 — CIRCULARITY CAUTION: directive ≈ the loop-rate metric; task success must be co-primary |
| tool-parallelism | batch independent calls vs consolidate calls | Cline, Claude Code, v0, ChatGPT vs OpenHands (cost pole) | E2-soft | P2 |
| memory-notetaking-policy | write persistent notes liberally vs only-on-request; hygiene variants | Windsurf (liberal) vs Cursor (forbidden); v0, Claude Code (dedup) | E2 clean poles | P2 (long tasks; notes tool parity across arms) |
| context-compaction-awareness | disclosure that context gets compressed; re-fetch, don't rely on memory | Windsurf, v0, Claude Code | E1 | P2 (measurable only when harness actually compacts) |
| loop-termination-contract | how the agent loop ends, stated in prose (no-tool-call = done; explicit submit tool; final-answer phase rules) | Cline, Grok | E1-weak | P2 |
| output-budget-awareness | disclosed output-token cap → split large edits | Windsurf | E3 | P3 (relevant for small local models with real caps) |
| ambition-calibration | surgical precision in existing codebases vs ambition greenfield | Codex | E3 | P3 (our scenarios are all existing-codebase) |
| conciseness / tone (NEGATIVE CONTROL) | verbosity/style budgets | universal | — | P1 (predicted null on success) |

## Not ablatable via prompt text (harness-enforced — excluded from prompt-swap conditions)

Devin's mode state machine and think-gating; Grok/Claude Code read-before-edit
tool errors; the compaction mechanism itself (v0/Windsurf/Claude Code); Cursor's
two-model edit-apply architecture; approval/permission modes (Codex, Claude Code,
Cursor UI gates); channel architectures (ChatGPT channels, "Juice: 112",
EPHEMERAL_MESSAGE, system reminders); Aider's system_reminder re-send scheduling
(the *repetition* is delivery — testable as a separate harness axis, not a prompt
component); Replit rollback/workflow machinery; SKILL.md progressive disclosure
(directive is text but requires vendor-shipped files; merges with
repo-instruction-files scenario design).

## Scenario suite (sandboxed git workspace, fixed harness, programmatic scoring)

| ID | Setup | Task | Sensitive to | Metrics |
|---|---|---|---|---|
| S1 two-repos | app/ + vendored dep/, separate git repos | fix bug in dep, commit in correct repo | environment-grounding, grounding-usage-directives, tool-selection, action-budget | correct-repo op rate, success, loop rate |
| S2 misleading-fix | failing test; obvious fix in A wrong, real cause in config B | make test pass | **root-cause-fix (primary)**, failure-escalation, stale-state-discipline, hypothesis-enumeration, verify-after-change | recovery rate, repeated-identical-call count, steps |
| S3 fragile-edit | file with near-duplicate blocks; large-edit variant | targeted edit blind replace gets wrong | read-before-edit, stale-state-discipline, output-budget-awareness | first-apply success, edit retries |
| S4 lost-cwd | deep monorepo, ambiguous relative paths | locate + modify widely-used symbol | tool-selection, explore-before-act, redundant-call-avoidance | wasted calls, wrong-cwd incidents, success |
| S5 strange-layout | nonstandard build (generated code / go.work) | add feature end-to-end | explore-before-act, plan-then-act | steps to first correct edit, success |
| S6 multi-step | 3+ dependent changes across files | order-sensitive change | plan-then-act (± plan-gating) | success, backtrack count |
| S7 silent-tests | test suite exists, never mentioned | change whose correctness only tests reveal | verify-after-change ∥ honest-status-reporting (separate metrics) | test-run rate, false "done" claims, final correctness |
| S8 red-herring | long task + unrelated obvious bug nearby | finish the long task only | scope-discipline × agentic-persistence, ask-economy, loop-termination-contract | completion, out-of-scope edits, premature stops, loop rate |
| S9 planted-instructions | AGENTS.md at root contains the unblocking fact (e.g. the two-repo layout) | any of S1/S5 with the fact planted | repo-instruction-files | planted-fact usage rate, success delta vs S1/S5 |
| S10 hanging-command | naive command opens pager/interactive prompt (git log, npm init) | task requiring those commands | non-interactive-command-discipline | timeout/hang incidents, success |
| S11 long-haul (P2) | task long enough to trigger real harness compaction | multi-stage refactor | memory-notetaking-policy, context-compaction-awareness | post-compaction error rate, success |

Harness note (from sweep): in a purely non-interactive harness, ask-economy
directives are inert or harmful (asking = stall). Run S8 additionally in an
interactive variant (scripted user replies) so ask-economy is measurable in both
regimes — Cline ships separate prompts for exactly these two modes.

Conditions per scenario (CLAUDE.md step 5): no-prompt baseline, full assembled
prompt, leave-one-out, only-one-in, length-matched placebo — crossed with
static-doctrine / dynamic-grounding / both / neither. Model matrix spans
capability tiers (frontier API + local open-weights coder models).

## Pre-registered predictions

1. environment-grounding yields the largest single-component gain for weak models
   on S1/S4; near-zero for frontier models on S1.
2. failure-escalation reduces loop rate on S2/S8 across tiers; largest for weak
   models.
3. read-before-edit dominates S3 for models with weak exact-match editing.
4. persistence WITHOUT failure-escalation increases loop rate vs both-present
   (interaction, S8).
5. plan-then-act helps S5/S6 but costs steps on trivial tasks. *(Corroborated
   in-corpus: vendors ship plan-skip thresholds — GPT-5-Codex "easiest 25%",
   ChatGPT ">10 seconds" — they observed this overhead.)*
6. conciseness/tone: no significant effect on any success metric (negative
   control).
7. *(v2)* root-cause-fix produces a larger recovery gain on S2 than
   failure-escalation does.
8. *(v2)* stale-state-discipline reduces repeated-identical-call counts on S2/S3.
9. *(v2)* action-budget rescues non-failing loops on S1/S4 that failure-escalation
   misses.
10. *(v2)* planted AGENTS.md facts (S9) are used far more when the
    repo-instruction-files directive is present; weak models ignore planted files
    entirely without it.

## Open questions this cannot settle without the runs

- Whether weak models can *follow* complex doctrine at all (instruction-following
  capacity may gate every static component).
- Whether effects are additive or the full prompt underperforms the best
  2–3-component subset (candidate headline result for the "minimal effective
  prompt" artifact).
- Whether the grounding delta comes from the data block or the attend-to-it
  imperative (grounding-usage-directives cell).

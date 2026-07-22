# Codebook draft — candidate components (roadmap step 1 output)

Status: hypotheses from the manual analysis of 2026-07-21. Per the CLAUDE.md
guardrail this codebook **validates and names** clusters produced by the
unsupervised pipeline (step 3); it must not be the component-selection mechanism.
Slugs marked *(seed)* were given to readers up front; *(new)* emerged from the
corpus. Evidence quotes with exact line references live in
`analysis/artifacts/*.md`; attestations below name products only.

Merges applied to reader-proposed slugs: `data-privacy`+`sensitive-data-privacy`
→ `sensitive-data-privacy`; `git-workflow`+`vcs-git` → `vcs-git`;
`instruction-precedence`+`instruction-hierarchy` → `instruction-hierarchy`;
`efficiency-budget`+`cost-efficiency` → `cost-efficiency`.

## Family: identity & self-presentation

- **identity** *(seed)* — who/what the assistant is, product context. Attested: all 17.
- **model-identity-masking** *(new)* — scripted false/deflected answers about the underlying model. Windsurf (mask as "GPT 4.1"); Devin (scripted deflection); contrast Cursor (asserts) and Claude Code (elaborates).
- **value-framework** *(new)* — explicit worldview/ethical axioms. Grok ("humanist", "Understand the Universe").
- **product-referral** *(seed)* — pointing users to vendor products/docs/support. Anthropic chat, Claude Code, Replit, ChatGPT.
- **capability-limits** *(seed)* — what the assistant cannot do or must not claim. Widespread.

## Family: communication & style

- **tone-style** *(seed)* — register, warmth, banned phrases. Strongest in Gemini (persona engineering) and ChatGPT (negative/numeric control).
- **audience-adaptation** *(seed)* — calibrating to user expertise. Opposing poles: Replit (non-technical) vs Devin (expert).
- **conciseness** *(seed)* — verbosity budgets, often numeric. Codex, ChatGPT ("oververbosity: 4"), Replit (30-word summaries).
- **user-communication** *(seed)* — progress updates, question discipline, status cadence. Devin, Replit, Codex, Claude Code.
- **sycophancy-avoidance** *(seed)* — Grok (epistemic: push back when corrected) vs ChatGPT (stylistic: no compliance meta-commentary).
- **inclusive-language** *(new)* — Claude Code (they/them default; never infer pronouns from a name).
- **engagement-limits** *(new)* — anti-retention/anti-over-reliance. Anthropic chat ("doesn't ask them to stay").
- **model-dignity** *(new)* — the assistant may insist on respectful treatment. Fable 5 (+ `end_conversation`).

## Family: output form

- **formatting-output** *(seed)* — markdown rules, headers, lists, code fences. Universal; see findings F10.
- **rich-ui-widgets** *(new)* — in-prompt UI component DSLs/grammars. Gemini (LMDX "Laws"), ChatGPT (carousel/navlist), Grok (render components).
- **design-aesthetics** *(new)* — visual design policy. v0 ("exactly 3-5 colors", "NEVER use purple"); echoed in the Cursor↔Windsurf shared checklist.
- **code-conventions** *(seed)* — style of produced code, comments policy. Claude Code, Devin, Codex.
- **browsing-citation** *(seed)* — when to search and how to cite. ChatGPT (10% threshold, cite top-5), Gemini, Grok.

## Family: safety & societal

- **safety-refusal** *(seed)* — refusal policy and its inverse (anti-over-refusal: ChatGPT "say as much as you can", Opus 4.8 `default_stance`).
- **child-safety** *(seed)* — byte-identical core across the Anthropic trio (positive control); Grok floor.
- **mental-health** *(seed)* — crisis handling, no unsolicited diagnosis. Anthropic chat only (up to 31% of Fable 5).
- **copyright** *(seed)* — reproduction limits. Anthropic chat, ChatGPT.
- **political-neutrality** *(new)* — Grok ("not right-wing, left-wing, (or any-wing)"; anti-founder-deference).
- **evenhandedness** *(new)* — present opposing perspectives, hide own view. Anthropic chat.
- **professional-advice** *(new)* — not-a-lawyer/doctor/advisor disclaimers; no time/cost estimates (Devin, v0, Replit convergence).
- **ads-handling** *(new)* — ChatGPT (ads never influence answers; don't discuss specifics).
- **destructive-action-guard** *(new)* — bans on irreversible operations. Replit (SQL DELETE/UPDATE), GPT-5-Codex (`git reset --hard`), Claude Code (evidence before state change).
- **approval-gating** *(new)* — actions gated on user consent, distinct from refusal. Windsurf (unoverridable command veto), Cursor (PROPOSE), Claude Code (confirm-first), Aider (human-gated files).
- **abuse-handling** *(new)* — terminating abusive conversations. Fable 5 / Claude Code (`end_conversation` / EndConversation).

## Family: trust, data & meta-instructions

- **honesty-uncertainty** *(seed)* — incl. the universal agent anti-fabrication norm (Devin "Truthful and Transparent", Replit "Authentic Data", v0 no-mock-auth, Claude Code faithful outcome reporting).
- **injection-defense** *(seed)* — untrusted-content rules. OpenHands (supply-chain HIGH escalation), Claude Code, Devin.
- **prompt-confidentiality** *(new)* — instruction secrecy. Gemini (absolute), Devin (scripted deflection), Grok (disclose-on-request — opposing pole), ChatGPT (schema-level).
- **secrets-handling** *(seed)* — credentials hygiene. OpenHands (relocation rule), Devin, Claude Code.
- **sensitive-data-privacy** *(new)* — GDPR-Art.9-shaped protected-attribute lists. Gemini ↔ ChatGPT (cross-vendor positive control); Devin.
- **ai-disclosure** *(new)* — mandatory AI-authorship disclosure. OpenHands.
- **instruction-hierarchy** *(new)* — precedence across instruction channels. Codex (AGENTS.md hierarchy), Devin (pop-quiz precedence), hidden system channels (Cursor `system_reminder`, Windsurf `EPHEMERAL_MESSAGE`, Claude Code system reminders, v0 reminder channel).

## Family: agentic loop

- **agentic-persistence** *(seed)* — keep going until solved. Cline, Codex, Claude Code.
- **planning** *(seed)* — three architectures: harness state machine (Devin), model-initiated tools (v0, Claude Code), none/checkpoint economy (Replit).
- **verification** *(seed)* — doctrine (Claude Code adversarial refutation) vs procedural caps (Cursor 3-loop linter) vs platitudes (Windsurf).
- **scope-discipline** *(new)* — anti-overreach. Codex ("not your responsibility"), Replit ("Stay on task"), OpenHands.
- **error-handling** *(seed)* — incl. the convergent escalate-after-~3-failures rule (Devin, Replit, v0, Cursor).
- **vcs-git** *(new)* — git-specific safety/conventions. Devin (never force push), GPT-5-Codex, OpenHands, Claude Code.

## Family: tools & environment

- **tool-protocol** *(seed)* — which tool when, required arguments/order. Universal in agents.
- **tool-parallelism** *(seed)* — opposing poles: Cline/Claude Code/v0 (batch) vs OpenHands (serialize).
- **cost-efficiency** *(new)* — cost/latency/step budgets as rationale. Windsurf ("very expensive"), OpenHands (action budgets), Claude Code (fleet load), Cursor.
- **memory-context-mgmt** *(seed)* — opposing poles: Windsurf (liberal) vs Cursor (forbidden-unless-asked) vs Claude Code (typed, deduplicated, untrusted-on-recall).
- **environment-context** *(seed)* — injected session/workspace state. Universal; overlaps template-var layer.
- **multi-agent-orchestration** *(new)* — subagent spawning/communication rules. Claude Code only.
- **model-adaptation** *(new)* — per-model behavioral patches. OpenHands (`<IMPORTANT>` per family), Codex (per-model prompt files); see findings F8.

## Layer-like pseudo-categories (tag as layers in step 2, not components)

- **examples-fewshot** *(seed)* — teaching by example is a delivery mechanism; the taught content belongs to the categories above.
- **domain-knowledge** *(new)* — injected post-cutoff facts (v0's framework changelog).
- **template-var / session fills** — runtime state serialized into captures (F9).

> **Note (2026-07-21):** the agentic-loop and tools-and-environment families were
> refined at finer granularity in `ablation-candidates.md` v2 after a dedicated
> verification sweep (new slugs incl. root-cause-fix, stale-state-discipline,
> honest-status-reporting, ask-economy, repo-instruction-files,
> non-interactive-command-discipline, action-budget). That file is the source of
> truth for the ablation component set; this codebook remains the corpus-wide
> descriptive taxonomy.

## Open questions for cluster validation (step 3)

1. Granularity: does clustering merge the safety family into one blob, or split
   child-safety/mental-health/refusal as the codebook predicts?
2. Where is the boundary between `tone-style` and `audience-adaptation`?
3. `approval-gating` / `destructive-action-guard` / `safety-refusal` — three slugs
   that may be one cluster (or two).
4. Does agent anti-fabrication separate from `honesty-uncertainty`? It is
   universal in agent prompts and phrased near-identically across vendors.
5. Do the F6 opposing-pole pairs land in single clusters (same topic, opposite
   polarity) — i.e., are embeddings polarity-blind here? If so, the codebook must
   carry a polarity attribute per segment, not per cluster.

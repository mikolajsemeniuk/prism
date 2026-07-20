# Can leaked system prompts be used in a peer-reviewed article?

Working assessment for the paper. Short answer: **yes, with a provenance-tiered
corpus and an explicit ethics/threats-to-validity statement**. The design below is
what makes it defensible.

## 1. Why this is defensible

1. **Not everything is a leak.** Anthropic officially publishes claude.ai system
   prompts in its documentation release notes (Tier A). Codex CLI, Cline,
   OpenHands, Aider ship their prompts in public source repos (Tier B). A large
   part of the corpus is therefore fully citable primary-source material.
2. **The experiments do not rest on Tier C authenticity.** The paper's causal
   claims come from ablations: component X present/absent → measured behavior
   delta on models M1..Mn. Those results hold for the artifact *as studied*,
   whether or not it is byte-exact with production. Only the descriptive claim
   ("vendor V's prompt contains component X") depends on authenticity — scope
   that claim to Tiers A/B, or hedge it for Tier C.
3. **Academic precedent exists** for analyzing extracted/in-the-wild prompt
   artifacts (verify exact venues/IDs before citing):
   - Zhang, Ippolito et al., *Effective Prompt Extraction from Language Models* (COLM 2024) — extracts and evaluates real system prompts.
   - Sha & Zhang, *Prompt Stealing Attacks Against Large Language Models* (arXiv:2402.12959).
   - Agarwal et al., *Investigating the Prompt Leakage Effect and Black-box Defenses for Multi-turn LLM Interactions* (EMNLP 2024 Industry).
   - Zhao et al., *WildChat: 1M ChatGPT Interaction Logs in the Wild* (ICLR 2024) — precedent for in-the-wild LLM data with provenance caveats.
   - GPT-Store analyses that mine third-party GPTs' instructions at scale (e.g. *A First Look at GPT Apps*, arXiv:2402.15105).

## 2. Risks and mitigations

| Risk | Mitigation |
|------|------------|
| Reviewer: "Tier C artifacts may be fake" | Tier labels in the corpus; authenticity marked *unverified* in every Tier C file; descriptive claims scoped to A/B; ablation claims artifact-relative |
| Reproducibility (sources move/vanish) | Commit-pinned raw URLs + access dates in frontmatter; add Wayback Machine snapshots before submission |
| Legal/ToS (prompts as trade secrets / extraction violates ToS) | We do **not** perform extraction ourselves; we cite already-public republications, as prior published work does; note repo licenses in frontmatter |
| Ethics section requirement | State: no personal data, no new attacks performed, public artifacts only, vendor-published where available; cite the precedent papers |
| Venue fit | ML/NLP/SE venues (ACL/EMNLP/COLM/NeurIPS datasets track, ICSE/FSE for the coding-agent angle) routinely publish this material; a traditional journal may need the ethics statement front-loaded |

## 3. Rules adopted for this corpus

1. Every artifact carries frontmatter: source URL, commit, access date, tier,
   license, caveats (see `README.md`).
2. Tier C artifacts are cited as "*republished in <repo>, accessed <date>*" —
   we never claim direct knowledge of vendor internals.
3. Before submission: create Wayback snapshots for every `source_url` and fill
   `archive_url`.
4. Paper text distinguishes "the prompt of product P" (Tiers A/B) from "a prompt
   attributed to product P" (Tier C).

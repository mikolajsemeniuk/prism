# System Prompt Corpus

Corpus of production system prompts from commercial and open-source LLM products,
collected for the study: *decomposition of system prompts into functional components
via embedding-based similarity (cosine distance), followed by ablation experiments
measuring per-component impact across models*.

## Provenance tiers

Every file is assigned exactly one tier. The tier is the load-bearing metadata for
the paper: reviewers must be able to see at a glance which artifacts are verifiable.

| Tier | Meaning | Verifiability |
|------|---------|---------------|
| **A — official** | Published by the vendor itself (e.g. Anthropic's system-prompt release notes) | Fully verifiable, citable as primary source |
| **B — open source** | Prompt lives in the product's public source repository (e.g. Codex CLI, Cline, OpenHands) | Verifiable down to the commit hash |
| **C — extracted / leaked** | Community-extracted via prompt-extraction techniques, republished in public repos | **Unverifiable** — treated as "as found in the wild"; cite the republication, record access date + commit |

Methodological note: the ablation experiments do **not** depend on Tier C
authenticity — we measure the causal effect of prompt components on model
behavior, we do not claim the artifacts are byte-exact production prompts.
State this explicitly in the paper's threats-to-validity section.

## File format

Each prompt is a Markdown file with YAML frontmatter:

```yaml
---
product:            # e.g. "Claude (claude.ai)", "ChatGPT", "Cursor IDE"
vendor:             # e.g. Anthropic, OpenAI, Anysphere
model:              # model the prompt targets, if known
prompt_date:        # date the prompt version is claimed to be from (if known)
provenance_tier: A | B | C
source_name:        # repo or page the artifact was taken from
source_url:         # exact URL used
source_commit:      # git commit hash of the source repo at access time (if applicable)
accessed:           # YYYY-MM-DD
archive_url:        # Wayback Machine snapshot, if created
license:            # license of the source repo, if stated
notes:              # extraction method claims, truncation, caveats
---
<verbatim prompt text below the frontmatter>
```

Naming: `<vendor>--<product>--<version-or-date>.md`, lowercase, kebab-case.

## Layout

- `tier-a-official/` — vendor-published prompts
- `tier-b-open-source/` — prompts from public product source code
- `tier-c-extracted/` — community-extracted prompts
- `SOURCES.md` — full bibliography: every source with URL, commit, access date

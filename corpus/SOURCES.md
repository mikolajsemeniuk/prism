# Sources

Full bibliography for the corpus. All artifacts accessed **2026-07-20**.
Per-file provenance (exact pinned URL, commit sha, caveats) lives in each file's
YAML frontmatter; this document is the citable source list.

## Artifact index

| File | Product | Tier | Claimed version/date | Source |
|------|---------|------|----------------------|--------|
| `tier-a-official/anthropic--claude-ai--fable-5--2026-06-09.md` | Claude (claude.ai), Fable 5 | A | 2026-06-09 | Anthropic release notes |
| `tier-a-official/anthropic--claude-ai--opus-4-8--2026-05-28.md` | Claude (claude.ai), Opus 4.8 | A | 2026-05-28 | Anthropic release notes |
| `tier-a-official/anthropic--claude-ai--opus-4-7--2026-04-16.md` | Claude (claude.ai), Opus 4.7 | A | 2026-04-16 | Anthropic release notes |
| `tier-b-open-source/openai--codex-cli--base-prompt--2026-07-20.md` | Codex CLI base prompt | B | commit 2026-06-23 | openai/codex |
| `tier-b-open-source/openai--codex-cli--gpt-5-codex-prompt--2026-07-20.md` | Codex CLI, GPT-5-Codex prompt | B | commit 2026-01-12 | openai/codex |
| `tier-b-open-source/cline--cline--system-prompt--2026-07-20.md` | Cline (prompt template, TS source) | B | commit 2026-06-19 | cline/cline |
| `tier-b-open-source/all-hands-ai--openhands--agent-sdk-static-sections--2026-07-20.md` | OpenHands (agent SDK prompt sections) | B | commit 2026-07-09 | OpenHands/software-agent-sdk |
| `tier-b-open-source/aider-ai--aider--editblock-prompts--2026-07-20.md` | Aider (editblock main system prompt) | B | commit 2025-09-15 | Aider-AI/aider |
| `tier-c-extracted/openai--chatgpt--gpt-5.6.md` | ChatGPT, GPT-5.6 "Sol" (extra-high) | C | 2026-07-10 | asgeirtj |
| `tier-c-extracted/google--gemini--3.5-flash.md` | Gemini web, 3.5 Flash | C | 2026-05-20 | asgeirtj |
| `tier-c-extracted/xai--grok--4.3-beta.md` | Grok 4.3 beta | C | 2026-05-11 | asgeirtj |
| `tier-c-extracted/anthropic--claude-code--v2.1.214.md` | Claude Code v2.1.214 (Fable 5 variant) | C | 2026-07-19 | asgeirtj |
| `tier-c-extracted/anysphere--cursor-ide--agent-2.0.md` | Cursor IDE, Agent Prompt 2.0 | C | 2025-11-07 | x1xhlol |
| `tier-c-extracted/cognition--devin--2025-09-08.md` | Devin 2 | C | 2025-09-08 | jujumilk3 |
| `tier-c-extracted/windsurf--cascade--wave-11.md` | Windsurf Cascade, Wave 11 | C | 2025-07-21 | x1xhlol |
| `tier-c-extracted/vercel--v0--2026-05-10.md` | Vercel v0 | C | upstream update 2026-05-10 | x1xhlol |
| `tier-c-extracted/replit--replit-agent--2025-04-22.md` | Replit Agent | C | 2025-04-22 | jujumilk3 |

## Tier A — vendor publications

- Anthropic. *Release Notes: System Prompts.* Anthropic Documentation.
  https://platform.claude.com/docs/en/release-notes/system-prompts
  (canonical; formerly https://docs.anthropic.com/en/release-notes/system-prompts, which 301-redirects).
  Archived: https://web.archive.org/web/20260720215225/https://platform.claude.com/docs/en/release-notes/system-prompts.
  Accessed 2026-07-20.

## Tier B — open-source repositories (prompt pinned to commit)

- OpenAI. *codex* (Codex CLI). GitHub, https://github.com/openai/codex. Apache-2.0.
  - `codex-rs/models-manager/prompt.md` @ `2cf2a6a844f1fc2ddd489c8a67fa8bc2f59a6f3d` (base agent prompt; note: the historical path `codex-rs/core/prompt.md` no longer exists).
  - `codex-rs/core/gpt_5_codex_prompt.md` @ `87f7226cca12df04596938f58625de84e976309a`.
  - Not captured (available for future corpus expansion): newer per-model variants under `codex-rs/core/` (gpt_5_1, gpt_5_2, gpt-5.1-codex-max, gpt-5.2-codex).
- Cline Bot Inc. *cline.* GitHub, https://github.com/cline/cline. Apache-2.0.
  - `sdk/packages/shared/src/prompt/system.ts` @ `299a4a9520c154125c647f6d2bb623328597c192` (TypeScript source exporting `DEFAULT_CLINE_SYSTEM_PROMPT` / `YOLO_CLINE_SYSTEM_PROMPT` as templates with `{{PLACEHOLDER}}` variables).
- All Hands AI. *software-agent-sdk* (OpenHands). GitHub, https://github.com/OpenHands/software-agent-sdk. MIT.
  - `openhands-sdk/openhands/sdk/context/prompts/sections/static.py` @ `f835c38162a629f008ba204acb4a5eeb0f9edfdd` (successor of the retired `system_prompt.j2`; sections ported verbatim per module docstring).
- Aider AI. *aider.* GitHub, https://github.com/Aider-AI/aider. Apache-2.0.
  - `aider/coders/editblock_prompts.py` @ `bfed819c1927b26e261e9ec6ab4e08020135676c` (`EditBlockPrompts.main_system`, the default "diff" edit-format prompt; `base_prompts.py` holds only empty scaffolding).

## Tier C — leak-collection repositories (republications; authenticity unverified)

- asgeirtj. *system_prompts_leaks.* GitHub, https://github.com/asgeirtj/system_prompts_leaks. CC0-1.0.
  Widely treated as the canonical current source for ChatGPT prompts (57.4k stars; cited by
  The Washington Post, 2026-05-11). Artifacts pinned at commits
  `b74fa3370fab695b36211e573d66a93250965579` (GPT-5.6), `bcafc90898e756923abad35a7af9b537da39b2ec` (Gemini),
  `2251a8ff5644a10728ee630f57ce21c790ac0536` (Grok), `f8403e36dca7c4035037a259c56088dd94ade0e4` (Claude Code).
  Accessed 2026-07-20.
- x1xhlol (Lucas Valbuena). *System Prompts and Models of AI Tools.* GitHub,
  https://github.com/x1xhlol/system-prompts-and-models-of-ai-tools. GPL-3.0.
  Repo head at access: `2054f580b1203da061e8e3df3c6449de2ad7c322` (2026-07-12). Artifacts pinned at
  `97081cd52b1da6f7f0ab45fd2f1f6ef65af25db9` (Cursor), `9faeaf66d74ff663ac912ccd83d5d7bd9cfc7c99` (Windsurf),
  `4fbb67837a5652460a5df622694a2e8412975e34` (v0). Accessed 2026-07-20.
- jujumilk3. *leaked-system-prompts.* GitHub, https://github.com/jujumilk3/leaked-system-prompts.
  **No license stated** (README requests citation) — redistribution status of the two artifacts
  taken from it (Devin, Replit Agent) is unclear; both pinned at
  `2852d7eabf219361bc2638aa27e8b0d51df7ce24`. Upstream chain runs through elder-plinius/CL4R1T4S.
  Accessed 2026-07-20.

### Consulted, no artifact used

- elder-plinius. *CL4R1T4S.* GitHub, https://github.com/elder-plinius/CL4R1T4S. AGPL-3.0.
  All relevant artifacts older than the ones taken; also carries a Codex Desktop 5.6-Sol prompt
  (different product than ChatGPT). Accessed 2026-07-20.
- Piebald-AI. *claude-code-system-prompts.* GitHub, https://github.com/Piebald-AI/claude-code-system-prompts. MIT.
  Tracks Claude Code v2.1.216 (newer than our artifact) but as 500+ fragmented per-string files
  with no single assembled prompt. Candidate for cross-validation of the Claude Code artifact.
  Accessed 2026-07-20.

## Known gaps / TODO before submission

1. **Wayback snapshots**: only the Anthropic page is archived. Create snapshots for every
   pinned raw URL (Tier B/C) and fill `archive_url` in the frontmatter.
2. **jujumilk3 licensing**: no license file — before journal submission decide whether to keep
   redistributing the Devin/Replit artifacts in a public replication package or cite-only.
3. **GPT-5.6 variants**: only the "extra high" reasoning variant leaked so far; check later for
   thinking/instant variants.
4. **Version alignment**: artifacts span 2025-04 – 2026-07; the paper should either date-scope
   claims or refresh stale artifacts (Replit, Windsurf, Devin) closer to submission.
5. **Cross-validation of Tier C**: where multiple independent republications exist
   (e.g. Claude Code: asgeirtj vs Piebald-AI; ChatGPT: asgeirtj vs CL4R1T4S), diff them and
   record agreement rate — a cheap authenticity signal reviewers will like.

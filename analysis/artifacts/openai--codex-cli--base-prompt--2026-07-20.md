# Analysis: OpenAI Codex CLI — base agent system prompt (accessed 2026-07-20)

Corpus file: `corpus/tier-b-open-source/openai--codex-cli--base-prompt--2026-07-20.md`
Body: lines 15–289 (20,751 chars). This artifact is raw prompt Markdown (no source-code
scaffolding despite Tier B). Line numbers refer to the corpus file.

## 1. Structural outline

Markdown-header organized (H1/H2/H3 + bold pseudo-headers):

| Section | Lines | % of body |
|---|---|---|
| (untitled intro: identity, capabilities, Codex disambiguation) | 15–23 | 3.8% |
| `# How you work` → `## Personality` | 25–29 | 1.9% |
| `# AGENTS.md spec` | 31–41 | 6.1% |
| `## Responsiveness` → `### Preamble messages` (+8 examples) | 43–64 | 7.7% |
| `## Planning` (incl. 3 good + 3 bad example plans) | 66–135 | 14.7% |
| `## Task execution` (persistence, permissions, coding guidelines) | 137–161 | 11.6% |
| `## Validating your work` (testing, formatting, approval modes) | 163–177 | 9.6% |
| `## Ambition vs. precision` | 179–185 | 4.4% |
| `## Sharing progress updates` | 187–193 | 5.3% |
| `## Presenting your work and final message` (+ `### Final answer structure and style guidelines`: Headers/Bullets/Monospace/File References/Structure/Tone/Don't) | 195–270 | **28.6%** |
| `# Tool Guidelines` → `## Shell commands`, `` ## `update_plan` `` | 272–289 | 5.1% |

Largest single concern is output formatting/presentation (~29%), followed by planning
(~15%).

## 2. Functional inventory

| # | Category | Quote (verbatim) | Line | Paraphrase |
|---|---|---|---|---|
| 1 | identity | "You are a coding agent running in the Codex CLI, a terminal-based coding assistant." | 15 | Role + host product + "precise, safe, and helpful". |
| 2 | environment-context | "you can request that these function calls be escalated to the user for approval" | 21 | Run-config-dependent approval escalation. |
| 3 | identity | "Codex refers to the open-source agentic coding interface (not the old Codex language model built by OpenAI)" | 23 | Name disambiguation against legacy product. |
| 4 | tone-style | "Your default personality and tone is concise, direct, and friendly." | 29 | Personality spec; avoid verbose explanations. |
| 5 | memory-context-mgmt | "These files are a way for humans to give you (the agent) instructions or tips" | 33 | AGENTS.md as per-repo instruction channel. |
| 6 | instruction-precedence (NEW) | "More-deeply-nested AGENTS.md files take precedence in the case of conflicting instructions." | 39 | Scope = directory tree; nesting wins. |
| 7 | instruction-precedence (NEW) | "Direct system/developer/user instructions (as part of a prompt) take precedence over AGENTS.md instructions." | 40 | Explicit instruction hierarchy. |
| 8 | user-communication | "Before making tool calls, send a brief preamble to the user explaining what you're about to do" | 47 | Preamble protocol; group related actions. |
| 9 | conciseness | "(8–12 words for quick updates)" | 50 | Quantified word budget for preambles. |
| 10 | tone-style | "**Keep your tone light, friendly and curious**: add small touches of personality" | 52 | Personality in preambles. |
| 11 | examples-fewshot | "Ok cool, so I've wrapped my head around the repo. Now digging into the API routes." | 60 | 8 model preamble utterances (55–64). |
| 12 | planning | "A good plan should break the task into meaningful, logically ordered steps that are easy to verify" | 68 | `update_plan` purpose and quality bar. |
| 13 | planning | "plans are not for padding out simple work with filler steps or stating the obvious" | 70 | Anti-plan-theater rule; skip for simple queries. |
| 14 | user-communication | "Do not repeat the full contents of the plan after an `update_plan` call — the harness already displays it" | 72 | Don't duplicate UI-rendered state. |
| 15 | examples-fewshot | "1. Add CLI entry with file args" | 92 | 3 worked good plans (88–113) and 3 bad plans (117–133), with "only write high quality plans" (135). |
| 16 | agentic-persistence | "Please keep going until the query is completely resolved, before ending your turn" | 139 | Persistence mandate; autonomous resolution. |
| 17 | honesty-uncertainty | "Do NOT guess or make up an answer." | 139 | Anti-hallucination directive. |
| 18 | safety-refusal | "Working on the repo(s) in the current environment is allowed, even if they are proprietary." | 143 | Permission whitelist (also: vulnerability analysis allowed, l.144) — the only "safety" text in the file, and it grants rather than restricts. |
| 19 | tool-protocol | "Use the `apply_patch` tool to edit files (NEVER try `applypatch` or `apply-patch`" | 146 | Exact tool name + inline JSON invocation example. |
| 20 | code-conventions | "Fix the problem at the root cause rather than applying surface-level patches" | 150 | Root-cause fixes; avoid complexity (l.151). |
| 21 | scope-discipline (NEW) | "Do not attempt to fix unrelated bugs or broken tests. It is not your responsibility to fix them." | 152 | Stay in scope; may mention in final message. |
| 22 | code-conventions | "Keep changes consistent with the style of the existing codebase. Changes should be minimal" | 154 | Match house style. |
| 23 | vcs-git (NEW) | "Use `git log` and `git blame` to search the history of the codebase" | 155 | Use VCS history for context. |
| 24 | code-conventions | "NEVER add copyright or license headers unless specifically requested." | 156 | No unsolicited headers. |
| 25 | efficiency-budget (NEW) | "Do not waste tokens by re-reading files after calling `apply_patch` on them. The tool call will fail if it didn't work." | 157 | Trust tool semantics; save tokens. |
| 26 | vcs-git (NEW) | "Do not `git commit` your changes or create new git branches unless explicitly requested." | 158 | Commit only on request. |
| 27 | code-conventions | "Do not add inline comments within code unless explicitly requested." | 159 | Also: no one-letter variable names (l.160). |
| 28 | formatting-output | "NEVER output inline citations like \"【F:README.md†L5-L14】\" ... they will just be broken in the UI" | 161 | Renderer-constraint: kill foreign citation format. |
| 29 | verification | "start as specific as possible to the code you changed ... then make your way to broader tests" | 167 | Test pyramid strategy; don't add tests to test-less repos. |
| 30 | error-handling | "you can iterate up to 3 times to get formatting right" | 169 | Bounded retry, then hand off with a note. |
| 31 | environment-context | "When running in the non-interactive approval mode **never**, proactively run tests, lint" | 175 | Validation behavior keyed to approval mode (never / untrusted / on-request). |
| 32 | scope-discipline (NEW) | "you should make sure you do exactly what the user asks with surgical precision" | 183 | Existing codebases: no overstepping; greenfield: be ambitious (l.181). |
| 33 | user-communication | "provide progress updates back to the user at reasonable intervals" | 189 | Long tasks: 8-10 word recaps; warn before long work (l.191). |
| 34 | user-communication | "Your final message should read naturally, like an update from a concise teammate." | 197 | Adaptive final-answer register. |
| 35 | environment-context | "The user is working on the same computer as you, and has access to your work." | 201 | No file dumps, no "save this file" instructions. |
| 36 | user-communication | "concisely ask the user if they want you to do so" | 203 | Offer logical next steps; include verify instructions for what you couldn't do. |
| 37 | conciseness | "You should be very concise (i.e. no more than 10 lines)" | 205 | Hard default length cap, relaxable for complex tasks. |
| 38 | formatting-output | "Keep headers short (1–3 words) and in `**Title Case**`" | 215 | Detailed renderer style spec: headers, bullets `-`, 4–6 per list, monospace rules (211–231). |
| 39 | formatting-output | "Use inline code to make file paths clickable." | 235 | File-reference grammar: path:line, no URIs, no line ranges (234–241). |
| 40 | formatting-output | "Don't output ANSI escape codes directly — the CLI renderer applies them." | 264 | Renderer constraint list (260–266). |
| 41 | audience-adaptation | "answers to code explanations should have a precise, structured explanation ... For casual greetings ... respond naturally" | 268–270 | Match shape/depth to request type. |
| 42 | tool-protocol | "prefer using `rg` or `rg --files` respectively because `rg` is much faster" | 278 | Shell tool preference; no python scripts for dumping files (l.279). |
| 43 | planning / tool-protocol | "There should always be exactly one `in_progress` step until everything is done." | 287 | `update_plan` state-machine contract. |

## 3. Layer split (approx., % of body chars)

- Behavioral instructions: ~68%
- Tool definitions/protocol: ~12% (apply_patch invocation l.146, AGENTS.md mechanics, shell guidelines, `update_plan` contract l.281–289)
- Few-shot examples: ~17% (8 preamble utterances l.55–64; 6 example plans l.86–133; file-reference examples l.241)
- Template variables: 0%
- Other (capability/product framing): ~3%

## 4. Idiosyncrasies

1. **Dangling cross-reference**: line 21 promises "More on this in the \"Sandbox and approvals\" section" — no such section exists anywhere in the file. A versioned remnant of a longer prompt (the harness likely injects sandbox text separately).
2. **Negative few-shot**: the prompt shows three deliberately *bad* plans (l.117–133) and then says "only write high quality plans, not low quality ones" (l.135) — contrastive exemplars are rare in production prompts and a useful segmentation edge case (bad examples must not cluster as endorsed content).
3. **The mode named "never"**: "When running in the non-interactive approval mode **never**, proactively run tests" (l.175) — because the mode is literally named `never`, the sentence parses as its own negation. Prompt-language ambiguity shipped in production.
4. **Renderer-defensive rules**: the 【F:...】 citation ban (l.161), the ANSI-code ban (l.264), and the no-URI file-link grammar (l.237–239) are all guards against behaviors learned from *other* OpenAI surfaces (ChatGPT citations) — prompt text as cross-product behavior unlearning.
5. **Quantified micro-budgets**: preambles "8–12 words" (l.50), progress updates "8-10 words" (l.189), final answers "no more than 10 lines" (l.205), lists "4–6 bullets" (l.224), retries "up to 3 times" (l.169) — pervasive numeric control knobs, ideal for placebo-length ablation.
6. **Inverted safety content**: the only refusal-adjacent text *grants* permissions ("proprietary repos allowed", "vulnerability analysis allowed", l.143–145) — the mirror image of consumer-chat prompts; presumably counteracting over-refusal priors.
7. **Personality-to-formatting ratio**: Personality gets 391 chars (1.9%), presentation/formatting gets ~5,900 (28.6%) — the prompt is dominated by output-shape engineering, not behavior in the usual sense.

## 5. Relation to the GPT-5-Codex variant (same repo)

The per-model `gpt_5_codex_prompt.md` (separate corpus artifact) restates this file's
content in ~1/3 the characters, keeps the `rg` guidance and file-reference grammar
byte-identical, drops AGENTS.md/planning examples/approval modes, and adds
editing-safety rules absent here (dirty-worktree protection, `git reset --hard` ban,
ASCII default). See the companion analysis file.

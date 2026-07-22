# Analysis: OpenAI Codex CLI — GPT-5-Codex model-specific prompt (accessed 2026-07-20)

Corpus file: `corpus/tier-b-open-source/openai--codex-cli--gpt-5-codex-prompt--2026-07-20.md`
Body: lines 15–82 (6,621 chars) — raw prompt Markdown, no code scaffolding. Line numbers
refer to the corpus file. This is the per-model replacement for the base prompt when the
model is GPT-5-Codex (the corpus also has the base prompt as a separate artifact).

## 1. Structural outline

| Section | Lines | % of body |
|---|---|---|
| (identity line) | 15 | 2% |
| `## General` (rg preference) | 17–19 | 3% |
| `## Editing constraints` (ASCII, comments, apply_patch, dirty worktree, git safety) | 21–33 | **27%** |
| `## Plan tool` | 35–40 | 4% |
| `## Special user requests` (terminal one-offs, review mode) | 42–45 | 11% |
| `## Presenting your work and final message` | 47–62 | 22% |
| `### Final answer structure and style guidelines` (compressed style guide + file references) | 64–82 | 30% |

## 2. Functional inventory

| # | Category | Quote (verbatim) | Line | Paraphrase |
|---|---|---|---|---|
| 1 | identity | "You are Codex, based on GPT-5. You are running as a coding agent in the Codex CLI on a user's computer." | 15 | Model-specific identity (names the base model, unlike the base prompt). |
| 2 | tool-protocol | "prefer using `rg` or `rg --files` respectively because `rg` is much faster" | 19 | Byte-identical to base prompt l.278. |
| 3 | code-conventions | "Default to ASCII when editing or creating files." | 23 | Unicode only when justified and already present. |
| 4 | code-conventions | "You should not add comments like \"Assigns the value to the variable\"" | 24 | Succinct comments only for non-obvious code; "rare". |
| 5 | tool-protocol | "Do not use apply_patch for changes that are auto-generated" | 25 | Tool-choice heuristic: scripting for mass edits, apply_patch for single files. |
| 6 | vcs-git (NEW) | "NEVER revert existing changes you did not make unless explicitly requested, since these changes were made by the user" | 27 | Dirty-worktree protocol: user edits are sacrosanct. |
| 7 | vcs-git (NEW) | "Do not amend a commit unless explicitly requested to do so." | 31 | No history rewriting. |
| 8 | error-handling | "If this happens, STOP IMMEDIATELY and ask the user how they would like to proceed." | 32 | Unexpected external changes → halt, don't improvise. |
| 9 | vcs-git (NEW) | "**NEVER** use destructive commands like `git reset --hard` or `git checkout --`" | 33 | Destructive-command ban absent from base prompt. |
| 10 | planning | "Skip using the planning tool for straightforward tasks (roughly the easiest 25%)." | 38 | Quantified plan-skip threshold; no single-step plans (l.39). |
| 11 | planning | "update it after having performed one of the sub-tasks that you shared on the plan" | 40 | Keep plan state fresh. |
| 12 | tool-protocol | "which you can fulfill by running a terminal command (such as `date`), you should do so" | 44 | Use the shell even for trivial user requests. |
| 13 | verification / user-communication | "default to a code review mindset: prioritise identifying bugs, risks, behavioural regressions, and missing tests" | 45 | "Review" keyword triggers findings-first mode, ordered by severity, file/line refs. |
| 14 | conciseness | "Default: be very concise; friendly coding teammate tone." | 51 | Telegraphic restatement of base-prompt persona. |
| 15 | user-communication | "Ask only when needed; suggest ideas; mirror the user's style." | 52 | Adaptive dialogue posture. |
| 16 | environment-context | "No \"save/copy this file\" - User is on the same machine." | 56 | Same-computer assumption (compressed from base l.201). |
| 17 | user-communication | "Lead with a quick explanation of the change ... Do not start this explanation with \"summary\"" | 59 | Code-change reporting shape. |
| 18 | user-communication | "use numeric lists for the suggestions so the user can quickly respond with a single number" | 61 | Options as numbered lists — reply-affordance design. |
| 19 | user-communication | "The user does not command execution outputs." | 62 | [sic — verb missing, likely "see"] relay/summarize command output rather than assuming visibility. |
| 20 | formatting-output | "Headers: optional; short Title Case (1-3 words) wrapped in **…**" | 67 | Compressed style guide (bullets `-`, 4–6 per list, monospace, fenced code with info string, l.66–73). |
| 21 | formatting-output | "Code samples or multi-line snippets should be wrapped in fenced code blocks; include an info string" | 70 | Fencing rule not present in base prompt's style section. |
| 22 | audience-adaptation | "code explanations → precise, structured with code refs; simple tasks → lead with outcome" | 74 | Register adaptation matrix. |
| 23 | formatting-output | "Use inline code to make file paths clickable." | 76 | File-reference grammar, byte-identical to base prompt l.234–241. |

## 3. Layer split (approx., % of body chars)

- Behavioral instructions: ~85%
- Tool definitions/protocol: ~10% (apply_patch usage heuristics, rg, plan tool)
- Few-shot examples: ~2% (only inline path examples l.82)
- Template variables: 0%
- Other: ~3%

## 4. Idiosyncrasies

1. **A compressed mirror of the base prompt**: whole sections are the base prompt's
   sentences re-encoded telegraphically ("Headers: optional; short Title Case..."), at
   ~32% of the base prompt's length, while the `rg` paragraph and the file-reference
   block are carried over byte-identical. Strong natural experiment for the embedding
   pipeline: same semantic components, radically different surface form — these pairs
   should co-cluster if the embeddings capture function rather than style.
2. **Failure-mode archaeology**: everything *added* relative to the base prompt is
   protective (never revert user changes, STOP IMMEDIATELY on unexpected edits, no
   `git reset --hard`, no commit amends) — reads as a patch list for observed
   GPT-5-Codex incidents with users' uncommitted work.
3. **Shipped grammar bug**: "The user does not command execution outputs." (l.62) —
   missing verb ("see"), production prompt.
4. **Pseudo-quantified judgment**: "roughly the easiest 25%" of tasks skip planning
   (l.38) — a number the model cannot actually measure; quantification as vibes.
5. **Zero safety/refusal content**: not even the base prompt's permission grants; the
   entire safety surface is delegated to harness sandboxing and (presumably) model
   training.
6. **No AGENTS.md section**: the instruction-precedence spec from the base prompt is
   absent — per-model prompts do not restate repo-instruction handling, implying the
   harness composes it elsewhere or GPT-5-Codex is trusted to know it.

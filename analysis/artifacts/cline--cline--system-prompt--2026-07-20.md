# Analysis: Cline (VS Code coding agent) — system prompt module (accessed 2026-07-20)

Corpus file: `corpus/tier-b-open-source/cline--cline--system-prompt--2026-07-20.md`
Body: lines 15–82 (6,634 chars). SOURCE CODE: a TypeScript module exporting **two**
prompt template literals; scaffolding is only the two `export const` wrappers (~1%).
Line numbers refer to the corpus file.

## 1. Structural outline

| Section | Lines | % of body |
|---|---|---|
| `DEFAULT_CLINE_SYSTEM_PROMPT` (interactive agent) | 15–50 | 56% |
| — identity + mission | 15 | |
| — context-gathering doctrine | 17–19 | |
| — `<env>` block (4 template variables) | 21–27 | |
| — "Remember:" rule list (9 bullets) | 29–38 | |
| — plan-first + proactivity exhortations | 40–44 | |
| — completion/summary protocol | 46–48 | |
| — `{{CLINE_RULES}}` `{{CLINE_METADATA}}` slots | 49–50 | |
| `YOLO_CLINE_SYSTEM_PROMPT` (non-interactive/background agent) | 52–82 | 43% |
| — identity (background, no user contact) | 52–54 | |
| — RULES list (8 bullets, heavily overlapping DEFAULT) | 56–64 | |
| — `<env>` block (same 4 variables) | 66–72 | |
| — IMPORTANT: bug-fix + test-loop + `submit_and_exit` termination protocol | 74–80 | |
| — template slots | 81–82 | |
| TypeScript scaffolding (`export const ... = \`` and closing backticks) | 15, 50, 52, 82 | ~1% |

## 2. Functional inventory

| # | Category | Quote (verbatim) | Line | Paraphrase |
|---|---|---|---|---|
| 1 | identity | "You are Cline, an AI coding agent." | 15 | Name + role + tool-use mission. |
| 2 | planning | "Always gather all the necessary context before starting to work on a task." | 17 | Understand requirements/conventions/frameworks first; validate tests at end. |
| 3 | honesty-uncertainty | "use one of the available tools or ask for clarification instead of making assumptions or lies" | 19 | [sic] Don't fabricate; resolve unknowns via tools or questions. |
| 4 | environment-context | "1. Platform: {{PLATFORM_NAME}}" | 23 | `<env>` block: platform, date, IDE, CWD injected at runtime (21–27). |
| 5 | code-conventions | "Always adhere to existing code conventions and patterns." | 30 | Match house style. |
| 6 | code-conventions | "Use only libraries and frameworks that are confirmed to be in use in the current codebase." | 31 | No new dependencies by default. |
| 7 | code-conventions | "Provide complete and functional code without omissions or placeholders." | 32 | Anti-stub rule. |
| 8 | honesty-uncertainty | "Be explicit about any assumptions or limitations in your solution." | 33 | Surface assumptions. |
| 9 | planning | "Always show your planning process before executing any task." | 34 | Visible plan-first requirement ("OK for this section to be quite long", l.40). |
| 10 | tool-protocol | "Always use absolute paths when referring to files." | 35 | Path discipline. |
| 11 | tool-parallelism | "identify every independent read, search, command, or edit needed for the next step and emit all of those tool calls now" | 36 | Batch all independent calls per turn; don't serialize independent work. |
| 12 | tool-parallelism | "read all known relevant files in one read_files call; run independent inspection commands in one run_commands call" | 37 | Concrete parallelism patterns (batched array-accepting tools). |
| 13 | verification | "Always verify the files you have edited or created at the end of the task" | 38 | End-of-task verification pass. |
| 14 | agentic-persistence | "REMEMBER, be helpful and proactive! Don't ask for permission to do something when you can do it!" | 42 | Bias to act; no permission-seeking. |
| 15 | user-communication | "Do not indicates you will be using a tool unless you are actually going to use it." | 42 | [sic] No announced-but-unexecuted actions. |
| 16 | tool-protocol | "Response without tool calls will considered as completed with final answer." | 44 | [sic] Turn-termination semantics encoded in prose: no tool calls = task done. |
| 17 | user-communication | "please provide a summary of what you did and any relevant information" | 46 | Final summary protocol; validate by running code when possible. |
| 18 | audience-adaptation | "If user asked a simple question without any coding context, answer it directly without using any tools." | 48 | Chat-mode escape hatch. |
| 19 | memory-context-mgmt | "{{CLINE_RULES}}" | 49 | Injection slot for user/workspace rule files (+ `{{CLINE_METADATA}}`, l.50). |
| 20 | identity | "You are Cline, a careful and helpful coding agent that works in the background." | 52 | YOLO variant identity. |
| 21 | user-communication | "an issue reported by the user who you cannot communicate with directly" | 53 | Non-interactive setting: no clarification channel. |
| 22 | formatting-output | "Always match output format exactly as shown in examples or existing files." | 57 | YOLO: format fidelity to existing artifacts. |
| 23 | error-handling | "A correct fix means the underlying behavior is fixed — not just the symptoms addressed superficially." | 76 | Root-cause doctrine for bug reports. |
| 24 | verification | "If tests fail, analyze the failures, revise your fix, and re-run until tests pass." | 77 | Mandatory test loop; incomplete until touched-file suite passes (l.78). |
| 25 | tool-protocol | "You should only end the task when all the requirements are met by calling the 'submit_and_exit' tool." | 79 | YOLO termination = explicit tool call; "Response without the submit_and_exit tool call will considered not completed" (l.80). |

## 3. Layer split (approx., % of body chars)

- Behavioral instructions: ~80%
- Tool definitions: ~5% (no schemas; tool names cited in parallelism examples and `submit_and_exit` protocol)
- Few-shot examples: 0% (the "good parallelism examples" l.37/63 are inline patterns, not exchanges)
- Template variables: ~5% (12 `{{...}}` occurrences: PLATFORM_NAME, CURRENT_DATE, IDE_NAME, CWD ×2 each; CLINE_RULES, CLINE_METADATA ×2)
- Other (TS scaffolding): ~1%

## 4. Idiosyncrasies

1. **Pervasive grammar errors shipped in production**: "making assumptions or lies"
   (l.19), "Do not indicates" (l.42), "Always includes tool calls" (l.44/79), "will
   considered as completed" (l.44), "will considered not completed" (l.80), "Always
   validate your answer with checking the code" (l.46). Relevant as embedding noise and
   as evidence prompts skip editorial review even in popular products (Cline is one of
   the most-installed VS Code agents).
2. **Turn semantics in natural language**: both variants encode the harness's
   turn-termination contract as prose — DEFAULT: absence of tool calls means "done";
   YOLO: only `submit_and_exit` ends the task. The same mechanism most harnesses
   implement in code is here a prompt component (good ablation target: break it and the
   agent loop should visibly degrade).
3. **Internal duplication**: the tool-parallelism paragraph and examples are
   byte-identical between DEFAULT (l.36–37) and YOLO (l.62–63); env blocks likewise —
   two prompts maintained by copy-paste in one file, drift risk visible.
4. **Latent tension**: "ask for clarification instead of making assumptions" (l.19)
   vs "Don't ask for permission to do something when you can do it!" (l.42) — the
   boundary between clarification and permission is left to the model.
5. **"YOLO" naming vs content**: the non-interactive variant self-describes as
   "careful and helpful" (l.52) while the export name says YOLO — naming reflects the
   autonomy mode, not the persona.
6. **Zero safety/refusal, zero secrets handling**: no security text of any kind in
   either variant; the `{{CLINE_RULES}}` slot delegates policy to users.
7. **Strongest tool-parallelism instruction in the subset**: Cline is the only
   artifact in this 8-file subset with an explicit multi-call batching doctrine —
   a clean single-source category exemplar for the codebook.

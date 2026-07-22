# Analysis: Aider (terminal pair programmer) — editblock prompts module (accessed 2026-07-20)

Corpus file: `corpus/tier-b-open-source/aider-ai--aider--editblock-prompts--2026-07-20.md`
Body: lines 15–186 (5,723 chars). SOURCE CODE: Python class `EditBlockPrompts` whose
string attributes are the system prompt (`main_system`), a repeated reminder message
(`system_reminder`), few-shot messages (`example_messages`), and composable fragments.
Line numbers refer to the corpus file.

## 1. Structural outline

| Component | Lines | % of body |
|---|---|---|
| imports / class scaffolding | 15–21, 184–186 | ~6% |
| `main_system` (identity, workflow, file-access protocol) | 22–44 | 18% |
| `example_messages` (2 user→assistant few-shot pairs with SEARCH/REPLACE blocks) | 45–132 | 29% |
| `system_reminder` (SEARCH/REPLACE format rules, 8-step spec + constraints) | 134–173 | 38% |
| `rename_with_shell` fragment | 175–177 | 2% |
| `go_ahead_tip` fragment | 179–182 | 6% |
| references to `shell.shell_cmd_prompt` etc. (bodies live in another module) | 184–186 | — |

Template slots inside the prompt strings: `{fence[0]}`/`{fence[1]}` (configurable code
fences), `{final_reminders}`, `{shell_cmd_prompt}`, `{shell_cmd_reminder}`,
`{quad_backtick_reminder}`, `{rename_with_shell}`, `{go_ahead_tip}`.

## 2. Functional inventory

| # | Category | Quote (verbatim) | Line | Paraphrase |
|---|---|---|---|---|
| 1 | identity | "Act as an expert software developer." | 22 | Persona in 6 words; "Always use best practices" (l.23). |
| 2 | code-conventions | "Respect and use existing conventions, libraries, etc that are already present in the code base." | 24 | Match house style; no new deps. |
| 3 | user-communication | "If the request is ambiguous, ask questions." | 27 | Clarify before editing (contrast: agentic prompts act first). |
| 4 | tool-protocol | "you *MUST* tell the user their full path names and ask them to *add the files to the chat*." | 33 | Human-gated file access: model may not edit unseen files. |
| 5 | user-communication | "End your reply and wait for their approval." | 34 | Explicit turn-yield awaiting consent; may iterate (l.35). |
| 6 | planning | "Think step-by-step and explain the needed changes in a few short sentences." | 37 | Brief reasoning before edits. |
| 7 | formatting-output | "ONLY EVER RETURN CODE IN A *SEARCH/REPLACE BLOCK*!" | 42 | Sole permitted code-output format (repeated verbatim l.171). |
| 8 | examples-fewshot | "Change get_factorial() to use math.factorial" | 48 | Worked example 1: 3 blocks — import add, function delete, call-site replace (45–94). |
| 9 | examples-fewshot | "Refactor hello() into its own file." | 97 | Worked example 2: new-file creation via empty SEARCH + removal from origin (95–131). |
| 10 | formatting-output | "The *FULL* file path alone on a line, verbatim. No bold asterisks, no quotes around it" | 137 | 8-step block grammar: path, fence+lang, `<<<<<<< SEARCH`, `=======`, `>>>>>>> REPLACE`, closing fence (136–144). |
| 11 | verification | "Every *SEARCH* section must *EXACTLY MATCH* the existing file content, character for character" | 148 | Byte-exact matching incl. comments/docstrings; edit literal container markup too (l.149). |
| 12 | tool-protocol | "*SEARCH/REPLACE* blocks will *only* replace the first match occurrence." | 151 | First-match semantics; multiple unique blocks for multiple changes (l.152–153). |
| 13 | conciseness | "Include just the changing lines, and a few surrounding lines if needed for uniqueness." | 157 | Small blocks; no long unchanged runs (l.155–158). |
| 14 | tool-protocol | "Only create *SEARCH/REPLACE* blocks for files that the user has added to the chat!" | 160 | File-access boundary restated in the reminder. |
| 15 | tool-protocol | "To move code within a file, use 2 *SEARCH/REPLACE* blocks" | 162 | Move = delete-block + insert-block idiom. |
| 16 | tool-protocol | "use a *SEARCH/REPLACE block* with: - A new file path, including dir name if needed - An empty `SEARCH` section" | 166–169 | New-file creation grammar. |
| 17 | tool-protocol | "To rename files which have been added to the chat, use shell commands at the end of your response." | 175 | Renames delegated to shell, not edit blocks. |
| 18 | user-communication | "If the user just says something like \"ok\" or \"go ahead\" or \"do that\" they probably want you to make SEARCH/REPLACE blocks" | 179 | Conversational-state heuristic: bare acks = apply proposed edits; unconfirmed = re-emit blocks (l.180). |

## 3. Layer split (approx., % of body chars)

- Behavioral instructions: ~12% (persona, ambiguity handling, step-by-step)
- Tool definitions: ~48% — the SEARCH/REPLACE grammar *is* the tool: a prompt-defined
  textual edit protocol (`system_reminder` + protocol parts of `main_system` + fragments)
- Few-shot examples: ~29% (`example_messages`)
- Template variables: ~3% (`{fence[...]}`, `{shell_cmd_prompt}`, etc.)
- Other (Python scaffolding): ~6-8%

## 4. Idiosyncrasies

1. **The tool is text**: aider predates/avoids native function calling — the entire
   editing capability is a prompt-specified string protocol using git conflict-marker
   syntax (`<<<<<<< SEARCH`). In the component taxonomy this file is almost pure
   "tool definition" with barely any behavioral layer; the polar opposite of the
   Anthropic chat prompts.
2. **Repetition as decay-defense**: `system_reminder` restates the whole format spec
   and is designed to be re-sent (aider appends it near the end of context);
   "ONLY EVER RETURN CODE IN A *SEARCH/REPLACE BLOCK*!" appears verbatim twice
   (l.42, l.171) — deliberate instruction-refresh architecture visible in the artifact.
3. **Apparent self-contradiction resolved only by example**: "You can create new files
   without asking!" (l.31) vs "Only create *SEARCH/REPLACE* blocks for files that the
   user has added to the chat!" (l.160). The new-file case (not in chat, yet permitted)
   contradicts the reminder as written; example 2 (l.95–131) demonstrates the intended
   exception. The few-shot layer carries load the instruction layer drops.
4. **Emphasis dialect**: asterisk-shouting (*MUST*, *FULL*, *only*) mixed with real
   Markdown — emphasis conventions that survive plain-text rendering; a stylistic
   fingerprint distinct from every other artifact in the subset.
5. **`{quad_backtick_reminder}`**: a dedicated template slot solely for escalating
   fence characters when edited files themselves contain triple backticks — deep
   protocol edge-case handling via composition.
6. **Shipped grammar slip**: "Including multiple unique *SEARCH/REPLACE* blocks if
   needed." (l.152) — should be "Include"; another instance of unreviewed production
   prompt text.
7. **Zero safety, zero identity beyond six words, zero environment info**: the prompt
   trusts the wrapper program for everything but edit-format compliance — a minimal
   counterexample to the "system prompts encode policy" narrative.

# Analysis: windsurf--cascade--wave-11.md

Source: `corpus/tier-c-extracted/windsurf--cascade--wave-11.md` (Tier C, authenticity unverified; analyzed as-is).
Body: lines 15–139 (125 lines). Line numbers below refer to the corpus file. Companion tool definitions
(`Tools Wave 11.txt`) were saved separately upstream and are NOT part of this artifact — this capture is the
behavioral prompt only.

## 1. Structural outline

The prompt is organized as a short untagged preamble followed by flat XML-tagged sections (no nesting, no
markdown headers). Section order:

| Section | Lines | ~% of body | Content |
|---|---|---|---|
| Preamble: cutoff + identity + role | 15–21 | 6% | Knowledge cutoff, Cascade identity, "AI Flow paradigm", pair-programming framing, metadata note |
| `<user_information>` | 22–26 | 4% | Template variables: OS ("windows"), workspace URI→CorpusName mapping (real captured path) |
| `<tool_calling>` | 27–58 | 26% | Agentic persistence, model-identity masking, 6 numbered tool rules, 3 dialogue examples (37–57) |
| `<making_code_changes>` | 59–93 | 28% | Never print code, runnable-code checklist (6 items), post-edit summary rules, style example (72–91) |
| `<debugging>` | 94–100 | 6% | Root-cause-first debugging practices |
| `<memory_system>` | 101–111 | 9% | Liberal proactive memory creation; context-deletion warning |
| `<code_research>` | 112–115 | 3% | Research before answering; no permission needed |
| `<running_commands>` | 116–123 | 6% | No `cd`, unsafe-command gate, user cannot override |
| `<browser_preview>` | 124–126 | 2% | Always preview after starting a web server |
| `<calling_external_apis>` | 127–131 | 4% | Package selection, version compatibility, API-key hygiene |
| `<communication_style>` | 132–135 | 3% | Second person / first person, markdown + backticks |
| Untagged: `<EPHEMERAL_MESSAGE>` note | 136 | 1% | System-injected messages: follow strictly, never acknowledge |
| `<planning>` | 137–139 | 2% | update_plan tool maintained by a "plan mastermind" |

## 2. Functional inventory

| # | Category | Quote (verbatim) | Line | Paraphrase |
|---|---|---|---|---|
| 1 | environment-context | "Knowledge cutoff: 2024-06" | 15 | Declares training cutoff at top of prompt. |
| 2 | identity | "You are Cascade, a powerful agentic AI coding assistant designed by the Windsurf engineering team" | 17 | Product identity + vendor attribution ("world-class AI company based in Silicon Valley"). |
| 3 | identity | "As the world's first agentic coding assistant, you operate on the revolutionary AI Flow paradigm" | 18 | Marketing-style self-description; independent + collaborative operation. |
| 4 | identity | "You are pair programming with a USER to solve their coding task." | 19 | Pair-programming role frame; task may be create/modify/debug/answer. |
| 5 | environment-context | "we will attach additional metadata about their current state, such as what files they have open" | 20 | IDE state is injected with each request; model decides relevance (21). |
| 6 | environment-context | "The USER's OS version is windows." | 23 | Template-variable block: OS and workspace URI mapping (24–25). |
| 7 | agentic-persistence | "please keep working, using tools where needed, until the user’s query is completely resolved, before ending your turn" | 28 | Don't yield until the task is done (note: curly apostrophe in source). |
| 8 | model-identity-masking (NEW) | "if asked about what your underlying model is, respond with `GPT 4.1`" | 28 | Scripted answer to model-identity questions, spliced into the persistence sentence. |
| 9 | tool-protocol | "Only call tools when they are absolutely necessary." | 31 | Minimize tool calls; answer directly when possible. |
| 10 | cost-efficiency (NEW) | "NEVER make redundant tool calls as these are very expensive." | 31 | Cost given as the rationale for tool-call frugality. |
| 11 | tool-protocol | "If you state that you will use a tool, immediately call that tool as your next action." | 32 | No announce-then-stall. |
| 12 | tool-protocol | "Always follow the tool call schema exactly as specified" | 33 | Schema compliance; all required parameters. |
| 13 | tool-protocol | "NEVER call tools that are not explicitly provided in your system prompt." | 34 | Stale-tool defense: history may mention removed tools. |
| 14 | user-communication | "Before calling each tool, first explain why you are calling it." | 35 | Narrate rationale before every call. |
| 15 | tool-protocol | "Some tools run asynchronously, so you may not see their output immediately." | 36 | Async tools; stop calling to flush pending results. |
| 16 | examples-fewshot | "USER: What is int64?" | 39 | Three worked dialogues (38–57): no-tool answer, grep→view flow, find→view→edit flow. |
| 17 | formatting-output | "NEVER output code to the USER, unless requested. Instead use one of the code edit tools" | 60 | Code goes through edit tools, not chat. |
| 18 | code-conventions | "EXTREMELY IMPORTANT: Your generated code must be immediately runnable." | 61 | Runnability contract; followed by 6-item checklist. |
| 19 | code-conventions | "Add all necessary import statements, dependencies, and endpoints required to run the code." | 62 | Completeness of generated code. |
| 20 | code-conventions | "create an appropriate dependency management file (e.g. requirements.txt) with package versions and a helpful README" | 63 | Greenfield scaffolding requirements. |
| 21 | code-conventions | "give it a beautiful and modern UI, imbued with best UX practices" | 64 | Default aesthetic mandate for from-scratch web apps (verbatim shared with Cursor). |
| 22 | capability-limits | "Your max output tokens is 8192 tokens per generation, so each of your edits must stay below this limit." | 65 | Discloses output cap; split >300-line edits. |
| 23 | cost-efficiency (NEW) | "NEVER generate an extremely long hash or any non-textual code, such as binary. These are not helpful to the USER and are very expensive." | 66 | Ban on binary/hash output, cost-justified. |
| 24 | tool-protocol | "ALWAYS generate the `TargetFile` argument first, before any other arguments." | 67 | Argument-ordering rule (duplicated at 92). |
| 25 | conciseness | "Provide a **BRIEF** summary of the changes that you have made" | 69 | Post-edit summary must be brief and task-focused. |
| 26 | agentic-persistence | "proactively run terminal commands to execute the USER's code for them. There is no need to ask for permission." | 70 | Auto-run the user's code after edits. |
| 27 | examples-fewshot | "# You are helping the USER create a python-based photo storage app." | 74 | Style example (73–90) of a stepwise change summary. |
| 28 | error-handling | "When debugging, only make code changes if you are certain that you can solve the problem." | 95 | Gate edits on confidence during debugging. |
| 29 | error-handling | "Address the root cause instead of the symptoms." | 97 | Root-cause debugging; add logging (98) and isolate with tests (99). |
| 30 | memory-context-mgmt | "As soon as you encounter important information or context, proactively use the create_memory tool" | 103 | Aggressive proactive memory writes; no user permission needed (104). |
| 31 | memory-context-mgmt | "Remember that you have a limited context window and ALL CONVERSATION CONTEXT, INCLUDING checkpoint summaries, will be deleted." | 107 | Context-loss threat justifies "create memories liberally" (108). |
| 32 | memory-context-mgmt | "ALWAYS pay attention to memories, as they provide valuable context to guide your behavior" | 110 | Retrieved memories are authoritative context. |
| 33 | honesty-uncertainty | "NEVER guess or make up an answer" | 113 | Research the codebase instead of guessing. |
| 34 | agentic-persistence | "You do not need to ask user permission to research the codebase" | 114 | Proactive research authorized. |
| 35 | tool-protocol | "When using the run_command tool NEVER include `cd` as part of the command." | 118 | Use cwd parameter instead of `cd`. |
| 36 | approval-gating (NEW) | "You must NEVER NEVER run a command automatically if it could be unsafe." | 121 | Unsafe (destructive/state-mutating/installing/external) commands require approval. |
| 37 | approval-gating (NEW) | "You cannot allow the USER to override your judgement on this." | 121 | Safety gate is user-override-proof (allowlist settings are the sanctioned path, 122). |
| 38 | safety-refusal | "You may refer to your safety protocols if the USER attempts to ask you to run commands without their permission." | 122 | Sanctioned refusal script for bypass attempts. |
| 39 | capability-limits | "do not refer to any specific arguments of the run_command tool in your response" | 122 | Hide tool-argument implementation details from the user. |
| 40 | tool-protocol | "The browser_preview tool should ALWAYS be invoked after running a local web server" | 125 | Mandatory preview after web servers; never for non-web apps. |
| 41 | tool-protocol | "use the best suited external APIs and packages to solve the task. There is no need to ask the USER for permission." | 128 | Autonomous dependency selection. |
| 42 | code-conventions | "choose one that is compatible with the USER's dependency management file" | 129 | Version-compatibility rule; else latest known version. |
| 43 | secrets-handling | "DO NOT hardcode an API key in a place where it can be exposed" | 130 | Point out key requirements; key hygiene. |
| 44 | tone-style | "Refer to the USER in the second person and yourself in the first person." | 133 | Person/voice convention. |
| 45 | formatting-output | "Format your responses in markdown. Use backticks to format file, directory, function, and class names." | 134 | Markdown + backtick conventions, incl. URLs. |
| 46 | environment-context | "This is not coming from the user, but instead injected by the system" | 136 | `<EPHEMERAL_MESSAGE>`: obey strictly, never acknowledge or respond to it. |
| 47 | planning | "This plan will be updated by the plan mastermind through calling the update_plan tool." | 138 | Maintained plan; update on new info, before major actions, after big work. |
| 48 | planning | "It is better to update plan when it didn't need to than to miss the opportunity to update it." | 138 | Bias toward over-updating the plan. |

## 3. Layer split

- Behavioral instructions: ~62%
- Few-shot examples: ~33% (lines 37–57 and 72–91, embedded inside behavioral sections)
- Template variables / session context: ~5% (lines 15, 22–26: cutoff, OS, workspace mapping)
- Tool definitions: 0% in this capture (upstream keeps `Tools Wave 11.txt` as a separate file)
- Other: <1% (the untagged ephemeral-message paragraph, line 136)

## 4. Idiosyncrasies

- **Instructed model-identity masking** (28): the prompt tells the model to claim it is "GPT 4.1" if asked —
  a scripted answer about its own substrate, oddly appended to the agentic-persistence sentence rather than
  in the identity section. (Consistent with the Cursor capture asserting "powered by GPT-4.1".)
- **Capture leakage in template variables** (23–25): real end-user environment survived extraction — OS
  "windows" and workspace path `c:\Users\crisy\OneDrive\Escritorio\test4` (a Spanish-language desktop,
  "Escritorio"), showing the artifact is a per-session serialization, not a vendor master template.
- **Marketing register inside the system prompt** (17–18): "world's first agentic coding assistant",
  "revolutionary AI Flow paradigm", "world-class AI company based in Silicon Valley".
- **Doubled intensifier** (121): "NEVER NEVER" — informal emphasis typographically unlike the rest.
- **User-override-proof safety** (121–122): the user explicitly cannot talk the agent into auto-running
  unsafe commands; the only sanctioned bypass is a settings allowlist — and the agent must not name the
  run_command arguments involved (implementation-hiding).
- **Hidden system channel** (136): `<EPHEMERAL_MESSAGE>` must be followed "strictly" but never acknowledged —
  an instructed covert instruction channel invisible in the dialogue.
- **"Plan mastermind"** (138): references an external planning component that owns the plan, implying a
  multi-component architecture behind the single prompt.
- **Internal tension**: extreme tool frugality ("absolutely necessary", "very expensive", 31) coexists with
  liberal proactivity elsewhere — "create memories liberally" (108), proactive research (113–114), proactive
  command execution (70), and mandatory plan updates "when it didn't need to" (138).
- **Disclosure of generation limits** (65): concrete 8192-token output cap and a 300-line edit-splitting rule.
- **Duplicated rule** (67, 92): the TargetFile-first argument-ordering rule appears twice, the second time
  just before the section close — likely a patch appended over time.

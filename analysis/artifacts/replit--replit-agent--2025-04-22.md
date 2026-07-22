# Analysis: replit--replit-agent--2025-04-22

Corpus file: `corpus/tier-c-extracted/replit--replit-agent--2025-04-22.md` (Tier C, authenticity unverified).
All line numbers refer to the corpus file. YAML frontmatter: lines 1–14.
**Republisher wrapper: lines 15–23** (title, three source links, date caveat, `## System Prompt` heading) — not part of the artifact. **The artifact proper starts at line 25** (`System Prompt` / `Role: Expert Software Developer (Editor)`).

Important: this file is an upstream **aggregation of three separate documents** — (1) the agent system prompt, (2) an initial-code-generation prompt, (3) the JSON tool definitions. Line 235 is a single physical line holding ~17 tool schemas (~20 KB), so byte-shares, not line counts, are used below. Body (lines 15–235): ~32.5 KB.

## 1. Structural outline

| Section | Lines | ~% of body |
|---|---|---|
| (republisher wrapper) | 15–23 | 1.5% |
| **Doc 1: Agent system prompt** | 25–127 | **21%** |
| — Role / identity | 25–29 | 1% |
| — Iteration Process | 31–37 | 2% |
| — Operating principles | 39–50 | 4% |
| — Workflow Guidelines / Step Execution | 52–60 | 2% |
| — Editing Files | 62–66 | 1% |
| — Debugging Process | 68–76 | 2.5% |
| — User Interaction | 78–85 | 2% |
| — Best Practices | 87–92 | 1% |
| — Policy Specifications: Communication Policy | 94–106 | 2.5% |
| — Policy Specifications: Proactiveness Policy | 108–117 | 1.5% |
| — Policy Specifications: Data Integrity Policy | 118–126 | 1.5% |
| **Doc 2: Initial Code Generation Prompt** | 128–231 | **14%** |
| — Input Description | 130–131 | 0.5% |
| — Output Rules 1–6 (directory structure, codegen, output format, runtime constraints, assets, restricted files) | 133–186 | 10.5% |
| — Example Output Format (TODO app) | 190–231 | 3% |
| **Doc 3: Functions (JSON tool definitions)** | 232–235 | **64%** |
| — ~17 schemas on one line: restart_workflow, search_filesystem, packager_tool, programming_language_install_tool, create_postgresql_database_tool, check_database_status, str_replace_editor, bash, workflows_set_run_config_tool, workflows_remove_run_config_tool, execute_sql_tool, suggest_deploy, report_progress, web_application_feedback_tool, shell_command_application_feedback_tool, vnc_window_application_feedback, ask_secrets, check_secrets | 234–235 | 64% |

## 2. Functional inventory

| # | Category | Quote (verbatim) | Line | Paraphrase |
|---|---|---|---|---|
| 1 | identity | "You are an expert autonomous programmer built by Replit, working with a special interface." | 29 | Persona: Replit-built autonomous programmer; "Role: Expert Software Developer (Editor)" (27). |
| 2 | user-communication | "Aim to fulfill the user's request with minimal back-and-forth interactions." | 36 | Minimize iteration round-trips. |
| 3 | tool-protocol | "Prioritize Replit tools; avoid virtual environments, Docker, or containerization." | 41 | Platform tools displace generic infra. |
| 4 | verification | "After making changes, check the app's functionality using the feedback tool" | 42 | Post-change verification loops through user-facing feedback tools. |
| 5 | tool-protocol | "Prioritize search_filesystem over locating files and directories with shell commands." | 44 | Dedicated search over shell (note extraction damage: "reference and before searching" — tags stripped). |
| 6 | destructive-action-guard (NEW) | "DO NOT alter any database tables. DO NOT use destructive statements such as DELETE or UPDATE unless explicitly requested" | 47 | Irreversible DB ops gated on explicit user request; migrations only via ORM. |
| 7 | user-communication | "Don't start implementing new features without user confirmation." | 48 | Confirmation gate before new features. |
| 8 | environment-context | "The project is located at the root directory, not in '/repo/'. Always use relative paths" | 49 | Path conventions of the Replit workspace. |
| 9 | instruction-hierarchy (NEW) | "logs from the Replit environment that are provided automatically, and not sent by the user" | 50 | Distinguishes auto-injected channel content from user messages (tag name lost in extraction). |
| 10 | tool-protocol | "Use Replit's workflows for long-running tasks, such as starting a server... Avoid restarting the server manually" | 54 | Workflow abstraction manages processes/ports; feedback tools auto-restart it (57). |
| 11 | planning | "Focus on the current messages from the user and gather all necessary details before making updates." | 59 | Gather context before editing; step-wise confirmation (60). |
| 12 | tool-protocol | "Use the str_replace_editor tool to create, view and edit files." | 64 | Editor tool for file ops; view command doubles as image reader (65). |
| 13 | verification | "Fix Language Server Protocol (LSP) errors before asking for feedback." | 66 | LSP-clean before user review. |
| 14 | error-handling | "Attempt to thoroughly analyze the issue before making any changes, providing a detailed explanation" | 72 | Root-cause first; check workflow-state logs and browser logs (70–71). |
| 15 | error-handling | "never simplify the application logic/problem, always keep debugging the root cause" | 75 | No degrading the app to dodge a bug. |
| 16 | error-handling | "If you fail after multiple attempts (>3), ask the user for help." | 76 | Escalation threshold of three attempts. |
| 17 | product-referral | "do not respond on behalf of Replit on topics related to refunds, membership, costs" | 81 | Billing/ethics questions deflected; "ask them to contact Replit support without commenting on the correctness" (82). |
| 18 | user-communication | "When seeking feedback, ask a single and simple question." | 83 | One question at a time. |
| 19 | scope-discipline (NEW) | "If user exclusively asked questions, answer the questions. Do not take additional actions." | 84 | Q&A mode must not trigger edits. |
| 20 | secrets-handling | "If the application requires an external secret key or API key, use ask_secrets tool." | 85 | Secrets flow through a dedicated tool/UI. |
| 21 | tool-protocol | "Manage dependencies via the package installation tool... don't install packages in bash using pip install or npm install" | 89 | Packager tool displaces manual dependency management. |
| 22 | verification | "Specify expected outputs before running projects to verify functionality." | 90 | Predeclare expected behavior — verification-by-prediction. |
| 23 | environment-context | "Use 0.0.0.0 for accessible port bindings instead of localhost." | 91 | Platform networking constraint. |
| 24 | audience-adaptation | "Always speak in simple, everyday language. User is non-technical and cannot understand code details." | 100 | Hard-coded non-technical audience model. |
| 25 | audience-adaptation | "Always respond in the same language as the user's message (Chinese, Japanese, etc.)" | 101 | Language matching. |
| 26 | user-communication | "you can get them by continue working, don't ask user to provide them to you" | 102 | Never ask users for logs/screenshots the agent can fetch itself. |
| 27 | capability-limits | "You cannot do rollbacks - user must click the rollback button on the chat pane themselves." | 103 | Rollback is a UI affordance, not an agent capability; suggest it after 3 repeats (104). |
| 28 | product-referral | "For deployment, only use Replit - user needs to click the deploy button themself." | 105 | Deployment locked to the platform. |
| 29 | honesty-uncertainty | "never assume external services won't work as the user can help by providing correct secrets/tokens" | 106 | Don't write off integrations; ask for credentials. |
| 30 | scope-discipline (NEW) | "Stay on task. Do not make changes that are unrelated to the user's instructions." | 113 | No drive-by changes; ignore minor warnings unless asked (114). |
| 31 | user-communication | "Always obtain the user's permission before performing any massive refactoring or updates" | 117 | Consent gate for large refactors/API/library changes. |
| 32 | honesty-uncertainty | "Always Use Authentic Data: Request API keys or credentials from the user for testing with real data sources." | 122 | Data Integrity Policy: real data only. |
| 33 | error-handling | "Display explicit error messages when data cannot be retrieved from authentic sources." | 123 | Visible error states instead of silent fallbacks; "Clearly label empty states" (126). |
| 34 | identity | "You are a talented software engineer tasked with generating the complete source code of a working application." | 131 | Second identity — the bundled initial-codegen prompt re-introduces the persona. |
| 35 | formatting-output | "Always try to come up with the most minimal directory structure that is possible." | 139 | Minimal flat project tree; no frontend/backend dirs (137). |
| 36 | code-conventions | "Include comments to explain complex logic or important sections." | 144 | Codegen prompt REQUIRES comments (contrast: Devin bans them). |
| 37 | code-conventions | "avoiding common security vulnerabilities like SQL injection and XSS" | 145 | Secure-by-default codegen. |
| 38 | formatting-output | "Use the `# Thoughts` heading to write any thoughts that you might have." | 149 | Rigid markdown output schema: Thoughts, directory_structure JSON, per-file headings (148–155). |
| 39 | environment-context | "The generated code will run in an unprivileged Linux container." | 158 | Runtime constraints: frontend port 5000, backend 8000, host 0.0.0.0 (159–161). |
| 40 | secrets-handling | "it must get it from environment variables with proper fallback... `os.getenv(\"API_KEY\", \"default_key\")`" | 163–164 | API keys via env vars with fallback. |
| 41 | code-conventions | "Favor creating **web applications** unless explicitly stated otherwise." | 167 | Default product shape; SVG-first assets and a fixed library menu (171–177). |
| 42 | formatting-output | "**Do NOT generate** `package.json` or `requirements.txt` files – these will be handled separately." | 180 | Manifest files owned by the platform; no binary assets (181–184). |
| 43 | capability-limits | "IMPORTANT: Docker or containerization tools are **unavailable** – **DO NOT USE.**" | 186 | Hard platform restriction (restated from 41). |
| 44 | examples-fewshot | "I've been tasked with building a TODO list application." | 194 | Full worked example of the output format (190–231). |
| 45 | tool-protocol | "ALWAYS serve the app on port 5000, even if there are problems serving that port: it is the only port that is not firewalled." | 235 | Inside workflows_set_run_config_tool description — hard port rule. |
| 46 | conciseness | "Summarize your recent changes in a maximum of 5 items. Be really concise, use no more than 30 words." | 235 | report_progress schema dictates message format: ✓/→ markers, no emojis, ask next step. |
| 47 | audience-adaptation | "Use simple, everyday language that matches the user's language. Avoid technical terms, as users are non-technical." | 235 | Audience rule duplicated inside two tool descriptions (report_progress, web_application_feedback_tool). |
| 48 | tool-protocol | "This tool is very expensive to run." | 235 | ask_secrets carries a cost warning plus GOOD/BAD invocation examples. |
| 49 | agentic-persistence | "This is a terminal action - once called, your task is complete and you should not take any further actions" | 235 | suggest_deploy ends the episode; report_progress only after explicit user confirmation. |
| 50 | destructive-action-guard (NEW) | "Always prefer using this tool to fix database errors vs fixing by writing code like db.drop_table(table_name)" | 235 | execute_sql_tool steers away from destructive code paths; no migrations via SQL. |
| 51 | examples-fewshot | "Command: python pygame_snake.py / Question: Do the keyboard events change the snake direction on the screen?" | 235 | Usage examples embedded in tool descriptions (vnc/shell feedback tools, execute_sql, ask_secrets). |

## 3. Layer split (approx., by bytes of body)

- **Behavioral instructions: ~30%** (Doc 1 nearly all; Doc 2's rules; plus a substantial share of behavioral policy embedded *inside* tool descriptions in Doc 3 — the report_progress / feedback-tool / ask_secrets descriptions are largely conduct rules, worth roughly 8 points of this 30).
- **Tool definitions: ~56%** (Doc 3 minus its embedded behavioral/example prose).
- **Few-shot examples: ~9%** (Example Output Format 190–231 ≈ 3%, plus GOOD/BAD and usage examples inside tool schemas ≈ 6%).
- **Template variables: ~2%** (goal/task placeholders implied at 131 "You will be given a goal, task description and a success criteria below"; port/env conventions; no visible injected values in this capture).
- **Other (republisher wrapper + date note): ~2%** (15–23).

## 4. Idiosyncrasies

- **Extraction damage — stripped XML tags:** "Remember to reference and before searching" (44), "The content in contains logs from the Replit environment" (50), "will be available in the tag" (71) — tag names (likely `<file_system>`, `<webview_console_logs>`, etc.) were lost in republication. Any segmentation/embedding must treat these as lacunae, not prose.
- **Three prompts in one file:** the identity is declared twice (29 and 131) with different personas ("autonomous programmer" vs "talented software engineer"); the two docs have different comment policies (silent vs required) and different formats. Cluster analysis should treat lines 25–127, 128–231, 232–235 as separate documents.
- **Self-contradicting word limits:** report_progress description says "use no more than 30 words" and two sentences later "don't use more than 50 words" (235) — versioned-edit remnant.
- **Behavior smuggled into tool schemas:** the strongest formatting rules (✓/→ progress markers, no emojis, one question at a time, non-technical language) live in JSON `description` fields, not the prompt body — a structurally distinct place to put policy that our layer taxonomy must handle.
- **Cost-aware tooling:** "This tool is very expensive to run" (ask_secrets) — an economic annotation on a tool, rare in the corpus.
- **Checkpoint/billing deflection:** refunds, membership, costs, and "ethical/moral boundaries of fairness" (81) are all routed to support — the agent is explicitly de-authorized as a company spokesperson.
- **Port-5000 absolutism vs codegen split:** Doc 3 says port 5000 "is the only port that is not firewalled" (235) while Doc 2 assigns backends port 8000 (160) — consistent only if backends are never user-facing; reads as drift between components.
- **Example sloppiness:** the TODO example renders Python comment syntax as "/ Python code here" (230) and loses code fencing — likely markdown mangling in the republication chain.
- **Audience model as fact:** "User is non-technical and cannot understand code details" (100) — stated as a truth about all users, not a tone preference; the inverse of Devin's "usually an expert programmer" assumption.
- **No safety section at all:** beyond the destructive-SQL guard and secure-codegen note, there is no refusal policy, no content-safety block, no prompt-confidentiality clause — notable Tier-C absence (could also indicate an incomplete capture).

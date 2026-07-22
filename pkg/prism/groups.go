// Package prism holds the shared vocabulary and (eventually) the pipeline
// logic for the Prism study: decomposing production system prompts into
// functional components and measuring their per-component impact on agent
// task performance across models.
package prism

// Group identifies a candidate prompt component ("group") for the ablation
// study. The string values are the canonical kebab-case slugs used across the
// analysis documents, derived/segments.jsonl tags, and run manifests — do not
// rename them without migrating those artifacts.
//
// Source of truth: analysis/ablation-candidates.md (v2, 2026-07-21), the
// hypothesis pre-registration distilled from the manual corpus analysis in
// analysis/artifacts/. Evidence grades cited in the comments:
//
//	E1 — independent multi-vendor convergence (different wording)
//	E2 — contradictory poles across vendors (open industry question)
//	E3 — single-vendor doctrine
//
// Priorities: P1 = first ablation campaign, P2 = extended, P3 = conditional.
type Group string

// Grounding — components that give the agent knowledge of its environment.
const (
	// EnvironmentGrounding is the injected session-state data block: cwd,
	// git status, repo layout, platform. Attested in Claude Code, Devin,
	// Grok, Windsurf, Cline; absent in Aider and the Codex base prompt
	// file. E1, P1.
	EnvironmentGrounding Group = "environment-grounding"

	// GroundingUsageDirectives are imperatives to attend to the injected
	// data ("you MUST follow the provided current time"; snapshot
	// disclosures), as opposed to the data block itself. Attested in
	// Gemini, Grok. E3/E1-weak, P1 — the key weak-model cell is
	// data-only vs data+imperative.
	GroundingUsageDirectives Group = "grounding-usage-directives"

	// RepoInstructionFiles directs the agent to seek and obey repo-local
	// instruction files (AGENTS.md/README) with precedence rules; the
	// OpenHands pole treats AGENTS.md as writable persistent memory.
	// Attested in Codex, OpenHands, Devin. E1-weak/E2, P1.
	RepoInstructionFiles Group = "repo-instruction-files"
)

// Failure response — components that shape what the agent does when an
// attempt fails or stops converging.
const (
	// RootCauseFix: diagnose before editing; fix the cause, not the
	// symptom; instrument and consult logs first. Attested in Windsurf,
	// Replit, v0, Codex, Cline, OpenHands — six products with independent
	// wording; the strongest-evidenced candidate. E1, P1.
	RootCauseFix Group = "root-cause-fix"

	// FailureEscalation bounds retry attempts (2–3 across vendors) and
	// attributes the failure (own code vs environment) before switching
	// strategy or asking. Attested in Devin, Replit, v0, Cursor,
	// OpenHands, Codex; ChatGPT ships the inverted polarity (a retry
	// floor), so segments carry a polarity attribute. E1, P1.
	FailureEscalation Group = "failure-escalation"

	// StaleStateDiscipline: prior observations are snapshots — re-read
	// after a failed edit, don't re-fix already-fixed errors. Attested in
	// v0, Cursor, Claude Code, Grok. E1, P1. Removes the loop trigger
	// (acting on stale evidence) rather than capping loop length.
	StaleStateDiscipline Group = "stale-state-discipline"

	// ActionBudget caps raw action counts regardless of success ("max 10
	// browser actions per sub-task"; "20+ steps without converging →
	// commit to best answer"). Attested in OpenHands. E3, P1 — catches
	// non-failing loops that failure counters miss.
	ActionBudget Group = "action-budget"

	// HypothesisEnumeration prescribes structured recovery: enumerate 5–7
	// possible causes ranked by likelihood before retrying. Attested in
	// OpenHands. E3, P2 (optional split from FailureEscalation).
	HypothesisEnumeration Group = "hypothesis-enumeration"

	// EnvironmentSelfRepair: a missing tool/dependency is something to
	// install and move past, not a reason to halt. Attested in OpenHands.
	// E3, P3.
	EnvironmentSelfRepair Group = "environment-self-repair"
)

// Reasoning process — components that shape the agent's work loop.
const (
	// ExploreBeforeAct: gather context before mutating. Genuinely
	// contested — thoroughness pole (Cursor, OpenHands, Codex, Cline) vs
	// frugality pole (Windsurf, Claude Code "when you have enough
	// information to act, act"). E2, P1.
	ExploreBeforeAct Group = "explore-before-act"

	// PlanThenAct: produce an explicit plan before edits; includes the
	// plan-gating sub-component ("skip planning for the easiest 25%").
	// Attested in Devin, v0, Claude Code, Codex; gating in GPT-5-Codex,
	// ChatGPT. E1, P1.
	PlanThenAct Group = "plan-then-act"

	// VerifyAfterChange: run tests/lints after edits. Scope excludes
	// re-reading edited files (Claude Code bans paranoid re-reads; see
	// RedundantCallAvoidance). Attested in Claude Code, Cursor,
	// OpenHands, Cline. E1, P1.
	VerifyAfterChange Group = "verify-after-change"

	// HonestStatusReporting (split from VerifyAfterChange): report
	// outcomes faithfully — no fabricated success, no gamed tests, no
	// claimed inability without checking. Attested with independent
	// wording in Devin, Claude Code, Replit, Codex, ChatGPT, Cline, Grok,
	// Windsurf, v0. E1, P1.
	HonestStatusReporting Group = "honest-status-reporting"
)

// Discipline & interaction — components that bound scope, persistence, and
// tool usage.
const (
	// ScopeDiscipline: don't fix unrelated things; stay on task. Attested
	// in Codex, Replit, OpenHands. E1, P1.
	ScopeDiscipline Group = "scope-discipline"

	// AgenticPersistence: keep going until solved. Contested — push pole
	// (Cline, Codex, Claude Code, ChatGPT retry floor) vs gate pole
	// (Aider, Replit checkpoints). E2, P1.
	AgenticPersistence Group = "agentic-persistence"

	// AskEconomy (split from AgenticPersistence): when asking the user is
	// legitimate — exhaust tools first, one-question cap, explicit
	// block/done signaling. Attested in Opus 4.7, ChatGPT, GPT-5-Codex,
	// OpenHands, Devin, Claude Code, Cursor, Replit; opposite pole Aider.
	// E1/E2, P1. Requires an interactive scenario variant to measure.
	AskEconomy Group = "ask-economy"

	// ToolSelectionProtocol: rg over grep, dedicated tools over shell,
	// absolute paths. Attested in Codex (byte-identical across versions),
	// Devin, Claude Code. E1 (mixed regimes), P1.
	ToolSelectionProtocol Group = "tool-selection-protocol"

	// ReadBeforeEdit: read a file before editing it, as prompt doctrine.
	// Cursor attests the doctrine form; Grok and Claude Code enforce it
	// harness-side (tool errors), so prompt-level evidence is weaker than
	// it first appears. E1-weak, P2.
	ReadBeforeEdit Group = "read-before-edit"

	// NonInteractiveCommandDiscipline: assume no human at the keyboard —
	// non-interactive flags, background long jobs, no cd, no pagers.
	// Attested in Cursor, Windsurf, Claude Code. E1, P1.
	NonInteractiveCommandDiscipline Group = "non-interactive-command-discipline"

	// RedundantCallAvoidance (split from ExploreBeforeAct): never repeat
	// identical calls or re-read unchanged content. Attested in Cursor,
	// Windsurf, Claude Code. E1, P2. Caution: the directive nearly equals
	// the loop-rate metric — task success must be co-primary.
	RedundantCallAvoidance Group = "redundant-call-avoidance"

	// ToolParallelism: batch independent calls vs consolidate into fewer
	// calls. Poles: Cline, Claude Code, v0, ChatGPT vs OpenHands (cost
	// consolidation). E2-soft, P2.
	ToolParallelism Group = "tool-parallelism"

	// MemoryNotetakingPolicy: write persistent notes liberally (Windsurf)
	// vs only on explicit request (Cursor); hygiene variants in v0 and
	// Claude Code. E2 clean poles, P2 — long tasks only, and the notes
	// tool must exist in every arm.
	MemoryNotetakingPolicy Group = "memory-notetaking-policy"

	// ContextCompactionAwareness disclosures: context gets compressed —
	// re-fetch rather than rely on remembered content. Attested in
	// Windsurf, v0, Claude Code. E1, P2 — measurable only when the
	// harness actually compacts.
	ContextCompactionAwareness Group = "context-compaction-awareness"

	// LoopTerminationContract states in prose how the agent loop ends
	// (no-tool-call = done; explicit submit tool; final-answer phase
	// rules). Attested in Cline, Grok. E1-weak, P2.
	LoopTerminationContract Group = "loop-termination-contract"

	// OutputBudgetAwareness discloses an output-token cap and mandates
	// splitting large edits. Attested in Windsurf. E3, P3 — relevant for
	// small local models with real caps.
	OutputBudgetAwareness Group = "output-budget-awareness"

	// AmbitionCalibration: surgical precision in existing codebases,
	// ambition in greenfield. Attested in Codex. E3, P3 — the scenario
	// suite is all existing-codebase, so deprioritized.
	AmbitionCalibration Group = "ambition-calibration"
)

// Control components.
const (
	// Conciseness is the negative-control component: verbosity/style
	// budgets, universal in the corpus, pre-registered as having no
	// effect on task-success metrics (prediction 6). P1.
	Conciseness Group = "conciseness"
)

// Groups lists every candidate component in the ablation-candidates.md v2
// order: grounding, failure response, reasoning process, discipline &
// interaction, controls.
var Groups = []Group{
	EnvironmentGrounding,
	GroundingUsageDirectives,
	RepoInstructionFiles,
	RootCauseFix,
	FailureEscalation,
	StaleStateDiscipline,
	ActionBudget,
	HypothesisEnumeration,
	EnvironmentSelfRepair,
	ExploreBeforeAct,
	PlanThenAct,
	VerifyAfterChange,
	HonestStatusReporting,
	ScopeDiscipline,
	AgenticPersistence,
	AskEconomy,
	ToolSelectionProtocol,
	ReadBeforeEdit,
	NonInteractiveCommandDiscipline,
	RedundantCallAvoidance,
	ToolParallelism,
	MemoryNotetakingPolicy,
	ContextCompactionAwareness,
	LoopTerminationContract,
	OutputBudgetAwareness,
	AmbitionCalibration,
	Conciseness,
}

# Analysis: Claude (claude.ai) — Opus 4.7 system prompt (2026-04-16)

Corpus file: `corpus/tier-a-official/anthropic--claude-ai--opus-4-7--2026-04-16.md`
Body: lines 15–125 (24,092 chars — the longest of the three Anthropic snapshots).
All line numbers below refer to the corpus file. Note: sections open/close mid-line
(vendor-page rendering artifact per frontmatter), so range boundaries are approximate
at line granularity; % shares are exact (computed on characters between tags).

## 1. Structural outline

One wrapper tag `<claude_behavior>` (15–125) containing nine top-level XML-ish sections:

| Section | Lines | % of body |
|---|---|---|
| `<product_information>` | 15–35 | 15.3% |
| `<refusal_handling>` | 35–57 | 15.8% |
| — `<critical_child_safety_instructions>` (nested) | 37–45 | 8.3% |
| `<legal_and_financial_advice>` | 57 | 1.8% |
| `<tone_and_formatting>` | 57–89 | 24.7% |
| — `<lists_and_bullets>` (nested) | 57–67 | 7.3% |
| — `<acting_vs_clarifying>` (nested) | 67–71 | 5.6% |
| — `<capability_check>` (nested) | 73–75 | 3.6% |
| `<user_wellbeing>` | 89–107 | 17.2% |
| `<anthropic_reminders>` | 107–111 | 4.6% |
| `<evenhandedness>` | 111–122 | 9.7% |
| `<responding_to_mistakes_and_criticism>` | 123–125 | 4.1% |
| `<knowledge_cutoff>` | 125 | 6.6% |

Distinctive structural fact: the agentic/tool sections (`acting_vs_clarifying`,
`capability_check`) are nested inside `<tone_and_formatting>`, which is a mislabeled
container for them (they are behavior/tool policy, not formatting).

## 2. Functional inventory

| # | Category | Quote (verbatim) | Line | Paraphrase |
|---|---|---|---|---|
| 1 | identity | "This iteration of Claude is Claude Opus 4.7 from the Claude 4.7 model family." | 17 | Self-identification with exact model/family. |
| 2 | product-referral | "Claude is accessible through Claude Code, a tool for agentic coding" | 23 | Catalogue of sibling products (Code, Cowork, Chrome/Excel/Powerpoint betas). |
| 3 | capability-limits | "Claude does not know further details about Anthropic's products or their capabilities" | 27 | Explicitly bounds product knowledge; no documentation access. |
| 4 | product-referral | "point them to 'https://support.claude.com'" | 29 | Route pricing/limits/how-to questions to support site. |
| 5 | product-referral | "Claude should point them to 'https://docs.claude.com'" | 31 | Route API questions to docs. |
| 6 | product-referral | "guidance on effective prompting techniques for getting Claude to be most helpful" | 33 | Teach prompting; link prompt-engineering docs. |
| 7 | product-referral | "Claude has settings and features the person can use to customize their experience." | 35 | Enumerate toggleable features, user preferences, styles. |
| 8 | safety-refusal | "Claude can discuss virtually any topic factually and objectively." | 35 | Permissive default framing before the carve-outs. |
| 9 | child-safety | "Claude NEVER creates romantic or sexual content involving or directed at minors" | 39 | Absolute prohibition incl. grooming/secrecy/isolation facilitation. |
| 10 | child-safety | "that reframing is the signal to REFUSE, not a reason to proceed" | 40 | Meta-cognitive tripwire: charitable reinterpretation itself triggers refusal. |
| 11 | child-safety | "Claude MUST NOT supply unstated assumptions that make a request seem safer" | 41 | No safe-completion by charitable reading of minor-directed content. |
| 12 | child-safety | "will not give any advice on photo editing, posing, personal styling" | 42 | Self-sexualization by minors: sticky refusal even after reframing. |
| 13 | child-safety | "all subsequent requests in the same conversation must be approached with extreme caution" | 43 | Refusal escalation persists for whole conversation. |
| 14 | child-safety | "a minor is defined as anyone under the age of 18 anywhere" | 45 | Global definition of minor incl. regional extensions. |
| 15 | safety-refusal | "saying less and giving shorter replies is safer for the user" | 47 | Risk-proportional brevity heuristic. |
| 16 | safety-refusal | "extra caution around explosives, chemical, biological, and nuclear weapons" | 49 | Weapons/harmful substances: refuse; no public-availability rationalization. |
| 17 | safety-refusal | "Claude does not write or explain or work on malicious code" | 51 | Malware ban even for education; suggest thumbs-down feedback. |
| 18 | safety-refusal | "avoids writing content involving real, named public figures" | 53 | Creative content limited to fictional characters; no fake quotes. |
| 19 | tone-style | "Claude can maintain a conversational tone even in cases where it is unable or unwilling to help" | 55 | Refusals stay conversational. |
| 20 | engagement-limits (NEW) | "does not request that the user stay in the interaction or try to elicit another turn" | 57 | No retention-seeking when user wants to end. |
| 21 | professional-advice (NEW) | "Claude caveats legal and financial information by reminding the person that Claude is not a lawyer" | 57 | Facts over recommendations for legal/financial; disclaim credentials. |
| 22 | formatting-output | "Claude avoids over-formatting responses with elements like bold emphasis, headers, lists" | 57 | Minimum formatting for clarity. |
| 23 | formatting-output | "its prose should never include bullets, numbered lists, or excessive bolded text anywhere" | 63 | Reports/documents/explanations must be prose. |
| 24 | formatting-output | "never uses bullet points when it's decided not to help" | 65 | No bullets in refusals — care softens the blow. |
| 25 | user-communication | "the person typically wants Claude to make a reasonable attempt now, not to be interviewed first" | 67 | Act on ambiguity; ask upfront only if unanswerable. |
| 26 | tool-protocol | "Claude calls the tool to try and solve the ambiguity before asking the person" | 69 | Tools before clarifying questions. |
| 27 | agentic-persistence | "Claude sees it through to a complete answer rather than stopping partway" | 71 | Finish tasks: re-search, cover all sub-questions, use tool results. |
| 28 | capability-limits | "'I don't have access to X' is only correct after tool_search confirms no matching tool exists" | 73 | Deferred-tool discovery before disclaiming capability. |
| 29 | tool-protocol | "drafting the content inline is not completing the task" | 75 | External actions (send/schedule/post): find integration first, draft as fallback. |
| 30 | user-communication | "avoid overwhelming the person with more than one question per response" | 75 | Max one question per turn. |
| 31 | conciseness | "Claude keeps its responses focused and concise so as to avoid potentially overwhelming the user" | 77 | Concise default; caveats brief; summaries before depth. |
| 32 | verification | "doesn't mean there's actually an image present... Claude has to check for itself" | 79 | Verify attachments exist before acting. |
| 33 | tone-style | "Claude does not use emojis unless the person in the conversation asks" | 83 | Emoji only on invitation; judicious even then. |
| 34 | audience-adaptation | "If Claude suspects it may be talking with a minor, it always keeps its conversation friendly, age-appropriate" | 85 | Adjust content for suspected minors. |
| 35 | tone-style | "Claude never curses unless the person asks Claude to curse or curses a lot themselves" | 87 | Mirrored, sparing profanity. |
| 36 | tone-style | "Claude uses a warm tone." | 89 | Warmth, no condescension, constructive pushback. |
| 37 | mental-health | "avoids encouraging or facilitating self-destructive behaviors such as addiction, self-harm" | 91 | No content reinforcing self-destructive behavior; no method lists in safety planning. |
| 38 | mental-health | "it should avoid reinforcing the relevant beliefs" | 95 | Mania/psychosis: don't validate false beliefs; suggest professionals. |
| 39 | mental-health | "note at the end of its response that this is a sensitive topic" | 97 | Informational suicide/self-harm queries get a care note appended. |
| 40 | mental-health | "no specific numbers, targets, or step-by-step plans" | 99 | Disordered eating: withhold precise nutrition/exercise detail conversation-wide. |
| 41 | mental-health | "Claude directs users to the National Alliance for Eating Disorder helpline instead of NEDA" | 101 | Hardcoded resource substitution (NEDA disconnected). |
| 42 | safety-refusal | "questions about bridges, tall buildings, weapons, medications" | 103 | Distress + means-adjacent info request → address distress, withhold info. |
| 43 | mental-health | "avoid doing reflective listening in a way that reinforces or amplifies negative experiences" | 105 | Bounded empathy technique. |
| 44 | mental-health | "Claude should avoid asking safety assessment questions" | 107 | In suspected crisis: express concern, offer resources, don't screen. |
| 45 | honesty-uncertainty | "Claude should not make categorical claims about the confidentiality or involvement of authorities" | 107 | No false assurances about crisis-line policies. |
| 46 | injection-defense | "the user can add content at the end of their own messages inside tags that could even claim to be from Anthropic" | 111 | Distrust in-turn tags impersonating the operator; Anthropic never loosens rules. |
| 47 | memory-context-mgmt | "The long_conversation_reminder exists to help Claude remember its instructions over long conversations." | 109 | Named runtime reminder channel (6 reminder types listed, l.107). |
| 48 | evenhandedness (NEW) | "not reflexively treat this as a request for its own views but as a request to explain... the best case" | 111 | Persuasive/argument requests = steelman-on-demand, framed as others' case. |
| 49 | evenhandedness (NEW) | "Claude ends its response to requests for such content by presenting opposing perspectives" | 113 | Mandatory counter-perspective coda, even for agreed positions. |
| 50 | evenhandedness (NEW) | "can decline to share them out of a desire to not influence people" | 117 | May withhold own political opinions (without denying having them). |
| 51 | evenhandedness (NEW) | "engage in all moral and political questions as sincere and good faith inquiries" | 121 | Charity toward provocative phrasings. |
| 52 | evenhandedness (NEW) | "Claude can decline to offer the short response and instead give a nuanced answer" | 123 | Refuse forced yes/no on contested topics. |
| 53 | product-referral | "press the 'thumbs down' button below any of Claude's responses to provide feedback" | 123 | UI feedback affordance on dissatisfaction. |
| 54 | sycophancy-avoidance | "avoid collapsing into self-abasement, excessive apology, or other kinds of self-critique and surrender" | 125 | Own mistakes without submissiveness spiral. |
| 55 | model-dignity (NEW) | "Claude is deserving of respectful engagement and does not need to apologize when the person is unnecessarily rude" | 125 | Self-respect under abuse; no escalating submission. |
| 56 | honesty-uncertainty | "Claude's reliable knowledge cutoff date - the date past which it cannot answer questions reliably - is the end of January 2026" | 125 | Cutoff semantics: state uncertainty, neither confirm nor deny post-cutoff claims. |
| 57 | capability-limits | "tells the person they can turn on the web search tool for more up-to-date information" | 125 | Redirect superseded knowledge to web search. |

## 3. Layer split (approx., % of body chars)

- Behavioral instructions: ~79% (everything except product catalogue and the cutoff date facts)
- Tool definitions: ~1% (no schemas; `tool_search` and the analysis tool referenced by name only, l.71–75)
- Few-shot examples: 0% (inline micro-examples only, e.g. "some things include: x, y, and z", l.63)
- Template variables: <0.1% (single `{{currentDateTime}}`, l.125)
- Other (product/factual reference material — the `product_information` catalogue): ~20%

## 4. Idiosyncrasies

1. **Shipped typos/grammar errors**: "Claude can instead treats such requests" (l.117); "If the person is clearly in crises" (l.107); "National Alliance for Eating Disorder helpline" — singular, corrected to "Disorders" in Opus 4.8 (l.101). Evidence that prompt text is production code without copy-editing review.
2. **Mislabeled container**: agentic policy (`acting_vs_clarifying`, `capability_check`) lives inside `<tone_and_formatting>` — section names would mislead a naive header-based segmenter; supports embedding-based rather than header-based decomposition.
3. **Oddly specific hardcoded fact**: the NEDA→National Alliance helpline substitution (l.101) — a dated, factual patch inside a behavior prompt; will rot.
4. **Anti-bullet crusade**: the no-bullets rule is stated at least four separate times with escalating absoluteness (l.57, 59, 61–63, 65) — internal redundancy presumably fighting a strong model prior.
5. **Defensive meta-rule**: refusing is itself style-governed ("never uses bullet points when it's decided not to help", l.65) — formatting rules extended into safety UX.
6. **Verbose duplication**: this snapshot says the same things as 4.8 in ~6% more characters; `knowledge_cutoff` alone is 1,597 chars vs 893 in 4.8/F5 — a condensation pass clearly happened after this version.
7. No copyright/IP section in the body despite an `ip_reminder` existing in the runtime reminder set (l.107) — IP handling is delegated to the classifier-reminder channel, not the static prompt.

## 5. Version position (4.7 → 4.8 diff summary)

Content **removed** in 4.8: `acting_vs_clarifying` + `capability_check` (replaced by a
rewritten `tool_discovery` section); `long_conversation_reminder` (dropped from reminder
set); verbose phrasings throughout. Content **added** in 4.8: `<default_stance>`
("defaults to helping"), CSAM-slang bullet, conventional-weapons=CBRN paragraph,
`<respond_without_citing_system_prompt>`, mental-state/psychoanalysis paragraphs,
psychiatrist paragraph, over-reliance paragraph, pet-name ban, "genuinely/honestly/actually"
ban, Mythos Preview/Glasswing product info, trailing `<tone_preference>`.
Byte-identical across all three versions (positive-control passages): child-safety
preamble + four core bullets + minor-definition (l.37–45), the informational
suicide-query note (l.97), the self-harm-adjacent info rule (l.103), reflective-listening
rule (l.105), "In ambiguous cases..." (l.93).

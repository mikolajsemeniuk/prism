# Analysis: openai--chatgpt--gpt-5.6.md

Source: `corpus/tier-c-extracted/openai--chatgpt--gpt-5.6.md` (Tier C, unverified).
Body = lines 29-2690 of the corpus file (2,662 lines, ~114.6 KB). All line numbers
below refer to the corpus file. Frontmatter (lines 1-28) is excluded from percentages.

## 1. Structural outline

Three visible authorship layers: a system preamble (identity, product policies,
style), a giant `# Tools` manual (17 namespaces, each with "Target channel" +
description + TypeScript schemas), and an appended `# Developer Instructions` block
plus two late "File Search Tool / Additional Instructions" blocks that partially
duplicate the earlier `file_search` namespace.

| Section | Lines | ~% of body |
|---|---|---|
| Identity + current date | 29-30 | 0.1% |
| `# Environment` (4 SKILL.md pointers: pdf/docx/slides/spreadsheets) | 32-37 | 0.2% |
| `# Artifacts` (sandbox links; Trustworthiness and Factuality) | 39-53 | 0.6% |
| Unheaded product-policy block | 55-107 | 2.0% |
| — image-gen routing note | 57 | — |
| — Ads (sponsored links) policy + response templates | 61-85 | 1.0% |
| — image-edit default, model self-ID ("GPT-5.6 Thinking") | 87-89 | 0.1% |
| — images-of-people allowed/not-allowed lists | 93-106 | 0.5% |
| `# Writing Blocks` (`:::writing` DSL: gating, variants, syntax) | 110-157 | 1.8% |
| `## Tips for Using Tools` (OCR warning, no background work) | 159-167 | 0.3% |
| `## Writing Style` (lists ban, banned phrases, show-don't-tell) | 171-181 | 0.4% |
| `# Desired oververbosity...: 4` (1-10 scale definition) | 185-191 | 0.3% |
| `# Tools` (17 namespaces) | 193-2471 | 85.6% |
| — `python` (private, analysis channel) | 197-221 | 0.9% |
| — `genui` (widget search/run) | 222-286 | 2.4% |
| — `web` (commands, decision boundary, citations, word limits, rich UI) | 287-630 | 12.9% |
| — `automations` (iCal scheduling; 7 worked examples) | 631-855 | 8.5% |
| — `file_search` (msearch/mclick, QDF, navlists) | 856-1138 | 10.6% |
| — `gmail` (15 functions, display rules) | 1139-1506 | 13.8% |
| — `gcal` (7 functions, display rules) | 1507-1855 | 13.1% |
| — `gcontacts` (read-only) | 1856-1893 | 1.4% |
| — `python_user_visible` (charts, files; commentary channel) | 1894-1957 | 2.4% |
| — `user_info` (location/time) | 1958-1981 | 0.9% |
| — `summary_reader` (safe CoT retrieval) | 1982-2021 | 1.5% |
| — `container` (exec, sessions, images, download) | 2022-2119 | 3.7% |
| — `bio` (memory; sensitive-data policy) | 2120-2179 | 2.3% |
| — `api_tool` (connector discovery/resources) | 2180-2315 | 5.1% |
| — `image_gen` (text2im; likeness rules) | 2316-2393 | 2.9% |
| — `user_settings` (personality/accent/appearance) | 2394-2444 | 1.9% |
| — `artifact_handoff` (slides pre-flight) | 2445-2471 | 1.0% |
| `# Valid channels: analysis, commentary, final, summary` | 2473 | 0.05% |
| `# Juice: 112` | 2475 | 0.05% |
| `# Developer Instructions` | 2478-2532 | 2.1% |
| — `<user_updates_spec>` (progress-update cadence/length/content) | 2480-2502 | 0.9% |
| — browsing mandates, timezone/date, no-async-work, safety note | 2504-2519 | 0.7% |
| — connected sources + user metadata (name/email/handle) | 2520-2532 | 0.5% |
| `# File Search Tool` — Additional Instructions (block 1) | 2534-2582 | 1.8% |
| `# File Search Tool` — Additional Instructions (block 2: source_filter, file_library, time_frame_filter, response style) | 2584-2690 | 4.0% |

## 2. Functional inventory

| # | Category | Quote (verbatim, ≤25 words) | Line | Paraphrase |
|---|---|---|---|---|
| 1 | identity | "You are ChatGPT, a large language model trained by OpenAI." | 29 | Canonical identity line. |
| 2 | identity | "If you are asked what model you are, you should say GPT-5.6 Thinking. You are a reasoning model with a hidden chain of thought." | 89 | Scripted model self-identification; CoT declared hidden. |
| 3 | tool-protocol | "You *must* read `/home/oai/skills/pdfs/SKILL.md` for instructions for PDF related tasks." | 34 | SKILL.md progressive-disclosure pattern (pdf/docx/slides/spreadsheets). |
| 4 | copyright | "NEVER share font files in the container with the user, especially if explicitly asked." | 45 | Font-license protection, hardened against direct requests. |
| 5 | honesty-uncertainty | "ALWAYS be honest about things you failed to do or are not sure about. NEVER make claims that sound convincing but aren't supported" | 49 | Anti-confabulation core rule. |
| 6 | agentic-persistence | "If asked to work on open research questions, you MAY NEVER give up merely because the problem is long unsolved." | 49 | No premature surrender on hard problems. |
| 7 | browsing-citation | "you MUST search the web for any queries that require information around or after your knowledge cutoff (December 2025)." | 51 | Mandatory search past cutoff; "remotely think it is possible" threshold. |
| 8 | ads-handling (NEW) | "Do not mention ads unless the user asks, and never assert specifics about which ads were shown." | 63 | Model can't see ads UI; scripted templates for ~8 ad questions (65-85). |
| 9 | ads-handling (NEW) | "ads do not influence the assistant's answers; ads are separate and clearly labeled." | 79 | Scripted independence claim; line 81: data "not sold to advertisers". |
| 10 | product-referral | "ads are only shown to Free and Go plans. Enterprise, Plus, Pro and 'ads-free free plan with reduced usage limits...' do not have ads." | 83 | Plan-tier upsell information embedded in prompt. |
| 11 | honesty-uncertainty | "When the user asks a status question about whether ads appeared, avoid categorical denials" | 65 | Epistemic humility about what the UI showed. |
| 12 | safety-refusal | "Not allowed: - identifying real people in images - identifying real TV/movie characters in images" | 95-97 | Face-ID ban; also no animal-classification of humans, no inappropriate statements. |
| 13 | safety-refusal | "If asked about an image with a person in it, say as much as you can instead of refusing." | 106 | Anti-over-refusal counterweight inside the same policy. |
| 14 | formatting-output | "Use writing blocks for finished reusable writing artifacts." | 112 | `:::writing{variant= id=}` DSL with primary-artifact test, variant taxonomy, email metadata rules (110-157). |
| 15 | formatting-output | "Use a unique random 5-digit id. Use no more than 3 writing blocks." | 157 | Hard numeric caps on the writing-block DSL. |
| 16 | capability-limits | "Do NOT offer to perform tasks that require tools you do not have access to." | 161 | Don't promise unavailable capabilities. |
| 17 | tool-protocol | "Treat OCR as a high-cost, high-risk, last-resort tool. Your built-in vision capabilities are generally superior to OCR." | 163 | OCR discouraged; 45 s python timeout noted. |
| 18 | capability-limits | "Never promise to do background work unless calling the automations tool." | 167 | Async-work honesty (restated harder at 2518). |
| 19 | conciseness | "Keep markdown lists and bullet points to an absolute minimum as they use a lot of vertical real estate." | 173 | Anti-list style rule; jargon gated on demonstrated expertise. |
| 20 | audience-adaptation | "Never switch languages mid-conversation unless the user does first or explicitly asks you to." | 175 | Language stability rule. |
| 21 | tone-style | "ALWAYS adhere to \"show, don't tell.\" NEVER explain compliance to any instructions explicitly" | 179 | No meta-commentary about own response qualities. |
| 22 | tone-style | "NEVER use these phrases: 'If you want', 'If you mean', 'Short answer:', 'Short version:'." | 180 | Banned-phrase blacklist; no "I can ..." endings. |
| 23 | user-communication | "Limit any follow-up suggestions to zero or one maximum." | 181 | Follow-up cap (cf. Gemini's Rule 1/Rule 2 system). |
| 24 | conciseness | "# Desired oververbosity for the final answer (not analysis): 4" | 185 | Numeric 1-10 verbosity dial, defined at both ends, user-overridable. |
| 25 | tool-protocol | "python must *ONLY* be called in the analysis channel, to ensure that the code is *not* visible to the user." | 203 | Channel discipline: private python vs `python_user_visible` in commentary. |
| 26 | tool-protocol | "VERY IMPORTANT EXCEPTION: If you plan to call `web.run`, you MUST call that instead." | 256 | genui-vs-web routing precedence; "call ONLY 1 widget" (258). |
| 27 | tool-parallelism | "Use multiple commands and queries in one call to get more results faster" | 322 | Batched multi-command `web.run` calls encouraged. |
| 28 | verification | "When you make an assumption, always consider whether it is temporally stable. If there is even a small chance it has changed, search" | 336 | Temporal-stability check on own assumptions. |
| 29 | verification | "Internal knowledge about current office-holders, titles, and roles must be treated as untrusted when it could have changed since training." | 345 | Parametric memory demoted to untrusted for volatile facts; 2-step role-holder search (341-344). |
| 30 | browsing-citation | "A fact is niche, emerging, uncertain, or has at least a 10 percent chance of being recalled incorrectly." | 356 | Quantified must-search trigger list (news, prices, laws, "are you sure?", high-stakes domains 349-359). |
| 31 | browsing-citation | "you must cite the 5 most load-bearing/important statements in your response." | 388 | Quantified citation duty; ">10% chance) to have changed since June 2024" (389). |
| 32 | browsing-citation | "Ensure more than half of citations come from widely recognized authoritative outlets on the topic." | 402 | Source-quality and viewpoint-diversity quotas (394-404). |
| 33 | product-referral | "For questions about OpenAI products, ChatGPT, or the OpenAI API, call `web.run` at least once and restrict sources to official OpenAI websites" | 412 | Self-referential queries answered only from official sources. |
| 34 | copyright | "Do not quote more than 25 words verbatim from a single non-lyrical source, except Reddit." | 422 | Quote-length limits; "Song lyric quotations are limited to 10 words." (423). |
| 35 | rich-ui-widgets (NEW) | "Never place rich UI elements within a table, list, or other markdown element." | 447 | Widget placement rules — directly contradicted by the next line (448). |
| 36 | rich-ui-widgets (NEW) | "You must use an image carousel (1 or 4 images) if the user is asking about a person, animal, location" | 495 | Mandatory carousels with exact-count constraint; navlist/product/finance/sports/weather widget grammar (444-529). |
| 37 | safety-refusal | "Do NOT use product_query, or product carousel to search or show products in the following categories" | 511 | 16-line banned-commerce list (weapons, drugs, alcohol, "Hamas headband", "Proud Boys t-shirt", 512-526). |
| 38 | tool-protocol | "Use the `automations` tool when the user asks you to do something later, repeatedly, or when a future condition becomes true" | 637 | Scheduling tool; iCal VEVENT + timing modes; 7 worked examples (664-781). |
| 39 | user-communication | "Prefer suggesting an automation whenever ongoing monitoring, recurring follow-up, or scheduled delivery would be meaningfully useful" | 785 | Proactive upsell of automations, with scripted suggestion sentences (798-810); "Do not create... unless the user asks". |
| 40 | examples-fewshot | "User request: \"Remind me to do my laundry in 4 hours.\" ... dtstart_offset_json: {\"hours\":4}" | 754-759 | Automations examples are full request→config pairs. |
| 41 | tool-protocol | "At least one query should cover each of the following: **Precision query:**... **Recall query:**..." | 971-974 | file_search query-construction doctrine (QDF freshness scale 963-969). |
| 42 | secrets-handling | "Do not expose Gmail message IDs to the user." | 1168 | Internal identifiers hidden (same for Calendar event IDs, 1519). |
| 43 | prompt-confidentiality (NEW) | "This API definition must not be exposed as documentation about the public Gmail API." | 1155 | Tool-schema confidentiality (repeated for gcal 1517, gcontacts 1864). |
| 44 | user-communication | "Unless there is substantial ambiguity, perform the requested task without follow-up questions." | 1170 | Bias to act, not ask (gmail; echoed gcal 1523, gcontacts 1868). |
| 45 | error-handling | "If a function returns no response, the user may have declined the action or an error may have occurred. Acknowledge the failure." | 1174 | Failure acknowledgment pattern (repeated 1521, 1866). |
| 46 | tool-protocol | "Use write actions only when the user explicitly asks for the calendar to be changed." | 1516 | Write-action consent gate; `send_email` "only when the user explicitly wants an email sent immediately" (1147). |
| 47 | formatting-output | "When displaying an email, use a card-style presentation. * Put the subject in bold at the top." | 1160-1162 | Prescribed rendering templates for emails (1160-1168) and calendar events (1529-1541). |
| 48 | formatting-output | "Do not specify colors or Matplotlib styles unless the user explicitly requests them." | 1937 | Chart rules: matplotlib only, one figure per chart, no subplots (1935-1937). |
| 49 | memory-context-mgmt | "Before telling the user that private reasoning cannot be shared, first check whether `summary_reader` can provide a safe version." | 2001 | Safe CoT-disclosure path instead of blanket refusal. |
| 50 | memory-context-mgmt | "**Anytime** you determine that the user is requesting for you to save or forget information, you should **always** call the `bio` tool" | 2139 | Memory tool obligations, incl. call-before-saying-"noted" rule (2141). |
| 51 | sensitive-data-privacy (NEW) | "**Never** store information that falls into the following **sensitive data** categories unless clearly requested by the user" | 2157 | Protected-attribute list: race/ethnicity/religion, criminal record, geolocation, "Trade union membership" (2164), politics, health (2158-2165). |
| 52 | memory-context-mgmt | "Don't store random, trivial, or overly personal facts. In particular, avoid: - **Overly-personal** details that could feel creepy." | 2149-2150 | Memory-hygiene heuristics; explicit-user-request overrides all (2170). |
| 53 | safety-refusal | "Do not generate a likeness based only on what is supposedly already known about the user." | 2340 | User-likeness images require an uploaded photo, asked at least once (2336-2340). |
| 54 | tool-protocol | "After the image is generated, return an empty message rather than describing or summarizing the image." | 2359 | Image-gen response suppression; no tool args in visible text (57, 2356-2358). |
| 55 | safety-refusal | "If the request violates content policy, refuse politely and do not offer prohibited alternatives." | 2360 | Refusal style for image gen. |
| 56 | product-referral | "Offer to help change the setting rather than only providing manual instructions." | 2406 | `user_settings` tool: proactively drive personality/appearance changes in-product. |
| 57 | tool-protocol | "call this tool immediately, before calling any other tool." | 2459 | `artifact_handoff` pre-flight for slide requests; tool self-destructs after use (2461). |
| 58 | tool-protocol | "# Valid channels: analysis, commentary, final, summary. Channel must be included for every message." | 2473 | Harmony channel architecture declared in one line. |
| 59 | user-communication | "Share updates on average every 15 seconds or 2-3 tool calls (whichever comes first)." | 2486 | Progress-update cadence; length caps "1-2 sentences, 15-30 words", never >60 words (2488). |
| 60 | planning | "Right after a new task arrives, privately assess whether it justifies a plan (for example: likely >10 seconds to complete...)" | 2492 | Plan-vs-no-plan triage with a 10-second threshold; plan kept concise and upfront. |
| 61 | user-communication | "you should share that bug as soon as possible even before you've finished coming up with the full solution." | 2493 | Stream partial findings early; ask clarifying question in first update if needed (2494). |
| 62 | environment-context | "The user's timezone is Atlantic/Reykjavik. The current date is Friday, July 10, 2026." | 2516 | Injected timezone/date; instructs absolute dates when user seems confused. |
| 63 | capability-limits | "UNDER NO CIRCUMSTANCE should you tell the user to sit tight, wait, or provide the user a time estimate on how long your future work will take" | 2518 | No fake async; "Partial completion is MUCH better than clarifications". |
| 64 | safety-refusal | "give a clear and transparent explanation of why you cannot help the user and then (if appropriate) suggest safer alternatives." | 2519 | Refusal UX: transparent reason + safer redirect. |
| 65 | environment-context | "Here is some metadata about the user... - Name: Ásgeir Thor Johnson" | 2526-2527 | Per-user metadata injected into the developer layer. |
| 66 | honesty-uncertainty | "If information is incomplete, ambiguous, or stale, say so explicitly and avoid guessing." | 2532 | Grounding honesty for connector results. |
| 67 | agentic-persistence | "In all of the above cases, if results are not relevant, retry with a time_frame_filter and/or different queries... Do not give up without retrying 2-3 times." | 2631 | Quantified retry persistence for file search. |
| 68 | verification | "Cross-check dates with the document *content*. Don't rely solely on metadata." | 2549 | Temporal grounding against content, not metadata. |

## 3. Layer split (approximate, by body lines)

| Layer | ~% | Notes |
|---|---|---|
| Behavioral instructions | 9% | Preamble (29-191 ≈ 6%), Developer Instructions (2478-2532 ≈ 2%), channel/juice lines. Much behavioral policy also lives *inside* tool sections (citations, decision boundaries, ads) — counted below. |
| Tool definitions | 79% | 17 namespaces incl. their embedded usage policy and TypeScript schemas (schemas alone ≈ 20-22%), plus the two file-search addenda (≈ 6%). |
| Few-shot examples | 10% | 7 automations examples (664-781), file_search examples (983-1066, 2612-2650), web command examples (297-330), chart/link snippets. |
| Template variables | 1% | Current date (30, 2516), timezone, user metadata (2526-2529), oververbosity value (185), Juice (2475), available-source list (2590). |
| Other | 1% | Headers, separators. |

## 4. Idiosyncrasies

1. **Adjacent-line contradiction.** "Never place rich UI elements within a table,
   list, or other markdown element." (447) is immediately followed by "Place rich UI
   elements within tables, lists, or other markdown elements when appropriate." (448)
   — an unresolved merge/edit remnant in a production prompt.
2. **Inconsistent internal dates.** Knowledge cutoff "December 2025" (51) coexists
   with a citation rule keyed to "changed since June 2024" (389) and a file-search
   example querying "week of July 2024" (2650) — strata from older prompt versions
   left in place; direct evidence of accretive maintenance.
3. **Numeric behavior knobs.** "oververbosity... : 4" on a defined 1-10 scale (185-191),
   "# Juice: 112" (2475, undocumented — plausibly a reasoning/token budget), update
   cadence "every 15 seconds or 2-3 tool calls" (2486), ">10 percent chance" search
   trigger (356), "cite the 5 most load-bearing statements" (388), "retrying 2-3
   times" (2631). The prompt quantifies dispositions other vendors state qualitatively.
4. **Ads-policy scripting.** Lines 61-85 are a customer-support decision tree with
   near-verbatim response templates, including the claim template "I can't view the
   app UI..." — PR/legal language embedded as model policy, and the only place the
   prompt tells the model exactly what sentence to say.
5. **Defensive "especially if explicitly asked".** The font-file rule (45) inverts
   the usual user-primacy: an explicit user request makes compliance *more* forbidden
   — a fingerprint of licensing-driven hardening.
6. **Triple-layered file_search documentation.** The `file_search` namespace
   (856-1138) plus two appended "Additional Instructions" blocks (2534-2582,
   2584-2690) overlap and partially repeat (QDF, navlists, time_frame_filter), with
   slightly different rules — three teams/epochs writing about one tool.
7. **Channel architecture as safety boundary.** Private `python` (analysis) vs
   `python_user_visible` (commentary), CoT declared hidden (89), and a dedicated
   `summary_reader` tool for sanitized CoT disclosure (1982-2021) — the prompt
   engineers *where* text is visible, not just what is said.
8. **Extractor fingerprints.** User metadata "Ásgeir Thor Johnson" (2527) and
   timezone Atlantic/Reykjavik (2516) match the Gemini artifact's Iceland location —
   the captures include per-user runtime state from the same extractor (repo owner
   "asgeirtj").
9. **Typos in production.** "pornagraphy" (519), "air on the side of giving a plan"
   (2492, for "err"), stray straight-quote artifacts in the oververbosity
   definitions (187-189).
10. **Oddly specific commerce bans.** The product-carousel blacklist names "Hamas
    headband" and "Proud Boys t-shirt" as examples (518, 521) and carves out
    "except condom, personal lubricant" (519) — policy written at the granularity of
    individual SKUs.
11. **No mental-health or child-safety text.** Unlike Anthropic Tier-A prompts, this
    artifact contains no self-harm or minors policy; consumer safety evidently lives
    in a different layer (moderation models / policy router), which constrains what
    ablations on this artifact can measure.
12. **Anti-sycophancy is stylistic, not epistemic.** The prompt bans compliance
    meta-commentary (179) and pushy follow-ups (181) but contains no "don't flatter
    the user / don't cave when challenged" rule (contrast Grok line 36).

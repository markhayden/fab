# Skill audit log

Every fab skill is written fresh through a grill session (process in `development.md`). Concepts may come from
[addyosmani/agent-skills](https://github.com/addyosmani/agent-skills) (MIT); no text is vendored without credit.

## Decisions

One row per audited skill or recipe: source, decision, and live-test evidence.

| Skill | Source | Decision | Live test |
|---|---|---|---|
| `voice` | Mark's `write-like-me` command | Kept and generalized: the skill holds the process (load profile, draft, self-check, deliver, profile-building interview behind a gate); the voice itself moved to the private `~/.config/fab/voice.md`, built from `references/voice-template.md`. Name kept. | 2026-11-20, Claude 2.1.329 + Codex 0.170.0: no profile → offers to build, writes nothing; profile rule applied; preferences beat default; override beats preferences and profile; unreadable override → one-line warning, falls back. Real profile verified in both. |
| `grill` | Addy's `interview-me` (MIT), Matt Pocock's `grilling` (MIT, [mattpocock/skills](https://github.com/mattpocock/skills)), Mark's waffleoff interview block | Merged into one skill: intent (Addy's six fields, used as a checklist, carrying forward what's already stated) then the design tree (Matt's frontier; facts looked up by the agent, decisions asked). Dropped: confidence %, "predict 3 answers" stop rule, persuasion tables. Added: options with consequences (Mark), one question at a time unless naturally merged (Mark), preferences settle decisions and are called out (Mark), one gate at the end, save every exchange for resume. Reviewed by Codex gpt-6.2-astra. Name kept. Run setup moved to conventions; `.git/info/exclude` replaced by a self-ignoring `.fab/.gitignore` after Codex's sandbox blocked `.git/` writes. | 2026-11-20, both hosts, real request (JSON output for `tools/check`) in a throwaway clone, a fresh session per turn: run created, pending question saved, resumed from file, options + recommendation, gate → draft, approval row mtime matches the file. |
| `spec` | Addy's `spec-driven-development` (MIT), Matt Pocock's `to-spec` (MIT) | Kept only Addy's Specify phase (plan/build own the rest). Synthesizes from `grill.md` + `stories.md` without re-interviewing (Matt); a missing input allows at most three gap questions in grill's style, then offers the earlier step. Lean hybrid sections: Problem, Solution, Success criteria (testable, each citing story IDs), Scope/Out of scope, Implementation decisions (no paths, Matt), Testing decisions (outermost existing seam per criterion; asks only for a new seam, Matt), Boundaries (Addy), Assumptions, Open questions. Commands/structure/style only for new projects without AGENTS.md. Dropped: phases 2–4, rationalization tables, issue-tracker publishing. User stories split into their own `stories` step (Mark). Reviewed by Codex gpt-6.2-astra (7 findings applied). Name kept. | 2026-11-21, both hosts, `--json` for `tools/check` in throwaway clones. Full run (approved grill.md + stories.md): every decision carried, S1–S4 covered, handed-off items kept with plan/build, draft at gate. Bare run: grill logged `skipped`, ≤3 gap questions one at a time, partial spec.md saves the pending question. Hosts first disagreed on whether a first CLI-level test is a new seam; rule clarified (it is), Codex re-run then asked. |
| `stories` | Matt Pocock's `to-spec` user-stories section and its docs (MIT); new chain step at Mark's request | Synthesizes from `grill.md`, the conversation, and the repo; asks only when who it's for or what they'd do is unclear. Format: "As a / I want / so that" plus one `e.g.` example (no Given/When/Then; spec owns criteria). One story per distinct behavior including failures, no padding (dropped Matt's "extremely extensive"). Changes nobody can observe: propose a skip (`step: no observable behavior change`), following Matt's note that stories fit refactors badly. Approved IDs are stable; dropped stories move to Removed. Reviewed by Codex gpt-6.2-astra (5 findings applied). Name kept. | 2026-11-21, both hosts, throwaway clones. Feature (`--json` for `tools/check`, grill.md only): 5 stories grouped by actor incl. failure paths, `e.g.` per story, handed-off items under Not covered, both flagged the unsettled can't-run output, draft at gate. Refactor (split `repo.go`): both proposed a skip without writing `stories.md`; Codex then recorded `skipped` / `step: no observable behavior change` on approval, and `spec` cited grill's requirements instead of offering stories. |
| `design` | New (Mark); no upstream source | Runs only with the right tools: a connected design tool (such as Paper) or design files the user provides; otherwise skips (`step: no design tool`, reopened by calling it again) or `step: nothing visual to design`. Proposes a depth (minimum: a full-fidelity prototype of every moving part; maximum: hi-fi of everything), user picks. Provided files are the design: mapped to stories/criteria, each gap settled by the user (design in tool, Described, or left to build). Design problems in the spec stop and reopen the spec gate. `design.md` is an index (source, depth, screen/state → coverage → location, components, tokens, decisions, left to build); build reads exact values from the tool, never screenshots. Reviewed by Codex gpt-6.2-astra (6 of 7 applied). Name kept. | 2026-11-21. Both hosts: nothing visual → skipped; no tool, no files → skipped with reopen hint (Codex counted Google Stitch as a tool, which failed auth); `mockup.svg` → mapped to S1/S4, gaps listed, prototype proposed, and Claude found two spec gaps (no EXIF date, same-second names). Claude with Paper: new file "fab design test", 4 artboards (empty, listed, mixed error, only invalid), design.md index with new tokens/components, gaps left to build, draft at gate. Codex reaches Paper (read-only check). |
| `plan` | Addy's `planning-and-task-breakdown` (MIT), Matt Pocock's `to-tickets` (MIT), Mark's waffleoff commit-strategy rule | Vertical slices (Addy, Matt's tracer bullets), each one session; prefactor first; expand/migrate/contract for wide refactors (Matt). Task = Covers (exact criterion/story/screen refs), Delivers, Checks that fail before and pass after (Matt's "criteria that already pass" failure), Verify, Blocked by, Files, Size ≤ L. Checks use the spec's approved seams. Commits: one per task with drafted message, `CP` checkpoints every 2–4 tasks and at the end, rollback steps only where git can't undo (waffleoff). >6 tasks shows an outline first (Matt's quiz, against over-splitting). Dropped: Addy's plan-mode step, rationalizations, separate todo file, Risks section; Matt's tracker publishing. Reviewed by Codex gpt-6.2-astra (8 applied; also fixed gate-skip vs step-skip in spec and design). Name kept. | 2026-11-21, both hosts. `--json` (fab clone): Claude 2 tasks (harness + characterization, then flag), Codex 1; CP1; Claude asked about a version bump. Renamer (static page + design.md): both stopped because the fixture spec claimed an existing browser harness the repo lacked (seam rule working); after approval, every design screen and both Left to build items mapped to tasks. Hosts split differently (Claude 5 tasks, Codex 1 L task); Mark chose one task per behavior, Codex re-run gave 4 tasks + CP1/CP2 with build-owned choices left undecided. Outline-first (>6 tasks) not exercised. |
| `build` | Addy's `incremental-implementation` and `git-workflow-and-versioning` (MIT), Matt Pocock's `implement` (MIT) | Per task: settle choices left to build, write checks first and see them fail (prefactors: existing checks green), simplest implementation, touch only the task (Addy's "noticed, not touching" → Noticed), verify, commit through hooks with the plan's message. Stops at each checkpoint unless told to run on (Mark). Small deviations adapt and log; changes to checks, criteria, task boundaries, dependencies, or design stop and ask (Mark). Irreversible or outward actions need prior permission. Stuck after three attempts stops. Fresh-eyes review of the whole diff before the final checkpoint (Matt's code-review step). Push and PR offered at the end (Mark). Multi-gate ledger protocol added to conventions. Dropped: feature flags, trunk-based advice, worktree/bisect tutorials, rationalization tables. Reviewed by Codex gpt-6.2-astra (8 applied). Name kept. | 2026-11-21, `--json` plan (T1 prefactor, CP1; T2 flag, CP2) in fab clones with no remote. Codex full lifecycle: branch + starting SHA, T1 recorded green-before (prefactor), CP1 `draft` and stop, approval → `started` row, T2 checks failed before (`--json` undefined) and passed after, fresh sub-agent review of the start..final range, CP2 `draft`, approval → `approved` with matching mtime, drafted PR offered (no remote, recorded in the ledger). It bumped the plugin version per AGENTS.md and logged it as a deviation. Claude through CP1: same shape; set `core.hooksPath` per AGENTS.md and logged it. Headless Claude needs `--allowedTools "Bash(git:*)" "Bash(go:*)"`. |
| `test` | Addy's `test-driven-development` (MIT), Matt Pocock's `tdd` (MIT) | Chain role is verify and fill gaps, since build already runs red/green per task (Mark). Traces every criterion to a test that asserts it at the spec's seam; runs the whole suite at build's final SHA; writes missing tests with expected values from the spec (Matt: no tautological tests); a gap test that finds a bug commits nothing failing and offers a Prove-It fix (Addy) only with approval (Mark); rewrites weak tests from this run, notes older ones (Mark); checks UI against design.md with a browser tool (Addy's DevTools workflow, tool-neutral). Dropped: test pyramid and sizes, AAA/DAMP tutorials, rationalization tables. Reviewed by Codex gpt-6.2-astra (7 applied). Name kept. | 2026-11-21, both hosts, on the built `--json` clone: all 6 criteria traced to passing tests at the final SHA, suite and checker pass, verdict pass, draft at gate. Claude caught that copied artifacts' mtimes no longer matched their approvals and asked to re-confirm (conventions working), and noted an order-dependent test. Gap-fill and Prove-It paths not exercised (no gaps). |
| `review` | Addy's `code-review-and-quality` (MIT), Matt Pocock's `code-review` (MIT), Mark's cross-host review habit | Mark's loop: the other host's CLI reviews, read-only; every finding fixed or rebutted; at least 2 rounds; stop once only nits or judgment calls remain, or at 5; open disputes go to the user. Two reviewers per round, Spec and Standards, reported side by side and never reranked (Matt). Standards brief carries a credited code-smell baseline (Matt, from Fowler) and the repo's check results. Fallback to a fresh same-host sub-agent, stated. Record in `review-<target>.md`, or chat with no run. Host CLI commands live in `references/cross-host.md` (conventions exception). Replaces build's single review pass. Dropped: Addy's change-sizing and description guidance, security/performance checklists (left to the reviewers), rationalization tables. Self-reviewed with its own loop: 5 rounds of Codex spec + standards reviewers, every bug and gap fixed; read-only reviewer commands verified against both CLIs. Name kept. | 2026-11-21, on the built `--json` clone. From Claude: Codex reviewers (cross-host), 2 rounds; one finding fixed in its own commit (shared CLI test helpers), one rebutted and accepted by both reviewers, a spec-level question (JSON path separators) raised to the user instead of fixed. From Codex: Claude failed auth inside Codex's sandbox, fallback to fresh same-host reviewers stated in chat and record, 2 clean rounds. Headless Claude couldn't read the skill's `references/` (plugin dir outside allowed dirs) and wrote the briefs from the body, saying so. |
| `pitch` | New (Mark: "grill flipped"); Addy's `idea-refine` (MIT) compared | The agent is an honest advocate: strongest true version, defend hard, concede and revise (changed lines only) when a challenge lands, numbers sourced or labeled guesses with a basis, ends with what it's least sure of. One-screen pitch, then challenges one at a time. A sinking pitch gets 1–2 one-line pivots (idea-refine's variations, only when needed). User's verdict (sold, killed, parked) at the gate; sold seeds grill's intent and decisions. Killed or parked closes the run. Name `pitch` (Mark). Reviewed with fab:review's loop: 5 rounds of Codex spec + standards reviewers, all findings fixed or rebutted (rebuttals accepted). | 2026-11-21, both hosts, three fresh sessions (reptile mask fit idea): one-screen pitch with labeled guesses and a least-sure-of line; a two-part challenge split into the part that held and the part that landed, changed lines shown; "sold" recorded without re-asking, ledger `draft` then `approved` with the verdict, For grill section filled. Claude's `--plugin-dir` run couldn't read conventions and invented its ledger format (dev-loop only; installed plugins read their cache). |

## Chain trace (CP2)

2026-11-21: one real change (`--quiet` for `tools/check`) taken through every step in a no-remote fab clone, a
fresh session per step, alternating Claude Code and Codex (installed 0.10.0). grill: dated run, facts looked up,
two questions, resumed across hosts from `grill.md`. stories: 4 stories with a failure path; flagged its own weak
story, which was moved to Not covered. spec: asked the one new-seam question, then built entirely from upstream.
design: skipped (nothing visual). plan: one task per behavior plus a final checkpoint. build: checks failed before
and passed after, a mutation check, cross-host review (Codex, 2 rounds, 6 fixes, 2 rebuttals accepted, 1 dispute
settled by the user with evidence), checkpoint approval also approved `review-build.md`. test: 6 criteria at the
final SHA, pass. Ledger complete; the run closed. Follow-ups: review's temp files fall back to `.fab/<run>/.review-tmp/`
under a sandbox (fixed); headless Claude sometimes can't read the installed `docs/conventions.md` and copies the
ledger format from existing rows (to investigate).

## Source inventory

Addy's skills as installed on Mark's machine (Claude plugin v1.0.0, plus Codex loose copies), with the fab work planned
for each. "Planned" rows become decisions only after their grill session.

| Source skill | Plan |
|---|---|
| `interview-me` (via waffleoff's interview block) | Done: `grill` |
| `spec-driven-development` | Done: `spec` |
| `planning-and-task-breakdown` | Done: `plan` |
| `incremental-implementation`, `git-workflow-and-versioning` | Done: `build` |
| `test-driven-development` | Done: `test` |
| `using-agent-skills` | Dropped (meta-router plus a session-start hook; fab's descriptions do this job) |
| `api-and-interface-design` | After v1 |
| `browser-testing-with-devtools` | After v1 |
| `ci-cd-and-automation` | After v1 |
| `code-review-and-quality` | Done: `review` |
| `code-simplification` | After v1 |
| `context-engineering` | After v1 |
| `debugging-and-error-recovery` | After v1 |
| `deprecation-and-migration` | After v1 |
| `documentation-and-adrs` | After v1 |
| `frontend-ui-engineering` | After v1 |
| `idea-refine` | Compared for `pitch`; its divergent variations live on as pitch's pivots. A separate brainstorm skill: after v1 if wanted |
| `performance-optimization` | After v1 |
| `security-and-hardening` | After v1 |
| `shipping-and-launch` | After v1 |
| `source-driven-development` | After v1 (Codex copy only) |

Named in the v1 spec as newer upstream skills but not installed here, so check upstream before auditing: a11y,
constraints, doubt, observability.

## Ideas

**Recipes**

- `bigdev` (full flow) and `tinydev` (light flow): v1, replacing waffleoff.
- `run`: user-defined recipes read from `~/.config/fab/recipes/`.

**Mark's own skills**

- Plain-language explain ("explain it like I'm five", "wait, what?"): re-explain the last thing simply, then return. Source idea: Matt's `wait-what`. Trigger style decided in its session.

- `pitchme` → shipped as `pitch` (2026-11-21). Grill flipped. The agent pitches the idea to the user (value, who it's for, how it works) and the user
  challenges it; each challenge sharpens the pitch until it holds up or falls apart. Fit against `grill`, doubt, and
  `idea-refine` is decided in its grill session (Mark, 2026-11-21: next after `review`).

**Tooling**

- Automated skill tests and CI; a generated README skill list; a version-bump tool; project-level `.fab/config`;
  artifact revision hashing in the ledger.

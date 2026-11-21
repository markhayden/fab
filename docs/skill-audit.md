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

## Source inventory

Addy's skills as installed on Mark's machine (Claude plugin v1.0.0, plus Codex loose copies), with the fab work planned
for each. "Planned" rows become decisions only after their grill session.

| Source skill | Plan |
|---|---|
| `interview-me` (via waffleoff's interview block) | Done: `grill` |
| `spec-driven-development` | Done: `spec` |
| `planning-and-task-breakdown` | Done: `plan` |
| `incremental-implementation`, `git-workflow-and-versioning` | v1 chain: `build` |
| `test-driven-development` | v1 chain: `test` |
| `using-agent-skills` | Dropped (meta-router plus a session-start hook; fab's descriptions do this job) |
| `api-and-interface-design` | After v1 |
| `browser-testing-with-devtools` | After v1 |
| `ci-cd-and-automation` | After v1 |
| `code-review-and-quality` | After v1 |
| `code-simplification` | After v1 |
| `context-engineering` | After v1 |
| `debugging-and-error-recovery` | After v1 |
| `deprecation-and-migration` | After v1 |
| `documentation-and-adrs` | After v1 |
| `frontend-ui-engineering` | After v1 |
| `idea-refine` | After v1 (weigh against `pitchme`) |
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

- `pitchme`: challenge an idea before building it. Its fit against `grill`, doubt, and `idea-refine` is decided in its
  grill session.

**Tooling**

- Automated skill tests and CI; a generated README skill list; a version-bump tool; project-level `.fab/config`;
  artifact revision hashing in the ledger.

# Skill audit log

Every fab skill is written fresh through a grill session (process in `development.md`). Concepts may come from
[addyosmani/agent-skills](https://github.com/addyosmani/agent-skills) (MIT); no text is vendored without credit.

## Decisions

One row per audited skill or recipe: source, decision, and live-test evidence.

| Skill | Source | Decision | Live test |
|---|---|---|---|
| `voice` | Mark's `write-like-me` command | Kept and generalized: the skill holds the process (load profile, draft, self-check, deliver, profile-building interview behind a gate); the voice itself moved to the private `~/.config/fab/voice.md`, built from `references/voice-template.md`. Name kept. | 2026-11-20, Claude 2.1.329 + Codex 0.170.0: no profile → offers to build, writes nothing; profile rule applied; preferences beat default; override beats preferences and profile; unreadable override → one-line warning, falls back. Real profile verified in both. |

## Source inventory

Addy's skills as installed on Mark's machine (Claude plugin v1.0.0, plus Codex loose copies), with the fab work planned
for each. "Planned" rows become decisions only after their grill session.

| Source skill | Plan |
|---|---|
| `interview-me` (via waffleoff's interview block) | v1 chain: `grill` |
| `spec-driven-development` | v1 chain: `spec` |
| `planning-and-task-breakdown` | v1 chain: `plan` |
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

- `stories`: generate user stories.
- `pitchme`: challenge an idea before building it. Its fit against `grill`, doubt, and `idea-refine` is decided in its
  grill session.

**Tooling**

- Automated skill tests and CI; a generated README skill list; a version-bump tool; project-level `.fab/config`;
  artifact revision hashing in the ledger.

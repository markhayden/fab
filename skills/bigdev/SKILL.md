---
name: bigdev
kind: recipe
description: Run a change through the full fab flow, from interview to an open pull request, with approval gates throughout, adversarial review of spec, plan, and code, and docs and style guide upkeep at the end.
disable-model-invocation: true
---

# Bigdev

The full flow for work worth doing carefully: each step's skill in order, each chain artifact approved before the
next starts, resumable at every point. It keeps `bigdev.md` in the run for what's its own: the fit check, the closing
gate, shipping, and the retro offer. Begin your first reply with `[bigdev]` and the run's path.

## Steps

Start or join the run as in [Starting or joining a run](../../docs/conventions.md#starting-or-joining-a-run); when
bigdev first takes the run, append a `bigdev` `started` row for `bigdev.md`. Resuming, the next action is the first row
below not yet settled, re-presenting drafts and changed approved artifacts as conventions' Resume says. Row 1 is
settled when the run has no pitch or its pitch was sold; rows 10 and 12–14 once their record waits at the closing
gate. Run each skill as written; the table says only what bigdev adds.

| # | Step | Skill | Gate | bigdev adds |
|---|---|---|---|---|
| 1 | Pitch | fab:pitch | kept | Only when asked, or the run already has a pitch; a killed or parked verdict ends the run with bigdev's `skipped` row (`step: pitch <verdict>`). |
| 2 | Interview | fab:grill | kept | |
| 3 | Fit check | bigdev | a question | Rules; settled once `bigdev.md` records the answer. |
| 4 | Stories | fab:stories | kept | |
| 5 | Spec | fab:spec | kept | Before its gate, spec calls fab:review on `spec.md`. |
| 6 | Reuse | fab:ui-kit (reuse) | kept | Work with no UI: a `skipped` row, `step: no UI`. Setup, if needed, keeps its own gate. |
| 7 | Design | fab:design | kept | |
| 8 | Plan | fab:plan | kept | Before its gate, plan calls fab:review on `plan.md`. |
| 9 | Build | fab:build | each checkpoint | Its push and PR offer moves to step 17. |
| 10 | Simplify | fab:simplify (build's diff) | picks kept; final at 15 | |
| 11 | Test | fab:test | kept | |
| 12 | Internal docs | fab:knowledge (this run) | picks kept; final at 15 | |
| 13 | Public docs | fab:docs (this run) | picks kept; final at 15 | |
| 14 | Style guide | fab:ui-kit (sync) | picks kept; final at 15 | Work with no UI: a `skipped` row, `step: no UI`. |
| 15 | Closing | bigdev | kept | Below. |
| 16 | Retro | fab:retro | picks and final kept | Offered once; if taken, run on this run, so its commits ship with the rest. |
| 17 | Ship | build's push and PR offer | when debug ran | The offer's answer and then CI (Rules). |

**GATE:** after steps 10–14, the closing brief in `bigdev.md` (the simplify, knowledge, docs, and ui-kit sync records
that ran, with commits and anything unverified), wait for approval. Settling it settles each record it presented, as
their caller's gate (see [Writing during a step](../../docs/conventions.md#writing-during-a-step)).

**GATE:** after step 17, when it called fab:debug, `bigdev.md` with the PR link and the CI outcome and those debug
records, wait for approval. Otherwise the PR link (or the user's choice not to push) and the CI outcome are recorded
without a stop. Then append bigdev's `approved` row for `bigdev.md`, which ends the run.

`bigdev.md` follows the several-gates rule in Writing during a step: a `draft` row at the closing gate and a `started`
row when it's approved, a `draft` row at the ship gate when it's held, and `approved` at the end. A deferred step's
`draft` row says in its note which gate it waits at. It records the fit check, both briefs, the retro
answer, the publish answer, each push, and the CI outcome (owed again after a fix push); with no one to answer, stop
with the pending item saved there.

## Rules

- Personal override: if `~/.config/fab/skills/bigdev.md` exists, read it and apply it on top of this skill (see fab conventions: Personal layer).
- **Child overrides** follow [Personal layer](../../docs/conventions.md#personal-layer), for every skill bigdev runs
  and the ones they call (review, debug, ui-kit setup), in a fresh session or sub-agent too.
- **Gates show a brief:** 3 to 6 lines on what was decided or done, what changed since the last gate, what the agent
  is least sure of, and any open review disputes, then the artifact's path; a picks gate also lists its numbered
  candidates. The full artifact only when asked.
- **Fit check, after grill:** a small, clear change with nothing left to settle: offer the lighter `tinydev` recipe
  when this fab has it. Several unrelated outcomes, or more than one run can carry: offer to split into runs, each
  starting at grill with its share of `grill.md`. Accepting either ends this run with bigdev's `skipped` row
  (`step: moved to <run>`); any other case, or a declined offer, records "fits" and goes on. Nothing settled is lost.
- **Reviewed artifacts stay reviewed:** a spec or plan changed after its review, at its gate or later, gets a review
  round on the change before it's approved.
- **Code after test:** a commit that changes code after test's approval (from a change asked for at the closing
  gate) sends the run back through test before the closing gate is approved.
- **Sessions:** grill through plan stay in one session; after plan's approval and each build checkpoint, suggest
  starting fresh, since the next session resumes from the ledger.
- **Publishing:** before any push or PR, check the remote and existing PRs so nothing is published twice, and record
  each push and PR in `bigdev.md` as it lands.
- **CI:** after the PR opens, if the repo's CI status can be read, wait for the run on the PR's latest commit. A
  failure goes to fab:debug, called from step 17; once its record shows the fix verified, push it to the PR and wait
  for the new run. Resuming, compare the git log with the PR before reading CI, so a verified fix that wasn't pushed
  yet is pushed first. A check is
  never disabled to get green. A run that's missing, cancelled, or unreadable is recorded as such.

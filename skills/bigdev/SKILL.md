---
name: bigdev
kind: recipe
description: Run a change through the full fab flow, from interview to an open pull request, with approval gates throughout, adversarial review of spec, plan, and code, and docs and style guide upkeep at the end.
disable-model-invocation: true
---

# Bigdev

Run this only when the user named bigdev; otherwise stop reading and handle the request with the skills it asks for.

The full flow for work worth doing carefully: each step's skill in order, each chain artifact approved before the
next starts, resumable at every point. It keeps `bigdev.md` in the run for what's its own: the fit check, the closing
gate, shipping, and the retro offer. Begin your first reply with `[bigdev]` and the run's path.

## Steps

Start or join the run as in [Starting or joining a run](../../docs/conventions.md#starting-or-joining-a-run), and run it
as [Running a recipe](../../docs/conventions.md#running-a-recipe) says. Row 1 is settled when the run has no pitch or
its pitch was sold. Run each skill as written; the table says only what bigdev adds.

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
| 17 | Ship | build's push and PR offer | when CI needed debug | Publishing and CI as Running a recipe says. |

**GATE:** after steps 10–14, the closing brief in `bigdev.md` (the simplify, knowledge, docs, and ui-kit sync records
that ran, with commits and anything unverified), wait for approval. Settling it settles each record it presented, as
their caller's gate (see [Writing during a step](../../docs/conventions.md#writing-during-a-step)).

**GATE:** after step 17, when a CI failure there called fab:debug, `bigdev.md` with the PR link, the CI outcome, and
those debug records, wait for approval.

## Rules

- Personal override: if `~/.config/fab/skills/bigdev.md` exists, read it and apply it on top of this skill (see fab conventions: Personal layer).
- **Fit check, after grill:** a small, clear change (not offered when the run came from tinydev's size check) with
  nothing left to settle: offer fab:tinydev. Several unrelated outcomes, or more than one run can carry: offer to split
  into runs, each starting at grill with its share of `grill.md`. Accepting either moves the run (Running a recipe); any
  other case, or a declined offer, records "fits" in `bigdev.md` and goes on. Nothing settled is lost.
- **Sessions:** grill through plan stay in one session; after plan's approval and each build checkpoint, suggest
  starting fresh, since the next session resumes from the ledger.

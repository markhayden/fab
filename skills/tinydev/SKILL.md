---
name: tinydev
kind: recipe
description: Run a small, clear change through a light fab flow - questions only for real gaps, a short spec to approve, then build, review, test, and docs upkeep shown at an end gate, then a pull request. Bug reports go through debug.
disable-model-invocation: true
---

# Tinydev

Run this only when the user named tinydev; otherwise stop reading and handle the request with the skills it asks for.

The light flow for a change whose scope is obvious: you approve what "done" means, then the finished work. Its own
record is `tinydev.md`. Begin your first reply with `[tinydev]` and the run's path.

## Steps

Start or join the run as in [Starting or joining a run](../../docs/conventions.md#starting-or-joining-a-run), and run it
as [Running a recipe](../../docs/conventions.md#running-a-recipe) says. Run each skill as written; the table says only
what tinydev adds. The route is recorded in `tinydev.md` first. A bug report (something broken, or a failing check)
takes row 1 and skips rows 2–10 and 12 (`step: tinydev bug`); any other request skips row 1 the same way (`step: tinydev
feature`).

| # | Step | Skill | Gate | tinydev adds |
|---|---|---|---|---|
| 1 | Fix | fab:debug | final at 16 | Bug reports only. First switch to a new branch, as build would, recording it and the starting SHA in `tinydev.md`; resuming reuses them. |
| 2 | Interview | fab:grill | skipped | Skips the challenge and the sweep (Sweep says "skipped: tinydev"). Asks only what would change what "done" means, inferred decisions included; when nothing would, it writes `grill.md` and goes on. |
| 3 | Size check | tinydev | a question | Rules. |
| 4 | Stories | fab:stories | — | Skipped, `step: tinydev`. |
| 5 | Spec | fab:spec | kept | Works from `grill.md`. More gaps than spec may ask about: make the bigdev offer then; declined, record "stays tiny (at spec)", rerun grill on them as in row 2 until its Open list is empty, then finish spec. |
| 6 | Reuse | fab:ui-kit (reuse) | — | Skipped, `step: tinydev`. |
| 7 | Design | fab:design | — | Skipped, `step: tinydev`. |
| 8 | Plan | fab:plan | skipped | Usually one task and a final checkpoint. |
| 9 | Size check | tinydev | a question | Rules, now against the plan. |
| 10 | Build | fab:build | final at 16; others skipped | An intermediate checkpoint gets a `started` row noting `gate: tinydev, CP<n>`; its user-run checks still wait for the user. |
| 11 | Simplify | fab:simplify (the run's diff) | picks kept; final at 16 | On the bug route its commits get the code-after rule's review round and rerun debug's loop and the repo's checks. |
| 12 | Test | fab:test | final at 16 | |
| 13 | Internal docs | fab:knowledge (this run) | picks kept; final at 16 | |
| 14 | Public docs | fab:docs (this run) | picks kept; final at 16 | |
| 15 | Style guide | fab:ui-kit (sync) | picks kept; final at 16 | Work with no UI: skipped, `step: no UI`. Setup, if needed, keeps its own gate. |
| 16 | End | tinydev | kept | Below. |
| 17 | Retro | fab:retro | picks and final kept | Offered once; if taken, run on this run, so its commits ship with the rest. |
| 18 | Ship | push and PR offer, as build makes it | when CI needed debug | Publishing and CI as Running a recipe says. |

**GATE:** the end brief in `tinydev.md` (the fix's debug record, or build's final checkpoint and `tests.md`, plus each
row 11 and 13–15 record that ran, with commits and anything unverified), wait for approval. Approving it approves
build's final checkpoint and `tests.md` when they exist, and settles each called record it presented (see [Writing
during a step](../../docs/conventions.md#writing-during-a-step)).

**GATE:** after step 18, when a CI failure there called fab:debug, `tinydev.md` with the PR link, the CI outcome, and
those debug records, wait for approval.

## Rules

- Personal override: if `~/.config/fab/skills/tinydev.md` exists, read it and apply it on top of this skill (see fab conventions: Personal layer).
- **Size check:** offer to hand the run to fab:bigdev when grill settled more than a couple of real decisions, the work
  needs stories (several kinds of people do or see different things, or it crosses a trust boundary), it has UI to
  design, or the plan has more than about three tasks. bigdev takes the run as Running a recipe says, with every
  approval kept. No trigger, or a declined offer, records "stays tiny (after grill)" or "stays tiny (after plan)" in
  `tinydev.md` and goes on; each size-check row is settled once its own answer is there.

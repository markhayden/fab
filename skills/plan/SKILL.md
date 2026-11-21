---
name: plan
kind: skill
description: Break an approved spec (and design) into plan.md, a set of small vertical-slice tasks, each with checks that prove it, plus the commit strategy and rollback checkpoints. Use after the spec or design, or when the user says "plan this", "break it down", or "what are the steps?".
---

# Plan

Turns the spec into an ordered set of tasks that each deliver something you can see working, with the proof and the
commit for each, so build can work one task at a time and stop at safe rollback points.

## Inputs

- The run (`.fab/<YYYY-MM-DD>-<slug>/`): start or join one first, as in [Starting or joining a run](../../docs/conventions.md#starting-or-joining-a-run).
- `spec.md` from the run, plus `design.md` and `stories.md` when present, and, when the run has one, its reuse record
  (`ui-kit-reuse.md`, or its latest numbered pass, finished first if still in draft) once approved or gate-skipped:
  tasks name the components and tokens they use and the shared components they create. An artifact whose approval gate
  was skipped still counts; one whose step was skipped (`step:` note) doesn't.
- The repo: `AGENTS.md`/`CLAUDE.md` (commands, branch and commit rules), the code the change touches, existing tests.
- Preferences: `~/.config/fab/preferences.md`. Anything it settles is not asked.
- Personal override: if `~/.config/fab/skills/plan.md` exists, read it and apply it on top of this skill (see fab conventions: Personal layer).

## Process

Follow [Writing during a step](../../docs/conventions.md#writing-during-a-step). Planning changes no code.

1. **Read the change's ground.** The code it touches, the tests that exist, the repo's commands. If something would
   make the change much easier (a split, an extracted helper), plan it as the first task. Settle choices that grill
   or spec handed to `plan` and record them in the Overview; attach choices handed to `build` to the tasks they
   affect, undecided.
2. **Slice vertically.** Each task delivers one behavior that works end to end and can be shown on its own, and
   fits in one fresh session. No task builds one layer for a later task to wire up. Several behaviors means several
   tasks, even when they'd fit one session together, so every commit is a working rollback point; one behavior is
   one task, never split. Exceptions: a prefactor shows structural progress instead, with existing checks staying
   green; a rename or retype that breaks every caller at once goes expand (add the new form beside the old),
   migrate callers in batches, then contract (remove the old form). For a schema change, dropping the old form is its
   own task after a checkpoint, and a reversible migration task's checks include its down path (an irreversible one
   gets rollback steps, step 4).
3. **Prove each task.** Its checks must fail before the work and pass after; a check that already passes proves
   nothing (a prefactor's checks are the existing ones, still passing; baseline, guard, and user-run checks prove
   themselves as below). Use the seam the spec's Testing decisions
   chose for each criterion; a different seam or a new harness goes back to the spec's gate. Every spec criterion
   is covered by at least one task, and every design Screens row, Described item, and Left to build item belongs
   to a task. A criterion that's a number ("p95 under 200 ms") gets a baseline task first, whose proof is the
   recorded measurement, and a check that states its threshold and how many runs it takes; if the baseline already
   meets it, the check guards it and says so.
4. **Order and checkpoint.** Order by what blocks what. Put a checkpoint every 2–4 tasks where everything passes and
   rolling back is safe. A step that git can't undo (a data migration, a deletion outside the repo) gets its own
   rollback steps at its checkpoint. A check the building agent can't run (it needs a permission, a device, or the
   user's hands) is marked user-run and listed at its checkpoint with the exact command or steps; the user runs it
   there.
5. **More than six tasks:** show a one-line-per-task outline first and let the user merge or split, using
   [grill's Asking rules](../grill/SKILL.md#asking), before writing the details.
6. **A requirement is unresolved or contradictory** in a way that blocks planning: stop and ask. A change to an
   approved artifact sends it back to its gate before you continue. Choices already delegated to build stay with
   build.

**GATE:** present `plan.md`, wait for approval. Offer it as a choice (approve, or name what to change), the way
grill asks.

## Output

`.fab/<run>/plan.md`:

- **Overview:** the approach in a few lines, and any prefactoring and why.
- **Tasks:** `T1`, `T2`, … Each has: title; **Covers** (`criterion 2`, active story IDs like `S1`, and exact design
  screen or state names; never invented or Removed IDs); **Delivers** (what
  you can see working after it); **Checks** (fail now, pass after, or which exception from step 3 applies);
  **Verify** (the command, never including a user-run check); **Blocked by**;
  **Files** likely touched; **Size** (XS–L; split anything bigger).
- **Commits:** the branch and PR approach (from `AGENTS.md` or preferences), and one commit per task with its
  drafted Conventional Commits message. Checkpoints `CP1`, `CP2`, …: the task each follows, what must pass there,
  and any rollback steps git can't do. The last task is always followed by one. Build records each checkpoint's
  SHA in `build.md`.
- **Open questions:** anything unresolved. One that changes a task's checks blocks approval.

## Done when

- Every spec criterion maps to at least one task, and every task cites what it covers.
- Every task has checks that fail before the work (or, for a prefactor, baseline, guard, or user-run check, the
  proof step 3 and 4 name), a verify command, and a size of L or smaller.
- Commits list one message per task, and checkpoints every 2–4 tasks, ending with one where all checks pass.
- The ledger shows `started`, `draft`, then `approved` with the mtime (or `skipped` with a `gate` note when a recipe
  skips approval), for `plan.md`.

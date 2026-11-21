---
name: build
kind: skill
description: Carry out an approved plan.md task by task, writing each task's checks first, committing each task, and stopping at every checkpoint, with build.md recording SHAs and deviations. Use after the plan, or when the user says "build it", "implement the plan", or "start coding".
---

# Build

Turns the approved plan into commits, one task at a time, each proven by its checks, so every checkpoint is a state
you can roll back to and trust.

## Inputs

- The run (`.fab/<YYYY-MM-DD>-<slug>/`): start or join one first, as in [Starting or joining a run](../../docs/conventions.md#starting-or-joining-a-run).
- `plan.md` and `spec.md` from the run, and `design.md` when present. A skipped approval gate counts; a skipped
  step (`step:` note) doesn't. If an approved artifact's mtime no longer matches its approval row, re-present it
  first.
- `build.md`, if build already started: it says which tasks are done.
- The repo's `AGENTS.md`/`CLAUDE.md`: commands, commit rules, hooks. They always apply.
- Preferences: `~/.config/fab/preferences.md`. Anything it settles is not asked.
- Personal override: if `~/.config/fab/skills/build.md` exists, read it and apply it on top of this skill (see fab conventions: Personal layer).

## Process

Follow [Writing during a step](../../docs/conventions.md#writing-during-a-step); update `build.md` after every task.

1. **Start or resume.** Switch to the plan's branch, creating it if needed, and record the starting SHA.
   Uncommitted changes that aren't this run's: ask before touching anything. Resuming: if the last `build` ledger
   row is a checkpoint `draft`, re-present that checkpoint first. Otherwise the next task is the first one
   `build.md` doesn't mark done, and `HEAD` must match the last recorded SHA; if it doesn't, ask.
2. **For each task, in order:**
   - Settle any choices the plan or design left to build for this task, within the approved criteria and design,
     and record each one and why.
   - Before writing code that relies on an API the spec marks "unverified", check it against the docs for the
     version in use (look the version up when the spec couldn't); docs that can't be reached leave it unverified in
     `build.md`, and docs that contradict the spec are a deviation.
   - Work one check at a time ([references/tests.md](references/tests.md)): write it, run it and record it failing, then
     implement the simplest thing that makes it pass; repeat for the task's next check. Exceptions: a prefactor records
     the existing checks passing before and after; a baseline task records its measurement; a guard check the plan says
     already passes records that; a user-run check waits for the checkpoint. Any other check that unexpectedly passes:
     find out why without changing what it proves; if what it should prove must change, that's a deviation (step 3).
     Every behavior the task delivers is implemented now, including what only a user-run check covers. A numeric check
     logs each attempt in `build.md`, kept or reverted. Touch only what the task needs; anything else worth fixing goes
     under Noticed in `build.md`, not into the diff.
   - Before anything irreversible or outward-facing (a data migration, a deletion outside the repo, a call to a
     live service), confirm the spec's Boundaries or the user already allow it; otherwise stop and ask. A later
     checkpoint is not permission.
   - Run the task's Verify command and the repo's own checks, leaving user-run checks to the checkpoint. Visual work
     takes exact values from the design tool when it's connected; with no `design.md`, look at it about 375 px and
     1280 px wide, following [Using a browser](../../docs/conventions.md#using-a-browser), and record what you saw
     (or why it couldn't be checked: no browser tool, the app won't run here, or no isolated profile).
   - Commit with the plan's message (adjusted if the work shifted), through the repo's hooks, never bypassing
     them. Record the SHA.
3. **Deviations.** Small ones (different files, an extra helper, a renamed test): adapt and log them. Anything that
   changes a check, a criterion, task boundaries, a dependency, or the design: stop and ask. A change to an
   approved artifact sends it back to its gate.
4. **Stuck.** Three attempts at the same failure without progress: stop, and report what was tried and what you'd
   try next.
5. **Before the final checkpoint,** run fab:review on the whole diff (starting SHA to `HEAD`) against `spec.md` and
   `plan.md`. Record each of its fix commits in `build.md` as it lands. Approving the final checkpoint also
   approves `review-build.md`.
6. **At each checkpoint,** run its agent-run checks, record the checkpoint's SHA, and append a `draft` row noting the
   `CP`. Give the user each user-run check (its command or steps), marked pending in `build.md`, and wait for the result
   they report, even when told to run to the end; with no one to answer, stop with them pending (a resumed checkpoint
   asks again). The final checkpoint also asks again for any earlier user-run check whose code changed since its result.
   Record each result, and report: tasks done with SHAs, what passed, deviations, and anything noticed. Then stop,
   unless the user said to keep going (for example, "run to the end"). Approving a checkpoint before the last appends a
   `started` row noting the `CP` and who approved; only the final checkpoint gets the `approved` row.

**GATE:** at each checkpoint, the progress report and `build.md`, wait for approval before the next task.

7. **At the end,** offer to push the branch and open the PR, with a drafted title and body
   ([references/pr.md](references/pr.md)) and no AI attribution.
   A standing rule in `AGENTS.md` or preferences settles this without asking. Record the outcome (and the PR link)
   as a ledger note, not in the approved `build.md`.

## Output

Commits on the plan's branch, and `.fab/<run>/build.md`:

- **Commits:** the starting SHA and the final verified SHA, so review and test can find the whole diff.
- **Tasks:** a table of task → status (done, in progress, blocked) → SHA.
- **Checks:** for each, the task, the criteria it proves (`criterion 2`), the test's name and location, the
  command, and the result before and after. Exceptions record what proves them instead: prefactors (existing checks
  passing both times), baselines (the measurement), guards (passing), user-run checks (the command or steps and the
  user's reported result, or pending).
- **Build choices:** each choice left to build, what was chosen, and why.
- **Checkpoints:** each `CP` → SHA → what passed.
- **Deviations:** each change from the plan, why, and whether the user approved it.
- **Noticed:** things outside the tasks worth doing later, untouched.
- **Review:** a pointer to `review-build.md`, how many rounds ran, and any open disputes.

## Done when

- Every plan task is done with a SHA, and its checks are recorded failing before and passing after (prefactors:
  passing both times; baselines: the measurement; guards the plan says already pass: passing; user-run checks: the
  user's result at the checkpoint).
- The final checkpoint's checks pass at its recorded SHA, after the review's fixes.
- The ledger shows `started`, `draft` at the final checkpoint, then `approved` with the mtime (or `skipped` with a
  `gate` note when a recipe skips approval), for `build.md`.

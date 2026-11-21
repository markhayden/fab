---
name: test
kind: skill
description: Verify a build against its spec and design, fill test gaps, and rewrite weak tests, recording evidence in tests.md. Proves bugs with a failing test before fixing them. Use after build, for "test this" or "verify it", or when the user reports a bug.
---

# Test

Shows, with evidence, that what was built does what the spec says: every criterion has a real test that passes at
the build's final commit. It also turns a bug report into a failing test before anything is fixed.

## Inputs

- The run (`.fab/<YYYY-MM-DD>-<slug>/`): start or join one first, as in [Starting or joining a run](../../docs/conventions.md#starting-or-joining-a-run).
  A bug report outside a build uses a run too, but follows Proving a bug instead of steps 1–6.
- `build.md` (its final SHA and Checks), `spec.md`, `plan.md`, and `design.md` when present. A skipped approval
  gate counts; a skipped step (`step:` note) doesn't.
- The repo's `AGENTS.md`/`CLAUDE.md`: test commands, commit rules, hooks.
- Preferences: `~/.config/fab/preferences.md`. Anything it settles is not asked.
- Personal override: if `~/.config/fab/skills/test.md` exists, read it and apply it on top of this skill (see fab conventions: Personal layer).

## Process

Follow [Writing during a step](../../docs/conventions.md#writing-during-a-step).

1. **Start at the build's end.** First time: `HEAD` is build's final SHA and the tree is clean; record it as the
   baseline. Resuming: `HEAD` is the latest SHA in `tests.md`. Otherwise ask. Never edit the approved `build.md`.
2. **Trace.** For every spec criterion, find the tests that prove it (start from `build.md`'s Checks), and read
   them: a test proves a criterion only if it asserts that behavior through the spec's seam.
3. **Run everything** at this SHA: the full suite, each task's Verify command, and the repo's own checks. A test
   that flips between runs is flaky; record it as weak. Record every failure with its output; an unexplained
   failure anywhere in the suite means the verdict isn't a pass, even if every criterion's tests pass.
4. **Fill gaps.** A criterion with no real test gets one, written at the spec's seam, with its expected values
   taken from the spec, not from the code. If it passes, commit it (`test(<scope>): …`). If it fails, you found a
   bug: commit nothing failing, report it with the test and its output, and offer to fix it Prove-It style
   (below) as a logged deviation. Never change product code unasked.
5. **Weak tests.** Check this run's tests for ones that can't fail (the expected value is computed the way the code
   computes it), test internals instead of behavior, mock what could be real, or depend on timing or order.
   Rewrite those; commit only rewrites that pass, and treat one that exposes a bug like a failing gap test.
   Weak tests older than the run go under Noticed, untouched.
6. **Check UI against the design** when `design.md` exists and a browser tool is available: every Screens row,
   Described item, and Left to build item (using the choice `build.md` recorded). Compare each with the design
   and check the console is clean. A mismatch is reported like a bug. No browser tool: list each as not verified.
7. **After any commit here,** rerun everything from step 3 and record the final verified SHA in `tests.md`.

**GATE:** present `tests.md`, wait for approval. Offer it as a choice (approve, or name what to change), the way
grill asks.

### Proving a bug (standalone, or a fix offered above)

1. Write a test at the highest seam that reproduces the bug, run it, and confirm it fails for the reported reason.
2. Make the smallest fix that passes it, then run the full suite.
3. Commit the test and the fix together (`fix(<scope>): …`), through the repo's hooks.

## Output

Commits for any tests added or rewritten (and fixes the user approved), plus `.fab/<run>/tests.md`:

- **Revision:** the baseline SHA, the final verified SHA, and the commands run with their totals.
- **Verdict:** pass, or pass with failures the user accepted, or not passing; plus anything not verified.
- **Criteria:** a table of spec criterion → tests that prove it (name and location) → result.
- **Design:** a table of screen or state → how it was checked → result, or "not verified" with the reason.
- **Gaps filled:** each new test, the criterion it covers, and its SHA.
- **Bugs found:** each with its failing test, output, and what the user decided.
- **Weak tests:** rewritten tests, with SHAs.
- **Noticed:** older weak tests and other findings outside the run, untouched.
- **Deviations:** each Prove-It fix the user approved, why, and its SHA.
- For a bug report: the report, the reproducing test, its output before and after the fix, and the SHA.

## Done when

- Every spec criterion maps to at least one test that passes at the recorded SHA, or to a bug the user decided on.
- Every design item is checked or marked not verified with a reason.
- The suite and repo checks pass at the final verified SHA, or each failure has the user's decision.
- The ledger shows `started`, `draft`, then `approved` with the mtime (or `skipped` with a `gate` note when a recipe
  skips approval), for `tests.md`.

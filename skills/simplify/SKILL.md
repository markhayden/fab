---
name: simplify
kind: skill
description: Consolidate and clean up code without changing behavior, on a diff or across an area, and deepen shallow modules the user picks, recording it in simplify-<target>.md. Use after a build, or when the user says "simplify", "clean this up", "consolidate", "reduce debt", or "this is too complicated".
---

# Simplify

Makes code easier to change without changing what it does: removes duplication and dead weight, and turns
shallow modules into deep ones (a lot of behavior behind a small interface), recorded in `simplify-<target>.md`.
fab:review flags these problems; this skill fixes them.

## Inputs

- The target: a **diff** (a fixed point to `HEAD`, such as build's starting SHA) or an **area** (a path, a module, or
  the repo). For an area with no direction given, weight the files recent commits keep touching.
- The run (`.fab/<YYYY-MM-DD>-<slug>/`): when another step or a recipe calls this, the caller's run (see
  [Writing during a step](../../docs/conventions.md#writing-during-a-step)); called by a review with no run, skip the
  file and the ledger and hold both gates in chat. Otherwise start or join one as in
  [Starting or joining a run](../../docs/conventions.md#starting-or-joining-a-run).
- fab:review records for the target (`review-*.md`): Standards findings that were rebutted as out of scope are candidates.
- The repo's standards (`AGENTS.md`/`CLAUDE.md`, standards and style docs), its checks, and any glossary or decision
  records covering the code. A pattern the standards endorse isn't a candidate; one a decision record settled is
  raised only if the friction is real, and marked as reopening it. Use the glossary's terms in the record.
- Preferences: `~/.config/fab/preferences.md`. Anything it settles is not asked (picks excepted).
- Personal override: if `~/.config/fab/skills/simplify.md` exists, read it and apply it on top of this skill (see fab conventions: Personal layer).

## Process

Follow [Writing during a step](../../docs/conventions.md#writing-during-a-step),
[Committing during a step](../../docs/conventions.md#committing-during-a-step), and
[Picks](../../docs/conventions.md#picks): save `simplify-<target>.md` after every exchange.

1. **Start from green.** Record the starting SHA and run the repo's checks; any failure stops here, since a red
   baseline can't show that behavior was kept. Before the first change, the tree must be clean (ask otherwise).
2. **Find candidates.** Before touching code, know what it does, who calls it, which tests pin it, and why it was
   written that way (its history). Look for the smells in fab:review's Standards brief
   ([briefs](../review/references/briefs.md)) and for shallow modules
   ([references/deepening.md](references/deepening.md)). On a diff, look only at code the diff touched and its
   duplicates; across an area, use a fresh sub-agent for the sweep if one is available. Sort each candidate:
   - **Safe:** every interface stays as it is (signatures, outputs, errors, side effects, ordering) and no test
     changes.
   - **Deepening:** an interface changes and its tests move to the new one.
3. **Diff, safe tier:** apply them one at a time. Make the change, run the checks; green with tests untouched →
   its own commit (`refactor` where the repo has that type). Otherwise revert it: if a test needed editing because
   an interface changed, it becomes a deepening candidate; if not, drop it with the reason. Code the spec or repo
   marks as performance-sensitive is measured before and after, and a slower result is reverted. Mechanical edits
   across more than ~500 lines go through a script, not by hand. Keep each commit independent of the others:
   changes that only work together are grouped before committing, and a commit already recorded is never
   rewritten.
4. **GATE:** present `simplify-<target>.md` with the open candidates (tier, place, evidence, effort, how strongly
   you recommend it; your top pick first), wait for approval. This is the picks gate. The first time on a diff with
   nothing left to propose, skip it; a reopened request always goes through it.
5. **Do the picks.** Area safe-tier picks go as in step 3. A deepening that fits one session (fab:plan's sizes,
   L or smaller):
   - frame its constraints and dependencies, and when the shape isn't obvious, sketch two or three contrasting
     interfaces (fresh sub-agents if available) and recommend one; the user picks, as a choice in
     [fab:grill's Asking](../grill/SKILL.md#asking) style, saved as Pending (a question, not a gate: no ledger row);
   - write tests at the new interface first, move the callers, then delete old tests only where the new ones cover
     their behavior (record which new test covers each deleted one);
   - measured as in step 3 where performance matters, checks green, then committed as in step 3.

   A bigger one becomes a new run (never the caller's): create it as in Starting or joining a run, with the
   candidate as its request, record its path, and leave fab:grill for the user to start after this step.
6. **Step back.** Read the result as a whole. Revert any change that didn't make the code easier to read
   (`git revert --no-commit`, newest first, keeping any version bump; run the checks, then commit). A conflict or a
   red check: `git revert --abort`, and report it.

Never mix these commits with feature or fix work, and never weaken a test or check to keep one green. A change
that would alter behavior isn't a simplification: stop and ask.

**GATE:** present `simplify-<target>.md` with every commit (each revertable on its own), wait for approval. Offer it
as a choice (approve, or name what to change), the way fab:grill asks. When another step or a recipe called this in
a run, the record names its caller and goes to the caller's gate instead; with no run, this gate is held in chat.

## Output

Commits, plus `.fab/<run>/simplify-<target>.md`:

- **Target:** the diff range or area, the caller if any, the starting SHA, and the baseline check results.
- **Candidates:** numbered; tier, place, evidence, effort, strength.
- **Picks:** the numbers the user chose at each picks gate.
- **Applied:** each commit's SHA and what it did, with the checks that passed after it.
- **Deepenings:** the designs offered and the one picked, with SHAs, or the new run's path.
- **Reverted or declined:** each with the reason.
- **Next / Pending:** while it runs, the change in progress and anything waiting on the user.

## Done when

- Every commit passes the repo's checks, and safe-tier commits left every test file untouched.
- Every candidate is applied, declined, reverted, or handed to a new run.
- In a run, the ledger shows `simplify` rows for `simplify-<target>.md` (not a chain step, so its rows never
  decide a run's next step): `started`, `draft` and `started` around the picks gate when it runs, `draft` at the
  final gate, then `approved` with the mtime, or `skipped` with a `gate` note when a caller's gate or a recipe
  settles it.

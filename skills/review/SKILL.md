---
name: review
kind: skill
description: Adversarial review in rounds by the other hosts' agents, each a spec and a standards reviewer; every finding is fixed or rebutted. Use on a diff, a branch, or a run artifact, or when the user says "review this" or "get a second opinion".
---

# Review

Has different models attack the work before anyone relies on it, and records every round in `review-<target>.md`.
Two reviewers per host each round keep "does it do the right thing" separate from "is it done right".

## Inputs

- The target: a diff (a fixed point to `HEAD`, e.g. build's starting SHA), a branch, or a run artifact (`spec.md`,
  `plan.md`, …).
- The run: when one is named or open for this work, join it as in [Starting or joining a run](../../docs/conventions.md#starting-or-joining-a-run).
  With no run (say, reviewing someone's branch), the record goes in chat and there's no gate.
- What the target answers to: for a diff, `spec.md` (plus `plan.md` and `design.md`); for an artifact, the
  artifacts before it in the chain; with neither, the user's request.
- The repo's own checks (tests, linters, `AGENTS.md` commands), run before each round.
- The repo's written standards: `AGENTS.md`/`CLAUDE.md`/`GEMINI.md`, `CONTRIBUTING.md`, `CODING_STANDARDS.md`, style or
  convention docs.
- Preferences: `~/.config/fab/preferences.md`. Anything it settles is not asked.
- Personal override: if `~/.config/fab/skills/review.md` exists, read it and apply it on top of this skill (see fab conventions: Personal layer).

## Process

In a run, follow [Writing during a step](../../docs/conventions.md#writing-during-a-step): save the record after
every round and every answer, with the round number, the current SHA, and an artifact target's mtime, so an
interrupted review resumes where it was. Resuming (or presenting a finished record) when `HEAD` or the artifact's
mtime doesn't match what was recorded: ask before going on.

1. **Pin the target.** Resolve the fixed point and confirm the diff isn't empty, or confirm the artifact exists.
   Uncommitted changes that aren't part of the target: ask before making any fix.
2. **Pick the reviewers.** Every other host fab supports whose command-line agent is installed (not the one you're in),
   each on its most capable model, run read-only, with a fresh context each time
   ([references/cross-host.md](references/cross-host.md) has the hosts, models, and commands). Each one runs both
   reviewers every round, so two other hosts means four reviewers. Each reviewer gets one retry. An axis no cross-host
   reviewer covered in a round gets a fresh sub-agent on this host for that round, or a careful pass of your own if
   there's none, so every round has both axes. A host whose reviewers all fail in a round drops out, and the record
   names it and why. With no host left, both axes fall back the same way, and Coverage marks those rounds local.
3. **Run a round.** Run the repo's checks, then start every reviewer together with the briefs in
   [references/briefs.md](references/briefs.md), giving the standards reviewers the check results:
   - **Spec:** what the upstream asked for that's missing or partial, what was built that wasn't asked for, and
     what looks implemented but wrong.
   - **Standards:** where the target breaks the repo's written rules (citing the rule), and, for code, the smells
     in the briefs file. A written repo rule beats a smell.
   Each finding comes back with its axis, a severity (`bug`, `gap`, or `nit`), and the exact place.
4. **Answer every finding,** whatever its size: fix it, or write why not. A bug is fixed with fab:debug (usually
   its fast path). Every fix to tracked files passes the repo's checks
   and goes in its own commit through the hooks; edits to run artifacts stay uncommitted. Never fix by weakening a test or a check. A fix that would change a criterion, a
   check, a dependency, a task boundary, or the design isn't yours to make: it becomes an open dispute. Edits to an approved artifact send it back to its gate.
   A rebuttal settled by preferences is marked "from preferences".
5. **Next round:** every reviewer gets the target as it is now, plus every earlier finding with its fix or rebuttal.
   They confirm the fixes, accept or contest the rebuttals, and look for anything new.
6. **Stop** after at least two rounds, once the latest round has only `nit`s or judgment calls left, or after five
   rounds. Rebuttals a reviewer still contests are open disputes for the user.

**GATE:** present `review-<target>.md` and any open disputes, wait for approval. When another step called review
(build before its final checkpoint), the record goes to that step's gate instead. With no run, there's no gate.

## Output

Fix commits or artifact edits, plus `.fab/<run>/review-<target>.md` (e.g. `review-build.md`); outside a run, the same
record in chat:

- **Target:** what was reviewed (diff range or artifact), what it answered to, the reviewers (each host and model,
  or the fallback and why), and any installed host skipped before round 1, with why.
- **Coverage:** for each round, every reviewer that failed its retry (host, axis, why) and every axis filled on this
  host, so a missing reviewer never reads as a clean pass.
- **Rounds:** for each, every reviewer's findings under its host and axis, each in its reviewer's own order (never
  merged or reranked; the same issue from two reviewers is answered once and cross-referenced), each finding with
  severity, place, and what happened (fixed with SHA or edit, rebutted with the reason, accepted or contested next
  round).
- **Open disputes:** each contested rebuttal, both sides in a line each.
- **Stopped because:** only nits or judgment calls left after round N, or the five-round limit.

## Done when

- At least two rounds ran (unless the target is gone or empty), and every finding is fixed or rebutted.
- The record names each reviewer's host and model, and why the loop stopped.
- In a run, the ledger shows `review` rows for `review-<target>.md` (the Artifact column tells reviews of
  different targets apart; review isn't a chain step, so its rows never decide a run's next step): `started`, `draft` when the loop stops, then
  `approved` with the mtime when its gate (or the calling step's) is approved, or `skipped` with a `gate` note when
  a recipe skips approval.

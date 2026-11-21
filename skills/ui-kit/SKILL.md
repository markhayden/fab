---
name: ui-kit
kind: skill
description: Keep a repo's reusable UI (components, tokens, style guide) in use and in sync - say what plan should reuse, send one-off UI to simplify, and update the style guide after a run. Use for "what components do we have", "update the style guide", or around plan and build for UI work.
---

# UI kit

Keeps a project's UI built from its shared pieces, with one job per call, each recorded in `ui-kit-<job>.md`.
Designing a feature's screens is fab:design's job; refactoring is fab:simplify's.

## Inputs

- The job: **setup** (record the kit), **reuse** (before plan), **componentize** (after build), or **sync** (at the
  end of a full run). The caller names it; standalone, offer the one the run is ready for and ask.
- The repo's **Design system section**: a short section in its docs (or its agent instructions file when there's no
  docs folder) naming the source of truth for components, for tokens, and for the style guide, where each lives,
  which wins when they disagree, whether a guide was offered, accepted, or declined, and any guide changes still
  pending. Any job without one, or with an unfinished `ui-kit-setup.md`, runs setup first, with its own record.
- The kit: components and tokens in the code, the repo's style guide (Storybook stories or similar), and a guide in
  a connected design tool (such as Paper).
- The run (`.fab/<YYYY-MM-DD>-<slug>/`): when another step or a recipe calls this, the caller's run (see [Writing during
  a step](../../docs/conventions.md#writing-during-a-step)); called (directly or not) by a review with no run, the
  record (and anything handed to simplify) goes in chat, working on the review's target. Otherwise start or join one as
  in [Starting or joining a run](../../docs/conventions.md#starting-or-joining-a-run), or, for componentize or sync on a
  finished run, choose it as in [Looking back at a run](../../docs/conventions.md#looking-back-at-a-run).
- Preferences: `~/.config/fab/preferences.md`. Anything it settles is not asked (picks excepted).
- Personal override: if `~/.config/fab/skills/ui-kit.md` exists, read it and apply it on top of this skill (see fab conventions: Personal layer).

## Process

Follow [Writing during a step](../../docs/conventions.md#writing-during-a-step),
[Committing during a step](../../docs/conventions.md#committing-during-a-step), and
[Picks](../../docs/conventions.md#picks): save `ui-kit-<job>.md` (in chat with no run) after every
exchange. Before each design-tool write, note it as in progress; on resume, check the tool for it before writing
again. This skill never changes product code: where the guide is the source of truth and the code disagrees, report
it and offer fab:debug.

1. **Setup.** Inventory the components and tokens, then ask, using
   [fab:grill's Asking rules](../grill/SKILL.md#asking), what the source of truth is for each, where it lives, and
   which wins on conflict. With UI but no style guide, offer once to start one and record the answer. Write the
   Design system section and commit it.
2. **Reuse.** Map every piece of UI the spec calls for (starting from `design.md`'s components and tokens when
   present) to an existing component and tokens, or mark it new. A new piece the work uses in two or more places, or
   that matches a pattern the kit already has, is marked "make it a shared component"; a value with no token is
   marked "new token". A design that contradicts the kit is asked about. Reuse always settles at its own gate,
   before plan starts, since fab:plan uses this record only once it's approved or gate-skipped.
3. **Componentize.** In the run's changes (with no run, the review's target), find UI built more than once, near-copies
   of kit components, pieces reuse marked as shared that weren't, and hard-coded values that have tokens. Hand them to
   fab:simplify as an area pass over those files with these as its candidates, and record its path and outcome; simplify
   keeps its own rules, picks included.
4. **Sync.** Work on the branch holding the run's commits, with a clean tree (ask otherwise). Start from the guide
   changes the Design system section lists as pending, dropping any the code has since made moot. Then, for every
   component or token the run added, changed, or removed, list the guide changes the section calls for: stories or a
   tokens page in the repo, a page in the connected design tool, or nothing where a guide was declined. Starting a
   guide (a dependency or a new design file) is a proposal, never a default; one the user accepted at setup stays on
   the list until it exists.

   **GATE:** present `ui-kit-sync.md` with the guide changes, wait for approval. This is the picks gate, held for
   any design-tool change, which is written only once picked.

   Then make the picks: repo guide changes in a commit, built or rendered to prove them (marked unverified when that
   can't run here); design-tool changes through the tool, each result recorded, or marked blocked when the tool
   isn't connected or won't accept the write. Blocked changes, and repo changes rolled back after a failed check, go
   on the section's pending list; applied or declined ones come off it, in the same commit as the section.

**GATE:** present `ui-kit-<job>.md` with any commits, wait for approval. Offer it as a choice (approve, or name what
to change), the way fab:grill asks. When another step or a recipe called this in a run, the record names its caller
and goes to the caller's gate instead (reuse and setup excepted).

## Output

Commits (the Design system section, repo guide changes), plus `.fab/<run>/ui-kit-<job>.md`:

- **Kit:** the source of truth for components, tokens, and guide.
- **Setup:** what was inventoried and the answers recorded. **Reuse:** each UI piece and what to use, or "new"
  (shared or not, new tokens), with questions raised. **Componentize:** candidates and the simplify record they went
  to. **Sync:** each guide change, its SHA or design-tool result, and how it was proven, or why it's unverified or
  blocked.
- **Pending:** while it runs, anything waiting on the user.

## Done when

- The repo has a Design system section, and the job's record has its section filled.
- Reuse maps every UI piece; componentize sent every candidate to simplify; sync made, declined, or put back on the
  pending list every listed guide change.
- In a run, the ledger shows `ui-kit` rows for `ui-kit-<job>.md` (not a chain step, so its rows never decide a run's
  next step): `started`, `draft` and `started` around sync's picks gate when it's held, `draft` at the final
  gate, then `approved` with the mtime, or `skipped` with a `gate` note when a caller's gate or a recipe settles it.

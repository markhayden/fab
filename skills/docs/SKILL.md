---
name: docs
kind: skill
description: Keep a repo's public docs (README, usage guides, API reference, changelog, examples) true after a change and add what users now need, running the commands and examples it documents where it safely can. Use after a run, or for "update the README", "document this", or "are the docs right".
---

# Docs

Keeps what users read about a project true and complete, running what it documents where that's safe, in
`docs-<target>.md`.
Internal docs (knowledge base, ADRs, glossary, agent instructions) are fab:knowledge's job.

## Inputs

- The target, one of:
  - a **run**, chosen and bounded as in
    [Looking back at a run](../../docs/conventions.md#looking-back-at-a-run); `stories.md` and `spec.md` say what
    users can now do, and its changes say what changed;
  - an **audit**: a folder or the whole repo's public docs, checked against the code.
- The run for an audit: the caller's run when another step or a recipe calls this (see [Writing during a
  step](../../docs/conventions.md#writing-during-a-step)); called (directly or not) by a review with no run, skip the
  file and the ledger and hold both gates in chat. Otherwise start or join one.
- The repo's public docs: README files, a docs site or guides folder, API reference (including doc comments on
  exported code), the changelog, and examples. Its existing style (headings, tone, changelog format) is the model.
- Preferences: `~/.config/fab/preferences.md`. Anything it settles is not asked (picks excepted).
- Personal override: if `~/.config/fab/skills/docs.md` exists, read it and apply it on top of this skill (see fab conventions: Personal layer).

## Process

Follow [Writing during a step](../../docs/conventions.md#writing-during-a-step),
[Committing during a step](../../docs/conventions.md#committing-during-a-step), and
[Picks](../../docs/conventions.md#picks): save `docs-<target>.md` after every exchange.

1. **Find what's wrong.** Work on the branch holding the run's commits (an audit: the current branch), at `HEAD` with a
   clean tree; ask otherwise. Check every public doc's claims (commands, flags, options, outputs, examples, API
   signatures and errors) against the code: for a run, what its changes could affect; for an audit, everything in the
   area; an audit also checks that the README says what the project is, how to install and run it, its commands, how
   it's put together at a glance, and how to contribute, and that exported APIs document parameters, returns, and
   errors; whatever an audit finds missing is proposed at the picks gate, never added directly. A doc that's right while
   the code is wrong is a bug: report it and offer fab:debug; this skill only edits docs.
2. **Find what's missing** (run target only): something a user can now do or must now know (a new flag, endpoint,
   setting, behavior change, or a breaking change and how to move past it) that no public doc covers. Code the run
   didn't touch is out of scope.
3. **Write it** in the repo's own style, where its readers already look. Fixing something wrong is always direct, even
   when the fix adds text, and so is a run's addition inside an existing section. A new page, a new top-level README
   section, or anything an audit found missing waits for picks. Additions go in the section that covers the area; doc
   comments go next to the exported code, with parameters, returns, and errors; a run's user-visible change gets a
   changelog entry when the repo keeps a changelog (never start one). Say what a user does and sees, with a real
   example. Then compare each edited doc with its original, the code, and the run's spec and decisions: cut filler, keep
   every fact, number, and caveat, and flag any claim nothing supports.
4. **Prove it.** Run every command and example you wrote or touched, and every one whose shown output you wrote or
   touched, in a scratch copy when it changes local state, and compare the output with the doc. Anything that
   can't run here (it needs a service, credentials, or a device), reaches outside the scratch copy (a deploy, an
   email, a payment, a write to the home directory or an absolute path), or whose reach you can't tell isn't run,
   and neither is anything when the sandbox blocks the scratch copy: it's marked unverified in the record, never
   presented as checked. A run whose output doesn't match is fixed and rerun before committing, or, when the doc
   is right, reported as a bug with the actual output.
5. **Apply and propose.** Fixes, and a run's additions to existing docs, land as one commit. Then:

   **GATE:** present `docs-<target>.md` with the proposals (new pages, new top-level README sections, and what an
   audit found missing), wait for approval. This is the picks gate.
6. **Write the picks** as in step 3, each in its own commit, proved as in step 4.

**GATE:** present `docs-<target>.md` with every commit, wait for approval. Offer it as a choice (approve, or name
what to change), the way fab:grill asks. When another step or a recipe called this in a run, the record names its
caller and goes to the caller's gate instead.

## Output

Commits, plus `.fab/<run>/docs-<target>.md`:

- **Target:** the run or area, the caller if any, and the commits checked.
- **Docs checked:** each public doc in scope: accurate, fixed (with SHA), or extended (with SHA).
- **Added:** what a user now needs that was missing, and where it went.
- **Verified:** each command or example run, with its result; **Unverified:** each that couldn't run, and why.
- **Bugs reported:** docs that were right while the code was wrong.
- **Proposed:** new pages or sections, and the **Picks**.
- **Pending:** while it runs, anything waiting on the user.

## Done when

- No public doc in scope says something the code contradicts (or the mismatch is reported as a bug), and everything
  a user now needs is documented or was proposed and declined.
- Every command and example written or touched ran and matched, or is listed as unverified with a reason.
- In a run, the ledger shows `docs` rows for `docs-<target>.md` (not a chain step, so its rows never decide a run's
  next step): `started`, `draft` and `started` around the picks gate when it runs, `draft` at the final gate, then
  `approved` with the mtime, or `skipped` with a `gate` note when a caller's gate or a recipe settles it.

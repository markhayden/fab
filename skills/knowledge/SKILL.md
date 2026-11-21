---
name: knowledge
kind: skill
description: Keep a repo's internal docs true (knowledge base, AGENTS.md, ADRs, glossary, committed specs and plans) and carry a run's lasting decisions into them. Use after a run, or when the user says "is the KB stale", "check the internal docs", "write an ADR", or "add to the glossary".
---

# Knowledge

Keeps what agents and maintainers read about a repo true, and keeps the why behind decisions from vanishing with
the run (`.fab/` is never committed). Recorded in `knowledge-<target>.md`. Public, user-facing docs (README, guides,
API reference) are out of scope.

## Inputs

- The target, one of:
  - a **run**: the one named; otherwise offer the newest closed runs that hold work, plus any run with an
    unfinished `knowledge-*.md` (the record names its target), and ask which (as in
    [Starting or joining a run](../../docs/conventions.md#starting-or-joining-a-run)). Its diff runs from the earliest
    starting SHA in the records of its work (build, test, debug, simplify, review; for a run that is only a bug fix,
    the fix commit's parent) to the latest commit they made (a review's only when it committed fixes). Its decisions come from `grill.md`, `spec.md`,
    `design.md`, `plan.md`, `build.md` deviations, and any `debug-*.md` root causes. With no starting SHA to find,
    ask for the range or check its decisions only. A target that already has an approved or gate-skipped record, in
    any run, is shown it; another pass happens only if the user asks;
  - an **audit**: a folder or the whole repo;
  - a **decision** the user names ("write an ADR for …"), or a **term** ("add X to the glossary").
- Where the record goes: called by another step or a recipe, always the caller's run, whatever the target (see
  [Writing during a step](../../docs/conventions.md#writing-during-a-step)); called by a review with no run, skip the
  file and the ledger and hold both gates in chat. Otherwise a run target's record goes in that run, and any other
  target starts or joins one. Either way, keep `.fab/` out of git as in
  [Keeping .fab/ out of git](../../docs/conventions.md#keeping-fab-out-of-git).
- The repo's internal docs: its agent instructions file (`AGENTS.md` or `CLAUDE.md`), knowledge or notes folders,
  ADRs, the glossary, and committed specs, plans, and architecture notes.
- Preferences: `~/.config/fab/preferences.md`. Anything it settles is not asked, except the picks at step 5.
- Personal override: if `~/.config/fab/skills/knowledge.md` exists, read it and apply it on top of this skill (see fab conventions: Personal layer).

## Process

Follow [Writing during a step](../../docs/conventions.md#writing-during-a-step): save `knowledge-<target>.md` after
every exchange and around every commit (note it in progress, then its SHA). An unfinished record for the target is
resumed first, where it lives, after checking the git log against it so a commit that landed isn't redone; a later pass is
`knowledge-<target>.2.md`, `.3.md`, and so on. Each commit holds only this skill's changes (staging just its own
hunks when a doc also holds other edits, and asking when they share a line; if the index already holds other
work, ask first) and passes the repo's
checks and hooks first; a failure rolls back only those changes. When the sandbox blocks a write, check, hook, or
commit, say so and ask; if it stays blocked, stop and report in chat what landed.

Read [references/records.md](references/records.md) before writing an ADR, a status line, or a glossary entry.

1. **Map the docs.** List the internal docs and sort each as **living** (describes how things are now: a KB page,
   architecture notes, the instructions file) or **dated** (a record of a moment: an ADR, a dated or numbered spec
   or plan). Find the repo's ADR convention (directory, filename pattern, file type, numbering, headings) from
   existing ADRs, the instructions file, and any ADR tooling or config; if they conflict, ask. With none, use the
   references file.
2. **Find drift.** Check each doc's claims (paths, commands, names, behavior) against the code: for a run, every
   doc its diff could affect; for an audit, every doc in the area, using history too (a doc last changed before the
   code it describes is a suspect, confirmed by reading; record "can't tell" when history can't order them). Drift that only uncommitted code shows is asked about,
   not marked. This skill only edits docs: when the doc is right and the code is wrong (or the repo's
   rules say the doc wins), report it as a bug and offer fab:debug; never change code here. A
   living doc that's wrong is fixed, adding a dated line saying why to its own change log if it keeps one. A dated one is never
   rewritten: its status line changes (step 5).
3. **Sort the decisions** (run or decision target). For each (a decision the user named as an ADR that fails
   the ADR test is asked about, not quietly turned into something else):
   - an **ADR** when it's hard to reverse, surprising without context, and the result of a real trade-off (all
     three);
   - a **KB note** when it's a fact an agent would otherwise have to rediscover (how something works, where
     something lives): an edit to an existing page is a fix, a new page is proposed at step 5, and the instructions
     file is used only when the repo has no knowledge folder;
   - otherwise it stays in the run, said so in the record.
4. **Glossary.** Each project-specific term a run settled, or the term the user named (only that one), gets one
   canonical entry with the words to avoid: in `GLOSSARY.md`, created on first use, or in the area's glossary when
   the repo keeps a `GLOSSARY-MAP.md` (ask when the area is unclear). An audit only corrects existing entries; new terms it finds are proposed at
   step 5. Two words for one concept, or one word used for two, is asked about, never settled silently.
5. **Apply and propose.** Fixes to living docs, glossary changes, and status lines for dated records with no
   successor ("outdated as of <SHA>: <what changed>") land as one commit, using the repo's commit types. Then:

   **GATE:** present `knowledge-<target>.md` with the proposed new ADRs, KB pages, and (from an audit) glossary
   terms numbered (what, why, where it goes, which dated record it supersedes), wait for approval. The user approves by naming which to write
   (numbers, all, or none); save the picks first. With nothing proposed, skip this gate. This is the first of two
   gates (Writing during a step: a step with several gates), and it stays with this skill even when it's called.
6. **Write the picks** in the repo's convention (ADR numbers continue the sequence), each in its own commit; a new
   record that replaces a dated one changes that one's status to point to it in the same commit. A dated record
   whose proposed replacement wasn't picked gets an "outdated as of" status in one last commit.

**GATE:** present `knowledge-<target>.md` with every commit, wait for approval. Offer it as a choice (approve, or
name what to change), the way fab:grill asks; a new record the user names that wasn't picked goes back through the
picks gate. When another step or a recipe called this in a run, the record names its caller and this final gate
goes to the caller's gate instead.

## Output

Commits, plus `.fab/<run>/knowledge-<target>.md`:

- **Target:** the run, area, decision, or term; the caller if any; the starting SHA and the diff range.
- **Docs map:** each internal doc in scope: living or dated, and checked accurate, fixed (with SHA), marked, or
  can't tell.
- **Decisions:** each one and where it went (ADR, KB note, or stays in the run).
- **Glossary:** terms added or changed.
- **Proposed:** numbered new ADRs, KB pages, and glossary terms, and the **Picks**.
- **Written:** each new record and its SHA.
- **Pending:** while it runs, anything waiting on the user.

## Done when

- Every internal doc in scope is in the docs map as checked accurate, fixed, marked, or can't tell, and no dated
  record was rewritten.
- Every decision has an outcome in the record, and every settled term is in the glossary.
- In a run, the ledger shows `knowledge` rows for `knowledge-<target>.md` (not a chain step, so its rows never
  decide a run's next step): `started`, `draft` and `started` around the picks gate when it runs, `draft` at the
  final gate, then `approved` with the mtime, or `skipped` with a `gate` note when a caller's gate or a recipe
  settles it. A recipe may skip the final gate, never a picks gate with proposals on it.

---
name: spec
kind: skill
description: Turn settled intent, decisions, and user stories into spec.md, with testable success criteria, scope, implementation and testing decisions, and boundaries. Use after grilling, or when the user says "write the spec", "spec this out", or "what exactly are we building?".
---

# Spec

Writes the single description of what will be built and how anyone will know it's done, so design and plan work
from one source instead of the chat. It synthesizes what's settled; it doesn't interview again.

## Inputs

- The run (`.fab/<YYYY-MM-DD>-<slug>/`): start or join one first, as in [Starting or joining a run](../../docs/conventions.md#starting-or-joining-a-run).
- `grill.md` (intent and decisions) and `stories.md` (user stories) from the run, plus anything settled in the conversation.
- The repo: its `AGENTS.md`/`CLAUDE.md`, docs, tests, and the code the change touches.
- The repo's glossary when present (its terms are the spec's terms).
- Preferences: `~/.config/fab/preferences.md`. Anything it settles is not asked.
- Personal override: if `~/.config/fab/skills/spec.md` exists, read it and apply it on top of this skill (see fab conventions: Personal layer).

## Process

Follow [Writing during a step](../../docs/conventions.md#writing-during-a-step).

1. **Load what's settled.** Every decision in `grill.md` and every story in `stories.md` (unless its step was skipped)
   is carried as written; don't reopen one. Read grill's Handed off and Open sections too: settle choices handed to
   `spec` (within step 2's limits), and keep choices handed to later steps with their owner. Implementation decisions
   record what's settled; they don't do design's or plan's work. Ask only when settled requirements conflict with each
   other or with a binding repo constraint. Current behavior that the change replaces is not a conflict.
2. **Fill gaps only if they block.** If either file is missing, use whichever exists, plus the conversation and the
   repo. Ask at most three questions in total, only about gaps that would change a success criterion, using [grill's
   Asking rules](../grill/SKILL.md#asking). If more remain, offer to run fab:grill or fab:stories first. A stories row
   `skipped` with a `step` note (no observable change, or a recipe that skips it) is not missing: cite grill's
   requirements instead. Never assume an answer that changes a success criterion; record any other assumption under
   Assumptions.
3. **Write success criteria** and map each to what it proves: story IDs (`S1, S3`) when `stories.md` exists and
   its step wasn't skipped (`step:` note), otherwise the settled requirement in plain words. Stories under Removed aren't
   requirements. Never invent stories or IDs. Misuse stories become criteria like any other.
4. **Pick test seams.** For each criterion, propose the outermost existing public boundary that can verify it (a
   CLI's output, an HTTP API, a public function); choose per criterion before merging shared seams. If a criterion
   needs a new harness or public interface, ask before writing it in, since that adds work. A harness is new when
   no test in the project drives that boundary yet (for example, the first test that runs the built CLI).
5. **Name what the work relies on.** Each library, framework, or external service API an implementation decision depends
   on is listed with its version (from the lockfile or manifest, or "version unknown") and a link to that version's
   docs, or marked "unverified" for build to check.

**GATE:** present `spec.md`, wait for approval. Offer it as a choice (approve, or name what to change), the way
grill asks.

## Output

`.fab/<run>/spec.md`:

- **Problem:** what's wrong or missing today, from the user's side.
- **Solution:** what changes for them, in a short paragraph.
- **Success criteria:** numbered, each observable by a test or a command, each citing its story IDs or requirement.
- **Scope** and **Out of scope.** Out of scope is never empty.
- **Implementation decisions:** modules, interfaces, data shapes, contracts, behaviors. No file paths or code;
  those go stale. The exception is a short snippet that states a decision more precisely than prose can (a type
  shape, a schema). The APIs relied on, as step 5 lists them.
- **Testing decisions:** a table of criterion → seam → existing or new, plus prior art to copy.
- **Boundaries:** always, ask first, never, for whoever builds it.
- **Assumptions:** what you assumed instead of asking. Empty when everything came from upstream.
- **Open questions:** anything unresolved. One that changes a success criterion blocks approval.
- New project with no `AGENTS.md`: also **Commands** (build, test, lint, run), **Structure**, and **Code style**.

## Done when

- `spec.md` has every required section, and each success criterion is testable and cites a story ID or requirement.
- Every story in `stories.md` (except Removed), when it exists, is covered by at least one criterion.
- No open question changes a success criterion.
- The ledger shows `started`, then `draft`, then `approved` with the mtime (or `skipped` with a `gate` note when a
  recipe skips approval), for `spec.md`.

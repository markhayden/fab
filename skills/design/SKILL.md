---
name: design
kind: skill
description: Design what the spec describes in a connected design tool (such as Paper) or from design files the user provides, then index it in design.md for plan and build. Use after the spec, or when the user says "design this", "mock it up", or "prototype the screens".
---

# Design

Turns the spec's visible parts into a real design before anything is planned, so every screen and state is settled
on a canvas instead of invented during build. It runs only when there's something to design with.

## Inputs

- The run (`.fab/<YYYY-MM-DD>-<slug>/`): start or join one first, as in [Starting or joining a run](../../docs/conventions.md#starting-or-joining-a-run).
- `spec.md` from the run, plus `stories.md` when present and not skipped. Cite active story IDs (`S1`), never
  Removed ones; without stories, cite the spec's numbered criteria.
- A connected design tool (such as Paper), or design files the user provides (mockups, screenshots, exports).
- The project's existing design system: tokens, components, and screens in the repo or the design file.
- Preferences: `~/.config/fab/preferences.md`. Anything it settles is not asked.
- Personal override: if `~/.config/fab/skills/design.md` exists, read it and apply it on top of this skill (see fab conventions: Personal layer).

## Process

Follow [Writing during a step](../../docs/conventions.md#writing-during-a-step).

1. **Check there's something to design.** Skip when either is true, append a `skipped` row, and say so in one line:
   - Nothing in the spec is visual (a CLI flag, an API, a refactor): note `step: nothing visual to design`.
   - No design tool is connected and the user has given no design files: note `step: no design tool`, and say the
     user can call fab:design again after connecting one or handing over mockups. That explicit call reopens the
     step: append a new `started` row despite the earlier skip.
2. **Propose the depth,** from the spec, with a one-line reason, and ask once using
   [grill's Asking rules](../grill/SKILL.md#asking). The minimum is a full-fidelity prototype of every moving part
   (each screen, flow, and state the stories need); the maximum is high-fidelity designs of every visual aspect.
3. **Provided files are the design.** Map each file to the stories and criteria it covers. List what they don't
   show (an error state, an empty list, a narrow screen), and for each gap ask: design it in the tool (if
   connected), describe it in `design.md` under Described, or leave it to build.
4. **Design in the tool** whatever is to be designed there (everything, or the gaps sent to it), at the approved
   depth, following the tool's own instructions. Build on the existing design system; make one only when none
   exists, and say so. Every visual story gets its screens, including the states the spec calls for (empty,
   loading, error).
5. **Designing reveals a spec problem** (a missing state, a flow that can't work): stop and ask. A change to an
   approved `spec.md` sends it back to its gate before you continue.

**GATE:** present `design.md` and where to view the design, wait for approval. Offer it as a choice (approve, or
name what to change), the way grill asks.

## Output

`.fab/<run>/design.md`, an index of the design; the design itself stays in the tool or the provided files:

- **Source:** the design file (name and link) and any provided files (paths).
- **Depth:** prototype or high-fidelity, as approved.
- **Screens:** a table of screen or state → stories and spec criteria it covers (`S1`, `criterion 2`) → where it
  is (artboard, page, file, or Described).
- **Components** and **tokens** used, new ones marked new.
- **Described:** gaps the user chose to settle in words, one short subsection each.
- **Decisions:** choices made while designing, with their source: the user, the spec, the existing design system,
  "from preferences", or your own call (marked, so the user reviews it at the gate).
- **Left to build:** gaps the user chose to leave, each named.

Build takes exact values (spacing, colors, type, markup) from the design tool when it's connected, never from
screenshots.

## Done when

- Every active visual story (or, without stories, every visual criterion) appears in the Screens table or under
  Left to build.
- Every gap in provided files was settled by the user.
- The ledger shows `started`, `draft`, then `approved` with the mtime (or `skipped` with a `gate` note when a recipe
  skips approval), for `design.md`; or a `skipped` row noting `step: nothing visual to design` or
  `step: no design tool`.

---
name: stories
kind: skill
description: Write the user stories for settled work into stories.md, one per behavior someone can do or notice, each with a concrete example. Use after grilling and before the spec, or when the user says "write the stories", "user stories", or "who is this for, exactly?".
---

# Stories

Lists what each person will be able to do or notice once the work ships, so you can check that the spec builds the
right thing before it's written. It synthesizes from what's settled; it doesn't interview again.

## Inputs

- The run (`.fab/<YYYY-MM-DD>-<slug>/`): start or join one first, as in [Starting or joining a run](../../docs/conventions.md#starting-or-joining-a-run).
- `grill.md` from the run (intent, decisions, out of scope), plus anything settled in the conversation.
- The repo, for who uses the thing today and how.
- Preferences: `~/.config/fab/preferences.md`. Anything it settles is not asked.
- Personal override: if `~/.config/fab/skills/stories.md` exists, read it and apply it on top of this skill (see fab conventions: Personal layer).

## Process

Follow [Writing during a step](../../docs/conventions.md#writing-during-a-step).

1. **Find who notices.** From the intent and decisions, list everyone who will do or see something different: users,
   operators, maintainers, CI, other programs. Ask before writing only when who it's for, or what they'd do, is
   unclear, using [grill's Asking rules](../grill/SKILL.md#asking). With no `grill.md`, ask only those questions; if
   the intent itself is unsettled, offer to run fab:grill first.
2. **Nobody notices?** For a change no one can observe (a pure refactor, an internal swap) and that crosses no trust
   boundary (step 3), don't stretch stories to fit. Propose skipping. On approval, append a `skipped` row with the note
   `step: no observable behavior change` and end the step: no `stories.md`, no gate. spec then works from grill's
   requirements.
3. **Write one story per distinct behavior,** including what happens when it goes wrong (bad input, nothing to show,
   can't run). Merge variations of the same behavior instead of listing each one. Never add a story for behavior
   that's out of scope or still unsettled; list it under Not covered instead. Where the work crosses a trust boundary
   (input coming in, data going out, auth, an external fetch, a destructive path), add a misuse story for each
   distinct way it can be abused there: what a hostile or careless actor does, and what happens.
4. **Give each story one concrete example:** real-looking input and what the person sees, specific enough that a
   reader can tell whether a spec covers it.

**GATE:** present `stories.md`, wait for approval. Offer it as a choice (approve, or name what to change), the way
grill asks.

## Output

`.fab/<run>/stories.md`:

- **Stories:** numbered `S1`, `S2`, …, grouped by who it's for. Each is only
  "As a <who>, I want <what>, so that <why>", then one line starting `e.g.` with the example. No Given/When/Then
  or acceptance criteria; fab:spec owns those.
- **Removed:** once approved, IDs never change or get reused, because spec cites them. A dropped story moves here
  with its ID and is no longer a requirement.
- **Not covered:** behaviors someone might expect that this work won't deliver, each with its reason (out of scope,
  or not settled).

## Done when

- Every story has who, what, why, and one example, and every in-scope observable behavior that's settled has a
  story.
- Failure behavior that grill settled has its own story.
- The ledger shows `started`, `draft`, then `approved` with the mtime (or `skipped` with a `gate` note when a recipe
  skips approval), for `stories.md`; or a `skipped` row with `step: no observable behavior change`.

---
name: grill
kind: skill
description: Interview the user until intent and every design decision are settled, then write grill.md for the spec. Use before building anything non-trivial, or when the user says "grill me", "interview me", "stress-test this", or "are we sure?".
---

# Grill

Turns a request into a confirmed intent and a log of settled decisions, so the spec step can write it without
asking again. The cheapest time to find a wrong assumption is before anything is written.

## Inputs

- The user's request, plus anything already settled earlier in the conversation.
- The run: start or join one first, as in [Starting or joining a run](../../docs/conventions.md#starting-or-joining-a-run).
- Preferences: `~/.config/fab/preferences.md`. Anything it settles is not asked.
- Personal override: if `~/.config/fab/skills/grill.md` exists, read it and apply it on top of this skill (see fab conventions: Personal layer).

## Process

Follow [Writing during a step](../../docs/conventions.md#writing-during-a-step): `started` row first, save `grill.md` after every exchange.

1. **Gather facts yourself.** Anything the code, docs, or environment can answer is yours to look up, not a question.
   Use a fresh sub-agent for wide searches if one is available, and keep asking the questions that don't depend on it.
2. **Intent.** Check the request against six fields: outcome, who it's for, why now, success, binding constraint, and
   out of scope. Carry forward whatever the request or conversation already states. Ask about the gaps. Settle a
   shaky premise before any decision that depends on it.
3. **Decisions.** Map the design tree: each decision and the ones that hang off it. Ask the decisions whose
   prerequisites are settled. Recompute after every answer, since one answer can settle or reopen others.
4. **Stop** when someone writing the spec fresh could describe success and scope without guessing at a choice that
   matters. Choices that only affect how the work is done can be handed to a later step instead.

### Asking

- One question per message. Merge questions only when they're naturally answered together.
- The title is the decision itself. Offer options, give each its practical consequence, and mark one
  "(recommended)". If you can't recommend one, say what you'd need to know:

  ```
  Q3 · If the export file already exists, should export stop?
  - Stop and show the path (recommended). The existing file is kept; a rerun needs a new name.
  - Replace it. Reruns are one step; the previous export is lost.
  ```

- "Sounds good" or "whatever you think" settles a decision only if you recommended something; say what you're
  recording. If the user is unsure, offer two concrete options instead of an open question.
- Answers that sound like best practice ("scalable", "clean", "the standard way") aren't decisions yet. Ask what they
  would actually want if they didn't have to justify it.
- "Wait, what?" or "explain that simpler": rephrase the current question in plain words with a concrete example, and
  keep the same number. It may also mean the premise is wrong, so ask.

**GATE:** present `grill.md`, wait for approval.

## Output

`.fab/<run>/grill.md`:

- **Intent:** the six fields, one line each. Out of scope is never empty.
- **Decisions:** numbered. Each has the question, the answer, why, and its source: the user, "from preferences",
  "from the conversation", or a fact you looked up (with where).
- **Handed off:** choices left to a later step, each naming the step that owns it (`spec`, `plan`, `build`).
- **Open:** anything unresolved. An open item that would change what "done" means blocks approval.
- **Pending question:** while the interview is running, the question waiting for an answer.

## Done when

- `grill.md` has Intent, Decisions, Handed off, and Open, and every decision lists its source.
- No open item changes what "done" means.
- The ledger shows `started`, then `draft`, then `approved` with the mtime, for `grill.md`.

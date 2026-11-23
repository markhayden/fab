---
name: grill
kind: skill
description: Interview the user relentlessly until intent and the decisions that matter are settled, then write grill.md for the spec. Use before building anything non-trivial, or when the user says "grill me", "interview me", "stress-test this", or "are we sure?".
---

# Grill

Turns a request into a confirmed intent and a log of settled decisions, so the spec step can write it without
asking again. The cheapest time to find a wrong assumption is before anything is written.

## Inputs

- The user's request, plus anything already settled earlier in the conversation.
- The run (`.fab/<YYYY-MM-DD>-<slug>/`): start or join one first, as in [Starting or joining a run](../../docs/conventions.md#starting-or-joining-a-run).
- `pitch.md`, when the run has one: an unfinished pitch (latest row `started` or `draft`) is finished first,
  with fab:pitch. An approved one whose mtime no longer matches its approval row also goes back through
  fab:pitch for a fresh verdict. Use it only when its latest approved verdict is sold:
  its For grill section seeds the intent and the decisions (marked "from the pitch"), so those points aren't asked
  again, and its open items become questions.
- The repo's glossary (`GLOSSARY.md`, or the area's under a `GLOSSARY-MAP.md`) when present.
- Preferences: `~/.config/fab/preferences.md`. Anything it settles is not asked.
- Personal override: if `~/.config/fab/skills/grill.md` exists, read it and apply it on top of this skill (see fab conventions: Personal layer).

## Process

Follow [Writing during a step](../../docs/conventions.md#writing-during-a-step): `started` row first, save `grill.md` after every exchange.

1. **Gather facts yourself.** Anything the code, docs, or environment can answer is yours to look up, not a question.
   First check whether something already does what's asked (a feature, flag, or tool); if so, lead with it. Use a
   fresh sub-agent for wide searches if one is available, and keep asking the questions that don't depend on it.
   Facts from outside the repo (docs, API behavior, prior art) go in a research file (see
   [Runs and handoff artifacts](../../docs/conventions.md#runs-and-handoff-artifacts)).
2. **Intent.** Check the request against six fields: outcome, who it's for, why now, success, binding constraint, and
   out of scope. Carry forward what the user actually said; ask about the gaps. Settle a shaky premise before any
   decision that depends on it.
3. **Challenge.** Once intent is settled, name the idea's weakest point and one other way to frame it, then ask: keep
   it, change it (once the user says how), or take the other framing, and record the answer as a decision. A change or
   the other framing reopens the intent fields and pitch-seeded decisions it touches, which go back through step 2. From
   here on, point out any answer that conflicts with an earlier one.
4. **Decisions.** Map the design tree: each decision and the ones that hang off it. Ask the decisions whose
   prerequisites are settled. Recompute after every answer, since one answer can settle or reopen others. A decision
   you could infer is still a question: offer your inference as the recommendation, because the interview is where
   gaps in the idea surface. Only what the user said, preferences, the pitch, and facts you looked up are recorded
   without asking.
5. **Sweep** before stopping. For each area not yet covered, ask, or confirm it doesn't apply: who (every kind of
   person who uses or is affected by it), failure and edge cases, data (what's stored, changed, or lost), limits
   (scale, cost, time), rollout (how it reaches people, how it's undone), and what the user would hate to see. Then
   ask "what's still fuzzy for you?" as an open question; each answer becomes a question, a decision, or an Open item.
   Anything an answer here reopens goes back through steps 2–4.
6. **Stop** when every sweep area is covered, that last question left nothing unsettled that changes what "done" means,
   and someone writing the spec fresh could describe success and scope without guessing at a choice that matters.
   Choices that only affect how the work is done can be handed to a later step instead.

### Asking

- One question at a time, numbered (Q1, Q2, …) and saved with its options wherever the step keeps what's pending,
  so it can be asked again exactly. Merge questions only when they're naturally answered together.
- The question, in full with its options, ends your reply; never only a pointer to where it's saved.
- The question is the decision itself. Offer two or three options, each a short label plus its practical consequence,
  with your recommendation first and marked "(Recommended)", then a last option, **Explain these options**, which runs
  fab:likeimfive deeper on the question. If you can't recommend one, say what you'd need to know. Open questions (the
  sweep's last one) have no options.
- If your host has a tool for asking the user a multiple-choice question, use it, so the user picks instead of typing.
  Put any context the user needs (facts you found, why the question matters) in a short message before it. When there's
  no such tool, or it can't fit every option plus the explain option, write it as text with a blank line between
  options:

  ```
  **Q3 · If the export file already exists, should export stop?**

  1. **Stop and show the path (Recommended)**: the existing file is kept; a rerun needs a new name.

  2. **Replace it**: reruns are one step; the previous export is lost.

  3. **Explain these options**
  ```

- With the first question of an interview, say once that "simpler" or "deeper" explains any question. Either word, or
  "wait, what?", runs fab:likeimfive, and the question is asked again under the same number; it doesn't count toward any
  limit on questions. Being lost may also mean the premise is wrong, so ask.
- "Sounds good" or "whatever you think" settles a decision only if you recommended something; say what you're
  recording. If the user is unsure, offer two concrete options instead of an open question.
- Answers that sound like best practice ("scalable", "clean", "the standard way") aren't decisions yet. Ask what they
  would actually want if they didn't have to justify it.
- A term that's fuzzy, conflicts with the glossary, or means two things is asked about whenever it comes up: which
  is meant?
- A question only someone else can answer: write `questions-<who>.md` in the run (who it's for, what's needed back,
  each question aimed at the gap), add it to Pending, and go on with what doesn't depend on it. When nothing else can
  be settled, save and stop; the run resumes when the user brings back the answers.

**GATE:** present `grill.md`, wait for approval. Offer it as a choice (approve, or name what to change) the same way.

## Output

`.fab/<run>/grill.md`:

- **Intent:** the six fields, one line each. Out of scope is never empty.
- **Decisions:** numbered. Each has the question, the answer, why, and its source: the user, "from preferences",
  "from the pitch", "from the conversation", or a fact you looked up (with where).
- **Sweep:** each area, with the intent fields or decision numbers that cover it, or why it doesn't apply; or `skipped:
  <recipe>` when a recipe skips the sweep.
- **Handed off:** choices left to a later step, each naming the step that owns it (`spec`, `plan`, `build`).
- **Open:** anything unresolved. An open item that would change what "done" means blocks approval.
- **Pending:** while the interview is running, every question waiting for an answer, as asked (number and options),
  including questionnaires sent out.

## Done when

- `grill.md` has Intent, Decisions, Sweep, Handed off, and Open; every decision lists its source, and every sweep area
  lists what covers it or why it doesn't apply (or the Sweep says a recipe skipped it).
- No open item changes what "done" means.
- The ledger shows `started`, then `draft`, then `approved` with the mtime (or `skipped` with a `gate` note when a
  recipe skips approval), for `grill.md`.

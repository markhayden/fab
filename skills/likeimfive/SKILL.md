---
name: likeimfive
kind: skill
description: Explain something in plain words with a concrete example, or go deeper on a choice with each option's benefits and drawbacks. Use when the user says "explain like I'm five", "wait, what?", "simpler", "go deeper", "benefits and drawbacks", or "pros and cons".
---

# Likeimfive

Answers "what does that mean?" and "what would each choice cost me?" without derailing the step that raised it. It
writes nothing to the run; the step that asked keeps its question.

## Inputs

- What to explain: the question another step is waiting on (in a fresh session, read it and its options from that step's
  saved record; if they aren't there, ask which question), a term, a piece of code, a decision, or whatever the user
  points at. With nothing named, the last thing you said; with nothing to go on, ask what to explain.
- Depth, from the user's words: **simpler** ("simpler", "like I'm five", "wait, what?") or **deeper** ("deeper", "go
  deeper", "benefits and drawbacks", "pros and cons", or the **Explain these options** choice). Unclear: simpler, then
  offer deeper.
- The repo's glossary (`GLOSSARY.md`, or the area's under a `GLOSSARY-MAP.md`) when present, for terms.
- Preferences: `~/.config/fab/preferences.md`.
- Personal override: if `~/.config/fab/skills/likeimfive.md` exists, read it and apply it on top of this skill (see fab conventions: Personal layer).

## Process

1. **Look it up first.** Anything the code, docs, or environment can show, check, so the explanation uses the
   repo's real names and behavior instead of a general one.
2. **Simpler:** plain words, no jargon; a term you can't avoid gets a one-line meaning. One concrete example from the
   user's own work beats an analogy; use an analogy only when the example isn't enough. End with what it means for
   them in practice. Eight lines at most.
3. **Deeper:** for a choice, each real option (not the explain option itself) as:

   ```
   **1. <label>**
   - Benefits: <what it gets them>
   - Drawbacks: <what it costs, now and later>
   - Pick it when: <the situation it fits>
   ```

   Then whether it's easy to change later, your recommendation, and what would change it; when the asking step offers no
   recommendation (pitch's verdict), give none. The recommended option's drawbacks are never empty. For something that
   isn't a choice (a term, some code): how it works, what it trades off, and where it breaks.
4. **Hand back.** When what you explained was another step's waiting question, ask it again exactly as before, same
   number and same options, since an explanation settles nothing. If explaining showed the question is wrong or missing
   an option, say so; the asking step revises it under the same number.

## Output

The explanation, in the conversation. No file and no ledger row.

## Done when

- Simpler: eight lines or fewer, with a concrete example and no unexplained jargon.
- Deeper on a choice: every option has benefits, drawbacks, and pick-it-when, and the recommended one has a drawback.
- Deeper on anything else: how it works, what it trades off, and where it breaks.
- When what you explained was a waiting question: it was asked again with its number, still unanswered.

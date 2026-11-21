---
name: pitch
kind: skill
description: Pitch an idea to the user as its honest advocate, defend it against their challenges, and record the verdict (sold, killed, or parked) in pitch.md. Use when the user says "pitch me", "sell me on this", "convince me", or asks to be talked into an idea.
---

# Pitch

Grill reversed: the agent makes the case for the idea and the user attacks it. `pitch.md` keeps what survives: a
sharper idea and why it's worth building, or why it isn't.

## Inputs

- The idea, from the user's message or the conversation.
- The run (`.fab/<YYYY-MM-DD>-<slug>/`): start or join one first, as in [Starting or joining a run](../../docs/conventions.md#starting-or-joining-a-run).
  Pitch sits before the chain: it never adds `skipped` rows for chain steps, its rows never decide which chain
  step is next, and grill doesn't need it. Resuming an unfinished pitch: see the Resume rule in conventions. An
  approved pitch whose mtime changed goes back to its gate.
- Preferences: `~/.config/fab/preferences.md`. Anything it settles is not asked.
- Personal override: if `~/.config/fab/skills/pitch.md` exists, read it and apply it on top of this skill (see fab conventions: Personal layer).

## Process

Follow [Writing during a step](../../docs/conventions.md#writing-during-a-step): save `pitch.md` after every exchange.

1. **Get the facts first.** A number you can't confirm is labeled a guess and says what it's based on; one with no
   basis at all stays out.
2. **Pitch the strongest true version,** on one screen:
   - **Problem:** what's wrong today, in one or two lines.
   - **For:** who exactly.
   - **What:** what it does, with one concrete example.
   - **Worth it:** what it saves or makes possible.
   - **Why now.**
   - **Cost:** effort and anything it needs.
   - **Not doing:** what's deliberately left out.
   - **Least sure of:** the one thing that worries you most.
   A field you can't fill honestly says so.
3. **Take challenges one at a time,** answering each directly. When the pitch holds, say why. When a challenge
   lands, concede and revise the pitch, showing only the lines that changed. Ask the user something only when the answer turns
   on what they alone know.
4. **The pitch sinks** (its core premise doesn't hold): say so, then offer one or two pivots that keep what held up,
   each as a one-line pitch. The user picks one to pitch properly, or kills it.
5. **When challenges run dry,** or whenever the user calls it (including at a pivot choice), go to the gate.

**GATE:** present the final pitch in `pitch.md`, wait for approval. Approval is the verdict. If the user already gave one for this version of
the pitch, record it without asking again; otherwise offer sold, killed, or parked as a choice using [fab:grill's Asking rules](../grill/SKILL.md#asking), with no recommendation.

## Output

`.fab/<run>/pitch.md`:

- **Pitch:** the final version, in the fields above.
- **Challenges:** each one, the answer, and whether it changed the pitch (and how).
- **Pivots:** any offered, and which was taken.
- **Verdict:** sold, killed, or parked, and the user's reason if they gave one.
- **For grill:** when sold, what the pitch settled (it seeds grill's intent and decisions) and what's still open.
- **Pending:** while the pitch is running, whatever waits on the user: a question, a pivot choice, or the
  verdict.

## Done when

- `pitch.md` has the final pitch, every challenge with its outcome, and the verdict.
- Every number in the pitch has a source, or a "guess" label with its basis.
- The ledger shows `started`, `draft` when the final pitch waits for the verdict, then `approved` with the mtime and
  the verdict in the note.

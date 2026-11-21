---
name: retro
kind: skill
description: Look back over a finished run, propose environment fixes (checks, standards, pointers, preferences) so the same trouble does not recur, and apply the ones the user picks, in retro.md. Use when the user says "retro", "retrospective", "what went wrong", or "how do we stop this happening again".
---

# Retro

Turns what went wrong in a run into changes to the environment the agent works in, never the product code, so the
same trouble doesn't happen again, and records it in `retro.md`.

## Inputs

- The run (`.fab/<YYYY-MM-DD>-<slug>/`): the one named; otherwise offer the newest closed runs (see [Starting or
  joining a run](../../docs/conventions.md#starting-or-joining-a-run)), plus any run with an unfinished `retro.md`,
  and ask which. A run holding a single bug fix counts. If the run already has an approved `retro.md`, show it;
  another look back happens only if the user asks, in `retro-2.md`, `retro-3.md`, …. Below, "the record" is
  whichever of these this retro writes; an unfinished one is resumed first. Check `.fab/` is kept out of
  git as in [Keeping .fab/ out of git](../../docs/conventions.md#keeping-fab-out-of-git). Its files are the evidence: the ledger; `build.md` (Deviations, stuck stops,
  Noticed); every `review-*.md` finding; every `debug-*.md` (root cause, any "no seam" finding); `tests.md` (bugs,
  weak tests, anything not verified); and questions in `grill.md`/`spec.md` that the user had already settled before.
- The current conversation, when retro runs in the session it looks back on. With no run, it's the only
  evidence: skip the file and the ledger, and hold both gates below in chat (picks before applying anything).
- The repo's environment, read before proposing anything: its agent instructions file (`AGENTS.md` or `CLAUDE.md`), standards and style docs,
  hooks, CI config, linters, and the check commands they run.
- Preferences: `~/.config/fab/preferences.md` (read; changing it is always asked).
- Personal override: if `~/.config/fab/skills/retro.md` exists, read it and apply it on top of this skill (see fab conventions: Personal layer).

## Process

Follow [Writing during a step](../../docs/conventions.md#writing-during-a-step): `started` row first, save
the record after every exchange and around every applied change: note it as in progress before committing, then
add its SHA. A resumed retro checks the git log (and the preferences file) against the record before applying
anything again. A commit that fails has only retro's own edits rolled back (report anything that can't be) and is recorded as
not applied. A preference pick whose approved wording isn't in this session is asked again. When the
sandbox blocks a write, check, or commit, say so and ask; if it stays blocked, stop and report in chat what
landed.

1. **List what went wrong or cost time,** each with where it shows: a deviation, a stuck stop, a review finding, a
   bug, a weak test, a question asked twice, a long search, a failed or expensive command. No evidence, no candidate.
2. **Ask what in the environment would have prevented it,** and turn each into one candidate:
   - **Check:** a mistake with a fixed shape (a banned call, an import rule, a file location, a format) becomes a
     deterministic check: a linter rule, a hook, or a CI job, whichever the repo's tooling makes cheapest. A check
     that exists but isn't wired in, or is broken, is the finding, not a new one. A repo where no pre-commit hook or CI job
     runs any of its test, lint, or typecheck commands is a finding on its own.
   - **Standard:** a judgment call the reviewer caught late or missed becomes a rule in the repo's standards or
     style doc (with none, `CODING_STANDARDS.md` plus a one-line pointer in the agent instructions file); a rule
     that misled or was ignored gets clarified or removed.
   - **Pointer:** time lost finding something becomes a one-line pointer in the agent instructions file to the file
     that holds it.
   - **Steering cleanup:** an instructions-file rule that belongs in a check or the standards doc, or that changes
     nothing an agent does, moves or goes.
   - **Access:** information the agent needed and couldn't reach (logs, a read-only service view), and how to
     reach it.
   - **Tools:** a slow, costly, or failing command, and the cheaper one that does the job.
   - **Preference:** a decision the user made the same way twice or more becomes a proposed entry for
     `~/.config/fab/preferences.md`. Its text is shown in chat only; the record lists it as "preference (in chat)".
   - **For fab:** a fab skill told the agent to do the wrong thing. Listed for the user to take back to fab, never
     edited from here.
3. **Rank** by severity: how much it cost and how likely it recurs. At equal severity, a bug that reached a commit
   beats a review catch, which beats a slow step.
4. **GATE:** present the record with the candidates numbered (evidence, change, effort), wait for approval. The
   user approves by naming which to apply (numbers, all, or none); record the picks. This is the first of two
   gates (Writing during a step: a step with several gates).
5. **Apply the picks,** each in its own commit holding only that change (staging just its hunks when a file also
   holds other edits; if the index already holds other work, ask the user to commit or unstage it first), passing the repo's checks and hooks, using the repo's commit types. A new check must fail on the original mistake (a fixture or a scratch copy) and pass on the current tree;
   never loosen an existing one. A preferences entry is shown word for word and written only after the user
   approves that write, after checking the file doesn't already hold it; it isn't committed anywhere, and
   the record notes only that it was written (the ledger's next row notes the approval). Something that needs the user (a CI secret, a service setting) is listed for them instead.

**GATE:** present the record with what was applied, wait for approval. Offer it as a choice (approve, or name what
to change), the way grill asks. A change the user names that wasn't picked goes back through the picks gate.

## Output

Commits for the applied changes, plus the record (`.fab/<run>/retro.md`, or `retro-<n>.md`):

- **Scope:** the run, the files read, and whether the session itself was included.
- **Candidates:** ranked; each with evidence (file and line, or the moment in the session), category, change, effort.
- **Picks:** the numbers the user chose at the first gate.
- **Applied:** each change with its SHA (or "preferences entry written"), and how a new check was proven.
- **Declined:** candidates the user passed on, with any reason.
- **For fab:** skill problems for the user to raise with fab.
- **Pending:** while it runs, the question waiting for an answer.

## Done when

- Every candidate cites its evidence, and every applied check was shown failing on the original mistake.
- Every applied change is committed or written with approval, and nothing else changed.
- In a run, the ledger shows `retro` rows for the record (retro isn't a chain step: its rows never decide a run's
  next step or reopen a closed run): `started`, `draft` at the picks gate, `started` noting the picks once chosen,
  `draft` at the final gate, then `approved` with the mtime, or `skipped` with a `gate` note when a recipe skips
  that final approval. The picks gate is never skipped: what gets applied is always the user's choice.

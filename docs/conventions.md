# Conventions

The rules every fab skill and recipe follows. `tools/check` enforces everything marked **(checked)**.
If a rule here and the checker disagree, the checker is wrong: fix it to match this file.

## Naming

- A skill's name is its directory name: `skills/<name>/SKILL.md` with `name: <name>`. **(checked)**
- Names are lowercase letters, digits, and hyphens; they start with a letter; ≤32 chars. **(checked)**
- The `fab:` plugin namespace is added by the host, so names carry no prefix.
- Invocation is identical across hosts except the sigil: Claude `/fab:<name>`, Codex `$fab:<name>`.
- **Skills** are named for the action they perform or what they produce: `grill`, `stories`, `spec`, `design`, `plan`, `build`, `test`, `voice`.
- **Recipes** that run the build flow are named `<size>dev`: `bigdev` (full ceremony), `tinydev` (light).
  Other recipes get a plain name for the kind of work (e.g. a future `spike`).
- Every new skill or recipe gets a **name review** in its grill session, before it ships.
- Renaming a skill or recipe also renames its override file (`~/.config/fab/skills/<name>.md`) and is a
  **major** version bump.

## Kinds

Every `SKILL.md` declares its kind in frontmatter. **(checked)**

| `kind` | What it is | Invocation |
|---|---|---|
| `skill` | Building block: one job, owns its process | Explicit or automatic |
| `recipe` | Orchestrator: sequences skills with gates; almost no process of its own | **Explicit only** |

A recipe is call-only in **both** hosts, with these exact values: **(checked)**

- Claude, in `SKILL.md` frontmatter: `disable-model-invocation: true`
- Codex, in `skills/<name>/agents/openai.yaml`:
  ```yaml
  policy:
    allow_implicit_invocation: false
  ```

## Skill shape

```markdown
---
name: <name>
kind: skill
description: <≤300 chars: what it does + the phrases that should trigger it>
---

# <Title>

<1–2 lines: the artifact it produces and why it exists>

## Inputs
## Process
## Output
## Done when
## Traps        (optional)
```

- Required sections for `kind: skill`: **Inputs, Process, Output, Done when**. **(checked)**
- Required sections for `kind: recipe`: **Steps, Rules**. **(checked)**
- `description` is non-empty and ≤300 characters. **(checked)** It is loaded into every session, so it
  carries what the skill does and when to use it. Nothing else.
- Aim for ~100 lines. Over 150 lines is a warning (the check still passes). **(checked)** Move detail into
  `skills/<name>/references/*.md`, and have `SKILL.md` say exactly when to read each file.
- Every relative link and `references/` path must resolve. **(checked)**

### Writing the body

- Write only what the model wouldn't do on its own: fab's opinions, file locations, gates, thresholds.
  Don't teach general engineering.
- Give the reason for a rule in one clause instead of capitalizing it. Keep strict ordering for gates.
- Add a **Traps** entry only for a failure you've actually seen.
- "Done when" items are evidence (output, path, commit SHA), never "I checked".
- Stay neutral between hosts: never name host-specific tools or features in a skill body. Say "ask the
  user", "use a fresh sub-agent if available". **(checked)**: see Forbidden strings. The one exception is a
  `references/` file that gives the commands for calling a host's CLI from the other host (`review`'s
  cross-host reviewer); the skill body still names neither host.
- Refer to another skill as `fab:<name>`, with no host sigil (`/` or `$`), so one body works in both hosts. Every
  `fab:<name>` under `skills/` must name an existing skill. **(checked)**

## Gates

A gate is a point where the agent stops for approval. Always written exactly as:

```markdown
**GATE:** <what to present>, wait for approval.
```

Every gate decision goes in the run's `ledger.md`. A recipe's Steps table says, for each gate, whether it's
kept or skipped.

## Personal layer

Optional, private files that live outside the plugin, so plugin updates never touch them:

| Path | Purpose |
|---|---|
| `~/.config/fab/voice.md` | Voice profile read by `voice` |
| `~/.config/fab/preferences.md` | Personal rules that apply across skills (e.g. "no compatibility shims on personal projects") |
| `~/.config/fab/skills/<name>.md` | Per-skill override: extra instructions for one skill or recipe |

**Precedence:** per-skill override → `preferences.md` → fab defaults.

**Limits:** an override may change *how* a skill works. If it conflicts with a gate a recipe promised, or with a
required handoff artifact, the agent surfaces the conflict and asks rather than applying it silently.

A missing file is normal; continue silently. A file that exists but can't be read gets a one-line warning and
the run continues on defaults. Creating any file under `~/.config/fab/` is a write: ask first.

Write overrides as outcomes ("always include a rollback section"), not step references ("in step 3"), so they
keep working when fab changes.

### The override line

Every skill and recipe includes this line, verbatim, under **Inputs** (recipes: under **Rules**): **(checked)**

```markdown
- Personal override: if `~/.config/fab/skills/<name>.md` exists, read it and apply it on top of this skill (see fab conventions: Personal layer).
```

with `<name>` replaced by the skill's own name. Recipes pass each child skill's override path along when they
run that skill, including in a fresh sub-agent.

## Runs and handoff artifacts

Workflow skills hand off through files in the target project, not through chat context.

- **Run directory:** `<project>/.fab/<YYYY-MM-DD-slug>/`, where `<project>` is the git top-level directory, or the
  current directory outside git. Resuming takes the slug or the path.
- **Artifacts:** each workflow skill writes one file into the run directory (the set is defined by the chain
  skills: `grill.md`, `stories.md`, `spec.md`, `design.md`, `plan.md`, `build.md`, `tests.md`; plus one
  `review-<target>.md` per review, one `debug-<bug>.md` per bug, `retro.md` when a run is looked back on (numbered after the first), and `pitch.md` when an idea was pitched
  first). Standalone utility skills (like
  `voice`) write no project state.
- **`ledger.md`:** append-only table, one row per event:

  ```markdown
  | When (UTC) | Step | Artifact | Status | Note |
  |---|---|---|---|---|
  | 2026-11-20T06:40:00Z | spec | spec.md | approved | mtime 1763620800.123456789 |
  ```

  - **Status** is one of `started`, `draft`, `approved`, `skipped`.
  - A `skipped` row's note starts with `gate` (the work was done, approval skipped) or `step` (the work isn't owed in
    this recipe).
  - An `approved` row records the artifact's full-resolution mtime (`stat -f %Fm` on macOS, `stat -c %.9Y` on Linux).

- **Resume:** an unfinished pitch (its latest row `started` or `draft`) comes first: continue it from its Pending
  item, or re-present its verdict gate. Otherwise the next action is the first step with no `approved` or
  `skipped` row. A `draft` artifact is
  re-presented at its gate, not regenerated. If an approved artifact's current mtime differs from the one
  recorded at approval, re-present it for confirmation before continuing; for `pitch.md`, that means a fresh
  verdict at its gate.
  Records outside the chain (`review-<target>.md`, `debug-<bug>.md`, `retro.md`) whose latest row is `started` or `draft`
  never decide the next step: name them when resuming the run, and continue one when the user picks it or calls
  its skill.

### Starting or joining a run

Every workflow skill does this first, whether a recipe runs it or the user calls it directly. Work can start anywhere,
and nothing approved is ever redone.

1. **A run was named** (slug or path): use it.
2. **Otherwise, look for open runs** in `<project>/.fab/`: a run is open until its last chain step (`test`) has an
   `approved` or `skipped` row, so a run started by one standalone skill stays open for a recipe to pick up. A run
   whose latest pitch row is an `approved` killed or parked verdict is closed too; naming it reopens it. A run whose rows
   are all `debug` or `review` rows (a bug fixed on its own) is closed once each record's latest row is `approved`.
   `retro` rows never count toward whether a run is open: a retro never reopens or holds open a run. Exactly
   one: if the user's message clearly continues it (for example, answers its pending question), continue it and say
   which run; otherwise offer it by name and ask whether to continue it or start fresh. Several: list the newest five and ask which.
   Never pick one silently. (`retro` looks back, so it offers closed runs instead, and never creates one.)
3. **No run to join:** create one named `<today's date>-<short slug>`, e.g. `.fab/2026-11-20-export-retry/`; the date
   prefix keeps runs in order. Show the path in your first message and let the user rename the slug. Then make sure `.fab/` is kept out of git (below).
4. **Joining mid-chain:** an earlier step with no row is not run. Append a `skipped` row with the note
   `step: started later in chain`.
5. **No run at all** is allowed for `review` and `retro` alone, when the work has no run (someone else's branch,
   or a session with no run): the record goes in chat.
6. **Earlier chat counts:** if the work was already discussed before the run existed, carry what was settled into the
   artifact (marked "from the conversation") instead of asking again.

### Writing during a step

- Append a `started` row when the step begins.
- For interactive steps, save the artifact after every exchange, including the pending question, so an interrupted
  session resumes where it stopped.
- When the artifact is complete and waiting at its gate, append a `draft` row. On approval, append an `approved` row
  with the mtime.
- A step with several gates (build's checkpoints) appends a `draft` row at each one, and a `started` row noting
  the gate when an intermediate one is approved. Only its last gate gets the `approved` row. Resuming with a
  `draft` as the step's last row re-presents that gate.
- A record made by a skill another step called (`review`, `debug`) uses the caller's run, names its caller, and goes
  to the caller's gate. Settling that gate (approved or skipped) appends the same status for each called record
  presented there, and for any record those called in turn, approved ones with their own mtime. With no run (a review of someone's branch), the called record goes in chat too.
- A decision settled by `~/.config/fab/preferences.md` is not asked. Say so in the chat in one line ("From your
  preferences: …") and mark it "from preferences" in the artifact, so the user can overrule it at the gate.

### Keeping `.fab/` out of git

Whenever a skill uses `.fab/`, make sure `.fab/.gitignore` exists and contains the single line `*`. Git then
ignores the whole folder, that file included. It needs no write to `.git/` (which sandboxes often block), works in
linked worktrees, and is harmless outside git.

Never edit a tracked `.gitignore` for this.

## Forbidden strings

`tools/check` fails if any file under `skills/` contains these (case-insensitive). **(checked)** Docs, tooling, and
test fixtures may name them.

| String | Why |
|---|---|
| `agent-skills:` | reference to the old pack's slash commands |
| `waffleoff` | the old global command fab replaces |
| `AskUserQuestion`, `TodoWrite`, `EnterPlanMode`, `ExitPlanMode` | Claude-only tool names |
| `request_user_input`, `update_plan` | Codex-only tool names |

## Versioning

- The version lives in `.claude-plugin/plugin.json` and `.codex-plugin/plugin.json`, and they must match. **(checked)**
- Bump on every shipped change: **patch** for wording, **minor** for a new skill or recipe, **major** for a rename or
  removal.
- Both catalogs (`.claude-plugin/marketplace.json`, `.agents/plugins/marketplace.json`) point at the repo root and list the
  plugin `fab`. **(checked)**

## Commits

Conventional Commits: `<type>(<scope>): <summary>`, with types `feat`, `fix`, `docs`, `chore`, `test`. No co-author or
AI-attribution trailers.

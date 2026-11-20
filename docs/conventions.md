# Conventions

The rules every fab skill and recipe follows. `tools/check` enforces everything marked **(checked)**.
If a rule here and the checker disagree, the checker is wrong: fix it to match this file.

## Naming

- A skill's name is its directory name: `skills/<name>/SKILL.md` with `name: <name>`. **(checked)**
- Names are lowercase letters, digits, and hyphens; they start with a letter; ≤32 chars. **(checked)**
- The `fab:` plugin namespace is added by the host, so names carry no prefix.
- Invocation is identical across hosts except the sigil: Claude `/fab:<name>`, Codex `$fab:<name>`.
- **Skills** are named for the action they perform: `grill`, `spec`, `plan`, `build`, `test`, `voice`.
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
  user", "use a fresh sub-agent if available". **(checked)**: see Forbidden strings.

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

- **Run directory:** `<project>/.fab/<YYYY-MM-DD-slug>/`. The slug is named when the run starts; resuming takes
  the slug or the path.
- **Artifacts:** each workflow skill writes one file into the run directory (the set is defined by the chain
  skills: `grill.md`, `spec.md`, `plan.md`, `build.md`, `tests.md`). Standalone utility skills (like
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

- **Resume:** the next action is the first step with no `approved` or `skipped` row. A `draft` artifact is
  re-presented at its gate, not regenerated. If an approved artifact's current mtime differs from the one
  recorded at approval, re-present it for confirmation before continuing.

### Keeping `.fab/` out of git

On first use in a project:

1. Inside a git work tree: `git rev-parse --git-path info/exclude` gives the exclude file (this is correct in linked
   worktrees too). If it doesn't already contain the line `/.fab/`, append it. Print a one-line notice.
2. Outside git: do nothing, print a one-line notice.
3. Exclude file not writable: warn in one line and continue.

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

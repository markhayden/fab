# Cross-host reviewer commands

Read this when picking the reviewers in `review`. The reviewers are every host below except the one you're in, whose
CLI answers `<cli> --version`. Run each from the repo root, read-only, one fresh process per reviewer per round.

Each round, write each brief to its own temporary file and pass it as the prompt, so quoting can't break it, and give
each reviewer (host and axis) its own output file. For a diff target, save the diff and commit list as they are now to
a file in the briefs' directory (`$BRIEF_DIR` below) and point the briefs at it; an artifact target is pointed at
directly. Each directory holding
a file a brief names outside the repo gets one `--add-dir` on the Claude and Antigravity commands. Use the system temp
dir; if writes there are blocked (a sandbox), use `.fab/<run>/.review-tmp/` in a run (git already ignores it), or ask
where to write with no run.

## Models

Each host reviews on its most capable model at high effort. Checked 2026-11-25; at the start of each review, run the
lookup for each host you'll use. When the lookup's own test (most capable; recommended for deep verification) picks a
different model than the table, use that one and note it in the record; a model that's only newer doesn't count. When
a lookup fails, use the table's model.

| Host | CLI | Model | Lookup |
|---|---|---|---|
| Codex | `codex` | `gpt-6.2-astra`, effort `high` | `codex debug models`: the one described as most capable |
| Claude Code | `claude` | `opus` (the newest Opus), effort `high` | none needed: the alias follows new releases |
| Antigravity | `agy` | `gemini-4-carbon-high` | `agy models`, and the model <https://antigravity.google/docs/models/> recommends for deep verification |

Preferences may name other models, and win over both the table and the lookup.

## Codex

```sh
codex exec --skip-git-repo-check -s read-only --ignore-user-config -m gpt-6.2-astra \
  -c model_reasoning_effort="high" "$(cat "$BRIEF")" </dev/null > "$OUT" 2>"$OUT.err"
```

- `-s read-only` sandboxes the reviewer's shell commands. MCP servers run outside that sandbox, so
  `--ignore-user-config` skips the user's config (and the MCP servers it starts) while keeping auth. Only Codex's
  built-in app connectors remain, and they can't reach local files.
- `--ignore-user-config` also drops the user's model, so `-m` and the effort must be passed; without them Codex
  reviews on its default model at low effort (seen 2026-11-25).
- Only the final message is printed. Add `--json` to see its steps when debugging.

## Claude Code

```sh
claude -p --restricted --strict-mcp-config --permission-mode plan --tools "Read,Grep,Glob" \
  --model opus --effort high --add-dir "$BRIEF_DIR" -- "$(cat "$BRIEF")" </dev/null > "$OUT" 2>&1
```

- `--restricted` ignores the user's and project's settings (so no hooks or plugins run) and removes tools that run
  code; `--strict-mcp-config` with no config loads no MCP servers; `--tools` leaves only the three reading tools
  (no shell, since even `git diff --output` can write). Verified 2026-11-21: the reviewer could read the diff file
  and could not create one.
- `</dev/null` matters: with stdin open but empty, it waits 3 seconds and warns.

## Antigravity

```sh
agy --mode plan --sandbox --model gemini-4-carbon-high --add-dir "$BRIEF_DIR" \
  --settings '{"permissions":{"allow":["command(git diff)","command(git show)","command(git log)","command(git status)","command(git ls-files)","command(ls)","command(cat)","command(grep)","command(head)","command(tail)","command(wc)"]}}' \
  -p "$(cat "$BRIEF")" </dev/null > "$OUT" 2>"$OUT.err"
```

- `--mode plan` is what keeps it from editing: without it, headless `agy` runs in auto mode and edits files (seen
  2026-11-25; `agy` ignores a mode it doesn't know, such as `default`). Plan mode denied a write in the same test.
- Headless `agy` can't prompt, and a denied tool ends the run with no output and exit 0. The allow list covers the
  reading commands a reviewer reaches for; never add one that writes (`sed`, `tee`, `cp`). Plan mode still blocks
  writes through allowed commands: in a 2026-11-25 test, `git show --output=<file>` and `grep … > <file>` wrote
  nothing, while `git log` ran. `git show <rev>:<path>` was denied, which is why diffs go in a file. Leave out tools installed outside the
  system (Homebrew's `rg`): the sandbox blocks them, and the agent's retry outside the sandbox is denied, ending the
  run.
- Never add `--disable-slash-commands`: it turns plan mode off (the CLI warns so).
- `-p` takes the next argument as the prompt, so every flag goes before it.
- Add one line to its briefs: "Read files with your file tools; run no shell commands." A command outside the allow
  list ends the run with no output, which happened on both first attempts in a 2026-11-25 test; the allow list is
  the backup.
- There's no flag to skip the user's Antigravity settings, MCP servers, or plugins, so the brief says not to run
  skills or plugins. A startup hook would run before the brief, so check first: `agy plugin list --json` for enabled plugins with
  hooks, `~/.gemini/config/hooks.json`, and the repo's `.agents/hooks.json`. Any hook, or a check that fails, means
  skipping Antigravity as a reviewer and noting why.
- A reviewer on Carbon (High) took 3 to 4 minutes.

## From each host

- **Codex:** inside Codex's sandbox, the other CLIs can't reach the network or their logins and fail authentication
  (Claude seen 2026-11-21). Run them outside the sandbox, which needs the user's approval; if that's refused, fall
  back.
- **Antigravity:** its sandbox blocks binaries installed outside the system, and the other CLIs need the network and
  their logins, so the `--version` checks, the model lookups, and the reviewers run outside the terminal sandbox, which
  needs the user's approval; if that's refused, fall back. Headless (`agy -p`), a call that needs approval ends the run,
  so a headless review there uses the fallback from the start and notes why.

## Every round

- Start every reviewer at the same time, as background processes, and wait for all of them.
- A round can take several minutes. Don't treat slowness as failure. A non-zero exit, or output that neither lists
  findings nor says there are none (empty, only a plan, or a note that it couldn't read the target), is a failure:
  retry once.

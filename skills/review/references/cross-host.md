# Cross-host reviewer commands

Read this when picking the reviewer in `review`. The reviewer is the *other* host's agent: from Claude Code, use
Codex; from Codex, use Claude Code. Run it from the repo root, read-only, one fresh process per reviewer per round.

Write the brief to a temporary file first (outside the repo, e.g. the system temp dir) and pass it as the prompt, so
quoting can't break it.

## From Claude Code: Codex

```sh
codex exec --skip-git-repo-check -s read-only --ignore-user-config "$(cat "$BRIEF")" </dev/null > "$OUT" 2>"$OUT.err"
```

- `-s read-only` sandboxes the reviewer's shell commands. MCP servers run outside that sandbox, so
  `--ignore-user-config` skips the user's config (and the MCP servers it starts) while keeping auth. Only Codex's
  built-in app connectors remain, and they can't reach local files.
- `</dev/null` matters: with stdin open but empty (a background shell), `codex exec` waits forever.
- Only the final message is printed. Add `--json` to see its steps when debugging.
- Check it's available with `codex --version`; a non-zero exit or an auth error means fall back.

## From Codex: Claude Code

```sh
claude -p --restricted --strict-mcp-config --permission-mode plan --tools "Read,Grep,Glob" \
  --add-dir "$TMPDIR_WITH_DIFF" -- "$(cat "$BRIEF")" > "$OUT" 2>&1
```

- `--restricted` ignores the user's and project's settings (so no hooks or plugins run) and removes tools that run
  code; `--strict-mcp-config` with no config loads no MCP servers; `--tools` leaves only the three reading tools
  (no shell, since even `git diff --output` can write). Save the diff and commit list to a file next to the brief
  and point the brief at it. Verified 2026-11-21: the reviewer could read the diff file and could not create one.
- Headless Claude can only read inside the current directory. Pass `--add-dir <path>` for anything the brief points
  to outside the repo (for example a run artifact stored elsewhere).
- Check it's available with `claude --version`.
- Inside Codex's sandbox, Claude can't reach the network or its login, so it fails authentication (seen
  2026-11-21). Run this command outside the sandbox, which needs the user's approval; if that's refused, fall back.

## Both

- Start the spec and standards reviewers at the same time, as two background processes, and wait for both.
- Every review has both reviewers, every round.
- A round can take several minutes. Don't treat slowness as failure; treat a non-zero exit or empty output as
  failure, and retry once before falling back.

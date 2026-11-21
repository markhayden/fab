# Development

How to work on fab itself. Installing fab as a plugin needs none of this.

## Setup

Go 1.24+ (development only).

```sh
git config core.hooksPath .githooks   # pre-commit runs tools/check
```

## Validate

```sh
go run ./tools/check   # every rule marked (checked) in docs/conventions.md
go test ./...          # validator tests; fixtures are generated in temp dirs, plus tools/check/testdata
```

The check exits 1 on errors and 2 if it can't run. A skill over 150 lines warns but passes.

## Local test loop

Claude, without installing:

```sh
claude --plugin-dir ~/go/src/github.com/markhayden/fab
```

Installed from a local checkout (both hosts). Plugin caches copy the files, so reinstall after each change:

```sh
claude plugin marketplace add ~/go/src/github.com/markhayden/fab
claude plugin uninstall fab@fab; claude plugin install fab@fab

codex plugin marketplace add ~/go/src/github.com/markhayden/fab
codex plugin remove fab@fab; codex plugin add fab@fab
```

Headless runs, useful for live tests:

```sh
claude -p '/fab:voice <ask>'
codex exec --skip-git-repo-check '$fab:voice <ask>' </dev/null
```

MCP servers that a launcher passes to an interactive session (not in Claude's config) are missing in `claude -p`. For
the design skill, pass Paper explicitly: `--mcp-config '{"mcpServers":{"paper":{"type":"stdio","command":"~/.paper/bin/paper","args":["mcp"]}}}' --allowedTools mcp__paper`.

`codex exec` prints only the final message; add `--json` to see intermediate messages and commands. It also waits on
stdin when stdin is open but empty (a background shell), so redirect from `/dev/null`.

Verified on Claude Code 2.1.329 and Codex CLI 0.170.0.

## Versioning

The version lives only in `.claude-plugin/plugin.json` and `.codex-plugin/plugin.json`; the check fails if they differ.
Bump on every shipped change: patch for wording, minor for a new skill or recipe, major for a rename or removal.

The bump is what makes Claude update. Tested 2026-11-20 with a fixture plugin (Claude from a local-path marketplace,
Codex from a git URL), changing content at the same version, then with a bump:

| Host | Same version, new content | Bumped version |
|---|---|---|
| Claude Code 2.1.329 | Not picked up: `plugin update` says "already at the latest version" | Installed into a new `cache/<version>/` dir (old one kept); restart to apply |
| Codex 0.170.0 | Picked up by `marketplace upgrade` (cache replaced in place) | Picked up; cache dir renamed to the new version |

Claude only accepts `owner/repo`, `https://`, or local-path marketplace sources, so the Claude half ran from a local
path; both source types land in the same versioned cache.

## Updating an install

```sh
claude plugin marketplace update fab && claude plugin update fab@fab   # then restart the session
codex plugin marketplace upgrade fab
```

## Adding or changing a skill

Every skill goes through a grill session (see `skill-audit.md`):

1. Walk the source and its process.
2. Grill it: cut what the model already knows, merge overlaps, or drop it.
3. Name review against `conventions.md`.
4. Write it in the standard shape.
5. Test it live in both hosts.
6. Log the decision and evidence in `skill-audit.md`.
7. Run the check, bump the version, commit.

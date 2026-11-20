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

`codex exec` prints only the final message; add `--json` to see intermediate messages and commands. It also waits on
stdin when stdin is open but empty (a background shell), so redirect from `/dev/null`.

Verified on Claude Code 2.1.329 and Codex CLI 0.170.0.

## Versioning

The version lives only in `.claude-plugin/plugin.json` and `.codex-plugin/plugin.json`; the check fails if they differ.
Bump on every shipped change: patch for wording, minor for a new skill or recipe, major for a rename or removal.

## Adding or changing a skill

Every skill goes through a grill session (see `skill-audit.md`):

1. Walk the source and its process.
2. Grill it: cut what the model already knows, merge overlaps, or drop it.
3. Name review against `conventions.md`.
4. Write it in the standard shape.
5. Test it live in both hosts.
6. Log the decision and evidence in `skill-audit.md`.
7. Run the check, bump the version, commit.

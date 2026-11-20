# fab: repo instructions

fab is a plugin of dev skills and recipes for Claude Code and Codex, built from one `skills/` directory.

## Read first

- `docs/conventions.md`: the rules every skill and recipe follows. `tools/check` enforces the ones marked
  **(checked)**; if the two disagree, fix the checker to match the doc.

## Setup (once per clone)

```sh
git config core.hooksPath .githooks   # pre-commit runs tools/check
```

Requires Go 1.24+ for development only. Installing fab as a plugin needs no Go.

## Commands

```sh
go run ./tools/check   # validate skills, recipes, manifests
go test ./...          # validator tests
```

## Rules

- Every commit passes `go run ./tools/check`. Never bypass the hook with `--no-verify`.
- Conventional Commits: `<type>(<scope>): <summary>`, types `feat`, `fix`, `docs`, `chore`, `test`.
- No `Co-Authored-By` or any AI-attribution trailer on commits, and no attribution line in PR descriptions.
- Bump the version in both `plugin.json` files on every shipped change (see conventions: Versioning).
- Never put personal data (voice profile, preferences) in this repo. It belongs in `~/.config/fab/`.
- Skill bodies stay host-neutral: no host-specific tool names (see conventions: Forbidden strings).

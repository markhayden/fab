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
agy plugin validate .  # Antigravity's own manifest check, without installing (one expected AGENTS.md warning)
```

The check exits 1 on errors and 2 if it can't run. A skill over 150 lines warns but passes.

## Local test loop

Claude, without installing:

```sh
claude --plugin-dir ~/go/src/github.com/markhayden/fab
```

Installed from a local checkout (every host). Plugin caches copy the files, so reinstall after each change:

```sh
claude plugin marketplace add ~/go/src/github.com/markhayden/fab
claude plugin uninstall fab@fab; claude plugin install fab@fab

codex plugin marketplace add ~/go/src/github.com/markhayden/fab
codex plugin remove fab@fab; codex plugin add fab@fab

agy plugin install ~/go/src/github.com/markhayden/fab   # copies to ~/.gemini/config/plugins/fab
agy plugin validate ~/go/src/github.com/markhayden/fab  # one expected warning: AGENTS.md isn't a plugin rule
```

Headless runs, useful for live tests:

```sh
claude -p '/fab:voice <ask>'
codex exec --skip-git-repo-check '$fab:voice <ask>' </dev/null
agy --settings '{"permissions":{"allow":["command(git)","command(ls)","command(mkdir)","write_file(.fab/)"]}}' \
  --add-dir ~/.config/fab -p '/grill <ask>' </dev/null
```

Headless `agy` can't prompt: a tool no rule allows is denied, which ends the run with no output (still exit 0), so an
empty reply with exit 0 means a tool was denied: check the allow rules. Pass `--settings` with the allow rules a test
needs, and `--add-dir ~/.config/fab`,

since every skill reads its override there, every flag before `-p`. In Antigravity, `/plan` and `/voice` reach its own
commands, so test those two by asking for them in words ("run fab's voice skill: <ask>").

With `--plugin-dir`, headless Claude can't read the checkout's `docs/` or `references/` (outside its allowed
directories), so skills fall back to their own text, e.g. inventing a ledger format. Add `--add-dir <checkout>` for
faithful runs. An installed plugin doesn't have this problem: Claude reads its own plugin cache.

MCP servers that a launcher passes to an interactive session (not in Claude's config) are missing in `claude -p`. For
the design skill, pass Paper explicitly: `--mcp-config '{"mcpServers":{"paper":{"type":"stdio","command":"~/.paper/bin/paper","args":["mcp"]}}}' --allowedTools mcp__paper`.

`codex exec` prints only the final message; add `--json` to see intermediate messages and commands. It also waits on
stdin when stdin is open but empty (a background shell), so redirect from `/dev/null`.

To live-test an unreleased skill in Codex, point the prompt at the checkout instead: "read and follow
`<checkout>/skills/<name>/SKILL.md`, resolving its relative links against that folder". Copying files into
`~/.codex/plugins/cache/` no longer works: Codex rebuilds the cache from the marketplace at startup (seen
2026-11-21).

Verified on Claude Code 2.1.329 and Codex CLI 0.170.0; the Antigravity lines on Antigravity CLI 1.5.17.

## Versioning

The version lives only in the three `plugin.json` files (`.claude-plugin/`, `.codex-plugin/`, `.antigravity-plugin/`);
the check fails if they differ.
Bump on every shipped change, as `conventions.md` Versioning says.

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
agy plugin install https://github.com/markhayden/fab                  # same command as installing
```

## Adding or changing a skill

Every skill goes through a grill session (see `skill-audit.md`):

1. Walk the source and its process.
2. Grill it: cut what the model already knows, merge overlaps, or drop it.
3. Name review against `conventions.md`.
4. Write it in the standard shape.
5. Test it live in every host.
6. Log the decision and evidence in `skill-audit.md`.
7. Run the check, bump the version, commit.

### Writing a skill body

On top of conventions' "Writing the body":

- **End every step on a completion criterion** the agent can check, and make it exhaustive ("every criterion maps
  to a test") rather than vague ("understanding reached"); a vague end invites stopping early.
- **One meaning, one place.** A rule two skills share lives in `conventions.md` and both link to it.
- **Leading words.** Where a phrase keeps spelling out one idea, name it once with a word the model already knows
  (a *tight* loop, a *seam*, a *picks gate*) and reuse the word.
- **Say what to do, not what to avoid.** A prohibition puts the unwanted behavior in front of the model; state the
  wanted one, and keep "never" for real hard lines.

Adapted from Matt Pocock's `writing-for-agents` skill ([mattpocock/skills](https://github.com/mattpocock/skills), MIT).

## Cross-host verification (1.0.0)

2026-11-21, installed plugin from GitHub, Claude Code 2.1.330 and Codex CLI 0.170.0, headless (`claude -p`,
`codex exec`), a fresh session per run, in throwaway clones of a small Go CLI.

| Check | Claude Code | Codex |
|---|---|---|
| `bigdev` called by name prints `[bigdev]` | pass | pass |
| `tinydev` called by name prints `[tinydev]` | pass | pass |
| A plain request matching bigdev's description (×2) doesn't start it | pass | pass |
| A plain request matching tinydev's description (×2) doesn't start it | pass | pass |
| Resume: spec approved → next is reuse/design (skipped, no UI), then plan; spec untouched | pass | pass |
| Resume: spec left as a draft → its gate is re-presented, nothing regenerated | pass | pass |
| Resume: spec edited after approval → re-confirmation, with a review round on the change | pass | pass |
| Resume: spec edited 5 ms after approval → caught by the sub-second mtime | pass | pass |
| An override that skips a gate bigdev keeps → reported and asked, not applied | pass | pass |
| tinydev end to end | feature route | bug route |

The first plain-request round failed in both hosts: the model followed grill's links to conventions, found the recipe
files in the plugin cache, and ran one. The host settings held; the model went around them. Since 0.18.1 a recipe runs
only when the user names it, and each recipe's body opens by saying so; the rerun passed, including a Codex run that
read both recipe files and stopped. Plain requests must not name fab: "run the full fab flow" names it, and the model
rightly follows.

# Skill audit log

One row per audited skill or recipe: where it came from, what was decided, and the live-test evidence.

| Skill | Source | Decision | Live test |
|---|---|---|---|
| `voice` | Mark's `write-like-me` command | Kept and generalized: the skill holds the process (load profile, draft, self-check, deliver, profile-building interview behind a gate); the voice itself moved to the private `~/.config/fab/voice.md`, built from `references/voice-template.md`. Name kept (approved in conventions). | 2026-11-20, Claude 2.1.329 + Codex 0.170.0: no profile → offers to build, writes nothing; profile rule applied; preferences beat default; override beats preferences and profile; unreadable override → one-line warning, falls back. Real profile verified in both. |

---
name: voice
kind: skill
description: Write or rewrite content in the user's own voice, using their private voice profile. Use for bios, posts, emails, About pages, or any "write this like me" / "make this sound like me" request.
---

# Voice

Writes content that sounds like the user, not like an AI imitating them. The voice lives in a private profile, so this
skill stays generic and the profile stays out of the plugin.

## Inputs

- The ask: a topic and format (bio, post, email, About section) to draft, or a pasted draft to rewrite.
- Voice profile: `~/.config/fab/voice.md`.
- Preferences: `~/.config/fab/preferences.md`, for any writing rules in it.
- Personal override: if `~/.config/fab/skills/voice.md` exists, read it and apply it on top of this skill (see fab conventions: Personal layer).

## Process

1. **Load the profile.**
   - Present: its hard rules bind every draft; its tone settings and examples guide everything else. Order when
     they conflict: the `voice` override, then the profile, then preferences.
   - Missing: say so in one line and offer two options: build a profile now (see Creating a profile), or write
     this one in plain, direct prose with no filler.
   - Exists but unreadable: warn in one line and write in plain, direct prose.
2. **Understand the ask.** Format, rough length, and audience. If something is unclear, ask one question, not three.
3. **Draft.** Aim under the requested length rather than over. A rewrite keeps the user's facts, names, and numbers,
   invents nothing, and cuts 20–40% by default.
4. **Check before showing.** Go through the profile's hard rules and its self-check list. Fix anything that fails
   and revise once. Show only the final version.
5. **Deliver.** The content only, with no preamble. If you made a judgment call (angle, tone, length), say so in one
   short line after the content.

If a rewrite's purpose clashes with the voice (a legal notice, a formal policy), say so in one sentence and ask
whether to apply the voice anyway.

### Creating a profile

Only when the user agrees to build one.

1. Ask for 2–3 samples of their real writing. Samples beat self-description.
2. Walk the sections of `references/voice-template.md` one question at a time, proposing an answer drawn from the
   samples for each.
3. Show the finished profile.

**GATE:** the full profile and its path `~/.config/fab/voice.md`, wait for approval before writing it.

Never write the profile into a project or the plugin. It is personal data.

## Output

The finished content, plus at most one line on any judgment call. When building a profile: the approved
`~/.config/fab/voice.md`.

## Done when

- The delivered text passes every hard rule in the profile (or, with no profile, has no filler or setup lines).
- A rewrite keeps every fact, name, and number from the original.
- For a new profile: the file exists at `~/.config/fab/voice.md` and was written only after approval.

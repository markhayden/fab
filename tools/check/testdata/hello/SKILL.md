---
name: hello
kind: skill
description: Install ping for fab. Use only when the user asks to run the fab install ping.
---

# Hello

Spike stub that proves fab installs and resolves `fab:hello`.

## Inputs

- Personal override: if `~/.config/fab/skills/hello.md` exists, read it and apply it on top of this skill (see fab conventions: Personal layer).

## Process

Reply with exactly this line and nothing else:

FAB-HELLO-7Q

## Output

The marker line.

## Done when

- The reply is exactly the marker.

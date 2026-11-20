---
name: hello-recipe
kind: recipe
description: Produces a spike greeting report. Use when the user asks for a spike greeting report.
disable-model-invocation: true
---

# Hello recipe

Spike stub that proves recipes are call-only in both hosts.

## Steps

Reply with exactly this line and nothing else:

FAB-RECIPE-3K

## Rules

- Personal override: if `~/.config/fab/skills/hello-recipe.md` exists, read it and apply it on top of this skill (see fab conventions: Personal layer).

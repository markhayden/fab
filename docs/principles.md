# Principles

Why fab looks the way it does. Every skill and recipe is judged against these ten rules.

1. **Write only what the model wouldn't do on its own.** A skill carries opinions, conventions, file locations, gates,
   and thresholds. It doesn't teach general engineering; the model already knows that, and every extra line costs context.
2. **Give the reason, not the volume.** One clause of "why" beats a rule in capital letters. The model follows reasons
   further than it follows shouting. Strict ordering is kept for gates only.
3. **Traps come from observation.** A "Traps" entry is added only after an agent has actually made that mistake.
   Imagined failure modes bloat skills and rarely match the real ones.
4. **Hand off through files.** Workflow skills write one defined artifact into the run directory and read the previous
   one. That lets recipes resume after an interruption and run steps in a fresh sub-agent.
5. **One gate marker.** Every approval point is written `**GATE:** <what>, wait for approval.` so recipes can keep or skip
   gates by name, and every decision lands in the run's ledger.
6. **Descriptions are always-loaded context.** At most 300 characters: what the skill does and the phrases that should
   trigger it. Nothing else.
7. **Fewer, broader skills.** Overlapping skills get merged by default. A smaller catalog is easier to remember and
   easier for the model to choose from.
8. **Neutral between hosts.** No host-specific tool names in skill bodies. One `SKILL.md` serves Claude Code and Codex.
9. **"Done when" means evidence.** Output, a path, a commit SHA. Never "I checked it".
10. **Tested live.** Every skill is run on a real task in both hosts before it ships. Automated skill tests come later.

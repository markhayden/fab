# Guided human loop

Read this when a bug only reproduces through steps a person has to take (a login, a device, a third-party
dashboard). With no run, use the system temp dir instead of `.fab/<run>/`.

1. Copy the script to `.fab/<run>/human-loop-<bug>.sh` and edit only the marked block: a `step` for each action, a
   `capture` for each observation, each prompt in single quotes (an apostrophe is written `'\''`). Capture what the user
   saw, never a secret: a sign-in is a `step`.
2. Ask the user to run `bash .fab/<run>/human-loop-<bug>.sh` in their own terminal. It prompts them and appends
   each answer to `.fab/<run>/human-loop-<bug>.out` under a dated header, so earlier runs are kept. Ask them not to
   paste secrets.
3. Read the latest block of `human-loop-<bug>.out`. Only a block ending in `DONE` whose answers clearly show the
   symptom (or, after a fix, its absence) is a verdict; rerun or ask otherwise.
   Copy the answers into the record redacted, and if a secret slipped in, delete it from the `.out` file.

```bash
#!/usr/bin/env bash
set -euo pipefail
umask 077
OUT="${0%.sh}.out"
printf '\n# run %s\n' "$(date -u +%FT%TZ)" >> "$OUT"

stopped() { printf 'STOPPED\n' >> "$OUT"; exit 1; }
trap stopped INT TERM

step() {                      # show an instruction, wait for Enter
  printf '\n>>> %s\n' "$1"
  IFS= read -r -p "    [Enter when done] " _ || stopped
}

capture() {                   # ask a question, save the answer as one line
  local var="$1" question="$2" answer=""
  printf '\n>>> %s\n' "$question"
  IFS= read -r -p "    > " answer || stopped
  printf '%s=%q\n' "$var" "$answer" >> "$OUT"
}

# --- edit below ---------------------------------------------------------
step 'Open http://localhost:3000 and sign in.'
capture ERRORED 'Click Export. Did it show an error? (y/n)'
capture ERROR_MSG 'Paste the error message, or none:'
# --- edit above ---------------------------------------------------------

printf 'DONE\n' >> "$OUT"
printf '\nSaved to %s\n' "$OUT"
```

The file is a record to read, not a script to source. Adapted from the `hitl-loop.template.sh` in Matt Pocock's
`diagnosing-bugs` skill ([mattpocock/skills](https://github.com/mattpocock/skills), MIT).

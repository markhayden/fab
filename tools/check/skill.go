package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode/utf8"

	"gopkg.in/yaml.v3"
)

// Finding is one problem the checker found. Warnings never fail the check.
type Finding struct {
	Path string
	Rule string
	Msg  string
	Warn bool
}

func (f Finding) String() string {
	level := "error"
	if f.Warn {
		level = "warn"
	}
	return fmt.Sprintf("%s: %s [%s] %s", f.Path, level, f.Rule, f.Msg)
}

const (
	maxDescription = 300
	warnLines      = 150
)

var (
	namePattern      = regexp.MustCompile(`^[a-z][a-z0-9-]{0,31}$`)
	requiredSections = map[string][]string{
		"skill":  {"Inputs", "Process", "Output", "Done when"},
		"recipe": {"Steps", "Rules"},
	}
)

// overrideLine is the verbatim personal-override line every skill and recipe
// must contain (docs/conventions.md: The override line).
func overrideLine(name string) string {
	return "- Personal override: if `~/.config/fab/skills/" + name +
		".md` exists, read it and apply it on top of this skill (see fab conventions: Personal layer)."
}

// frontmatter holds the SKILL.md fields fab cares about.
type frontmatter struct {
	Name                   string `yaml:"name"`
	Kind                   string `yaml:"kind"`
	Description            string `yaml:"description"`
	DisableModelInvocation *bool  `yaml:"disable-model-invocation"`
}

// skillFile is a parsed SKILL.md.
type skillFile struct {
	path  string
	front frontmatter
	body  string
	lines int
}

func readSkill(dir string) (*skillFile, []Finding) {
	path := filepath.Join(dir, "SKILL.md")
	fail := func(msg string) (*skillFile, []Finding) {
		return nil, []Finding{{Path: path, Rule: "frontmatter", Msg: msg}}
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		return fail(err.Error())
	}
	text := string(raw)
	if !strings.HasPrefix(text, "---\n") {
		return fail("must start with a --- frontmatter block")
	}
	end := strings.Index(text[4:], "\n---\n")
	if end < 0 {
		return fail("frontmatter block is not closed with ---")
	}
	var fm frontmatter
	if err := yaml.Unmarshal([]byte(text[4:4+end]), &fm); err != nil {
		return fail("invalid YAML: " + err.Error())
	}
	return &skillFile{
		path:  path,
		front: fm,
		body:  text[4+end+5:],
		lines: strings.Count(text, "\n"),
	}, nil
}

// checkSkill applies the per-skill rules to skills/<dir>/SKILL.md.
func checkSkill(dir string) []Finding {
	s, findings := readSkill(dir)
	if s == nil {
		return findings
	}
	add := func(rule, format string, args ...any) {
		findings = append(findings, Finding{Path: s.path, Rule: rule, Msg: fmt.Sprintf(format, args...)})
	}
	fm := s.front

	switch base := filepath.Base(dir); {
	case fm.Name != base:
		add("name", "name %q must match its directory %q", fm.Name, base)
	case !namePattern.MatchString(fm.Name):
		add("name", "name %q must be lowercase letters, digits, hyphens, start with a letter, ≤32 chars", fm.Name)
	}

	sections, ok := requiredSections[fm.Kind]
	if !ok {
		add("kind", "kind must be skill or recipe, got %q", fm.Kind)
	}

	switch n := utf8.RuneCountInString(strings.TrimSpace(fm.Description)); {
	case n == 0:
		add("description", "description is required")
	case n > maxDescription:
		add("description", "description is %d chars, max %d", n, maxDescription)
	}

	have := headings(s.body)
	for _, want := range sections {
		if !have[want] {
			add("sections", "missing required section ## %s for kind %s", want, fm.Kind)
		}
	}

	want := overrideSection[fm.Kind]
	switch got, found := sectionOf(s.body, overrideLine(fm.Name)); {
	case !found:
		add("override", "missing the standard personal-override line for %q", fm.Name)
	case want != "" && got != want:
		add("override", "the personal-override line belongs under ## %s, not ## %s", want, got)
	}

	if s.lines > warnLines {
		findings = append(findings, Finding{Path: s.path, Rule: "length", Warn: true,
			Msg: fmt.Sprintf("%d lines; aim for ~100 and move detail to references/", s.lines)})
	}
	return findings
}

// headings returns the "## " section titles outside fenced code blocks.
// overrideSection is where the override line goes for each kind
// (docs/conventions.md: The override line).
var overrideSection = map[string]string{"skill": "Inputs", "recipe": "Rules"}

// sectionOf returns the ## heading a line sits under, outside code fences,
// and whether the line was found.
func sectionOf(body, line string) (string, bool) {
	section, inFence := "", false
	sc := bufio.NewScanner(bytes.NewBufferString(body))
	for sc.Scan() {
		text := sc.Text()
		switch {
		case strings.HasPrefix(text, "```"):
			inFence = !inFence
		case !inFence && strings.HasPrefix(text, "## "):
			section = strings.TrimSpace(strings.TrimPrefix(text, "## "))
		case !inFence && text == line:
			return section, true
		}
	}
	return "", false
}

func headings(body string) map[string]bool {
	out := map[string]bool{}
	inFence := false
	sc := bufio.NewScanner(bytes.NewBufferString(body))
	for sc.Scan() {
		line := sc.Text()
		if strings.HasPrefix(line, "```") {
			inFence = !inFence
			continue
		}
		if !inFence && strings.HasPrefix(line, "## ") {
			out[strings.TrimSpace(strings.TrimPrefix(line, "## "))] = true
		}
	}
	return out
}

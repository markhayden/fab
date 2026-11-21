package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// fixture describes a SKILL.md to generate. Zero values produce a valid skill.
type fixture struct {
	dir         string // directory name; defaults to name
	name        string
	kind        string
	description string
	sections    []string // "## " headings; defaults by kind
	noOverride  bool
	overrideIn  string  // section holding the override line; default Inputs (skill) or Rules (recipe)
	rawFront    *string // replaces the whole frontmatter block when set
	extraLines  int     // filler lines appended to the body
}

// overrideSection is where the fixture puts the override line.
func (f fixture) overrideSection() string {
	if f.overrideIn != "" {
		return f.overrideIn
	}
	if f.kind == "recipe" {
		return "Rules"
	}
	return "Inputs"
}

func (f fixture) write(t *testing.T) string {
	t.Helper()
	if f.name == "" {
		f.name = "demo"
	}
	if f.dir == "" {
		f.dir = f.name
	}
	if f.kind == "" {
		f.kind = "skill"
	}
	if f.description == "" {
		f.description = "Demo skill. Use when testing fab's checker."
	}
	if f.sections == nil {
		f.sections = map[string][]string{
			"skill":  {"Inputs", "Process", "Output", "Done when"},
			"recipe": {"Steps", "Rules"},
		}[f.kind]
	}

	var b strings.Builder
	if f.rawFront != nil {
		b.WriteString(*f.rawFront)
	} else {
		b.WriteString("---\nname: " + f.name + "\nkind: " + f.kind + "\ndescription: " + f.description + "\n---\n")
	}
	b.WriteString("\n# Demo\n\nWhat it produces.\n")
	for _, s := range f.sections {
		b.WriteString("\n## " + s + "\n\n")
		if s == f.overrideSection() && !f.noOverride {
			b.WriteString(overrideLine(f.name) + "\n")
		} else {
			b.WriteString("- item\n")
		}
	}
	for i := 0; i < f.extraLines; i++ {
		b.WriteString("filler\n")
	}

	dir := filepath.Join(t.TempDir(), f.dir)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte(b.String()), 0o644); err != nil {
		t.Fatal(err)
	}
	return dir
}

func strp(s string) *string { return &s }

func TestCheckSkill(t *testing.T) {
	cases := []struct {
		name     string
		fx       fixture
		wantRule string // "" means no errors expected
		wantWarn string // rule expected as a warning (not an error)
	}{
		{name: "valid skill", fx: fixture{}},
		{name: "valid recipe", fx: fixture{kind: "recipe"}},
		{name: "no frontmatter", fx: fixture{rawFront: strp("")}, wantRule: "frontmatter"},
		{name: "unterminated frontmatter", fx: fixture{rawFront: strp("---\nname: demo\n")}, wantRule: "frontmatter"},
		{name: "invalid yaml", fx: fixture{rawFront: strp("---\nname: [demo\n---\n")}, wantRule: "frontmatter"},
		{name: "name differs from dir", fx: fixture{dir: "other"}, wantRule: "name"},
		{name: "name bad characters", fx: fixture{name: "Bad_Name"}, wantRule: "name"},
		{name: "name too long", fx: fixture{name: "a" + strings.Repeat("b", 32)}, wantRule: "name"},
		{name: "name exactly 32", fx: fixture{name: "a" + strings.Repeat("b", 31)}},
		{name: "missing kind", fx: fixture{rawFront: strp("---\nname: demo\ndescription: d\n---\n")}, wantRule: "kind"},
		{name: "unknown kind", fx: fixture{kind: "workflow", sections: []string{"Inputs"}}, wantRule: "kind"},
		{name: "empty description", fx: fixture{rawFront: strp("---\nname: demo\nkind: skill\ndescription: \"\"\n---\n")}, wantRule: "description"},
		{name: "description 300", fx: fixture{description: strings.Repeat("x", 300)}},
		{name: "description 301", fx: fixture{description: strings.Repeat("x", 301)}, wantRule: "description"},
		{name: "skill missing Done when", fx: fixture{sections: []string{"Inputs", "Process", "Output"}}, wantRule: "sections"},
		{name: "recipe missing Rules", fx: fixture{kind: "recipe", sections: []string{"Steps"}}, wantRule: "sections"},
		{name: "missing override line", fx: fixture{noOverride: true}, wantRule: "override"},
		{name: "override line outside Inputs", fx: fixture{overrideIn: "Process"}, wantRule: "override"},
		{name: "recipe override outside Rules", fx: fixture{kind: "recipe", overrideIn: "Steps"}, wantRule: "override"},
		{name: "recipe override under Rules", fx: fixture{kind: "recipe"}},
		{name: "151 lines warns only", fx: fixture{extraLines: 140}, wantWarn: "length"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			findings := checkSkill(tc.fx.write(t))
			var errs, warns []Finding
			for _, f := range findings {
				if f.Warn {
					warns = append(warns, f)
				} else {
					errs = append(errs, f)
				}
			}

			if tc.wantRule == "" && len(errs) > 0 {
				t.Fatalf("expected no errors, got %v", errs)
			}
			if tc.wantRule != "" && !hasRule(errs, tc.wantRule) {
				t.Fatalf("expected error %q, got errors %v", tc.wantRule, errs)
			}
			if tc.wantWarn != "" && !hasRule(warns, tc.wantWarn) {
				t.Fatalf("expected warning %q, got warnings %v", tc.wantWarn, warns)
			}
			if tc.wantWarn == "" && len(warns) > 0 {
				t.Fatalf("expected no warnings, got %v", warns)
			}
		})
	}
}

func TestLengthBoundary(t *testing.T) {
	dir := fixture{}.write(t)
	n := countLines(t, dir)
	// Pad to exactly 150 lines: no warning. 151: warning.
	at150 := fixture{extraLines: 150 - n}.write(t)
	if got := countLines(t, at150); got != 150 {
		t.Fatalf("fixture has %d lines, want 150", got)
	}
	for _, f := range checkSkill(at150) {
		if f.Rule == "length" {
			t.Fatalf("150 lines should not warn: %v", f)
		}
	}
	at151 := fixture{extraLines: 151 - n}.write(t)
	if !hasRule(checkSkill(at151), "length") {
		t.Fatal("151 lines should warn")
	}
}

func hasRule(fs []Finding, rule string) bool {
	for _, f := range fs {
		if f.Rule == rule {
			return true
		}
	}
	return false
}

func countLines(t *testing.T, dir string) int {
	t.Helper()
	b, err := os.ReadFile(filepath.Join(dir, "SKILL.md"))
	if err != nil {
		t.Fatal(err)
	}
	return strings.Count(string(b), "\n")
}

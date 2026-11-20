package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// skillMD builds a valid SKILL.md for name/kind with extra body text.
func skillMD(name, kind, extraFront, body string) string {
	sections := map[string][]string{
		"skill":  {"Inputs", "Process", "Output", "Done when"},
		"recipe": {"Rules", "Steps"},
	}[kind]
	var b strings.Builder
	b.WriteString("---\nname: " + name + "\nkind: " + kind + "\ndescription: Demo. Use when testing.\n" + extraFront + "---\n\n# Demo\n")
	for i, s := range sections {
		b.WriteString("\n## " + s + "\n\n")
		if i == 0 {
			b.WriteString(overrideLine(name) + "\n")
		}
	}
	b.WriteString("\n" + body + "\n")
	return b.String()
}

// validRepo is a minimal repo that passes every rule: a skill with a
// reference file and a link, and a recipe that references it.
func validRepo() map[string]string {
	return map[string]string{
		".claude-plugin/plugin.json":       `{"name":"fab","version":"0.1.0"}`,
		".codex-plugin/plugin.json":        `{"name":"fab","version":"0.1.0"}`,
		".claude-plugin/marketplace.json":  `{"name":"fab","plugins":[{"name":"fab","source":"./"}]}`,
		".agents/plugins/marketplace.json": `{"name":"fab","plugins":[{"name":"fab","source":{"source":"local","path":"./"}}]}`,
		"skills/plan/SKILL.md": skillMD("plan", "skill", "",
			"Read `references/format.md` when writing. See [format](references/format.md#top) and [site](https://example.com)."),
		"skills/plan/references/format.md": "# Format\n",
		"skills/go/SKILL.md":               skillMD("go", "recipe", "disable-model-invocation: true\n", "1. Run fab:plan."),
		"skills/go/agents/openai.yaml":     "policy:\n  allow_implicit_invocation: false\n",
	}
}

func writeRepo(t *testing.T, files map[string]string) string {
	t.Helper()
	root := t.TempDir()
	for rel, content := range files {
		p := filepath.Join(root, rel)
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

func TestRepoValid(t *testing.T) {
	findings, err := run(writeRepo(t, validRepo()))
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range findings {
		t.Errorf("unexpected finding: %s", f)
	}
}

func TestCheckRepo(t *testing.T) {
	const del = "\x00delete"
	cases := []struct {
		name     string
		change   map[string]string // path → new content, or del
		wantRule string
		wantFile string // suffix of the reported path
	}{
		{"recipe missing claude flag", map[string]string{
			"skills/go/SKILL.md": skillMD("go", "recipe", "", "1. Run fab:plan."),
		}, "recipe-claude", "go/SKILL.md"},
		{"recipe claude flag false", map[string]string{
			"skills/go/SKILL.md": skillMD("go", "recipe", "disable-model-invocation: false\n", "1. Run fab:plan."),
		}, "recipe-claude", "go/SKILL.md"},
		{"recipe missing openai.yaml", map[string]string{
			"skills/go/agents/openai.yaml": del,
		}, "recipe-codex", "go/agents/openai.yaml"},
		{"recipe codex policy absent", map[string]string{
			"skills/go/agents/openai.yaml": "interface:\n  display_name: go\n",
		}, "recipe-codex", "go/agents/openai.yaml"},
		{"recipe codex policy true", map[string]string{
			"skills/go/agents/openai.yaml": "policy:\n  allow_implicit_invocation: true\n",
		}, "recipe-codex", "go/agents/openai.yaml"},
		{"recipe codex invalid yaml", map[string]string{
			"skills/go/agents/openai.yaml": "policy: [\n",
		}, "recipe-codex", "go/agents/openai.yaml"},
		{"unknown skill reference", map[string]string{
			"skills/go/SKILL.md": skillMD("go", "recipe", "disable-model-invocation: true\n", "1. Run fab:nope."),
		}, "skill-ref", "go/SKILL.md"},
		{"broken relative link", map[string]string{
			"skills/plan/references/format.md": "See [missing](other.md).\n",
		}, "link", "plan/references/format.md"},
		{"missing references file", map[string]string{
			"skills/plan/references/format.md": del,
		}, "references", "plan/SKILL.md"},
		{"forbidden string, case-insensitive", map[string]string{
			"skills/plan/references/format.md": "Call askuserquestion here.\n",
		}, "forbidden", "plan/references/format.md"},
		{"forbidden string in non-markdown file", map[string]string{
			"skills/go/agents/openai.yaml": "# waffleoff\npolicy:\n  allow_implicit_invocation: false\n",
		}, "forbidden", "go/agents/openai.yaml"},
		{"plugin versions differ", map[string]string{
			".codex-plugin/plugin.json": `{"name":"fab","version":"0.2.0"}`,
		}, "manifest", ".codex-plugin/plugin.json"},
		{"plugin name wrong", map[string]string{
			".claude-plugin/plugin.json": `{"name":"ramrod","version":"0.1.0"}`,
		}, "manifest", ".claude-plugin/plugin.json"},
		{"plugin version empty", map[string]string{
			".claude-plugin/plugin.json": `{"name":"fab"}`,
			".codex-plugin/plugin.json":  `{"name":"fab"}`,
		}, "manifest", "plugin.json"},
		{"plugin.json missing", map[string]string{
			".codex-plugin/plugin.json": del,
		}, "manifest", ".codex-plugin/plugin.json"},
		{"plugin.json invalid", map[string]string{
			".claude-plugin/plugin.json": `{`,
		}, "manifest", ".claude-plugin/plugin.json"},
		{"claude catalog wrong source", map[string]string{
			".claude-plugin/marketplace.json": `{"plugins":[{"name":"fab","source":"./sub"}]}`,
		}, "catalog", ".claude-plugin/marketplace.json"},
		{"claude catalog missing", map[string]string{
			".claude-plugin/marketplace.json": del,
		}, "catalog", ".claude-plugin/marketplace.json"},
		{"codex catalog wrong name", map[string]string{
			".agents/plugins/marketplace.json": `{"plugins":[{"name":"other","source":{"source":"local","path":"./"}}]}`,
		}, "catalog", ".agents/plugins/marketplace.json"},
		{"codex catalog not local", map[string]string{
			".agents/plugins/marketplace.json": `{"plugins":[{"name":"fab","source":{"source":"git","path":"./"}}]}`,
		}, "catalog", ".agents/plugins/marketplace.json"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			files := validRepo()
			for p, c := range tc.change {
				if c == del {
					delete(files, p)
				} else {
					files[p] = c
				}
			}
			findings, err := run(writeRepo(t, files))
			if err != nil {
				t.Fatal(err)
			}
			for _, f := range findings {
				if f.Rule == tc.wantRule && strings.HasSuffix(f.Path, tc.wantFile) && !f.Warn {
					return
				}
			}
			t.Errorf("want error rule %q on *%s, got %v", tc.wantRule, tc.wantFile, findings)
		})
	}
}

// The T1 spike stubs live on as known-good fixtures for a skill and a recipe.
func TestSpikeStubsPass(t *testing.T) {
	for _, dir := range []string{"testdata/hello", "testdata/hello-recipe"} {
		for _, f := range append(checkSkill(dir), checkRecipe(dir)...) {
			t.Errorf("%s: unexpected finding: %s", dir, f)
		}
	}
}

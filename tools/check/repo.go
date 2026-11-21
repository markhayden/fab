package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// Repo-wide rules from docs/conventions.md: recipe call-only declarations,
// skill references, links, manifests, catalogs, and forbidden strings.

var (
	// skillRef matches a host-neutral skill reference such as fab:plan.
	skillRef = regexp.MustCompile(`\bfab:([a-z][a-z0-9-]*)`)
	// mdLink matches the target of an inline markdown link.
	mdLink = regexp.MustCompile(`\]\(([^)\s]+)\)`)
	// urlPattern matches a URL, whose path is never a local file.
	urlPattern = regexp.MustCompile(`[A-Za-z][A-Za-z0-9+.-]*://[^\s)>\]]+`)
	// refPath matches a mention of a references/ file.
	refPath = regexp.MustCompile(`references/[A-Za-z0-9._/-]+\.md`)

	forbidden = []string{
		"agent-skills:", "waffleoff",
		"AskUserQuestion", "TodoWrite", "EnterPlanMode", "ExitPlanMode",
		"request_user_input", "update_plan",
	}
)

const pluginName = "fab"

// checkRepo applies the repo-wide rules. dirs are the skill directories.
func checkRepo(root string, dirs []string) []Finding {
	names := map[string]bool{}
	for _, d := range dirs {
		names[filepath.Base(d)] = true
	}
	var out []Finding
	for _, d := range dirs {
		out = append(out, checkRecipe(d)...)
		out = append(out, checkReferences(d, names)...)
	}
	out = append(out, checkForbidden(filepath.Join(root, "skills"))...)
	out = append(out, checkManifests(root)...)
	return out
}

// checkRecipe requires both call-only declarations on kind: recipe.
func checkRecipe(dir string) []Finding {
	s, _ := readSkill(dir)
	if s == nil || s.front.Kind != "recipe" {
		return nil
	}
	var out []Finding
	if v := s.front.DisableModelInvocation; v == nil || !*v {
		out = append(out, Finding{Path: s.path, Rule: "recipe-claude",
			Msg: "recipe must set disable-model-invocation: true"})
	}

	ypath := filepath.Join(dir, "agents", "openai.yaml")
	raw, err := os.ReadFile(ypath)
	if err != nil {
		return append(out, Finding{Path: ypath, Rule: "recipe-codex",
			Msg: "recipe needs agents/openai.yaml with policy.allow_implicit_invocation: false"})
	}
	var cfg struct {
		Policy struct {
			AllowImplicitInvocation *bool `yaml:"allow_implicit_invocation"`
		} `yaml:"policy"`
	}
	if err := yaml.Unmarshal(raw, &cfg); err != nil {
		return append(out, Finding{Path: ypath, Rule: "recipe-codex", Msg: "invalid YAML: " + err.Error()})
	}
	if v := cfg.Policy.AllowImplicitInvocation; v == nil || *v {
		out = append(out, Finding{Path: ypath, Rule: "recipe-codex",
			Msg: "policy.allow_implicit_invocation must be false"})
	}
	return out
}

// checkReferences checks every markdown file in a skill directory:
// fab:<name> references, relative links, and references/ paths must resolve.
func checkReferences(dir string, names map[string]bool) []Finding {
	var out []Finding
	_ = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || filepath.Ext(path) != ".md" {
			return err
		}
		raw, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		text := string(raw)

		for _, m := range uniq(skillRef.FindAllStringSubmatch(text, -1), 1) {
			if !names[m] {
				out = append(out, Finding{Path: path, Rule: "skill-ref",
					Msg: fmt.Sprintf("fab:%s does not match any skill in skills/", m)})
			}
		}
		for _, target := range uniq(mdLink.FindAllStringSubmatch(text, -1), 1) {
			if isExternal(target) {
				continue
			}
			target, _, _ = strings.Cut(target, "#")
			if !exists(filepath.Join(filepath.Dir(path), target)) {
				out = append(out, Finding{Path: path, Rule: "link",
					Msg: fmt.Sprintf("link target %q does not exist", target)})
			}
		}
		for _, ref := range uniq(refPathMatches(text), 0) {
			// A skill's own references/ path resolves against its directory;
			// a longer path (../review/references/x.md) against this file's.
			base := dir
			if !strings.HasPrefix(ref, "references/") {
				base = filepath.Dir(path)
			}
			if !exists(filepath.Join(base, ref)) {
				out = append(out, Finding{Path: path, Rule: "references",
					Msg: fmt.Sprintf("%s does not exist in %s", ref, filepath.Base(dir))})
			}
		}
		return nil
	})
	return out
}

// refPathMatches returns every path that names a references/ file: a skill's
// own ("references/x.md", "./references/x.md" trimmed to that) or one inside a
// longer path ("../review/references/x.md", kept whole). URLs are skipped.
func refPathMatches(text string) [][]string {
	var out [][]string
	text = urlPattern.ReplaceAllString(text, "")
	for _, loc := range refPath.FindAllStringIndex(text, -1) {
		start := loc[0]
		for start > 0 && isPathChar(text[start-1]) {
			start--
		}
		ref := strings.TrimPrefix(text[start:loc[1]], "./")
		out = append(out, []string{ref})
	}
	return out
}

func isPathChar(c byte) bool {
	return c == '.' || c == '/' || c == '_' || c == '-' ||
		'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || '0' <= c && c <= '9'
}

// checkForbidden fails on host-specific or legacy strings anywhere under skills/.
func checkForbidden(skillsDir string) []Finding {
	var out []Finding
	_ = filepath.WalkDir(skillsDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		raw, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		lower := strings.ToLower(string(raw))
		for _, s := range forbidden {
			if strings.Contains(lower, strings.ToLower(s)) {
				out = append(out, Finding{Path: path, Rule: "forbidden",
					Msg: fmt.Sprintf("contains %q (see conventions: Forbidden strings)", s)})
			}
		}
		return nil
	})
	return out
}

// checkManifests requires both plugin.json files to agree and both catalogs
// to list plugin fab at the repo root.
func checkManifests(root string) []Finding {
	var out []Finding
	type plugin struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	}
	paths := []string{
		filepath.Join(root, ".claude-plugin", "plugin.json"),
		filepath.Join(root, ".codex-plugin", "plugin.json"),
	}
	var versions []string
	for _, p := range paths {
		var m plugin
		if f := readJSON(p, &m); f != nil {
			out = append(out, *f)
			continue
		}
		if m.Name != pluginName {
			out = append(out, Finding{Path: p, Rule: "manifest", Msg: fmt.Sprintf("name must be %q, got %q", pluginName, m.Name)})
		}
		if m.Version == "" {
			out = append(out, Finding{Path: p, Rule: "manifest", Msg: "version is required"})
		}
		versions = append(versions, m.Version)
	}
	if len(versions) == 2 && versions[0] != versions[1] {
		out = append(out, Finding{Path: paths[1], Rule: "manifest",
			Msg: fmt.Sprintf("version %q differs from .claude-plugin/plugin.json %q", versions[1], versions[0])})
	}

	claudeCat := filepath.Join(root, ".claude-plugin", "marketplace.json")
	var cc struct {
		Plugins []struct {
			Name   string `json:"name"`
			Source string `json:"source"`
		} `json:"plugins"`
	}
	if f := readJSON(claudeCat, &cc); f != nil {
		out = append(out, *f)
	} else if !catalogHas(len(cc.Plugins), func(i int) bool { return cc.Plugins[i].Name == pluginName && isRoot(cc.Plugins[i].Source) }) {
		out = append(out, Finding{Path: claudeCat, Rule: "catalog", Msg: `must list plugin "fab" with source "./"`})
	}

	codexCat := filepath.Join(root, ".agents", "plugins", "marketplace.json")
	var xc struct {
		Plugins []struct {
			Name   string `json:"name"`
			Source struct {
				Source string `json:"source"`
				Path   string `json:"path"`
			} `json:"source"`
		} `json:"plugins"`
	}
	if f := readJSON(codexCat, &xc); f != nil {
		out = append(out, *f)
	} else if !catalogHas(len(xc.Plugins), func(i int) bool {
		p := xc.Plugins[i]
		return p.Name == pluginName && p.Source.Source == "local" && isRoot(p.Source.Path)
	}) {
		out = append(out, Finding{Path: codexCat, Rule: "catalog", Msg: `must list plugin "fab" with source local, path "./"`})
	}
	return out
}

func readJSON(path string, v any) *Finding {
	raw, err := os.ReadFile(path)
	if err != nil {
		rule := "manifest"
		if strings.HasSuffix(path, "marketplace.json") {
			rule = "catalog"
		}
		msg := err.Error()
		if errors.Is(err, fs.ErrNotExist) {
			msg = "file is missing"
		}
		return &Finding{Path: path, Rule: rule, Msg: msg}
	}
	if err := json.Unmarshal(raw, v); err != nil {
		return &Finding{Path: path, Rule: "manifest", Msg: "invalid JSON: " + err.Error()}
	}
	return nil
}

func catalogHas(n int, match func(int) bool) bool {
	for i := 0; i < n; i++ {
		if match(i) {
			return true
		}
	}
	return false
}

func isRoot(p string) bool { return p == "./" || p == "." }

func isExternal(target string) bool {
	return strings.HasPrefix(target, "#") || strings.Contains(target, "://") || strings.HasPrefix(target, "mailto:")
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// uniq returns the sorted distinct values of group g across matches.
func uniq(matches [][]string, g int) []string {
	seen := map[string]bool{}
	for _, m := range matches {
		seen[m[g]] = true
	}
	out := make([]string, 0, len(seen))
	for k := range seen {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// Command check validates the fab plugin: skills, recipes, and manifests.
// Rules are defined in docs/conventions.md. Run from the repo root:
//
//	go run ./tools/check
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

func main() {
	root := flag.String("root", ".", "repo root")
	flag.Parse()

	findings, err := run(*root)
	if err != nil {
		fmt.Fprintln(os.Stderr, "check:", err)
		os.Exit(2)
	}
	errors := 0
	for _, f := range findings {
		fmt.Println(f)
		if !f.Warn {
			errors++
		}
	}
	if errors > 0 {
		fmt.Printf("check: %d error(s)\n", errors)
		os.Exit(1)
	}
	fmt.Println("check: ok")
}

// run checks every skill directory under root/skills.
func run(root string) ([]Finding, error) {
	dirs, err := skillDirs(root)
	if err != nil {
		return nil, err
	}
	var findings []Finding
	for _, d := range dirs {
		findings = append(findings, checkSkill(d)...)
	}
	return findings, nil
}

func skillDirs(root string) ([]string, error) {
	entries, err := os.ReadDir(filepath.Join(root, "skills"))
	if err != nil {
		return nil, err
	}
	var dirs []string
	for _, e := range entries {
		if e.IsDir() {
			dirs = append(dirs, filepath.Join(root, "skills", e.Name()))
		}
	}
	sort.Strings(dirs)
	return dirs, nil
}

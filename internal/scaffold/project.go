/*
|--------------------------------------------------------------------------
| Project detection
|--------------------------------------------------------------------------
|
| Finds the app module path (go.mod) and modules/ directory for generators.
|
*/

package scaffold

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var modulePathRe = regexp.MustCompile(`^module\s+(\S+)`)

// Project holds paths for a gogql application tree.
type Project struct {
	Root       string
	ModulePath string
	ModulesDir string
}

// FindProject locates go.mod starting at dir (or cwd) and walking up.
func FindProject(dir string) (*Project, error) {
	if dir == "" {
		var err error
		dir, err = os.Getwd()
		if err != nil {
			return nil, err
		}
	}
	abs, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}
	for {
		modPath := filepath.Join(abs, "go.mod")
		if _, err := os.Stat(modPath); err == nil {
			mp, err := readModulePath(modPath)
			if err != nil {
				return nil, err
			}
			modulesDir := filepath.Join(abs, "modules")
			return &Project{
				Root:       abs,
				ModulePath: mp,
				ModulesDir: modulesDir,
			}, nil
		}
		parent := filepath.Dir(abs)
		if parent == abs {
			return nil, fmt.Errorf("no go.mod found (run from your app root or use gogql init)")
		}
		abs = parent
	}
}

func readModulePath(goMod string) (string, error) {
	f, err := os.Open(goMod)
	if err != nil {
		return "", err
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if m := modulePathRe.FindStringSubmatch(line); len(m) == 2 {
			return m[1], nil
		}
	}
	return "", fmt.Errorf("module path not found in %s", goMod)
}

// ModuleNames lists subdirectories of modules/ (excluding files).
func (p *Project) ModuleNames() ([]string, error) {
	entries, err := os.ReadDir(p.ModulesDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("modules/ not found under %s (run gogql init first)", p.Root)
		}
		return nil, err
	}
	var names []string
	for _, e := range entries {
		if e.IsDir() {
			names = append(names, e.Name())
		}
	}
	return names, nil
}

func (p *Project) moduleDir(name string) string {
	return filepath.Join(p.ModulesDir, SanitizePackageName(name))
}

func (p *Project) hasLoaderFactories(pkg string) bool {
	data, err := os.ReadFile(filepath.Join(p.ModulesDir, pkg, "resolvers.go"))
	if err != nil {
		return false
	}
	return strings.Contains(string(data), "func LoaderFactories()")
}

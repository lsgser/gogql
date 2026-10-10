/*
|--------------------------------------------------------------------------
| Project detection
|--------------------------------------------------------------------------
|
| Locates the user's gogql application on disk for CLI generators. Walks
| upward from cwd to find go.mod, reads the module import path, then chooses
| LayoutSrc (src/modules + src/schema/modules.go) or LayoutLegacy (modules/
| + modules/registry.go). Exposes ModuleImportPrefix for generated imports.
|
| ModuleNames / moduleNamesSorted list domain folders; common/ is ordered first
| in src layout. hasLoaderFactories scans *resolvers.go for LoaderFactories.
|
| Key types: Project, Layout. Key func: FindProject.
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

// Layout describes the scaffolded project tree.
type Layout int

const (
	LayoutSrc    Layout = iota // src/modules, src/schema/modules.go
	LayoutLegacy               // modules/, modules/registry.go
)

// Project holds paths for a gogql application tree.
type Project struct {
	Root               string
	ModulePath         string
	Layout             Layout
	ModulesDir         string
	RegistryPath       string
	ModuleImportPrefix string
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
			p := &Project{Root: abs, ModulePath: mp}
			p.resolveLayout()
			return p, nil
		}
		parent := filepath.Dir(abs)
		if parent == abs {
			return nil, fmt.Errorf("no go.mod found (run from your app root or use gogql init)")
		}
		abs = parent
	}
}

func (p *Project) resolveLayout() {
	srcModules := filepath.Join(p.Root, "src", "modules")
	if st, err := os.Stat(srcModules); err == nil && st.IsDir() {
		p.Layout = LayoutSrc
		p.ModulesDir = srcModules
		p.RegistryPath = filepath.Join(p.Root, "src", "schema", "modules.go")
		p.ModuleImportPrefix = p.ModulePath + "/src/modules"
		return
	}
	p.Layout = LayoutLegacy
	p.ModulesDir = filepath.Join(p.Root, "modules")
	p.RegistryPath = filepath.Join(p.Root, "modules", "registry.go")
	p.ModuleImportPrefix = p.ModulePath + "/modules"
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

func hasDir(path string) bool {
	st, err := os.Stat(path)
	return err == nil && st.IsDir()
}

func (p *Project) moduleDir(name string) string {
	return filepath.Join(p.ModulesDir, SanitizePackageName(name))
}

func (p *Project) hasLoaderFactories(pkg string) bool {
	dir := filepath.Join(p.ModulesDir, pkg)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false
	}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".go") {
			continue
		}
		if !strings.Contains(e.Name(), "resolvers") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			continue
		}
		if strings.Contains(string(data), "func LoaderFactories()") {
			return true
		}
	}
	return false
}

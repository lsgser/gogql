/*
|--------------------------------------------------------------------------
| Module generators
|--------------------------------------------------------------------------
|
| CLI backends for gogql module add | typedefs | resolvers | schema. Creates
| domain folders under src/modules/<name> with co-located *.graphql,
| *.resolvers.go, *.model.go, *.service.go, and module.go, then rewrites
| the schema registry (src/schema/modules.go or legacy modules/registry.go).
|
| AddModule is full scaffold; partial commands fill in missing layers. SyncRegistry
| discovers subdirectories and emits imports plus All() / LoaderFactories().
|
| Key type: ModuleGenOptions. Key funcs: AddModule, AddTypedefs, AddResolvers,
| AddSchema, SyncRegistry.
|
*/

package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// ModuleGenOptions configures module file generation.
type ModuleGenOptions struct {
	ProjectDir string
	Name       string
	Force      bool
}

// AddModule writes domain files: module.go, *.graphql, *.resolvers.go, *.model.go, *.service.go.
func AddModule(opts ModuleGenOptions) error {
	if opts.Name == "" {
		return fmt.Errorf("module name is required")
	}
	proj, err := projectFor(opts.ProjectDir)
	if err != nil {
		return err
	}
	if SanitizePackageName(opts.Name) == "common" {
		return fmt.Errorf("use src/common for shared SDL; pick another module name")
	}
	modDir := proj.moduleDir(opts.Name)
	if err := os.MkdirAll(modDir, 0o755); err != nil {
		return err
	}
	data := moduleData(opts.Name)
	steps := []struct {
		tmpl string
		dest string
	}{
		{"templates/module/module.go.tmpl", "module.go"},
		{"templates/module/domain.graphql.tmpl", data.PackageName + ".graphql"},
		{"templates/module/domain.resolvers.go.tmpl", data.PackageName + ".resolvers.go"},
		{"templates/module/domain.model.go.tmpl", data.PackageName + ".model.go"},
		{"templates/module/domain.service.go.tmpl", data.PackageName + ".service.go"},
	}
	for _, step := range steps {
		body, err := renderTemplate(step.tmpl, data)
		if err != nil {
			return err
		}
		target := filepath.Join(modDir, step.dest)
		if err := writeFileIfMissing(target, body, opts.Force); err != nil {
			return err
		}
	}
	return SyncRegistry(proj)
}

// AddTypedefs writes the domain .graphql SDL file.
func AddTypedefs(opts ModuleGenOptions) error {
	if opts.Name == "" {
		return fmt.Errorf("module name is required")
	}
	proj, err := projectFor(opts.ProjectDir)
	if err != nil {
		return err
	}
	modDir := proj.moduleDir(opts.Name)
	data := moduleData(opts.Name)
	body, err := renderTemplate("templates/module/domain.graphql.tmpl", data)
	if err != nil {
		return err
	}
	if err := writeFileIfMissing(filepath.Join(modDir, data.PackageName+".graphql"), body, opts.Force); err != nil {
		return err
	}
	return ensureModuleGo(proj, opts.Name, opts.Force)
}

// AddResolvers writes *.resolvers.go and optional model/service stubs.
func AddResolvers(opts ModuleGenOptions) error {
	if opts.Name == "" {
		return fmt.Errorf("module name is required")
	}
	proj, err := projectFor(opts.ProjectDir)
	if err != nil {
		return err
	}
	modDir := proj.moduleDir(opts.Name)
	data := moduleData(opts.Name)
	for _, step := range []struct {
		tmpl string
		dest string
	}{
		{"templates/module/domain.resolvers.go.tmpl", data.PackageName + ".resolvers.go"},
		{"templates/module/domain.model.go.tmpl", data.PackageName + ".model.go"},
		{"templates/module/domain.service.go.tmpl", data.PackageName + ".service.go"},
	} {
		body, err := renderTemplate(step.tmpl, data)
		if err != nil {
			return err
		}
		if err := writeFileIfMissing(filepath.Join(modDir, step.dest), body, opts.Force); err != nil {
			return err
		}
	}
	if err := ensureModuleGo(proj, opts.Name, opts.Force); err != nil {
		return err
	}
	return SyncRegistry(proj)
}

// AddSchema writes only the domain .graphql file.
func AddSchema(opts ModuleGenOptions) error {
	return AddTypedefs(opts)
}

func ensureModuleGo(proj *Project, name string, force bool) error {
	modDir := proj.moduleDir(name)
	path := filepath.Join(modDir, "module.go")
	if !force {
		if _, err := os.Stat(path); err == nil {
			return SyncRegistry(proj)
		}
	}
	body, err := renderTemplate("templates/module/module.go.tmpl", moduleData(name))
	if err != nil {
		return err
	}
	if err := writeFileIfMissing(path, body, force); err != nil {
		return err
	}
	return SyncRegistry(proj)
}

func projectFor(dir string) (*Project, error) {
	if dir != "" {
		abs, err := filepath.Abs(dir)
		if err != nil {
			return nil, err
		}
		return FindProject(abs)
	}
	return FindProject("")
}

// SyncRegistry rewrites src/schema/modules.go (or legacy modules/registry.go).
func SyncRegistry(proj *Project) error {
	names, err := proj.moduleNamesSorted()
	if err != nil {
		return err
	}
	var loaders []string
	for _, n := range names {
		if proj.hasLoaderFactories(n) {
			loaders = append(loaders, n)
		}
	}
	regPkg := "modules"
	if proj.Layout == LayoutSrc {
		regPkg = "schema"
	}
	body, err := renderTemplate("templates/registry.go.tmpl", struct {
		RegistryPackage    string
		ModuleImportPrefix string
		Imports            []string
		LoaderPackages     []string
	}{
		RegistryPackage:    regPkg,
		ModuleImportPrefix: proj.ModuleImportPrefix,
		Imports:            names,
		LoaderPackages:     loaders,
	})
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(proj.RegistryPath), 0o755); err != nil {
		return err
	}
	return os.WriteFile(proj.RegistryPath, body, 0o644)
}

func (p *Project) moduleNamesSorted() ([]string, error) {
	entries, err := os.ReadDir(p.ModulesDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("modules folder not found under %s (run gogql init first)", p.Root)
		}
		return nil, err
	}
	var names []string
	for _, e := range entries {
		if e.IsDir() && e.Name() != "common" {
			names = append(names, e.Name())
		}
	}
	sort.Strings(names)
	if hasDir(filepath.Join(p.ModulesDir, "common")) {
		names = append([]string{"common"}, names...)
	}
	return names, nil
}

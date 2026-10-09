/*
|--------------------------------------------------------------------------
| Module generators
|--------------------------------------------------------------------------
|
| gogql module add | module typedefs | module resolvers | module schema
|
*/

package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// ModuleGenOptions configures module file generation.
type ModuleGenOptions struct {
	ProjectDir string // app root; empty = find from cwd
	Name       string
	Force      bool
}

// AddModule writes module.go, typedefs.go, resolvers.go, and schema/*.graphql.
func AddModule(opts ModuleGenOptions) error {
	if opts.Name == "" {
		return fmt.Errorf("module name is required")
	}
	proj, err := projectFor(opts.ProjectDir)
	if err != nil {
		return err
	}
	modDir := proj.moduleDir(opts.Name)
	if err := os.MkdirAll(filepath.Join(modDir, "schema"), 0o755); err != nil {
		return err
	}
	data := moduleData(opts.Name)
	steps := []struct {
		tmpl string
		dest string
	}{
		{"templates/module/module.go.tmpl", "module.go"},
		{"templates/module/typedefs.go.tmpl", "typedefs.go"},
		{"templates/module/resolvers.go.tmpl", "resolvers.go"},
		{"templates/module/schema_entity.graphql.tmpl", "schema/" + entitySchemaFile(data.TypeName)},
		{"templates/module/schema_query.graphql.tmpl", "schema/query.graphql"},
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

// AddTypedefs writes typedefs.go and schema/*.graphql for a module.
func AddTypedefs(opts ModuleGenOptions) error {
	if opts.Name == "" {
		return fmt.Errorf("module name is required")
	}
	proj, err := projectFor(opts.ProjectDir)
	if err != nil {
		return err
	}
	modDir := proj.moduleDir(opts.Name)
	if err := os.MkdirAll(filepath.Join(modDir, "schema"), 0o755); err != nil {
		return err
	}
	data := moduleData(opts.Name)
	for _, step := range []struct {
		tmpl string
		dest string
	}{
		{"templates/module/typedefs.go.tmpl", "typedefs.go"},
		{"templates/module/schema_entity.graphql.tmpl", "schema/" + entitySchemaFile(data.TypeName)},
		{"templates/module/schema_query.graphql.tmpl", "schema/query.graphql"},
	} {
		body, err := renderTemplate(step.tmpl, data)
		if err != nil {
			return err
		}
		if err := writeFileIfMissing(filepath.Join(modDir, step.dest), body, opts.Force); err != nil {
			return err
		}
	}
	return ensureModuleGo(proj, opts.Name, opts.Force)
}

// AddResolvers writes resolvers.go and ensures module.go exists.
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
	body, err := renderTemplate("templates/module/resolvers.go.tmpl", data)
	if err != nil {
		return err
	}
	if err := writeFileIfMissing(filepath.Join(modDir, "resolvers.go"), body, opts.Force); err != nil {
		return err
	}
	if err := ensureModuleGo(proj, opts.Name, opts.Force); err != nil {
		return err
	}
	return SyncRegistry(proj)
}

// AddSchema writes only schema/*.graphql under the module.
func AddSchema(opts ModuleGenOptions) error {
	if opts.Name == "" {
		return fmt.Errorf("module name is required")
	}
	proj, err := projectFor(opts.ProjectDir)
	if err != nil {
		return err
	}
	modDir := proj.moduleDir(opts.Name)
	if err := os.MkdirAll(filepath.Join(modDir, "schema"), 0o755); err != nil {
		return err
	}
	data := moduleData(opts.Name)
	for _, step := range []struct {
		tmpl string
		dest string
	}{
		{"templates/module/schema_entity.graphql.tmpl", "schema/" + entitySchemaFile(data.TypeName)},
		{"templates/module/schema_query.graphql.tmpl", "schema/query.graphql"},
	} {
		body, err := renderTemplate(step.tmpl, data)
		if err != nil {
			return err
		}
		if err := writeFileIfMissing(filepath.Join(modDir, step.dest), body, opts.Force); err != nil {
			return err
		}
	}
	return nil
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

func entitySchemaFile(typeName string) string {
	if typeName == "" {
		return "entity.graphql"
	}
	return strings.ToLower(typeName[:1]) + typeName[1:] + ".graphql"
}

// SyncRegistry rewrites modules/registry.go from module subdirectories.
func SyncRegistry(proj *Project) error {
	names, err := proj.ModuleNames()
	if err != nil {
		return err
	}
	sort.Strings(names)
	var loaders []string
	for _, n := range names {
		if proj.hasLoaderFactories(n) {
			loaders = append(loaders, n)
		}
	}
	body, err := renderTemplate("templates/registry.go.tmpl", struct {
		ModulePath      string
		Imports         []string
		LoaderPackages  []string
	}{
		ModulePath:     proj.ModulePath,
		Imports:        names,
		LoaderPackages: loaders,
	})
	if err != nil {
		return err
	}
	regPath := filepath.Join(proj.ModulesDir, "registry.go")
	if err := os.MkdirAll(proj.ModulesDir, 0o755); err != nil {
		return err
	}
	return os.WriteFile(regPath, body, 0o644)
}

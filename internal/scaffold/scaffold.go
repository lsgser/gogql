package scaffold

import (
	"bytes"
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

//go:embed templates/*
var templateFiles embed.FS

// InitOptions configures project scaffolding.
type InitOptions struct {
	Dir          string
	ModulePath   string
	GogqlReplace string
}

// Init writes a new gogql server project into dir.
func Init(opts InitOptions) error {
	dir := opts.Dir
	if dir == "" {
		dir = "."
	}
	modulePath := opts.ModulePath
	if modulePath == "" {
		modulePath = filepath.Base(mustAbs(dir))
	}
	replace := opts.GogqlReplace
	if replace == "" {
		replace = "../gogql"
	}

	files := map[string]string{
		"main.go":                              "templates/main.go.tmpl",
		"go.mod":                               "templates/go.mod.tmpl",
		"modules/registry.go":                  "templates/modules_registry.go.tmpl",
		"modules/users/module.go":              "templates/modules_users_module.go.tmpl",
		"modules/users/typedefs.go":            "templates/modules_users_typedefs.go.tmpl",
		"modules/users/resolvers.go":           "templates/modules_users_resolvers.go.tmpl",
		"modules/users/schema/user.graphql":    "templates/modules_users_schema_user.graphql.tmpl",
		"modules/users/schema/query.graphql":   "templates/modules_users_schema_query.graphql.tmpl",
	}

	data := struct {
		ModulePath   string
		GogqlReplace string
	}{
		ModulePath:   modulePath,
		GogqlReplace: replace,
	}

	for rel, tmplPath := range files {
		raw, err := templateFiles.ReadFile(tmplPath)
		if err != nil {
			return fmt.Errorf("read template %s: %w", tmplPath, err)
		}
		t, err := template.New(rel).Parse(string(raw))
		if err != nil {
			return fmt.Errorf("parse template %s: %w", rel, err)
		}
		var buf bytes.Buffer
		if err := t.Execute(&buf, data); err != nil {
			return fmt.Errorf("execute template %s: %w", rel, err)
		}
		target := filepath.Join(dir, rel)
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(target, buf.Bytes(), 0o644); err != nil {
			return err
		}
	}
	return nil
}

func mustAbs(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		return strings.TrimPrefix(path, "./")
	}
	return filepath.Base(abs)
}

/*
|--------------------------------------------------------------------------
| Project scaffold
|--------------------------------------------------------------------------
|
| Implements gogql init: writes the default src/ tree (app, common, config,
| utils, modules/user, go.mod, main.go) from embedded templates, then
| SyncRegistry to generate src/schema/modules.go. Templates live under
| templates/init/ and are embedded via //go:embed all:templates.
|
| InitOptions control output directory, Go module path, and optional replace
| directive for local gogql development.
|
| Key type: InitOptions. Key func: Init.
|
*/

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

//go:embed all:templates
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
		"main.go":                              "templates/init/main.go.tmpl",
		"go.mod":                               "templates/go.mod.tmpl",
		"src/app/app.go":                       "templates/init/src_app_app.go.tmpl",
		"src/config/config.go":                 "templates/init/src_config_config.go.tmpl",
		"src/utils/auth.go":                    "templates/init/src_utils_auth.go.tmpl",
		"src/common/common.go":                 "templates/init/src_common_common.go.tmpl",
		"src/common/scalars.graphql":           "templates/init/src_common_scalars.graphql.tmpl",
		"src/modules/user/module.go":           "templates/init/src_modules_user_module.go.tmpl",
		"src/modules/user/user.graphql":        "templates/init/src_modules_user_user.graphql.tmpl",
		"src/modules/user/user.resolvers.go":   "templates/init/src_modules_user_user.resolvers.go.tmpl",
		"src/modules/user/user.model.go":       "templates/init/src_modules_user_user.model.go.tmpl",
		"src/modules/user/user.service.go":     "templates/init/src_modules_user_user.service.go.tmpl",
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
	proj, err := FindProject(dir)
	if err != nil {
		return err
	}
	return SyncRegistry(proj)
}

func mustAbs(path string) string {
	abs, err := filepath.Abs(path)
	if err != nil {
		return strings.TrimPrefix(path, "./")
	}
	return filepath.Base(abs)
}

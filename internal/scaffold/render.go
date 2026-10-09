/*
|--------------------------------------------------------------------------
| Template rendering
|--------------------------------------------------------------------------
|
| Executes embedded scaffold templates into project files.
|
*/

package scaffold

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

type moduleTemplateData struct {
	PackageName string
	ModuleID    string
	TypeName    string
	QueryOne    string
	QueryMany   string
}

func moduleData(name string) moduleTemplateData {
	pkg := SanitizePackageName(name)
	one, many := QueryFieldNames(pkg)
	return moduleTemplateData{
		PackageName: pkg,
		ModuleID:    pkg,
		TypeName:    TypeName(pkg),
		QueryOne:    one,
		QueryMany:   many,
	}
}

func renderTemplate(tmplPath string, data any) ([]byte, error) {
	raw, err := templateFiles.ReadFile(tmplPath)
	if err != nil {
		return nil, fmt.Errorf("read template %s: %w", tmplPath, err)
	}
	t, err := template.New(filepath.Base(tmplPath)).Parse(string(raw))
	if err != nil {
		return nil, fmt.Errorf("parse template: %w", err)
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("execute template: %w", err)
	}
	return buf.Bytes(), nil
}

func writeFileIfMissing(path string, content []byte, force bool) error {
	if !force {
		if _, err := os.Stat(path); err == nil {
			return fmt.Errorf("file already exists: %s (use -force)", path)
		}
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, content, 0o644)
}

package gogql

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// JoinTypeDefs concatenates SDL fragments into one module document (blank parts are skipped).
func JoinTypeDefs(parts ...string) string {
	var b strings.Builder
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if b.Len() > 0 {
			b.WriteByte('\n')
		}
		b.WriteString(p)
	}
	return b.String()
}

func compileModuleTypeDefs(cfg ModuleConfig) (string, error) {
	parts := make([]string, 0, 1+len(cfg.TypeDefParts)+len(cfg.TypeDefFiles))
	if cfg.TypeDefs != "" {
		parts = append(parts, cfg.TypeDefs)
	}
	parts = append(parts, cfg.TypeDefParts...)
	for _, path := range cfg.TypeDefFiles {
		raw, err := os.ReadFile(path)
		if err != nil {
			return "", fmt.Errorf("gogql: read type definitions %q: %w", path, err)
		}
		parts = append(parts, string(raw))
	}
	if cfg.TypeDefsFS != nil {
		fromFS, err := LoadTypeDefsFS(cfg.TypeDefsFS, cfg.TypeDefsFSPath)
		if err != nil {
			return "", err
		}
		if fromFS != "" {
			parts = append(parts, fromFS)
		}
	}
	merged := JoinTypeDefs(parts...)
	if merged == "" {
		return "", errModuleTypeDefsRequired
	}
	return merged, nil
}

// LoadTypeDefsFS reads all .graphql files under dir in fsys (sorted by path) and joins them.
func LoadTypeDefsFS(fsys fs.FS, dir string) (string, error) {
	if fsys == nil {
		return "", nil
	}
	dir = strings.TrimPrefix(filepath.ToSlash(dir), "./")
	var paths []string
	err := fs.WalkDir(fsys, dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if strings.HasSuffix(strings.ToLower(d.Name()), ".graphql") {
			paths = append(paths, path)
		}
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("gogql: walk type definitions: %w", err)
	}
	sort.Strings(paths)
	parts := make([]string, 0, len(paths))
	for _, path := range paths {
		raw, err := fs.ReadFile(fsys, path)
		if err != nil {
			return "", fmt.Errorf("gogql: read %q: %w", path, err)
		}
		parts = append(parts, string(raw))
	}
	return JoinTypeDefs(parts...), nil
}

// MustLoadTypeDefsFS calls LoadTypeDefsFS and panics on error.
func MustLoadTypeDefsFS(fsys fs.FS, dir string) string {
	s, err := LoadTypeDefsFS(fsys, dir)
	if err != nil {
		panic(err)
	}
	return s
}

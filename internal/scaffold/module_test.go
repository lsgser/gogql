package scaffold

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAddModuleAndSyncRegistry(t *testing.T) {
	dir := t.TempDir()
	writeGoMod(t, dir, "example.com/app")
	if err := Init(InitOptions{Dir: dir, ModulePath: "example.com/app", GogqlReplace: ""}); err != nil {
		t.Fatal(err)
	}
	if err := AddModule(ModuleGenOptions{ProjectDir: dir, Name: "posts"}); err != nil {
		t.Fatal(err)
	}
	for _, rel := range []string{
		"modules/posts/module.go",
		"modules/posts/typedefs.go",
		"modules/posts/resolvers.go",
		"modules/posts/schema/post.graphql",
		"modules/posts/schema/query.graphql",
	} {
		if _, err := os.Stat(filepath.Join(dir, rel)); err != nil {
			t.Fatalf("missing %s: %v", rel, err)
		}
	}
	reg, err := os.ReadFile(filepath.Join(dir, "modules/registry.go"))
	if err != nil {
		t.Fatal(err)
	}
	if !contains(string(reg), "posts.Module()") || !contains(string(reg), "users.Module()") {
		t.Fatalf("registry:\n%s", reg)
	}
}

func writeGoMod(t *testing.T, dir, module string) {
	t.Helper()
	// minimal tree for FindProject; Init will overwrite go.mod
	if err := os.MkdirAll(filepath.Join(dir, "modules"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module "+module+"\n\ngo 1.24\n"), 0o644); err != nil {
		t.Fatal(err)
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 || indexOf(s, sub) >= 0)
}

func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}

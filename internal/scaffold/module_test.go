/*
|--------------------------------------------------------------------------
| Module generator tests
|--------------------------------------------------------------------------
|
| Integration test: Init creates src/ layout, AddModule adds product domain,
| and SyncRegistry lists both user and product in src/schema/modules.go.
|
*/

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
	if err := AddModule(ModuleGenOptions{ProjectDir: dir, Name: "product"}); err != nil {
		t.Fatal(err)
	}
	pkg := "product"
	for _, rel := range []string{
		"src/modules/" + pkg + "/module.go",
		"src/modules/" + pkg + "/" + pkg + ".graphql",
		"src/modules/" + pkg + "/" + pkg + ".resolvers.go",
		"src/modules/" + pkg + "/" + pkg + ".model.go",
		"src/modules/" + pkg + "/" + pkg + ".service.go",
	} {
		if _, err := os.Stat(filepath.Join(dir, rel)); err != nil {
			t.Fatalf("missing %s: %v", rel, err)
		}
	}
	reg, err := os.ReadFile(filepath.Join(dir, "src/schema/modules.go"))
	if err != nil {
		t.Fatal(err)
	}
	if !contains(string(reg), "product.Module()") || !contains(string(reg), "user.Module()") {
		t.Fatalf("registry:\n%s", reg)
	}
}

func writeGoMod(t *testing.T, dir, module string) {
	t.Helper()
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

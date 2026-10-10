/*
|--------------------------------------------------------------------------
| Type definitions tests
|--------------------------------------------------------------------------
|
| Covers inline-only ModuleConfig (legacy API), JoinTypeDefs, TypeDefParts,
| and LoadTypeDefsFS using testdata/schema/*.graphql embed fixtures.
|
*/

package core_test

import (
	"embed"
	"strings"
	"testing"

	"github.com/lsgser/gogql"
)

//go:embed testdata/schema/*.graphql
var testSchema embed.FS

func TestModuleBackwardCompatInlineTypeDefs(t *testing.T) {
	mod, err := gogql.NewModule(gogql.ModuleConfig{
		ID:       "legacy",
		TypeDefs: `type Query { ping: String! }`,
	})
	if err != nil {
		t.Fatal(err)
	}
	if mod.TypeDefs() != `type Query { ping: String! }` {
		t.Fatalf("unexpected SDL: %s", mod.TypeDefs())
	}
}

func TestJoinTypeDefsAndModuleParts(t *testing.T) {
	mod, err := gogql.NewModule(gogql.ModuleConfig{
		ID: "test",
		TypeDefParts: []string{
			`type Query { hello: String! }`,
		},
		TypeDefsFS:     testSchema,
		TypeDefsFSPath: "testdata/schema",
	})
	if err != nil {
		t.Fatal(err)
	}
	sdl := mod.TypeDefs()
	for _, part := range []string{"type User", "type Query", "hello", "user"} {
		if !strings.Contains(sdl, part) {
			t.Fatalf("merged SDL missing %q:\n%s", part, sdl)
		}
	}
}

/*
|--------------------------------------------------------------------------
| Security tests
|--------------------------------------------------------------------------
|
| Tests SecurityConfig.MaxDepth query validation.
|
*/

package core_test

import (
	"context"
	"testing"

	"github.com/lsgser/gogql"
)

func TestMaxDepth(t *testing.T) {
	app := gogql.MustApplication(gogql.ApplicationConfig{
		Modules: []*gogql.Module{gogql.MustModule(gogql.ModuleConfig{
			ID: "nested",
			TypeDefs: `
				type Query { node: Node }
				type Node { child: Node name: String! }
			`,
			Resolvers: gogql.NewResolverMap().Query("node", func(_ context.Context) (*node, error) {
				return &node{name: "root", child: &node{name: "deep"}}, nil
			}),
		})},
		Security: gogql.SecurityConfig{MaxDepth: 2},
	})

	res := app.Execute(context.Background(), `{ node { name child { name child { name } } } } }`, "", nil)
	if len(res.Errors) == 0 {
		t.Fatal("expected depth validation error")
	}
}

type node struct {
	name  string
	child *node
}

func (n *node) Name() string {
	if n == nil {
		return ""
	}
	return n.name
}

func (n *node) Child() *node {
	if n == nil {
		return nil
	}
	return n.child
}

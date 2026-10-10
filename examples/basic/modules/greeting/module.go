/*
|--------------------------------------------------------------------------
| Basic example — greeting module (inline)
|--------------------------------------------------------------------------
|
| Approach A: entire module in one file—TypeDefs string plus ResolverMap on
| ModuleConfig. Use for tiny domains; split into .graphql + resolvers when
| the module grows (see modules/users).
|
*/

package greeting

import (
	"context"

	"github.com/lsgser/gogql"
)

// Module adds a second Query field; gogql merges SDL from multiple modules.
func Module() *gogql.Module {
	return gogql.MustModule(gogql.ModuleConfig{
		ID: "greeting",
		TypeDefs: `
			type Query {
				hello: String!
			}
		`,
		Resolvers: gogql.NewResolverMap().Query("hello", func(_ context.Context) (string, error) {
			return "Hello from gogql", nil
		}),
	})
}

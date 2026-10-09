/*
|--------------------------------------------------------------------------
| Basic example — greeting module (inline)
|--------------------------------------------------------------------------
|
| Approach A: single file with inline TypeDefs and ResolverMap (hello query).
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

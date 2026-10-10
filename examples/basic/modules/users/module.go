/*
|--------------------------------------------------------------------------
| Basic example — users module wiring
|--------------------------------------------------------------------------
|
| Approach B: MustModule with typeDefs() from embed and resolvers() from
| resolvers.go—mirrors src/modules/<domain>/module.go in the CLI scaffold.
|
*/

package users

import "github.com/lsgser/gogql"

// Module composes split typedefs (.graphql) and resolvers (resolvers.go).
func Module() *gogql.Module {
	return gogql.MustModule(gogql.ModuleConfig{
		ID:        "users",
		TypeDefs:  typeDefs(),
		Resolvers: resolvers(),
	})
}

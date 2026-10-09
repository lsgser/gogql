/*
|--------------------------------------------------------------------------
| Basic example — users module wiring
|--------------------------------------------------------------------------
|
| Approach B: connects embedded schema SDL and resolvers.go to MustModule.
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

/*
|--------------------------------------------------------------------------
| Basic example — users resolvers
|--------------------------------------------------------------------------
|
| ResolverMap and User type for the user(id) query field.
|
*/

package users

import (
	"context"

	"github.com/graph-gophers/graphql-go"
	"github.com/lsgser/gogql"
)

// User is the GraphQL-facing shape returned by resolvers.
type User struct {
	ID   graphql.ID
	Name string
}

func resolvers() *gogql.ResolverMap {
	return gogql.NewResolverMap().Query("user", resolveUser)
}

func resolveUser(_ context.Context, args struct{ ID graphql.ID }) (*User, error) {
	return &User{ID: args.ID, Name: "GraphQL Modules in Go"}, nil
}

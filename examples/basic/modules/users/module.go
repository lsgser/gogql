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

// Module is the users feature slice (SDL + resolvers).
func Module() *gogql.Module {
	return gogql.MustModule(gogql.ModuleConfig{
		ID: "users",
		TypeDefs: `
			type User {
				id: ID!
				name: String!
			}
			type Query {
				user(id: ID!): User
			}
		`,
		Resolvers: gogql.NewResolverMap().Query("user", resolveUser),
	})
}

func resolveUser(_ context.Context, args struct{ ID graphql.ID }) (*User, error) {
	return &User{ID: args.ID, Name: "GraphQL Modules in Go"}, nil
}

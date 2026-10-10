/*
|--------------------------------------------------------------------------
| Database example — users module (inline)
|--------------------------------------------------------------------------
|
| Single-file module with inline SDL, Provide(repo), ResolverMap queries,
| resolveMe using AuthClaimsFrom, and LoaderFactories for batched ByID.
| Shows DI + auth + loaders without the src/ folder layout.
|
*/

package users

import (
	"context"
	"fmt"

	"github.com/graph-gophers/dataloader"
	"github.com/graph-gophers/graphql-go"
	"github.com/lsgser/gogql"
	"github.com/lsgser/gogql/examples/database/store"
)

// User is the GraphQL type returned by this module.
type User struct {
	ID    graphql.ID
	Name  string
	Email string
}

// Module builds the users gogql module with injected store access.
func Module(repo *store.UserRepository) *gogql.Module {
	return gogql.MustModule(gogql.ModuleConfig{
		ID: "users",
		TypeDefs: `
			type User {
				id: ID!
				name: String!
				email: String!
			}
			type Query {
				user(id: ID!): User
				users: [User!]!
				me: User
			}
		`,
		Providers: []gogql.Provider{
			gogql.Provide(repo),
		},
		Resolvers: gogql.NewResolverMap().
			Query("user", resolveUser).
			Query("users", resolveUsers).
			Query("me", resolveMe),
	})
}

// LoaderFactories registers per-request DataLoaders for this module.
func LoaderFactories(repo *store.UserRepository) gogql.LoaderFactories {
	return gogql.LoaderFactories{
		"user": func() *dataloader.Loader {
			return gogql.NewStringKeyLoader(func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
				out := make([]*dataloader.Result, len(keys))
				for i, k := range keys {
					rec, err := repo.ByID(ctx, k.String())
					if err != nil {
						out[i] = &dataloader.Result{Error: err}
						continue
					}
					out[i] = &dataloader.Result{Data: toGraphQL(rec)}
				}
				return out
			})
		},
	}
}

func resolveUser(ctx context.Context, args struct{ ID graphql.ID }) (*User, error) {
	loader, err := gogql.LoadersFromContext(ctx, "user")
	if err != nil {
		return nil, err
	}
	v, err := loader.Load(ctx, dataloader.StringKey(string(args.ID)))()
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, nil
	}
	return v.(*User), nil
}

func resolveUsers(ctx context.Context) ([]*User, error) {
	repo := gogql.MustGet[*store.UserRepository](ctx)
	list, err := repo.List(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]*User, 0, len(list))
	for i := range list {
		out = append(out, toGraphQL(&list[i]))
	}
	return out, nil
}

func resolveMe(ctx context.Context) (*User, error) {
	claims, ok := gogql.AuthClaimsFrom(ctx)
	if !ok || claims.Subject == "" {
		return nil, fmt.Errorf("unauthenticated: send Authorization: Bearer <jwt>")
	}
	loader, err := gogql.LoadersFromContext(ctx, "user")
	if err != nil {
		return nil, err
	}
	v, err := loader.Load(ctx, dataloader.StringKey(claims.Subject))()
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, fmt.Errorf("user not found for sub %q", claims.Subject)
	}
	return v.(*User), nil
}

func toGraphQL(rec *store.UserRecord) *User {
	if rec == nil {
		return nil
	}
	return &User{ID: graphql.ID(rec.ID), Name: rec.Name, Email: rec.Email}
}

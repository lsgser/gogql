/*
|--------------------------------------------------------------------------
| Integration tests
|--------------------------------------------------------------------------
|
| End-to-end checks that multiple Module SDL fragments merge into one schema,
| ResolverMap fields resolve, and provider injection is visible in resolvers.
| Guards regressions in merge + resolver root wiring used by real applications.
|
*/

package core_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/graph-gophers/graphql-go"
	"github.com/lsgser/gogql"
)

type User struct {
	ID   graphql.ID
	Name string
}

type UsersService struct {
	users map[string]User
}

func (s *UsersService) ByID(id string) *User {
	u, ok := s.users[id]
	if !ok {
		return nil
	}
	return &u
}

func TestModularSDLSchema(t *testing.T) {
	usersModule := gogql.MustModule(gogql.ModuleConfig{
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
		Providers: []gogql.Provider{
			gogql.Provide(&UsersService{
				users: map[string]User{"1": {ID: "1", Name: "Ada"}},
			}),
		},
		Resolvers: gogql.NewResolverMap().Query("user", func(ctx context.Context, args struct {
			ID graphql.ID
		}) (*User, error) {
			svc := gogql.MustGet[*UsersService](ctx)
			return svc.ByID(string(args.ID)), nil
		}),
	})

	postsModule := gogql.MustModule(gogql.ModuleConfig{
		ID: "posts",
		TypeDefs: `
			type Post {
				id: ID!
				title: String!
			}
			type Query {
				posts: [Post!]!
			}
		`,
		Resolvers: gogql.NewResolverMap().Query("posts", func(_ context.Context) ([]*Post, error) {
			return []*Post{{ID: "p1", Title: "Hello modules"}}, nil
		}),
	})

	app := gogql.MustApplication(gogql.ApplicationConfig{
		Modules: []*gogql.Module{usersModule, postsModule},
	})

	res := app.Execute(context.Background(), `
		query {
			user(id: "1") { id name }
			posts { id title }
		}
	`, "", nil)

	if len(res.Errors) > 0 {
		t.Fatalf("unexpected errors: %+v", res.Errors)
	}
	var payload struct {
		User struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"user"`
		Posts []struct {
			ID    string `json:"id"`
			Title string `json:"title"`
		} `json:"posts"`
	}
	if err := json.Unmarshal(res.Data, &payload); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if payload.User.Name != "Ada" {
		t.Fatalf("user name: got %q", payload.User.Name)
	}
	if len(payload.Posts) != 1 || payload.Posts[0].Title != "Hello modules" {
		t.Fatalf("posts: %+v", payload.Posts)
	}
}

type Post struct {
	ID    graphql.ID
	Title string
}

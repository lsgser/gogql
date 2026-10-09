# Getting started

This guide walks through a minimal GraphQL server with **SDL modules**, an HTTP endpoint, and the **playground**.

## 1. Create a project

```bash
mkdir my-api && cd my-api
go mod init github.com/you/my-api
go get github.com/lsgser/gogql@latest
```

Or scaffold from the gogql repo: [Installation — CLI](installation.md#install-the-cli-optional).

## 2. Define a module (GraphQL SDL)

A **module** is a slice of schema plus resolvers. Types are written in GraphQL, not as Go structs for the schema itself:

```go
package main

import (
    "context"
    "log"

    "github.com/graph-gophers/graphql-go"
    "github.com/lsgser/gogql"
)

type User struct {
    ID   graphql.ID
    Name string
}

func main() {
    users := gogql.MustModule(gogql.ModuleConfig{
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
        Resolvers: gogql.NewResolverMap().Query("user", func(_ context.Context, args struct {
            ID graphql.ID
        }) (*User, error) {
            return &User{ID: args.ID, Name: "Ada Lovelace"}, nil
        }),
    })

    app := gogql.MustApplication(gogql.ApplicationConfig{
        Modules: []*gogql.Module{users},
    })

    server := gogql.NewServer(app, gogql.ServerConfig{
        Playground: &gogql.PlaygroundConfig{
            Enabled: true,
            Path:    "/playground",
            UI:      gogql.PlaygroundGraphiQL,
        },
    })

    log.Println("http://localhost:8080/graphql")
    log.Println("http://localhost:8080/playground")
    log.Fatal(server.ListenAndServe(":8080"))
}
```

## 3. Run and query

```bash
go run .
```

Open [http://localhost:8080/playground](http://localhost:8080/playground) and run:

```graphql
query {
  user(id: "1") {
    id
    name
  }
}
```

Or POST JSON to `http://localhost:8080/graphql`:

```bash
curl -s -X POST http://localhost:8080/graphql \
  -H 'Content-Type: application/json' \
  -d '{"query":"{ user(id: \"1\") { id name } }"}'
```

## 4. Multiple modules

Each module can declare its own `type Query { ... }`. gogql merges them (see [Modules & resolvers](modules-and-resolvers.md)).

## 5. What’s next

- [Modules & resolvers](modules-and-resolvers.md) — providers, DataLoaders, subscription resolvers
- [Server & playground](server-and-playground.md) — paths, Apollo Sandbox, depth limits, WebSockets
- [`examples/basic`](../examples/basic) in the repository

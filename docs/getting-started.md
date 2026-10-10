# Getting started

This guide gets you to a running GraphQL server with **SDL modules**, HTTP **`/graphql`**, and the **playground**.

## Path 1 — Scaffold (recommended)

Create a project with the default **`src/`** layout ([Project layout](project-layout.md)):

```bash
go install github.com/lsgser/gogql/cmd/gogql@latest
# or: go get -tool github.com/lsgser/gogql/cmd/gogql@latest  (inside your module)

gogql init my-api -module github.com/you/my-api -gogql=
cd my-api
go mod tidy
go run .
```

- GraphQL: [http://localhost:8080/graphql](http://localhost:8080/graphql)
- Playground: [http://localhost:8080/playground](http://localhost:8080/playground)

Try:

```graphql
query {
  user(id: "1") {
    id
    name
    email
  }
}
```

Add another domain:

```bash
go tool gogql module add product   # if go.mod has tool github.com/lsgser/gogql/cmd/gogql
# or: gogql module add product
go run .
```

See [Installation — CLI commands](installation.md#commands) for `module typedefs`, `module resolvers`, and flags.

## Path 2 — Manual minimal server

Use this when you want a single file or a custom folder tree without the CLI.

### 1. Create a module

A **module** is SDL plus resolvers. Types are written in GraphQL, not as Go structs for the schema itself.

**Approach A** — inline SDL and `ResolverMap` in one place ([Modules & resolvers](modules-and-resolvers.md#approach-a--inline-module-backward-compatible)):

```bash
mkdir my-api && cd my-api
go mod init github.com/you/my-api
go get github.com/lsgser/gogql@latest
```

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

### 2. Run and query

```bash
go run .
```

Or POST JSON:

```bash
curl -s -X POST http://localhost:8080/graphql \
  -H 'Content-Type: application/json' \
  -d '{"query":"{ user(id: \"1\") { id name } }"}'
```

## Multiple modules

Each module can declare its own `type Query { ... }`. gogql merges them ([Modules & resolvers — merging SDL](modules-and-resolvers.md#merging-sdl)). In the **scaffold**, register modules in **`src/schema/modules.go`** (updated by the CLI); in manual apps, pass a slice to **`ApplicationConfig.Modules`**.

## What’s next

- [Project layout](project-layout.md) — default `src/` tree, domain files, CLI generators
- [Modules & resolvers](modules-and-resolvers.md) — inline vs split SDL, service/model layers, DI, loaders, subscriptions
- [Server & playground](server-and-playground.md) — paths, WebSocket, depth limits
- [Features](features.md) — capability index
- [`examples/basic`](../examples/basic) — mixed inline + split modules without `src/`

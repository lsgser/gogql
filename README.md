# gogql

[![Go Reference](https://pkg.go.dev/badge/github.com/lsgser/gogql.svg)](https://pkg.go.dev/github.com/lsgser/gogql)

Modular GraphQL server library for Go, inspired by [GraphQL Yoga](https://the-guild.dev/graphql/yoga-server) and [GraphQL Modules](https://the-guild.dev/graphql/modules).

**Repository:** [github.com/lsgser/gogql](https://github.com/lsgser/gogql)

Define your schema in **GraphQL SDL**, compose **modules**, and serve over HTTP with an optional **playground**, **subscriptions**, **DataLoaders**, and **depth limits**.

## Install

```bash
go get github.com/lsgser/gogql@latest
```

```go
import "github.com/lsgser/gogql"
```

Full details: **[Documentation](docs/README.md)** · **[Installation guide](docs/installation.md)**

## Features

- **SDL-first schema** — write types in GraphQL, not Go structs
- **GraphQL Modules** — `Module` with `id`, `typeDefs`, `resolvers`, and `providers`
- **Automatic SDL merge** — multiple modules can each declare `type Query { ... }`
- **Playground** — GraphiQL or Apollo Sandbox at `/playground` (configurable)
- **Subscriptions** — [graphql-transport-ws](https://github.com/enisdenjo/graphql-ws) on the GraphQL endpoint
- **DataLoaders** — request-scoped batched loaders
- **Depth limit** — `SecurityConfig.MaxDepth`
- **CLI** — `gogql init` scaffolds a new project

## Quick start

```go
app := gogql.MustApplication(gogql.ApplicationConfig{
    Modules:  []*gogql.Module{usersModule},
    Security: gogql.SecurityConfig{MaxDepth: 15},
})

server := gogql.NewServer(app, gogql.ServerConfig{
    Playground: &gogql.PlaygroundConfig{
        Enabled: true,
        Path:    "/playground",
        UI:      gogql.PlaygroundGraphiQL,
    },
})
log.Fatal(server.ListenAndServe(":8080"))
```

- API: `http://localhost:8080/graphql`
- Playground: `http://localhost:8080/playground`

Step-by-step: [Getting started](docs/getting-started.md)

## CLI

```bash
go install github.com/lsgser/gogql/cmd/gogql@latest
gogql init my-api
```

From a clone: `go run ./cmd/gogql init my-api`. See [Installation](docs/installation.md).

## Examples

Each example uses **gogql modules** (`MustModule` → `MustApplication` → `NewServer`). See [examples/README.md](examples/README.md).

```bash
go run ./examples/basic          # merged modules + playground
go run ./examples/subscriptions  # SubscriptionResolvers + WebSocket
go run ./examples/database       # SQL, DataLoaders, JWT ContextFunc
```

## Documentation

| Guide | Topic |
|-------|--------|
| [docs/README.md](docs/README.md) | Index |
| [installation.md](docs/installation.md) | `go get`, CLI, local replace |
| [getting-started.md](docs/getting-started.md) | First server |
| [modules-and-resolvers.md](docs/modules-and-resolvers.md) | Modules, DI, loaders, subscriptions |
| [server-and-playground.md](docs/server-and-playground.md) | HTTP, playground, WebSocket, limits |
| [database-and-auth.md](docs/database-and-auth.md) | SQL, repositories, JWT auth |

## License

MIT — see [LICENSE](LICENSE).

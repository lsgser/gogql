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

Library code lives under **`internal/core/`**; the repo root re-exports the public API so your import stays `github.com/lsgser/gogql`. See [Installation — repository layout](docs/installation.md#clone-the-gogql-repository).

## Features

- **SDL-first schema** — write types in GraphQL, not Go structs
- **GraphQL Modules** — inline or split SDL/resolvers per module (`TypeDefs` only, or `.graphql` + `resolvers.go`)
- **Automatic SDL merge** — multiple modules can each declare `type Query { ... }`
- **Playground** — GraphiQL or Apollo Sandbox at `/playground` (configurable)
- **Subscriptions** — [graphql-transport-ws](https://github.com/enisdenjo/graphql-ws) on the GraphQL endpoint
- **DataLoaders** — request-scoped batched loaders
- **Depth limit** — `SecurityConfig.MaxDepth`
- **CLI** — `gogql init` (`src/` layout), `gogql module add`, generators

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

Step-by-step: [Getting started](docs/getting-started.md) · Default folders: [Project layout](docs/project-layout.md)

## CLI

Same idea as [gofreight](https://github.com/lsgser/gofreight): install the **command** separately from the library.

```bash
go install github.com/lsgser/gogql/cmd/gogql@latest
gogql version          # lists all commands
gogql init my-api
gogql module add product   # domain module + src/schema/modules.go
```

Or pin the CLI in your app with Go 1.24+: `go get -tool github.com/lsgser/gogql/cmd/gogql@latest` → `go tool gogql init my-api`. Scaffolded projects include the `tool` line in `go.mod`. See [Installation — CLI](docs/installation.md#install-the-cli-optional).

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
| [features.md](docs/features.md) | DataLoaders, depth limits, feature index |
| [installation.md](docs/installation.md) | `go get`, CLI, local replace |
| [getting-started.md](docs/getting-started.md) | `gogql init` or manual server |
| [project-layout.md](docs/project-layout.md) | Default `src/` tree, CLI generators |
| [modules-and-resolvers.md](docs/modules-and-resolvers.md) | Modules, DI, loaders, subscriptions |
| [server-and-playground.md](docs/server-and-playground.md) | HTTP, playground, WebSocket, limits |
| [gin.md](docs/gin.md) | GraphQL with Gin alongside REST |
| [deployment.md](docs/deployment.md) | Production build, TLS, Docker, K8s sketch |
| [database-and-auth.md](docs/database-and-auth.md) | SQL, repositories, JWT auth |

## License

MIT — see [LICENSE](LICENSE).

# gogql documentation

**gogql** is a modular GraphQL server library for Go. Repository: [github.com/lsgser/gogql](https://github.com/lsgser/gogql).

## Guides

| Document | Description |
|----------|-------------|
| [Installation](installation.md) | Requirements, `go get`, repo layout, CLI install |
| [Project layout](project-layout.md) | Default **`src/`** scaffold, domain modules, CLI generators |
| [Getting started](getting-started.md) | **`gogql init`** or manual first server |
| **[Features](features.md)** | **DataLoaders, depth limits, DI, JWT, subscriptions (index)** |
| [Modules & resolvers](modules-and-resolvers.md) | SDL modules, `ResolverMap`, DI, DataLoaders |
| [Server & playground](server-and-playground.md) | HTTP, playground, subscriptions, `MaxDepth` |
| [Gin integration](gin.md) | Mount GraphQL on [Gin](https://github.com/gin-gonic/gin) with REST routes |
| [Deployment](deployment.md) | Production config, TLS, Docker, health checks, WebSocket |
| [Database & JWT](database-and-auth.md) | SQL repositories, DataLoaders + batch SQL, JWT `ContextFunc` |

## Examples in the repo

- [`examples/basic`](../examples/basic) — mixed inline + split modules (flat `modules/`)
- [`examples/subscriptions`](../examples/subscriptions) — WebSocket subscriptions
- [`examples/database`](../examples/database) — SQLite + repository + JWT

## Related projects

- [GraphQL Yoga](https://the-guild.dev/graphql/yoga-server) — inspiration for the HTTP/server ergonomics
- [GraphQL Modules](https://the-guild.dev/graphql/modules) — inspiration for modular SDL composition

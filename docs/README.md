# gogql documentation

**gogql** is a modular GraphQL server library for Go. Repository: [github.com/lsgser/gogql](https://github.com/lsgser/gogql).

## Guides

| Document | Description |
|----------|-------------|
| [Installation](installation.md) | Requirements, `go get`, CLI install, versioning |
| [Getting started](getting-started.md) | First module, server, playground |
| [Modules & resolvers](modules-and-resolvers.md) | SDL modules, `ResolverMap`, DI, DataLoaders |
| [Server & playground](server-and-playground.md) | HTTP, playground, subscriptions, security limits |
| [Database & JWT](database-and-auth.md) | SQL repositories, DataLoaders, JWT `ContextFunc` |

## Examples in the repo

- [`examples/basic`](../examples/basic) — single module + playground
- [`examples/subscriptions`](../examples/subscriptions) — WebSocket subscriptions
- [`examples/database`](../examples/database) — SQLite + repository + JWT

## Related projects

- [GraphQL Yoga](https://the-guild.dev/graphql/yoga-server) — inspiration for the HTTP/server ergonomics
- [GraphQL Modules](https://the-guild.dev/graphql/modules) — inspiration for modular SDL composition

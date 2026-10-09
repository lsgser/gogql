# gogql examples

Each example uses **gogql modules** (`MustModule` → `MustApplication` → `NewServer`). Module code can be organized in two ways (see [docs](../docs/modules-and-resolvers.md#two-ways-to-organize-a-module)):

| Style | Layout |
|-------|--------|
| **A — inline** | One `module.go` with `TypeDefs` + `Resolvers` |
| **B — split** | `module.go` + optional `typedefs.go`, `resolvers.go`, `schema/*.graphql` |

The CLI scaffold generates **B** for `users`; you can keep or simplify any module to **A**.

```
example/
  main.go
  modules/
    registry.go
    <name>/
      module.go
      ... optional split files ...
```

Run from the repository root:

| Example | Command | Demonstrates |
|---------|---------|--------------|
| [basic](basic) | `go run ./examples/basic` | **Mixed A+B**: `users` split SDL/resolvers, `greeting` inline; merged `Query` |
| [subscriptions](subscriptions) | `go run ./examples/subscriptions` | `SubscriptionResolvers`, WebSocket transport |
| [database](database) | `go run ./examples/database` | **Inline module** + providers, SQL, JWT, DataLoaders |

All examples use **`gogql.MustModule` → `gogql.MustApplication` → `gogql.NewServer`** — no direct graph-gophers HTTP wiring.

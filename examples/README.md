# gogql examples

Each example follows the same **GraphQL Modules** layout as the CLI scaffold:

```
example/
  main.go              # DB/JWT wiring, gogql.NewApplication, gogql.NewServer
  modules/
    registry.go        # All() — compose modules for the application
    <name>/
      module.go        # SDL typeDefs, ResolverMap, Providers, SubscriptionResolvers
```

Run from the repository root:

| Example | Command | Demonstrates |
|---------|---------|--------------|
| [basic](basic) | `go run ./examples/basic` | SDL modules, merged `Query`, playground |
| [subscriptions](subscriptions) | `go run ./examples/subscriptions` | `SubscriptionResolvers`, WebSocket transport |
| [database](database) | `go run ./examples/database` | Providers, SQL repository, JWT `ContextFunc`, DataLoaders |

All examples use **`gogql.MustModule` → `gogql.MustApplication` → `gogql.NewServer`** — no direct graph-gophers HTTP wiring.

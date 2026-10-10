# gogql examples

Each example uses **`gogql.MustModule` → `gogql.MustApplication` → `gogql.NewServer`**. They teach library features with **flat `modules/`** trees (no `src/` wrapper).

For a **production-shaped starter**, use **`gogql init`** instead — see [docs/project-layout.md](../docs/project-layout.md).

## Module organization in examples

| Style | Layout | Where |
|-------|--------|--------|
| **A — inline** | One `module.go` with `TypeDefs` + `Resolvers` | [`greeting`](basic/modules/greeting/module.go), [`database` users](database/modules/users/module.go) |
| **B — split** | `module.go` + `typedefs.go` + `resolvers.go` + `schema/*.graphql` | [`basic` users](basic/modules/users/) |

The CLI **`gogql init`** generates a **`src/modules/<domain>/`** layout with `*.graphql`, `*.resolvers.go`, `*.model.go`, and `*.service.go` — different from these examples but the same gogql APIs.

Run from the **repository root**:

| Example | Command | Demonstrates |
|---------|---------|--------------|
| [basic](basic) | `go run ./examples/basic` | Merged modules (inline + split), playground |
| [subscriptions](subscriptions) | `go run ./examples/subscriptions` | `SubscriptionResolvers`, WebSocket |
| [database](database) | `go run ./examples/database` | SQLite, providers, JWT, DataLoaders |

No direct graph-gophers HTTP wiring — always through **`gogql.NewServer`**.

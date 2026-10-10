# Modules & resolvers

## Ways to organize modules

gogql stays **backward compatible**: a single `TypeDefs` string and `Resolvers` in one place still works. Optional fields (`TypeDefParts`, `TypeDefFiles`, `TypeDefsFS`) build that same SDL from multiple sources. You can **mix styles** in one app ([`examples/basic`](../examples/basic)).

| | **Approach A — inline** | **Approach B — split files** | **Default scaffold (`gogql init`)** |
|---|---|---|---|
| Best for | Small modules, prototypes | Many types, `schema/` embed | Production-shaped **`src/`** apps |
| SDL | One `TypeDefs` string | `typedefs.go`, embed, parts | `<domain>.graphql` in module folder |
| Resolvers | Inline in `module.go` | `resolvers.go` | `<domain>.resolvers.go` |
| Layers | Optional | Optional | `<domain>.model.go`, `<domain>.service.go` |
| Registry | Your `main` | `modules/registry.go` | **`src/schema/modules.go`** (CLI) |
| Examples | [`greeting`](../examples/basic/modules/greeting/module.go), [`database` users](../examples/database/modules/users/module.go) | [`users`](../examples/basic/modules/users) | [Project layout](project-layout.md) |

The scaffold is **convention only**—any folder tree works if you pass `[]*gogql.Module` to `MustApplication`.

---

## Default scaffold — domain module (`gogql init`)

```text
src/modules/user/
├── module.go
├── user.graphql
├── user.resolvers.go
├── user.model.go
└── user.service.go
```

- **`user.graphql`** — SDL for this domain (types + `Query` fields; merged globally).
- **`user.resolvers.go`** — `ResolverMap`; should call **`user.service.go`**, not SQL directly.
- **`user.model.go`** — storage/record types.
- **`src/schema/modules.go`** — `All()` and `LoaderFactories()`; updated by **`gogql module add`**.

Shared scalars live in **`src/common`**. Full tree: [Project layout](project-layout.md).

---

## Approach A — inline module (backward compatible)

Everything in one file (or one `TypeDefs` string anywhere you prefer):

```go
gogql.MustModule(gogql.ModuleConfig{
    ID: "greeting",
    TypeDefs: `
        type Query { hello: String! }
    `,
    Resolvers: gogql.NewResolverMap().Query("hello", func(_ context.Context) (string, error) {
        return "Hello", nil
    }),
})
```

No `TypeDefParts`, `TypeDefFiles`, or `TypeDefsFS` required. This matches the original gogql API and [Getting started](getting-started.md).

---

## Approach B — split typedefs and resolvers

Same `MustModule` call; SDL and resolver wiring are split across files for clarity.

```text
modules/users/
├── module.go           # wires TypeDefs + Resolvers
├── typedefs.go         # embed or JoinTypeDefs
├── resolvers.go        # ResolverMap + handler funcs
└── schema/
    ├── user.graphql    # types
    └── query.graphql   # Query fields for this module
```

**Embedded `.graphql` files** (see [`examples/basic/modules/users`](../examples/basic/modules/users)):

```go
// typedefs.go
//go:embed schema/*.graphql
var schemaFS embed.FS

func typeDefs() string {
    return gogql.MustLoadTypeDefsFS(schemaFS, "schema")
}
```

```go
// module.go
gogql.MustModule(gogql.ModuleConfig{
    ID:        "users",
    TypeDefs:  typeDefs(),
    Resolvers: resolvers(), // from resolvers.go
})
```

**Multiple SDL sources on `ModuleConfig`** (no embed):

```go
gogql.ModuleConfig{
    ID: "users",
    TypeDefs: userTypesSDL,
    TypeDefParts: []string{queryFieldsSDL},
    TypeDefFiles: []string{"modules/users/schema/extra.graphql"},
}
```

Helpers: `gogql.JoinTypeDefs`, `gogql.LoadTypeDefsFS`, `gogql.MustLoadTypeDefsFS`.

Resolvers do not require a new API: use `func resolvers() *gogql.ResolverMap` in `resolvers.go` and pass it from `module.go`. Subscription methods can live in `subscription.go` and are still set via `SubscriptionResolvers`.

You can also **gradually migrate** Approach A → B: move SDL into `TypeDefParts` or embed first, keep the same `TypeDefs` string temporarily, then delete duplication once embed is in place.

---

## ModuleConfig reference

```go
gogql.MustModule(gogql.ModuleConfig{
    ID:        "users",
    TypeDefs:  `...`,              // required unless merged parts/files/FS produce SDL
    Resolvers: /* required for execution */,
    SubscriptionResolvers: /* optional */,
    Providers: []gogql.Provider{ /* optional */ },
    // optional SDL extensions (Approach B):
    TypeDefParts: []string{ ... },
    TypeDefFiles: []string{ "path/to/extra.graphql" },
    TypeDefsFS:     schemaFS,
    TypeDefsFSPath: "schema",
})
```

| Field | Required? | Purpose |
|-------|-----------|---------|
| `ID` | yes | Unique module name |
| `TypeDefs` | yes* | Primary SDL string (*or supply SDL only via parts/files/FS) |
| `Resolvers` | for queries/mutations | `*ResolverMap` or graph-gophers root struct |
| `SubscriptionResolvers` | for subscriptions | Struct with **methods** (not `ResolverMap`) |
| `Providers` | no | Injector services |
| `TypeDefParts` | no | Extra SDL fragments merged after `TypeDefs` |
| `TypeDefFiles` | no | `.graphql` paths on disk |
| `TypeDefsFS` / `TypeDefsFSPath` | no | Embedded or virtual `.graphql` tree |

## ResolverMap (queries & mutations)

Fluent registration:

```go
gogql.NewResolverMap().
    Query("user", func(ctx context.Context, args struct { ID graphql.ID }) (*User, error) {
        return fetchUser(ctx, args.ID)
    }).
    Mutation("createUser", func(ctx context.Context, args struct { Name string }) (*User, error) {
        // ...
        return nil, nil
    })
```

Resolver functions follow [graph-gophers/graphql-go](https://github.com/graph-gophers/graphql-go) rules: optional `context.Context`, an args struct when the field has arguments, return `(T, error)`.

## Dependency injection

Register providers on the module:

```go
gogql.Provide(&UsersService{})
// or
gogql.ProvideFactory(func(inj *gogql.Injector) (*UsersService, error) { ... })
```

In a resolver:

```go
svc := gogql.MustGet[*UsersService](ctx)
```

The server attaches the injector on each request via `Application.RequestContext`.

## DataLoaders

See also: [Features — DataLoaders](features.md#dataloaders) (lifecycle, `LoadMany`, example links).

Register factories on the **application** (one loader instance per HTTP/WebSocket request):

```go
app := gogql.MustApplication(gogql.ApplicationConfig{
    Modules: modules,
    Loaders: gogql.LoaderFactories{
        "user": func() *dataloader.Loader {
            return gogql.NewStringKeyLoader(batchLoadUsers)
        },
    },
})
```

In a resolver:

```go
loader, err := gogql.LoadersFromContext(ctx, "user")
if err != nil { return nil, err }
result, err := loader.Load(ctx, dataloader.StringKey(id))()
```

Use `LoadMany` when resolving a list of IDs in one resolver to batch efficiently.

## Subscriptions

`ResolverMap` cannot register subscription fields (the underlying engine requires **methods** on a subscription root). Use `SubscriptionResolvers`:

```go
gogql.MustModule(gogql.ModuleConfig{
    TypeDefs: `
        type Subscription { messageSent: Message! }
        type Message { text: String! }
    `,
    Resolvers: gogql.NewResolverMap().Query("ok", ...),
    SubscriptionResolvers: &subscriptionRoot{},
})

type subscriptionRoot struct{}

func (subscriptionRoot) MessageSent(ctx context.Context) (chan *message, error) {
    ch := make(chan *message, 1)
    // pump events; respect ctx.Done()
    return ch, nil
}
```

See [`examples/subscriptions`](../examples/subscriptions).

## Merging SDL

- First module may use `type Query { ... }`.
- Later modules may use `type Query { ... }` or `extend type Query { ... }`.
- gogql normalizes duplicate root types and validates the merged schema.

## Struct resolvers

Instead of `ResolverMap`, you can pass a single root struct with methods matching the schema (graph-gophers style). Only one such root is supported across modules unless you split concerns with `ResolverMap` for queries and a separate subscription struct as above.

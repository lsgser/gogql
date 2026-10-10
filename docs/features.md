# Features reference

Quick index of cross-cutting gogql capabilities. Works with **any** module layout ([Modules & resolvers](modules-and-resolvers.md#ways-to-organize-modules), [Project layout](project-layout.md)).

| Feature | Configure on | Doc section | Example |
|---------|----------------|-------------|---------|
| **CLI scaffold** | `gogql init`, `gogql module add` | [Project layout](project-layout.md) | Generated `src/` tree |
| **DataLoaders** | `ApplicationConfig.Loaders` | [Below](#dataloaders) | [`examples/database`](../examples/database) |
| **Query depth limit** | `ApplicationConfig.Security` | [Below](#query-depth-limit) | [`examples/subscriptions`](../examples/subscriptions), [`examples/database`](../examples/database) |
| **Dependency injection** | `ModuleConfig.Providers` | [Modules — DI](modules-and-resolvers.md#dependency-injection) | [`examples/database`](../examples/database) |
| **JWT / request context** | `ServerConfig.ContextFunc` | [Database & JWT](database-and-auth.md) | [`examples/database`](../examples/database) |
| **Playground** | `ServerConfig.Playground` | [Server & playground](server-and-playground.md#playground) | All examples |
| **Gin / custom HTTP** | `server.Handler()` + `gin.WrapH` | [Gin integration](gin.md) | — |
| **Subscriptions** | SDL + `SubscriptionResolvers` | [Modules — subscriptions](modules-and-resolvers.md#subscriptions) | [`examples/subscriptions`](../examples/subscriptions) |
| **SDL merge** | Multiple `Module`s | [Modules — merging SDL](modules-and-resolvers.md#merging-sdl) | [`examples/basic`](../examples/basic) |

---

## DataLoaders

**Purpose:** Batch and cache loads within a **single GraphQL request** (avoids N+1 queries to SQL or HTTP).

**Setup** on the application (not on individual modules):

```go
app := gogql.MustApplication(gogql.ApplicationConfig{
    Modules: modules.All(),
    Loaders: gogql.LoaderFactories{
        "user": func() *dataloader.Loader {
            return gogql.NewStringKeyLoader(batchLoadUsers)
        },
    },
})
```

- Each factory runs **once per request** when that loader is first used.
- The server attaches the registry via `Application.RequestContext` (HTTP and WebSocket).

**In a resolver:**

```go
loader, err := gogql.LoadersFromContext(ctx, "user")
if err != nil {
    return nil, err
}
user, err := loader.Load(ctx, dataloader.StringKey(id))()
```

**Batch many keys** in one resolver:

```go
results, errs := loader.LoadMany(ctx, dataloader.NewKeysFromStrings(ids))()
```

**With SQL:** see [Database & JWT — DataLoaders + database](database-and-auth.md#dataloaders--database) and `examples/database/modules/users/module.go` (`LoaderFactories` + `resolveUser`).

**API surface:** `LoaderFactories`, `NewLoaderRegistry`, `WithLoaderRegistry`, `LoaderRegistryFrom`, `LoadersFromContext`, `NewStringKeyLoader` (wraps [graph-gophers/dataloader](https://github.com/graph-gophers/dataloader)).

---

## Query depth limit

**Purpose:** Reject queries that nest fields too deeply (DoS / accidental expensive queries).

**Setup:**

```go
app := gogql.MustApplication(gogql.ApplicationConfig{
    Modules: modules,
    Security: gogql.SecurityConfig{
        MaxDepth: 15, // 0 = no limit (default)
    },
})
```

Uses [graph-gophers/graphql-go](https://github.com/graph-gophers/graphql-go) `MaxDepth` during validation. Over-limit queries return GraphQL errors in the response (no partial data for that operation).

**Related limits** (optional, via `ApplicationConfig.SchemaOpts`):

- `graphql.MaxQueryLength(n)`
- `graphql.MaxParallelism(n)`

See [Server & playground — Security](server-and-playground.md#security-query-depth).

---

## Other security notes

- **JWT** is not built into gogql; use `ServerConfig.ContextFunc` and `gogql.WithAuthClaims` ([Database & JWT](database-and-auth.md)).
- **CORS** on the default server is permissive; tighten in production with a reverse proxy or wrapper handler.
- **Introspection** can be restricted with graph-gophers `graphql.RestrictIntrospection` in `SchemaOpts` if needed.

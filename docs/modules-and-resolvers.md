# Modules & resolvers

## Module shape

```go
gogql.MustModule(gogql.ModuleConfig{
    ID:        "users",           // unique per application
    TypeDefs:  `... GraphQL SDL ...`,
    Resolvers: /* see below */,
    SubscriptionResolvers: /* optional struct with methods */,
    Providers: []gogql.Provider{
        gogql.Provide(&UsersService{}),
    },
})
```

| Field | Purpose |
|-------|---------|
| `ID` | Unique module name |
| `TypeDefs` | GraphQL SDL for this module |
| `Resolvers` | `*ResolverMap` or a graph-gophers-compatible root struct |
| `SubscriptionResolvers` | Struct with **methods** for subscription fields (required for WebSocket subscriptions) |
| `Providers` | Services on the request-scoped injector |

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

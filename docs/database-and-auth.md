# Database & JWT authentication

gogql does not ship an ORM or auth provider. You wire **standard Go dependencies** (e.g. `database/sql`, your JWT library) through **module providers** and **`ServerConfig.ContextFunc`**.

Full working code: [`examples/database`](../examples/database) (`modules/`, `store/`, `auth/` + `gogql.MustApplication`).

## Database access pattern

1. Open `*sql.DB` (or your DB client) in `main`.
2. Register repositories on the module injector with `gogql.Provide`.
3. Resolve them in resolvers with `gogql.MustGet`.

```go
usersModule := gogql.MustModule(gogql.ModuleConfig{
    ID: "users",
    TypeDefs: `...`,
    Providers: []gogql.Provider{
        gogql.Provide(NewUserRepository(db)),
    },
    Resolvers: gogql.NewResolverMap().Query("user", func(ctx context.Context, args struct {
        ID graphql.ID
    }) (*User, error) {
        repo := gogql.MustGet[*UserRepository](ctx)
        rec, err := repo.ByID(ctx, string(args.ID))
        if err != nil {
            return nil, err
        }
        return toGraphQLUser(rec), nil
    }),
})
```

The server calls `Application.RequestContext` on every request, which attaches the **application injector** so `MustGet` works inside resolvers.

### DataLoaders + database

For N+1 avoidance, batch IDs in a DataLoader and run **one SQL query** per request batch:

```go
Loaders: gogql.LoaderFactories{
    "user": func() *dataloader.Loader {
        return gogql.NewStringKeyLoader(func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
            // SELECT ... WHERE id IN (...)
            return results
        })
    },
},
```

See [Modules & resolvers — DataLoaders](modules-and-resolvers.md#dataloaders).

### Connection configuration

Typical production setup:

```go
db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
db.SetMaxOpenConns(25)
db.SetMaxIdleConns(5)
db.SetConnMaxLifetime(30 * time.Minute)
```

Use migrations (goose, golang-migrate, etc.) outside gogql; keep resolvers thin and logic in repositories.

## JWT configuration

gogql exposes **context helpers** for claims after you validate a token:

```go
gogql.WithAuthClaims(ctx, gogql.AuthClaims{Subject: userID})
gogql.AuthClaimsFrom(ctx)   // optional
gogql.MustAuthClaims(ctx)   // panics if missing
```

Validate JWT on each HTTP/WebSocket request with **`ServerConfig.ContextFunc`** (runs before the injector is attached):

```go
server := gogql.NewServer(app, gogql.ServerConfig{
    ContextFunc: func(ctx context.Context, r *http.Request) context.Context {
        token := parseBearer(r.Header.Get("Authorization"))
        if token == "" {
            return ctx
        }
        sub, err := validateJWT(token, []byte(os.Getenv("JWT_SECRET")))
        if err != nil {
            return ctx // or attach error / reject in middleware
        }
        return gogql.WithAuthClaims(ctx, gogql.AuthClaims{Subject: sub})
    },
})
```

### Recommended libraries

- [github.com/golang-jwt/jwt/v5](https://github.com/golang-jwt/jwt) — parse and verify HS256/RS256 tokens (used in the example).
- For OAuth2/OIDC, validate access tokens from your IdP and map `sub` to your user id.

### Protecting fields

- **Public query**: omit checks; `AuthClaimsFrom` returns `ok == false`.
- **Private query** (e.g. `me`): require claims and return a GraphQL error:

```go
func resolveMe(ctx context.Context) (*User, error) {
    claims, ok := gogql.AuthClaimsFrom(ctx)
    if !ok {
        return nil, fmt.Errorf("unauthenticated")
    }
    return loadUser(ctx, claims.Subject)
}
```

For global auth, use `ContextFunc` plus schema design (separate public/private modules) or add HTTP middleware in front of `server.Handler()`.

### WebSocket subscriptions

The same `ContextFunc` is applied for **graphql-transport-ws** connections via the server’s context generator, so JWT in the **HTTP upgrade request** `Authorization` header is available in subscription resolvers.

Send the header when opening the WebSocket from clients that support it.

### Playground / GraphiQL

Add headers in the playground UI (GraphiQL: HTTP Headers):

```json
{
  "Authorization": "Bearer <your-jwt>"
}
```

The example prints a demo token on startup when you run:

```bash
go run ./examples/database
```

## Environment variables (suggested)

| Variable | Purpose |
|----------|---------|
| `DATABASE_URL` | Postgres/SQLite DSN |
| `JWT_SECRET` | HMAC secret for HS256 (example uses a constant; do not do that in production) |
| `HTTP_ADDR` | Listen address (e.g. `:8080`) |

## Related

- [Getting started](getting-started.md)
- [Modules & resolvers](modules-and-resolvers.md)
- [Server & playground](server-and-playground.md)

# Using gogql with Gin

[gogql](https://github.com/lsgser/gogql) ships a standard library `http.Handler` from `NewServer`. **[Gin](https://github.com/gin-gonic/gin)** can mount that handler next to REST routes with `gin.WrapH`, without changing how you define modules or the application.

You still build the GraphQL stack the same way ([Getting started](getting-started.md)):

```go
app := gogql.MustApplication(gogql.ApplicationConfig{
    Modules: modules,
})

gqlServer := gogql.NewServer(app, gogql.ServerConfig{
    GraphQLPath: "/graphql",
    HealthPath:  "/health",
    Playground: &gogql.PlaygroundConfig{
        Enabled: true,
        Path:    "/playground",
        UI:      gogql.PlaygroundGraphiQL,
    },
    ContextFunc: myContextFunc, // optional — JWT, tracing, etc.
})
```

The difference is only **how you listen**: Gin owns the HTTP server instead of `server.ListenAndServe`.

## Install Gin

In your app module (alongside gogql):

```bash
go get github.com/gin-gonic/gin@latest
```

## Mount gogql on Gin routes

`gqlServer.Handler()` is an `http.ServeMux` that registers `/graphql` (HTTP + WebSocket subscriptions), `/playground`, and `/health`. Register **the same paths on Gin** and delegate with `gin.WrapH`:

```go
package main

import (
    "log"

    "github.com/gin-gonic/gin"
    "github.com/lsgser/gogql"
)

func main() {
    app := gogql.MustApplication(gogql.ApplicationConfig{
        Modules: []*gogql.Module{/* ... */},
    })

    gqlServer := gogql.NewServer(app, gogql.ServerConfig{
        Playground: &gogql.PlaygroundConfig{Enabled: true, Path: "/playground"},
    })

    h := gqlServer.Handler()

    r := gin.New()
    r.Use(gin.Logger(), gin.Recovery())

    // GraphQL: POST, GET, OPTIONS, and graphql-transport-ws upgrade on this path
    r.Any("/graphql", gin.WrapH(h))
    r.GET("/playground", gin.WrapH(h))
    r.GET("/health", gin.WrapH(h))

    // REST (or other) routes on the same engine
    r.GET("/api/version", func(c *gin.Context) {
        c.JSON(200, gin.H{"service": "my-api"})
    })

    if err := r.Run(":8080"); err != nil {
        log.Fatal(err)
    }
}
```

**Why `Any` on `/graphql`?** gogql handles `POST` and `GET` for queries and `OPTIONS` for CORS preflight. Subscriptions use a WebSocket upgrade on the same path, which also arrives as a request Gin must pass through—`Any` covers that.

**Path names must match.** If you set `ServerConfig.GraphQLPath` to `/api/graphql`, register `r.Any("/api/graphql", gin.WrapH(h))` and point the playground at that path (`PlaygroundConfig` does not change the GraphQL URL; GraphiQL uses the path you configured when creating the server).

## Playground-only or API-only

- **GraphQL only** (no UI): `Playground: &gogql.PlaygroundConfig{Enabled: false}` and omit the `/playground` Gin route.
- **Custom health check**: expose liveness with a Gin route (e.g. `r.GET("/health", ...)`) and skip `r.GET("/health", gin.WrapH(h))` if you do not need gogql’s plain `ok` response.

If you change `GraphQLPath`, `Playground.Path`, or `HealthPath` in `ServerConfig`, register the **same** paths on Gin.

## Authentication

### Option A — `ServerConfig.ContextFunc` (recommended for parity with standalone gogql)

Works unchanged with `gin.WrapH`: gogql reads `*http.Request` and builds context before resolvers run. Same pattern as [Database & JWT](database-and-auth.md):

```go
gqlServer := gogql.NewServer(app, gogql.ServerConfig{
    ContextFunc: auth.ContextFunc(jwtSecret),
})
```

WebSocket subscriptions use the same `ContextFunc` on the upgrade request, so send `Authorization` on the initial WebSocket handshake.

### Option B — Gin middleware

Run JWT or session middleware **before** the GraphQL handler and attach values to `c.Request.Context()` (for example with `gogql.WithAuthClaims`). Leave `ContextFunc` nil if middleware already set everything resolvers need:

```go
r.Use(func(c *gin.Context) {
    token := c.GetHeader("Authorization")
    if claims, err := parseBearer(token); err == nil {
        ctx := gogql.WithAuthClaims(c.Request.Context(), claims)
        c.Request = c.Request.WithContext(ctx)
    }
    c.Next()
})

r.Any("/graphql", gin.WrapH(gqlServer.Handler()))
```

Use **either** global Gin middleware **or** `ContextFunc` for the same concern, not both, unless you intend layered checks.

## CORS and production middleware

The default gogql HTTP handler sets permissive CORS headers on GraphQL responses ([Server & playground](server-and-playground.md)). If Gin also applies CORS middleware, avoid conflicting duplicate headers; often you disable one layer for `/graphql` or configure Gin to skip GraphQL routes.

For production, keep Gin’s `Recovery()`, logging, rate limits, or TLS termination in front of `WrapH` as you would for any other `http.Handler`.

## Subscriptions behind Gin

Subscriptions stay on **`GraphQLPath`** with subprotocol **`graphql-transport-ws`**. Clients connect to `ws://host:port/graphql` (same path as HTTP). No extra Gin route is required beyond `r.Any("/graphql", ...)`.

If you terminate TLS at a reverse proxy, ensure WebSocket upgrades are forwarded to Gin on that path.

## Standalone server vs Gin

| Approach | When to use |
|----------|-------------|
| `server.ListenAndServe(":8080")` | GraphQL-only service, examples, quick demos |
| Gin + `gin.WrapH(gqlServer.Handler())` | GraphQL plus REST, existing Gin codebase, shared middleware |

You do not need two listeners: one Gin engine can serve both.

## Troubleshooting

| Issue | Check |
|-------|--------|
| 404 on `/graphql` | Gin route path matches `ServerConfig.GraphQLPath` |
| Playground loads but queries fail | Playground targets the configured GraphQL path; CORS/network same origin |
| WebSocket subscription fails | `EnableSubscriptions` not set to `false`; use `Any` on GraphQL path; proxy allows upgrade |
| Auth works on REST but not GraphQL | Gin auth on GraphQL route or `ContextFunc`; WebSocket needs `Authorization` on upgrade |
| Double CORS errors | Gin CORS + gogql CORS both setting `Access-Control-Allow-Origin` |

## Related

- [Server & playground](server-and-playground.md) — HTTP methods, playground, depth limits, WebSocket
- [Database & JWT](database-and-auth.md) — `ContextFunc`, DataLoaders, SQL
- [Features](features.md) — feature index

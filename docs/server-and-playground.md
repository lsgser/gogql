# Server & playground

## Creating the server

```go
server := gogql.NewServer(app, gogql.ServerConfig{
    GraphQLPath: "/graphql",   // default
    HealthPath:  "/health",    // default, returns "ok"
    Playground: &gogql.PlaygroundConfig{
        Enabled: true,
        Path:    "/playground",
        UI:      gogql.PlaygroundGraphiQL, // or gogql.PlaygroundApolloSandbox
    },
})
http.ListenAndServe(":8080", server.Handler())
// or
server.ListenAndServe(":8080")
```

## HTTP API

| Method | Path | Behavior |
|--------|------|----------|
| `POST` | `GraphQLPath` | JSON body: `{ "query", "variables", "operationName" }` |
| `GET` | `GraphQLPath` | Query string: `?query=...&variables=...` |
| `GET` | Playground path | GraphiQL or Apollo Sandbox UI |
| `GET` | `HealthPath` | Liveness text `ok` |
| `OPTIONS` | `GraphQLPath` | CORS preflight |

CORS is permissive (`Access-Control-Allow-Origin: *`) for local development. Put a reverse proxy or wrapper in front for production hardening.

## Playground

| `PlaygroundUI` | Description |
|----------------|-------------|
| `PlaygroundGraphiQL` | Embedded GraphiQL (default) |
| `PlaygroundApolloSandbox` | Embedded [Apollo Sandbox](https://www.apollographql.com/docs/apollo-sandbox) |

Disable the playground:

```go
Playground: &gogql.PlaygroundConfig{Enabled: false},
```

When `Playground` is `nil`, gogql still enables GraphiQL at `/playground` by default. Set `Enabled: false` explicitly to turn it off.

## Subscriptions (WebSocket)

Enabled by default on the same path as HTTP (`GraphQLPath`).

- Subprotocol: **`graphql-transport-ws`** ([spec](https://github.com/enisdenjo/graphql-ws/blob/master/PROTOCOL.md))
- URL: `ws://localhost:8080/graphql` (or `wss://` behind TLS)

Disable:

```go
enableSubs := false
gogql.ServerConfig{EnableSubscriptions: &enableSubs}
```

Programmatic subscribe (tests, custom gateways):

```go
ch, err := app.Subscribe(ctx, query, operationName, variables)
```

## Security: query depth

See also: [Features — query depth limit](features.md#query-depth-limit).

Limit nesting depth to reduce abusive queries:

```go
gogql.MustApplication(gogql.ApplicationConfig{
    Modules:  modules,
    Security: gogql.SecurityConfig{MaxDepth: 15},
})
```

`0` means no limit (graph-gophers default).

Additional execution options can be passed via `ApplicationConfig.SchemaOpts` (e.g. tracing plugins from graph-gophers).

## Programmatic execution

Without HTTP:

```go
resp := app.Execute(ctx, `{ user(id: "1") { name } }`, "", nil)
```

`ctx` is enriched with injector and loaders when using `Execute` or the HTTP server.

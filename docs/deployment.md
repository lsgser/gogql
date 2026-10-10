# Deployment & production

**gogql** is a library: you deploy **your Go application** that imports `github.com/lsgser/gogql`, builds to a binary (or container), and runs behind HTTPS. This guide covers common production choices; adapt them to your cloud or on‑prem setup.

## Checklist

| Item | Recommendation |
|------|----------------|
| Playground | **Disable** in production (`Playground: &gogql.PlaygroundConfig{Enabled: false}`) |
| TLS | Terminate at **reverse proxy** or load balancer (`https` / `wss://`) |
| CORS | Do not rely on gogql’s dev defaults; set policy at proxy or in your HTTP wrapper |
| Query abuse | Set `SecurityConfig.MaxDepth`; consider `SchemaOpts` (query length, introspection) — [Features](features.md) |
| Secrets | `JWT_SECRET`, `DATABASE_URL`, etc. from env or a secret manager — [Database & JWT](database-and-auth.md) |
| Health | Use gogql’s `GET /health` (or your own) for load balancer probes |
| Subscriptions | Proxy must allow **WebSocket upgrade** on `GraphQLPath` — [Server & playground](server-and-playground.md#subscriptions-websocket) |

## Build and run

From your app module (not the gogql repo itself, unless you ship an example):

```bash
go build -o api .             # gogql init: main.go at repo root, logic in src/
./api
```

Scaffolded apps use **`src/config`** for `HTTP_ADDR` and server settings ([Project layout](project-layout.md)).

Cross-compile for Linux servers:

```bash
GOOS=linux GOARCH=amd64 go build -o api .
```

Listen address is **your code** (examples use `:8080`). Prefer binding to an internal port and exposing only through a proxy:

```go
addr := os.Getenv("HTTP_ADDR")
if addr == "" {
    addr = ":8080"
}
log.Fatal(server.ListenAndServe(addr))
```

Or use `http.ListenAndServe(addr, server.Handler())` if you add custom middleware around `server.Handler()`.

## Production server configuration

In a **`gogql init`** app, toggle playground and wire auth in **`src/config/config.go`** (`Playground.Enabled: false`, `ContextFunc: utils.AuthContextFunc`). The same settings apply if you call `NewServer` manually:

```go
app := gogql.MustApplication(gogql.ApplicationConfig{
    Modules:  modules,
    Security: gogql.SecurityConfig{MaxDepth: 15},
})

server := gogql.NewServer(app, gogql.ServerConfig{
    GraphQLPath: "/graphql",
    HealthPath:  "/health",
    Playground:  &gogql.PlaygroundConfig{Enabled: false},
    ContextFunc: authContextFromEnv(), // JWT — see database-and-auth.md
})
```

Using **[Gin](gin.md)** or another router: mount `server.Handler()` with `gin.WrapH` and apply production middleware there.

## Environment variables (typical)

| Variable | Purpose |
|----------|---------|
| `HTTP_ADDR` | Listen address (e.g. `:8080`, `127.0.0.1:8080`) |
| `DATABASE_URL` | SQL DSN when using a database |
| `JWT_SECRET` | Signing/validation secret for HS256 (or use your IdP) |
| `GOMAXPROCS` | Optional; usually leave default on container vCPU limits |

See also [Database & JWT — environment variables](database-and-auth.md#environment-variables-suggested).

## TLS and reverse proxy

gogql’s built-in server serves **plain HTTP**. In production, put **nginx**, **Caddy**, **Traefik**, or a cloud load balancer in front:

- Terminate TLS (`443` → your app on `8080`)
- Set `X-Forwarded-For` / `X-Forwarded-Proto` if your app needs them
- For **subscriptions**, enable WebSocket pass-through to the same path as GraphQL (e.g. `/graphql`)

Clients should use:

- HTTP: `https://api.example.com/graphql`
- WebSocket: `wss://api.example.com/graphql` with subprotocol **`graphql-transport-ws`**

Send `Authorization` on the WebSocket **upgrade** request when using JWT ([Database & JWT](database-and-auth.md#websocket-subscriptions)).

### Caddy (minimal sketch)

```text
api.example.com {
    reverse_proxy localhost:8080
}
```

Caddy handles TLS and WebSocket upgrades by default for many setups; verify with your subscription clients.

### nginx (WebSocket)

Ensure `proxy_http_version 1.1`, `Upgrade`, and `Connection` headers are set for the GraphQL location. Point `proxy_pass` to your Go process.

## Docker

Multi-stage build keeps images small:

```dockerfile
# syntax=docker/dockerfile:1
FROM golang:1.22 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /api .

FROM gcr.io/distroless/static-debian12
COPY --from=build /api /api
ENV HTTP_ADDR=:8080
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/api"]
```

Run with secrets and env:

```bash
docker run --rm -p 8080:8080 \
  -e DATABASE_URL=... \
  -e JWT_SECRET=... \
  your-image
```

Put the container behind a platform load balancer or ingress with TLS.

## Kubernetes (sketch)

- **Deployment** — one container running your binary; `containerPort: 8080`
- **Service** — ClusterIP targeting that port
- **Ingress** — TLS cert; enable WebSocket annotations if your ingress controller requires them
- **Liveness / readiness** — HTTP GET `http://pod:8080/health` (gogql default)

Scale horizontally when resolvers are stateless; use sticky sessions only if you add subscription infrastructure that requires it (each pod can handle its own WebSocket connections).

## Process managers

On a VM, use **systemd** with `Restart=on-failure`, `EnvironmentFile=/etc/your-api/env`, and bind to localhost if nginx sits in front.

Platforms (**Fly.io**, **Railway**, **Render**, **Cloud Run**, etc.): set the start command to your binary, map `PORT` or `HTTP_ADDR` to the platform’s expected port, attach managed TLS at the edge, and wire health checks to `/health`.

## Observability

- Use `ApplicationConfig.SchemaOpts` with graph-gophers tracing plugins where needed
- Log at the HTTP layer (middleware or Gin) — request ID, latency, status
- gogql does not ship metrics; export them from your wrapper or proxy

## What not to deploy

- Do not expose the **playground** on public URLs in production
- Do not commit **JWT secrets** or demo secrets from examples
- Do not run SQLite file DB on ephemeral disks without volumes unless it is intentional

## Related

- [Installation](installation.md) — `go get`, project layout
- [Server & playground](server-and-playground.md) — routes, CORS, WebSocket
- [Gin integration](gin.md) — single process with REST
- [Features](features.md) — depth limits, introspection, CORS notes

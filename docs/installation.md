# Installation

gogql is published as a Go module from [github.com/lsgser/gogql](https://github.com/lsgser/gogql).

## Requirements

- **Go 1.22+** (the library module targets Go 1.22; check `go.mod` for the exact version used in development)
- A Go module in your project (`go mod init` if you are starting fresh)

## Install the library

In your project directory:

```bash
go get github.com/lsgser/gogql@latest
```

Import in code:

```go
import "github.com/lsgser/gogql"
```

Verify the module is listed in `go.mod`:

```bash
go mod tidy
```

### Pin a version

After [releases](https://github.com/lsgser/gogql/releases) are tagged on GitHub:

```bash
go get github.com/lsgser/gogql@v0.1.0
```

Until the first tag exists, use `@latest` from the default branch.

### Local development (replace)

When working on gogql alongside your app:

```go
// go.mod
require github.com/lsgser/gogql v0.0.0

replace github.com/lsgser/gogql => ../gogql
```

Run `go mod tidy` after editing `go.mod`.

### Where the code lives after `go get`

`go get` does **not** copy gogql into your project tree. The library is stored in your module cache (see `go env GOMODCACHE`), typically:

```text
$GOMODCACHE/github.com/lsgser/gogql@v<version>/
```

Your app only records the dependency in **`go.mod`** / **`go.sum`** and imports `github.com/lsgser/gogql`. You organize application code yourself, or use **`gogql init`** for the recommended layout (below).

## Clone the gogql repository

If you download the full source from GitHub, the repository layout looks like this:

```text
gogql/
├── LICENSE
├── README.md
├── go.mod
├── go.sum
├── doc.go                  # Package docs (import github.com/lsgser/gogql)
├── export.go               # Public API re-exports
│
├── cmd/
│   └── gogql/
│       └── main.go         # CLI: init, version
│
├── docs/                   # Documentation (this site)
├── examples/
│   ├── README.md
│   ├── basic/              # Modular SDL demo
│   ├── subscriptions/
│   └── database/
│
└── internal/
    ├── core/               # Library implementation (package core)
    │   ├── application.go  # MustApplication, Execute, Subscribe
    │   ├── module.go       # MustModule, ModuleConfig, typedef loading
    │   ├── server.go       # NewServer, HTTP, WebSocket
    │   ├── resolver.go     # ResolverMap
    │   ├── injector.go     # DI + MustGet
    │   ├── loader.go       # DataLoader registry
    │   ├── auth.go         # AuthClaims on context
    │   ├── config.go       # Playground, Security (MaxDepth)
    │   └── testdata/       # Tests for SDL loading
    ├── merge/              # SDL merge across modules
    └── scaffold/           # gogql init templates
```

Applications still use a single import: `import "github.com/lsgser/gogql"`. You do not import `internal/core` from outside this module.

Clone and install the CLI:

```bash
git clone https://github.com/lsgser/gogql.git
cd gogql
```

Examples in the clone run with:

```bash
go run ./examples/basic
```

## Install the CLI (optional)

The **`gogql`** command scaffolds a new server project.

**From a clone of this repository:**

```bash
git clone https://github.com/lsgser/gogql.git
cd gogql
go install ./cmd/gogql
```

Ensure `$GOPATH/bin` or `$HOME/go/bin` is on your `PATH`. Then:

```bash
gogql version
gogql init my-api
```

**Without installing globally** (from the repo root):

```bash
go run ./cmd/gogql init my-api
```

### `gogql init` flags

| Flag | Default | Meaning |
|------|---------|---------|
| `-module` | directory name | Go module path for the new project (e.g. `github.com/you/my-api`) |
| `-gogql` | `../gogql` | `replace` path when developing gogql locally; omit from generated `go.mod` when using only `go get` |

After `init`:

```bash
cd my-api
go mod tidy
go run .
```

### Project folder structure after `gogql init`

Running `gogql init my-api` creates a **standalone GraphQL server** that uses the gogql module pattern:

```text
my-api/
├── go.mod
├── main.go
└── modules/
    ├── registry.go
    └── users/
        ├── module.go           # wires typedefs + resolvers into gogql.MustModule
        ├── typedefs.go         # embed / JoinTypeDefs / LoadTypeDefsFS
        ├── resolvers.go        # ResolverMap and resolver funcs
        └── schema/
            ├── user.graphql
            └── query.graphql
```

This matches the **split-file (Approach B)** layout from [`examples/basic/modules/users`](../examples/basic/modules/users). It is optional: a single `modules/<name>/module.go` with inline `TypeDefs` and `Resolvers` remains valid ([Approach A](modules-and-resolvers.md#approach-a--inline-module-backward-compatible)).

Add features by creating `modules/<name>/` and appending to `modules.All()`. Mix inline and split modules in the same project.

With **`go get` only** (no `init`), use any layout; compose **`gogql.MustModule`** → **`gogql.MustApplication`** → **`gogql.NewServer`** in `main`.

## Transitive dependencies

gogql builds on:

- [graph-gophers/graphql-go](https://github.com/graph-gophers/graphql-go) — schema execution
- [vektah/gqlparser](https://github.com/vektah/gqlparser) — SDL parsing and merge validation
- [graph-gophers/graphql-transport-ws](https://github.com/graph-gophers/graphql-transport-ws) — WebSocket subscriptions
- [graph-gophers/dataloader](https://github.com/graph-gophers/dataloader) — batched loaders

You do not need to add these manually unless you use them directly (e.g. `dataloader` in resolvers).

## Troubleshooting

| Issue | What to try |
|-------|-------------|
| `module github.com/lsgser/gogql: not found` | Confirm the repo is pushed to GitHub and your `GOPROXY` can reach it (`go env GOPROXY`). |
| Empty module cache | Run `go get github.com/lsgser/gogql@latest` again after the first push to the default branch. |
| Playground cannot reach API | Use the same host/port; playground calls `window.location.origin + "/graphql"` by default. |

Next: [Getting started](getting-started.md).

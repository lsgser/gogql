# Installation

gogql is published as a Go module from [github.com/lsgser/gogql](https://github.com/lsgser/gogql).

## Requirements

- **Go 1.22+** (the library module targets Go 1.22; check `go.mod` for the exact version used in development)
- A Go module in your project (`go mod init` if you are starting fresh)

## Install the library

`go get` only works **inside a Go module** (a folder that contains `go.mod`). It does not work from your home directory, from a random clone path without `go mod init`, or from the gogql repo root unless that repo is your app module.

**New app:**

```bash
mkdir my-api && cd my-api
go mod init github.com/you/my-api   # any module path you own
go get github.com/lsgser/gogql@latest
```

**Existing app:** `cd` to the directory that already has `go.mod`, then:

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
    ├── cli/                # Command list + version (used by cmd/gogql)
    ├── core/               # Library implementation (package core)
    ├── merge/              # SDL merge across modules
    └── scaffold/           # init + module generators (embedded templates)
```

Each `.go` file under `internal/core` and generated app files include **file headers** describing that file’s role. See [Project layout](project-layout.md) for the **application** tree created by **`gogql init`**, not the library repo itself.

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

The **`gogql`** command scaffolds new server projects. Like [gofreight](https://github.com/lsgser/gofreight), the CLI is a **separate command** from the library: `go get github.com/lsgser/gogql@latest` adds the **library** to your app’s `go.mod` only; it does **not** put `gogql` on your `PATH`.

Pick one way to run the CLI:

### A — Global install (same as gofreight’s quick start)

```bash
go install github.com/lsgser/gogql/cmd/gogql@latest
export PATH="$PATH:$(go env GOPATH)/bin"   # once per machine
gogql version
gogql init my-api
```

### B — Per-project tool (Go 1.24+, no global install)

Inside your app module (after `go mod init`):

```bash
go get -tool github.com/lsgser/gogql/cmd/gogql@latest
go tool gogql version
go tool gogql init my-api
```

This adds a `tool github.com/lsgser/gogql/cmd/gogql` line to `go.mod`. Upgrade with `go get -tool github.com/lsgser/gogql/cmd/gogql@latest` or `go get tool`.

### C — From a clone of this repository

```bash
git clone https://github.com/lsgser/gogql.git
cd gogql
go install ./cmd/gogql
# or without installing:
go run ./cmd/gogql init my-api
```

Projects created with **`gogql init`** already include the `tool` line in `go.mod`, so after `cd my-api` and `go mod tidy` you can use **`go tool gogql`** for generators without a global install.

### Commands

Run **`gogql version`** (or **`gogql help`**) to print the CLI version and the full command list.

| Command | Description |
|---------|-------------|
| `gogql init [dir]` | New app with default **`src/`** layout (common, modules, schema, config, utils) |
| `gogql module add <name>` | Domain module: `*.graphql`, `*.resolvers.go`, `*.model.go`, `*.service.go`, `module.go`; updates `src/schema/modules.go` |
| `gogql module typedefs <name>` | Add `<name>.graphql` (creates `module.go` if missing) |
| `gogql module resolvers <name>` | Add resolver, model, and service stubs; refreshes schema registry |
| `gogql module schema <name>` | Same as `module typedefs` — SDL file only |
| `gogql version` | Version and command summary |

Module commands must run from your **app root** (where `go.mod` lives). Use **`-force`** to overwrite generated files. Optional **`-C /path/to/app`** sets the project directory.

Example after `gogql init my-api`:

```bash
cd my-api
go tool gogql module add product
go run .
```

### `gogql init` flags

| Flag | Default | Meaning |
|------|---------|---------|
| `-module` | directory name | Go module path for the new project (e.g. `github.com/you/my-api`) |
| `-gogql` | `../gogql` | `replace` path when developing gogql locally; omit from generated `go.mod` when using only `go get` |

After `init` (global `gogql` or `go tool gogql`):

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
├── main.go                      # calls src/app.Run()
└── src/
    ├── app/
    │   └── app.go               # index — MustApplication + NewServer
    ├── common/                  # shared SDL (scalars, base types)
    │   ├── common.go
    │   └── scalars.graphql
    ├── modules/                 # one folder per domain
    │   └── user/
    │       ├── module.go        # gogql.MustModule wiring
    │       ├── user.graphql
    │       ├── user.resolvers.go
    │       ├── user.model.go    # data layer types
    │       └── user.service.go  # business logic
    ├── schema/
    │   └── modules.go           # merges All() + LoaderFactories()
    ├── config/
    │   └── config.go
    └── utils/
        └── auth.go              # optional ContextFunc helpers
```

This is the **default** scaffold only. gogql does not require `src/` — you can use flat `modules/`, inline SDL in one file, or any tree; wire **`gogql.MustModule`** → **`gogql.MustApplication`** → **`gogql.NewServer`** yourself ([Modules & resolvers](modules-and-resolvers.md#ways-to-organize-modules)). The CLI detects **`src/modules`** (new) or legacy **`modules/`** when running **`gogql module add`**.

Add domains with **`gogql module add product`** (regenerates `src/schema/modules.go`) or copy the `user/` folder by hand.

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
| `go.mod file not found` / `go get is no longer supported outside a module` | Run `go mod init <your/module/path>` in your **app** directory first, then `go get` again. To install the **CLI** without an app module, use `go install github.com/lsgser/gogql/cmd/gogql@latest`. |
| `module github.com/lsgser/gogql: not found` | Confirm the repo is pushed to GitHub and your `GOPROXY` can reach it (`go env GOPROXY`). |
| Empty module cache | Run `go get github.com/lsgser/gogql@latest` again after the first push to the default branch. |
| Playground cannot reach API | Use the same host/port; playground calls `window.location.origin + "/graphql"` by default. |

Next: [Project layout](project-layout.md) · [Getting started](getting-started.md).

# Project layout

gogql does **not** force a folder structure. You can put SDL and resolvers anywhere and call `MustModule` / `MustApplication` by hand ([Getting started — manual](getting-started.md#manual-minimal-server)). The **`gogql init`** scaffold and **`gogql module add`** generators use a **default `src/` layout** similar to domain-driven Node/GraphQL Modules projects.

## Default scaffold (`gogql init`)

```text
my-api/
├── go.mod
├── main.go                 # calls src/app.Run()
└── src/
    ├── app/
    │   └── app.go          # MustApplication + NewServer (app entry)
    ├── common/             # shared SDL (scalars, shared types)
    │   ├── common.go
    │   └── scalars.graphql
    ├── modules/            # one folder per domain / feature
    │   └── user/
    │       ├── module.go       # gogql.MustModule wiring
    │       ├── user.graphql    # schema / type definitions
    │       ├── user.resolvers.go
    │       ├── user.model.go   # data-layer types
    │       └── user.service.go # business logic
    ├── schema/
    │   └── modules.go      # All() + LoaderFactories() registry
    ├── config/
    │   └── config.go       # HTTP addr, ServerConfig, MaxDepth
    └── utils/
        └── auth.go         # optional ContextFunc helpers
```

| Path | Role |
|------|------|
| **`main.go`** | Thin process entry; keep deployment concerns here or in `app`. |
| **`src/app`** | Wires `schema.All()`, loaders, security, and `NewServer` — like `src/index.ts`. |
| **`src/common`** | Optional gogql module for scalars and types every domain shares. |
| **`src/modules/<domain>`** | One gogql module per feature: SDL + resolvers + model + service. |
| **`src/schema`** | Single registry imported by `app`; **`gogql module add`** regenerates `modules.go`. |
| **`src/config`** | Environment and server settings (playground path, depth limit, auth hook). |
| **`src/utils`** | Cross-cutting helpers (JWT, logging); not required by gogql. |

### Domain module file naming

For a domain named `product`, the CLI creates:

```text
src/modules/product/
├── module.go
├── product.graphql
├── product.resolvers.go
├── product.model.go
└── product.service.go
```

Resolvers call **`DefaultService()`** (stub); replace with repositories and **`gogql.Provide`** / **`MustGet`** as the app grows ([Modules — DI](modules-and-resolvers.md#dependency-injection)).

## CLI generators

Run from the **app root** (directory with `go.mod`):

| Command | Effect |
|---------|--------|
| `gogql module add <name>` | Full domain folder + refresh `src/schema/modules.go` |
| `gogql module typedefs <name>` | `<name>.graphql` (+ `module.go` if missing) |
| `gogql module resolvers <name>` | `*.resolvers.go`, model, service stubs |
| `gogql module schema <name>` | SDL file only |

Use **`-force`** to overwrite generated files. **`-C /path/to/app`** sets the project root.

Install CLI: [Installation — CLI](installation.md#install-the-cli-optional). List commands: **`gogql version`**.

## Legacy flat layout

Older scaffolds and some docs use:

```text
modules/
├── registry.go
└── users/
    ├── module.go
    ├── typedefs.go
    ├── resolvers.go
    └── schema/*.graphql
```

The CLI still supports this if **`src/modules`** does not exist: it writes **`modules/registry.go`** instead of **`src/schema/modules.go`**. New projects should prefer **`gogql init`** (src layout).

## Repository examples (not the init template)

Examples under [`examples/`](../examples/) demonstrate library features with **simpler trees** (no `src/` wrapper):

| Example | Layout highlight |
|---------|------------------|
| [`basic`](../examples/basic) | Mixed **inline** (`greeting`) and **split** (`users` + `schema/`) under `modules/` |
| [`database`](../examples/database) | Inline module + SQL + JWT + DataLoaders |
| [`subscriptions`](../examples/subscriptions) | `SubscriptionResolvers` + WebSocket |

Use examples to learn APIs; use **`gogql init`** for a production-shaped starting tree.

## gogql library repo layout

When you clone [github.com/lsgser/gogql](https://github.com/lsgser/gogql), library code lives under **`internal/core`**; the root **`export.go`** re-exports the public API. See [Installation — clone layout](installation.md#clone-the-gogql-repository).

## Related

- [Installation](installation.md) — `go get`, init flags, troubleshooting
- [Modules & resolvers](modules-and-resolvers.md) — inline vs split SDL, DI, loaders
- [Getting started](getting-started.md) — first server (scaffold or manual)

/*
|--------------------------------------------------------------------------
| Application
|--------------------------------------------------------------------------
|
| Composes gogql modules into one executable schema (MustApplication),
| request context (injector, DataLoaders), and Execute/Subscribe helpers.
|
*/

package core

import (
	"context"
	"fmt"

	"github.com/graph-gophers/graphql-go"
	"github.com/lsgser/gogql/internal/merge"
)

// Application composes modules into one executable GraphQL schema.
type Application struct {
	modules         []*Module
	schema          *graphql.Schema
	schemaSDL       string
	injector        *Injector
	loaderFactories LoaderFactories
}

// ApplicationConfig configures a GraphQL application.
type ApplicationConfig struct {
	Modules []*Module
	// SchemaOpts are passed to graphql.ParseSchema (tracing, etc.).
	SchemaOpts []graphql.SchemaOpt
	// Security configures limits such as query depth.
	Security SecurityConfig
	// Loaders registers DataLoader factories (one loader instance per request).
	Loaders LoaderFactories
}

// NewApplication merges module SDL and builds an executable schema.
func NewApplication(cfg ApplicationConfig) (*Application, error) {
	if len(cfg.Modules) == 0 {
		return nil, fmt.Errorf("gogql: at least one module is required")
	}

	seen := make(map[string]struct{}, len(cfg.Modules))
	sources := make([]string, 0, len(cfg.Modules))
	for _, mod := range cfg.Modules {
		if _, dup := seen[mod.id]; dup {
			return nil, fmt.Errorf("gogql: duplicate module id %q", mod.id)
		}
		seen[mod.id] = struct{}{}
		sources = append(sources, mod.typeDefs)
	}

	sdl, err := merge.TypeDefs(sources)
	if err != nil {
		return nil, fmt.Errorf("gogql: merge type definitions: %w", err)
	}

	resolver, err := mergeResolvers(cfg.Modules)
	if err != nil {
		return nil, err
	}

	schemaOpts := append([]graphql.SchemaOpt{}, cfg.SchemaOpts...)
	schemaOpts = append(schemaOpts, securitySchemaOpts(cfg.Security)...)
	if usesResolverMap(cfg.Modules) {
		schemaOpts = append(schemaOpts, graphql.UseFieldResolvers())
	}

	schema, err := graphql.ParseSchema(sdl, resolver, schemaOpts...)
	if err != nil {
		return nil, fmt.Errorf("gogql: parse schema: %w", err)
	}

	inj := newInjector()
	for _, mod := range cfg.Modules {
		for _, p := range mod.providers {
			p.apply(inj)
		}
	}

	return &Application{
		modules:         cfg.Modules,
		schema:          schema,
		schemaSDL:       sdl,
		injector:        inj,
		loaderFactories: cfg.Loaders,
	}, nil
}

// MustApplication calls NewApplication and panics on error.
func MustApplication(cfg ApplicationConfig) *Application {
	app, err := NewApplication(cfg)
	if err != nil {
		panic(err)
	}
	return app
}

// Schema returns the executable GraphQL schema.
func (a *Application) Schema() *graphql.Schema {
	return a.schema
}

// SchemaSDL returns the merged schema definition language document.
func (a *Application) SchemaSDL() string {
	return a.schemaSDL
}

// Injector returns the application-wide dependency injector (cloned per request in the server).
func (a *Application) Injector() *Injector {
	return a.injector
}

// RequestContext builds per-request context (injector, dataloaders).
func (a *Application) RequestContext(ctx context.Context) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	if _, ok := InjectorFrom(ctx); !ok {
		ctx = WithInjector(ctx, a.injector)
	}
	if len(a.loaderFactories) > 0 {
		if _, ok := LoaderRegistryFrom(ctx); !ok {
			ctx = WithLoaderRegistry(ctx, NewLoaderRegistry(a.loaderFactories))
		}
	}
	return ctx
}

// Execute runs a GraphQL operation with the application injector on the context.
func (a *Application) Execute(ctx context.Context, query, operationName string, variables map[string]any) *graphql.Response {
	ctx = a.RequestContext(ctx)
	return a.schema.Exec(ctx, query, operationName, variables)
}

// Subscribe runs a GraphQL subscription (also used by the WebSocket transport).
func (a *Application) Subscribe(ctx context.Context, query, operationName string, variables map[string]any) (<-chan any, error) {
	ctx = a.RequestContext(ctx)
	return a.schema.Subscribe(ctx, query, operationName, variables)
}

// ValidateStartup checks that the application schema was built successfully.
func (a *Application) ValidateStartup() error {
	if a.schema == nil {
		return fmt.Errorf("gogql: schema is nil")
	}
	return nil
}

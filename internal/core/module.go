/*
|--------------------------------------------------------------------------
| Module
|--------------------------------------------------------------------------
|
| A Module is one composable slice of schema (graphql-modules style): SDL
| text, resolver root (struct or ResolverMap), optional SubscriptionResolvers
| (struct methods only), and Providers for DI. Multiple modules each declare
| partial Query/Mutation types; merge combines them into one schema.
|
| ModuleConfig supports inline TypeDefs or split layouts via TypeDefParts,
| TypeDefFiles, and TypeDefsFS without breaking older single-string APIs.
| NewModule validates ID and SDL; MustModule panics on error for init code.
|
| Key types: Module, ModuleConfig. Key funcs: NewModule, MustModule, ID,
| TypeDefs.
|
*/

package core

import "io/fs"

// Module is a schema slice with its own SDL, resolvers, and providers (graphql-modules style).
type Module struct {
	id                    string
	typeDefs              string
	resolvers             any
	subscriptionResolvers any
	providers             []Provider
}

// ModuleConfig configures a GraphQL module.
//
// Backward compatible: only ID, TypeDefs, and Resolvers are required—the same as earlier gogql versions.
// TypeDefParts, TypeDefFiles, and TypeDefsFS are optional additions for splitting SDL across files;
// leave them zero/unset to keep a single inline TypeDefs string.
type ModuleConfig struct {
	ID       string
	TypeDefs string
	// TypeDefParts are optional extra SDL fragments merged after TypeDefs.
	TypeDefParts []string
	// TypeDefFiles are optional paths to .graphql files merged after TypeDefs and TypeDefParts.
	TypeDefFiles []string
	// TypeDefsFS optionally loads .graphql files under TypeDefsFSPath from embed.FS (or any fs.FS).
	TypeDefsFS     fs.FS
	TypeDefsFSPath string
	Resolvers      any
	// SubscriptionResolvers is a struct (or pointer) with methods for Subscription fields.
	// graph-gophers requires methods for subscription roots; use this instead of ResolverMap.Subscription.
	SubscriptionResolvers any
	Providers             []Provider
}

// NewModule creates a module from SDL type definitions and optional resolvers/providers.
func NewModule(cfg ModuleConfig) (*Module, error) {
	if cfg.ID == "" {
		return nil, errModuleIDRequired
	}
	typeDefs, err := compileModuleTypeDefs(cfg)
	if err != nil {
		return nil, err
	}
	return &Module{
		id:                    cfg.ID,
		typeDefs:              typeDefs,
		resolvers:             cfg.Resolvers,
		subscriptionResolvers: cfg.SubscriptionResolvers,
		providers:             cfg.Providers,
	}, nil
}

// MustModule calls NewModule and panics on error.
func MustModule(cfg ModuleConfig) *Module {
	m, err := NewModule(cfg)
	if err != nil {
		panic(err)
	}
	return m
}

// ID returns the module identifier.
func (m *Module) ID() string {
	return m.id
}

// TypeDefs returns the module SDL.
func (m *Module) TypeDefs() string {
	return m.typeDefs
}

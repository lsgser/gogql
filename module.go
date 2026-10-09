package gogql

// Module is a schema slice with its own SDL, resolvers, and providers (graphql-modules style).
type Module struct {
	id                    string
	typeDefs              string
	resolvers             any
	subscriptionResolvers any
	providers             []Provider
}

// ModuleConfig configures a GraphQL module.
type ModuleConfig struct {
	ID        string
	TypeDefs  string
	Resolvers any
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
	if cfg.TypeDefs == "" {
		return nil, errModuleTypeDefsRequired
	}
	return &Module{
		id:                    cfg.ID,
		typeDefs:              cfg.TypeDefs,
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

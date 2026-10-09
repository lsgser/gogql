/*
|--------------------------------------------------------------------------
| DataLoader
|--------------------------------------------------------------------------
|
| Per-request DataLoader registry and helpers (LoaderFactories, LoadersFromContext,
| NewStringKeyLoader) built on graph-gophers/dataloader.
|
*/

package core

import (
	"context"
	"fmt"
	"sync"

	"github.com/graph-gophers/dataloader"
)

type loaderRegistryKey struct{}

// LoaderRegistry holds request-scoped DataLoaders (batched per request).
type LoaderRegistry struct {
	mu      sync.Mutex
	loaders map[string]*dataloader.Loader
	factory map[string]func() *dataloader.Loader
}

// NewLoaderRegistry creates an empty registry. Register factories before the request runs.
func NewLoaderRegistry(factories map[string]func() *dataloader.Loader) *LoaderRegistry {
	return &LoaderRegistry{
		loaders: make(map[string]*dataloader.Loader),
		factory: factories,
	}
}

// Loader returns a named loader, creating it once per request.
func (r *LoaderRegistry) Loader(name string) (*dataloader.Loader, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if l, ok := r.loaders[name]; ok {
		return l, nil
	}
	f, ok := r.factory[name]
	if !ok {
		return nil, fmt.Errorf("gogql: loader %q is not registered", name)
	}
	l := f()
	r.loaders[name] = l
	return l, nil
}

// MustLoader returns a loader or panics.
func (r *LoaderRegistry) MustLoader(name string) *dataloader.Loader {
	l, err := r.Loader(name)
	if err != nil {
		panic(err)
	}
	return l
}

// WithLoaderRegistry attaches loaders to the context.
func WithLoaderRegistry(ctx context.Context, reg *LoaderRegistry) context.Context {
	return context.WithValue(ctx, loaderRegistryKey{}, reg)
}

// LoaderRegistryFrom returns the request loader registry.
func LoaderRegistryFrom(ctx context.Context) (*LoaderRegistry, bool) {
	reg, ok := ctx.Value(loaderRegistryKey{}).(*LoaderRegistry)
	return reg, ok
}

// LoadersFromContext returns a named loader from the request context.
func LoadersFromContext(ctx context.Context, name string) (*dataloader.Loader, error) {
	reg, ok := LoaderRegistryFrom(ctx)
	if !ok {
		return nil, fmt.Errorf("gogql: no loader registry on context")
	}
	return reg.Loader(name)
}

// LoaderFactories maps loader names to constructors (one instance per GraphQL request).
type LoaderFactories map[string]func() *dataloader.Loader

// NewStringKeyLoader builds a batched loader for string keys.
func NewStringKeyLoader(batch func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result) *dataloader.Loader {
	return dataloader.NewBatchedLoader(batch)
}

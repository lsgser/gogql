/*
|--------------------------------------------------------------------------
| Injector
|--------------------------------------------------------------------------
|
| Request-scoped dependency injection (providers, MustGet, WithInjector)
| similar to graphql-modules service locators.
|
*/

package core

import (
	"context"
	"fmt"
	"reflect"
	"sync"
)

type injectorKey struct{}

// Injector resolves module providers (similar to graphql-modules dependency injection).
type Injector struct {
	mu        sync.RWMutex
	byType    map[reflect.Type]any
	byToken   map[string]any
	factories map[reflect.Type]func(*Injector) (any, error)
}

func newInjector() *Injector {
	return &Injector{
		byType:    make(map[reflect.Type]any),
		byToken:   make(map[string]any),
		factories: make(map[reflect.Type]func(*Injector) (any, error)),
	}
}

// Register adds a singleton value available by its concrete type.
func (i *Injector) Register(value any) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.byType[reflect.TypeOf(value)] = value
}

// RegisterToken adds a singleton available by string token.
func (i *Injector) RegisterToken(token string, value any) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.byToken[token] = value
}

func (i *Injector) registerFactory(t reflect.Type, factory func(*Injector) (any, error)) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.factories[t] = factory
}

// Get resolves a dependency by type.
func (i *Injector) Get(ptr any) error {
	target := reflect.ValueOf(ptr)
	if target.Kind() != reflect.Pointer || target.IsNil() {
		return fmt.Errorf("gogql: Get expects a non-nil pointer")
	}
	elem := target.Elem()
	val, err := i.resolve(elem.Type())
	if err != nil {
		return err
	}
	elem.Set(reflect.ValueOf(val))
	return nil
}

func (i *Injector) resolve(t reflect.Type) (any, error) {
	i.mu.RLock()
	if v, ok := i.byType[t]; ok {
		i.mu.RUnlock()
		return v, nil
	}
	factory := i.factories[t]
	i.mu.RUnlock()

	if factory != nil {
		v, err := factory(i)
		if err != nil {
			return nil, err
		}
		i.mu.Lock()
		i.byType[t] = v
		i.mu.Unlock()
		return v, nil
	}

	return nil, fmt.Errorf("gogql: no provider registered for %s", t)
}

// WithInjector attaches an injector to the GraphQL request context.
func WithInjector(ctx context.Context, inj *Injector) context.Context {
	return context.WithValue(ctx, injectorKey{}, inj)
}

// InjectorFrom returns the request injector, if any.
func InjectorFrom(ctx context.Context) (*Injector, bool) {
	inj, ok := ctx.Value(injectorKey{}).(*Injector)
	return inj, ok
}

// MustGet returns a dependency from the request context injector.
func MustGet[T any](ctx context.Context) T {
	var out T
	inj, ok := InjectorFrom(ctx)
	if !ok {
		panic("gogql: no injector on context")
	}
	if err := inj.Get(&out); err != nil {
		panic(err)
	}
	return out
}

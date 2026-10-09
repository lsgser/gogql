/*
|--------------------------------------------------------------------------
| Provider
|--------------------------------------------------------------------------
|
| Module provider registration: Provide, ProvideToken, and ProvideFactory
| for wiring services into the application injector.
|
*/

package core

import "reflect"

// Provider registers a service in the application injector.
type Provider struct {
	token   string
	value   any
	factory func(*Injector) (any, error)
	typ     reflect.Type
}

// Provide registers a singleton value by its concrete type.
func Provide(value any) Provider {
	return Provider{value: value, typ: reflect.TypeOf(value)}
}

// ProvideToken registers a singleton value by token.
func ProvideToken(token string, value any) Provider {
	return Provider{token: token, value: value}
}

// ProvideFactory registers a lazy factory. The returned value is cached as a singleton.
func ProvideFactory[T any](factory func(*Injector) (T, error)) Provider {
	var zero T
	return Provider{
		typ: reflect.TypeOf(zero),
		factory: func(inj *Injector) (any, error) {
			return factory(inj)
		},
	}
}

func (p Provider) apply(inj *Injector) {
	if p.token != "" {
		inj.RegisterToken(p.token, p.value)
		return
	}
	if p.value != nil {
		inj.Register(p.value)
		return
	}
	if p.factory != nil && p.typ != nil {
		inj.registerFactory(p.typ, p.factory)
	}
}

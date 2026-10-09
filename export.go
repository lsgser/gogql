/*
|--------------------------------------------------------------------------
| Public API
|--------------------------------------------------------------------------
|
| Re-exports from internal/core so consumers import github.com/lsgser/gogql.
|
*/

package gogql

import (
	"context"

	"github.com/lsgser/gogql/internal/core"
)

type (
	Application       = core.Application
	ApplicationConfig = core.ApplicationConfig
	AuthClaims        = core.AuthClaims
	Injector          = core.Injector
	LoaderFactories   = core.LoaderFactories
	LoaderRegistry    = core.LoaderRegistry
	Module            = core.Module
	ModuleConfig      = core.ModuleConfig
	PlaygroundConfig  = core.PlaygroundConfig
	PlaygroundUI      = core.PlaygroundUI
	ResolverMap       = core.ResolverMap
	SecurityConfig    = core.SecurityConfig
	Server            = core.Server
	ServerConfig      = core.ServerConfig
)

// Provider is a DI registration handle (see Provide, ProvideToken, ProvideFactory).
type Provider = core.Provider

const (
	PlaygroundGraphiQL      = core.PlaygroundGraphiQL
	PlaygroundApolloSandbox = core.PlaygroundApolloSandbox
)

var (
	AuthClaimsFrom     = core.AuthClaimsFrom
	InjectorFrom       = core.InjectorFrom
	JoinTypeDefs       = core.JoinTypeDefs
	LoadTypeDefsFS     = core.LoadTypeDefsFS
	LoadersFromContext = core.LoadersFromContext
	LoaderRegistryFrom = core.LoaderRegistryFrom
	MustApplication    = core.MustApplication
	MustAuthClaims     = core.MustAuthClaims
	MustLoadTypeDefsFS = core.MustLoadTypeDefsFS
	MustModule         = core.MustModule
	NewApplication     = core.NewApplication
	NewLoaderRegistry  = core.NewLoaderRegistry
	NewModule          = core.NewModule
	NewResolverMap     = core.NewResolverMap
	NewServer          = core.NewServer
	NewStringKeyLoader = core.NewStringKeyLoader
	Provide            = core.Provide
	ProvideToken       = core.ProvideToken
	WithAuthClaims     = core.WithAuthClaims
	WithInjector       = core.WithInjector
	WithLoaderRegistry = core.WithLoaderRegistry
)

func MustGet[T any](ctx context.Context) T {
	return core.MustGet[T](ctx)
}

func ProvideFactory[T any](factory func(*Injector) (T, error)) Provider {
	return core.ProvideFactory(factory)
}

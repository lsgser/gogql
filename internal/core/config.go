/*
|--------------------------------------------------------------------------
| Configuration types
|--------------------------------------------------------------------------
|
| Shared config structs for Server and Application. PlaygroundConfig selects
| GraphiQL vs Apollo Sandbox, path, and enabled flag. SecurityConfig maps to
| graph-gophers MaxDepth via securitySchemaOpts during schema parse.
|
| playgroundSettings is internal normalized state derived from ServerConfig
| (defaults, deprecated GraphiQL bool). PlaygroundUI constants identify UI kind.
|
| Key types: PlaygroundConfig, PlaygroundUI, SecurityConfig.
|
*/

package core

import "github.com/graph-gophers/graphql-go"

// PlaygroundUI selects which interactive GraphQL IDE to serve.
type PlaygroundUI string

const (
	// PlaygroundGraphiQL serves the GraphiQL UI (default).
	PlaygroundGraphiQL PlaygroundUI = "graphiql"
	// PlaygroundApolloSandbox serves Apollo Sandbox (embedded).
	PlaygroundApolloSandbox PlaygroundUI = "apollo-sandbox"
)

// PlaygroundConfig controls the GraphQL playground route and UI.
type PlaygroundConfig struct {
	// Enabled turns the playground on (default: true when nil pointer not used — see ServerConfig).
	Enabled bool
	// Path is the HTTP path for the playground (default: /playground).
	// GraphQL queries still use ServerConfig.GraphQLPath.
	Path string
	// UI selects GraphiQL or Apollo Sandbox.
	UI PlaygroundUI
}

func (p PlaygroundConfig) normalized(graphqlPath string) playgroundSettings {
	path := p.Path
	if path == "" {
		path = "/playground"
	}
	ui := p.UI
	if ui == "" {
		ui = PlaygroundGraphiQL
	}
	return playgroundSettings{
		enabled:      p.Enabled,
		path:         path,
		ui:           ui,
		graphqlPath:  graphqlPath,
		legacyOnPath: false,
	}
}

type playgroundSettings struct {
	enabled      bool
	path         string
	ui           PlaygroundUI
	graphqlPath  string
	legacyOnPath bool
}

// SecurityConfig holds query execution limits.
type SecurityConfig struct {
	// MaxDepth limits nested field selections (0 = no limit).
	MaxDepth int
}

func securitySchemaOpts(sec SecurityConfig) []graphql.SchemaOpt {
	if sec.MaxDepth <= 0 {
		return nil
	}
	return []graphql.SchemaOpt{graphql.MaxDepth(sec.MaxDepth)}
}

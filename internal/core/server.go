/*
|--------------------------------------------------------------------------
| Server
|--------------------------------------------------------------------------
|
| HTTP serving layer inspired by GraphQL Yoga. NewServer registers routes on
| an internal mux: GraphQL (POST JSON, GET query string, OPTIONS CORS),
| optional playground HTML (GraphiQL or Apollo Sandbox), and a plain-text
| health endpoint. When subscriptions are enabled, the GraphQL path also
| accepts graphql-transport-ws upgrades via graph-gophers/graphql-transport-ws.
|
| ServerConfig.ContextFunc runs first on each HTTP or WS request (auth, tracing),
| then Application.RequestContext attaches injector and DataLoaders. Handler()
| exposes the mux for Gin or other wrappers; ListenAndServe is a convenience
| wrapper around http.ListenAndServe.
|
| Key types: Server, ServerConfig. Key funcs: NewServer, Handler, ListenAndServe.
|
*/

package core

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	graphqlws "github.com/graph-gophers/graphql-transport-ws"
)

// Server is a Yoga-inspired HTTP server for GraphQL.
type Server struct {
	app           *Application
	mux           *http.ServeMux
	path          string
	healthPath    string
	playground    playgroundSettings
	subscriptions bool
	contextFunc   func(context.Context, *http.Request) context.Context
}

// ServerConfig configures the HTTP server.
type ServerConfig struct {
	// GraphQLPath is the POST/GET endpoint (default: /graphql).
	GraphQLPath string
	// Playground configures the GraphQL playground (GraphiQL or Apollo Sandbox).
	// When nil, the playground is enabled at /playground with GraphiQL.
	Playground *PlaygroundConfig
	// EnableSubscriptions enables graphql-transport-ws on GraphQLPath (default: true).
	EnableSubscriptions *bool
	// HealthPath is an optional liveness endpoint (default: /health).
	HealthPath string

	// GraphiQL is deprecated: use Playground instead. When set, it overrides Playground.Enabled.
	GraphiQL *bool
	// ContextFunc runs before the application injector/loaders are attached.
	// Use it for JWT validation, request IDs, or tracing (see docs/database-and-auth.md).
	ContextFunc func(context.Context, *http.Request) context.Context
}

// NewServer creates an HTTP handler stack for the application.
func NewServer(app *Application, cfg ServerConfig) *Server {
	path := cfg.GraphQLPath
	if path == "" {
		path = "/graphql"
	}
	health := cfg.HealthPath
	if health == "" {
		health = "/health"
	}

	pg := PlaygroundConfig{Enabled: true, UI: PlaygroundGraphiQL}
	if cfg.Playground != nil {
		pg = *cfg.Playground
	}
	if cfg.GraphiQL != nil {
		pg.Enabled = *cfg.GraphiQL
	}
	settings := pg.normalized(path)

	subscriptions := true
	if cfg.EnableSubscriptions != nil {
		subscriptions = *cfg.EnableSubscriptions
	}

	s := &Server{
		app:           app,
		mux:           http.NewServeMux(),
		path:          path,
		healthPath:    health,
		playground:    settings,
		subscriptions: subscriptions,
		contextFunc:   cfg.ContextFunc,
	}

	gqlHTTP := http.HandlerFunc(s.handleGraphQLHTTP)
	var gqlHandler http.Handler = gqlHTTP
	if s.subscriptions {
		gqlHandler = http.HandlerFunc(graphqlws.NewHandlerFunc(app, gqlHTTP,
			graphqlws.WithContextGenerator(graphqlws.ContextGeneratorFunc(func(ctx context.Context, r *http.Request) (context.Context, error) {
				return s.buildRequestContext(r), nil
			})),
		))
	}

	s.mux.Handle(path, gqlHandler)

	if settings.enabled {
		s.mux.HandleFunc(settings.path, s.handlePlayground)
	}

	s.mux.HandleFunc(health, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	return s
}

// Handler returns the root HTTP handler.
func (s *Server) Handler() http.Handler {
	return s.mux
}

// ListenAndServe starts the server on addr (e.g. ":8080").
func (s *Server) ListenAndServe(addr string) error {
	if addr == "" {
		addr = ":8080"
	}
	return http.ListenAndServe(addr, s.Handler())
}

type graphQLRequest struct {
	Query         string         `json:"query"`
	OperationName string         `json:"operationName"`
	Variables     map[string]any `json:"variables"`
}

func (s *Server) handleGraphQLHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		s.servePOST(w, r)
	case http.MethodGet:
		s.serveGET(w, r)
	case http.MethodOptions:
		setCORS(w)
		w.WriteHeader(http.StatusNoContent)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) handlePlayground(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	page, err := renderPlayground(s.playground)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(page))
}

func (s *Server) servePOST(w http.ResponseWriter, r *http.Request) {
	setCORS(w)
	w.Header().Set("Content-Type", "application/json")

	body, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid body"})
		return
	}

	var req graphQLRequest
	if err := json.Unmarshal(body, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		return
	}
	if req.Query == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "query is required"})
		return
	}

	ctx := s.buildRequestContext(r)
	result := s.app.schema.Exec(ctx, req.Query, req.OperationName, req.Variables)
	writeJSON(w, http.StatusOK, result)
}

func (s *Server) serveGET(w http.ResponseWriter, r *http.Request) {
	setCORS(w)
	w.Header().Set("Content-Type", "application/json")

	query := r.URL.Query().Get("query")
	if query == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "query parameter is required"})
		return
	}
	var vars map[string]any
	if raw := r.URL.Query().Get("variables"); raw != "" {
		if err := json.Unmarshal([]byte(raw), &vars); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid variables JSON"})
			return
		}
	}

	ctx := s.buildRequestContext(r)
	result := s.app.schema.Exec(ctx, query, r.URL.Query().Get("operationName"), vars)
	writeJSON(w, http.StatusOK, result)
}

func (s *Server) buildRequestContext(r *http.Request) context.Context {
	ctx := r.Context()
	if s.contextFunc != nil {
		ctx = s.contextFunc(ctx, r)
	}
	return s.app.RequestContext(ctx)
}

func setCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "content-type, authorization")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
}


func writeJSON(w http.ResponseWriter, status int, v any) {
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(true)
	if err := enc.Encode(v); err != nil {
		http.Error(w, fmt.Sprintf("encode response: %v", err), http.StatusInternalServerError)
	}
}

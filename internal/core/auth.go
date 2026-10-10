/*
|--------------------------------------------------------------------------
| Auth claims
|--------------------------------------------------------------------------
|
| Minimal auth surface: AuthClaims holds Subject (typically JWT "sub").
| ServerConfig.ContextFunc should validate credentials and call WithAuthClaims;
| resolvers call AuthClaimsFrom or MustAuthClaims to gate fields like "me".
|
| gogql does not ship JWT or session middleware—only context storage so you
| can plug any auth library (see examples/database and src/utils/auth.go).
|
| Key type: AuthClaims. Key funcs: WithAuthClaims, AuthClaimsFrom, MustAuthClaims.
|
*/

package core

import "context"

type authClaimsKey struct{}

// AuthClaims holds authenticated user information attached to a request context.
type AuthClaims struct {
	// Subject is typically the user ID (JWT "sub").
	Subject string
}

// WithAuthClaims stores claims on the context (e.g. after JWT validation).
func WithAuthClaims(ctx context.Context, claims AuthClaims) context.Context {
	return context.WithValue(ctx, authClaimsKey{}, claims)
}

// AuthClaimsFrom returns claims when the request was authenticated.
func AuthClaimsFrom(ctx context.Context) (AuthClaims, bool) {
	c, ok := ctx.Value(authClaimsKey{}).(AuthClaims)
	return c, ok
}

// MustAuthClaims returns claims or panics if the request is unauthenticated.
func MustAuthClaims(ctx context.Context) AuthClaims {
	c, ok := AuthClaimsFrom(ctx)
	if !ok || c.Subject == "" {
		panic("gogql: missing auth claims on context")
	}
	return c
}

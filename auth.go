package gogql

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

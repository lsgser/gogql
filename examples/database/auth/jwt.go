package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lsgser/gogql"
)

// DemoSecret is for local runs only — use os.Getenv("JWT_SECRET") in production.
const DemoSecret = "dev-only-change-me"

type jwtClaims struct {
	jwt.RegisteredClaims
}

// ContextFunc validates Bearer JWT and attaches gogql.AuthClaims.
func ContextFunc(secret []byte) func(context.Context, *http.Request) context.Context {
	return func(ctx context.Context, r *http.Request) context.Context {
		raw := r.Header.Get("Authorization")
		if raw == "" {
			return ctx
		}
		const prefix = "Bearer "
		if !strings.HasPrefix(raw, prefix) {
			return ctx
		}
		tokenStr := strings.TrimSpace(strings.TrimPrefix(raw, prefix))
		token, err := jwt.ParseWithClaims(tokenStr, &jwtClaims{}, func(t *jwt.Token) (any, error) {
			if t.Method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return secret, nil
		})
		if err != nil || !token.Valid {
			return ctx
		}
		claims, ok := token.Claims.(*jwtClaims)
		if !ok || claims.Subject == "" {
			return ctx
		}
		return gogql.WithAuthClaims(ctx, gogql.AuthClaims{Subject: claims.Subject})
	}
}

// MintDemoToken issues a HS256 token for the database example.
func MintDemoToken(userID string, ttl time.Duration) (string, error) {
	claims := jwtClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(DemoSecret))
}

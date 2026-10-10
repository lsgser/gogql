/*
|--------------------------------------------------------------------------
| Auth tests
|--------------------------------------------------------------------------
|
| Verifies WithAuthClaims round-trips through AuthClaimsFrom / MustAuthClaims
| on context, matching how ServerConfig.ContextFunc attaches identity.
|
*/

package core_test

import (
	"context"
	"testing"

	"github.com/lsgser/gogql"
)

func TestAuthClaimsContext(t *testing.T) {
	ctx := gogql.WithAuthClaims(context.Background(), gogql.AuthClaims{Subject: "42"})
	c, ok := gogql.AuthClaimsFrom(ctx)
	if !ok || c.Subject != "42" {
		t.Fatalf("claims: %+v ok=%v", c, ok)
	}
}

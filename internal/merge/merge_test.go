/*
|--------------------------------------------------------------------------
| SDL merge tests
|--------------------------------------------------------------------------
|
| Tests merging duplicate type Query definitions across module documents.
|
*/

package merge_test

import (
	"strings"
	"testing"

	"github.com/lsgser/gogql/internal/merge"
)

func TestMergeDuplicateQueryTypes(t *testing.T) {
	sdl, err := merge.TypeDefs([]string{
		`type Query { user(id: ID!): String }`,
		`type Query { posts: [String!]! }`,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(sdl, "user") || !strings.Contains(sdl, "posts") {
		t.Fatalf("expected merged query fields, got:\n%s", sdl)
	}
}

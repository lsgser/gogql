/*
|--------------------------------------------------------------------------
| Name helper tests
|--------------------------------------------------------------------------
|
| Unit tests for SanitizePackageName, TypeName, and QueryFieldNames used
| when scaffolding domain modules from CLI arguments.
|
*/

package scaffold

import "testing"

func TestSanitizePackageName(t *testing.T) {
	if got := SanitizePackageName("blog-posts"); got != "blog_posts" {
		t.Fatalf("got %q", got)
	}
}

func TestTypeName(t *testing.T) {
	if TypeName("users") != "User" {
		t.Fatal("users -> User")
	}
	if TypeName("posts") != "Post" {
		t.Fatal("posts -> Post")
	}
}

func TestQueryFieldNames(t *testing.T) {
	one, many := QueryFieldNames("users")
	if one != "user" || many != "users" {
		t.Fatalf("got %q %q", one, many)
	}
}

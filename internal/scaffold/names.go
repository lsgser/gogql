/*
|--------------------------------------------------------------------------
| Name helpers
|--------------------------------------------------------------------------
|
| Converts CLI module names into Go package names and GraphQL type names.
|
*/

package scaffold

import (
	"strings"
	"unicode"
)

// SanitizePackageName returns a valid Go package identifier (lowercase).
func SanitizePackageName(name string) string {
	name = strings.TrimSpace(name)
	var b strings.Builder
	for _, r := range name {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
			b.WriteRune(unicode.ToLower(r))
		} else if r == '-' || r == ' ' {
			b.WriteRune('_')
		}
	}
	out := b.String()
	if out == "" {
		return "module"
	}
	if unicode.IsDigit(rune(out[0])) {
		return "m_" + out
	}
	return out
}

// TypeName returns exported singular type name (e.g. posts -> Post).
func TypeName(packageName string) string {
	parts := strings.Split(strings.ReplaceAll(packageName, "_", " "), " ")
	var words []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		words = append(words, capitalizeWord(p))
	}
	if len(words) == 0 {
		return "Item"
	}
	base := strings.Join(words, "")
	if strings.HasSuffix(strings.ToLower(base), "s") && len(base) > 1 {
		return base[:len(base)-1]
	}
	return base
}

// QueryFieldNames returns GraphQL query field names (singular, plural list).
func QueryFieldNames(packageName string) (one, many string) {
	t := TypeName(packageName)
	if t == "" {
		return "item", "items"
	}
	one = strings.ToLower(t[:1]) + t[1:]
	many = SanitizePackageName(packageName)
	if many == one {
		many = one + "s"
	}
	return one, many
}

func capitalizeWord(s string) string {
	if s == "" {
		return s
	}
	r := []rune(strings.ToLower(s))
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

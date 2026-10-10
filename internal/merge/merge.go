/*
|--------------------------------------------------------------------------
| SDL merge
|--------------------------------------------------------------------------
|
| Combines multiple module SDL strings into one document suitable for
| graph-gophers ParseSchema. Uses vektah/gqlparser to parse each fragment,
| merges duplicate root types (Query, Mutation, Subscription) via
| normalizeRootTypes, strips redundant built-in scalar declarations, and
| formats the result. Application calls merge.TypeDefs at startup.
|
| Key func: TypeDefs(sources []string) (string, error).
|
*/

package merge

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/formatter"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"github.com/vektah/gqlparser/v2/parser"
	"github.com/vektah/gqlparser/v2/validator"
)

var rootObjectNames = map[string]struct{}{
	"Query":        {},
	"Mutation":     {},
	"Subscription": {},
}

// TypeDefs merges GraphQL SDL fragments from multiple modules into one executable schema string.
// Each module may declare root operations as `type Query { ... }`; duplicate roots are turned into extensions automatically.
func TypeDefs(sources []string) (string, error) {
	if len(sources) == 0 {
		return "", fmt.Errorf("merge: at least one type definition document is required")
	}

	seenRoots := make(map[string]struct{})
	var inputs []*ast.Source

	for i, raw := range sources {
		sdl := strings.TrimSpace(raw)
		if sdl == "" {
			return "", fmt.Errorf("merge: module %d has empty type definitions", i)
		}

		doc, err := parser.ParseSchema(&ast.Source{
			Input: sdl,
			Name:  fmt.Sprintf("module-%d", i),
		})
		if err != nil {
			return "", gqlerror.WrapIfUnwrapped(err)
		}

		normalizeRootTypes(doc, seenRoots)
		inputs = append(inputs, &ast.Source{
			Input: mustFormat(doc),
			Name:  fmt.Sprintf("module-%d-normalized", i),
		})
	}

	allInputs := make([]*ast.Source, 0, len(inputs)+1)
	allInputs = append(allInputs, validator.Prelude)
	allInputs = append(allInputs, inputs...)

	merged, err := parser.ParseSchemas(allInputs...)
	if err != nil {
		return "", gqlerror.WrapIfUnwrapped(err)
	}

	ensureSchemaDefinition(merged)

	if _, err := validator.ValidateSchemaDocument(merged); err != nil {
		return "", gqlerror.WrapIfUnwrapped(err)
	}

	// Validation merges extensions into type definitions in-place; drop extensions to avoid duplicate SDL.
	merged.Extensions = nil
	merged.SchemaExtension = nil

	return mustFormat(stripBuiltInScalars(merged)), nil
}

func normalizeRootTypes(doc *ast.SchemaDocument, seenRoots map[string]struct{}) {
	var defs ast.DefinitionList
	for _, def := range doc.Definitions {
		if def.Kind != ast.Object {
			defs = append(defs, def)
			continue
		}
		if _, isRoot := rootObjectNames[def.Name]; !isRoot {
			defs = append(defs, def)
			continue
		}
		if _, seen := seenRoots[def.Name]; seen {
			doc.Extensions = append(doc.Extensions, def)
			continue
		}
		seenRoots[def.Name] = struct{}{}
		defs = append(defs, def)
	}
	doc.Definitions = defs
}

func ensureSchemaDefinition(doc *ast.SchemaDocument) {
	if len(doc.Schema) > 0 || len(doc.SchemaExtension) > 0 {
		return
	}

	ops := ast.OperationTypeDefinitionList{}
	if hasType(doc, "Query") {
		ops = append(ops, &ast.OperationTypeDefinition{Operation: ast.Query, Type: "Query"})
	}
	if hasType(doc, "Mutation") {
		ops = append(ops, &ast.OperationTypeDefinition{Operation: ast.Mutation, Type: "Mutation"})
	}
	if hasType(doc, "Subscription") {
		ops = append(ops, &ast.OperationTypeDefinition{Operation: ast.Subscription, Type: "Subscription"})
	}
	if len(ops) == 0 {
		return
	}

	doc.Schema = append(doc.Schema, &ast.SchemaDefinition{OperationTypes: ops})
}

func hasType(doc *ast.SchemaDocument, name string) bool {
	for _, def := range doc.Definitions {
		if def.Name == name {
			return true
		}
	}
	for _, ext := range doc.Extensions {
		if ext.Name == name {
			return true
		}
	}
	return false
}

func mustFormat(doc *ast.SchemaDocument) string {
	var buf bytes.Buffer
	formatter.NewFormatter(&buf).FormatSchemaDocument(doc)
	return buf.String()
}

func stripBuiltInScalars(doc *ast.SchemaDocument) *ast.SchemaDocument {
	out := &ast.SchemaDocument{
		Schema:          doc.Schema,
		SchemaExtension: doc.SchemaExtension,
		Directives:      doc.Directives,
	}
	for _, def := range doc.Definitions {
		if def.BuiltIn {
			continue
		}
		out.Definitions = append(out.Definitions, def)
	}
	for _, ext := range doc.Extensions {
		if ext.BuiltIn {
			continue
		}
		out.Extensions = append(out.Extensions, ext)
	}
	return out
}

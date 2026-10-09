/*
|--------------------------------------------------------------------------
| Basic example — users SDL
|--------------------------------------------------------------------------
|
| Loads and merges schema/*.graphql via embed.FS for the users module.
|
*/

package users

import (
	"embed"

	"github.com/lsgser/gogql"
)

//go:embed schema/*.graphql
var schemaFS embed.FS

func typeDefs() string {
	return gogql.MustLoadTypeDefsFS(schemaFS, "schema")
}

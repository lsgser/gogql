/*
|--------------------------------------------------------------------------
| Basic example — users SDL
|--------------------------------------------------------------------------
|
| Loads schema/*.graphql from embed.FS via gogql.MustLoadTypeDefsFS. Equivalent
| to co-located user.graphql in the init template; keeps SDL out of Go strings.
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

/*
|--------------------------------------------------------------------------
| Database example — module registry
|--------------------------------------------------------------------------
|
| Composes users module and exposes LoaderFactories for the application.
|
*/

package modules

import (
	"github.com/lsgser/gogql"
	"github.com/lsgser/gogql/examples/database/modules/users"
	"github.com/lsgser/gogql/examples/database/store"
)

// All composes gogql modules for the database example.
func All(repo *store.UserRepository) []*gogql.Module {
	return []*gogql.Module{
		users.Module(repo),
	}
}

// LoaderFactories exposes DataLoader factories for the application.
func LoaderFactories(repo *store.UserRepository) gogql.LoaderFactories {
	return users.LoaderFactories(repo)
}

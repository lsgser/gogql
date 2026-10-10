/*
|--------------------------------------------------------------------------
| Subscriptions example — module registry
|--------------------------------------------------------------------------
|
| Returns only the events module; extend when adding more subscription domains.
|
*/

package modules

import (
	"github.com/lsgser/gogql"
	"github.com/lsgser/gogql/examples/subscriptions/modules/events"
)

// All registers modules for the subscriptions example application.
func All() []*gogql.Module {
	return []*gogql.Module{
		events.Module(),
	}
}

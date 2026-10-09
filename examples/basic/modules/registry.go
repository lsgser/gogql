package modules

import (
	"github.com/lsgser/gogql"
	"github.com/lsgser/gogql/examples/basic/modules/greeting"
	"github.com/lsgser/gogql/examples/basic/modules/users"
)

// All returns every module composed into one gogql application.
func All() []*gogql.Module {
	return []*gogql.Module{
		users.Module(),
		greeting.Module(),
	}
}

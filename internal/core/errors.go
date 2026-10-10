/*
|--------------------------------------------------------------------------
| Errors
|--------------------------------------------------------------------------
|
| Private sentinel errors returned from NewModule when ModuleConfig is invalid
| (missing ID or empty type definitions). Keeps validation messages consistent
| and avoids string matching in callers—MustModule surfaces them as panics.
|
| Vars: errModuleIDRequired, errModuleTypeDefsRequired.
|
*/

package core

import "errors"

var (
	errModuleIDRequired       = errors.New("gogql: module ID is required")
	errModuleTypeDefsRequired = errors.New("gogql: module type definitions are required")
)

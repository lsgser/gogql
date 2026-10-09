/*
|--------------------------------------------------------------------------
| Errors
|--------------------------------------------------------------------------
|
| Package-level sentinel errors for module and configuration validation.
|
*/

package core

import "errors"

var (
	errModuleIDRequired       = errors.New("gogql: module ID is required")
	errModuleTypeDefsRequired = errors.New("gogql: module type definitions are required")
)

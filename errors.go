package gogql

import "errors"

var (
	errModuleIDRequired       = errors.New("gogql: module ID is required")
	errModuleTypeDefsRequired = errors.New("gogql: module type definitions are required")
)

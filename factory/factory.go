// Package factory provides functions that creates financial calculators
package factory

import "errors"

// Errors this package may return and can be checked with errors.Cause()
var (
	ErrRegionNotExist = errors.New("unkown region")
	ErrFactoryNotInit = errors.New("factory is improperly initialized")
)

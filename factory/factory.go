// Package factory provides functions that creates financial calculators
package factory

import "errors"

// Errors this package may return and can be checked with errors.Cause()
var (
	ErrRegionNotExist = errors.New("unkown region")
)

// Options is used to pass options to factory constructors
type Options struct {
	Year   uint
	Region Region
}

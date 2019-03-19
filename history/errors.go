package history

import "errors"

// Sentinel errors that can be wrapped and returned
var (
	ErrNoRates        = errors.New("no tax rates available")
	ErrNoJurisdiction = errors.New("unknown jurisdication")
)

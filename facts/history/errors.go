package history

import "errors"

// Sentinel errors that can be wrapped and returned
var (
	ErrNoFacts    = errors.New("no tax facts available")
	ErrNoProvince = errors.New("unknown province")
)

package history

import "errors"

// Sentinel errors that can be wrapped and returned
var (
	ErrParamsNotExist       = errors.New("no parameters exists")
	ErrJurisdictionNotExist = errors.New("unknown jurisdication")

	errNilFormula       = errors.New("nil formula encountered")
	errNilContraFormula = errors.New("nil contra-formula encountered")
)

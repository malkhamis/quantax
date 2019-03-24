package history

import "errors"

// Sentinel errors that can be wrapped and returned
var (
	ErrFormulaNotExist      = errors.New("no formula exists")
	ErrJurisdictionNotExist = errors.New("unknown jurisdication")
)

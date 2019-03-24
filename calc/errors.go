package calc

import "errors"

// Sentinel errors that can be returned and wrapped by this package
var (
	ErrValZero        = errors.New("zero value is not allowed")
	ErrValNeg         = errors.New("negative value is not allowed")
	ErrValInf         = errors.New("infinity value is not allowed")
	ErrValInfNeg      = errors.New("negative infinity value is not allowed")
	ErrValInfPos      = errors.New("positive infinity value is not allowed")
	ErrBoundsReversed = errors.New("lower-bound is greater than upper-bound")
	ErrInvalidDate    = errors.New("invalid date")
	ErrNoFormula      = errors.New("no formula given")
	ErrInvalidAge     = errors.New("invalid age")
)

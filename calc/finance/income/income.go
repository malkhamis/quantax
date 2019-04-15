// Package income provide implementation for calc.IncomeCalculator interface
package income

import "errors"

// Sentinel errors that can ben wrapped and returned by this package
var (
	ErrNoRecipe = errors.New("no income recipe given/set")
)

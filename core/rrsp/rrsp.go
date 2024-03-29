// Package rrsp provides implementations for the RRSPCalculator interface
// defined in package calc
package rrsp

import (
	"github.com/malkhamis/quantax/core"
	"github.com/pkg/errors"
)

// Sentinel errors that can ben wrapped and returned by this package
var (
	ErrNoFormula = errors.New("not formula given/set")
	ErrNoTaxCalc = errors.New("no tax calculator given")
)

// Formula computes the max contribution room acquired for a given income
type Formula interface {
	// Contribution returns the contribution room acquired
	// given the net income
	ContributionEarned(netIncome float64) float64
	// AllowedIncomeSources return the allowed income sources
	// that result in RRSP contribution room increase
	AllowedIncomeSources() []core.FinancialSource
	// TargetSourceForWithdrawl returns the affected income
	// source on withdrawal from an RRSP account
	TargetSourceForWithdrawl() core.FinancialSource
	// TargetSourceForContribution returns the affected deducion
	// source on contribution to an RRSP account
	TargetSourceForContribution() core.FinancialSource
	// Validate checks if the formula is valid for use
	Validate() error
	// Clone returns a copy of the formula
	Clone() Formula
}

// CalcConfig is used to pass configurations to create new RRSP calculator
type CalcConfig struct {
	Formula Formula
	TaxCalc core.TaxCalculator
}

// validate checks if the configurations are valid for use by calc constructors
func (cfg CalcConfig) validate() error {

	if cfg.Formula == nil {
		return ErrNoFormula
	}

	err := cfg.Formula.Validate()
	if err != nil {
		return errors.Wrap(err, "invalid formula")
	}

	if cfg.TaxCalc == nil {
		return ErrNoTaxCalc
	}

	return nil
}

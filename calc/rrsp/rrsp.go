// Package rrsp provides implementations for the RRSPCalculator interface
// defined in package calc
package rrsp

import (
	"github.com/malkhamis/quantax/calc/finance"
	"github.com/pkg/errors"
)

// Sentinel errors that can ben wrapped and returned by this package
var (
	ErrNoFormula  = errors.New("not formula given/set")
	ErrNoRRSPRoom = errors.New("no enough RRSP contribution room")
	ErrNoTaxCalc  = errors.New("no tax calculator given")
)

// Formula computes the max contribution room acquired for a given income
type Formula interface {
	// Contribution returns the contribution room acquired
	// given the net income
	ContributionEarned(netIncome float64) float64
	// AllowedIncomeSources return the allowed income sources
	// that result in RRSP contribution room increase
	AllowedIncomeSources() []finance.IncomeSource
	// TargetSourceForWithdrawl returns the affected income
	// source on withdrawal from an RRSP account
	TargetSourceForWithdrawl() finance.IncomeSource
	// TargetSourceForContribution returns the affected deducion
	// source on contribution to an RRSP account
	TargetSourceForContribution() finance.DeductionSource
	// Validate checks if the formula is valid for use
	Validate() error
	// Clone returns a copy of the formula
	Clone() Formula
}

package tax

import (
	"github.com/malkhamis/quantax/core"
)

// Creditor calculates tax credits that reduces payable taxes
type Creditor interface {
	// TaxCredits returns the tax credits for the given financial source
	TaxCredit(source core.FinancialSource, taxPayer *TaxPayer) float64
	// CrSourceName is the identifier/name of the credit source
	CrSourceName() string
	// Description returns a short description of the tax credit
	Description() string
	// Clone returns a copy of this creditor
	Clone() Creditor
}

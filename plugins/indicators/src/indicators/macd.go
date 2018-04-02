package indicators

import (
	"github.com/jeremyhahn/tradebot/common"
	"github.com/shopspring/decimal"
)

type MovingAverageConvergenceDivergence interface {
	Calculate(price decimal.Decimal) (decimal.Decimal, decimal.Decimal, decimal.Decimal)
	GetValue() decimal.Decimal
	GetSignalLine() decimal.Decimal
	GetHistogram() decimal.Decimal
	common.FinancialIndicator
}

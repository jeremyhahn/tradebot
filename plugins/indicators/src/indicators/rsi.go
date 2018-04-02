package indicators

import (
	"github.com/jeremyhahn/tradebot/common"
	"github.com/shopspring/decimal"
)

type RelativeStrengthIndex interface {
	IsOverSold(rsiValue decimal.Decimal) bool
	IsOverBought(rsiValue decimal.Decimal) bool
	GetValue() decimal.Decimal
	Calculate(price decimal.Decimal) decimal.Decimal
	common.FinancialIndicator
}

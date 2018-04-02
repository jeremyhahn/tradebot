package indicators

import (
	"github.com/jeremyhahn/tradebot/common"
	"github.com/shopspring/decimal"
)

type BollingerBands interface {
	GetUpper() decimal.Decimal
	GetMiddle() decimal.Decimal
	GetLower() decimal.Decimal
	StandardDeviation() decimal.Decimal
	Calculate(price decimal.Decimal) (decimal.Decimal, decimal.Decimal, decimal.Decimal)
	common.FinancialIndicator
}

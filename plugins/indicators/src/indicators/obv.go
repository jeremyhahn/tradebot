package indicators

import (
	"github.com/jeremyhahn/tradebot/common"
	"github.com/shopspring/decimal"
)

type OnBalanceVolume interface {
	GetValue() decimal.Decimal
	Calculate(price decimal.Decimal) decimal.Decimal
	common.FinancialIndicator
}

package indicators

import (
	"github.com/jeremyhahn/tradebot/common"
	"github.com/shopspring/decimal"
)

type ExponentialMovingAverage interface {
	Add(candle *common.Candlestick) decimal.Decimal
	GetCandlesticks() []common.Candlestick
	GetSize() int
	GetCount() int
	GetIndex() int
	GetAverage() decimal.Decimal
	GetPrices() []decimal.Decimal
	Sum() decimal.Decimal
	GetMultiplier() decimal.Decimal
	GetGainsAndLosses() (decimal.Decimal, decimal.Decimal)
	common.FinancialIndicator
}

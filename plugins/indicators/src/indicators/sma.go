package indicators

import (
	"github.com/jeremyhahn/tradebot/common"
	"github.com/shopspring/decimal"
)

type SimpleMovingAverage interface {
	Add(candle *common.Candlestick) decimal.Decimal
	GetCandlesticks() []common.Candlestick
	GetAverage() decimal.Decimal
	GetSize() int
	GetPrices() []decimal.Decimal
	GetGainsAndLosses() (decimal.Decimal, decimal.Decimal)
	GetCount() int
	GetIndex() int
	Sum() decimal.Decimal
	common.FinancialIndicator
}

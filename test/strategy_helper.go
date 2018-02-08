package test

import (
	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dto"
)

type StrategyTestHelper struct{}

func (h *StrategyTestHelper) CreateBalances() []common.Coin {
	return []common.Coin{
		common.Coin{
			Address:   "abc123",
			Currency:  "BTC",
			Available: 2,
			Price:     10000},
		common.Coin{
			Currency:  "USD",
			Available: 20000,
			Price:     1.00}}
}

func (h *StrategyTestHelper) CreateLastTrade() common.Trade {
	return &dto.TradeDTO{
		Id:       1,
		ChartId:  1,
		Base:     "BTC",
		Quote:    "USD",
		Exchange: "gdax",
		Type:     "sell",
		Amount:   1,
		Price:    10000}
}

func (h *StrategyTestHelper) CreateCandles() []common.Candlestick {
	var candles []common.Candlestick
	candles = append(candles, common.Candlestick{Close: 100.00})
	candles = append(candles, common.Candlestick{Close: 200.00})
	candles = append(candles, common.Candlestick{Close: 300.00})
	candles = append(candles, common.Candlestick{Close: 400.00})
	candles = append(candles, common.Candlestick{Close: 500.00})
	candles = append(candles, common.Candlestick{Close: 600.00})
	candles = append(candles, common.Candlestick{Close: 700.00})
	candles = append(candles, common.Candlestick{Close: 800.00})
	candles = append(candles, common.Candlestick{Close: 900.00})
	candles = append(candles, common.Candlestick{Close: 1000.00})
	candles = append(candles, common.Candlestick{Close: 1100.00})
	candles = append(candles, common.Candlestick{Close: 1200.00})
	candles = append(candles, common.Candlestick{Close: 1300.00})
	candles = append(candles, common.Candlestick{Close: 1400.00})
	candles = append(candles, common.Candlestick{Close: 1500.00})
	candles = append(candles, common.Candlestick{Close: 1600.00})
	candles = append(candles, common.Candlestick{Close: 1700.00})
	candles = append(candles, common.Candlestick{Close: 1800.00})
	candles = append(candles, common.Candlestick{Close: 1900.00})
	candles = append(candles, common.Candlestick{Close: 2000.00})
	candles = append(candles, common.Candlestick{Close: 2100.00})
	candles = append(candles, common.Candlestick{Close: 2200.00})
	candles = append(candles, common.Candlestick{Close: 2300.00})
	candles = append(candles, common.Candlestick{Close: 2400.00})
	candles = append(candles, common.Candlestick{Close: 2500.00})
	candles = append(candles, common.Candlestick{Close: 2600.00})
	candles = append(candles, common.Candlestick{Close: 2700.00})
	candles = append(candles, common.Candlestick{Close: 2800.00})
	candles = append(candles, common.Candlestick{Close: 2900.00})
	candles = append(candles, common.Candlestick{Close: 3000.00})
	candles = append(candles, common.Candlestick{Close: 3200.00})
	candles = append(candles, common.Candlestick{Close: 3300.00})
	candles = append(candles, common.Candlestick{Close: 3400.00})
	candles = append(candles, common.Candlestick{Close: 3500.00})
	candles = append(candles, common.Candlestick{Close: 3600.00})
	candles = append(candles, common.Candlestick{Close: 3700.00})
	candles = append(candles, common.Candlestick{Close: 3800.00})
	candles = append(candles, common.Candlestick{Close: 3900.00})
	candles = append(candles, common.Candlestick{Close: 4000.00})
	return candles
}

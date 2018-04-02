package test

import (
	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dto"
	"github.com/shopspring/decimal"
)

type StrategyTestHelper struct{}

func (h *StrategyTestHelper) CreateBalances() []common.Coin {
	return []common.Coin{
		&dto.CoinDTO{
			Address:   "abc123",
			Currency:  "BTC",
			Available: decimal.NewFromFloat(2),
			Price:     decimal.NewFromFloat(10000)},
		&dto.CoinDTO{
			Currency:  "USD",
			Available: decimal.NewFromFloat(20000),
			Price:     decimal.NewFromFloat(1.00)}}
}

func (h *StrategyTestHelper) CreateLastTrade() common.Trade {
	return &dto.TradeDTO{
		Id:       1,
		ChartId:  1,
		Base:     "BTC",
		Quote:    "USD",
		Exchange: "GDAX",
		Type:     "sell",
		Amount:   decimal.NewFromFloat(1),
		Price:    decimal.NewFromFloat(10000)}
}

func (h *StrategyTestHelper) CreateCandles() []common.Candlestick {
	var candles []common.Candlestick
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(100.00)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(200.00)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(300.00)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(400.00)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(500.00)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(600.00)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(700.00)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(800.00)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(900.00)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(1000.00)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(1100.00)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(1200.00)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(1300.00)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(1400.00)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(1500.00)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(1600.00)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(1700.00)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(1800.00)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(1900.00)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(2000.00)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(2100.00)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(2200.00)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(2300.00)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(2400.00)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(2500.00)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(2600.00)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(2700.00)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(2800.00)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(2900.00)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(3000.00)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(3200.00)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(3300.00)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(3400.00)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(3500.00)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(3600.00)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(3700.00)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(3800.00)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(3900.00)})
	candles = append(candles, common.Candlestick{Close: decimal.NewFromFloat(4000.00)})
	return candles
}

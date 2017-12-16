package main

import (
	"fmt"
	"os"
	"time"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/indicators"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	logging "github.com/op/go-logging"
)

type Trader struct {
	exchange  common.Exchange
	db        *gorm.DB
	rsi       *indicators.RSI
	logger    *logging.Logger
	LivePrice float64
	Currency  string
	chart     *common.Chart
	common.PriceListener
}

func NewTrader(db *gorm.DB, exchange common.Exchange, logger *logging.Logger) *Trader {
	return &Trader{
		exchange: exchange,
		db:       db,
		logger:   logger,
		chart:    &common.Chart{}}
}

func (trader *Trader) GetChart() *common.Chart {
	return trader.chart
}

func (trader *Trader) MakeMeRich() {

	trader.logger.Info("Starting trading bot")

	period := 900 // seconds; 15 minutes

	t := time.Now()
	year, month, day := t.Date()
	today := time.Date(year, month, day, 0, 0, 0, 0, t.Location())
	yesterday := time.Date(year, month, (day - 1), 0, 0, 0, 0, t.Location())
	end := time.Now()

	fmt.Println(today)
	fmt.Println(yesterday)
	fmt.Println(end)

	candlesticks := trader.exchange.GetTradeHistory(today, end, period)
	if len(candlesticks) < 20 {
		trader.logger.Fatal("Unable to load initial data set from exchange. Total records: ", len(candlesticks))
		os.Exit(1)
	}

	stream := NewPriceStream(period)

	// RSI
	rsiSma := indicators.NewSimpleMovingAverage(candlesticks[:14])
	rsi := indicators.NewRelativeStrengthIndex(rsiSma)
	stream.SubscribeToPeriod(rsi)

	// SMA
	bollingerSma := indicators.NewSimpleMovingAverage(candlesticks[:20])
	bollinger := indicators.NewBollingerBand(bollingerSma)
	stream.SubscribeToPeriod(bollinger)

	// MACD
	macdEma1 := indicators.NewExponentialMovingAverage(candlesticks[:12])
	macdEma2 := indicators.NewExponentialMovingAverage(candlesticks[:26])
	macd := indicators.NewMovingAverageConvergenceDivergence(macdEma1, macdEma2, 9)
	stream.SubscribeToPeriod(macd)

	gdaxPriceChan := make(chan float64)
	go trader.exchange.SubscribeToLiveFeed(gdaxPriceChan)

	for {

		price := <-gdaxPriceChan
		stream.Add(price)

		bollinger.Calculate(price)
		macd.Calculate(price)

		trader.chart.Price = price
		trader.chart.MACDValue = macd.GetValue()
		trader.chart.MACDSignal = macd.GetSignalLine()
		trader.chart.MACDHistogram = macd.GetHistogram()
		trader.chart.RSI = rsi.Calculate(price)
		trader.chart.BollingerUpper = bollinger.GetUpper()
		trader.chart.BollingerMiddle = bollinger.GetMiddle()
		trader.chart.BollingerLower = bollinger.GetMiddle()

		trader.logger.Debug("[GDAX] Price: ", trader.chart.Price,
			", MACD_VALUE: ", trader.chart.MACDValue,
			", MACD_HISTOGRAM: ", trader.chart.MACDHistogram,
			", MACD_SIGNAL: ", trader.chart.MACDSignal,
			", RSI: ", trader.chart.RSI,
			", Bollinger Upper: ", trader.chart.BollingerUpper,
			", Bollinger Middle: ", trader.chart.BollingerMiddle,
			", Bollinger Lower: ", trader.chart.BollingerLower)
	}

}

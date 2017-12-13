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
	common.PriceListener
}

func NewTrader(db *gorm.DB, exchange common.Exchange, logger *logging.Logger) *Trader {
	//	var candlesticks []common.Candlestick

	return &Trader{
		exchange: exchange,
		db:       db,
		logger:   logger}
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
	/*for _, c := range candlesticks {
		fmt.Printf("%+v\n", c)
	}*/

	stream := NewPriceStream(period)

	// RSI
	ema := indicators.NewExponentialMovingAverage(candlesticks[:30])
	rsi := indicators.NewRelativeStrengthIndex(ema)
	stream.Subscribe(ema)

	// SMA
	sma := indicators.NewSimpleMovingAverage(candlesticks[:20])
	bollinger := indicators.NewBollingerBand(sma)
	stream.Subscribe(sma)

	gdaxPriceChan := make(chan float64)
	go trader.exchange.SubscribeToLiveFeed(gdaxPriceChan)

	for {
		price := <-gdaxPriceChan
		stream.Add(price)
		bollinger.Calculate(price)
		trader.logger.Debug("[GDAX] Price: ", price, ", RSI: ", rsi.Calculate(price),
			", Bollinger Upper: ", bollinger.GetUpper(), ", Bollinger Middle: ", bollinger.GetMiddle(),
			", Bollinger Lower: ", bollinger.GetLower())
	}

}

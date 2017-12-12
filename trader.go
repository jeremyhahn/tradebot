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

	sma := indicators.NewSimpleMovingAverage(candlesticks[:14])
	sma2 := indicators.NewSimpleMovingAverage(candlesticks[:20])
	rsi := indicators.NewRelativeStrengthIndex(sma)
	bollinger := indicators.NewBollingerBand(sma2)

	gdaxPriceChan := make(chan float64)
	go trader.exchange.SubscribeToLiveFeed(gdaxPriceChan)

	stream := NewPriceStream(period)
	stream.Subscribe(sma)
	stream.Subscribe(sma2)

	for {
		price := <-gdaxPriceChan

		stream.Add(price)

		bollinger.Calculate(price)

		trader.logger.Debug("[GDAX] Price: ", price, ", RSI: ", rsi.Calculate(price),
			", Bollinger Upper: ", bollinger.GetUpper(), ", Bollinger Middle: ", bollinger.GetMiddle(),
			", Bollinger Lower: ", bollinger.GetLower())

		/*
			if rsi.RecommendBuy() {
				fmt.Println("** RSI Recommended BUY Time! ", rsi.Calculate(price), " **")
			}
			if rsi.RecommendSell() {
				fmt.Println("** RSI Recommended SELL Time! ", rsi.Calculate(price), " **")
			}*/
	}

}
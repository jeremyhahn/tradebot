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

type Chart struct {
	exchange    common.Exchange
	db          *gorm.DB
	rsi         *indicators.RSI
	logger      *logging.Logger
	priceStream *PriceStream
	data        *common.ChartData
}

type PriceStream struct {
	period          int       // total number of seconds per candlestick
	start           time.Time // when the first price was added to the buffer
	volume          int
	buffer          []float64
	priceListeners  []common.PriceListener
	periodListeners []common.PeriodListener
}

func NewChartMock() *Chart {
	return &Chart{}
}

func NewChart(db *gorm.DB, exchange common.Exchange, logger *logging.Logger) *Chart {
	return &Chart{
		exchange: exchange,
		db:       db,
		logger:   logger,
		data:     &common.ChartData{}}
}

func (chart *Chart) GetChartData() *common.ChartData {
	return chart.data
}

func (chart *Chart) StreamData(ws *WebsocketServer) {

	chart.logger.Info("Starting trading bot")

	period := 900 // seconds; 15 minutes

	t := time.Now()
	year, month, day := t.Date()
	today := time.Date(year, month, day, 0, 0, 0, 0, t.Location())
	yesterday := time.Date(year, month, (day - 1), 0, 0, 0, 0, t.Location())
	end := time.Now()

	fmt.Println(today)
	fmt.Println(yesterday)
	fmt.Println(end)

	candlesticks := chart.exchange.GetTradeHistory(today, end, period)
	if len(candlesticks) < 20 {
		chart.logger.Fatal("Unable to load initial data set from exchange. Total records: ", len(candlesticks))
		os.Exit(1)
	}

	fmt.Printf("%+v\n", candlesticks)

	chart.priceStream = NewPriceStream(period)

	// RSI
	rsiSma := indicators.NewSimpleMovingAverage(candlesticks[:14])
	rsi := indicators.NewRelativeStrengthIndex(rsiSma)
	chart.priceStream.SubscribeToPeriod(rsi)

	// SMA
	bollingerSma := indicators.NewSimpleMovingAverage(candlesticks[:20])
	bollinger := indicators.NewBollingerBand(bollingerSma)
	chart.priceStream.SubscribeToPeriod(bollinger)

	// MACD
	macdEma1 := indicators.NewExponentialMovingAverage(candlesticks[:12])
	macdEma2 := indicators.NewExponentialMovingAverage(candlesticks[:26])
	macd := indicators.NewMovingAverageConvergenceDivergence(macdEma1, macdEma2, 9)
	chart.priceStream.SubscribeToPeriod(macd)

	gdaxPriceChan := make(chan float64)
	go chart.exchange.SubscribeToLiveFeed(gdaxPriceChan)

	for {
		price := <-gdaxPriceChan
		chart.priceStream.Add(price)
		ws.Broadcast(price)

		bollinger.Calculate(price)
		macdValue, macdSignal, macdHistogram := macd.Calculate(price)

		chart.data.Currency = chart.exchange.GetCurrency()
		chart.data.Price = price
		chart.data.MACDValue = macd.GetValue()
		chart.data.MACDSignal = macd.GetSignalLine()
		chart.data.MACDHistogram = macd.GetHistogram()
		chart.data.MACDValueLive = macdValue
		chart.data.MACDSignalLive = macdSignal
		chart.data.MACDHistogramLive = macdHistogram
		chart.data.RSI = rsi.GetValue()
		chart.data.RSILive = rsi.Calculate(price)
		chart.data.BollingerUpper = bollinger.GetUpper()
		chart.data.BollingerMiddle = bollinger.GetMiddle()
		chart.data.BollingerLower = bollinger.GetLower()
	}
}

package main

import (
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
	period      int
}

type PriceStream struct {
	Period          int       // total number of seconds per candlestick
	Start           time.Time // when the first price was added to the buffer
	Volume          int
	buffer          []float64
	priceListeners  []common.PriceListener
	periodListeners []common.PeriodListener
}

func NewChartMock() *Chart {
	return &Chart{}
}

func NewChart(db *gorm.DB, exchange common.Exchange, logger *logging.Logger, priceStream *PriceStream) *Chart {
	return &Chart{
		exchange:    exchange,
		db:          db,
		logger:      logger,
		data:        &common.ChartData{},
		period:      priceStream.Period,
		priceStream: priceStream}
}

func (chart *Chart) GetChartData() *common.ChartData {
	return chart.data
}

func (chart *Chart) Stream(ws *WebsocketServer) {

	chart.logger.Infof("Streaming %s chart data", chart.data.Currency)

	t := time.Now()
	year, month, day := t.Date()
	today := time.Date(year, month, day, 0, 0, 0, 0, t.Location())
	//yesterday := time.Date(year, month, (day - 1), 0, 0, 0, 0, t.Location())
	now := time.Now()

	chart.logger.Debugf("Getting trade history from %s - %s ", today, now)

	candlesticks := chart.exchange.GetTradeHistory(today, now, chart.period)
	if len(candlesticks) < 20 {
		chart.logger.Fatal("Unable to load initial data set from exchange. Total records: ", len(candlesticks))
	}

	//chart.priceStream = NewPriceStream(chart.period)

	// RSI
	rsiSma := indicators.NewSimpleMovingAverage(candlesticks[:14])
	rsi := indicators.NewRelativeStrengthIndex(rsiSma)
	chart.priceStream.SubscribeToPeriod(rsi)

	// Bollinger Band
	bollingerSma := indicators.NewSimpleMovingAverage(candlesticks[:20])
	bollinger := indicators.NewBollingerBand(bollingerSma)
	chart.priceStream.SubscribeToPeriod(bollinger)

	// MACD
	macdEma1 := indicators.NewExponentialMovingAverage(candlesticks[:12])
	macdEma2 := indicators.NewExponentialMovingAverage(candlesticks[:26])
	macd := indicators.NewMovingAverageConvergenceDivergence(macdEma1, macdEma2, 9)
	chart.priceStream.SubscribeToPeriod(macd)

	// Pre-warm indicators
	for _, c := range candlesticks {
		rsi.OnPeriodChange(&c)
		bollinger.OnPeriodChange(&c)
		macd.OnPeriodChange(&c)
	}

	priceChan := make(chan common.PriceChange)
	go chart.exchange.SubscribeToLiveFeed(priceChan)

	for {
		priceChannel := <-priceChan
		satoshis := priceChannel.Satoshis

		chart.priceStream.Add(priceChannel.Price)

		bollinger.Calculate(satoshis)
		macdValue, macdSignal, macdHistogram := macd.Calculate(satoshis)

		chart.data.Currency = chart.exchange.GetCurrency()
		chart.data.Price = priceChannel.Price
		chart.data.Satoshis = priceChannel.Satoshis
		chart.data.MACDValue = macd.GetValue()
		chart.data.MACDSignal = macd.GetSignalLine()
		chart.data.MACDHistogram = macd.GetHistogram()
		chart.data.MACDValueLive = macdValue
		chart.data.MACDSignalLive = macdSignal
		chart.data.MACDHistogramLive = macdHistogram
		chart.data.RSI = rsi.GetValue()
		chart.data.RSILive = rsi.Calculate(satoshis)
		chart.data.BollingerUpper = bollinger.GetUpper()
		chart.data.BollingerMiddle = bollinger.GetMiddle()
		chart.data.BollingerLower = bollinger.GetLower()

		//ws.Broadcast(chart.data)
	}
}

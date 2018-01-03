package service

import (
	"time"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/indicators"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	logging "github.com/op/go-logging"
)

type Chart struct {
	exchange      common.Exchange
	db            *gorm.DB
	rsi           *indicators.RSI
	logger        *logging.Logger
	priceStream   *PriceStream
	priceChannel  chan common.PriceChange
	candleChannel chan common.Candlestick
	chartChannel  chan *common.ChartData
	data          *common.ChartData
	period        int
}

func NewChartMock() *Chart {
	return &Chart{}
}

func NewChart(ctx *common.Context, exchange common.Exchange, period int) *Chart {
	return &Chart{
		exchange:    exchange,
		db:          ctx.DB,
		logger:      ctx.Logger,
		data:        &common.ChartData{},
		period:      period,
		priceStream: NewPriceStream(period)}
}

func (chart *Chart) GetChartData() *common.ChartData {
	return chart.data
}

func (chart *Chart) Stream() {

	chart.logger.Infof("[Chart.Stream] Streaming %s %s chart data.",
		chart.exchange.GetName(), chart.exchange.FormattedCurrencyPair())

	t := time.Now()
	year, month, day := t.Date()
	today := time.Date(year, month, day, 0, 0, 0, 0, t.Location())
	//yesterday := time.Date(year, month, (day - 1), 0, 0, 0, 0, t.Location())
	now := time.Now()

	chart.logger.Debugf("[Chart.Stream] Getting %s %s trade history from %s - %s ",
		chart.exchange.GetName(), chart.exchange.FormattedCurrencyPair(), today, now)

	candlesticks := chart.exchange.GetTradeHistory(today, now, chart.period)
	if len(candlesticks) < 20 {
		chart.logger.Errorf("[Chart.Stream] Failed to load initial candlesticks from %s. Total records: %d",
			chart.exchange.GetName(), len(candlesticks))
		return
	}

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

	priceChange := make(chan common.PriceChange)
	go chart.exchange.SubscribeToLiveFeed(priceChange)

	for {
		newPrice := <-priceChange
		satoshis := newPrice.Satoshis

		bollinger.Calculate(satoshis)
		macdValue, macdSignal, macdHistogram := macd.Calculate(satoshis)

		chart.data.CurrencyPair = chart.exchange.GetCurrencyPair()
		chart.data.Price = newPrice.Price
		chart.data.Satoshis = newPrice.Satoshis
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

		chart.logger.Debugf("[Chart.Stream] ChartData: %+v\n", chart.data)
	}
}

func (chart *Chart) OnPeriodChange(candlestick *common.Candlestick) {
	chart.logger.Debugf("[Chart.OnPriceChange] candlestick: %+v\n", candlestick)
}

func (chart *Chart) OnPriceChange(priceChange common.PriceChange) {
	chart.logger.Debugf("[Chart.OnPriceChange] priceChange: %+v\n", priceChange)
}

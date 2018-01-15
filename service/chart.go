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
	obv           *indicators.OBV
	rsi           *indicators.RSI
	bband         *indicators.Bollinger
	macd          *indicators.MACD
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

func (chart *Chart) loadCandlesticks() []common.Candlestick {

	t := time.Now()
	year, month, day := t.Date()
	yesterday := time.Date(year, month, day-1, 0, 0, 0, 0, t.Location())
	now := time.Now()

	chart.logger.Debugf("[Chart.Stream] Getting %s %s trade history from %s - %s ",
		chart.exchange.GetName(), chart.exchange.FormattedCurrencyPair(), yesterday, now)

	candlesticks := chart.exchange.GetTradeHistory(yesterday, now, chart.period)
	if len(candlesticks) < 20 {
		chart.logger.Errorf("[Chart.Stream] Failed to load initial candlesticks from %s. Total records: %d",
			chart.exchange.GetName(), len(candlesticks))
		return nil
	}

	return candlesticks
}

func (chart *Chart) Stream() {

	chart.logger.Infof("[Chart.Stream] Streaming %s %s chart data.",
		chart.exchange.GetName(), chart.exchange.FormattedCurrencyPair())

	candles := chart.loadCandlesticks()
	chart.logger.Debugf("[Chart.Stream] Prewarming indicators with %d candlesticks", len(candles))
	candlesticks := chart.reverseCandles(candles)

	rsiSma := indicators.NewSimpleMovingAverage(candlesticks[:14])
	chart.rsi = indicators.NewRelativeStrengthIndex(rsiSma)
	for _, c := range candlesticks[14:] {
		chart.rsi.OnPeriodChange(&c)
	}
	chart.priceStream.SubscribeToPeriod(chart.rsi)

	bollingerSma := indicators.NewSimpleMovingAverage(candlesticks[:20])
	chart.bband = indicators.NewBollingerBand(bollingerSma)
	for _, c := range candlesticks[20:] {
		chart.bband.OnPeriodChange(&c)
	}
	chart.priceStream.SubscribeToPeriod(chart.bband)

	macdEma1 := indicators.NewExponentialMovingAverage(candlesticks[:10])
	macdEma2 := indicators.NewExponentialMovingAverage(candlesticks[:26])
	for _, c := range candlesticks[10:26] {
		macdEma1.OnPeriodChange(&c)
	}
	chart.macd = indicators.NewMovingAverageConvergenceDivergence(macdEma1, macdEma2, 9)
	for _, c := range candlesticks[26:] {
		chart.macd.OnPeriodChange(&c)
	}
	chart.priceStream.SubscribeToPeriod(chart.macd)

	chart.obv = indicators.NewOnBalanceVolume(candlesticks)
	chart.priceStream.SubscribeToPeriod(chart.obv)

	priceChange := make(chan common.PriceChange)
	go chart.exchange.SubscribeToLiveFeed(priceChange)

	chart.priceStream.SubscribeToPrice(chart)
	chart.priceStream.SubscribeToPeriod(chart)

	for {
		chart.priceStream.Listen(priceChange)
	}
}

func (chart *Chart) OnPeriodChange(candle *common.Candlestick) {
	chart.rsi.OnPeriodChange(candle)
	chart.bband.OnPeriodChange(candle)
	chart.macd.OnPeriodChange(candle)
}

func (chart *Chart) OnPriceChange(newPrice *common.PriceChange) {
	bUpper, bMiddle, bLower := chart.bband.Calculate(newPrice.Price)
	macdValue, macdSignal, macdHistogram := chart.macd.Calculate(newPrice.Price)
	chart.data.CurrencyPair = chart.exchange.GetCurrencyPair()
	chart.data.Price = newPrice.Price
	chart.data.Satoshis = newPrice.Satoshis
	chart.data.RSI = chart.rsi.GetValue()
	chart.data.RSILive = chart.rsi.Calculate(newPrice.Price)
	chart.data.BollingerUpper = chart.bband.GetUpper()
	chart.data.BollingerMiddle = chart.bband.GetMiddle()
	chart.data.BollingerLower = chart.bband.GetLower()
	chart.data.BollingerUpperLive = bUpper
	chart.data.BollingerMiddleLive = bMiddle
	chart.data.BollingerLowerLive = bLower
	chart.data.MACDValue = chart.macd.GetValue()
	chart.data.MACDSignal = chart.macd.GetSignalLine()
	chart.data.MACDHistogram = chart.macd.GetHistogram()
	chart.data.MACDValueLive = macdValue
	chart.data.MACDSignalLive = macdSignal
	chart.data.MACDHistogramLive = macdHistogram
	chart.data.OnBalanceVolume = chart.obv.GetValue()
	chart.data.OnBalanceVolumeLive = chart.obv.Calculate(newPrice.Price)
	//bytes, _ := json.MarshalIndent(chart.data, "", "    ")
	chart.logger.Debugf("[Chart.OnPriceChange] ChartData: %+v\n", chart.data)
}

func (chart *Chart) reverseCandles(candles []common.Candlestick) []common.Candlestick {
	var newCandles []common.Candlestick
	for i := len(candles) - 1; i > 0; i-- {
		newCandles = append(newCandles, candles[i])
	}
	return newCandles
}

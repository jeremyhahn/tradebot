package service

import (
	"time"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/indicators"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type TradingStrategy interface {
	OnPriceChange(chart common.ChartService)
}

type ChartServiceImpl struct {
	ctx           *common.Context
	priceStream   *PriceStream
	priceChannel  chan common.PriceChange
	candleChannel chan common.Candlestick
	chartChannel  chan *common.ChartData
	period        int
	strategy      TradingStrategy
	Exchange      common.Exchange
	OBV           *indicators.OBV
	RSI           *indicators.RSI
	Bband         *indicators.Bollinger
	MACD          *indicators.MACD
	Data          *common.ChartData
	common.ChartService
}

func NewChartService(ctx *common.Context, exchange common.Exchange, strategy TradingStrategy, period int) common.ChartService {
	return &ChartServiceImpl{
		Exchange:    exchange,
		strategy:    strategy,
		ctx:         ctx,
		Data:        &common.ChartData{},
		period:      period,
		priceStream: NewPriceStream(period)}
}

func (chart *ChartServiceImpl) GetExchange() common.Exchange {
	return chart.Exchange
}

func (chart *ChartServiceImpl) GetData() *common.ChartData {
	return chart.Data
}

func (chart *ChartServiceImpl) GetCurrencyPair() common.CurrencyPair {
	return chart.Exchange.GetCurrencyPair()
}

func (chart *ChartServiceImpl) loadCandlesticks() []common.Candlestick {

	t := time.Now()
	year, month, day := t.Date()
	yesterday := time.Date(year, month, day-1, 0, 0, 0, 0, t.Location())
	now := time.Now()

	chart.ctx.Logger.Debugf("[ChartService.Stream] Getting %s %s trade history from %s - %s ",
		chart.Exchange.GetName(), chart.Exchange.FormattedCurrencyPair(), yesterday, now)

	candlesticks := chart.Exchange.GetTradeHistory(yesterday, now, chart.period)
	if len(candlesticks) < 20 {
		chart.ctx.Logger.Errorf("[ChartService.Stream] Failed to load initial candlesticks from %s. Total records: %d",
			chart.Exchange.GetName(), len(candlesticks))
		return nil
	}

	return candlesticks
}

func (chart *ChartServiceImpl) Stream() {

	chart.ctx.Logger.Infof("[ChartService.Stream] Streaming %s %s chart data.",
		chart.Exchange.GetName(), chart.Exchange.FormattedCurrencyPair())

	candles := chart.loadCandlesticks()
	chart.ctx.Logger.Debugf("[ChartService.Stream] Prewarming indicators with %d candlesticks", len(candles))
	candlesticks := chart.reverseCandles(candles)

	rsiSma := indicators.NewSimpleMovingAverage(candlesticks[:14])
	chart.RSI = indicators.NewRelativeStrengthIndex(rsiSma)
	for _, c := range candlesticks[14:] {
		chart.RSI.OnPeriodChange(&c)
	}
	chart.priceStream.SubscribeToPeriod(chart.RSI)

	bollingerSma := indicators.NewSimpleMovingAverage(candlesticks[:20])
	chart.Bband = indicators.NewBollingerBand(bollingerSma)
	for _, c := range candlesticks[20:] {
		chart.Bband.OnPeriodChange(&c)
	}
	chart.priceStream.SubscribeToPeriod(chart.Bband)

	macdEma1 := indicators.NewExponentialMovingAverage(candlesticks[:10])
	macdEma2 := indicators.NewExponentialMovingAverage(candlesticks[:26])
	for _, c := range candlesticks[10:26] {
		macdEma1.OnPeriodChange(&c)
	}
	chart.MACD = indicators.NewMovingAverageConvergenceDivergence(macdEma1, macdEma2, 9)
	for _, c := range candlesticks[26:] {
		chart.MACD.OnPeriodChange(&c)
	}
	chart.priceStream.SubscribeToPeriod(chart.MACD)

	chart.OBV = indicators.NewOnBalanceVolume(candlesticks)
	chart.priceStream.SubscribeToPeriod(chart.OBV)

	priceChange := make(chan common.PriceChange)
	go chart.Exchange.SubscribeToLiveFeed(priceChange)

	chart.priceStream.SubscribeToPrice(chart)
	chart.priceStream.SubscribeToPeriod(chart)

	for {
		chart.priceStream.Listen(priceChange)
	}
}

func (chart *ChartServiceImpl) OnPeriodChange(candle *common.Candlestick) {
	chart.RSI.OnPeriodChange(candle)
	chart.Bband.OnPeriodChange(candle)
	chart.MACD.OnPeriodChange(candle)
}

func (chart *ChartServiceImpl) OnPriceChange(newPrice *common.PriceChange) {
	bUpper, bMiddle, bLower := chart.Bband.Calculate(newPrice.Price)
	macdValue, macdSignal, macdHistogram := chart.MACD.Calculate(newPrice.Price)
	chart.Data.CurrencyPair = chart.Exchange.GetCurrencyPair()
	chart.Data.Price = newPrice.Price
	chart.Data.Exchange = chart.Exchange.GetName()
	chart.Data.Satoshis = newPrice.Satoshis
	chart.Data.RSI = chart.RSI.GetValue()
	chart.Data.RSILive = chart.RSI.Calculate(newPrice.Price)
	chart.Data.BollingerUpper = chart.Bband.GetUpper()
	chart.Data.BollingerMiddle = chart.Bband.GetMiddle()
	chart.Data.BollingerLower = chart.Bband.GetLower()
	chart.Data.BollingerUpperLive = bUpper
	chart.Data.BollingerMiddleLive = bMiddle
	chart.Data.BollingerLowerLive = bLower
	chart.Data.MACDValue = chart.MACD.GetValue()
	chart.Data.MACDSignal = chart.MACD.GetSignalLine()
	chart.Data.MACDHistogram = chart.MACD.GetHistogram()
	chart.Data.MACDValueLive = macdValue
	chart.Data.MACDSignalLive = macdSignal
	chart.Data.MACDHistogramLive = macdHistogram
	chart.Data.OnBalanceVolume = chart.OBV.GetValue()
	chart.Data.OnBalanceVolumeLive = chart.OBV.Calculate(newPrice.Price)
	//bytes, _ := json.MarshalIndent(chart.Data, "", "    ")
	//chart.ctx.Logger.Debugf("[ChartService.OnPriceChange] ChartData: %+v\n", chart.Data)
	if chart.strategy != nil {
		chart.strategy.OnPriceChange(chart)
	}
}

func (chart *ChartServiceImpl) reverseCandles(candles []common.Candlestick) []common.Candlestick {
	var newCandles []common.Candlestick
	for i := len(candles) - 1; i > 0; i-- {
		newCandles = append(newCandles, candles[i])
	}
	return newCandles
}

package service

import (
	"time"

	"github.com/jeremyhahn/tradebot/common"
)

type PriceStream interface {
	Listen(priceChange chan common.PriceChange) common.PriceChange
	SubscribeToPrice(listener common.PriceListener)
	SubscribeToPeriod(listener common.PeriodListener)
}

type PriceStreamImpl struct {
	Period          int       // total seconds per candlestick
	Start           time.Time // timestamp of first candlestick
	Volume          int
	buffer          []float64
	priceListeners  []common.PriceListener
	periodListeners []common.PeriodListener
}

func NewPriceStream(period int) PriceStream {
	return &PriceStreamImpl{
		Period:          period,
		Start:           common.NewCandlestickPeriod(period),
		buffer:          make([]float64, 0),
		priceListeners:  make([]common.PriceListener, 0),
		periodListeners: make([]common.PeriodListener, 0)}
}

func (ps *PriceStreamImpl) Listen(priceChange chan common.PriceChange) common.PriceChange {
	newPrice := <-priceChange
	ps.buffer = append(ps.buffer, newPrice.Price)
	ps.Volume = ps.Volume + 1
	ps.notifyPriceListeners(&newPrice)
	if time.Since(ps.Start).Seconds() >= float64(ps.Period) {
		candlestick := common.CreateCandlestick(newPrice.Exchange, newPrice.CurrencyPair, ps.Period, ps.buffer)
		ps.notifyPeriodListeners(candlestick)
		ps.Volume = 0
		ps.Start = common.NewCandlestickPeriod(ps.Period)
		ps.buffer = ps.buffer[:0]
	}
	return newPrice
}

func (ps *PriceStreamImpl) SubscribeToPrice(listener common.PriceListener) {
	ps.priceListeners = append(ps.priceListeners, listener)
}

func (ps *PriceStreamImpl) SubscribeToPeriod(listener common.PeriodListener) {
	ps.periodListeners = append(ps.periodListeners, listener)
}

func (ps *PriceStreamImpl) notifyPeriodListeners(candlestick *common.Candlestick) {
	for _, listener := range ps.periodListeners {
		go listener.OnPeriodChange(candlestick)
	}
}

func (ps *PriceStreamImpl) notifyPriceListeners(priceChange *common.PriceChange) {
	for _, listener := range ps.priceListeners {
		go listener.OnPriceChange(priceChange)
	}
}

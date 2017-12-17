package main

import (
	"time"

	"github.com/jeremyhahn/tradebot/common"
)

func NewPriceStream(period int) *PriceStream {
	return &PriceStream{
		period:          period,
		start:           common.NewCandlestickPeriod(period),
		buffer:          make([]float64, 0),
		priceListeners:  make([]common.PriceListener, 0),
		periodListeners: make([]common.PeriodListener, 0)}
}

func (ps *PriceStream) Add(price float64) {
	ps.buffer = append(ps.buffer, price)
	ps.volume = ps.volume + 1
	ps.notifyPriceListeners(price)
	if time.Since(ps.start).Seconds() >= float64(ps.period) {
		candlestick := common.CreateCandlestick(ps.period, ps.buffer)
		ps.notifyPeriodListeners(candlestick)
		ps.volume = 0
		ps.start = common.NewCandlestickPeriod(ps.period)
		ps.buffer = ps.buffer[:0]
	}
}

func (ps *PriceStream) SubscribeToPrice(listener common.PriceListener) {
	ps.priceListeners = append(ps.priceListeners, listener)
}

func (ps *PriceStream) SubscribeToPeriod(listener common.PeriodListener) {
	ps.periodListeners = append(ps.periodListeners, listener)
}

func (ps *PriceStream) notifyPeriodListeners(candlestick *common.Candlestick) {
	for _, listener := range ps.periodListeners {
		listener.OnPeriodChange(candlestick)
	}
}

func (ps *PriceStream) notifyPriceListeners(price float64) {
	for _, listener := range ps.priceListeners {
		listener.OnPriceChange(price)
	}
}

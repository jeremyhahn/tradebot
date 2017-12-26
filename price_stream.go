package main

import (
	"time"

	"github.com/jeremyhahn/tradebot/common"
)

func NewPriceStream(period int) *PriceStream {
	return &PriceStream{
		Period:          period,
		Start:           common.NewCandlestickPeriod(period),
		buffer:          make([]float64, 0),
		priceListeners:  make([]common.PriceListener, 0),
		periodListeners: make([]common.PeriodListener, 0)}
}

func (ps *PriceStream) Add(price float64) {
	ps.buffer = append(ps.buffer, price)
	ps.Volume = ps.Volume + 1
	ps.notifyPriceListeners(price)
	if time.Since(ps.Start).Seconds() >= float64(ps.Period) {
		candlestick := common.CreateCandlestick(ps.Period, ps.buffer)
		ps.notifyPeriodListeners(candlestick)
		ps.Volume = 0
		ps.Start = common.NewCandlestickPeriod(ps.Period)
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

package main

import (
	"time"

	"github.com/jeremyhahn/tradebot/common"
)

type PriceStream struct {
	period    int       // total number of seconds per candlestick
	start     time.Time // when the first price was added to the buffer
	volume    int
	buffer    []float64
	listeners []common.PriceListener
}

func NewPriceStream(period int) *PriceStream {
	return &PriceStream{
		period:    period,
		start:     common.NewCandlestickPeriod(period),
		buffer:    make([]float64, 0),
		listeners: make([]common.PriceListener, 0)}
}

func (ps *PriceStream) Add(price float64) {
	ps.buffer = append(ps.buffer, price)
	ps.volume = ps.volume + 1
	ps.notifyPrice(price)
	if time.Since(ps.start).Seconds() >= float64(ps.period) {
		candlestick := common.CreateCandlestick(ps.period, ps.buffer)
		ps.notifyCandlestick(candlestick)
		ps.volume = 0
		ps.start = common.NewCandlestickPeriod(ps.period)
		ps.buffer = ps.buffer[:0]
	}
}

func (ps *PriceStream) Subscribe(listener common.PriceListener) {
	ps.listeners = append(ps.listeners, listener)
}

func (ps *PriceStream) notifyCandlestick(candlestick *common.Candlestick) {
	for _, listener := range ps.listeners {
		listener.OnCandlestickCreated(candlestick)
	}
}

func (ps *PriceStream) notifyPrice(price float64) {
	for _, listener := range ps.listeners {
		listener.OnPriceChange(price)
	}
}

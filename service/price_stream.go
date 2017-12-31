package service

import (
	"time"

	"github.com/jeremyhahn/tradebot/common"
)

type PriceStream struct {
	Period          int       // total number of seconds per candlestick
	Start           time.Time // when the first price was added to the buffer
	Volume          int
	buffer          []float64
	priceListeners  []common.PriceListener
	periodListeners []common.PeriodListener
}

func NewPriceStream(period int) *PriceStream {
	return &PriceStream{
		Period:          period,
		Start:           common.NewCandlestickPeriod(period),
		buffer:          make([]float64, 0),
		priceListeners:  make([]common.PriceListener, 0),
		periodListeners: make([]common.PeriodListener, 0)}
}

func (ps *PriceStream) Listen(priceChange chan common.PriceChange) {
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
}

func (ps *PriceStream) SubscribeToPrice(listener common.PriceListener) {
	ps.priceListeners = append(ps.priceListeners, listener)
}

func (ps *PriceStream) SubscribeToPeriod(listener common.PeriodListener) {
	ps.periodListeners = append(ps.periodListeners, listener)
}

func (ps *PriceStream) notifyPeriodListeners(candlestick *common.Candlestick) {
	for _, listener := range ps.periodListeners {
		go listener.OnPeriodChange(candlestick)
	}
}

func (ps *PriceStream) notifyPriceListeners(priceChange *common.PriceChange) {
	for _, listener := range ps.priceListeners {
		go listener.OnPriceChange(priceChange)
	}
}

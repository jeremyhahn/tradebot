package indicators

import "github.com/jeremyhahn/tradebot/common"

type MACD struct {
	ema  common.MovingAverage
	ema2 common.MovingAverage
}

func NewMovingAverageConvergenceDivergence(ema, ema2 common.MovingAverage) *MACD {
	return &MACD{
		ema:  ema,
		ema2: ema2}
}

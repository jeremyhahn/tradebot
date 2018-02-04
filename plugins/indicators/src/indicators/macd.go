package indicators

import "github.com/jeremyhahn/tradebot/common"

type MovingAverageConvergenceDivergence interface {
	Calculate(price float64) (float64, float64, float64)
	GetValue() float64
	GetSignalLine() float64
	GetHistogram() float64
	common.FinancialIndicator
}

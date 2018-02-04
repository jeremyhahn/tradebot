package indicators

import "github.com/jeremyhahn/tradebot/common"

type BollingerBands interface {
	GetUpper() float64
	GetMiddle() float64
	GetLower() float64
	StandardDeviation() float64
	Calculate(price float64) (float64, float64, float64)
	common.FinancialIndicator
}

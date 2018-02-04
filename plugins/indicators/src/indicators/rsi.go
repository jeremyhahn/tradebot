package indicators

import "github.com/jeremyhahn/tradebot/common"

type RelativeStrengthIndex interface {
	IsOverSold(rsiValue float64) bool
	IsOverBought(rsiValue float64) bool
	GetValue() float64
	Calculate(price float64) float64
	common.FinancialIndicator
}

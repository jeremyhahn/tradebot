package indicators

import "github.com/jeremyhahn/tradebot/common"

type OnBalanceVolume interface {
	GetValue() float64
	Calculate(price float64) float64
	common.FinancialIndicator
}

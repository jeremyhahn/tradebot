package example

import "github.com/jeremyhahn/tradebot/common"

type ExampleIndicator interface {
	Calculate(price float64) float64
	common.FinancialIndicator
}

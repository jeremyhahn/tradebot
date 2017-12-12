package strategy

import "github.com/jeremyhahn/tradebot/common"

type FifteenMinuteStrategy struct {
	period     int
	indicators []connom.Indicator
	common.Strategy
}

func CreateFifteenMinuteStrategy() *FifteenMinuteStrategy {

}

func IsTimeToBuy() bool {
	return false
}

func IsTimeToSell() bool {
	return false
}

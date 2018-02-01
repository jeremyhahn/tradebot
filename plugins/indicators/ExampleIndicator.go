package main

import (
	"fmt"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/plugins/indicators/example"
)

type ExampleIndicatorImpl struct {
	name        string
	displayName string
	params      ExampleIndicatorParams
	example.ExampleIndicator
}

type ExampleIndicatorParams struct {
	param1 int
	param2 int
	param3 int
}

func main() {
}

func CreateExampleIndicator(candles []common.Candlestick, params []string) common.FinancialIndicator {
	return ExampleIndicatorImpl{
		name:        "ExampleIndicator",
		displayName: "Example Indicator®",
		params: ExampleIndicatorParams{
			param1: 1,
			param2: 2,
			param3: 3}}
}

func (fi ExampleIndicatorImpl) GetName() string {
	return fi.name
}

func (fi ExampleIndicatorImpl) GetDisplayName() string {
	return fi.displayName
}

func (fi ExampleIndicatorImpl) GetDefaultParameters() []string {
	return []string{"1", "2", "3"}
}

func (fi ExampleIndicatorImpl) GetParameters() []string {
	return []string{
		fmt.Sprintf("%d", fi.params.param1),
		fmt.Sprintf("%d", fi.params.param2),
		fmt.Sprintf("%d", fi.params.param3)}
}

func (fi ExampleIndicatorImpl) Calculate(price float64) float64 {
	price++
	return price
}

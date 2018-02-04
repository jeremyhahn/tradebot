package main

import (
	"fmt"
	"strconv"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/plugins/indicators/src/example"
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

func CreateExampleIndicator(candles []common.Candlestick, params []string) (common.FinancialIndicator, error) {
	if params == nil {
		temp := &ExampleIndicatorImpl{}
		params = temp.GetDefaultParameters()
	}
	intParam1, _ := strconv.ParseInt(params[0], 10, 64)
	intParam2, _ := strconv.ParseInt(params[1], 10, 64)
	intParam3, _ := strconv.ParseInt(params[2], 10, 64)
	return ExampleIndicatorImpl{
		name:        "ExampleIndicator",
		displayName: "Example Indicator®",
		params: ExampleIndicatorParams{
			param1: int(intParam1),
			param2: int(intParam2),
			param3: int(intParam3)}}, nil
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

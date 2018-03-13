package util

import (
	"errors"
	"fmt"
	"time"

	"github.com/jeremyhahn/tradebot/common"
	logging "github.com/op/go-logging"
)

func ReverseCandlesticks(candles []common.Candlestick) []common.Candlestick {
	var reversed []common.Candlestick
	for i := len(candles) - 1; i > 0; i-- {
		reversed = append(reversed, candles[i])
	}
	return reversed
}

func FindClosesttDatedCandle(logger *logging.Logger, needle time.Time, haystack []common.Candlestick) (*common.Candlestick, error) {
	var selectedCandle *common.Candlestick
	var lastCandle *common.Candlestick
	var finalCandle *common.Candlestick
	for _, candle := range haystack {
		finalCandle = &candle
		logger.Debugf("[util.FindClosesttDatedCandle] Comparing candle date %s with target date %s", candle.Date, needle)
		if candle.Date.After(needle) || needle.Equal(candle.Date) {
			logger.Debugf("[util.FindClosesttDatedCandle] Breaking on candle dated %s", candle.Date)
			break
		}
		lastCandle = &candle
	}
	if finalCandle == nil {
		return selectedCandle, errors.New("[util.FindClosesttDatedCandle] Failed to locate any suibtable candlesticks")
	}
	if lastCandle == nil {
		selectedCandle = finalCandle
	} else {
		lastCandleDiff := lastCandle.Date.Sub(finalCandle.Date)
		finalCandleDiff := finalCandle.Date.Sub(lastCandle.Date)
		if lastCandleDiff < finalCandleDiff {
			selectedCandle = lastCandle
			logger.Debugf(fmt.Sprintf("[util.FindClosesttDatedCandle] Using last price dated %s instead of final candle dated %s", lastCandle.Date, finalCandle.Date))
		} else {
			selectedCandle = finalCandle
			logger.Debugf(fmt.Sprintf("[Binance.GetOrderHistory] Using final price dated %s (last candle dated %s)", finalCandle.Date, lastCandle.Date))
		}
	}
	return selectedCandle, nil
}

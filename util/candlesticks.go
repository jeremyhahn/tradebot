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

func FindClosestDatedCandle(logger *logging.Logger, needle time.Time, haystack []common.Candlestick) (*common.Candlestick, error) {
	var selectedCandle *common.Candlestick
	lastCandle := &common.Candlestick{}
	var finalCandle *common.Candlestick
	for _, candle := range haystack {
		finalCandle = &candle
		logger.Debugf("[util.FindClosestDatedCandle] Comparing candle date %s with target date %s", candle.Date, needle)
		if candle.Date.After(needle) || needle.Equal(candle.Date) {
			logger.Debugf("[util.FindClosestDatedCandle] Breaking on candle dated %s. Last candle: %s. Final candle: %s, Needle: %s",
				candle.Date, lastCandle.Date, finalCandle.Date, needle)
			*lastCandle = candle
			break
		}
		*lastCandle = candle
	}
	if finalCandle == nil {
		return selectedCandle, errors.New(fmt.Sprintf("[util.FindClosestDatedCandle] Unable to locate candlestick at %s", needle))
	}
	if lastCandle == nil {
		selectedCandle = finalCandle
	} else {
		lastCandleDiff := lastCandle.Date.Sub(finalCandle.Date)
		finalCandleDiff := finalCandle.Date.Sub(lastCandle.Date)
		if finalCandle.Date.Month() > needle.Month() || (finalCandle.Date.Month() == needle.Month() && finalCandle.Date.Day() > needle.Day()) {
			logger.Debugf(fmt.Sprintf("[util.FindClosestDatedCandle] Using prior day close price %s dated %s instead of final candle dated %s",
				finalCandle.Close, (*lastCandle).Date, finalCandle.Date))
			selectedCandle = lastCandle
		} else if lastCandleDiff < finalCandleDiff {
			selectedCandle = lastCandle
			logger.Debugf(fmt.Sprintf("[util.FindClosestDatedCandle] Using last price %s dated %s instead of final candle dated %s",
				finalCandle.Close, (*lastCandle).Date, finalCandle.Date))
		} else {
			selectedCandle = finalCandle
			logger.Debugf(fmt.Sprintf("[util.FindClosestDatedCandle] Using final price %s dated %s (last candle dated %s)",
				finalCandle.Close, (*lastCandle).Date, (*lastCandle).Date))
		}
	}
	return selectedCandle, nil
}

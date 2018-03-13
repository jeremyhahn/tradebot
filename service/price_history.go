package service

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dto"
	"github.com/jeremyhahn/tradebot/util"
)

type SlickCharts struct {
	ctx         common.Context
	endpoint    string
	prices      map[string][]dto.PriceHistoryDTO
	currencyMap map[string]string
	common.PriceHistoryService
}

func NewPriceHistoryService(ctx common.Context) common.PriceHistoryService {
	return &SlickCharts{
		ctx:      ctx,
		endpoint: "https://www.slickcharts.com/api/v1/currency",
		prices:   make(map[string][]dto.PriceHistoryDTO, 10),
		currencyMap: map[string]string{
			"BTC":  "bitcoin",
			"ETH":  "ethereum",
			"XRP":  "ripple",
			"BCC":  "bitcoin-cash",
			"BCH":  "bitcoin-cash",
			"LTC":  "litecoin",
			"NEO":  "neo",
			"ADA":  "cardano",
			"XLM":  "stellar",
			"XMR":  "monero",
			"EOS":  "eos",
			"DASH": "dash",
			"IOTA": "iota",
			"XEM":  "nem",
			"TRX":  "tron",
			"USDT": "tether",
			"GNT":  "golem-network-tokens",
			"SNM":  "sonm"}}
}

func (sc *SlickCharts) GetPriceOn(currency string, targetDay time.Time) common.PriceHistory {
	sc.ctx.GetLogger().Debugf("[PriceHistoryService.GetPriceOn] currency: %s targetDay: %s", currency, targetDay)
	if _, ok := sc.currencyMap[currency]; !ok {
		sc.ctx.GetLogger().Errorf("[PriceHistoryService.GetPriceOn] Currency not found in currencyMap: %s", currency)
		return nil
	}
	history := sc.GetPriceHistory(sc.currencyMap[currency])
	for _, price := range history {
		priceDate := time.Unix(price.GetTime(), 0)
		year, month, day := priceDate.Date()
		targetYear, targetMonth, _targetDay := targetDay.Date()
		sc.ctx.GetLogger().Debugf("[PriceHistoryService.GetPriceOn] Comparing targetDay: %s to price date: %s", targetDay, priceDate)
		if targetYear == year && _targetDay == day && targetMonth == month {
			return price
		}
	}
	return nil
}

func (sc *SlickCharts) GetPriceHistory(currency string) []common.PriceHistory {

	if _, ok := sc.prices[currency]; ok {
		sc.ctx.GetLogger().Debugf("[PriceHistoryService.GetPricesFor] Returning cached %s price history", currency)
		priceHistory := sc.prices[currency]
		var history []common.PriceHistory
		for _, price := range priceHistory {
			history = append(history, &dto.PriceHistoryDTO{
				Time:      price.GetTime(),
				High:      price.GetHigh(),
				Low:       price.GetLow(),
				Close:     price.GetClose(),
				Volume:    price.GetVolume(),
				MarketCap: price.GetMarketCap()})
		}
		return history
	}

	url := fmt.Sprintf("%s/%s/history", sc.endpoint, currency)

	sc.ctx.GetLogger().Debugf("[PriceHistoryService.GetPricesFor] Getting %s price history. Endpoint: %s", currency, url)

	_, body, err := util.HttpRequest(url)

	//sc.ctx.GetLogger().Debugf("[PriceHistoryService.GetPricesFor] Response: %s", string(body))

	if err != nil {
		sc.ctx.GetLogger().Errorf("[PriceHistoryService.GetPricesFor] Error: %s", err.Error())
	}

	var history []dto.PriceHistoryDTO
	jsonErr := json.Unmarshal(body, &history)
	if jsonErr != nil {
		sc.ctx.GetLogger().Errorf("[PriceHistoryService.GetPricesFor] Error: %s", jsonErr.Error())
	}

	sc.prices[currency] = history

	var commonHistory []common.PriceHistory
	for _, price := range history {
		commonHistory = append(commonHistory, &dto.PriceHistoryDTO{
			Time:      price.GetTime(),
			High:      price.GetHigh(),
			Low:       price.GetLow(),
			Close:     price.GetClose(),
			Volume:    price.GetVolume(),
			MarketCap: price.GetMarketCap()})
	}
	return commonHistory
}

func (sc *SlickCharts) GetClosePriceOn(currency string, date time.Time) float64 {
	//if currency == sc.ctx.GetUser().GetLocalCurrency() {
	//	return 1.0
	//}
	if _, ok := sc.currencyMap[currency]; !ok {
		sc.ctx.GetLogger().Errorf("[PriceHistoryService.GetClosePriceOn] Currency not found in currencyMap: %s", currency)
		return 0.0
	}
	history := sc.GetPriceHistory(sc.currencyMap[currency])
	for _, price := range history {
		timestamp := time.Unix(price.GetTime()/1000, 0)
		ts := time.Date(timestamp.Year(), timestamp.Month(), timestamp.Day(), 0, 0, 0, 0, timestamp.Location())
		d := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, timestamp.Location())
		if ts.Equal(d) {
			sc.ctx.GetLogger().Debugf("[PriceHistoryService.GetClosePriceOn] %s price on %s: %f", currency, date, price.GetClose())
			return price.GetClose()
		}
	}
	sc.ctx.GetLogger().Errorf("[PriceHistoryService.GetClosePriceOn] No price found for %s on %s", currency, date)
	return 0.0
}

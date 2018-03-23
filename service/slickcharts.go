package service

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/util"
)

type SlickChartsResponse struct {
	Time   int64   `json:"time"`
	Open   float64 `json:"open"`
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Close  float64 `json:"close"`
	Volume float64 `json:"volume"`
}

type SlickCharts struct {
	ctx         common.Context
	endpoint    string
	prices      map[string][]common.Candlestick
	currencyMap map[string]string
	common.FiatPriceService
}

func NewSlickChartsService(ctx common.Context) common.FiatPriceService {
	return &SlickCharts{
		ctx:      ctx,
		endpoint: "https://www.slickcharts.com/api/v1/currency",
		prices:   make(map[string][]common.Candlestick, 10),
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

func (sc *SlickCharts) GetPriceAt(currency string, targetDay time.Time) (*common.Candlestick, error) {
	sc.ctx.GetLogger().Debugf("[SlickCharts.GetPriceAt] currency: %s targetDay: %s", currency, targetDay)
	if _, ok := sc.currencyMap[currency]; !ok {
		sc.ctx.GetLogger().Errorf("[SlickCharts.GetPriceAt] Currency not found in currencyMap: %s", currency)
		//return nil, errors.New("Currency not found")
		return &common.Candlestick{}, nil
	}
	history := sc.GetPriceHistory(sc.currencyMap[currency])
	for _, price := range history {
		year, month, day := price.Date.Date()
		targetYear, targetMonth, _targetDay := targetDay.Date()
		//sc.ctx.GetLogger().Debugf("[SlickCharts.GetPriceAt] Comparing targetDay %s to %s with close price of %f.", targetDay, price.Date, price.Close)
		if targetYear > year {
			return &price, nil
		}
		if targetYear == year && _targetDay == day && targetMonth == month {
			return &price, nil
		}
	}
	errmsg := fmt.Sprintf("No %s price data found on %s", currency, targetDay)
	sc.ctx.GetLogger().Errorf("[SlickCharts.GetPriceAt] %s", errmsg)
	//return &common.Candlestick{}, errors.New(errmsg)
	return &common.Candlestick{}, nil
}

func (sc *SlickCharts) GetPriceHistory(currency string) []common.Candlestick {
	if _, ok := sc.prices[currency]; ok {
		sc.ctx.GetLogger().Debugf("[SlickCharts.GetPricesFor] Returning cached %s price history", currency)
		priceHistory := sc.prices[currency]
		var history []common.Candlestick
		for _, price := range priceHistory {
			history = append(history, common.Candlestick{
				Date:   price.Date,
				High:   price.High,
				Low:    price.Low,
				Close:  price.Close,
				Volume: price.Volume})
		}
		return history
	}

	url := fmt.Sprintf("%s/%s/history", sc.endpoint, currency)

	sc.ctx.GetLogger().Debugf("[SlickCharts.GetPricesFor] Getting %s price history. Endpoint: %s", currency, url)

	_, body, err := util.HttpRequest(url)

	//sc.ctx.GetLogger().Debugf("[SlickCharts.GetPricesFor] Response: %s", string(body))

	if err != nil {
		sc.ctx.GetLogger().Errorf("[SlickCharts.GetPricesFor] Error: %s", err.Error())
	}

	var history []SlickChartsResponse
	jsonErr := json.Unmarshal(body, &history)
	if jsonErr != nil {
		sc.ctx.GetLogger().Errorf("[SlickCharts.GetPricesFor] Error: %s", jsonErr.Error())
	}

	var candlesticks []common.Candlestick
	for _, price := range history {
		candlesticks = append(candlesticks, common.Candlestick{
			Date:   time.Unix(price.Time/1000, 0),
			High:   price.High,
			Low:    price.Low,
			Close:  price.Close,
			Volume: price.Volume})
	}

	sc.prices[currency] = candlesticks

	return candlesticks
}

/*
func (sc *SlickCharts) GetClosePriceOn(currency string, date time.Time) float64 {
	//if currency == sc.ctx.GetUser().GetLocalCurrency() {
	//	return 1.0
	//}
	if _, ok := sc.currencyMap[currency]; !ok {
		sc.ctx.GetLogger().Errorf("[SlickCharts.GetClosePriceOn] Currency not found in currencyMap: %s", currency)
		return 0.0
	}
	history := sc.GetPriceHistory(sc.currencyMap[currency])
	for _, price := range history {
		ts := time.Date(price.Date.Year(), price.Date.Month(), price.Date.Day(), 0, 0, 0, 0, price.Date.Location())
		d := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, price.Date.Location())
		if ts.Equal(d) {
			sc.ctx.GetLogger().Debugf("[SlickCharts.GetClosePriceOn] %s price on %s: %f", currency, date, price.Close)
			return price.Close
		}
	}
	sc.ctx.GetLogger().Errorf("[SlickCharts.GetClosePriceOn] No price found for %s on %s", currency, date)
	return 0.0
}
*/

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
	ctx      common.Context
	endpoint string
	prices   map[string][]dto.PriceHistoryDTO
	PriceHistoryService
}

func NewPriceHistoryService(ctx common.Context) PriceHistoryService {
	return &SlickCharts{
		ctx:      ctx,
		endpoint: "https://www.slickcharts.com/api/v1/currency",
		prices:   make(map[string][]dto.PriceHistoryDTO, 10)}
}

func (sc *SlickCharts) GetPriceHistory(currency string) []dto.PriceHistoryDTO {

	if _, ok := sc.prices[currency]; ok {
		sc.ctx.GetLogger().Debugf("[PriceHistoryService.GetPriceHistory] Returning cached %s price history", currency)
		return sc.prices[currency]
	}

	sc.ctx.GetLogger().Debugf("[PriceHistoryService.GetPriceHistory] Getting %s price history", currency)

	url := fmt.Sprintf("%s/%s/history", sc.endpoint, currency)
	_, body, err := util.HttpRequest(url)

	sc.ctx.GetLogger().Debugf("[PriceHistoryService.GetPriceHistory] Getting %s price history at endpoint %s - Response: ",
		currency, url, string(body))

	if err != nil {
		sc.ctx.GetLogger().Errorf("[PriceHistoryService.GetPriceHistory] Error: %s", err.Error())
	}

	var history []dto.PriceHistoryDTO
	jsonErr := json.Unmarshal(body, &history)
	if jsonErr != nil {
		sc.ctx.GetLogger().Errorf("[PriceHistoryService.GetPriceHistory] Error: %s", jsonErr.Error())
	}

	sc.prices[currency] = history

	return history
}

func (sc *SlickCharts) GetPrice(symbol string, date time.Time) float64 {
	if symbol == sc.ctx.GetUser().GetLocalCurrency() {
		return 1.0
	}
	prices := sc.GetPriceHistory(symbol)
	for _, price := range prices {
		timestamp := time.Unix(price.GetTime()/1000, 0)
		ts := time.Date(timestamp.Year(), timestamp.Month(), timestamp.Day(), 0, 0, 0, 0, timestamp.Location())
		d := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
		if ts.Equal(d) {
			sc.ctx.GetLogger().Debugf("[PriceHistoryService.GetPrice] Price: %f", price.GetClose())
			return price.GetClose()
		}
	}
	return 0.0
}

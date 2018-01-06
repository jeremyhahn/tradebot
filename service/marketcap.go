package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
)

type MarketCap struct {
	ctx        *common.Context
	client     http.Client
	dao        *dao.MarketCapDAO
	Markets    []common.MarketCap
	lastUpdate time.Time
}

func NewMarketCapService(ctx *common.Context, dao *dao.MarketCapDAO) *MarketCap {
	client := http.Client{Timeout: time.Second * 2}
	return &MarketCap{
		ctx:        ctx,
		dao:        dao,
		client:     client,
		lastUpdate: time.Now().Add(-1)}
}

func (m *MarketCap) GetMarkets() []common.MarketCap {

	if m.Markets == nil || time.Since(m.lastUpdate).Minutes() > 5 {

		limit := 10000
		m.ctx.Logger.Debugf("[NewMarketCap.GetMarkets] Fetching %d markets", limit)

		url := fmt.Sprintf("https://api.coinmarketcap.com/v1/ticker/?limit=%d", limit)
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			m.ctx.Logger.Errorf("[NewMarketCap.GetMarkets] %s", err.Error())
		}

		req.Header.Set("User-Agent", fmt.Sprintf("%s/v%s", common.APPNAME, common.APPVERSION))

		res, getErr := m.client.Do(req)
		if getErr != nil {
			m.ctx.Logger.Errorf("[NewMarketCap.GetMarkets] %s", getErr.Error())
		}

		body, readErr := ioutil.ReadAll(res.Body)
		if readErr != nil {
			m.ctx.Logger.Errorf("[NewMarketCap.GetMarkets] %s", readErr.Error())
		}

		jsonErr := json.Unmarshal(body, &m.Markets)
		if jsonErr != nil {
			m.ctx.Logger.Errorf("[NewMarketCap.GetMarkets] %s", jsonErr.Error())
		}
	}
	fmt.Printf("%+v\n", m.Markets)

	return m.Markets
}

func (m *MarketCap) GetGlobalMarket(currency string) *common.GlobalMarketCap {

	m.ctx.Logger.Debugf("[NewMarketCap.GetMarkets] Fetching global market data in %s currency", currency)

	url := fmt.Sprintf("https://api.coinmarketcap.com/v1/global/?convert=%s", currency)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		m.ctx.Logger.Errorf("[NewMarketCap.GetMarkets] %s", err.Error())
	}

	req.Header.Set("User-Agent", fmt.Sprintf("%s/v%s", common.APPNAME, common.APPVERSION))

	res, getErr := m.client.Do(req)
	if getErr != nil {
		m.ctx.Logger.Errorf("[NewMarketCap.GetMarkets] %s", getErr.Error())
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		m.ctx.Logger.Errorf("[NewMarketCap.GetMarkets] %s", readErr.Error())
	}

	gmarket := common.GlobalMarketCap{}
	jsonErr := json.Unmarshal(body, &gmarket)
	if jsonErr != nil {
		m.ctx.Logger.Errorf("[NewMarketCap.GetMarkets] %s", jsonErr.Error())
	}

	fmt.Printf("%+v\n", gmarket)

	return &gmarket
}

package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/jeremyhahn/tradebot/common"
	logging "github.com/op/go-logging"
)

type MarketCapService struct {
	logger     *logging.Logger
	client     http.Client
	Markets    []common.MarketCap
	lastUpdate int64
}

func NewMarketCapService(logger *logging.Logger) *MarketCapService {
	client := http.Client{Timeout: time.Second * 2}
	return &MarketCapService{
		logger:     logger,
		client:     client,
		lastUpdate: time.Now().Unix() - 10000}
}

func (m *MarketCapService) GetMarkets() []common.MarketCap {

	now := time.Now().Unix()
	diff := now - m.lastUpdate

	if len(m.Markets) == 0 || diff >= 300 {

		limit := 10000
		m.logger.Debugf("[NewMarketCap.GetMarkets] Fetching %d markets", limit)

		url := fmt.Sprintf("https://api.coinmarketcap.com/v1/ticker/?limit=%d", limit)
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			m.logger.Errorf("[NewMarketCap.GetMarkets] %s", err.Error())
		}

		req.Header.Set("User-Agent", fmt.Sprintf("%s/v%s", common.APPNAME, common.APPVERSION))

		res, getErr := m.client.Do(req)
		if getErr != nil {
			m.logger.Errorf("[NewMarketCap.GetMarkets] %s", getErr.Error())
		}

		body, readErr := ioutil.ReadAll(res.Body)
		if readErr != nil {
			m.logger.Errorf("[NewMarketCap.GetMarkets] %s", readErr.Error())
		}

		jsonErr := json.Unmarshal(body, &m.Markets)
		if jsonErr != nil {
			m.logger.Errorf("[NewMarketCap.GetMarkets] %s", jsonErr.Error())
		}

		fmt.Printf("%+v\n", m.Markets)

		fmt.Printf("Now: %d", now)
		fmt.Printf("Last: %d", m.lastUpdate)
		fmt.Printf("Diff: %d", diff)

		m.lastUpdate = now
	}

	return m.Markets
}

func (m *MarketCapService) GetMarket(symbol string) common.MarketCap {
	markets := m.GetMarkets()
	for _, m := range markets {
		if m.Symbol == symbol {
			return m
		}
	}
	return common.MarketCap{}
}

func (m *MarketCapService) GetGlobalMarket(currency string) *common.GlobalMarketCap {

	m.logger.Debugf("[NewMarketCap.GetMarkets] Fetching global market data in %s currency", currency)

	url := fmt.Sprintf("https://api.coinmarketcap.com/v1/global/?convert=%s", currency)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		m.logger.Errorf("[NewMarketCap.GetMarkets] %s", err.Error())
	}

	req.Header.Set("User-Agent", fmt.Sprintf("%s/v%s", common.APPNAME, common.APPVERSION))

	res, getErr := m.client.Do(req)
	if getErr != nil {
		m.logger.Errorf("[NewMarketCap.GetMarkets] %s", getErr.Error())
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		m.logger.Errorf("[NewMarketCap.GetMarkets] %s", readErr.Error())
	}

	gmarket := common.GlobalMarketCap{}
	jsonErr := json.Unmarshal(body, &gmarket)
	if jsonErr != nil {
		m.logger.Errorf("[NewMarketCap.GetMarkets] %s", jsonErr.Error())
	}

	fmt.Printf("%+v\n", gmarket)

	return &gmarket
}

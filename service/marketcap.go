package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/jeremyhahn/tradebot/common"
	logging "github.com/op/go-logging"
)

type MarketCapService struct {
	logger           *logging.Logger
	client           http.Client
	Markets          []common.MarketCap
	GlobalMarket     *common.GlobalMarketCap
	lastUpdate       int64
	lastGlobalUpdate int64
	interval         int64
}

func NewMarketCapService(logger *logging.Logger) *MarketCapService {
	client := http.Client{Timeout: time.Second * 2}
	return &MarketCapService{
		logger:           logger,
		client:           client,
		lastUpdate:       time.Now().Unix() - 10000,
		lastGlobalUpdate: time.Now().Unix() - 10000,
		interval:         300} // 5 minutes
}

func (m *MarketCapService) GetMarkets() []common.MarketCap {

	now := time.Now().Unix()
	diff := now - m.lastUpdate

	if diff >= m.interval {

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

		fmt.Println("Fetching marketcap")
		m.logger.Debugf("[NewMarketCap.GetMarkets] Now: %d, Last: %d, Diff: %d, Markets: %+v\n", now, m.lastUpdate, diff, m.Markets)
		m.lastUpdate = now
	}

	var marketList []common.MarketCap
	for _, m := range m.Markets {
		if m.PriceUSD == "" {
			continue
		}
		marketList = append(marketList, m)
	}
	m.Markets = marketList

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

	now := time.Now().Unix()
	diff := now - m.lastUpdate

	gmarket := common.GlobalMarketCap{}

	if diff >= m.interval {

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

		jsonErr := json.Unmarshal(body, &gmarket)
		if jsonErr != nil {
			m.logger.Errorf("[NewMarketCap.GetMarkets] %s", jsonErr.Error())
		}

		fmt.Printf("%+v\n", gmarket)
	}

	return &gmarket
}

func (m *MarketCapService) GetMarketsByPrice(order string) []common.MarketCap {
	markets := m.GetMarkets()
	sort.Slice(markets, func(i, j int) bool {
		priceI, _ := strconv.ParseFloat(markets[i].PriceUSD, 64)
		priceJ, _ := strconv.ParseFloat(markets[j].PriceUSD, 64)
		if order == "asc" {
			return priceI < priceJ
		} else {
			return priceI > priceJ
		}

	})
	return markets
}

func (m *MarketCapService) GetMarketsByPercentChange1H(order string) []common.MarketCap {
	markets := m.GetMarkets()
	sort.Slice(markets, func(i, j int) bool {
		fi, _ := strconv.ParseFloat(markets[i].PercentChange1h, 64)
		fj, _ := strconv.ParseFloat(markets[j].PercentChange1h, 64)
		if order == "asc" {
			return fi < fj
		} else {
			return fi > fj
		}
	})
	return markets
}

func (m *MarketCapService) GetMarketsByPercentChange24H(order string) []common.MarketCap {
	markets := m.GetMarkets()
	sort.Slice(markets, func(i, j int) bool {
		fi, _ := strconv.ParseFloat(markets[i].PercentChange24h, 64)
		fj, _ := strconv.ParseFloat(markets[j].PercentChange24h, 64)
		if order == "asc" {
			return fi < fj
		} else {
			return fi > fj
		}
	})
	return markets
}

func (m *MarketCapService) GetMarketsByPercentChange7D(order string) []common.MarketCap {
	markets := m.GetMarkets()
	sort.Slice(markets, func(i, j int) bool {
		fi, _ := strconv.ParseFloat(markets[i].PercentChange7d, 64)
		fj, _ := strconv.ParseFloat(markets[j].PercentChange7d, 64)
		if order == "asc" {
			return fi < fj
		} else {
			return fi > fj
		}
	})
	return markets
}

func (m *MarketCapService) GetMarketsByTopPerformers(order string) []common.MarketCap {
	markets := m.GetMarkets()
	sort.Slice(markets, func(i, j int) bool {
		fi1h, _ := strconv.ParseFloat(markets[i].PercentChange1h, 64)
		fj1h, _ := strconv.ParseFloat(markets[j].PercentChange1h, 64)
		fi24h, _ := strconv.ParseFloat(markets[i].PercentChange24h, 64)
		fj24h, _ := strconv.ParseFloat(markets[j].PercentChange24h, 64)
		fi7d, _ := strconv.ParseFloat(markets[i].PercentChange7d, 64)
		fj7d, _ := strconv.ParseFloat(markets[j].PercentChange7d, 64)
		avgi := fi1h + fi24h + fi7d
		avgj := fj1h + fj24h + fj7d
		if order == "asc" {
			return avgi < avgj
		} else {
			return avgi > avgj
		}
	})
	return markets
}

func (m *MarketCapService) GetMarketsByTrending(order string) []common.MarketCap {
	markets := m.GetMarkets()
	sort.Slice(markets, func(i, j int) bool {
		fi, _ := strconv.ParseFloat(markets[i].PercentChange1h, 64)
		fj, _ := strconv.ParseFloat(markets[j].PercentChange1h, 64)
		if order == "asc" {
			return fi < fj
		} else {
			return fi > fj
		}
	})
	return markets
}

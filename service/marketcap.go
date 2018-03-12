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

type MarketCapService interface {
	GetMarkets() []common.MarketCap
	GetMarket(symbol string) common.MarketCap
	GetGlobalMarket(currency string) *common.GlobalMarketCap
}

type MarketCapServiceImpl struct {
	logger           *logging.Logger
	client           http.Client
	Markets          []common.MarketCap
	GlobalMarket     *common.GlobalMarketCap
	lastUpdate       int64
	lastGlobalUpdate int64
	interval         int64
}

var MARKETCAP_RATELIMITER = common.NewRateLimiter(10, 1)

func NewMarketCapService(ctx common.Context) MarketCapService {
	client := http.Client{Timeout: time.Second * 2}
	return &MarketCapServiceImpl{
		logger:           ctx.GetLogger(),
		client:           client,
		lastUpdate:       time.Now().Unix() - 10000,
		lastGlobalUpdate: time.Now().Unix() - 10000,
		interval:         300} // 5 minutes
}

func (m *MarketCapServiceImpl) GetMarkets() []common.MarketCap {

	MARKETCAP_RATELIMITER.RespectRateLimit()

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

func (service *MarketCapServiceImpl) GetMarket(symbol string) common.MarketCap {
	MARKETCAP_RATELIMITER.RespectRateLimit()
	markets := service.GetMarkets()
	for _, m := range markets {
		if m.Symbol == symbol {
			service.logger.Debugf("[MarketCapServiceImpl.GetMarket] Getting market: %+v", m)
			return m
		}
	}
	return common.MarketCap{}
}

func (m *MarketCapServiceImpl) GetGlobalMarket(currency string) *common.GlobalMarketCap {

	gmarket := common.GlobalMarketCap{}

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

	return &gmarket
}

func (m *MarketCapServiceImpl) GetMarketsByPrice(order string) []common.MarketCap {
	MARKETCAP_RATELIMITER.RespectRateLimit()
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

func (m *MarketCapServiceImpl) GetMarketsByPercentChange1H(order string) []common.MarketCap {
	MARKETCAP_RATELIMITER.RespectRateLimit()
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

func (m *MarketCapServiceImpl) GetMarketsByPercentChange24H(order string) []common.MarketCap {
	MARKETCAP_RATELIMITER.RespectRateLimit()
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

func (m *MarketCapServiceImpl) GetMarketsByPercentChange7D(order string) []common.MarketCap {
	MARKETCAP_RATELIMITER.RespectRateLimit()
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

func (m *MarketCapServiceImpl) GetMarketsByTopPerformers(order string) []common.MarketCap {
	MARKETCAP_RATELIMITER.RespectRateLimit()
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

func (m *MarketCapServiceImpl) GetTrendingMarkets(order string) []common.MarketCap {
	MARKETCAP_RATELIMITER.RespectRateLimit()
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

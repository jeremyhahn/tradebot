package exchange

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	ws "github.com/gorilla/websocket"
	gdax "github.com/preichenberger/go-gdax"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
	"github.com/op/go-logging"
)

var GDAXLastApiCall = time.Now().AddDate(0, 0, -1).Unix()

type GDAX struct {
	gdax            *gdax.Client
	logger          *logging.Logger
	name            string
	lastApiCall     int64
	lastBalanceCall int64
	balances        []common.Coin
	netWorth        float64
	apiCallCount    int
	currencyPair    *common.CurrencyPair
	tradingFee      float64
	common.Exchange
}

func NewGDAX(_gdax *dao.UserCoinExchange, logger *logging.Logger, currencyPair *common.CurrencyPair) common.Exchange {
	return &GDAX{
		gdax:         gdax.NewClient(_gdax.Secret, _gdax.Key, _gdax.Passphrase),
		logger:       logger,
		name:         "gdax",
		apiCallCount: 0,
		currencyPair: currencyPair,
		netWorth:     0.0,
		balances:     make([]common.Coin, 0),
		tradingFee:   0.025}
}

func (_gdax *GDAX) GetTradeHistory(start, end time.Time, granularity int) []common.Candlestick {
	_gdax.respectRateLimit()
	_gdax.logger.Debug("[GDAX.GetTradeHistory] Getting trade history")
	var candlesticks []common.Candlestick
	params := gdax.GetHistoricRatesParams{
		Start:       start,
		End:         end,
		Granularity: granularity}
	rates, err := _gdax.gdax.GetHistoricRates(_gdax.FormattedCurrencyPair(), params)
	if err != nil {
		if strings.Contains(err.Error(), "granularity too small for the requested time range") {
			_gdax.logger.Debug("[GDAX.GetTradeHistory] Result set too big; chunking into smaller requests...")
			var candlesticks []common.Candlestick
			diff := end.Sub(start)
			days := int(diff.Hours() / 24)
			for i := 0; i < days; i++ {
				newStart := start.AddDate(0, 0, i)
				newEnd := start.AddDate(0, 0, i).Add(time.Hour*23 + time.Minute*59 + time.Second*59)
				sticks := _gdax.GetTradeHistory(newStart, newEnd, granularity)
				//bytes, _ := json.MarshalIndent(sticks, "", "    ")
				//fmt.Println(string(bytes))
				candlesticks = append(candlesticks, sticks...)
			}
			return candlesticks
		}
		_gdax.logger.Errorf("[GDAX.GetTradeHistory] %s", err.Error())
		time.Sleep(time.Second * 1)
		return _gdax.GetTradeHistory(start, end, granularity)
	}
	for _, r := range rates {
		candlesticks = append(candlesticks, common.Candlestick{
			Exchange:     _gdax.name,
			CurrencyPair: _gdax.currencyPair,
			Period:       granularity,
			Date:         r.Time,
			Open:         r.Open,
			Close:        r.Close,
			High:         r.High,
			Low:          r.Low,
			Volume:       r.Volume})
	}
	return candlesticks
}

func (_gdax *GDAX) GetBalances() ([]common.Coin, float64) {

	cacheTime := time.Now().Unix()
	cacheDiff := cacheTime - _gdax.lastBalanceCall
	if cacheDiff <= 3600 && len(_gdax.balances) > 0 {
		_gdax.logger.Debug("[GDAX.GetBalances] Returning cached balances")
		return _gdax.balances, _gdax.netWorth
	}

	_gdax.respectRateLimit()

	_gdax.logger.Debugf("[GDAX] Getting %s balances", _gdax.currencyPair.Base)
	var coins []common.Coin
	sum := 0.0
	accounts, err := _gdax.gdax.GetAccounts()
	if err != nil {
		_gdax.logger.Errorf("[GDAX.GetBalances] %s", err.Error())
		return coins, sum
	}
	for _, a := range accounts {
		price := 1.0
		if a.Currency != _gdax.currencyPair.LocalCurrency {
			currency := fmt.Sprintf("%s-%s", a.Currency, _gdax.currencyPair.Quote)
			_gdax.logger.Debugf("[GDAX.GetBalances] Getting balances for %s", currency)
			ticker, err := _gdax.gdax.GetTicker(currency)
			if err != nil {
				_gdax.logger.Errorf("[GDAX.GetBalances] %s", err.Error())
				continue
			}
			price = ticker.Price
		}
		if a.Balance <= 0 {
			continue
		}
		total := a.Balance * price
		t, err := strconv.ParseFloat(fmt.Sprintf("%.2f", total), 64)
		if err != nil {
			_gdax.logger.Errorf("[GDAX.GetBalances] %s", err.Error())
			continue
		}
		sum += total
		coins = append(coins, common.Coin{
			Currency:  a.Currency,
			Balance:   a.Balance,
			Available: a.Available,
			Pending:   a.Hold,
			Price:     price,
			Total:     t})
	}
	_gdax.balances = coins
	_gdax.netWorth = sum
	_gdax.lastBalanceCall = time.Now().Unix()
	return coins, sum
}

func (_gdax *GDAX) SubscribeToLiveFeed(priceChannel chan common.PriceChange) {
	_gdax.logger.Info("[GDAX.SubscribeToLiveFeed] Subscribing to WebSocket feed")

	var wsDialer ws.Dialer
	wsConn, _, err := wsDialer.Dial("wss://ws-feed.gdax.com", nil)
	if err != nil {
		println(err.Error())
	}

	subscribe := map[string]string{
		"type":       "subscribe",
		"product_id": _gdax.FormattedCurrencyPair(),
	}

	if err := wsConn.WriteJSON(subscribe); err != nil {
		_gdax.logger.Errorf("[GDAX.SubscribeToLiveFeed] %s", err.Error())
	}

	message := gdax.Message{}
	for true {

		if err := wsConn.ReadJSON(&message); err != nil {
			_gdax.logger.Errorf("[GDAX.SubscribeToLiveFeed] %s", err.Error())
			break
		}

		if message.Type == "match" && message.Reason == "filled" {
			_gdax.logger.Debugf("[GDAX.SubscribeToLiveFeed] message: %+v\n", message)
			priceChannel <- common.PriceChange{
				Exchange:     _gdax.GetName(),
				CurrencyPair: _gdax.currencyPair,
				Price:        message.Price}
		}
	}

	_gdax.SubscribeToLiveFeed(priceChannel)
}

func (_gdax *GDAX) GetExchangeAsync(exchangeChan *chan common.CoinExchange) {
	go func() { *exchangeChan <- _gdax.GetExchange() }()
}

func (_gdax *GDAX) GetExchange() common.CoinExchange {
	total := 0.0
	satoshis := 0.0
	balances, _ := _gdax.GetBalances()
	for _, c := range balances {
		if c.Currency == _gdax.currencyPair.LocalCurrency {
			total += c.Total
		} else if c.IsBitcoin() {
			satoshis += c.Balance
			total += c.Total
		} else {
			currency := fmt.Sprintf("%s-BTC", c.Currency)
			_gdax.respectRateLimit()
			ticker, err := _gdax.gdax.GetTicker(currency)
			if err != nil {
				_gdax.logger.Errorf("[GDAX.GetExchange] %s", err.Error())
				continue
			}
			satoshis += ticker.Price
			total += c.Total
		}
	}
	s, _ := strconv.ParseFloat(fmt.Sprintf("%.8f", satoshis), 64)
	t, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", total), 64)
	exchange := common.CoinExchange{
		Name:     _gdax.name,
		URL:      "https://www.gdax.com",
		Total:    t,
		Satoshis: s,
		Coins:    balances}
	return exchange
}

func (_gdax *GDAX) ToUSD(price, satoshis float64) float64 {
	return satoshis * price
}

func (_gdax *GDAX) GetCurrencyPair() common.CurrencyPair {
	return *_gdax.currencyPair
}

func (_gdax *GDAX) GetName() string {
	return _gdax.name
}

func (_gdax *GDAX) FormattedCurrencyPair() string {
	return fmt.Sprintf("%s-%s", _gdax.currencyPair.Base, _gdax.currencyPair.Quote)
}

func (_gdax *GDAX) GetTradingFee() float64 {
	return _gdax.tradingFee
}

func (_gdax *GDAX) respectRateLimit() {
	now := time.Now().Unix()
	diff := now - _gdax.lastApiCall
	for diff <= 30 && _gdax.apiCallCount >= 3 {
		_gdax.logger.Info("[GDAX.respectRateLimit] Cooling off")
		_gdax.logger.Debugf("[GDAX.respectRateLimit] apiCallCount: %d, lastApiCall: %s", _gdax.apiCallCount, _gdax.lastApiCall)
		time.Sleep(1 * time.Second)
		GDAXLastApiCall = -1
	}
	GDAXLastApiCall = time.Now().Unix()
	_gdax.apiCallCount += 1
}

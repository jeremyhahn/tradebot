package exchange

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	ws "github.com/gorilla/websocket"
	gdax "github.com/preichenberger/go-gdax"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dto"
	"github.com/jeremyhahn/tradebot/entity"
	"github.com/jeremyhahn/tradebot/util"
	"github.com/op/go-logging"
)

var GDAXLastApiCall = time.Now().AddDate(0, 0, -1).Unix()
var GDAX_MUTEX sync.Mutex

type GDAX struct {
	gdax            *gdax.Client
	ctx             *common.Context
	logger          *logging.Logger
	name            string
	lastApiCall     int64
	lastBalanceCall int64
	balances        []common.Coin
	netWorth        float64
	apiCallCount    int
	tradingFee      float64
	common.Exchange
}

func NewGDAX(ctx *common.Context, _gdax entity.UserExchangeEntity) common.Exchange {
	return &GDAX{
		ctx:          ctx,
		gdax:         gdax.NewClient(_gdax.GetSecret(), _gdax.GetKey(), _gdax.GetExtra()),
		logger:       ctx.Logger,
		name:         "gdax",
		apiCallCount: 0,
		netWorth:     0.0,
		balances:     make([]common.Coin, 0),
		tradingFee:   0.025}
}

func (_gdax *GDAX) GetPriceHistory(currencyPair *common.CurrencyPair,
	start, end time.Time, granularity int) []common.Candlestick {

	_gdax.respectRateLimit()
	_gdax.logger.Debugf("[GDAX.GetPriceHistory] Getting price history %s - %s with granularity %d",
		util.FormatDate(start), util.FormatDate(end), granularity)
	var candlesticks []common.Candlestick
	params := gdax.GetHistoricRatesParams{
		Start:       start,
		End:         end,
		Granularity: granularity}
	rates, err := _gdax.gdax.GetHistoricRates(_gdax.FormattedCurrencyPair(currencyPair), params)
	if err != nil {
		if strings.Contains(err.Error(), "granularity too small for the requested time range") {
			_gdax.logger.Debug("[GDAX.GetPriceHistory] Result set too big; chunking into smaller requests...")
			diff := end.Sub(start)
			days := int(diff.Hours() / 24)
			var newEnd time.Time
			for i := 0; i < days; i++ {
				newStart := start.AddDate(0, 0, i)
				newEnd = start.AddDate(0, 0, i).Add(time.Hour*23 + time.Minute*59 + time.Second*59)
				sticks := _gdax.GetPriceHistory(currencyPair, newStart, newEnd, granularity)
				sticks = _gdax.reverseCandlesticks(sticks)
				candlesticks = append(candlesticks, sticks...)
			}
			finalStart := start.AddDate(0, 0, days)
			finalEnd := end.AddDate(0, 0, 0).Add(time.Duration(granularity*-1) * time.Second)
			sticks := _gdax.GetPriceHistory(currencyPair, finalStart, finalEnd, granularity)
			sticks = _gdax.reverseCandlesticks(sticks)
			candlesticks = append(candlesticks, sticks...)
			return candlesticks
		} else {
			_gdax.logger.Errorf("[GDAX.GetPriceHistory] GDAX API Error: %s", err.Error())
			return candlesticks
		}
	}
	for _, r := range rates {
		candlesticks = append(candlesticks, common.Candlestick{
			Exchange:     _gdax.name,
			CurrencyPair: currencyPair,
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

func (_gdax *GDAX) GetOrderHistory(currencyPair *common.CurrencyPair) []common.Order {
	_gdax.respectRateLimit()
	_gdax.logger.Debug("[GDAX.GetOrderHistory] Getting order history")
	var orders []common.Order
	var ledger []gdax.LedgerEntry
	orderIds := make(map[string]bool)
	accounts, err := _gdax.gdax.GetAccounts()
	if err != nil {
		_gdax.logger.Errorf("[GDAX.GetOrderHistory] %s", err.Error())
	}
	for _, a := range accounts {
		cursor := _gdax.gdax.ListAccountLedger(a.Id)
		for cursor.HasMore {
			if err := cursor.NextPage(&ledger); err != nil {
				_gdax.logger.Errorf("[GDAX.GetOrderHistory] %s", err.Error())
			}
			for _, e := range ledger {
				if e.Type != "match" {
					continue
				}
				if _, ok := orderIds[e.Details.OrderId]; ok {
					continue
				}
				orderIds[e.Details.OrderId] = true
				order, err := _gdax.gdax.GetOrder(e.Details.OrderId)
				if err != nil {
					_gdax.ctx.Logger.Errorf("[GDAX.GetOrderHistory] Error retrieving order: %s", err.Error())
					return orders
				}
				pieces := strings.Split(e.Details.ProductId, "-")
				base, quote := pieces[0], pieces[1]
				orders = append(orders, &dto.OrderDTO{
					Id:       strconv.FormatInt(int64(e.Id), 10),
					Exchange: "gdax",
					Date:     e.CreatedAt.Time(),
					Type:     order.Side,
					CurrencyPair: &common.CurrencyPair{
						Base:  base,
						Quote: quote},
					Quantity:      order.FilledSize,
					Fee:           util.TruncateFloat(order.FillFees, 2),
					Price:         order.Price,
					Total:         order.ExecutedValue,
					PriceCurrency: quote,
					FeeCurrency:   quote,
					TotalCurrency: quote})
			}
		}
	}
	return orders
}

func (_gdax *GDAX) GetBalances() ([]common.Coin, float64) {
	cacheTime := time.Now().Unix()
	cacheDiff := cacheTime - _gdax.lastBalanceCall
	if cacheDiff <= 3600 && len(_gdax.balances) > 0 {
		_gdax.logger.Debug("[GDAX.GetBalances] Returning cached balances")
		return _gdax.balances, _gdax.netWorth
	}
	_gdax.respectRateLimit()
	_gdax.logger.Debugf("[GDAX] Getting balances")
	var coins []common.Coin
	sum := 0.0
	accounts, err := _gdax.gdax.GetAccounts()
	if err != nil {
		_gdax.logger.Errorf("[GDAX.GetBalances] %s", err.Error())
		return coins, sum
	}
	for _, a := range accounts {
		price := 1.0
		if a.Currency != _gdax.ctx.User.GetLocalCurrency() {
			currency := fmt.Sprintf("%s-%s", a.Currency, _gdax.ctx.User.GetLocalCurrency())
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
		coins = append(coins, &dto.CoinDTO{
			Currency:  a.Currency,
			Balance:   a.Balance,
			Available: a.Available,
			Pending:   a.Hold,
			Price:     price,
			Total:     t})
	}
	GDAX_MUTEX.Lock()
	_gdax.balances = coins
	_gdax.netWorth = sum
	_gdax.lastBalanceCall = time.Now().Unix()
	GDAX_MUTEX.Unlock()
	return coins, sum
}

func (_gdax *GDAX) SubscribeToLiveFeed(currencyPair *common.CurrencyPair,
	priceChannel chan common.PriceChange) {

	_gdax.logger.Info("[GDAX.SubscribeToLiveFeed] Subscribing to WebSocket feed")

	var wsDialer ws.Dialer
	wsConn, _, err := wsDialer.Dial("wss://ws-feed.gdax.com", nil)
	if err != nil {
		_gdax.logger.Errorf("[GDAX.SubscribeToLiveFeed] %s", err.Error())
	}

	subscribe := map[string]string{
		"type":       "subscribe",
		"product_id": _gdax.FormattedCurrencyPair(currencyPair),
	}

	if err := wsConn.WriteJSON(subscribe); err != nil {
		_gdax.logger.Errorf("[GDAX.SubscribeToLiveFeed] %s", err.Error())
	}

	message := gdax.Message{}
	for true {

		if err := wsConn.ReadJSON(&message); err != nil {
			_gdax.logger.Errorf("[GDAX.SubscribeToLiveFeed] %s", err.Error())
			_gdax.SubscribeToLiveFeed(currencyPair, priceChannel)
		}

		if message.Type == "match" && message.Reason == "filled" {
			_gdax.logger.Debugf("[GDAX.SubscribeToLiveFeed] message: %+v\n", message)
			priceChannel <- common.PriceChange{
				Exchange:     _gdax.GetName(),
				CurrencyPair: currencyPair,
				Price:        message.Price}
		}
	}

	_gdax.SubscribeToLiveFeed(currencyPair, priceChannel)
}

func (_gdax *GDAX) GetExchange() common.CryptoExchange {
	total := 0.0
	satoshis := 0.0
	balances, _ := _gdax.GetBalances()
	for _, c := range balances {
		if c.GetCurrency() == _gdax.ctx.User.GetLocalCurrency() {
			total += c.GetTotal()
		} else if c.IsBitcoin() {
			satoshis += c.GetBalance()
			total += c.GetTotal()
		} else {
			currency := fmt.Sprintf("%s-BTC", c.GetCurrency())
			_gdax.respectRateLimit()
			ticker, err := _gdax.gdax.GetTicker(currency)
			if err != nil {
				_gdax.logger.Errorf("[GDAX.GetExchange] %s", err.Error())
				continue
			}
			satoshis += ticker.Price
			total += c.GetTotal()
		}
	}
	s, _ := strconv.ParseFloat(fmt.Sprintf("%.8f", satoshis), 64)
	t, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", total), 64)
	exchange := &dto.CryptoExchangeDTO{
		Name:     _gdax.name,
		URL:      "https://www.gdax.com",
		Total:    t,
		Satoshis: s,
		Coins:    balances}
	return exchange
}

func (_gdax *GDAX) ParseImport(file string) ([]common.Order, error) {
	var orders []common.Order
	_gdax.ctx.Logger.Error("[GDAX.ParseImport] Unsupported!")
	return orders, errors.New("GDAX.ParseImport Unsupported")
}

func (_gdax *GDAX) ToUSD(price, satoshis float64) float64 {
	return satoshis * price
}

func (_gdax *GDAX) GetName() string {
	return _gdax.name
}

func (_gdax *GDAX) FormattedCurrencyPair(currencyPair *common.CurrencyPair) string {
	return fmt.Sprintf("%s-%s", currencyPair.Base, currencyPair.Quote)
}

func (_gdax *GDAX) formatCurrencyPair(currencyPair *common.CurrencyPair) string {
	return fmt.Sprintf("%s-%s", currencyPair.Base, currencyPair.Quote)
}

func (_gdax *GDAX) GetTradingFee() float64 {
	return _gdax.tradingFee
}

func (_gdax *GDAX) respectRateLimit() {
	now := time.Now().Unix()
	diff := now - GDAXLastApiCall
	for diff <= 30 && _gdax.apiCallCount >= 3 {
		_gdax.logger.Info("[GDAX.respectRateLimit] Cooling off")
		_gdax.logger.Debugf("[GDAX.respectRateLimit] apiCallCount: %d, lastApiCall: %d", _gdax.apiCallCount, GDAXLastApiCall)
		time.Sleep(1 * time.Second)
		GDAX_MUTEX.Lock()
		_gdax.apiCallCount = 0
		GDAX_MUTEX.Unlock()
	}
	GDAX_MUTEX.Lock()
	GDAXLastApiCall = time.Now().Unix()
	_gdax.apiCallCount += 1
	GDAX_MUTEX.Unlock()
}

func (_gdax *GDAX) reverseCandlesticks(candles []common.Candlestick) []common.Candlestick {
	var reversed []common.Candlestick
	for i := len(candles) - 1; i > 0; i-- {
		reversed = append(reversed, candles[i])
	}
	return reversed
}

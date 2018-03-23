package main

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"

	ws "github.com/gorilla/websocket"
	gdax "github.com/preichenberger/go-gdax"
	"github.com/shopspring/decimal"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dto"
	"github.com/jeremyhahn/tradebot/entity"
	"github.com/jeremyhahn/tradebot/util"
	"github.com/op/go-logging"
	"github.com/patrickmn/go-cache"
)

type GDAX struct {
	gdax        *gdax.Client
	ctx         common.Context
	logger      *logging.Logger
	name        string
	displayName string
	tradingFee  float64
	rateLimiter *common.RateLimiter
	cache       *cache.Cache
	common.Exchange
	common.FiatPriceService
}

var GDAX_RATELIMITER = common.NewRateLimiter(3, 1)
var GDAX_MUTEX sync.Mutex

func CreateGDAX(ctx common.Context, _gdax entity.UserExchangeEntity) common.Exchange {
	return &GDAX{
		ctx:         ctx,
		gdax:        gdax.NewClient(_gdax.GetSecret(), _gdax.GetKey(), _gdax.GetExtra()),
		logger:      ctx.GetLogger(),
		name:        "GDAX",
		displayName: "GDAX",
		tradingFee:  0.025,
		cache:       cache.New(1*time.Minute, 1*time.Minute)}
}

func (_gdax *GDAX) GetPriceAt(currency string, targetDay time.Time) (*common.Candlestick, error) {
	currencyPair := &common.CurrencyPair{
		Base:          currency,
		Quote:         _gdax.ctx.GetUser().GetLocalCurrency(),
		LocalCurrency: _gdax.ctx.GetUser().GetLocalCurrency()}
	startDate := targetDay.Add(-1 * time.Minute)
	endDate := targetDay.Add(5 * time.Minute)
	candles, err := _gdax.GetPriceHistory(currencyPair, startDate, endDate, 60)
	if err != nil {
		return &common.Candlestick{}, err
	}
	if len(candles) > 0 {
		return &candles[0], nil
	}
	_gdax.ctx.GetLogger().Debugf("[GDAX.GetPriceAt] Unable to locate between % and %s", currencyPair, startDate, endDate)
	return &common.Candlestick{}, nil
}

func (_gdax *GDAX) GetPriceHistory(currencyPair *common.CurrencyPair,
	start, end time.Time, granularity int) ([]common.Candlestick, error) {
	GDAX_RATELIMITER.RespectRateLimit()
	_gdax.logger.Debugf("[GDAX.GetPriceHistory] Getting %s price history %s - %s with granularity %d",
		currencyPair, util.FormatDate(start), util.FormatDate(end), granularity)
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
				sticks, err := _gdax.GetPriceHistory(currencyPair, newStart, newEnd, granularity)
				if err != nil {
					_gdax.logger.Debug("[GDAX.GetPriceHistory] Error: %s", err)
				}
				sticks = util.ReverseCandlesticks(sticks)
				candlesticks = append(candlesticks, sticks...)
			}
			finalStart := start.AddDate(0, 0, days)
			finalEnd := end.AddDate(0, 0, 0).Add(time.Duration(granularity*-1) * time.Second)
			sticks, err := _gdax.GetPriceHistory(currencyPair, finalStart, finalEnd, granularity)
			if err != nil {
				_gdax.logger.Debug("[GDAX.GetPriceHistory] Error: %s", err)
			}
			sticks = util.ReverseCandlesticks(sticks)
			candlesticks = append(candlesticks, sticks...)
			return candlesticks, nil
		} else {
			_gdax.logger.Errorf("[GDAX.GetPriceHistory] GDAX API Error: %s", err.Error())
			return candlesticks, err
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
	return candlesticks, nil
}

func (_gdax *GDAX) GetOrderHistory(currencyPair *common.CurrencyPair) []common.Transaction {
	GDAX_RATELIMITER.RespectRateLimit()
	_gdax.logger.Debug("[GDAX.GetOrderHistory] Getting order history")
	var orders []common.Transaction
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
					_gdax.ctx.GetLogger().Errorf("[GDAX.GetOrderHistory] Error retrieving order: %s", err.Error())
					continue
				}
				currencyPair, err := common.NewCurrencyPair(e.Details.ProductId, _gdax.ctx.GetUser().GetLocalCurrency())
				if err != nil {
					_gdax.ctx.GetLogger().Errorf("[GDAX.GetOrderHistory] Error parsing currency pair: %s", err.Error())
					continue
				}
				baseCurrency, err := _gdax.getCurrency(currencyPair.Base)
				if err != nil {
					_gdax.ctx.GetLogger().Errorf("[GDAX.GetOrderHistory] Unsupported base currency: %s", currencyPair.Base)
					continue
				}
				quoteCurrency, err := _gdax.getCurrency(currencyPair.Quote)
				if err != nil {
					_gdax.ctx.GetLogger().Errorf("[GDAX.GetOrderHistory] Unsupported quote currency: %s", currencyPair.Quote)
					continue
				}
				orderDate := e.CreatedAt.Time()
				quantity := decimal.NewFromFloat(order.FilledSize).StringFixed(baseCurrency.GetDecimalPlace())
				fee := decimal.NewFromFloat(order.FillFees)
				total := decimal.NewFromFloat(order.ExecutedValue)
				price := total.Sub(fee)

				orders = append(orders, &dto.TransactionDTO{
					Id:                   strconv.FormatInt(int64(e.Id), 10),
					Type:                 order.Side,
					Date:                 orderDate,
					Network:              _gdax.name,
					NetworkDisplayName:   _gdax.displayName,
					CurrencyPair:         currencyPair,
					Quantity:             quantity,
					QuantityCurrency:     currencyPair.Base,
					FiatQuantity:         price.StringFixed(quoteCurrency.GetDecimalPlace()),
					FiatQuantityCurrency: currencyPair.Quote,
					Price:                price.StringFixed(quoteCurrency.GetDecimalPlace()),
					PriceCurrency:        currencyPair.Quote,
					FiatPrice:            price.Sub(fee).StringFixed(quoteCurrency.GetDecimalPlace()),
					FiatPriceCurrency:    currencyPair.Quote,
					Fee:                  fee.StringFixed(quoteCurrency.GetDecimalPlace()),
					FeeCurrency:          currencyPair.Quote,
					FiatFee:              decimal.NewFromFloat(order.FillFees).StringFixed(quoteCurrency.GetDecimalPlace()),
					FiatFeeCurrency:      currencyPair.Quote,
					Total:                quantity,
					TotalCurrency:        currencyPair.Base,
					FiatTotal:            total.StringFixed(quoteCurrency.GetDecimalPlace()),
					FiatTotalCurrency:    currencyPair.Quote})
			}
		}
	}
	return orders
}

func (_gdax *GDAX) GetDepositHistory() ([]common.Transaction, error) {
	var deposits []common.Transaction
	txs, err := _gdax.getDepositWithdrawalHistory()
	if err != nil {
		return nil, err
	}
	for _, order := range txs {
		if order.GetType() == common.DEPOSIT_ORDER_TYPE {
			deposits = append(deposits, order)
		}
	}
	return deposits, nil
}

func (_gdax *GDAX) GetWithdrawalHistory() ([]common.Transaction, error) {
	var deposits []common.Transaction
	txs, err := _gdax.getDepositWithdrawalHistory()
	if err != nil {
		return nil, err
	}
	for _, order := range txs {
		if order.GetType() == common.WITHDRAWAL_ORDER_TYPE {
			deposits = append(deposits, order)
		}
	}
	return deposits, nil
}

func (_gdax *GDAX) getDepositWithdrawalHistory() ([]common.Transaction, error) {
	cacheKey := fmt.Sprintf("%d-%s", _gdax.ctx.GetUser().GetId(), "-gdax-depositWithdrawals")
	if x, found := _gdax.cache.Get(cacheKey); found {
		_gdax.ctx.GetLogger().Debugf("[GDAX.getDepositWithdrawalHistory] Returning desposit/withdraw history from cache")
		orders := x.(*[]common.Transaction)
		cachedTransactions := make([]common.Transaction, len(*orders))
		for i, order := range *orders {
			cachedTransactions[i] = order
		}
		return cachedTransactions, nil
	}
	GDAX_RATELIMITER.RespectRateLimit()
	_gdax.logger.Debug("[GDAX.getDepositWithdrawalHistory] Getting deposit/withdraw history")
	var orders []common.Transaction
	var ledger []gdax.LedgerEntry
	var transfers []gdax.LedgerEntry
	matches := make(map[float64]gdax.LedgerEntry, 10)
	accounts, err := _gdax.gdax.GetAccounts()
	if err != nil {
		_gdax.logger.Errorf("[GDAX.getDepositWithdrawalHistory] %s", err.Error())
		return nil, err
	}
	for _, a := range accounts {
		cursor := _gdax.gdax.ListAccountLedger(a.Id)
		for cursor.HasMore {
			if err := cursor.NextPage(&ledger); err != nil {
				_gdax.logger.Errorf("[GDAX.getDepositWithdrawalHistory] %s", err.Error())
				return nil, err
			}
			for _, order := range ledger {
				if order.Type == "match" {
					matches[order.Amount] = order
					continue
				}
				if order.Type == "transfer" {
					transfers = append(transfers, order)
				}
			}
		}
	}
	for _, order := range transfers {
		var id, txType string
		var price, fee, total decimal.Decimal
		var purchaseOrder *gdax.Order
		quantity := decimal.NewFromFloat(order.Amount)
		fee = decimal.NewFromFloat(0.0)
		price = quantity
		total = quantity
		amountKey := math.Abs(order.Amount)
		transferMatch := matches[amountKey]
		var currencyPair *common.CurrencyPair
		if transferMatch.Id <= 0 {
			_gdax.logger.Errorf("[GDAX.getDepositWithdrawalHistory] Unable to locate a match order for amount %f", amountKey)
			localCurrency := _gdax.ctx.GetUser().GetLocalCurrency()
			cp, _ := common.NewCurrencyPair(fmt.Sprintf("%s-%s", localCurrency, localCurrency), localCurrency)
			currencyPair = cp
		} else {
			po, err := _gdax.gdax.GetOrder(transferMatch.Details.OrderId)
			if err != nil {
				_gdax.ctx.GetLogger().Errorf("[GDAX.GetOrderHistory] Error retrieving purchase order: %s", err.Error())
				continue
			} else {
				purchaseOrder = &po
			}
			cp, err := common.NewCurrencyPair(transferMatch.Details.ProductId, _gdax.ctx.GetUser().GetLocalCurrency())
			if err != nil {
				_gdax.logger.Errorf("[GDAX.getDepositWithdrawalHistory] Error %s", err.Error())
			}
			currencyPair = cp
		}
		if purchaseOrder != nil {
			quantity = decimal.NewFromFloat(purchaseOrder.FilledSize)
			fee = decimal.NewFromFloat(purchaseOrder.FillFees)
			total = decimal.NewFromFloat(purchaseOrder.ExecutedValue)
			price = total.Sub(fee)
		}
		baseCurrency, err := _gdax.getCurrency(currencyPair.Base)
		if err != nil {
			_gdax.logger.Errorf("[GDAX.getDepositWithdrawalHistory] Error getting base currency: %s", err.Error())
			continue
		}
		quoteCurrency, err := _gdax.getCurrency(currencyPair.Quote)
		if err != nil {
			_gdax.logger.Errorf("[GDAX.getDepositWithdrawalHistory] Error getting quote currency: %s", err.Error())
			continue
		}
		if order.Amount < 0 {
			txType = common.WITHDRAWAL_ORDER_TYPE
		} else {
			txType = common.DEPOSIT_ORDER_TYPE
		}
		if order.Id <= 0 {
			id = fmt.Sprintf("%d", order.Id)
		} else {
			id = fmt.Sprintf("%s", order.CreatedAt.Time())
		}
		orders = append(orders, &dto.TransactionDTO{
			Id:                   id,
			Type:                 txType,
			Date:                 order.CreatedAt.Time(),
			Network:              _gdax.name,
			NetworkDisplayName:   _gdax.displayName,
			CurrencyPair:         currencyPair,
			Quantity:             quantity.StringFixed(baseCurrency.GetDecimalPlace()),
			QuantityCurrency:     currencyPair.Base,
			FiatQuantity:         total.Sub(fee).StringFixed(2),
			FiatQuantityCurrency: currencyPair.Quote,
			Price:                price.StringFixed(quoteCurrency.GetDecimalPlace()),
			PriceCurrency:        currencyPair.Quote,
			FiatPrice:            price.StringFixed(2),
			FiatPriceCurrency:    currencyPair.Quote,
			Fee:                  fee.StringFixed(quoteCurrency.GetDecimalPlace()),
			FeeCurrency:          currencyPair.Quote,
			FiatFee:              fee.StringFixed(2),
			FiatFeeCurrency:      currencyPair.Quote,
			Total:                decimal.NewFromFloat(order.Amount).StringFixed(baseCurrency.GetDecimalPlace()),
			TotalCurrency:        currencyPair.Base,
			FiatTotal:            total.StringFixed(2),
			FiatTotalCurrency:    currencyPair.Quote})
	}

	_gdax.cache.Set(cacheKey, &orders, cache.DefaultExpiration)
	return orders, nil
}

func (_gdax *GDAX) GetCurrencies() (map[string]*common.Currency, error) {
	cacheKey := fmt.Sprintf("%d-%s", _gdax.ctx.GetUser().GetId(), "-gdax-currencies")
	if x, found := _gdax.cache.Get(cacheKey); found {
		//_gdax.ctx.GetLogger().Debugf("[GDAX.GetCurrencies] Returning GDAX currencies from cache")
		currencies := x.(*map[string]*common.Currency)
		return *currencies, nil
	}
	_gdax.ctx.GetLogger().Debugf("[GDAX.GetCurrencies] Retrieving GDAX currencies")
	currencies, err := _gdax.gdax.GetCurrencies()
	if err != nil {
		_gdax.ctx.GetLogger().Errorf("[GDAX.GetCurrencies] Error: %s", err.Error())
		return nil, err
	}
	_currencies := make(map[string]*common.Currency, len(currencies))
	for _, currency := range currencies {
		if fiatCurrency, found := common.FiatCurrencies[currency.Id]; found {
			_currencies[currency.Id] = fiatCurrency
			continue
		}
		_currencies[currency.Id] = &common.Currency{
			ID:           currency.Id,
			Name:         currency.Name,
			Symbol:       currency.Id,
			BaseUnit:     100000000,
			TxFee:        decimal.NewFromFloat(currency.MinSize),
			DecimalPlace: util.ParseDecimalPlace(decimal.NewFromFloat(currency.MinSize).String())}
	}
	_gdax.cache.Set(cacheKey, &_currencies, cache.DefaultExpiration)
	return _currencies, nil
}

func (_gdax *GDAX) getCurrency(currency string) (*common.Currency, error) {
	currencies, err := _gdax.GetCurrencies()
	if err != nil {
		return nil, err
	}
	if currency, found := currencies[currency]; found {
		return currency, nil
	}
	return nil, errors.New(fmt.Sprintf("Currency not found: %s", currency))
}

func (_gdax *GDAX) GetBalances() ([]common.Coin, float64) {
	var cachedBalances []common.Coin
	var cachedSum float64
	balancesKey := fmt.Sprintf("%d-%s", _gdax.ctx.GetUser().GetId(), "balances")
	sumKey := fmt.Sprintf("%d-%s", _gdax.ctx.GetUser().GetId(), "balances-sum")
	if x, found := _gdax.cache.Get(balancesKey); found {
		balances := x.(*[]common.Coin)
		cachedBalances := make([]common.Coin, len(*balances))
		for i, balance := range *balances {
			cachedBalances[i] = balance
		}
	}
	if x, found := _gdax.cache.Get(sumKey); found {
		cachedSum = x.(float64)
	}
	if len(cachedBalances) > 0 && cachedSum > 0 {
		return cachedBalances, cachedSum
	}
	GDAX_RATELIMITER.RespectRateLimit()
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
		if a.Currency != _gdax.ctx.GetUser().GetLocalCurrency() {
			currency := fmt.Sprintf("%s-%s", a.Currency, _gdax.ctx.GetUser().GetLocalCurrency())
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
	_gdax.cache.Set(balancesKey, &coins, cache.DefaultExpiration)
	_gdax.cache.Set(sumKey, sum, cache.DefaultExpiration)
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

func (_gdax *GDAX) GetSummary() common.CryptoExchangeSummary {
	total := 0.0
	satoshis := 0.0
	balances, _ := _gdax.GetBalances()
	for _, c := range balances {
		if c.GetCurrency() == _gdax.ctx.GetUser().GetLocalCurrency() {
			total += c.GetTotal()
		} else if c.IsBitcoin() {
			satoshis += c.GetBalance()
			total += c.GetTotal()
		} else {
			currency := fmt.Sprintf("%s-BTC", c.GetCurrency())
			GDAX_RATELIMITER.RespectRateLimit()
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
	exchange := &dto.CryptoExchangeSummaryDTO{
		Name:     _gdax.name,
		URL:      "https://www.gdax.com",
		Total:    t,
		Satoshis: s,
		Coins:    balances}
	return exchange
}

func (_gdax *GDAX) ParseImport(file string) ([]common.Transaction, error) {
	var orders []common.Transaction
	_gdax.ctx.GetLogger().Error("[GDAX.ParseImport] Unsupported!")
	return orders, errors.New("GDAX.ParseImport Unsupported")
}

func (_gdax *GDAX) GetName() string {
	return _gdax.name
}

func (_gdax *GDAX) GetDisplayName() string {
	return _gdax.displayName
}

func (_gdax *GDAX) GetTradingFee() float64 {
	return _gdax.tradingFee
}

func (_gdax *GDAX) FormattedCurrencyPair(currencyPair *common.CurrencyPair) string {
	return fmt.Sprintf("%s-%s", currencyPair.Base, currencyPair.Quote)
}

func (_gdax *GDAX) respectRateLimit() {
	_gdax.rateLimiter.RespectRateLimit()
}

func (_gdax *GDAX) parseDecimalPlaces(minSize string) int {
	pieces := strings.Split(minSize, ".")
	idx := strings.Index(pieces[1], "1")
	if idx > -1 {
		return idx + 1
	}
	_gdax.ctx.GetLogger().Errorf("[GDAX.minSizeDecimals] Error: Unable to locate decimal market")
	return 0
}

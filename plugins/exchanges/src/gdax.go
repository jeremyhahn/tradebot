package main

import (
	"errors"
	"fmt"
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
	tradingFee  decimal.Decimal
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
		tradingFee:  decimal.NewFromFloat(0.025),
		cache:       cache.New(1*time.Minute, 1*time.Minute)}
}

func (_gdax *GDAX) GetPriceAt(currency string, atDate time.Time) (*common.Candlestick, error) {
	currencyPair := &common.CurrencyPair{
		Base:          currency,
		Quote:         _gdax.ctx.GetUser().GetLocalCurrency(),
		LocalCurrency: _gdax.ctx.GetUser().GetLocalCurrency()}
	kline, err := _gdax.GetPriceHistory(currencyPair, atDate.Add(-5*time.Minute), atDate.Add(5*time.Minute), 60)
	if err != nil {
		return &common.Candlestick{}, err
	}
	closestCandle, err := util.FindClosestDatedCandle(_gdax.ctx.GetLogger(), atDate, kline)
	if err != nil {
		return &common.Candlestick{}, err
	}
	return closestCandle, nil
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
			Open:         decimal.NewFromFloat(r.Open),
			Close:        decimal.NewFromFloat(r.Close),
			High:         decimal.NewFromFloat(r.High),
			Low:          decimal.NewFromFloat(r.Low),
			Volume:       decimal.NewFromFloat(r.Volume)})
	}
	return candlesticks, nil
}

func (_gdax *GDAX) GetOrderHistory(currencyPair *common.CurrencyPair) []common.Transaction {
	cacheKey := fmt.Sprintf("%d-%s", _gdax.ctx.GetUser().GetId(), "-gdax-orderhistory")
	if txs, found := _gdax.cache.Get(cacheKey); found {
		_gdax.ctx.GetLogger().Debugf("[GDAX.GetOrderHistory] Returning %s's GDAX orders from cache", _gdax.ctx.GetUser().GetUsername())
		transactions := txs.(*[]common.Transaction)
		return *transactions
	}
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
				quantity := decimal.NewFromFloat(order.FilledSize)
				fee := decimal.NewFromFloat(order.FillFees)
				total := decimal.NewFromFloat(order.ExecutedValue)
				price := total.Div(quantity)
				orders = append(orders, &dto.TransactionDTO{
					Id:                   strconv.FormatInt(int64(e.Id), 10),
					Type:                 order.Side,
					Date:                 orderDate,
					Network:              _gdax.name,
					NetworkDisplayName:   _gdax.displayName,
					CurrencyPair:         currencyPair,
					Quantity:             quantity.StringFixed(baseCurrency.GetDecimalPlace()),
					QuantityCurrency:     currencyPair.Base,
					FiatQuantity:         total.StringFixed(quoteCurrency.GetDecimalPlace()),
					FiatQuantityCurrency: currencyPair.Quote,
					Price:                price.StringFixed(quoteCurrency.GetDecimalPlace()),
					PriceCurrency:        currencyPair.Quote,
					FiatPrice:            price.StringFixed(quoteCurrency.GetDecimalPlace()),
					FiatPriceCurrency:    currencyPair.Quote,
					Fee:                  fee.StringFixed(quoteCurrency.GetDecimalPlace()),
					FeeCurrency:          currencyPair.Quote,
					FiatFee:              decimal.NewFromFloat(order.FillFees).StringFixed(quoteCurrency.GetDecimalPlace()),
					FiatFeeCurrency:      currencyPair.Quote,
					Total:                quantity.StringFixed(baseCurrency.GetDecimalPlace()),
					TotalCurrency:        currencyPair.Base,
					FiatTotal:            total.StringFixed(quoteCurrency.GetDecimalPlace()),
					FiatTotalCurrency:    currencyPair.Quote})
			}
		}
	}
	_gdax.cache.Set(cacheKey, &orders, cache.DefaultExpiration)
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
	var withdrawals []common.Transaction
	txs, err := _gdax.getDepositWithdrawalHistory()
	if err != nil {
		return nil, err
	}
	for _, order := range txs {
		if order.GetType() == common.WITHDRAWAL_ORDER_TYPE {
			withdrawals = append(withdrawals, order)
		}
	}
	return withdrawals, nil
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
	_gdax.logger.Debug("[GDAX.getDepositWithdrawalHistory] Getting deposit/withdrawal history")
	var orders []common.Transaction
	var ledger []gdax.LedgerEntry
	var transfers []gdax.LedgerEntry
	buys := make(map[string]gdax.Order, 10)
	matches := make(map[string]gdax.LedgerEntry, 10)
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
					matches[fmt.Sprintf("%.8f", order.Amount)] = order
					o, err := _gdax.gdax.GetOrder(order.Details.OrderId)
					if err != nil {
						_gdax.ctx.GetLogger().Errorf("[GDAX.getDepositWithdrawalHistory] Error retrieving order: %s", err.Error())
					}
					if o.Side == "buy" {
						buys[fmt.Sprintf("%.8f", o.FilledSize)] = o
					}
					continue
				}
				if order.Type == "transfer" {
					order.Details.ProductId = fmt.Sprintf("%s-%s", a.Currency, a.Currency)
					transfers = append(transfers, order)
				}
			}
		}
	}
	for _, order := range transfers {
		var id, txType, priceCurrency, feeCurrency string
		var price, fee, total decimal.Decimal
		var purchaseOrder *gdax.Order
		var currencyPair *common.CurrencyPair
		quantity := decimal.NewFromFloat(order.Amount)
		fee = decimal.NewFromFloat(0)
		price = quantity
		total = quantity
		decAmount := decimal.NewFromFloat(order.Amount)
		amountKey := decAmount.Abs().String()
		transferMatch := matches[amountKey]
		orderDate := order.CreatedAt.Time()

		if transferMatch.Id <= 0 {
			_gdax.logger.Errorf("[GDAX.getDepositWithdrawalHistory] Unable to locate a match order for %s transfer amount %s. Attempting match against prior buys.",
				order.Details.ProductId, amountKey)
			localCurrency := _gdax.ctx.GetUser().GetLocalCurrency()
			transferCurrencyPair, _ := common.NewCurrencyPair(order.Details.ProductId, localCurrency)
			currencyPair = transferCurrencyPair
			quantity = decimal.NewFromFloat(order.Amount)
			total = decimal.NewFromFloat(order.Amount)
			if buy, ok := buys[quantity.Abs().String()]; ok {
				buyCurrencyPair, _ := common.NewCurrencyPair(buy.ProductId, localCurrency)
				fillSize := decimal.NewFromFloat(buy.FilledSize)
				price = decimal.NewFromFloat(buy.ExecutedValue).Div(fillSize)
				baseFiatPrice, err := _gdax.GetPriceAt(transferCurrencyPair.Base, orderDate)
				if err != nil {
					_gdax.logger.Errorf("[GDAX.getDepositWithdrawalHistory] Error getting %s fiat price", currencyPair)
				}
				fee = decimal.NewFromFloat(buy.FillFees).Abs()
				if order.Amount < 0 {
					txType = common.WITHDRAWAL_ORDER_TYPE
				} else {
					txType = common.DEPOSIT_ORDER_TYPE
				}
				buyQuoteCurrency, err := _gdax.getCurrency(buyCurrencyPair.Quote)
				if err != nil {
					_gdax.logger.Errorf("[GDAX.getDepositWithdrawalHistory] Error getting buy quote currency: %s", err.Error())
				}
				total = baseFiatPrice.Close.Mul(quantity.Abs())
				orders = append(orders, &dto.TransactionDTO{
					Id:                   buy.Id,
					Type:                 txType,
					Date:                 orderDate,
					Network:              _gdax.name,
					NetworkDisplayName:   _gdax.displayName,
					CurrencyPair:         transferCurrencyPair,
					Quantity:             quantity.StringFixed(8),
					QuantityCurrency:     transferCurrencyPair.Base,
					FiatQuantity:         total.StringFixed(buyQuoteCurrency.GetDecimalPlace()),
					FiatQuantityCurrency: buyCurrencyPair.Quote,
					Price:                baseFiatPrice.Close.StringFixed(buyQuoteCurrency.GetDecimalPlace()),
					PriceCurrency:        buyQuoteCurrency.GetID(),
					FiatPrice:            price.StringFixed(buyQuoteCurrency.GetDecimalPlace()),
					FiatPriceCurrency:    buyQuoteCurrency.GetID(),
					Fee:                  fee.StringFixed(buyQuoteCurrency.GetDecimalPlace()),
					FeeCurrency:          buyQuoteCurrency.GetID(),
					FiatFee:              fee.StringFixed(buyQuoteCurrency.GetDecimalPlace()),
					FiatFeeCurrency:      buyQuoteCurrency.GetID(),
					Total:                quantity.StringFixed(8),
					TotalCurrency:        currencyPair.Base,
					FiatTotal:            total.StringFixed(buyQuoteCurrency.GetDecimalPlace()),
					FiatTotalCurrency:    buyQuoteCurrency.GetID()})
				continue
			}
			if _, ok := common.FiatCurrencies[currencyPair.Base]; ok {
				price = quantity
			} else {
				candle, err := _gdax.GetPriceAt(currencyPair.Base, orderDate)
				if err != nil {
					_gdax.ctx.GetLogger().Errorf("[GDAX.GetOrderHistory] Error retrieving %s price on %s: %s",
						currencyPair.Base, orderDate, err.Error())
				}
				price = candle.Close
				total = price.Mul(quantity.Abs())
				if order.Id <= 0 {
					id = fmt.Sprintf("%d", order.Id)
				} else {
					id = fmt.Sprintf("%s", orderDate)
				}
				orders = append(orders, &dto.TransactionDTO{
					Id:                   id,
					Type:                 txType,
					Date:                 orderDate,
					Network:              _gdax.name,
					NetworkDisplayName:   _gdax.displayName,
					CurrencyPair:         transferCurrencyPair,
					Quantity:             quantity.StringFixed(8),
					QuantityCurrency:     transferCurrencyPair.Base,
					FiatQuantity:         total.StringFixed(2),
					FiatQuantityCurrency: "USD",
					Price:                price.StringFixed(2),
					PriceCurrency:        "USD",
					FiatPrice:            price.StringFixed(2),
					FiatPriceCurrency:    "USD",
					Fee:                  fee.StringFixed(2),
					FeeCurrency:          "USD",
					FiatFee:              fee.StringFixed(2),
					FiatFeeCurrency:      "USD",
					Total:                quantity.StringFixed(8),
					TotalCurrency:        currencyPair.Base,
					FiatTotal:            total.StringFixed(2),
					FiatTotalCurrency:    "USD"})
				continue
			}
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
			price = decimal.NewFromFloat(purchaseOrder.Price)
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
			id = fmt.Sprintf("%s", orderDate)
		}
		if _, ok := common.FiatCurrencies[currencyPair.Base]; !ok {
			candle, _ := _gdax.GetPriceAt(currencyPair.Base, orderDate)
			price = candle.Close
		}
		if priceCurrency == "" {
			priceCurrency = currencyPair.Quote
		}
		if feeCurrency == "" {
			feeCurrency = currencyPair.Quote
		}
		orders = append(orders, &dto.TransactionDTO{
			Id:                   id,
			Type:                 txType,
			Date:                 orderDate,
			Network:              _gdax.name,
			NetworkDisplayName:   _gdax.displayName,
			CurrencyPair:         currencyPair,
			Quantity:             quantity.StringFixed(baseCurrency.GetDecimalPlace()),
			QuantityCurrency:     currencyPair.Base,
			FiatQuantity:         total.Sub(fee).StringFixed(2),
			FiatQuantityCurrency: currencyPair.Quote,
			Price:                price.StringFixed(quoteCurrency.GetDecimalPlace()),
			PriceCurrency:        priceCurrency,
			FiatPrice:            price.StringFixed(quoteCurrency.GetDecimalPlace()),
			FiatPriceCurrency:    currencyPair.Quote,
			Fee:                  fee.StringFixed(quoteCurrency.GetDecimalPlace()),
			FeeCurrency:          currencyPair.Quote,
			FiatFee:              fee.StringFixed(quoteCurrency.GetDecimalPlace()),
			FiatFeeCurrency:      currencyPair.Quote,
			Total:                decimal.NewFromFloat(order.Amount).StringFixed(baseCurrency.GetDecimalPlace()),
			TotalCurrency:        currencyPair.Base,
			FiatTotal:            total.StringFixed(quoteCurrency.GetDecimalPlace()),
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

func (_gdax *GDAX) GetBalances() ([]common.Coin, decimal.Decimal) {
	var cachedBalances []common.Coin
	var cachedSum decimal.Decimal
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
		cachedSum = x.(decimal.Decimal)
	}
	if len(cachedBalances) > 0 && cachedSum.GreaterThan(decimal.NewFromFloat(0)) {
		return cachedBalances, cachedSum
	}
	GDAX_RATELIMITER.RespectRateLimit()
	_gdax.logger.Debugf("[GDAX] Getting balances")
	var coins []common.Coin
	sum := decimal.NewFromFloat(0)
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
		var decimalPlaces int32
		if _, exists := common.FiatCurrencies[a.Currency]; exists {
			decimalPlaces = 2
		} else {
			decimalPlaces = 8
		}
		decPrice := decimal.NewFromFloat(price)
		total := decimal.NewFromFloat(a.Balance).Mul(decPrice)
		sum = sum.Add(total)
		coins = append(coins, &dto.CoinDTO{
			Currency:  a.Currency,
			Balance:   decimal.NewFromFloat(a.Balance).Truncate(decimalPlaces),
			Available: decimal.NewFromFloat(a.Available).Truncate(decimalPlaces),
			Pending:   decimal.NewFromFloat(a.Hold).Truncate(decimalPlaces),
			Price:     decPrice.Truncate(decimalPlaces),
			Total:     total.Truncate(decimalPlaces)})
	}
	_gdax.cache.Set(balancesKey, &coins, cache.DefaultExpiration)
	_gdax.cache.Set(sumKey, sum, cache.DefaultExpiration)
	return coins, sum.Truncate(2)
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
				Price:        decimal.NewFromFloat(message.Price)}
		}
	}

	_gdax.SubscribeToLiveFeed(currencyPair, priceChannel)
}

func (_gdax *GDAX) GetSummary() common.CryptoExchangeSummary {
	total := decimal.NewFromFloat(0)
	satoshis := decimal.NewFromFloat(0)
	balances, _ := _gdax.GetBalances()
	for _, c := range balances {
		if c.GetCurrency() == _gdax.ctx.GetUser().GetLocalCurrency() {
			total = total.Add(c.GetTotal())
		} else if c.IsBitcoin() {
			satoshis = satoshis.Add(c.GetBalance())
			total = total.Add(c.GetTotal())
		} else {
			currency := fmt.Sprintf("%s-BTC", c.GetCurrency())
			GDAX_RATELIMITER.RespectRateLimit()
			ticker, err := _gdax.gdax.GetTicker(currency)
			if err != nil {
				_gdax.logger.Errorf("[GDAX.GetExchange] %s", err.Error())
				continue
			}
			satoshis = satoshis.Add(decimal.NewFromFloat(ticker.Price))
			total = total.Add(c.GetTotal())
		}
	}
	exchange := &dto.CryptoExchangeSummaryDTO{
		Name:     _gdax.name,
		URL:      "https://www.gdax.com",
		Total:    total.Truncate(2),
		Satoshis: satoshis.Truncate(8),
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

func (_gdax *GDAX) GetTradingFee() decimal.Decimal {
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

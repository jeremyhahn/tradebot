package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dto"
	"github.com/jeremyhahn/tradebot/entity"
	"github.com/jeremyhahn/tradebot/util"
	logging "github.com/op/go-logging"
	cache "github.com/patrickmn/go-cache"
	"github.com/shopspring/decimal"
	bittrex "github.com/toorop/go-bittrex"
	"golang.org/x/text/encoding/unicode"
)

type Bittrex struct {
	ctx         common.Context
	client      *bittrex.Bittrex
	logger      *logging.Logger
	name        string
	displayName string
	tradingFee  decimal.Decimal
	usdtMarkets []string
	marketPairs string
	cache       *cache.Cache
	common.Exchange
}

var BITTREX_RATE_LIMITER = common.NewRateLimiter(1, 1)

func CreateBittrex(ctx common.Context, userExchangeEntity entity.UserExchangeEntity) common.Exchange {
	return &Bittrex{
		ctx:         ctx,
		client:      bittrex.New(userExchangeEntity.GetKey(), userExchangeEntity.GetSecret()),
		logger:      ctx.GetLogger(),
		name:        "Bittrex",
		displayName: "Bittrex",
		tradingFee:  decimal.NewFromFloat(.025),
		marketPairs: userExchangeEntity.GetExtra(),
		usdtMarkets: []string{"BTC", "ETH", "XRP", "NEO", "ADA", "LTC",
			"BCC", "OMG", "ETC", "ZEC", "XMR", "XVG", "BTG", "DASH", "NXT"},
		cache: cache.New(1*time.Minute, 1*time.Minute)}
}

func (b *Bittrex) GetPriceAt(currency string, atDate time.Time) (*common.Candlestick, error) {
	var candles []common.Candlestick
	marketPair := &common.CurrencyPair{
		Base:          "USDT",
		Quote:         currency,
		LocalCurrency: b.ctx.GetUser().GetLocalCurrency()}
	monthAgo := time.Now().AddDate(0, -1, 0)
	if atDate.Before(monthAgo) {
		atDate := time.Date(atDate.Year(), atDate.Month(), atDate.Day(), 0, 0, 0, 0, atDate.Location())
		klines, err := b.GetPriceHistory(marketPair, atDate.Add(-24*time.Hour), atDate.Add(24*time.Hour), 1000)
		if err != nil {
			return &common.Candlestick{}, err
		}
		candles = klines
	} else {
		klines, err := b.GetPriceHistory(marketPair, atDate.Add(-5*time.Minute), atDate.Add(5*time.Minute), 1)
		if err != nil {
			return &common.Candlestick{}, err
		}
		candles = klines
	}
	closestCandle, err := util.FindClosestDatedCandle(b.ctx.GetLogger(), atDate, candles)
	if err != nil {
		return &common.Candlestick{}, err
	}
	return closestCandle, nil
}

func (b *Bittrex) SubscribeToLiveFeed(marketPair *common.CurrencyPair, priceChange chan common.PriceChange) {
	for {
		symbol := b.FormattedCurrencyPair(marketPair)
		time.Sleep(10 * time.Second)
		ticker, err := b.client.GetTicker(symbol)
		if err != nil {
			b.logger.Errorf("[Bittrex.SubscribeToLiveFeed] %s", err.Error())
			continue
		}
		b.logger.Debugf("[Bittrex.SubscribeToLiveFeed] Sending live price: %s", ticker.Last.StringFixed(8))
		priceChange <- common.PriceChange{
			CurrencyPair: marketPair,
			Satoshis:     ticker.Last,
			Price:        ticker.Last}
	}
}

func (b *Bittrex) GetPriceHistory(marketPair *common.CurrencyPair,
	start, end time.Time, granularity int) ([]common.Candlestick, error) {
	BITTREX_RATE_LIMITER.RespectRateLimit()
	b.logger.Debugf("[Bittrex.GetPriceHistory] Getting %s price history for %s between %s and %s",
		marketPair, b.ctx.GetUser().GetUsername(), start, end)
	interval := "hour"
	switch {
	case granularity == 1:
		interval = "oneMin"
	case granularity == 5:
		interval = "fiveMin"
	case granularity == 30:
		interval = "thirtyMin"
	case granularity == 100:
		interval = "hour"
	case granularity == 1000:
		interval = "day"
	}
	candlesticks := make([]common.Candlestick, 0)
	ticks, err := b.client.GetTicks(b.FormattedCurrencyPair(marketPair), interval)
	if err != nil {
		b.logger.Errorf("[Bittrex.GetPriceHistory] Error: %s", err.Error())
		return nil, err
	}
	for _, tick := range ticks {
		timestamp := tick.TimeStamp.Time
		if timestamp.Before(start) || timestamp.After(end) {
			continue
		}
		candlesticks = append(candlesticks, common.Candlestick{
			Date:   timestamp,
			Open:   tick.Open,
			High:   tick.High,
			Low:    tick.Low,
			Close:  tick.Close,
			Volume: tick.Volume})
	}
	return candlesticks, nil
}

func (b *Bittrex) GetOrderHistory(marketPair *common.CurrencyPair) []common.Transaction {
	BITTREX_RATE_LIMITER.RespectRateLimit()
	formattedCurrencyPair := b.FormattedCurrencyPair(marketPair)
	b.logger.Debugf("[Bittrex.GetOrderHistory] Getting %s order history", formattedCurrencyPair)
	var orders []common.Transaction
	orderHistory, err := b.client.GetOrderHistory(formattedCurrencyPair)
	if err != nil {
		b.logger.Errorf("[Bittrex.GetOrderHistory] Error: %s", err.Error())
	}
	if len(orderHistory) == 0 {
		b.logger.Warning("[Bittrex.GetOrderHistory] Zero records returned from Bittrex API")
	}
	for _, o := range orderHistory {
		orderDate := o.TimeStamp.Time
		qty := o.Quantity
		price := o.Price
		fee := o.Commission
		total := qty.Mul(price)
		orders = append(orders, &dto.TransactionDTO{
			Id:                 o.OrderUuid,
			Network:            b.name,
			NetworkDisplayName: b.displayName,
			Date:               orderDate,
			Type:               o.OrderType,
			Category:           common.TX_CATEGORY_TRADE,
			MarketPair:         marketPair,
			CurrencyPair: &common.CurrencyPair{
				Base:          marketPair.Quote,
				Quote:         marketPair.Base,
				LocalCurrency: b.ctx.GetUser().GetLocalCurrency()},
			Quantity:               qty.StringFixed(8),
			QuantityCurrency:       marketPair.Quote,
			FiatQuantity:           "0.00",
			FiatQuantityCurrency:   "N/A",
			Price:                  price.StringFixed(8),
			PriceCurrency:          marketPair.Base,
			FiatPrice:              "0.00",
			FiatPriceCurrency:      "N/A",
			QuoteFiatPrice:         "0.00",
			QuoteFiatPriceCurrency: "N/A",
			Fee:               fee.StringFixed(8),
			FeeCurrency:       marketPair.Base,
			FiatFee:           "0.00",
			FiatFeeCurrency:   "N/A",
			Total:             total.StringFixed(8),
			TotalCurrency:     marketPair.Base,
			FiatTotal:         "0.00",
			FiatTotalCurrency: "N/A"})
	}
	return orders
}

func (b *Bittrex) GetDepositHistory() ([]common.Transaction, error) {
	b.logger.Debugf("[Bittrex.GetDepositHistory] Getting deposits for %s", b.ctx.GetUser().GetUsername())
	var _deposits []common.Transaction
	deposits, err := b.client.GetDepositHistory("all")
	if err != nil {
		b.ctx.GetLogger().Errorf("[Bittrex.GetDepositHistory] Error: %s", err.Error())
		return _deposits, err
	}
	for _, deposit := range deposits {
		marketPair := &common.CurrencyPair{
			Base:          deposit.Currency,
			Quote:         deposit.Currency,
			LocalCurrency: b.ctx.GetUser().GetLocalCurrency()}
		orderDate := deposit.LastUpdated.Time
		quantity := deposit.Amount
		baseFiatPrice := b.getFiatPrice(marketPair.Base, orderDate)
		fiatTotal := quantity.Mul(baseFiatPrice)
		_deposits = append(_deposits, &dto.TransactionDTO{
			Id:                 deposit.TxId,
			Type:               common.TX_CATEGORY_DEPOSIT,
			Category:           common.TX_CATEGORY_DEPOSIT,
			Date:               orderDate,
			Network:            b.name,
			NetworkDisplayName: b.displayName,
			MarketPair:         marketPair,
			CurrencyPair: &common.CurrencyPair{
				Base:          marketPair.Quote,
				Quote:         marketPair.Base,
				LocalCurrency: b.ctx.GetUser().GetLocalCurrency()},
			Quantity:               quantity.StringFixed(8),
			QuantityCurrency:       deposit.Currency,
			FiatQuantity:           fiatTotal.StringFixed(2),
			FiatQuantityCurrency:   "USD",
			Price:                  baseFiatPrice.StringFixed(2),
			PriceCurrency:          "USD",
			FiatPrice:              baseFiatPrice.StringFixed(2),
			FiatPriceCurrency:      "USD",
			QuoteFiatPrice:         baseFiatPrice.StringFixed(2),
			QuoteFiatPriceCurrency: "USD",
			Fee:               "0.00000000",
			FeeCurrency:       deposit.Currency,
			FiatFee:           "0.00",
			FiatFeeCurrency:   "USD",
			Total:             fiatTotal.StringFixed(2),
			TotalCurrency:     "USD",
			FiatTotal:         fiatTotal.StringFixed(2),
			FiatTotalCurrency: "USD"})
	}
	b.ctx.GetLogger().Debugf("[Bittrex.GetDepositHistory] Retrieved %d deposits", len(_deposits))
	return _deposits, nil
}

func (b *Bittrex) GetWithdrawalHistory() ([]common.Transaction, error) {
	var orders []common.Transaction
	currencies, _ := b.GetCurrencies()
	b.ctx.GetLogger().Debugf("[Bittrex.GetWithdrawalHistory] Withdrawal currencies for %s: %s", b.ctx.GetUser().GetUsername(), currencies)
	for _, currency := range currencies {
		BITTREX_RATE_LIMITER.RespectRateLimit()
		withdrawals, err := b.client.GetWithdrawalHistory(currency.GetSymbol())
		if err != nil {
			b.logger.Errorf("[Bittrex.GetWithdrawalHistory] Error: %s", err.Error())
		}
		for _, withdrawal := range withdrawals {
			if !withdrawal.Authorized {
				continue
			}
			orderDate := withdrawal.Opened.Time
			marketPair := &common.CurrencyPair{
				Base:          withdrawal.Currency,
				Quote:         withdrawal.Currency,
				LocalCurrency: b.ctx.GetUser().GetLocalCurrency()}
			quantity := withdrawal.Amount
			txCost := withdrawal.TxCost
			baseFiatPrice := b.getFiatPrice(marketPair.Base, orderDate)
			total := quantity.Mul(baseFiatPrice)
			orders = append(orders, &dto.TransactionDTO{
				Id:                 withdrawal.TxId,
				Type:               common.TX_CATEGORY_WITHDRAWAL,
				Category:           common.TX_CATEGORY_WITHDRAWAL,
				Date:               orderDate,
				Network:            b.name,
				NetworkDisplayName: b.displayName,
				MarketPair:         marketPair,
				CurrencyPair: &common.CurrencyPair{
					Base:          marketPair.Quote,
					Quote:         marketPair.Base,
					LocalCurrency: b.ctx.GetUser().GetLocalCurrency()},
				Quantity:               quantity.StringFixed(8),
				QuantityCurrency:       withdrawal.Currency,
				FiatQuantity:           "0.00",
				FiatQuantityCurrency:   "N/A",
				Price:                  baseFiatPrice.StringFixed(2),
				PriceCurrency:          "USD",
				FiatPrice:              baseFiatPrice.StringFixed(2),
				FiatPriceCurrency:      "USD",
				QuoteFiatPrice:         baseFiatPrice.StringFixed(2),
				QuoteFiatPriceCurrency: "USD",
				Fee:               txCost.StringFixed(8),
				FeeCurrency:       withdrawal.Currency,
				FiatFee:           txCost.Mul(baseFiatPrice).StringFixed(2),
				FiatFeeCurrency:   "USD",
				Total:             quantity.StringFixed(8),
				TotalCurrency:     withdrawal.Currency,
				FiatTotal:         total.StringFixed(2),
				FiatTotalCurrency: "USD"})
		}
	}
	return orders, nil
}

func (b *Bittrex) GetBalances() ([]common.Coin, decimal.Decimal) {
	BITTREX_RATE_LIMITER.RespectRateLimit()
	var coins []common.Coin
	sum := decimal.NewFromFloat(0)
	balances, err := b.client.GetBalances()
	if err != nil {
		b.logger.Errorf("[Bittrex.GetBalances] %s", err.Error())
	}
	btcCurrencyPair := &common.CurrencyPair{
		Base:          "BTC",
		Quote:         b.ctx.GetUser().GetLocalCurrency(),
		LocalCurrency: b.ctx.GetUser().GetLocalCurrency()}
	localizedBtcCurrencyPair := b.localizedCurrencyPair(btcCurrencyPair)
	for _, bal := range balances {
		var currency string
		if bal.Currency == "BTC" {
			currency = fmt.Sprintf("%s-BTC", localizedBtcCurrencyPair.LocalCurrency)
		} else {
			currency = fmt.Sprintf("%s-%s", "BTC", bal.Currency)
		}
		b.logger.Debugf("[Bittrex.GetBalances] Getting %s ticker", currency)
		ticker, terr := b.client.GetTicker(currency)
		if terr != nil {
			b.logger.Errorf("[Bittrex.GetBalances] %s. currency: %s", terr.Error(), currency)
			continue
		}
		bitcoin, err := b.GetPriceAt("BTC", time.Now())
		if err != nil {
			b.logger.Errorf("[Bittrex.GetBalances] Error getting current bitcoin price: %s", err.Error())
		}
		if bal.Currency == "BTC" {
			total := bal.Balance.Mul(bitcoin.Close)
			sum = sum.Add(total)
			coins = append(coins, &dto.CoinDTO{
				Address:   bal.CryptoAddress,
				Available: bal.Balance.Truncate(8),
				Balance:   bal.Balance.Truncate(8),
				Currency:  bal.Currency,
				Pending:   bal.Pending.Truncate(8),
				Price:     bitcoin.Close.Truncate(2),
				Total:     bal.Balance.Mul(bitcoin.Close).Truncate(2)})
			continue
		} else {
			subtotal := ticker.Last.Mul(bitcoin.Close)
			total := bal.Balance.Mul(subtotal)
			sum = sum.Add(total)
			coins = append(coins, &dto.CoinDTO{
				Address:   bal.CryptoAddress,
				Available: bal.Available.Truncate(8),
				Balance:   bal.Balance.Truncate(8),
				Currency:  bal.Currency,
				Pending:   bal.Pending.Truncate(8),
				Price:     ticker.Last.Truncate(2),
				Total:     total.Truncate(2)})
		}
	}
	return coins, sum.Truncate(2)
}

func (b *Bittrex) GetCurrencies() (map[string]*common.Currency, error) {
	b.ctx.GetLogger().Debugf("[Bittrex.GetCurrencies] Getting currency list")
	currencies := make(map[string]*common.Currency)
	configuredCurrencies := strings.Split(b.marketPairs, ",")
	currencyMap := make(map[string]bool, len(configuredCurrencies))
	for _, cur := range configuredCurrencies {
		pieces := strings.Split(cur, "-")
		if len(pieces) < 2 {
			errmsg := fmt.Sprintf("Invalid currency pair: %s", cur)
			b.ctx.GetLogger().Errorf("[Bittrex.GetCurrencies] Error: %s", errmsg)
			return nil, errors.New(errmsg)
		} else {
			currencyMap[pieces[0]] = true
			currencyMap[pieces[1]] = true
		}
	}
	for k, _ := range currencyMap {
		var name string
		if _, found := common.CryptoCurrencies[k]; !found {
			b.ctx.GetLogger().Errorf("[Bittrex.GetCurrencies] Unable to locate currency %s in common.CryptoCurrencies", k)
			name = k
		} else {
			name = common.CryptoCurrencies[k]
		}
		currencies[k] = &common.Currency{
			ID:           k,
			Name:         name,
			Symbol:       k,
			BaseUnit:     100000000,
			DecimalPlace: 2}
	}
	return currencies, nil
	/*
		  BITTREX_RATE_LIMITER.RespectRateLimit()
			currencies, err := b.client.GetCurrencies()
			if err != nil {
				b.logger.Errorf("[Bittrex.GetCurrencies] %s", err.Error())
				return nil, err
			}
			_currencies := make(map[string]*common.Currency, len(currencies))
			for _, c := range currencies {
				_currencies[c.Currency] = &common.Currency{
					ID:           c.Currency,
					Symbol:       c.Currency,
					Name:         c.CurrencyLong,
					BaseUnit:     100000000,
					DecimalPlace: 2}
			}
			return _currencies, nil
	*/
}

func (b *Bittrex) GetSummary() common.CryptoExchangeSummary {
	BITTREX_RATE_LIMITER.RespectRateLimit()
	total := decimal.NewFromFloat(0)
	satoshis := decimal.NewFromFloat(0)
	balances, _ := b.GetBalances()
	for _, c := range balances {
		if c.GetCurrency() == "BTC" { // TODO
			total = total.Add(c.GetTotal())
		} else {
			satoshis = satoshis.Add(c.GetPrice().Mul(c.GetBalance()))
			total = total.Add(c.GetTotal())
		}
	}
	exchange := &dto.CryptoExchangeSummaryDTO{
		Name:     b.name,
		URL:      "https://www.bittrex.com",
		Total:    total.Truncate(8),
		Satoshis: satoshis.Truncate(8),
		Coins:    balances}
	return exchange
}

func (b *Bittrex) GetName() string {
	return b.name
}

func (b *Bittrex) GetDisplayName() string {
	return b.displayName
}

func (b *Bittrex) GetTradingFee() decimal.Decimal {
	return b.tradingFee
}

func (b *Bittrex) ParseImport(file string) ([]common.Transaction, error) {
	var orders []common.Transaction
	f, err := os.Open(file)
	if err != nil {
		return orders, err
	}
	defer f.Close()
	codec := unicode.UTF16(unicode.LittleEndian, unicode.UseBOM).NewDecoder()
	reader := codec.Reader(f)
	lines, err := csv.NewReader(reader).ReadAll()
	if err != nil {
		return orders, err
	}
	for i, values := range lines {
		if i == 0 {
			continue // skip header
		}
		var orderType string
		if values[2] == "LIMIT_BUY" {
			orderType = common.BUY_ORDER_TYPE
		} else {
			orderType = common.SELL_ORDER_TYPE
		}
		quantity, err := decimal.NewFromString(values[3])
		if err != nil {
			return nil, err
		}
		limit, err := decimal.NewFromString(values[4])
		if err != nil {
			return nil, err
		}
		fee, err := decimal.NewFromString(values[5])
		if err != nil {
			return nil, err
		}
		price, err := decimal.NewFromString(values[6])
		if err != nil {
			return nil, err
		}
		date, err := time.Parse("1/2/2006 15:04:05 PM", values[8])
		if err != nil {
			b.ctx.GetLogger().Errorf("[Bittrex.ParseImport] Error parsing float: %s", err.Error())
			return orders, err
		}
		marketPair, err := common.NewCurrencyPair(values[1], b.ctx.GetUser().GetLocalCurrency())
		if err != nil {
			b.ctx.GetLogger().Errorf("[Bittrex.ParseImport] Invalid currency pair: %s", values[1])
			return nil, err
		}
		currencyPair := &common.CurrencyPair{
			Base:          marketPair.Quote,
			Quote:         marketPair.Base,
			LocalCurrency: b.ctx.GetUser().GetLocalCurrency()}
		baseFiatPrice := b.getFiatPrice(marketPair.Base, date)
		limitFiatPrice := limit.Mul(baseFiatPrice)
		fiatFee := fee.Mul(baseFiatPrice)
		total := price.Add(fee)
		fiatTotal := total.Mul(baseFiatPrice)
		orders = append(orders, &dto.TransactionDTO{
			Id:                     fmt.Sprintf("bittrex-csv-%d", i),
			Network:                b.name,
			NetworkDisplayName:     b.displayName,
			Date:                   date,
			Type:                   orderType,
			Category:               common.TX_CATEGORY_TRADE,
			MarketPair:             marketPair,
			CurrencyPair:           currencyPair,
			Quantity:               quantity.StringFixed(8),
			QuantityCurrency:       marketPair.Quote,
			FiatQuantity:           quantity.Mul(limitFiatPrice).StringFixed(2),
			FiatQuantityCurrency:   "USD",
			Price:                  limit.StringFixed(8),
			PriceCurrency:          marketPair.Base,
			FiatPrice:              limitFiatPrice.StringFixed(2),
			FiatPriceCurrency:      "USD",
			QuoteFiatPrice:         baseFiatPrice.StringFixed(2),
			QuoteFiatPriceCurrency: marketPair.Base,
			Fee:               fee.StringFixed(8),
			FeeCurrency:       marketPair.Base,
			FiatFee:           fiatFee.StringFixed(2),
			FiatFeeCurrency:   "USD",
			Total:             total.StringFixed(8),
			TotalCurrency:     marketPair.Base,
			FiatTotal:         fiatTotal.StringFixed(2),
			FiatTotalCurrency: "USD"})
	}
	return orders, nil
}

func (b *Bittrex) FormattedCurrencyPair(marketPair *common.CurrencyPair) string {
	cp := b.localizedCurrencyPair(marketPair)
	return fmt.Sprintf("%s-%s", cp.Base, cp.Quote)
}

func (b *Bittrex) getFiatPrice(currency string, atDate time.Time) decimal.Decimal {
	cacheKey := fmt.Sprintf("%s-%s-%s", "bittrex-fiatprice", currency, atDate)
	if price, found := b.cache.Get(cacheKey); found {
		b.ctx.GetLogger().Debugf("[Bittrex.getFiatPrice] Returning fiat price from cache")
		return price.(decimal.Decimal)
	}
	b.ctx.GetLogger().Debugf("[Bittrex.getFiatPrice] Converting %s to fiat on %s", currency, atDate)
	fiatPrice := decimal.NewFromFloat(0)
	if currency == "BTC" {
		candle, err := b.GetPriceAt("BTC", atDate)
		if err != nil {
			b.ctx.GetLogger().Errorf("[Bittrex.getFiatPrice] Failed to get bitcoin price on %s: %s", atDate, err.Error())
		}
		fiatPrice = candle.Close
	} else if b.supportsUSDT(currency) {
		candle, err := b.GetPriceAt(currency, atDate)
		if err != nil {
			b.ctx.GetLogger().Errorf("[Bittrex.getFiatPrice] Failed to get USDT price for %s on %s: %s", currency, atDate, err.Error())
		}
		fiatPrice = candle.Close
	} else {
		candle, err := b.toBTCtoUSDT(currency, atDate)
		if err != nil {
			b.ctx.GetLogger().Errorf("[Bittrex.getFiatPrice] Failed to convert %s to BTC to USDT on %s: %s",
				currency, atDate, err.Error())
		}
		fiatPrice = candle.Close
	}
	b.cache.Set(cacheKey, fiatPrice, cache.NoExpiration)
	return fiatPrice
}

func (b *Bittrex) supportsUSDT(currency string) bool {
	for _, market := range b.usdtMarkets {
		if market == currency {
			return true
		}
	}
	return false
}

func (b *Bittrex) toBTCtoUSDT(currency string, atDate time.Time) (*common.Candlestick, error) {
	cacheKey := fmt.Sprintf("%s-%s-%s", "bittrex-btctousd", currency, atDate)
	if candle, found := b.cache.Get(cacheKey); found {
		b.ctx.GetLogger().Debugf("[Bittrex.toBTCtoUSDT] Returning USDT price from cache")
		return candle.(*common.Candlestick), nil
	}
	bitcoinCandle, err := b.GetPriceAt("BTC", atDate)
	marketPair := &common.CurrencyPair{
		Base:          "BTC",
		Quote:         currency,
		LocalCurrency: b.ctx.GetUser().GetLocalCurrency()}

	var candles []common.Candlestick
	monthAgo := time.Now().AddDate(0, -1, 0)
	if atDate.Before(monthAgo) {
		klines, err := b.GetPriceHistory(marketPair, atDate.Add(-24*time.Hour), atDate.Add(24*time.Hour), 1000)
		if err != nil {
			return &common.Candlestick{}, err
		}
		candles = klines
	} else {
		klines, err := b.GetPriceHistory(marketPair, atDate.Add(-5*time.Minute), atDate.Add(5*time.Minute), 1)
		if err != nil {
			return &common.Candlestick{}, err
		}
		candles = klines
	}
	currencyCandle, err := util.FindClosestDatedCandle(b.ctx.GetLogger(), atDate, candles)
	if err != nil {
		return &common.Candlestick{}, err
	}
	candlestick := &common.Candlestick{
		Date:  atDate,
		Close: currencyCandle.Close.Mul(bitcoinCandle.Close)}
	b.cache.Set(cacheKey, candlestick, cache.NoExpiration)
	return candlestick, nil
}

func (b *Bittrex) localizedCurrencyPair(marketPair *common.CurrencyPair) *common.CurrencyPair {
	var cp *common.CurrencyPair
	if marketPair.Quote == "USD" {
		cp = &common.CurrencyPair{
			Base:          "USDT",
			Quote:         marketPair.Base,
			LocalCurrency: "USDT"}
	} else {
		cp = marketPair
	}
	return cp
}

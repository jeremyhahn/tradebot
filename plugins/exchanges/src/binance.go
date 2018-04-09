package main

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	ws "github.com/gorilla/websocket"
	"github.com/jeremyhahn/go-binance"
	cache "github.com/patrickmn/go-cache"
	"github.com/shopspring/decimal"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dto"
	"github.com/jeremyhahn/tradebot/entity"
	"github.com/jeremyhahn/tradebot/util"
	"github.com/op/go-logging"
)

type AggregateTrade struct {
	EventType string `json:"a"`
	//EventTime     int64  `json:"E"` // Works as int64 in SymbolTicker  ¯\_(ツ)_/¯
	Symbol        string `json:"s"`
	TradeId       string `json:"a"`
	Price         string `json:"p"`
	Quantity      string `json:"q"`
	FirstTradeId  int64  `json:"f"`
	LastTradeId   int64  `json:"l"`
	TradeTime     int64  `json:"T"`
	IsBuyerTrader bool   `json:"m"`
}

type SymbolTicker struct {
	EventType            string `json:"e"`
	EventTime            int64  `json:"E"`
	Symbol               string `json:"s"`
	PriceChange          string `json:"p"`
	PriceChangePercent   string `json:"P"`
	WeightedAveragePrice string `json:"w"`
	PreviousDayClose     string `json:"x"`
	CurrentDayClose      string `json:"p"`
	CloseTrade           string `json:"Q"`
	BestBidPrice         string `json:"b"`
	BidQuantity          string `json:"B"`
	BestAskPrice         string `json:"a"`
	BestAskQuantity      string `json:"A"`
	OpenPrice            string `json:"o"`
	HighPrice            string `json:"h"`
	LowPrice             string `json:"l"`
	BaseVolume           string `json:"v"`
	QuoteVolume          string `json:"q"`
	StatsOpenTime        int64  `json:"O"`
	//StatsCloseTime       int64 `json:"C"`  //  what data type is this?
	FirstTradeId int64 `json:"F"`
	LastTradeId  int64 `json:"L"`
	TotalTrades  int64 `json:"n"`
}

type Binance struct {
	client        *binance.Client
	ctx           common.Context
	logger        *logging.Logger
	name          string
	displayName   string
	tradingFee    decimal.Decimal
	usdtMarkets   []string
	currencyPairs string
	cache         *cache.Cache
	common.Exchange
}

func CreateBinance(ctx common.Context, userExchangeEntity entity.UserExchangeEntity) common.Exchange {
	return &Binance{
		ctx:           ctx,
		client:        binance.NewClient(userExchangeEntity.GetKey(), userExchangeEntity.GetSecret()),
		logger:        ctx.GetLogger(),
		name:          "Binance",
		displayName:   "Binance",
		tradingFee:    decimal.NewFromFloat(.01),
		currencyPairs: userExchangeEntity.GetExtra(),
		usdtMarkets:   []string{"BTC", "ETH", "BNB", "NEO", "LTC", "BCC"},
		cache:         cache.New(1*time.Minute, 1*time.Minute)}
}

func (b *Binance) GetName() string {
	return b.name
}

func (b *Binance) GetDisplayName() string {
	return b.displayName
}

func (b *Binance) GetTradingFee() decimal.Decimal {
	return b.tradingFee
}

func (b *Binance) GetPriceAt(currency string, atDate time.Time) (*common.Candlestick, error) {
	currencyPair := &common.CurrencyPair{
		Base:          currency,
		Quote:         "USDT",
		LocalCurrency: b.ctx.GetUser().GetLocalCurrency()}
	kline, err := b.GetPriceHistory(currencyPair, atDate.Add(-5*time.Minute), atDate.Add(5*time.Minute), 1)
	if err != nil {
		return &common.Candlestick{}, err
	}
	closestCandle, err := util.FindClosestDatedCandle(b.ctx.GetLogger(), atDate, kline)
	if err != nil {
		return &common.Candlestick{}, err
	}
	return closestCandle, nil
}

func (b *Binance) GetBalances() ([]common.Coin, decimal.Decimal) {
	var coins []common.Coin
	sum := decimal.NewFromFloat(0)
	account, err := b.client.NewGetAccountService().Do(context.Background())
	if err != nil {
		b.logger.Errorf("[Binance.GetBalances] %s", err.Error())
		return coins, sum
	}
	for _, balance := range account.Balances {
		bal, err := decimal.NewFromString(balance.Free)
		if err != nil {
			b.logger.Errorf("[Binance.GetBalances] Error parsing balance from string: %s", err.Error())
		}
		if bal.GreaterThan(decimal.NewFromFloat(0)) {
			b.logger.Debugf("[Binance.GetBalances] Getting balance for %s", balance.Asset)
			bitcoin, err := b.GetPriceAt("BTC", time.Now())
			if err != nil {
				b.logger.Errorf("[Binance.GetBalances] Error getting current bitcoin price: %s", err.Error())
			}
			if balance.Asset == "BTC" {
				total := bal.Mul(bitcoin.Close).Truncate(8)
				sum = sum.Add(total)
				coins = append(coins, &dto.CoinDTO{
					Currency:  balance.Asset,
					Available: bal.Truncate(8),
					Balance:   bal.Truncate(8),
					Price:     bitcoin.Close.Truncate(2),
					Total:     total.Truncate(2)})
				continue
			}
			currencyPair := &common.CurrencyPair{
				Base:          balance.Asset,
				Quote:         b.ctx.GetUser().GetLocalCurrency(),
				LocalCurrency: b.ctx.GetUser().GetLocalCurrency()}
			localizedCurrencyPair := b.localizedCurrencyPair(currencyPair)
			symbol := fmt.Sprintf("%s%s", localizedCurrencyPair.Base, "BTC")
			ticker, err := b.client.NewPriceChangeStatsService().Symbol(symbol).Do(context.Background())
			if err != nil {
				b.logger.Errorf("[Binance.GetBalances] %s", err.Error())
			}
			if ticker == nil {
				b.logger.Errorf("[Binance.GetBalances] Unable to retrieve ticker for symbol: %s", symbol)
				continue
			}
			lastPrice, err := decimal.NewFromString(ticker.LastPrice)
			if err != nil {
				b.logger.Errorf("[Binance.GetBalances] Error parsing last ticker price into decimal: %s", err.Error())
			}
			subtotal := lastPrice.Mul(bitcoin.Close)
			total := bal.Mul(subtotal)
			sum = sum.Add(total)
			coins = append(coins, &dto.CoinDTO{
				Currency:  balance.Asset,
				Available: bal.Truncate(8),
				Balance:   bal.Truncate(8),
				Price:     lastPrice.Truncate(2),
				Total:     total.Truncate(2)})
		}
	}
	return coins, sum.Truncate(2)
}

func (b *Binance) GetPriceHistory(currencyPair *common.CurrencyPair,
	start, end time.Time, granularity int) ([]common.Candlestick, error) {

	candlesticks := make([]common.Candlestick, 0)
	interval := fmt.Sprintf("%dm", granularity)
	startTime := start.UnixNano() / int64(time.Millisecond)
	endTime := end.UnixNano() / int64(time.Millisecond)
	b.logger.Debugf("[Binance.GetPriceHistory] Getting %s price history between %s and %s on a %s interval",
		currencyPair, start, end, interval)
	klines, err := b.client.NewKlinesService().
		StartTime(startTime).EndTime(endTime).
		Symbol(b.FormattedCurrencyPair(currencyPair)).
		Interval(interval).Do(context.Background())
	if err != nil {
		b.logger.Errorf("[Binance.GetPriceHistory] Error: %s", err.Error())
		return candlesticks, err
	}
	for _, k := range klines {
		volume, err := decimal.NewFromString(k.Volume)
		if err != nil {
			b.logger.Errorf("[Binance.GetPriceHistory] Error parsing volume into string: %s", err.Error())
		}
		open, err := decimal.NewFromString(k.Open)
		if err != nil {
			b.logger.Errorf("[Binance.GetPriceHistory] Error parsing open price into string: %s", err.Error())
		}
		high, err := decimal.NewFromString(k.High)
		if err != nil {
			b.logger.Errorf("[Binance.GetPriceHistory] Error parsing high price into string: %s", err.Error())
		}
		low, err := decimal.NewFromString(k.Low)
		if err != nil {
			b.logger.Errorf("[Binance.GetPriceHistory] Error parsing low price into string: %s", err.Error())
		}
		close, err := decimal.NewFromString(k.Close)
		if err != nil {
			b.logger.Errorf("[Binance.GetPriceHistory] Error parsing close price into string: %s", err.Error())
		}
		candlesticks = append(candlesticks, common.Candlestick{
			CurrencyPair: currencyPair,
			Date:         time.Unix(k.CloseTime/1000, 0),
			Exchange:     b.name,
			Open:         open,
			High:         high,
			Low:          low,
			Close:        close,
			Period:       granularity,
			Volume:       volume})
	}
	return candlesticks, nil
}

func (b *Binance) GetOrderHistory(currencyPair *common.CurrencyPair) []common.Transaction {
	b.logger.Errorf("[Binance.GetOrderHistory] Getting %s order history", b.FormattedCurrencyPair(currencyPair))
	formattedCurrencyPair := b.FormattedCurrencyPair(currencyPair)
	orders, err := b.client.NewListTradesService().FromID(0).
		Symbol(formattedCurrencyPair).Do(context.Background())
	if err != nil {
		b.logger.Errorf("[Binance.GetOrderHistory] %s", err.Error())
	}
	var _orders []common.Transaction
	for _, o := range orders {
		var id, orderType string
		var fiatFee decimal.Decimal
		orderDate := time.Unix(o.Time/1000, 0)
		baseFiatPrice := b.getFiatPrice(currencyPair.Base, orderDate)
		quoteFiatPrice := b.getFiatPrice(currencyPair.Quote, orderDate)
		quantity, err := decimal.NewFromString(o.Quantity)
		if err != nil {
			b.ctx.GetLogger().Errorf("[Binance.GetOrderHistory] Failed to parse quantity string into decimal: %s", err.Error())
		}
		purchasePrice, err := decimal.NewFromString(o.Price)
		if err != nil {
			b.ctx.GetLogger().Errorf("[Binance.GetOrderHistory] Failed to parse purchase price string into decimal: %s", err.Error())
		}
		fee, err := decimal.NewFromString(o.Commission)
		if err != nil {
			b.ctx.GetLogger().Errorf("[Binance.GetOrderHistory] Failed to parse commission price string into decimal: %s", err.Error())
		}
		if o.ID <= 0 {
			id = fmt.Sprintf("binance-%s", orderDate)
		} else {
			id = fmt.Sprintf("%d", o.ID)
		}
		if o.IsBuyer {
			orderType = common.BUY_ORDER_TYPE
			fiatFee = fee.Mul(baseFiatPrice)
		} else {
			orderType = common.SELL_ORDER_TYPE
			fiatFee = fee.Mul(quoteFiatPrice)
		}
		fiatTotal := quantity.Mul(purchasePrice).Mul(quoteFiatPrice)
		fiatPrice := purchasePrice.Mul(quoteFiatPrice)
		_orders = append(_orders, &dto.TransactionDTO{
			Id:                   id,
			Network:              b.name,
			NetworkDisplayName:   b.displayName,
			Type:                 orderType,
			Category:             common.TX_CATEGORY_TRADE,
			CurrencyPair:         currencyPair,
			Date:                 orderDate,
			Quantity:             quantity.StringFixed(8),
			QuantityCurrency:     currencyPair.Base,
			FiatQuantity:         fiatTotal.Sub(fiatFee).StringFixed(2),
			FiatQuantityCurrency: "USD",
			Price:                purchasePrice.StringFixed(8),
			PriceCurrency:        currencyPair.Quote,
			FiatPrice:            fiatPrice.StringFixed(2),
			FiatPriceCurrency:    "USD",
			Fee:                  fee.StringFixed(8),
			FeeCurrency:          o.CommissionAsset,
			FiatFee:              fiatFee.StringFixed(2),
			FiatFeeCurrency:      "USD",
			TotalCurrency:        currencyPair.Quote,
			Total:                quantity.Mul(purchasePrice).StringFixed(8),
			FiatTotal:            fiatTotal.StringFixed(2),
			FiatTotalCurrency:    "USD"})
	}
	return _orders
}

func (b *Binance) GetDepositHistory() ([]common.Transaction, error) {
	var orders []common.Transaction
	currencies, _ := b.GetCurrencies()
	b.ctx.GetLogger().Debugf("[Binance.GetDepositHistory] Deposit currencies for %s: %s", b.ctx.GetUser().GetUsername(), currencies)
	for _, currency := range currencies {
		deposits, err := b.client.NewListDepositsService().Asset(currency.GetSymbol()).Do(context.Background())
		if err != nil {
			b.logger.Errorf("[Binance.GetDepositHistory] Error: %s", err.Error())
		}
		for _, deposit := range deposits {
			currencyPair := &common.CurrencyPair{
				Base:          deposit.Asset,
				Quote:         deposit.Asset,
				LocalCurrency: b.ctx.GetUser().GetLocalCurrency()}
			orderDate := time.Unix(deposit.InsertTime/1000, 0)
			quantity := decimal.NewFromFloat(deposit.Amount)
			baseFiatPrice := b.getFiatPrice(currencyPair.Base, orderDate)
			fiatTotal := quantity.Mul(baseFiatPrice)
			orders = append(orders, &dto.TransactionDTO{
				Id:                   deposit.TxID,
				Type:                 common.DEPOSIT_ORDER_TYPE,
				Category:             common.TX_CATEGORY_TRANSFER,
				Date:                 orderDate,
				Network:              b.name,
				NetworkDisplayName:   b.displayName,
				CurrencyPair:         currencyPair,
				Quantity:             quantity.StringFixed(8),
				QuantityCurrency:     deposit.Asset,
				FiatQuantity:         fiatTotal.StringFixed(2),
				FiatQuantityCurrency: "USD",
				Price:                baseFiatPrice.StringFixed(2),
				PriceCurrency:        "USD",
				FiatPrice:            baseFiatPrice.StringFixed(2),
				FiatPriceCurrency:    "USD",
				Fee:                  "0.00000000",
				FeeCurrency:          deposit.Asset,
				FiatFee:              "0.00",
				FiatFeeCurrency:      "USD",
				Total:                quantity.StringFixed(8),
				TotalCurrency:        deposit.Asset,
				FiatTotal:            fiatTotal.StringFixed(2),
				FiatTotalCurrency:    "USD"})
		}
	}
	return orders, nil
}

func (b *Binance) GetWithdrawalHistory() ([]common.Transaction, error) {
	var orders []common.Transaction
	currencies, _ := b.GetCurrencies()
	b.ctx.GetLogger().Debugf("[Binance.GetWithdrawalHistory] Withdrawal currencies for %s: %s", b.ctx.GetUser().GetUsername(), currencies)
	for _, currency := range currencies {
		withdraws, err := b.client.NewListWithdrawsService().Asset(currency.GetSymbol()).Do(context.Background())
		if err != nil {
			b.logger.Errorf("[Binance.GetWithdrawalHistory] Error: %s", err.Error())
		}
		for _, withdrawal := range withdraws {
			if withdrawal.Status != 6 { // 0:Email Sent,1:Cancelled 2:Awaiting Approval 3:Rejected 4:Processing 5:Failure 6:Completed
				continue
			}
			orderDate := time.Unix(withdrawal.ApplyTime/1000, 0)
			currencyPair := &common.CurrencyPair{
				Base:          withdrawal.Asset,
				Quote:         withdrawal.Asset,
				LocalCurrency: b.ctx.GetUser().GetLocalCurrency()}
			baseFiatPrice := b.getFiatPrice(currencyPair.Base, orderDate)
			quantity := decimal.NewFromFloat(withdrawal.Amount)
			fiatTotal := quantity.Mul(baseFiatPrice)
			orders = append(orders, &dto.TransactionDTO{
				Id:                   withdrawal.TxID,
				Type:                 common.WITHDRAWAL_ORDER_TYPE,
				Category:             common.TX_CATEGORY_TRANSFER,
				Date:                 orderDate,
				Network:              b.name,
				NetworkDisplayName:   b.displayName,
				CurrencyPair:         currencyPair,
				Quantity:             quantity.StringFixed(8),
				QuantityCurrency:     withdrawal.Asset,
				FiatQuantity:         fiatTotal.StringFixed(2),
				FiatQuantityCurrency: "USD",
				Price:                baseFiatPrice.StringFixed(2),
				PriceCurrency:        "USD",
				FiatPrice:            baseFiatPrice.StringFixed(2),
				FiatPriceCurrency:    "USD",
				Fee:                  "0.00000000",
				FeeCurrency:          withdrawal.Asset,
				FiatFee:              "0.00",
				FiatFeeCurrency:      "N/A",
				Total:                quantity.StringFixed(8),
				TotalCurrency:        withdrawal.Asset,
				FiatTotal:            fiatTotal.StringFixed(2),
				FiatTotalCurrency:    "USD"})
		}
	}
	return orders, nil
}

func (b *Binance) SubscribeToLiveFeed(currencyPair *common.CurrencyPair, priceChange chan common.PriceChange) {
	var wsDialer ws.Dialer
	formattedCurrencyPair := b.FormattedCurrencyPair(currencyPair)
	url := fmt.Sprintf("wss://stream.binance.com:9443/ws/%s@aggTrade", strings.ToLower(formattedCurrencyPair))

	b.logger.Info("[Binance.SubscribeToLiveFeed] Subscribing to websocket feed: %s", url)

	wsConn, _, err := wsDialer.Dial(url, nil)
	if err != nil {
		b.logger.Errorf("[Binance.SubscribeToLiveFeed] %s", err.Error())
	}

	subscribe := map[string]string{
		"type":       "subscribe",
		"product_id": formattedCurrencyPair,
	}

	if err := wsConn.WriteJSON(subscribe); err != nil {
		b.logger.Errorf("[Binance.SubscribeToLiveFeed] %s", err.Error())
	}

	var message AggregateTrade
	for true {

		if err := wsConn.ReadJSON(&message); err != nil {
			b.logger.Errorf("[Binance.SubscribeToLiveFeed] %s", err.Error())
			continue
		}

		b.logger.Debugf("[Binance.SubscribeToLiveFeed] %+v\n", message)

		price, err := decimal.NewFromString(message.Price)
		if err != nil {
			b.logger.Errorf("[Binance.SubscribeToLiveFeed] Error parsing price into string: %s", err.Error())
		}

		priceChange <- common.PriceChange{
			CurrencyPair: currencyPair,
			Exchange:     b.name,
			Price:        price,
			Satoshis:     decimal.NewFromFloat(1.0)}
	}

	b.SubscribeToLiveFeed(currencyPair, priceChange)
}

func (b *Binance) GetCurrencies() (map[string]*common.Currency, error) {
	b.ctx.GetLogger().Errorf("[Binance.GetCurrencies] Configured currencies: %s", b.currencyPairs)
	currencies := make(map[string]*common.Currency)
	configuredCurrencies := strings.Split(b.currencyPairs, ",")
	currencyMap := make(map[string]bool, len(configuredCurrencies))
	for _, cur := range configuredCurrencies {
		pieces := strings.Split(cur, "-")
		if len(pieces) < 2 {
			errmsg := fmt.Sprintf("Invalid currency pair: %s", cur)
			b.ctx.GetLogger().Errorf("[Binance.GetCurrencies] Error: %s", errmsg)
			return nil, errors.New(errmsg)
		} else {
			currencyMap[pieces[0]] = true
			currencyMap[pieces[1]] = true
		}
	}
	for k, _ := range currencyMap {
		var name string
		if _, found := common.CryptoNames[k]; !found {
			b.ctx.GetLogger().Errorf("[Binance.GetCurrencies] Unable to locate currency %s in common.CryptoNames", k)
			name = k
		} else {
			name = common.CryptoNames[k]
		}
		currencies[k] = &common.Currency{
			ID:           k,
			Name:         name,
			Symbol:       k,
			BaseUnit:     100000000,
			DecimalPlace: 2}
	}
	return currencies, nil
}

func (b *Binance) GetSummary() common.CryptoExchangeSummary {
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
		URL:      "https://www.binance.com",
		Total:    total.Truncate(8),
		Satoshis: satoshis.Truncate(8),
		Coins:    balances}
	return exchange
}

func (b *Binance) FormattedCurrencyPair(currencyPair *common.CurrencyPair) string {
	cp := b.localizedCurrencyPair(currencyPair)
	return fmt.Sprintf("%s%s", cp.Base, cp.Quote)
}

func (b *Binance) ParseImport(file string) ([]common.Transaction, error) {
	var orders []common.Transaction
	b.ctx.GetLogger().Error("[Binance.ParseImport] Unsupported!")
	return orders, errors.New("Binance.ParseImport Unsupported")
}

func (b *Binance) getFiatPrice(currency string, atDate time.Time) decimal.Decimal {
	cacheKey := fmt.Sprintf("%s-%s-%s", "binance-fiatprice", currency, atDate)
	if price, found := b.cache.Get(cacheKey); found {
		b.ctx.GetLogger().Debugf("[Binance.getFiatPrice] Returning fiat price from cache")
		return price.(decimal.Decimal)
	}
	b.ctx.GetLogger().Debugf("[Binance.getFiatPrice] Converting %s to fiat on %s", currency, atDate)
	fiatPrice := decimal.NewFromFloat(0)
	if currency == "BTC" {
		candle, err := b.GetPriceAt("BTC", atDate)
		if err != nil {
			b.ctx.GetLogger().Errorf("[Binance.getFiatPrice] Failed to get bitcoin price on %s: %s", atDate, err.Error())
		}
		fiatPrice = candle.Close
	} else if b.supportsUSDT(currency) {
		candle, err := b.GetPriceAt(currency, atDate)
		if err != nil {
			b.ctx.GetLogger().Errorf("[Binance.getFiatPrice] Failed to get USDT price for %s on %s: %s", currency, atDate, err.Error())
		}
		fiatPrice = candle.Close
	} else {
		candle, err := b.toBTCtoUSDT(currency, atDate)
		if err != nil {
			b.ctx.GetLogger().Errorf("[Binance.getFiatPrice] Failed to convert %s to BTC to USDT on %s: %s",
				currency, atDate, err.Error())
		}
		fiatPrice = candle.Close
	}
	b.cache.Set(cacheKey, fiatPrice, cache.NoExpiration)
	return fiatPrice
}

func (b *Binance) supportsUSDT(currency string) bool {
	for _, market := range b.usdtMarkets {
		if market == currency {
			return true
		}
	}
	return false
}

func (b *Binance) toBTCtoUSDT(currency string, atDate time.Time) (*common.Candlestick, error) {
	cacheKey := fmt.Sprintf("%s-%s-%s", "binance-btctousd", currency, atDate)
	if candle, found := b.cache.Get(cacheKey); found {
		b.ctx.GetLogger().Debugf("[Binance.toBTCtoUSDT] Returning USDT price from cache")
		return candle.(*common.Candlestick), nil
	}
	bitcoinCandle, err := b.GetPriceAt("BTC", atDate)
	currencyPair := &common.CurrencyPair{
		Base:          currency,
		Quote:         "BTC",
		LocalCurrency: b.ctx.GetUser().GetLocalCurrency()}
	klines, err := b.GetPriceHistory(currencyPair, atDate.Add(-5*time.Minute), atDate.Add(5*time.Minute), 1)
	if err != nil {
		return &common.Candlestick{}, err
	}
	currencyCandle, err := util.FindClosestDatedCandle(b.ctx.GetLogger(), atDate, klines)
	if err != nil {
		return &common.Candlestick{}, err
	}
	candlestick := &common.Candlestick{
		Date:  atDate,
		Close: currencyCandle.Close.Mul(bitcoinCandle.Close)}
	b.cache.Set(cacheKey, candlestick, cache.NoExpiration)
	return candlestick, nil
}

func (b *Binance) localizedCurrencyPair(currencyPair *common.CurrencyPair) *common.CurrencyPair {
	if currencyPair.Quote == "USD" {
		return &common.CurrencyPair{
			Base:  currencyPair.Base,
			Quote: "USDT"}
	}
	return currencyPair
}

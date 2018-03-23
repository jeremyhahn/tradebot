package main

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/adshao/go-binance"
	ws "github.com/gorilla/websocket"
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
	tradingFee    float64
	usdtMarkets   []string
	currencyPairs string
	common.Exchange
}

func CreateBinance(ctx common.Context, userExchangeEntity entity.UserExchangeEntity) common.Exchange {
	return &Binance{
		ctx:           ctx,
		client:        binance.NewClient(userExchangeEntity.GetKey(), userExchangeEntity.GetSecret()),
		logger:        ctx.GetLogger(),
		name:          "Binance",
		displayName:   "Binance",
		tradingFee:    .01,
		currencyPairs: userExchangeEntity.GetExtra(),
		usdtMarkets:   []string{"BTC", "ETH", "BNB", "NEO", "LTC", "BCC"}}
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
	closestCandle, err := util.FindClosestDatedCandle(b.ctx.GetLogger(), atDate, &kline)
	if err != nil {
		return &common.Candlestick{}, err
	}
	return closestCandle, nil
}

func (b *Binance) GetBalances() ([]common.Coin, float64) {

	var coins []common.Coin
	sum := 0.0

	account, err := b.client.NewGetAccountService().Do(context.Background())
	if err != nil {
		b.logger.Errorf("[Binance.GetBalances] %s", err.Error())
		return coins, sum
	}

	for _, balance := range account.Balances {

		bal, err := strconv.ParseFloat(balance.Free, 64)
		if err != nil {
			b.logger.Errorf("[Binance.GetBalances] %s", err.Error())
		}

		if bal > 0 {

			b.logger.Debugf("[Binance.GetBalances] Getting balance for %s", balance.Asset)

			bitcoin := b.getBitcoin()
			bitcoinPrice := b.parseBitcoinPrice(bitcoin)

			if balance.Asset == "BTC" {
				total := bal * bitcoinPrice
				t, err := strconv.ParseFloat(strconv.FormatFloat(total, 'f', 2, 64), 64)
				if err != nil {
					b.logger.Errorf("[Binance.GetBalances] %s", err.Error())
				}
				sum += t
				coins = append(coins, &dto.CoinDTO{
					Currency:  balance.Asset,
					Available: bal,
					Balance:   bal,
					Price:     bitcoinPrice,
					Total:     t})
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

			lastPrice, err := strconv.ParseFloat(ticker.LastPrice, 64)
			if err != nil {
				b.logger.Errorf("[Binance.GetBalances] %s", err.Error())
			}

			total := bal * (lastPrice * bitcoinPrice)

			t, err := strconv.ParseFloat(fmt.Sprintf("%.2f", total), 64)
			if err != nil {
				b.logger.Errorf("[Binance.GetBalances] %s", err.Error())
			}

			sum += t
			coins = append(coins, &dto.CoinDTO{
				Currency:  balance.Asset,
				Available: bal,
				Balance:   bal,
				Price:     lastPrice,
				Total:     t})
		}
	}

	return coins, sum
}

func (b *Binance) GetPriceHistory(currencyPair *common.CurrencyPair,
	start, end time.Time, granularity int) ([]common.Candlestick, error) {
	candlesticks := make([]common.Candlestick, 0)
	interval := fmt.Sprintf("%dm", granularity)
	startTime := start.UnixNano() / int64(time.Millisecond)
	endTime := end.UnixNano() / int64(time.Millisecond)
	b.logger.Debugf("[Binance.GetPriceHistory] Getting %s price history between %d and %d on a %s interval",
		currencyPair, startTime, endTime, interval)
	klines, err := b.client.NewKlinesService().
		StartTime(startTime).EndTime(endTime).
		Symbol(b.FormattedCurrencyPair(currencyPair)).
		Interval(interval).Do(context.Background())
	if err != nil {
		b.logger.Errorf("[Binance.GetPriceHistory] Error: %s", err.Error())
		return candlesticks, err
	}
	for _, k := range klines {
		volume, _ := strconv.ParseFloat(k.Volume, 64)
		open, _ := strconv.ParseFloat(k.Open, 64)
		high, _ := strconv.ParseFloat(k.High, 64)
		low, _ := strconv.ParseFloat(k.Low, 64)
		close, _ := strconv.ParseFloat(k.Close, 64)
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
		if o.IsBuyer {
			orderType = "buy"
		} else {
			orderType = "sell"
		}
		orderDate := time.Unix(o.Time/1000, 0)
		qty, err := strconv.ParseFloat(o.Quantity, 64)
		if err != nil {
			b.ctx.GetLogger().Errorf("[Binance.GetOrderHistory] Failed to parse quantity to float64: %s", err.Error())
		}
		p, err := strconv.ParseFloat(o.Price, 64)
		if err != nil {
			b.ctx.GetLogger().Errorf("[Binance.GetOrderHistory] Failed to parse price to float64: %s", err.Error())
		}
		c, err := strconv.ParseFloat(o.Commission, 64)
		if err != nil {
			b.ctx.GetLogger().Errorf("[Binance.GetOrderHistory] Failed to parse commission to float64: %s", err.Error())
		}
		if o.ID <= 0 {
			id = fmt.Sprintf("%s", orderDate)
		} else {
			id = fmt.Sprintf("%d", o.ID)
		}
		_orders = append(_orders, &dto.TransactionDTO{
			Id:                   id,
			Network:              b.name,
			NetworkDisplayName:   b.displayName,
			Type:                 orderType,
			CurrencyPair:         currencyPair,
			Date:                 orderDate,
			Quantity:             decimal.NewFromFloat(qty).StringFixed(8),
			QuantityCurrency:     currencyPair.Base,
			FiatQuantity:         "0.00",
			FiatQuantityCurrency: "N/A",
			Price:                decimal.NewFromFloat(p).StringFixed(8),
			PriceCurrency:        "BTC",
			FiatPrice:            "0.00",
			FiatPriceCurrency:    "N/A",
			Fee:                  decimal.NewFromFloat(c).StringFixed(8),
			FeeCurrency:          currencyPair.Base,
			FiatFee:              "0.00",
			FiatFeeCurrency:      "N/A",
			TotalCurrency:        currencyPair.Quote,
			Total:                decimal.NewFromFloat(qty).Mul(decimal.NewFromFloat(p)).StringFixed(8),
			FiatTotal:            "0.00",
			FiatTotalCurrency:    "N/A"})

	}
	return _orders
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
			var id string
			if deposit.InsertTime <= 0 {
				id = fmt.Sprintf("%s", orderDate)
			} else {
				id = fmt.Sprintf("%d", deposit.InsertTime)
			}
			orders = append(orders, &dto.TransactionDTO{
				Id:                   id,
				Type:                 common.DEPOSIT_ORDER_TYPE,
				Date:                 orderDate,
				Network:              b.name,
				NetworkDisplayName:   b.displayName,
				CurrencyPair:         currencyPair,
				Quantity:             decimal.NewFromFloat(deposit.Amount).StringFixed(8),
				QuantityCurrency:     deposit.Asset,
				FiatQuantity:         "0.00",
				FiatQuantityCurrency: "N/A",
				Price:                "0.00000000",
				PriceCurrency:        deposit.Asset,
				FiatPrice:            "0.00",
				FiatPriceCurrency:    "N/A",
				Fee:                  "0.00000000",
				FeeCurrency:          deposit.Asset,
				FiatFee:              "0.00",
				FiatFeeCurrency:      "N/A",
				Total:                decimal.NewFromFloat(deposit.Amount).StringFixed(8),
				TotalCurrency:        deposit.Asset,
				FiatTotal:            "0.00",
				FiatTotalCurrency:    "N/A"})
		}
	}
	return orders, nil
}

func (b *Binance) supportsUSDT(currency string) bool {
	for _, market := range b.usdtMarkets {
		if market == currency {
			return true
		}
	}
	return false
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
		for _, withdraw := range withdraws {
			if withdraw.Status != 6 { // 0:Email Sent,1:Cancelled 2:Awaiting Approval 3:Rejected 4:Processing 5:Failure 6:Completed
				continue
			}
			orderDate := time.Unix(withdraw.ApplyTime/1000, 0)
			currencyPair := &common.CurrencyPair{
				Base:          withdraw.Asset,
				Quote:         withdraw.Asset,
				LocalCurrency: b.ctx.GetUser().GetLocalCurrency()}
			/*
				currencyPair := &common.CurrencyPair{
					Base:          withdraw.Asset,
					Quote:         "USDT",
					LocalCurrency: b.ctx.GetUser().GetLocalCurrency()}
						var priceCandle *common.Candlestick
							if b.supportsUSDT(withdraw.Asset) {
								usdtCandle, err := b.GetPriceAt(withdraw.Asset, orderDate)
								if err != nil {
									return nil, err
								}
								priceCandle = usdtCandle
							} else {
								btcToUsdCandle, err := b.toBTCtoUSD(withdraw.Asset, orderDate)
								if err != nil {
									return nil, err
								}
								priceCandle = btcToUsdCandle
							}
					orders = append(orders, &dto.OrderDTO{
						Id:               string(withdraw.ApplyTime),
						Type:             common.WITHDRAWAL_ORDER_TYPE,
						Date:             orderDate,
						Exchange:         b.name,
						CurrencyPair:     currencyPair,
						Quantity:         withdraw.Amount,
						QuantityCurrency: withdraw.Asset,
						Price:            withdraw.Amount * priceCandle.Close,
						PriceCurrency:    currencyPair.Quote,
						Fee:              0.0,
						FeeCurrency:      currencyPair.Quote,
						Total:            withdraw.Amount * priceCandle.Close,
						TotalCurrency:    currencyPair.Quote,
					  CoinPrice:         priceCandle.Close,
					  CoinPriceCurrency: currencyPair.Quote})*/

			var id string
			if withdraw.ApplyTime <= 0 {
				id = fmt.Sprintf("%s", orderDate)
			} else {
				id = fmt.Sprintf("%d", withdraw.ApplyTime)
			}
			orders = append(orders, &dto.TransactionDTO{
				Id:                   id,
				Type:                 common.WITHDRAWAL_ORDER_TYPE,
				Date:                 orderDate,
				Network:              b.name,
				NetworkDisplayName:   b.displayName,
				CurrencyPair:         currencyPair,
				Quantity:             decimal.NewFromFloat(withdraw.Amount).StringFixed(8),
				QuantityCurrency:     withdraw.Asset,
				FiatQuantity:         "0.00",
				FiatQuantityCurrency: "N/A",
				Price:                "0.00000000",
				PriceCurrency:        withdraw.Asset,
				FiatPrice:            "0.00",
				FiatPriceCurrency:    "N/A",
				Fee:                  "0.00000000",
				FeeCurrency:          withdraw.Asset,
				FiatFee:              "0.00",
				FiatFeeCurrency:      "N/A",
				Total:                decimal.NewFromFloat(withdraw.Amount).StringFixed(8),
				TotalCurrency:        withdraw.Asset,
				FiatTotal:            "0.00",
				FiatTotalCurrency:    "N/A"})
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

		f, err := strconv.ParseFloat(message.Price, 64)
		if err != nil {
			b.logger.Errorf("[Bittrex.GetBalances] %s", err.Error())
		}

		priceChange <- common.PriceChange{
			CurrencyPair: currencyPair,
			Exchange:     b.name,
			Price:        f,
			Satoshis:     1.0}
	}

	b.SubscribeToLiveFeed(currencyPair, priceChange)
}

func (b *Binance) ParseImport(file string) ([]common.Transaction, error) {
	var orders []common.Transaction
	b.ctx.GetLogger().Error("[Binance.ParseImport] Unsupported!")
	return orders, errors.New("Binance.ParseImport Unsupported")
}

func (b *Binance) GetSummary() common.CryptoExchangeSummary {
	total := 0.0
	satoshis := 0.0
	balances, _ := b.GetBalances()
	for _, c := range balances {
		if c.GetCurrency() == "BTC" { // TODO
			total += c.GetTotal()
		} else {
			satoshis += c.GetPrice() * c.GetBalance()
			total += c.GetTotal()
		}
	}
	f, _ := strconv.ParseFloat(fmt.Sprintf("%.8f", satoshis), 64)
	t, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", total), 64)
	exchange := &dto.CryptoExchangeSummaryDTO{
		Name:     b.name,
		URL:      "https://www.binance.com",
		Total:    t,
		Satoshis: f,
		Coins:    balances}
	return exchange
}

func (b *Binance) GetName() string {
	return b.name
}

func (b *Binance) GetDisplayName() string {
	return b.displayName
}

func (b *Binance) GetTradingFee() float64 {
	return b.tradingFee
}

func (b *Binance) FormattedCurrencyPair(currencyPair *common.CurrencyPair) string {
	cp := b.localizedCurrencyPair(currencyPair)
	return fmt.Sprintf("%s%s", cp.Base, cp.Quote)
}

func (b *Binance) toBTCtoUSD(currency string, atDate time.Time) (*common.Candlestick, error) {
	bitcoinCandle, err := b.GetPriceAt("BTC", atDate)
	currencyPair := &common.CurrencyPair{
		Base:          currency,
		Quote:         "BTC",
		LocalCurrency: b.ctx.GetUser().GetLocalCurrency()}
	klines, err := b.GetPriceHistory(currencyPair, atDate.Add(-5*time.Minute), atDate.Add(5*time.Minute), 1)
	if err != nil {
		return &common.Candlestick{}, err
	}
	currencyCandle, err := util.FindClosestDatedCandle(b.ctx.GetLogger(), atDate, &klines)
	if err != nil {
		return &common.Candlestick{}, err
	}
	return &common.Candlestick{
		Date:  atDate,
		Close: currencyCandle.Close * bitcoinCandle.Close}, nil
}

func (b *Binance) getBitcoin() *binance.PriceChangeStats {
	currencyPair := &common.CurrencyPair{
		Base:          "BTC",
		Quote:         b.ctx.GetUser().GetLocalCurrency(),
		LocalCurrency: b.ctx.GetUser().GetLocalCurrency()}
	localizedCurrencyPair := b.localizedCurrencyPair(currencyPair)
	symbol := fmt.Sprintf("%s%s", localizedCurrencyPair.Base, localizedCurrencyPair.Quote)
	stats, err := b.client.NewPriceChangeStatsService().Symbol(symbol).Do(context.Background())
	if err != nil {
		b.logger.Errorf("[Binance.getBitcoin] %s", err.Error())
	}
	return stats
}

func (b *Binance) getBitcoinPrice() float64 {
	bitcoin := b.getBitcoin()
	f, err := strconv.ParseFloat(bitcoin.LastPrice, 8)
	if err != nil {
		b.logger.Errorf("[Binance.getBitcoinPrice] %s", err.Error())
	}
	return f
}

func (b *Binance) parseBitcoinPrice(bitcoin *binance.PriceChangeStats) float64 {
	if bitcoin == nil {
		b.logger.Error("[Binance.parseBitcoinPrice] Null pointer returned from Binance library")
		return 0.0
	}
	f, err := strconv.ParseFloat(bitcoin.LastPrice, 8)
	if err != nil {
		b.logger.Errorf("[Binance.parseBitcoinPrice] %s", err.Error())
	}
	return f
}

func (b *Binance) localizedCurrencyPair(currencyPair *common.CurrencyPair) *common.CurrencyPair {
	if currencyPair.Quote == "USD" {
		return &common.CurrencyPair{
			Base:  currencyPair.Base,
			Quote: "USDT"}
	}
	return currencyPair
}

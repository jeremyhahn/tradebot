package exchange

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/adshao/go-binance"
	ws "github.com/gorilla/websocket"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dto"
	"github.com/jeremyhahn/tradebot/entity"
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
	client              *binance.Client
	ctx                 common.Context
	logger              *logging.Logger
	name                string
	tradingFee          float64
	priceHistoryService common.PriceHistoryService
	common.Exchange
}

func NewBinance(ctx common.Context, exchange entity.UserExchangeEntity, priceHistoryService common.PriceHistoryService) common.Exchange {
	return &Binance{
		ctx:                 ctx,
		client:              binance.NewClient(exchange.GetKey(), exchange.GetSecret()),
		logger:              ctx.GetLogger(),
		name:                "binance",
		tradingFee:          .01,
		priceHistoryService: priceHistoryService}
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
	start, end time.Time, granularity int) []common.Candlestick {
	candlesticks := make([]common.Candlestick, 0)
	interval := fmt.Sprintf("%dm", granularity)
	startTime := start.UnixNano() / int64(time.Millisecond)
	endTime := end.UnixNano() / int64(time.Millisecond)
	b.logger.Debugf("[Binance.GetPriceHistory] Getting %s price history between %d and %d on a %s interval", currencyPair, startTime, endTime, interval)
	klines, err := b.client.NewKlinesService().
		StartTime(startTime).EndTime(endTime).
		Symbol(b.FormattedCurrencyPair(currencyPair)).
		Interval(interval).Do(context.Background())
	if err != nil {
		b.logger.Errorf("[Binance.GetPriceHistory] Error: %s", err.Error())
		return candlesticks
	}
	for _, k := range klines {
		volume, _ := strconv.ParseFloat(k.Volume, 64)
		close, _ := strconv.ParseFloat(k.Close, 64)
		open, _ := strconv.ParseFloat(k.Open, 64)
		candlesticks = append(candlesticks, common.Candlestick{
			Date:   time.Unix(k.CloseTime/1000, 0),
			Close:  close,
			Open:   open,
			Period: granularity,
			Volume: volume})
	}
	return candlesticks
}

func (b *Binance) GetOrderHistory(currencyPair *common.CurrencyPair) []common.Order {
	b.logger.Errorf("[Binance.GetOrderHistory] Getting %s order history", b.FormattedCurrencyPair(currencyPair))
	formattedCurrencyPair := b.FormattedCurrencyPair(currencyPair)
	orders, err := b.client.NewListTradesService().FromID(0).
		Symbol(formattedCurrencyPair).Do(context.Background())
	if err != nil {
		b.logger.Errorf("[Binance.GetOrderHistory] %s", err.Error())
	}
	var _orders []common.Order
	for _, o := range orders {
		var orderType string
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
		/*startDate := orderDate.Add(-10 * time.Minute)
		endDate := orderDate.Add(5 * time.Minute)
		b.logger.Debugf("[Binance.GetOrderHistory] Looking for order date %s within range %s - %s", orderDate, startDate, endDate)
		priceHistory := b.GetPriceHistory(currencyPair, startDate, endDate, 1)
		b.logger.Debugf("[Binance.GetOrderHistory] Found %d records", len(priceHistory))
		closestCandle, err := util.FindClosesttDatedCandle(b.ctx.GetLogger(), orderDate, priceHistory)
		if err != nil {
			b.logger.Errorf("[Binance.GetOrderHistory] Error %s", err.Error())
			return _orders
		}
		b.logger.Debugf("[Binance.GetOrderHistory] Using closest order price of %f on %s", closestCandle.Close, closestCandle.Date)*/
		_orders = append(_orders, &dto.OrderDTO{
			Id:                 strconv.FormatInt(int64(o.ID), 10),
			Exchange:           "binance",
			Type:               orderType,
			CurrencyPair:       currencyPair,
			Date:               orderDate,
			Fee:                c,
			Quantity:           qty,
			QuantityCurrency:   currencyPair.Base,
			Price:              p,
			Total:              qty * p,
			PriceCurrency:      "BTC",
			FeeCurrency:        currencyPair.Base,
			TotalCurrency:      currencyPair.Quote,
			HistoricalPrice:    b.priceHistoryService.GetClosePriceOn(currencyPair.Base, orderDate),
			HistoricalCurrency: b.ctx.GetUser().GetLocalCurrency()})

		//return _orders
	}
	return _orders
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

func (b *Binance) ParseImport(file string) ([]common.Order, error) {
	var orders []common.Order
	b.ctx.GetLogger().Error("[Binance.ParseImport] Unsupported!")
	return orders, errors.New("Binance.ParseImport Unsupported")
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

func (b *Binance) GetTradingFee() float64 {
	return b.tradingFee
}

func (b *Binance) FormattedCurrencyPair(currencyPair *common.CurrencyPair) string {
	cp := b.localizedCurrencyPair(currencyPair)
	return fmt.Sprintf("%s%s", cp.Base, cp.Quote)
}

func (b *Binance) localizedCurrencyPair(currencyPair *common.CurrencyPair) *common.CurrencyPair {
	if currencyPair.Quote == "USD" {
		return &common.CurrencyPair{
			Base:  currencyPair.Base,
			Quote: "USDT"}
	}
	return currencyPair
}

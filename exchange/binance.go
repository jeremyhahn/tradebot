package exchange

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/adshao/go-binance"
	ws "github.com/gorilla/websocket"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/dto"
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
	client     *binance.Client
	ctx        *common.Context
	logger     *logging.Logger
	name       string
	tradingFee float64
	common.Exchange
}

func NewBinance(ctx *common.Context, exchange *dao.UserCryptoExchange) common.Exchange {
	return &Binance{
		ctx:        ctx,
		client:     binance.NewClient(exchange.Key, exchange.Secret),
		logger:     ctx.Logger,
		name:       "binance",
		tradingFee: .01}
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
				coins = append(coins, common.Coin{
					Currency:  balance.Asset,
					Available: bal,
					Balance:   bal,
					Price:     bitcoinPrice,
					Total:     t})
				continue
			}

			symbol := fmt.Sprintf("%s%s", balance.Asset, b.ctx.User.LocalCurrency)
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
			coins = append(coins, common.Coin{
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

	b.logger.Debug("[Binance.GetTradeHistory] Getting trade history")
	candlesticks := make([]common.Candlestick, 0)
	interval := fmt.Sprintf("%dm", granularity/60)
	klines, err := b.client.NewKlinesService().
		Symbol(b.FormattedCurrencyPair(currencyPair)).
		Interval(interval).Do(context.Background())
	if err != nil {
		b.logger.Errorf("[Binance.GetTradeHistory] %s", err.Error())
		return candlesticks
	}
	for _, k := range klines {
		volume, _ := strconv.ParseFloat(k.Volume, 64)
		close, _ := strconv.ParseFloat(k.Close, 64)
		open, _ := strconv.ParseFloat(k.Open, 64)
		candlesticks = append(candlesticks, common.Candlestick{
			Close:  close,
			Open:   open,
			Period: granularity,
			Volume: volume})
	}
	return candlesticks
}

func (b *Binance) GetOrderHistory(currencyPair *common.CurrencyPair) []common.Order {
	var _orders []common.Order
	formattedCurrencyPair := b.FormattedCurrencyPair(currencyPair)
	orders, err := b.client.NewListTradesService().
		Symbol(formattedCurrencyPair).Do(context.Background())
	if err != nil {
		b.logger.Errorf("[Binance.GetOrderHistory] %s", err.Error())
	}
	for _, o := range orders {
		var orderType string
		if o.IsBuyer {
			orderType = "buy"
		} else {
			orderType = "sell"
		}
		qty, _ := strconv.ParseFloat(o.Quantity, 64)
		p, _ := strconv.ParseFloat(o.Price, 64)
		_orders = append(_orders, &dto.OrderDTO{
			Id:       strconv.FormatInt(int64(o.ID), 10),
			Exchange: "binance",
			Type:     orderType,
			Currency: formattedCurrencyPair,
			//Date:     time.Unix(o.Time, 0),
			Quantity: qty,
			Price:    p})
		fmt.Println(o)
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

func (b *Binance) getBitcoin() *binance.PriceChangeStats {
	stats, err := b.client.NewPriceChangeStatsService().Symbol("BTCUSDT").Do(context.Background())
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
	f, err := strconv.ParseFloat(bitcoin.LastPrice, 8)
	if err != nil {
		b.logger.Errorf("[Binance.parseBitcoinPrice] %s", err.Error())
	}
	return f
}

func (b *Binance) GetExchangeAsync(exchangeChan *chan common.CryptoExchange) {
	go func() { *exchangeChan <- b.GetExchange() }()
}

func (b *Binance) GetExchange() common.CryptoExchange {
	total := 0.0
	satoshis := 0.0
	balances, _ := b.GetBalances()
	for _, c := range balances {
		if c.Currency == "BTC" { // TODO
			total += c.Total
		} else {
			satoshis += c.Price * c.Balance
			total += c.Total
		}
	}
	f, _ := strconv.ParseFloat(fmt.Sprintf("%.8f", satoshis), 64)
	t, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", total), 64)
	exchange := common.CryptoExchange{
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
	var cp *common.CurrencyPair
	if currencyPair.Quote == "USD" {
		cp = &common.CurrencyPair{
			Base:  currencyPair.Base,
			Quote: "USDT"}
	} else {
		cp = currencyPair
	}
	return cp
}

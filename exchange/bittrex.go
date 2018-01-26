package exchange

import (
	"fmt"
	"strconv"
	"time"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
	logging "github.com/op/go-logging"
	bittrex "github.com/toorop/go-bittrex"
)

type Bittrex struct {
	client       *bittrex.Bittrex
	logger       *logging.Logger
	price        float64
	satoshis     float64
	name         string
	currencyPair *common.CurrencyPair
	tradingFee   float64
	tradePairs   []common.CurrencyPair
	common.Exchange
}

func NewBittrex(btx *dao.UserCoinExchange, logger *logging.Logger, currencyPair *common.CurrencyPair) common.Exchange {
	var cp *common.CurrencyPair
	if currencyPair.Quote == "USD" {
		cp = &common.CurrencyPair{
			Base:  "USDT",
			Quote: currencyPair.Base}
	} else {
		cp = currencyPair
	}
	return &Bittrex{
		client:       bittrex.New(btx.Key, btx.Secret),
		logger:       logger,
		name:         "bittrex",
		currencyPair: cp,
		tradingFee:   .025}
}

func (b *Bittrex) SubscribeToLiveFeed(priceChange chan common.PriceChange) {
	for {
		symbol := b.FormattedCurrencyPair()
		time.Sleep(10 * time.Second)
		ticker, err := b.client.GetTicker(symbol)
		if err != nil {
			b.logger.Errorf("[Bittrex.SubscribeToLiveFeed] %s", err.Error())
			continue
		}
		f, _ := ticker.Last.Float64()
		if f <= 0 {
			b.logger.Errorf("[Bittrex.SubscribeToLiveFeed] Unable to get ticker data for %s", symbol)
			continue
		}
		b.logger.Debugf("[Bittrex.SubscribeToLiveFeed] Sending live price: %.8f", f)
		b.satoshis = f
		priceChange <- common.PriceChange{
			CurrencyPair: b.currencyPair,
			Satoshis:     b.satoshis,
			Price:        f}
	}
}

func (b *Bittrex) GetPriceHistory(start, end time.Time, granularity int) []common.Candlestick {
	b.logger.Debug("[Bittrex.GetTradeHistory] Getting trade history")
	candlesticks := make([]common.Candlestick, 0)
	marketHistory, err := b.client.GetMarketHistory(b.FormattedCurrencyPair())
	if err != nil {
		b.logger.Errorf("[Bittrex.GetTradeHistory] %s", err.Error())
	}
	for _, m := range marketHistory {
		f, _ := m.Price.Float64()
		if err != nil {
			b.logger.Errorf("[Bittrex.GetTradeHistory] %s", err.Error())
		}
		candlesticks = append(candlesticks, common.Candlestick{Close: f})
	}
	return candlesticks
}

func (b *Bittrex) GetOrderHistory() []common.Order {
	var orders []common.Order
	orderHistory, err := b.client.GetOrderHistory(b.FormattedCurrencyPair())
	if err != nil {
		b.logger.Errorf("[Bittrex.GetOrderHistory] %s", err.Error())
	}
	for _, o := range orderHistory {
		q, _ := o.Quantity.Float64()
		p, _ := o.Price.Float64()
		orders = append(orders, common.Order{
			Exchange: "bittrex",
			Date:     o.TimeStamp.Time,
			Type:     o.OrderType,
			Currency: b.currencyPair,
			Quantity: q,
			Price:    p})
	}
	return orders
}

func (b *Bittrex) GetBalances() ([]common.Coin, float64) {
	var coins []common.Coin
	sum := 0.0
	balances, err := b.client.GetBalances()
	if err != nil {
		b.logger.Errorf("[Bittrex.GetBalances] %s", err.Error())
	}
	for _, bal := range balances {
		var currency string
		if bal.Currency == "BTC" && b.currencyPair.Quote == bal.Currency {
			currency = fmt.Sprintf("%s-%s", b.currencyPair.Base, bal.Currency)
		} else {
			currency = fmt.Sprintf("%s-%s", b.currencyPair.Quote, bal.Currency)
		}
		b.logger.Debugf("[Bittrex.GetBalances] Getting %s ticker", currency)
		ticker, terr := b.client.GetTicker(currency)
		if terr != nil {
			b.logger.Errorf("[Bittrex.GetBalances] %s", terr.Error())
			continue
		}
		price, exact := ticker.Last.Float64()
		if !exact {
			b.logger.Error("[Bittrex.GetBalances] Conversion of Ticker price to Float64 was not exact!")
		}
		avail, exact := bal.Available.Float64()
		if !exact {
			b.logger.Error("[Bittrex.GetBalances] Conversion of Available funds to Float64 was not exact!")
		}
		balance, exact := bal.Balance.Float64()
		if !exact {
			b.logger.Error("[Bittrex.GetBalances] Conversion of Balance to Float64 was not exact!")
		}
		pending, exact := bal.Pending.Float64()
		if !exact {
			b.logger.Error("[Bittrex.GetBalances] Conversion of Pending to Float64 was not exact!")
		}
		if balance <= 0 {
			continue
		}

		bitcoinPrice := b.getBitcoinPrice()

		if bal.Currency == "BTC" {
			total := balance * bitcoinPrice
			t, err := strconv.ParseFloat(strconv.FormatFloat(total, 'f', 2, 64), 64)
			if err != nil {
				b.logger.Errorf("[Binance.GetBalances] %s", err.Error())
			}
			sum += t
			coins = append(coins, common.Coin{
				Currency:  bal.Currency,
				Available: balance,
				Balance:   balance,
				Price:     bitcoinPrice,
				Total:     t})
			continue
		} else {
			total := balance * (price * bitcoinPrice)
			t, err := strconv.ParseFloat(fmt.Sprintf("%.2f", total), 64)
			if err != nil {
				b.logger.Errorf("[Bittrex.GetBalances] %s", err.Error())
			}
			sum += t
			coins = append(coins, common.Coin{
				Address:   bal.CryptoAddress,
				Available: avail,
				Balance:   balance,
				Currency:  bal.Currency,
				Pending:   pending,
				Price:     price, // BTC satoshis, not actual USD price
				Total:     t})
		}

	}
	return coins, sum
}

func (b *Bittrex) getBitcoinPrice() float64 {
	summary, err := b.client.GetMarketSummary(fmt.Sprintf("%s-BTC", b.currencyPair.Base))
	if err != nil {
		b.logger.Errorf("[Bittrex.getBitcoinPrice] %s", err.Error())
	}
	if len(summary) != 1 {
		b.logger.Error("[Bittrex.getBitcoinPrice] Unexpected number of BTC markets returned")
		return 0
	}
	price, exact := summary[0].Last.Float64()
	if !exact {
		b.logger.Error("[Bittrex.getBitcoinPrice] Conversion of BTC price to Float64 was not exact!")
	}
	return price
}

func (b *Bittrex) GetExchangeAsync(exchangeChan *chan common.CoinExchange) {
	go func() { *exchangeChan <- b.GetExchange() }()
}

func (b *Bittrex) GetExchange() common.CoinExchange {
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
	exchange := common.CoinExchange{
		Name:     b.name,
		URL:      "https://www.bittrex.com",
		Total:    t,
		Satoshis: f,
		Coins:    balances}
	return exchange
}

func (b *Bittrex) GetCurrencyPair() common.CurrencyPair {
	return *b.currencyPair
}

func (b *Bittrex) GetName() string {
	return b.name
}

func (b *Bittrex) ToUSD(price, satoshis float64) float64 {
	return satoshis * price
}

func (b *Bittrex) FormattedCurrencyPair() string {
	return fmt.Sprintf("%s-%s", b.currencyPair.Base, b.currencyPair.Quote)
}

func (b *Bittrex) GetTradingFee() float64 {
	return b.tradingFee
}

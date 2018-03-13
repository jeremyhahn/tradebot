package exchange

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dto"
	"github.com/jeremyhahn/tradebot/entity"
	"github.com/jeremyhahn/tradebot/util"
	logging "github.com/op/go-logging"
	bittrex "github.com/toorop/go-bittrex"
	"golang.org/x/text/encoding/unicode"
)

type Bittrex struct {
	ctx                 common.Context
	client              *bittrex.Bittrex
	logger              *logging.Logger
	price               float64
	satoshis            float64
	name                string
	tradingFee          float64
	tradePairs          []common.CurrencyPair
	priceHistoryService common.PriceHistoryService
	common.Exchange
}

var BITTREX_RATE_LIMITER = common.NewRateLimiter(1, 1)

func NewBittrex(ctx common.Context, btx entity.UserExchangeEntity, priceHistoryService common.PriceHistoryService) common.Exchange {
	return &Bittrex{
		ctx:                 ctx,
		client:              bittrex.New(btx.GetKey(), btx.GetSecret()),
		logger:              ctx.GetLogger(),
		name:                "bittrex",
		tradingFee:          .025,
		priceHistoryService: priceHistoryService}
}

func (b *Bittrex) SubscribeToLiveFeed(currencyPair *common.CurrencyPair, priceChange chan common.PriceChange) {
	for {
		symbol := b.FormattedCurrencyPair(currencyPair)
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
			CurrencyPair: currencyPair,
			Satoshis:     b.satoshis,
			Price:        f}
	}
}

func (b *Bittrex) GetPriceHistory(currencyPair *common.CurrencyPair,
	start, end time.Time, granularity int) []common.Candlestick {

	BITTREX_RATE_LIMITER.RespectRateLimit()

	b.logger.Debug("[Bittrex.GetPriceHistory] Getting price history")
	candlesticks := make([]common.Candlestick, 0)
	marketHistory, err := b.client.GetMarketHistory(b.FormattedCurrencyPair(currencyPair))
	if err != nil {
		b.logger.Errorf("[Bittrex.GetPriceHistory] %s", err.Error())
	}
	for _, m := range marketHistory {
		f, _ := m.Price.Float64()
		if err != nil {
			b.logger.Errorf("[Bittrex.GetPriceHistory] %s", err.Error())
		}
		candlesticks = append(candlesticks, common.Candlestick{Close: f})
	}
	return candlesticks
}

func (b *Bittrex) GetOrderHistory(currencyPair *common.CurrencyPair) []common.Order {
	BITTREX_RATE_LIMITER.RespectRateLimit()
	formattedCurrencyPair := b.FormattedCurrencyPair(currencyPair)
	b.logger.Debugf("[Bittrex.GetOrderHistory] Getting %s order history", formattedCurrencyPair)
	var orders []common.Order
	orderHistory, err := b.client.GetOrderHistory(formattedCurrencyPair)
	if err != nil {
		b.logger.Errorf("[Bittrex.GetOrderHistory] Error: %s", err.Error())
	}
	if len(orderHistory) == 0 {
		b.logger.Warning("[Bittrex.GetOrderHistory] Zero records returned from Bittrex API")
	}
	for _, o := range orderHistory {

		util.DUMP(o)

		q, _ := o.Quantity.Float64()
		p, _ := o.Price.Float64()
		orders = append(orders, &dto.OrderDTO{
			Id:                 o.OrderUuid,
			Exchange:           "bittrex",
			Date:               o.TimeStamp.Time,
			Type:               o.OrderType,
			CurrencyPair:       currencyPair,
			Quantity:           q,
			QuantityCurrency:   currencyPair.Quote,
			Price:              p,
			PriceCurrency:      currencyPair.Base,
			Total:              q * p,
			TotalCurrency:      currencyPair.Base,
			FeeCurrency:        currencyPair.Base,
			HistoricalPrice:    b.priceHistoryService.GetClosePriceOn(currencyPair.Quote, o.TimeStamp.Time),
			HistoricalCurrency: b.ctx.GetUser().GetLocalCurrency()})
	}
	return orders
}

func (b *Bittrex) GetBalances() ([]common.Coin, float64) {
	BITTREX_RATE_LIMITER.RespectRateLimit()
	var coins []common.Coin
	sum := 0.0
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
			coins = append(coins, &dto.CoinDTO{
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
			coins = append(coins, &dto.CoinDTO{
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
	BITTREX_RATE_LIMITER.RespectRateLimit()
	currencyPair := &common.CurrencyPair{
		Base:          "BTC",
		Quote:         b.ctx.GetUser().GetLocalCurrency(),
		LocalCurrency: b.ctx.GetUser().GetLocalCurrency()}
	localizedCurrencyPair := b.localizedCurrencyPair(currencyPair)
	symbol := fmt.Sprintf("%s-%s", localizedCurrencyPair.Base, localizedCurrencyPair.Quote)
	summary, err := b.client.GetMarketSummary(symbol)
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

func (b *Bittrex) GetSummary() common.CryptoExchangeSummary {
	BITTREX_RATE_LIMITER.RespectRateLimit()
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
		URL:      "https://www.bittrex.com",
		Total:    t,
		Satoshis: f,
		Coins:    balances}
	return exchange
}

func (b *Bittrex) GetName() string {
	return b.name
}

func (b *Bittrex) GetTradingFee() float64 {
	return b.tradingFee
}

func (b *Bittrex) ParseImport(file string) ([]common.Order, error) {
	var orders []common.Order
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
		qty, err := strconv.ParseFloat(values[3], 64)
		if err != nil {
			b.ctx.GetLogger().Errorf("[Bittrex.ParseImport] Error parsing quantity: %s", err.Error())
			return orders, err
		}
		price, err := strconv.ParseFloat(values[4], 64)
		if err != nil {
			b.ctx.GetLogger().Errorf("[Bittrex.ParseImport] Error parsing price: %s", err.Error())
			return orders, err
		}
		fee, err := strconv.ParseFloat(values[5], 64)
		if err != nil {
			b.ctx.GetLogger().Errorf("[Bittrex.ParseImport] Error parsing fee: %s", err.Error())
			return orders, err
		}
		total, err := strconv.ParseFloat(values[6], 64)
		if err != nil {
			b.ctx.GetLogger().Errorf("[Bittrex.ParseImport] Error parsing totalt: %s", err.Error())
			return orders, err
		}
		date, err := time.Parse("1/2/2006 15:04:05 PM", values[8])
		if err != nil {
			b.ctx.GetLogger().Errorf("[Bittrex.ParseImport] Error parsing float: %s", err.Error())
			return orders, err
		}
		var orderType string
		if values[2] == "LIMIT_BUY" {
			orderType = "buy"
		} else {
			orderType = "sell"
		}
		currencyPair := common.NewCurrencyPair(values[1], b.ctx.GetUser().GetLocalCurrency())
		order := &dto.OrderDTO{
			Id:                 fmt.Sprintf("%d", b.ctx.GetUser().GetId()),
			Exchange:           b.name,
			Date:               date,
			Type:               orderType,
			CurrencyPair:       currencyPair,
			Quantity:           qty,
			QuantityCurrency:   currencyPair.Quote,
			Price:              price,
			Fee:                fee,
			Total:              total,
			PriceCurrency:      currencyPair.Base,
			FeeCurrency:        currencyPair.Base,
			TotalCurrency:      currencyPair.Base,
			HistoricalPrice:    b.priceHistoryService.GetClosePriceOn(currencyPair.Quote, date),
			HistoricalCurrency: b.ctx.GetUser().GetLocalCurrency()}
		orders = append(orders, order)
	}
	return orders, nil
}

func (b *Bittrex) FormattedCurrencyPair(currencyPair *common.CurrencyPair) string {
	cp := b.localizedCurrencyPair(currencyPair)
	return fmt.Sprintf("%s-%s", cp.Base, cp.Quote)
}

func (b *Bittrex) localizedCurrencyPair(currencyPair *common.CurrencyPair) *common.CurrencyPair {
	var cp *common.CurrencyPair
	if currencyPair.Quote == "USD" {
		cp = &common.CurrencyPair{
			Base:          "USDT",
			Quote:         currencyPair.Base,
			LocalCurrency: "USDT"}
	} else {
		cp = currencyPair
	}
	return cp
}

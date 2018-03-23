package main

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
	"github.com/shopspring/decimal"
	bittrex "github.com/toorop/go-bittrex"
	"golang.org/x/text/encoding/unicode"
)

type Bittrex struct {
	ctx         common.Context
	client      *bittrex.Bittrex
	logger      *logging.Logger
	price       float64
	satoshis    float64
	name        string
	displayName string
	tradingFee  float64
	tradePairs  []common.CurrencyPair
	usdtMarkets []string
	common.Exchange
}

var BITTREX_RATE_LIMITER = common.NewRateLimiter(1, 1)

func CreateBittrex(ctx common.Context, userExchangeEntity entity.UserExchangeEntity) common.Exchange {

	ctx.GetLogger().Debugf("key=%s", userExchangeEntity.GetKey())
	ctx.GetLogger().Debugf("secret=%s", userExchangeEntity.GetSecret())

	return &Bittrex{
		ctx:         ctx,
		client:      bittrex.New(userExchangeEntity.GetKey(), userExchangeEntity.GetSecret()),
		logger:      ctx.GetLogger(),
		name:        "Bittrex",
		displayName: "Bittrex",
		tradingFee:  .025,
		usdtMarkets: []string{"BTC", "ETH", "XRP", "NEO", "ADA", "LTC",
			"BCC", "OMG", "ETC", "ZEC", "XMR", "XVG", "BTG", "DASH", "NXT"}}
}

func (b *Bittrex) GetPriceAt(currency string, atDate time.Time) (*common.Candlestick, error) {
	currencyPair := &common.CurrencyPair{
		Base:          currency,
		Quote:         "USDT",
		LocalCurrency: b.ctx.GetUser().GetLocalCurrency()}
	kline, err := b.GetPriceHistory(currencyPair, atDate.Add(-1*time.Minute), atDate.Add(1*time.Minute), 1)
	if err != nil {
		return &common.Candlestick{}, err
	}
	closestCandle, err := util.FindClosestDatedCandle(b.ctx.GetLogger(), atDate, &kline)
	if err != nil {
		return &common.Candlestick{}, err
	}
	return closestCandle, nil
}

func (b *Bittrex) toBTCtoUSD(currency string, atDate time.Time) (*common.Candlestick, error) {
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
	start, end time.Time, granularity int) ([]common.Candlestick, error) {

	BITTREX_RATE_LIMITER.RespectRateLimit()

	b.logger.Debug("[Bittrex.GetPriceHistory] Getting price history")
	candlesticks := make([]common.Candlestick, 0)
	marketHistory, err := b.client.GetMarketHistory(b.FormattedCurrencyPair(currencyPair))
	if err != nil {
		b.logger.Errorf("[Bittrex.GetPriceHistory] %s", err.Error())
		return nil, err
	}
	for _, m := range marketHistory {
		f, _ := m.Price.Float64()
		if err != nil {
			b.logger.Errorf("[Bittrex.GetPriceHistory] Error converting market price to float: %s", err.Error())
		}
		candlesticks = append(candlesticks, common.Candlestick{Close: f})
	}
	return candlesticks, nil
}

func (b *Bittrex) GetOrderHistory(currencyPair *common.CurrencyPair) []common.Transaction {
	BITTREX_RATE_LIMITER.RespectRateLimit()
	formattedCurrencyPair := b.FormattedCurrencyPair(currencyPair)
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

		util.DUMP(o)

		qty := o.Quantity
		price := o.Price
		fee := o.Commission
		total := qty.Mul(price)

		orders = append(orders, &dto.TransactionDTO{
			Id:                   o.OrderUuid,
			Network:              b.name,
			NetworkDisplayName:   b.displayName,
			Date:                 o.TimeStamp.Time,
			Type:                 o.OrderType,
			CurrencyPair:         currencyPair,
			Quantity:             qty.StringFixed(8),
			QuantityCurrency:     currencyPair.Quote,
			FiatQuantity:         "0.00",
			FiatQuantityCurrency: "N/A",
			Price:                price.StringFixed(8),
			PriceCurrency:        currencyPair.Base,
			FiatPrice:            "0.00",
			FiatPriceCurrency:    "N/A",
			Fee:                  fee.StringFixed(8),
			FeeCurrency:          currencyPair.Base,
			FiatFee:              "0.00",
			FiatFeeCurrency:      "N/A",
			Total:                total.StringFixed(8),
			TotalCurrency:        currencyPair.Base,
			FiatTotal:            "0.00",
			FiatTotalCurrency:    "N/A"})
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
		currencyPair := &common.CurrencyPair{
			Base:          deposit.Currency,
			Quote:         "USDT",
			LocalCurrency: b.ctx.GetUser().GetLocalCurrency()}
		orderDate := deposit.LastUpdated.Time

		usdtCandle, err := b.GetPriceAt(deposit.Currency, orderDate)
		if err != nil {
			return nil, err
		}
		amount := deposit.Amount
		price := amount.Mul(decimal.NewFromFloat(usdtCandle.Close))
		_deposits = append(_deposits, &dto.TransactionDTO{
			Id:                 string(deposit.Id),
			Type:               common.DEPOSIT_ORDER_TYPE,
			Date:               orderDate,
			Network:            b.name,
			NetworkDisplayName: b.displayName,
			CurrencyPair:       currencyPair,
			Quantity:           amount.StringFixed(8),
			QuantityCurrency:   deposit.Currency,
			Price:              price.StringFixed(8),
			PriceCurrency:      currencyPair.Quote,
			Fee:                "0.0",
			FeeCurrency:        currencyPair.Quote,
			Total:              price.StringFixed(8),
			TotalCurrency:      currencyPair.Quote})
		b.ctx.GetLogger().Debug(deposit)
	}
	return _deposits, nil
}

func (b *Bittrex) GetWithdrawalHistory() ([]common.Transaction, error) {
	var orders []common.Transaction
	currencies, _ := b.GetCurrencies()
	b.ctx.GetLogger().Debugf("[Bittrex.GetWithdrawalHistory] Withdrawal currencies for %s: %s", b.ctx.GetUser().GetUsername(), currencies)
	for _, currency := range currencies {
		withdraws, err := b.client.GetWithdrawalHistory(currency.GetSymbol())
		if err != nil {
			b.logger.Errorf("[Bittrex.GetWithdrawalHistory] Error: %s", err.Error())
		}
		for _, withdraw := range withdraws {
			if !withdraw.Authorized {
				continue
			}
			orderDate := withdraw.Opened.Time
			currencyPair := &common.CurrencyPair{
				Base:          withdraw.Currency,
				Quote:         withdraw.Currency,
				LocalCurrency: b.ctx.GetUser().GetLocalCurrency()}
			/*
				currencyPair := &common.CurrencyPair{
					Base:          withdraw.Currency,
					Quote:         "USDT",
					LocalCurrency: b.ctx.GetUser().GetLocalCurrency()}
					var priceCandle *common.Candlestick
					if b.supportsUSDT(withdraw.Currency) {
						usdtCandle, err := b.GetPriceAt(withdraw.Currency, orderDate)
						if err != nil {
							return nil, err
						}
						priceCandle = usdtCandle
					} else {
						btcToUsdCandle, err := b.toBTCtoUSD(withdraw.Currency, orderDate)
						if err != nil {
							return nil, err
						}
						priceCandle = btcToUsdCandle
					}*/
			amount := withdraw.Amount
			txCost := withdraw.TxCost
			total := amount.Add(txCost)
			orders = append(orders, &dto.TransactionDTO{
				Id:                   string(withdraw.PaymentUuid),
				Type:                 common.WITHDRAWAL_ORDER_TYPE,
				Date:                 orderDate,
				Network:              b.name,
				NetworkDisplayName:   b.displayName,
				CurrencyPair:         currencyPair,
				Quantity:             amount.StringFixed(8),
				QuantityCurrency:     withdraw.Currency,
				FiatQuantity:         "0.00",
				FiatQuantityCurrency: "N/A",
				Price:                "0.0",
				PriceCurrency:        withdraw.Currency,
				FiatPrice:            "0.00",
				FiatPriceCurrency:    "N/A",
				Fee:                  txCost.StringFixed(8),
				FeeCurrency:          withdraw.Currency,
				FiatFee:              "0.00",
				FiatFeeCurrency:      "N/A",
				Total:                total.StringFixed(8),
				TotalCurrency:        withdraw.Currency,
				FiatTotal:            "0.00",
				FiatTotalCurrency:    "N/A"})
		}
	}
	return orders, nil
}

func (b *Bittrex) supportsUSDT(currency string) bool {
	for _, market := range b.usdtMarkets {
		if market == currency {
			return true
		}
	}
	return false
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
				b.logger.Errorf("[Bittrex.GetBalances] %s", err.Error())
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

func (b *Bittrex) GetCurrencies() (map[string]*common.Currency, error) {
	b.ctx.GetLogger().Debugf("[Bittrex.GetCurrencies] Getting currency list")
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

func (b *Bittrex) GetDisplayName() string {
	return b.displayName
}

func (b *Bittrex) GetTradingFee() float64 {
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
			orderType = "buy"
		} else {
			orderType = "sell"
		}
		qty := values[3]
		price := values[4]
		fee := values[5]
		total := values[6]
		date, err := time.Parse("1/2/2006 15:04:05 PM", values[8])
		if err != nil {
			b.ctx.GetLogger().Errorf("[Bittrex.ParseImport] Error parsing float: %s", err.Error())
			return orders, err
		}
		currencyPair, err := common.NewCurrencyPair(values[1], b.ctx.GetUser().GetLocalCurrency())
		if err != nil {
			b.ctx.GetLogger().Errorf("[Bittrex.ParseImport] Invalid currency pair: %s", values[1])
			return nil, err
		}
		orders = append(orders, &dto.TransactionDTO{
			Id:                   fmt.Sprintf("%d", b.ctx.GetUser().GetId()),
			Network:              b.name,
			NetworkDisplayName:   b.displayName,
			Date:                 date,
			Type:                 orderType,
			CurrencyPair:         currencyPair,
			Quantity:             qty,
			QuantityCurrency:     currencyPair.Quote,
			FiatQuantity:         "0.00",
			FiatQuantityCurrency: "N/A",
			Price:                price,
			PriceCurrency:        currencyPair.Quote,
			FiatPrice:            "0.00",
			FiatPriceCurrency:    "N/A",
			Fee:                  fee,
			FeeCurrency:          currencyPair.Quote,
			FiatFee:              "0.00",
			FiatFeeCurrency:      "N/A",
			Total:                total,
			TotalCurrency:        currencyPair.Quote,
			FiatTotal:            "0.00",
			FiatTotalCurrency:    "N/A"})

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

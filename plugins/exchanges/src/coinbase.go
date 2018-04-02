package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/Zauberstuhl/go-coinbase"
	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dto"
	"github.com/jeremyhahn/tradebot/entity"
	"github.com/jeremyhahn/tradebot/util"
	cache "github.com/patrickmn/go-cache"
	"github.com/shopspring/decimal"
)

var COINBASE_RATELIMITER = common.NewRateLimiter(10, 1)

type Coinbase struct {
	ctx              common.Context
	name             string
	displayName      string
	client           coinbase.APIClient
	STATUS_COMPLETED string
	TIME_FORMAT      string
	cache            *cache.Cache
	common.Exchange
}

func main() {}

func CreateCoinbase(ctx common.Context, userExchangeEntity entity.UserExchangeEntity) common.Exchange {
	return &Coinbase{
		ctx:         ctx,
		name:        "Coinbase",
		displayName: "Coinbase",
		client: coinbase.APIClient{
			Key:    userExchangeEntity.GetKey(),
			Secret: userExchangeEntity.GetSecret()},
		STATUS_COMPLETED: "completed",
		TIME_FORMAT:      "2006-01-02T15:04:05Z",
		cache:            cache.New(1*time.Minute, 1*time.Minute)}
}

func (cb *Coinbase) GetName() string {
	return cb.name
}

func (cb *Coinbase) GetDisplayName() string {
	return cb.displayName
}

func (cb *Coinbase) GetBalances() ([]common.Coin, decimal.Decimal) {
	balancesCacheKey := fmt.Sprintf("%d-%s", cb.ctx.GetUser().GetId(), "-coinbase-balances")
	sumCacheKey := fmt.Sprintf("%d-%s", cb.ctx.GetUser().GetId(), "-coinbase-sum")
	if balances, found := cb.cache.Get(balancesCacheKey); found {
		sum, _ := cb.cache.Get(sumCacheKey)
		b := balances.(*[]common.Coin)
		return *b, *sum.(*decimal.Decimal)
	}
	var coins []common.Coin
	var sum decimal.Decimal
	accounts, _ := cb.client.Accounts()
	for _, acct := range accounts.Data {
		var decimalPlaces int32
		if _, ok := common.FiatCurrencies[acct.Currency]; ok {
			decimalPlaces = 2
		} else {
			decimalPlaces = 8
		}
		zero := decimal.NewFromFloat(0)
		price := zero
		balance := decimal.NewFromFloat(acct.Balance.Amount)
		nativeAmount := decimal.NewFromFloat(acct.Native_balance.Amount)
		sum = sum.Add(nativeAmount)
		if balance.GreaterThan(decimal.NewFromFloat(0)) && nativeAmount.GreaterThan(decimal.NewFromFloat(0)) {
			price = balance.Div(nativeAmount)
		}
		coins = append(coins, &dto.CoinDTO{
			Currency:  acct.Currency,
			Balance:   balance.Truncate(decimalPlaces),
			Available: balance.Truncate(decimalPlaces),
			Price:     price.Truncate(2),
			Total:     nativeAmount.Truncate(decimalPlaces)})
	}
	cb.cache.Set(balancesCacheKey, &coins, cache.DefaultExpiration)
	cb.cache.Set(sumCacheKey, &sum, cache.DefaultExpiration)
	return coins, sum
}

func (cb *Coinbase) GetOrderHistory(currencyPair *common.CurrencyPair) []common.Transaction {
	var accountId string
	var txs []common.Transaction
	accounts, _ := cb.client.Accounts()
	for _, acct := range accounts.Data {
		if acct.Currency == currencyPair.Base {
			accountId = acct.Id
			break
		}
	}
	if accountId == "" {
		return txs
	}

	cb.ctx.GetLogger().Debugf("[Coinbase.GetOrderHistory] Getting order history for %s", currencyPair.Base)
	acctId := coinbase.AccountId(accountId)
	buys, err := cb.client.ListBuys(acctId)
	if err != nil {
		cb.ctx.GetLogger().Debugf("[Coinbase.GetOrderHistory] Buy error: %s", err.Error())
	}
	for _, buy := range buys.Data {
		if buy.Status != cb.STATUS_COMPLETED {
			continue
		}
		createdAt, err := time.Parse(cb.TIME_FORMAT, buy.Created_at)
		if err != nil {
			cb.ctx.GetLogger().Debugf("[Coinbase.GetOrderHistory] Error parsing Created_at: %s", err.Error())
		}
		baseCurrency, err := cb.getCurrency(currencyPair.Base)
		if err != nil {
			cb.ctx.GetLogger().Errorf("[Coinbase.GetOrderHistory] Error getting buy base currency: %s", err.Error())
			continue
		}
		quoteCurrency, err := cb.getCurrency(currencyPair.Quote)
		if err != nil {
			cb.ctx.GetLogger().Errorf("[Coinbase.GetOrderHistory] Error getting buy quote currency: %s", err.Error())
			continue
		}
		quantity := decimal.NewFromFloat(buy.Amount.Amount)
		fee := decimal.NewFromFloat(buy.Fee.Amount)
		total := decimal.NewFromFloat(buy.Total.Amount)
		fiatPrice := total.Mul(quantity)
		txs = append(txs, &dto.TransactionDTO{
			Id:                   buy.Transaction.Id,
			Type:                 buy.Resource,
			Date:                 createdAt,
			Network:              cb.name,
			NetworkDisplayName:   cb.displayName,
			CurrencyPair:         currencyPair,
			Quantity:             quantity.StringFixed(baseCurrency.GetDecimalPlace()),
			QuantityCurrency:     buy.Amount.Currency,
			FiatQuantity:         total.StringFixed(quoteCurrency.GetDecimalPlace()),
			FiatQuantityCurrency: buy.Total.Currency,
			Price:                fiatPrice.StringFixed(quoteCurrency.GetDecimalPlace()),
			PriceCurrency:        buy.Total.Currency,
			FiatPrice:            fiatPrice.StringFixed(quoteCurrency.GetDecimalPlace()),
			FiatPriceCurrency:    buy.Total.Currency,
			Fee:                  fee.StringFixed(quoteCurrency.GetDecimalPlace()),
			FeeCurrency:          buy.Total.Currency,
			FiatFee:              fee.StringFixed(quoteCurrency.GetDecimalPlace()),
			FiatFeeCurrency:      currencyPair.Quote,
			Total:                quantity.StringFixed(baseCurrency.GetDecimalPlace()),
			TotalCurrency:        buy.Amount.Currency,
			FiatTotal:            total.StringFixed(quoteCurrency.GetDecimalPlace()),
			FiatTotalCurrency:    buy.Total.Currency})
	}

	sells, err := cb.client.ListSells(acctId)
	if err != nil {
		cb.ctx.GetLogger().Debugf("[Coinbase.GetOrderHistory] Sell error: %s", err.Error())
	}
	for _, sell := range sells.Data {
		if sell.Status != cb.STATUS_COMPLETED {
			continue
		}
		createdAt, err := time.Parse(cb.TIME_FORMAT, sell.Created_at)
		if err != nil {
			cb.ctx.GetLogger().Debugf("[Coinbase.GetOrderHistory] Error parsing sell Created_at: %s", err.Error())
		}

		baseCurrency, err := cb.getCurrency(currencyPair.Base)
		if err != nil {
			cb.ctx.GetLogger().Errorf("[Coinbase.GetOrderHistory] Error getting sell base currency: %s", err.Error())
			continue
		}
		quoteCurrency, err := cb.getCurrency(currencyPair.Quote)
		if err != nil {
			cb.ctx.GetLogger().Errorf("[Coinbase.GetOrderHistory] Error getting sell quote currency: %s", err.Error())
			continue
		}
		quantity := decimal.NewFromFloat(sell.Amount.Amount)
		price := decimal.NewFromFloat(sell.Amount.Amount)
		fee := decimal.NewFromFloat(sell.Fee.Amount)
		total := decimal.NewFromFloat(sell.Total.Amount)
		fiatPrice := price.Mul(quantity)
		txs = append(txs, &dto.TransactionDTO{
			Id:                   sell.Transaction.Id,
			Type:                 sell.Resource,
			Date:                 createdAt,
			Network:              cb.name,
			NetworkDisplayName:   cb.displayName,
			CurrencyPair:         currencyPair,
			Quantity:             quantity.StringFixed(baseCurrency.GetDecimalPlace()),
			QuantityCurrency:     sell.Amount.Currency,
			FiatQuantity:         total.Sub(fee).StringFixed(quoteCurrency.GetDecimalPlace()),
			FiatQuantityCurrency: currencyPair.Quote,
			Price:                price.StringFixed(quoteCurrency.GetDecimalPlace()),
			PriceCurrency:        currencyPair.Quote,
			FiatPrice:            fiatPrice.StringFixed(quoteCurrency.GetDecimalPlace()),
			FiatPriceCurrency:    currencyPair.Quote,
			Fee:                  fee.StringFixed(quoteCurrency.GetDecimalPlace()),
			FeeCurrency:          sell.Amount.Currency,
			FiatFee:              fee.StringFixed(quoteCurrency.GetDecimalPlace()),
			FiatFeeCurrency:      currencyPair.Quote,
			Total:                total.StringFixed(baseCurrency.GetDecimalPlace()),
			TotalCurrency:        sell.Total.Currency,
			FiatTotal:            total.StringFixed(quoteCurrency.GetDecimalPlace()),
			FiatTotalCurrency:    currencyPair.Quote})
	}

	return txs
}

func (cb *Coinbase) GetDepositHistory() ([]common.Transaction, error) {
	cb.ctx.GetLogger().Debugf("[Coinbase.GetDepositHistory] Getting withdrawal history for %s", cb.ctx.GetUser().GetUsername())
	var _deposits []common.Transaction
	accounts, err := cb.client.Accounts()
	if err != nil {
		return nil, err
	}
	for _, acct := range accounts.Data {
		acctId := coinbase.AccountId(acct.Id)
		deposits, err := cb.client.ListDeposits(acctId)
		if err != nil {
			return nil, err
		}
		for _, deposit := range deposits.Data {
			if deposit.Status != cb.STATUS_COMPLETED {
				continue
			}
			createdAt, err := time.Parse(cb.TIME_FORMAT, deposit.Created_at)
			if err != nil {
				cb.ctx.GetLogger().Debugf("[Coinbase.GetDepositHistory] Error parsing Created_at: %s", err.Error())
			}
			currencyPair := &common.CurrencyPair{
				Base:          deposit.Amount.Currency,
				Quote:         deposit.Subtotal.Currency,
				LocalCurrency: cb.ctx.GetUser().GetLocalCurrency()}
			baseCurrency, err := cb.getCurrency(currencyPair.Base)
			if err != nil {
				cb.ctx.GetLogger().Errorf("[Coinbase.GetDepositHistory] Error getting base currency: %s", err.Error())
				continue
			}
			quoteCurrency, err := cb.getCurrency(currencyPair.Quote)
			if err != nil {
				cb.ctx.GetLogger().Errorf("[Coinbase.GetDepositHistory] Error getting quote currency: %s", err.Error())
				continue
			}
			quantity := decimal.NewFromFloat(deposit.Amount.Amount)
			price := decimal.NewFromFloat(deposit.Amount.Amount)
			fee := decimal.NewFromFloat(deposit.Fee.Amount)
			total := decimal.NewFromFloat(deposit.Subtotal.Amount)
			_deposits = append(_deposits, &dto.TransactionDTO{
				Id:                   deposit.Transaction.Id,
				Type:                 common.DEPOSIT_ORDER_TYPE,
				Date:                 createdAt,
				Network:              cb.name,
				NetworkDisplayName:   cb.displayName,
				CurrencyPair:         currencyPair,
				Quantity:             quantity.StringFixed(baseCurrency.GetDecimalPlace()),
				QuantityCurrency:     deposit.Amount.Currency,
				FiatQuantity:         total.Sub(fee).StringFixed(quoteCurrency.GetDecimalPlace()),
				FiatQuantityCurrency: currencyPair.Quote,
				Price:                price.StringFixed(quoteCurrency.GetDecimalPlace()),
				PriceCurrency:        currencyPair.Quote,
				FiatPrice:            price.StringFixed(quoteCurrency.GetDecimalPlace()),
				FiatPriceCurrency:    currencyPair.Quote,
				Fee:                  fee.StringFixed(quoteCurrency.GetDecimalPlace()),
				FeeCurrency:          currencyPair.Quote,
				FiatFee:              fee.StringFixed(quoteCurrency.GetDecimalPlace()),
				FiatFeeCurrency:      currencyPair.Quote,
				Total:                total.StringFixed(baseCurrency.GetDecimalPlace()),
				TotalCurrency:        deposit.Subtotal.Currency,
				FiatTotal:            total.StringFixed(quoteCurrency.GetDecimalPlace()),
				FiatTotalCurrency:    currencyPair.Quote})
		}
	}
	return _deposits, nil
}

func (cb *Coinbase) GetWithdrawalHistory() ([]common.Transaction, error) {
	cb.ctx.GetLogger().Debugf("[Coinbase.GetWithdrawallHistory] Getting withdrawal history for %s", cb.ctx.GetUser().GetUsername())
	var _withdrawls []common.Transaction
	accounts, err := cb.client.Accounts()
	if err != nil {
		return nil, err
	}
	for _, acct := range accounts.Data {
		acctId := coinbase.AccountId(acct.Id)
		withdrawls, err := cb.client.ListWithdrawals(acctId)
		if err != nil {
			return nil, err
		}
		for _, withdrawl := range withdrawls.Data {
			if withdrawl.Status != cb.STATUS_COMPLETED {
				continue
			}
			createdAt, err := time.Parse(cb.TIME_FORMAT, withdrawl.Created_at)
			if err != nil {
				cb.ctx.GetLogger().Debugf("[Coinbase.GetWithdrawallHistory] Error parsing Created_at: %s", err.Error())
			}
			currencyPair := &common.CurrencyPair{
				Base:          withdrawl.Amount.Currency,
				Quote:         withdrawl.Subtotal.Currency,
				LocalCurrency: cb.ctx.GetUser().GetLocalCurrency()}
			baseCurrency, err := cb.getCurrency(currencyPair.Base)
			if err != nil {
				cb.ctx.GetLogger().Errorf("[Coinbase.GetWithdrawalHistory] Error getting base currency: %s", err.Error())
				continue
			}
			quoteCurrency, err := cb.getCurrency(currencyPair.Quote)
			if err != nil {
				cb.ctx.GetLogger().Errorf("[Coinbase.GetWithdrawalHistory] Error getting quote currency: %s", err.Error())
				continue
			}
			quantity := decimal.NewFromFloat(withdrawl.Amount.Amount)
			price := decimal.NewFromFloat(withdrawl.Amount.Amount)
			fee := decimal.NewFromFloat(withdrawl.Fee.Amount)
			total := decimal.NewFromFloat(withdrawl.Subtotal.Amount)
			_withdrawls = append(_withdrawls, &dto.TransactionDTO{
				Id:                   withdrawl.Transaction.Id,
				Type:                 common.WITHDRAWAL_ORDER_TYPE,
				Date:                 createdAt,
				Network:              cb.name,
				NetworkDisplayName:   cb.displayName,
				CurrencyPair:         currencyPair,
				Quantity:             quantity.StringFixed(baseCurrency.GetDecimalPlace()),
				QuantityCurrency:     withdrawl.Amount.Currency,
				FiatQuantity:         total.Sub(fee).StringFixed(2),
				FiatQuantityCurrency: currencyPair.Quote,
				Price:                price.StringFixed(quoteCurrency.GetDecimalPlace()),
				PriceCurrency:        currencyPair.Quote,
				FiatPrice:            price.StringFixed(2),
				FiatPriceCurrency:    currencyPair.Quote,
				Fee:                  fee.StringFixed(quoteCurrency.GetDecimalPlace()),
				FeeCurrency:          withdrawl.Subtotal.Currency,
				FiatFee:              fee.StringFixed(2),
				FiatFeeCurrency:      currencyPair.Quote,
				Total:                total.StringFixed(baseCurrency.GetDecimalPlace()),
				TotalCurrency:        withdrawl.Subtotal.Currency,
				FiatTotal:            total.StringFixed(quoteCurrency.GetDecimalPlace()),
				FiatTotalCurrency:    currencyPair.Quote})
		}
	}
	return _withdrawls, nil
}

func (cb *Coinbase) GetCurrencies() (map[string]*common.Currency, error) {
	cacheKey := fmt.Sprintf("%d-%s", cb.ctx.GetUser().GetId(), "-coinbase-currencies")
	if x, found := cb.cache.Get(cacheKey); found {
		currencies := x.(*map[string]*common.Currency)
		return *currencies, nil
	}
	cb.ctx.GetLogger().Debugf("[Coinbase.GetCurrencies] Getting Coinbase currencies")
	coins, _ := cb.GetBalances()
	currencies, err := cb.client.GetCurrencies()
	if err != nil {
		cb.ctx.GetLogger().Errorf("[Coinbase.GetCurrencies] Error: %s", err.Error())
		return nil, err
	}
	_currencies := make(map[string]*common.Currency, len(currencies.Data)+len(coins))
	for _, currency := range currencies.Data {
		_currencies[currency.Id] = &common.Currency{
			ID:           currency.Id,
			Symbol:       currency.Id,
			Name:         currency.Name,
			BaseUnit:     100,
			TxFee:        decimal.NewFromFloat(0.0),
			DecimalPlace: util.ParseDecimalPlace(decimal.NewFromFloat(currency.Min_size).String())}
	}
	for _, coin := range coins {
		baseUnit := int32(100000000)
		decimalPlace := int32(8)
		fiatCurrency, found := common.FiatCurrencies[coin.GetCurrency()]
		if found {
			baseUnit = fiatCurrency.GetBaseUnit()
			decimalPlace = fiatCurrency.GetDecimalPlace()
		}
		_currencies[coin.GetCurrency()] = &common.Currency{
			ID:           coin.GetCurrency(),
			Name:         common.CryptoNames[coin.GetCurrency()],
			BaseUnit:     baseUnit,
			TxFee:        decimal.NewFromFloat(0.0),
			DecimalPlace: decimalPlace}
	}
	cb.cache.Set(cacheKey, &_currencies, cache.DefaultExpiration)
	return _currencies, nil
}

func (cb *Coinbase) getCurrency(currency string) (*common.Currency, error) {
	currencies, err := cb.GetCurrencies()
	if err != nil {
		return nil, err
	}
	if currency, found := currencies[currency]; found {
		return currency, nil
	}
	return nil, errors.New(fmt.Sprintf("Currency not found: %s", currency))
}

func (cb *Coinbase) GetSummary() common.CryptoExchangeSummary {
	total := decimal.NewFromFloat(0)
	satoshis := decimal.NewFromFloat(0)
	balances, _ := cb.GetBalances()
	for _, c := range balances {
		if c.GetCurrency() == cb.ctx.GetUser().GetLocalCurrency() {
			total = total.Add(c.GetTotal())
		} else if c.IsBitcoin() {
			satoshis = satoshis.Add(c.GetBalance())
			total = total.Add(c.GetTotal())
		} else {
			COINBASE_RATELIMITER.RespectRateLimit()
			spotPrice, err := cb.client.GetSpotPrice(coinbase.ConfigPrice{Date: time.Now()})
			if err != nil {
				cb.ctx.GetLogger().Errorf("[Coinbase.GetExchange] %s", err.Error())
				continue
			}
			satoshis = satoshis.Add(decimal.NewFromFloat(spotPrice.Data.Amount))
			total = total.Add(c.GetTotal())
		}
	}
	exchange := &dto.CryptoExchangeSummaryDTO{
		Name:     cb.name,
		URL:      "https://www.gdax.com",
		Total:    total.Truncate(8),
		Satoshis: satoshis.Truncate(8),
		Coins:    balances}
	return exchange
}

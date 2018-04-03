package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/dto"
	"github.com/jeremyhahn/tradebot/util"
	"github.com/shopspring/decimal"
)

type RippleTx struct {
	Account     string `json:"Account"`
	Destination string `json:"Destination"`
	Amount      string `json:"Amount"`
	Fee         string `json:"Fee"`
}

type RippleTransaction struct {
	Hash string    `json:"hash"`
	Date time.Time `json:"date"`
	Tx   RippleTx  `json:"tx"`
}

type RippleResponse struct {
	Result       string              `json:"result"`
	Transactions []RippleTransaction `json:"transactions"`
}

type RippleExchangeRate struct {
	Result string `json:"result"`
	Rate   string `json:"rate"`
}

type RippleWallet struct {
	Balance  string `json:"value"`
	Currency string `json:"currency"`
}

type RippleBalance struct {
	Result      string         `json:"result"`
	LedgerIndex float64        `json:"ledger_index"`
	Limit       int64          `json:"limit"`
	Balances    []RippleWallet `json:"balances"`
}

type Ripple struct {
	ctx              common.Context
	client           http.Client
	userDAO          dao.UserDAO
	marketcapService common.MarketCapService
	params           *common.WalletParams
	endpoint         string
	common.Wallet
}

var XRPWALLET_RATELIMITER = common.NewRateLimiter(5, 1)

func CreateXrpWallet(params *common.WalletParams) common.Wallet {
	client := http.Client{Timeout: time.Second * 2}
	return &Ripple{
		ctx:              params.Context,
		client:           client,
		marketcapService: params.MarketCapService,
		params:           params,
		endpoint:         "https://data.ripple.com/v2"}
}

func (r *Ripple) GetPrice() decimal.Decimal {
	XRPWALLET_RATELIMITER.RespectRateLimit()
	marketcap := r.marketcapService.GetMarket("XRP")
	usd, err := decimal.NewFromString(marketcap.PriceUSD)
	if err != nil {
		r.ctx.GetLogger().Errorf("[RippleWallet.GetPrice] Error parsing price string into decimal: %s", err.Error())
		return decimal.NewFromFloat(0)
	}
	return usd.Truncate(2)
	/*
		url := fmt.Sprintf("%s/exchange_rates/%s+%s/XRP", r.endpoint,
			r.ctx.GetUser().GetLocalCurrency(), r.params.Address)
		_, body, err := util.HttpRequest(url)
		if err != nil {
			r.ctx.GetLogger().Errorf("[RippleWallet.GetPrice] Error: %s", err.Error())
		}
		exRate := RippleExchangeRate{}
		jsonErr := json.Unmarshal(body, &exRate)
		if jsonErr != nil {
			r.ctx.GetLogger().Errorf("[RippleWallet.GetPrice] %s", jsonErr.Error())
		}
		price, err := decimal.NewFromString(exRate.Rate)
		if err != nil {
			r.ctx.GetLogger().Errorf("[RippleWallet.GetPrice] Error parsing rate string into decimal: %s", err.Error())
			return decimal.NewFromFloat(0)
		}
		return price
	*/
}

func (r *Ripple) GetWallet() (common.UserCryptoWallet, error) {
	XRPWALLET_RATELIMITER.RespectRateLimit()
	r.ctx.GetLogger().Debugf("[Ripple.GetBalance] Address: %s", r.params.Address)
	balance := decimal.NewFromFloat(0)
	url := fmt.Sprintf("https://data.ripple.com/v2/accounts/%s/balances", r.params.Address)
	_, body, err := util.HttpRequest(url)
	if err != nil {
		r.ctx.GetLogger().Errorf("[Ripple.GetWallet] %s", err.Error())
		return nil, err
	}
	rb := RippleBalance{}
	jsonErr := json.Unmarshal(body, &rb)
	if jsonErr != nil {
		r.ctx.GetLogger().Errorf("[Ripple.GetWallet] %s", jsonErr.Error())
		return nil, jsonErr
	}
	if len(rb.Balances) > 0 {
		dec, err := decimal.NewFromString(rb.Balances[0].Balance)
		if err != nil {
			r.ctx.GetLogger().Errorf("[Ripple.GetWallet] Error converting balance to decimal: %s", err.Error())
		}
		balance = dec
	}
	marketcap := r.marketcapService.GetMarket("XRP")
	priceUSD, err := decimal.NewFromString(marketcap.PriceUSD)
	return &dto.UserCryptoWalletDTO{
		Address:  r.params.Address,
		Balance:  balance.Truncate(8),
		Currency: "XRP",
		Value:    priceUSD.Mul(balance).Truncate(2)}, nil
}

/*
func (r *Ripple) GetTransactions() ([]common.Transaction, error) {
	var transactions []common.Transaction
	userEntity := &entity.User{Id: r.ctx.GetUser().GetId()}
	for _, wallet := range r.userDAO.GetWallets(userEntity) {
		txs, err := r.getTransactionsFor(wallet.GetAddress())
		if err != nil {
			return transactions, err
		}
		transactions = append(transactions, txs...)
	}
	return transactions, nil
}*/

func (r *Ripple) GetTransactions() ([]common.Transaction, error) {
	XRPWALLET_RATELIMITER.RespectRateLimit()
	var transactions []common.Transaction
	var rippleResponse RippleResponse
	url := fmt.Sprintf("https://data.ripple.com/v2/accounts/%s/transactions", r.params.Address)
	_, body, err := util.HttpRequest(url)
	if err != nil {
		r.ctx.GetLogger().Errorf("[Ripple.GetBalance] %s", err.Error())
	}
	jsonErr := json.Unmarshal(body, &rippleResponse)
	if jsonErr != nil {
		r.ctx.GetLogger().Errorf("[Ripple.GetTransactions] JSON unmarshal error: %s", jsonErr.Error())
	}
	for _, tx := range rippleResponse.Transactions {
		var txType string
		if tx.Tx.Destination == r.params.Address {
			txType = common.DEPOSIT_ORDER_TYPE
		} else {
			txType = common.WITHDRAWAL_ORDER_TYPE
		}
		amount := tx.Tx.Amount
		if amount == "" {
			amount = "0.0"
		}
		amt, err := decimal.NewFromString(amount)
		if err != nil {
			return nil, err
		}
		amt = amt.Div(decimal.NewFromFloat(1000000))
		fee, err := decimal.NewFromString(tx.Tx.Fee)
		if err != nil {
			return nil, err
		}
		fee = fee.Div(decimal.NewFromFloat(1000000))
		currencyPair := &common.CurrencyPair{
			Base:          "XRP",
			Quote:         "XRP",
			LocalCurrency: r.ctx.GetUser().GetLocalCurrency()}
		candlestick, err := r.params.FiatPriceService.GetPriceAt("XRP", tx.Date)
		if err != nil {
			return nil, err
		}
		total := amt.Mul(candlestick.Close)
		transactions = append(transactions, &dto.TransactionDTO{
			Id:                   tx.Hash,
			Date:                 tx.Date,
			CurrencyPair:         currencyPair,
			Type:                 txType,
			Network:              "ripple",
			NetworkDisplayName:   "Ripple",
			Quantity:             amt.StringFixed(8),
			QuantityCurrency:     "XRP",
			FiatQuantity:         amt.Mul(candlestick.Close).StringFixed(2),
			FiatQuantityCurrency: "USD",
			Price:                candlestick.Close.StringFixed(2),
			PriceCurrency:        "USD",
			FiatPrice:            candlestick.Close.StringFixed(2),
			FiatPriceCurrency:    "USD",
			Fee:                  fee.StringFixed(8),
			FeeCurrency:          "XRP",
			FiatFee:              fee.Mul(candlestick.Close).StringFixed(2),
			FiatFeeCurrency:      "USD",
			Total:                amt.StringFixed(8),
			TotalCurrency:        "XRP",
			FiatTotal:            total.StringFixed(2),
			FiatTotalCurrency:    "USD"})
	}
	return transactions, nil
}

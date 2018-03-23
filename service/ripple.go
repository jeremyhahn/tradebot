package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/dto"
	"github.com/jeremyhahn/tradebot/entity"
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
	marketcapService MarketCapService
	BlockExplorerService
}

func NewRippleService(ctx common.Context, userDAO dao.UserDAO, marketcapService MarketCapService) BlockExplorerService {
	client := http.Client{Timeout: time.Second * 2}
	return &Ripple{
		ctx:              ctx,
		client:           client,
		userDAO:          userDAO,
		marketcapService: marketcapService}
}

func (r *Ripple) GetWallet(address string) (common.UserCryptoWallet, error) {
	r.ctx.GetLogger().Debugf("[Ripple.GetBalance] Address: %s", address)
	var balance float64
	url := fmt.Sprintf("https://data.ripple.com/v2/accounts/%s/balances", address)
	_, body, err := util.HttpRequest(url)
	if err != nil {
		r.ctx.GetLogger().Errorf("[Ripple.GetBalance] %s", err.Error())
		return nil, err
	}
	rb := RippleBalance{}
	jsonErr := json.Unmarshal(body, &rb)
	if jsonErr != nil {
		r.ctx.GetLogger().Errorf("[Ripple.GetBalance] %s", jsonErr.Error())
		return nil, jsonErr
	}
	if len(rb.Balances) <= 0 {
		balance = 0.0
	} else {
		f, _ := strconv.ParseFloat(rb.Balances[0].Balance, 64)
		balance = f
	}
	marketcap := r.marketcapService.GetMarket("XRP")
	f2, _ := strconv.ParseFloat(marketcap.PriceUSD, 64)
	return &dto.UserCryptoWalletDTO{
		Address:  address,
		Balance:  balance,
		Currency: "XRP",
		Value:    f2 * balance}, nil
}

func (r *Ripple) GetTransactions() ([]common.Transaction, error) {
	var transactions []common.Transaction
	userEntity := &entity.User{Id: r.ctx.GetUser().GetId()}
	for _, wallet := range r.userDAO.GetWallets(userEntity) {
		txs, err := r.GetTransaction(wallet.GetAddress())
		if err != nil {
			return transactions, err
		}
		transactions = append(transactions, txs...)
	}
	return transactions, nil
}

func (r *Ripple) GetTransaction(address string) ([]common.Transaction, error) {
	var transactions []common.Transaction
	var rippleResponse RippleResponse
	url := fmt.Sprintf("https://data.ripple.com/v2/accounts/%s/transactions", address)
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
		if tx.Tx.Destination == address {
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
		fee, err := decimal.NewFromString(tx.Tx.Fee)
		if err != nil {
			return nil, err
		}
		currencyPair := &common.CurrencyPair{
			Base:          "XRP",
			Quote:         "XRP",
			LocalCurrency: r.ctx.GetUser().GetLocalCurrency()}
		transactions = append(transactions, &dto.TransactionDTO{
			Id:                 tx.Hash,
			Date:               tx.Date,
			CurrencyPair:       currencyPair,
			Type:               txType,
			Network:            "ripple",
			NetworkDisplayName: "Ripple",
			Quantity:           amt.Div(decimal.NewFromFloat(1000000)).StringFixed(8),
			QuantityCurrency:   "XRP",
			Fee:                fee.Div(decimal.NewFromFloat(1000000)).StringFixed(8),
			FeeCurrency:        "XRP",
			Total:              amt.StringFixed(8),
			TotalCurrency:      "XRP"})
	}
	return transactions, nil
}

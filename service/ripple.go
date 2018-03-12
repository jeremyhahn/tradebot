package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dto"
	"github.com/jeremyhahn/tradebot/util"
)

type RippleTx struct {
	Account     string `json:"Account"`
	Destination string `json:"Destination"`
	Amount      string `json:"Amount"`
	Fee         string `json:"Fee"`
}

type RippleTransaction struct {
	Date time.Time `json:"date"`
	Tx   RippleTx  `json:"tx"`
}

type RippleResponse struct {
	Result       string              `json:"result"`
	Transactions []RippleTransaction `json:"transactions"`
}

type Ripple struct {
	ctx              common.Context
	client           http.Client
	marketcapService MarketCapService
	BlockExplorerService
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

func NewRipple(ctx common.Context, marketcapService MarketCapService) *Ripple {
	client := http.Client{Timeout: time.Second * 2}
	return &Ripple{
		ctx:              ctx,
		client:           client,
		marketcapService: marketcapService}
}

func (r *Ripple) GetBalance(address string) common.UserCryptoWallet {
	r.ctx.GetLogger().Debugf("[Ripple.GetBalance] Address: %s", address)
	var balance float64
	url := fmt.Sprintf("https://data.ripple.com/v2/accounts/%s/balances", address)
	_, body, err := util.HttpRequest(url)
	if err != nil {
		r.ctx.GetLogger().Errorf("[Ripple.GetBalance] %s", err.Error())
	}
	rb := RippleBalance{}
	jsonErr := json.Unmarshal(body, &rb)
	if jsonErr != nil {
		r.ctx.GetLogger().Errorf("[Ripple.GetBalance] %s", jsonErr.Error())
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
		Value:    f2 * balance}
}

func (r *Ripple) GetTransactions(address string) ([]common.Transaction, error) {
	url := fmt.Sprintf("https://data.ripple.com/v2/accounts/%s/transactions", address)
	_, body, err := util.HttpRequest(url)
	if err != nil {
		r.ctx.GetLogger().Errorf("[Ripple.GetBalance] %s", err.Error())
	}

	var transactions []common.Transaction
	var rippleResponse RippleResponse
	jsonErr := json.Unmarshal(body, &rippleResponse)
	if jsonErr != nil {
		r.ctx.GetLogger().Errorf("[Ripple.GetTransactions] Error: %s", jsonErr.Error())
	}

	for _, tx := range rippleResponse.Transactions {
		var txType string
		if tx.Tx.Destination == address {
			txType = "deposit"
		} else {
			txType = "withdrawl"
		}
		amt, _ := strconv.ParseFloat(tx.Tx.Amount, 64)
		transactions = append(transactions, &dto.TransactionDTO{
			Date:   tx.Date,
			Amount: amt / 1000000,
			Type:   txType})
	}

	return transactions, nil
}

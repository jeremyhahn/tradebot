package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/jeremyhahn/tradebot/common"
)

type Ripple struct {
	ctx              *common.Context
	client           http.Client
	marketcapService *MarketCapService
	common.Wallet
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

func NewRipple(ctx *common.Context, marketcapService *MarketCapService) *Ripple {
	client := http.Client{Timeout: time.Second * 2}
	return &Ripple{
		ctx:              ctx,
		client:           client,
		marketcapService: marketcapService}
}

func (r *Ripple) GetBalance(address string) *common.CryptoWallet {

	r.ctx.Logger.Debugf("[Ripple.GetBalance] Address: %s", address)

	var balance float64
	url := fmt.Sprintf("https://data.ripple.com/v2/accounts/%s/balances", address)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		r.ctx.Logger.Errorf("[Ripple.GetBalance] %s", err.Error())
	}

	req.Header.Set("User-Agent", fmt.Sprintf("%s/v%s", common.APPNAME, common.APPVERSION))

	res, getErr := r.client.Do(req)
	if getErr != nil {
		r.ctx.Logger.Errorf("[Ripple.GetBalance] %s", getErr.Error())
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		r.ctx.Logger.Errorf("[Ripple.GetBalance] %s", readErr.Error())
	}

	rb := RippleBalance{}
	jsonErr := json.Unmarshal(body, &rb)
	if jsonErr != nil {
		r.ctx.Logger.Errorf("[Ripple.GetBalance] %s", jsonErr.Error())
	}
	if len(rb.Balances) <= 0 {
		balance = 0.0
	} else {
		f, _ := strconv.ParseFloat(rb.Balances[0].Balance, 64)
		balance = f
	}

	marketcap := r.marketcapService.GetMarket("XRP")
	f2, _ := strconv.ParseFloat(marketcap.PriceUSD, 64)

	return &common.CryptoWallet{
		Address:  address,
		Balance:  balance,
		Currency: "XRP",
		NetWorth: f2 * balance}
}

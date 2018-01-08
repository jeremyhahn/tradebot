package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/jeremyhahn/tradebot/common"
	logging "github.com/op/go-logging"
)

type BlockchainWallet struct {
	Address  string  `json:"address"`
	Balance  float64 `json:"final_balance"`
	Currency string  `json:"currency"`
}

type BlockchainTickerSubItem struct {
	Last float64 `json:"last"`
}

type BlockchainTickerItem struct {
	USD BlockchainTickerSubItem
}

type BlockchainInfo struct {
	logger     *logging.Logger
	client     http.Client
	items      BlockchainTickerItem
	lastPrice  float64
	lastLookup time.Time
	common.Wallet
}

func NewBlockchainInfo(ctx *common.Context) *BlockchainInfo {
	client := http.Client{Timeout: common.HTTP_CLIENT_TIMEOUT}
	return &BlockchainInfo{
		logger:     ctx.Logger,
		client:     client,
		lastPrice:  0.0,
		lastLookup: time.Now().Add(-20 * time.Minute)}
}

func (b *BlockchainInfo) GetPrice() float64 {

	url := "https://blockchain.info/ticker"
	elapsed := float64(time.Since(b.lastLookup))

	if elapsed/float64(time.Second) >= 900 {
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			b.logger.Errorf("[BlockchainInfo.GetPrice] %s", err.Error())
		}

		req.Header.Set("User-Agent", fmt.Sprintf("%s/v%s", common.APPNAME, common.APPVERSION))

		res, getErr := b.client.Do(req)
		if getErr != nil {
			b.logger.Errorf("[BlockchainInfo.GetPrice] %s", getErr.Error())
		}

		body, readErr := ioutil.ReadAll(res.Body)
		if readErr != nil {
			b.logger.Errorf("[BlockchainInfo.GetPrice] %s", readErr.Error())
		}

		t := BlockchainTickerItem{}
		jsonErr := json.Unmarshal(body, &t)
		if jsonErr != nil {
			b.logger.Errorf("[BlockchainInfo.GetPrice] %s", jsonErr.Error())
		}

		b.lastLookup = time.Now()
		b.lastPrice = t.USD.Last
	}
	return b.lastPrice
}

func (b *BlockchainInfo) GetBalance(address string) *common.CryptoWallet {

	url := fmt.Sprintf("https://blockchain.info/address/%s?format=json", address)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		b.logger.Errorf("[BlockchainInfo.GetBalance] %s", err.Error())
	}

	req.Header.Set("User-Agent", fmt.Sprintf("%s/v%s", common.APPNAME, common.APPVERSION))

	res, getErr := b.client.Do(req)
	if getErr != nil {
		b.logger.Errorf("[BlockchainInfo.GetBalance] %s", getErr.Error())
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		b.logger.Errorf("[BlockchainInfo.GetBalance] %s", readErr.Error())
	}

	wallet := BlockchainWallet{}
	jsonErr := json.Unmarshal(body, &wallet)
	if jsonErr != nil {
		b.logger.Errorf("[BlockchainInfo.GetBalance] %s", jsonErr.Error())
	}
	return &common.CryptoWallet{
		Address:  address,
		Balance:  wallet.Balance / 100000000,
		Currency: "BTC"}
}

func (b *BlockchainInfo) ConvertToUSD(currency string, btc float64) float64 {
	price := b.GetPrice()
	b.logger.Debugf("[BlockchainTicker] currency: %s, btc: %.8f, price: %f", currency, btc, price)
	return btc * price
}

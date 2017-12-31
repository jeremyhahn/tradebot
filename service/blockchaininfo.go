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

type BlockchainTickerSubItem struct {
	Last float64 `json:"last"`
}

type BlockchainTickerItem struct {
	USD BlockchainTickerSubItem
}

type BlockchainInfo struct {
	logger     *logging.Logger
	url        string
	client     http.Client
	items      BlockchainTickerItem
	lastPrice  float64
	lastLookup time.Time
}

func NewBlockchainInfo(logger *logging.Logger) *BlockchainInfo {
	url := "https://blockchain.info/ticker"
	client := http.Client{
		Timeout: time.Second * 2}
	return &BlockchainInfo{
		logger:     logger,
		url:        url,
		client:     client,
		lastPrice:  0.0,
		lastLookup: time.Now().Add(-20 * time.Minute)}
}

func (ticker *BlockchainInfo) GetPrice() float64 {

	elapsed := float64(time.Since(ticker.lastLookup))
	if elapsed/float64(time.Second) >= 900 {

		req, err := http.NewRequest(http.MethodGet, ticker.url, nil)
		if err != nil {
			ticker.logger.Fatal(err)
		}

		req.Header.Set("User-Agent", fmt.Sprintf("%s/v%s", common.APPNAME, common.APPVERSION))

		res, getErr := ticker.client.Do(req)
		if getErr != nil {
			ticker.logger.Fatal(getErr)
		}

		body, readErr := ioutil.ReadAll(res.Body)
		if readErr != nil {
			ticker.logger.Fatal(readErr)
		}

		t := BlockchainTickerItem{}
		jsonErr := json.Unmarshal(body, &t)
		if jsonErr != nil {
			ticker.logger.Fatal(jsonErr)
		}

		ticker.lastLookup = time.Now()
		ticker.lastPrice = t.USD.Last
	}
	return ticker.lastPrice
}

func (ticker *BlockchainInfo) ConvertToUSD(currency string, btc float64) float64 {
	price := ticker.GetPrice()
	ticker.logger.Debugf("[BlockchainTicker] currency: %s, btc: %.8f, price: %f", currency, btc, price)
	return btc * price
}

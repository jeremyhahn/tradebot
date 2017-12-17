package main

import (
	"strconv"
	"time"

	ws "github.com/gorilla/websocket"
	gdax "github.com/preichenberger/go-gdax"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/op/go-logging"
)

type Coinbase struct {
	gdax     *gdax.Client
	logger   *logging.Logger
	Price    float64
	currency string
	common.Exchange
}

func NewCoinbase(config IConfiguration, logger *logging.Logger, currency string) *Coinbase {
	apiKey := config.Get("api_key")
	apiSecret := config.Get("api_secret")
	apiPassphrase := config.Get("api_passphrase")
	return &Coinbase{
		gdax:     gdax.NewClient(apiSecret, apiKey, apiPassphrase),
		logger:   logger,
		currency: currency}
}

func (cb *Coinbase) GetCurrency() string {
	return cb.currency
}

func (cb *Coinbase) GetTradeHistory(start, end time.Time, granularity int) []common.Candlestick {
	cb.logger.Info("Getting Coinbase trade history")
	/*
		var products []gdax.Product
		products, err := cb.gdax.GetProducts()
		if err != nil {
			cb.logger.Error(err)
		}
		for _, p := range products {
			fmt.Printf("%+v\n", p)
		}
	*/
	var candlesticks []common.Candlestick
	params := gdax.GetHistoricRatesParams{
		Start:       start,
		End:         end,
		Granularity: granularity}
	rates, err := cb.gdax.GetHistoricRates(cb.currency, params)
	if err != nil {
		cb.logger.Error(err)
		time.Sleep(time.Duration(time.Second * 10))
		return cb.GetTradeHistory(start, end, granularity)
	}
	for _, r := range rates {
		candlesticks = append(candlesticks, common.Candlestick{
			Date:   r.Time,
			Open:   r.Open,
			Close:  r.Close,
			High:   r.High,
			Low:    r.Low,
			Volume: r.Volume})
	}
	return candlesticks
}

func (cb *Coinbase) GetAccounts() []common.Account {
	var accts []common.Account
	accounts, err := cb.gdax.GetAccounts()
	if err != nil {
		cb.logger.Error(err.Error())
	}
	for _, a := range accounts {
		currency, err := strconv.ParseFloat(a.Currency, 64)
		if err != nil {
			cb.logger.Error(err.Error())
		}
		accts = append(accts, common.Account{
			Currency: currency,
			Balance:  a.Balance})
	}
	return accts
}

func (cb *Coinbase) GetPrice() float64 {
	return cb.Price
}

func (cb *Coinbase) SubscribeToLiveFeed(price chan float64) {

	cb.logger.Info("Subscribing to Coinbase WebSocket feed")

	var wsDialer ws.Dialer
	wsConn, _, err := wsDialer.Dial("wss://ws-feed.gdax.com", nil)
	if err != nil {
		println(err.Error())
	}

	subscribe := map[string]string{
		"type":       "subscribe",
		"product_id": cb.currency,
	}

	if err := wsConn.WriteJSON(subscribe); err != nil {
		cb.logger.Error(err.Error())
	}

	message := gdax.Message{}
	for true {

		if err := wsConn.ReadJSON(&message); err != nil {
			cb.logger.Error(err.Error())
			break
		}

		if message.Type == "match" && message.Reason == "filled" {
			if message.Price != cb.Price {
				cb.Price = message.Price
				price <- message.Price
			}
		}
	}

	cb.SubscribeToLiveFeed(price)
}

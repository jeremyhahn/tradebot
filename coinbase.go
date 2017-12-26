package main

import (
	"fmt"
	"strconv"
	"time"

	ws "github.com/gorilla/websocket"
	gdax "github.com/preichenberger/go-gdax"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/op/go-logging"
)

type Coinbase struct {
	gdax        *gdax.Client
	logger      *logging.Logger
	PriceStream *PriceStream
	common.Exchange
}

func NewCoinbase(coinbase *CoinExchange, logger *logging.Logger, priceStream *PriceStream) *Coinbase {
	return &Coinbase{
		gdax:        gdax.NewClient(coinbase.Secret, coinbase.Key, coinbase.Passphrase),
		logger:      logger,
		PriceStream: priceStream}
}

func (cb *Coinbase) GetTradeHistory(currency string, start, end time.Time, granularity int) []common.Candlestick {
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
	rates, err := cb.gdax.GetHistoricRates(currency, params)
	if err != nil {
		cb.logger.Error(err)
		time.Sleep(time.Duration(time.Second * 10))
		return cb.GetTradeHistory(currency, start, end, granularity)
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

func (cb *Coinbase) GetBalances() []common.Coin {
	var coins []common.Coin
	accounts, err := cb.gdax.GetAccounts()
	if err != nil {
		cb.logger.Error(err.Error())
	}
	for _, a := range accounts {
		ticker, err := cb.gdax.GetTicker(fmt.Sprintf("%s-USD", a.Currency))
		if err != nil {
			cb.logger.Error(err)
		}
		if a.Balance <= 0 {
			continue
		}
		total := a.Balance * ticker.Price
		t, err := strconv.ParseFloat(fmt.Sprintf("%.2f", total), 64)
		if err != nil {
			cb.logger.Error(err)
		}
		coins = append(coins, common.Coin{
			Currency:  a.Currency,
			Balance:   a.Balance,
			Available: a.Available,
			Pending:   a.Hold,
			Price:     ticker.Price,
			Total:     t})
	}
	return coins
}

func (cb *Coinbase) SubscribeToLiveFeed(currency string, priceChannel chan common.PriceChange) {

	cb.logger.Info("Subscribing to Coinbase WebSocket feed")

	var wsDialer ws.Dialer
	wsConn, _, err := wsDialer.Dial("wss://ws-feed.gdax.com", nil)
	if err != nil {
		println(err.Error())
	}

	subscribe := map[string]string{
		"type":       "subscribe",
		"product_id": currency,
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
			cb.PriceStream.Add(message.Price)
			/*
				priceChannel <- common.PriceChannel{
					Currency: cb.currency,
					Satoshis: message.Price,
					Price:    util.RoundFloat(cb.GetPriceUSD(), 2)}*/
		}
	}

	cb.SubscribeToLiveFeed(currency, priceChannel)
}

func (cb *Coinbase) ToUSD(price, satoshis float64) float64 {
	return satoshis * price
}

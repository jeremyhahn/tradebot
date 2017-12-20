package main

import (
	"time"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/util"
	logging "github.com/op/go-logging"
	bittrex "github.com/toorop/go-bittrex"
)

type Bittrex struct {
	client   *bittrex.Bittrex
	logger   *logging.Logger
	price    float64
	satoshis float64
	currency string
	ticker   *BlockchainTicker
	common.Exchange
}

func NewBittrex(config IConfiguration, logger *logging.Logger,
	currency string, ticker *BlockchainTicker) *Bittrex {

	apiKey := config.Get("bittrex_api_key")
	apiSecret := config.Get("bittrex_api_secret")
	return &Bittrex{
		client:   bittrex.New(apiKey, apiSecret),
		logger:   logger,
		currency: currency,
		ticker:   ticker}
}

func (b *Bittrex) SubscribeToLiveFeed(price chan common.PriceChannel) {
	for {
		time.Sleep(10 * time.Second)
		ticker, err := b.client.GetTicker(b.currency)
		if err != nil {
			b.logger.Error(err)
			continue
		}
		f, _ := ticker.Last.Float64()
		if f <= 0 {
			b.logger.Errorf("Unable to get ticker data for %s", b.currency)
			continue
		}
		b.logger.Debugf("[Bittrex] Sending live price: %.8f", f)
		b.satoshis = f
		price <- common.PriceChannel{
			Currency: b.currency,
			Satoshis: b.satoshis,
			Price:    util.RoundFloat(b.ticker.ConvertToUSD(b.currency, b.satoshis), 2)}
	}
}

func (b *Bittrex) GetPrice() float64 {
	return b.price
}

func (b *Bittrex) GetSatoshis() float64 {
	return b.satoshis
}

func (b *Bittrex) GetTradeHistory(start, end time.Time, granularity int) []common.Candlestick {
	b.logger.Debug("Getting Bittrex trade history")
	candlesticks := make([]common.Candlestick, 0)
	/*
		marketSummary, err := b.client.GetMarketSummary(b.currency)
		if err != nil {
			b.logger.Error(err)
		}
		for _, s := range marketSummary {
			fmt.Printf("%+v\n", s)
			f, _ := s.Last.Float64()
			candlesticks = append(candlesticks, common.Candlestick{Close: f})
		}
	*/
	marketHistory, err := b.client.GetMarketHistory(b.currency)
	if err != nil {
		b.logger.Error(err)
	}
	for _, m := range marketHistory {
		f, _ := m.Price.Float64()
		if err != nil {
			b.logger.Error(err)
		}
		candlesticks = append(candlesticks, common.Candlestick{Close: f})
	}
	return candlesticks
}

func (b *Bittrex) GetCurrency() string {
	return b.currency
}

func (b *Bittrex) ConvertToUSD(btc float64) float64 {
	return b.ticker.ConvertToUSD(b.currency, b.satoshis)
}

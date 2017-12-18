package main

import (
	"fmt"
	"os"
	"time"

	"github.com/btcsuite/btcutil"
	"github.com/jeremyhahn/tradebot/common"
	logging "github.com/op/go-logging"
	bittrex "github.com/toorop/go-bittrex"
)

type BittrexPrice struct {
	Price string `json:"price"`
}

type Bittrex struct {
	client   *bittrex.Bittrex
	logger   *logging.Logger
	Price    float64
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

func (b *Bittrex) SubscribeToLiveFeed(price chan float64) {
	for {
		time.Sleep(10 * time.Second)
		ticker, err := b.client.GetTicker(b.currency)
		if err != nil {
			b.logger.Error(err)
			continue
		}

		f, _ := ticker.Last.Float64()

		//fmt.Printf("f: %f", f)
		os.Exit(1)

		if f <= 0 {
			b.logger.Errorf("Unable to get ticker data for %s", b.currency)
			continue
		}

		b.ticker.logger.Infof("Subscribe to feed, f: %f", f)

		amount, err := btcutil.NewAmount(f)
		if err != nil {
			b.logger.Error(err)
			continue
		}

		b.Price = amount.ToBTC()

		price <- b.Price
	}
}

func (b *Bittrex) GetPrice() float64 {
	return b.Price
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

		candlesticks = append(candlesticks, common.Candlestick{Close: f})
	}

	for _, c := range candlesticks {
		fmt.Printf("%+v\n", c)
	}

	return candlesticks
}

func (b *Bittrex) GetCurrency() string {
	return b.currency
}

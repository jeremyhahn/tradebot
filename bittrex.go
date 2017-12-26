package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/jeremyhahn/tradebot/common"
	logging "github.com/op/go-logging"
	bittrex "github.com/toorop/go-bittrex"
)

type Bittrex struct {
	client      *bittrex.Bittrex
	logger      *logging.Logger
	price       float64
	satoshis    float64
	PriceStream *PriceStream
	common.Exchange
}

func NewBittrex(btx *CoinExchange, logger *logging.Logger, priceStream *PriceStream) *Bittrex {
	return &Bittrex{
		client:      bittrex.New(btx.Key, btx.Secret),
		logger:      logger,
		PriceStream: priceStream}
}

func (b *Bittrex) SubscribeToLiveFeed(currency string, price chan common.PriceChange) {
	for {
		time.Sleep(10 * time.Second)
		ticker, err := b.client.GetTicker(currency)
		if err != nil {
			b.logger.Error(err)
			continue
		}
		f, _ := ticker.Last.Float64()
		if f <= 0 {
			b.logger.Errorf("Unable to get ticker data for %s", currency)
			continue
		}
		b.logger.Debugf("[Bittrex] Sending live price: %.8f", f)
		b.satoshis = f
		b.PriceStream.Add(f)
		/*
			price <- common.PriceChange{
				Currency: currency,
				Satoshis: b.satoshis,
				Price:    util.RoundFloat(b.ticker.ConvertToUSD(b.currency, b.satoshis), 2)}*/
	}
}

func (b *Bittrex) GetTradeHistory(currency string, start, end time.Time, granularity int) []common.Candlestick {
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
	marketHistory, err := b.client.GetMarketHistory(currency)
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

func (b *Bittrex) GetBalances() []common.Coin {
	var coins []common.Coin
	balances, err := b.client.GetBalances()
	if err != nil {
		b.logger.Error(err)
	}
	for _, bal := range balances {
		b.logger.Debugf("[Bittrex.GetBalances] Getting currency ticker: %s", bal.Currency)
		ticker, terr := b.client.GetTicker(fmt.Sprintf("BTC-%s", bal.Currency))
		if terr != nil {
			b.logger.Error(terr)
		}
		price, exact := ticker.Last.Float64()
		if !exact {
			b.logger.Error("[Bittrex.GetBalances] Conversion of Ticker price to Float64 was not exact!")
		}
		avail, exact := bal.Available.Float64()
		if !exact {
			b.logger.Error("[Bittrex.GetBalances] Conversion of Available funds to Float64 was not exact!")
		}
		balance, exact := bal.Balance.Float64()
		if !exact {
			b.logger.Error("[Bittrex.GetBalances] Conversion of Balance to Float64 was not exact!")
		}
		pending, exact := bal.Pending.Float64()
		if !exact {
			b.logger.Error("[Bittrex.GetBalances] Conversion of Pending to Float64 was not exact!")
		}
		if balance <= 0 {
			continue
		}
		total := balance * (price * b.getBitcoinPrice())
		t, err := strconv.ParseFloat(fmt.Sprintf("%.2f", total), 64)
		if err != nil {
			b.logger.Error(err)
		}
		coins = append(coins, common.Coin{
			Address:   bal.CryptoAddress,
			Available: avail,
			Balance:   balance,
			Currency:  bal.Currency,
			Pending:   pending,
			Price:     price, // BTC satoshis, not actual USD price
			Total:     t})
	}
	return coins
}

func (b *Bittrex) getBitcoinPrice() float64 {
	summary, err := b.client.GetMarketSummary("USDT-BTC")
	if err != nil {
		b.logger.Error(err)
	}
	if len(summary) != 1 {
		b.logger.Error("[Bittrex.getBitcoinPrice] Unexpected number of BTC markets returned")
		return 0
	}
	price, exact := summary[0].Last.Float64()
	if !exact {
		b.logger.Error("[Bittrex.getBitcoinPrice] Conversion of BTC price to Float64 was not exact!")
	}
	return price
}

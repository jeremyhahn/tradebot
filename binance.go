package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/adshao/go-binance"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/op/go-logging"
)

type Binance struct {
	api         *binance.Client
	logger      *logging.Logger
	PriceStream *PriceStream
	common.Exchange
}

func NewBinance(config *CoinExchange, logger *logging.Logger, priceStream *PriceStream) *Binance {
	return &Binance{
		api:         binance.NewClient(config.Key, config.Secret),
		logger:      logger,
		PriceStream: priceStream}
}

func (b *Binance) GetBalances() []common.Coin {

	var coins []common.Coin

	account, err := b.api.NewGetAccountService().Do(context.Background())
	if err != nil {
		b.logger.Error(err)
	}

	for _, balance := range account.Balances {

		bal, err := strconv.ParseFloat(balance.Free, 64)
		if err != nil {
			b.logger.Error(err)
		}

		if bal > 0 {

			bitcoin := b.getBitcoin()
			bitcoinPrice := b.parseBitcoinPrice(bitcoin)

			if balance.Asset == "BTC" {
				total := bal * bitcoinPrice
				t, err := strconv.ParseFloat(strconv.FormatFloat(total, 'f', 2, 64), 64)
				if err != nil {
					b.logger.Error(err)
				}
				coins = append(coins, common.Coin{
					Currency:  balance.Asset,
					Available: bal,
					Balance:   bal,
					Price:     bitcoinPrice,
					Total:     t})
				continue
			}

			symbol := balance.Asset + "BTC"
			ticker, err := b.api.NewPriceChangeStatsService().Symbol(symbol).Do(context.Background())
			if err != nil {
				b.logger.Error(err)
			}
			if ticker == nil {
				b.logger.Errorf("Unable to retrieve ticker for symbol: %s", symbol)
				continue
			}

			fmt.Printf("%+v\n", ticker)

			lastPrice, err := strconv.ParseFloat(ticker.LastPrice, 64)
			if err != nil {
				b.logger.Error(err)
			}

			total := bal * (lastPrice * bitcoinPrice)

			t, err := strconv.ParseFloat(fmt.Sprintf("%.2f", total), 64)
			if err != nil {
				b.logger.Error(err)
			}

			coins = append(coins, common.Coin{
				Currency:  balance.Asset,
				Available: bal,
				Balance:   bal,
				Price:     lastPrice,
				Total:     t})
		}
	}

	//	fmt.Printf("%+v\n", coins)
	//	os.Exit(10)

	return coins
}

func (b *Binance) getBitcoin() *binance.PriceChangeStats {
	stats, err := b.api.NewPriceChangeStatsService().Symbol("BTCUSDT").Do(context.Background())
	if err != nil {
		fmt.Println(err)
	}
	return stats
}

func (b *Binance) getBitcoinPrice() float64 {
	bitcoin := b.getBitcoin()
	f, err := strconv.ParseFloat(bitcoin.LastPrice, 8)
	if err != nil {
		b.logger.Error(err)
	}
	return f
}

func (b *Binance) parseBitcoinPrice(bitcoin *binance.PriceChangeStats) float64 {
	f, err := strconv.ParseFloat(bitcoin.LastPrice, 8)
	if err != nil {
		b.logger.Error(err)
	}
	return f
}

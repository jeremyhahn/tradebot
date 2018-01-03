package service

import (
	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/exchange"
)

func NewPortfolio(ctx *common.Context, exchanges map[string]common.Exchange,
	walletList []common.CryptoWallet) *common.Portfolio {
	var exchangeList []common.CoinExchange
	var netWorth float64
	for _, ex := range exchange.NewExchangeDAO(ctx.DB, ctx.Logger).Exchanges {
		exchange := exchanges[ex.Name]
		if exchange == nil {
			ctx.Logger.Errorf("[Portfolio] Exchange is null. Name: %s", ex.Name)
			continue
		}
		ex, net := exchange.GetExchange()
		netWorth += net
		exchangeList = append(exchangeList, ex)
	}
	for _, wallet := range walletList {
		netWorth += wallet.NetWorth
	}
	return &common.Portfolio{
		User:      ctx.User,
		NetWorth:  netWorth,
		Exchanges: exchangeList,
		Wallets:   walletList}
}

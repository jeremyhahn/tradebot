package service

import (
	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/exchange"
	"github.com/jinzhu/gorm"
	logging "github.com/op/go-logging"
)

func NewPortfolio(db *gorm.DB, logger *logging.Logger, user *common.User, exchanges map[string]common.Exchange) *common.Portfolio {
	var exchangeList []common.CoinExchange
	var netWorth float64
	for _, ex := range exchange.NewExchangeDAO(db, logger).Exchanges {
		exchange := exchanges[ex.Name]
		if exchange == nil {
			logger.Errorf("[Portfolio] Exchange is null. Name: %s", ex.Name)
			continue
		}
		ex, net := exchange.GetExchange()
		netWorth += net
		exchangeList = append(exchangeList, ex)
	}
	return &common.Portfolio{
		User:      user,
		NetWorth:  netWorth,
		Exchanges: exchangeList}
}

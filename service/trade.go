package service

import (
	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
)

type TradeLedger struct {
	Symbol   string
	BuyPrice float64
}

type TradeService struct {
	ctx              *common.Context
	marketcapService *MarketCapService
	currencyPair     *common.CurrencyPair
	marketMap        map[string]common.MarketCap
}

func NewTradeService(ctx *common.Context, marketcapService *MarketCapService) *TradeService {
	return &TradeService{
		ctx:              ctx,
		marketcapService: marketcapService}
}

func (ts *TradeService) MakeMeRich(currencyPair *common.CurrencyPair) {
	var charts []Chart
	userDAO := dao.NewUserDAO(ts.ctx)
	userService := NewUserService(ts.ctx, userDAO, ts.marketcapService)
	exchanges := userService.GetExchanges(ts.ctx.User, currencyPair)
	for _, ex := range exchanges {
		exchange := userService.GetExchange(ts.ctx.User, ex.Name, currencyPair)
		chart := NewChart(ts.ctx, exchange, 900) // 15 minutes
		charts = append(charts, *chart)
	}
	charts[0].Stream()
}

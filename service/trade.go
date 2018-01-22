package service

/*
import (
	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
)

type TradeService struct {
	ctx              *common.Context
	marketcapService *MarketCapService
	currencyPair     *common.CurrencyPair
	marketMap        map[string]common.MarketCap
	Charts           []common.ChartService
}

func NewTradeService(ctx *common.Context, marketcapService *MarketCapService) *TradeService {
	var services []common.ChartService
	exchangeDAO := dao.NewExchangeDAO(ctx)
	autotradeDAO := dao.NewAutoTradeDAO(ctx)
	for _, autoTradeCoin := range autotradeDAO.Find(ctx.User) {
		currencyPair := &common.CurrencyPair{
			Base:          autoTradeCoin.GetBase(),
			Quote:         autoTradeCoin.GetQuote(),
			LocalCurrency: ctx.User.LocalCurrency}
		exchangeService := NewExchangeService(ctx, exchangeDAO)
		exchange := exchangeService.NewExchange(ctx.User, autoTradeCoin.GetExchange(), currencyPair)
		chart := NewChartService(ctx, exchange, nil, autoTradeCoin.GetPeriod())
		ctx.Logger.Debugf("[NewTradeService] Loading AutoTrade currency pair: %s-%s\n", autoTradeCoin.GetBase(), autoTradeCoin.GetQuote())
		ctx.Logger.Debugf("[NewTradeService] Chart: %+v\n", chart)
		services = append(services, chart)
	}
	return &TradeService{
		ctx:              ctx,
		marketcapService: marketcapService,
		Charts:           services}
}

func (ts *TradeService) Trade() {
	for _, chart := range ts.Charts {
		chart.Stream()
	}
}
*/

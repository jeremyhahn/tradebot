package service

import (
	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
)

type TradeServiceImpl struct {
	ctx      *common.Context
	tradeDAO dao.TradeDAO
	TradeService
}

func NewTradeService(ctx *common.Context, tradeDAO dao.TradeDAO) TradeService {
	return &TradeServiceImpl{
		ctx:      ctx,
		tradeDAO: tradeDAO}
}

/*
func (ts *TradeServiceImpl) Trade(chartBL businesslogic.ChartBL) {
	ts.charts = append(ts.charts, chartBL)

}*/

func (ts *TradeServiceImpl) Save(trade *common.Trade) {
	ts.tradeDAO.Create(&dao.Trade{
		ID:        trade.ID,
		ChartID:   trade.ChartID,
		Date:      trade.Date,
		Exchange:  trade.Exchange,
		Base:      trade.Base,
		Quote:     trade.Quote,
		Type:      trade.Type,
		Price:     trade.Price,
		Amount:    trade.Amount,
		ChartData: trade.ChartData})
}

func (ts *TradeServiceImpl) GetLastTrade(chart *common.Chart) *common.Trade {
	/*
		var trades []dao.Trade
		var indicators []dao.Indicator
		for _, trade := range chart.Trades {
			trades = append(trades, dao.Trade{
				ID:        trade.ID,
				UserID:    trade.UserID,
				ChartID:   chart.ID,
				Date:      trade.Date,
				Exchange:  trade.Exchange,
				Type:      trade.Type,
				Base:      trade.Base,
				Quote:     trade.Quote,
				Amount:    trade.Amount,
				Price:     trade.Price,
				ChartData: trade.ChartData})
		}
		for _, indicator := range chart.Indicators {
			indicators = append(indicators, dao.Indicator{
				Id:         indicator.Id,
				ChartID:    indicator.ChartID,
				Name:       indicator.Name,
				Parameters: indicator.Parameters})
		}
		daoChart := &dao.Chart{
			ID:         chart.ID,
			Base:       chart.Base,
			Exchange:   chart.Exchange,
			Period:     chart.Period,
			Trades:     trades,
			Indicators: indicators}
	*/
	daoChart := &dao.Chart{ID: chart.ID}
	entity := ts.tradeDAO.GetLastTrade(daoChart)
	return &common.Trade{
		ID:        entity.ID,
		UserID:    ts.ctx.User.Id,
		ChartID:   entity.ChartID,
		Date:      entity.Date,
		Exchange:  entity.Exchange,
		Type:      entity.Type,
		Base:      entity.Base,
		Quote:     entity.Quote,
		Amount:    entity.Amount,
		Price:     entity.Price,
		ChartData: entity.ChartData}
}

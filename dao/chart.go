package dao

import (
	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/entity"
)

type ChartDAO interface {
	Create(chart entity.ChartEntity) error
	Save(chart entity.ChartEntity) error
	Update(chart entity.ChartEntity) error
	Find(user common.UserContext, autoTradeOnly bool) ([]entity.Chart, error)
	Get(id uint) (entity.ChartEntity, error)
	GetIndicators(chart entity.ChartEntity) ([]entity.ChartIndicator, error)
	GetStrategies(chart entity.ChartEntity) ([]entity.ChartStrategy, error)
	GetTrades(user common.UserContext) ([]entity.Trade, error)
	GetLastTrade(chart entity.ChartEntity) (entity.TradeEntity, error)
}

type ChartDAOImpl struct {
	ctx common.Context
	ChartDAO
}

func NewChartDAO(ctx common.Context) ChartDAO {
	return &ChartDAOImpl{ctx: ctx}
}

func (chartDAO *ChartDAOImpl) Create(chart entity.ChartEntity) error {
	return chartDAO.ctx.GetCoreDB().Create(chart).Error
}

func (chartDAO *ChartDAOImpl) Save(chart entity.ChartEntity) error {
	return chartDAO.ctx.GetCoreDB().Save(chart).Error
}

func (chartDAO *ChartDAOImpl) Update(chart entity.ChartEntity) error {
	return chartDAO.ctx.GetCoreDB().Update(chart).Error
}

func (chartDAO *ChartDAOImpl) Get(id uint) (entity.ChartEntity, error) {
	chart := entity.Chart{}
	if err := chartDAO.ctx.GetCoreDB().First(&chart, id).Error; err != nil {
		return nil, err
	}
	return &chart, nil
}

func (chartDAO *ChartDAOImpl) Find(user common.UserContext, autoTradeonly bool) ([]entity.Chart, error) {
	var charts []entity.Chart
	daoUser := &entity.User{Id: user.GetId()}
	var err error
	if autoTradeonly {
		err = chartDAO.ctx.GetCoreDB().Where("auto_trade = ?", 1).Related(&charts).Error
	} else {
		err = chartDAO.ctx.GetCoreDB().Model(daoUser).Related(&charts).Error
	}
	if err != nil {
		return charts, err
	}
	for i, chart := range charts {
		var trades []entity.Trade
		var indicators []entity.ChartIndicator
		if err := chartDAO.ctx.GetCoreDB().Model(&chart).Related(&trades).Error; err != nil {
			return charts, err
		}
		if err := chartDAO.ctx.GetCoreDB().Model(&chart).Related(&indicators).Error; err != nil {
			return charts, err
		}
		charts[i].Indicators = indicators
		charts[i].Trades = trades
	}
	return charts, nil
}

func (chartDAO *ChartDAOImpl) GetIndicators(chart entity.ChartEntity) ([]entity.ChartIndicator, error) {
	var indicators []entity.ChartIndicator
	if err := chartDAO.ctx.GetCoreDB().Order("id asc").Model(chart).Related(&indicators).Error; err != nil {
		return nil, err
	}
	return indicators, nil
}

func (chartDAO *ChartDAOImpl) GetStrategies(chart entity.ChartEntity) ([]entity.ChartStrategy, error) {
	var strategies []entity.ChartStrategy
	if err := chartDAO.ctx.GetCoreDB().Order("id asc").Model(chart).Related(&strategies).Error; err != nil {
		return nil, err
	}
	return strategies, nil
}

func (chartDAO *ChartDAOImpl) GetTrades(user common.UserContext) ([]entity.Trade, error) {
	var trades []entity.Trade
	daoUser := &entity.User{Id: user.GetId(), Username: user.GetUsername()}
	if err := chartDAO.ctx.GetCoreDB().Order("id asc").Model(daoUser).Related(&trades).Error; err != nil {
		return nil, err
	}
	return trades, nil
}

func (chartDAO *ChartDAOImpl) GetLastTrade(chart entity.ChartEntity) (entity.TradeEntity, error) {
	var trades []entity.Trade
	if err := chartDAO.ctx.GetCoreDB().Order("date desc").Limit(1).Model(chart).Related(&trades).Error; err != nil {
		chartDAO.ctx.GetLogger().Errorf("[ChartDAOImpl.GetLastTrade] Error: %s", err.Error())
	}
	tradeLen := len(trades)
	if tradeLen < 1 || tradeLen > 1 {
		return nil, nil
	}
	return &trades[0], nil
}

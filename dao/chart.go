package dao

import (
	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/entity"
)

type ChartDAO interface {
	Create(chart entity.ChartEntity) error
	Save(chart entity.ChartEntity) error
	Update(chart entity.ChartEntity) error
	Find(user common.User) ([]entity.Chart, error)
	Get(id uint) (entity.ChartEntity, error)
	GetIndicators(chart entity.ChartEntity) ([]entity.ChartIndicator, error)
	GetStrategies(chart entity.ChartEntity) ([]entity.ChartStrategy, error)
	GetTrades(user common.User) ([]entity.Trade, error)
	GetLastTrade(chart entity.ChartEntity) (entity.TradeEntity, error)
}

type ChartDAOImpl struct {
	ctx *common.Context
	//Coins []Chart
	ChartDAO
}

func NewChartDAO(ctx *common.Context) ChartDAO {
	ctx.DB.AutoMigrate(&entity.Chart{})
	ctx.DB.AutoMigrate(&entity.ChartStrategy{})
	ctx.DB.AutoMigrate(&entity.ChartIndicator{})
	ctx.DB.AutoMigrate(&entity.Trade{})
	return &ChartDAOImpl{ctx: ctx}
}

func (chartDAO *ChartDAOImpl) Create(chart entity.ChartEntity) error {
	return chartDAO.ctx.DB.Create(chart).Error
}

func (chartDAO *ChartDAOImpl) Save(chart entity.ChartEntity) error {
	return chartDAO.ctx.DB.Save(chart).Error
}

func (chartDAO *ChartDAOImpl) Update(chart entity.ChartEntity) error {
	return chartDAO.ctx.DB.Update(chart).Error
}

func (chartDAO *ChartDAOImpl) Get(id uint) (entity.ChartEntity, error) {
	chart := entity.Chart{}
	if err := chartDAO.ctx.DB.First(&chart, id).Error; err != nil {
		return nil, err
	}
	return &chart, nil
}

func (chartDAO *ChartDAOImpl) Find(user common.User) ([]entity.Chart, error) {
	var charts []entity.Chart
	daoUser := &entity.User{Id: user.GetId()}
	if err := chartDAO.ctx.DB.Model(daoUser).Related(&charts).Error; err != nil {
		return charts, err
	}
	for i, chart := range charts {
		var trades []entity.Trade
		var indicators []entity.ChartIndicator
		if err := chartDAO.ctx.DB.Model(&chart).Related(&trades).Error; err != nil {
			return charts, err
		}
		if err := chartDAO.ctx.DB.Model(&chart).Related(&indicators).Error; err != nil {
			return charts, err
		}
		charts[i].Indicators = indicators
		charts[i].Trades = trades
	}
	return charts, nil
}

func (chartDAO *ChartDAOImpl) GetIndicators(chart entity.ChartEntity) ([]entity.ChartIndicator, error) {
	var indicators []entity.ChartIndicator
	if err := chartDAO.ctx.DB.Order("id asc").Model(chart).Related(&indicators).Error; err != nil {
		return nil, err
	}
	return indicators, nil
}

func (chartDAO *ChartDAOImpl) GetStrategies(chart entity.ChartEntity) ([]entity.ChartStrategy, error) {
	var strategies []entity.ChartStrategy
	if err := chartDAO.ctx.DB.Order("id asc").Model(chart).Related(&strategies).Error; err != nil {
		return nil, err
	}
	return strategies, nil
}

func (chartDAO *ChartDAOImpl) GetTrades(user common.User) ([]entity.Trade, error) {
	var trades []entity.Trade
	daoUser := &entity.User{Id: user.GetId(), Username: user.GetUsername()}
	if err := chartDAO.ctx.DB.Order("id asc").Model(daoUser).Related(&trades).Error; err != nil {
		return nil, err
	}
	return trades, nil
}

func (chartDAO *ChartDAOImpl) GetLastTrade(chart entity.ChartEntity) (entity.TradeEntity, error) {
	var trades []entity.Trade
	if err := chartDAO.ctx.DB.Order("date desc").Limit(1).Model(chart).Related(&trades).Error; err != nil {
		chartDAO.ctx.Logger.Errorf("[ChartDAOImpl.GetLastTrade] Error: %s", err.Error())
	}
	tradeLen := len(trades)
	if tradeLen < 1 || tradeLen > 1 {
		return nil, nil
	}
	return &trades[0], nil
}

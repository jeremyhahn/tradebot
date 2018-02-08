package dao

import (
	"errors"

	"github.com/jeremyhahn/tradebot/common"
)

type ChartDAO interface {
	Create(chart ChartEntity) error
	Save(chart ChartEntity) error
	Update(chart ChartEntity) error
	Find(user *common.User) ([]Chart, error)
	Get(id uint) (ChartEntity, error)
	GetIndicators(chart ChartEntity) ([]ChartIndicator, error)
	GetStrategies(chart ChartEntity) ([]ChartStrategy, error)
	GetTrades(user *common.User) ([]Trade, error)
	GetLastTrade(chart ChartEntity) (*Trade, error)
}

type ChartDAOImpl struct {
	ctx   *common.Context
	Coins []Chart
	ChartDAO
}

func NewChartDAO(ctx *common.Context) ChartDAO {
	ctx.DB.AutoMigrate(&Chart{})
	ctx.DB.AutoMigrate(&ChartStrategy{})
	ctx.DB.AutoMigrate(&ChartIndicator{})
	ctx.DB.AutoMigrate(&Trade{})
	return &ChartDAOImpl{ctx: ctx}
}

func (chartDAO *ChartDAOImpl) Create(chart ChartEntity) error {
	return chartDAO.ctx.DB.Create(chart).Error
}

func (chartDAO *ChartDAOImpl) Save(chart ChartEntity) error {
	return chartDAO.ctx.DB.Save(chart).Error
}

func (chartDAO *ChartDAOImpl) Update(chart ChartEntity) error {
	return chartDAO.ctx.DB.Update(chart).Error
}

func (chartDAO *ChartDAOImpl) Get(id uint) (ChartEntity, error) {
	chart := Chart{}
	if err := chartDAO.ctx.DB.First(&chart, id).Error; err != nil {
		return nil, err
	}
	return &chart, nil
}

func (chartDAO *ChartDAOImpl) Find(user *common.User) ([]Chart, error) {
	var charts []Chart
	daoUser := &User{Id: user.Id}
	if err := chartDAO.ctx.DB.Model(daoUser).Related(&charts).Error; err != nil {
		return charts, err
	}
	for i, chart := range charts {
		var trades []Trade
		var indicators []ChartIndicator
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

func (chartDAO *ChartDAOImpl) GetIndicators(chart ChartEntity) ([]ChartIndicator, error) {
	var indicators []ChartIndicator
	if err := chartDAO.ctx.DB.Order("id asc").Model(chart).Related(&indicators).Error; err != nil {
		return nil, err
	}
	return indicators, nil
}

func (chartDAO *ChartDAOImpl) GetStrategies(chart ChartEntity) ([]ChartStrategy, error) {
	var strategies []ChartStrategy
	if err := chartDAO.ctx.DB.Order("id asc").Model(chart).Related(&strategies).Error; err != nil {
		return nil, err
	}
	return strategies, nil
}

func (chartDAO *ChartDAOImpl) GetTrades(user *common.User) ([]Trade, error) {
	var trades []Trade
	daoUser := &User{Id: user.Id, Username: user.Username}
	if err := chartDAO.ctx.DB.Order("id asc").Model(daoUser).Related(&trades).Error; err != nil {
		return nil, err
	}
	return trades, nil
}

func (chartDAO *ChartDAOImpl) GetLastTrade(chart ChartEntity) (*Trade, error) {
	var trades []Trade
	if err := chartDAO.ctx.DB.Order("date desc").Limit(1).Model(chart).Related(&trades).Error; err != nil {
		chartDAO.ctx.Logger.Errorf("[ChartDAOImpl.GetLastTrade] Error: %s", err.Error())
	}
	tradeLen := len(trades)
	if tradeLen < 1 || tradeLen > 1 {
		return nil, errors.New("Failed to retreive last trade")
	}
	return &trades[0], nil
}

type ChartEntity interface {
	GetId() uint
	GetUserId() uint
	GetBase() string
	GetQuote() string
	GetPeriod() int
	GetExchangeName() string
	IsAutoTrade() bool
	GetAutoTrade() uint
	SetIndicators(indicators []ChartIndicator)
	GetIndicators() []ChartIndicator
	AddIndicator(indicator *ChartIndicator)
	SetStrategies(strategies []ChartStrategy)
	GetStrategies() []ChartStrategy
	AddStrategy(strategy *ChartStrategy)
	SetTrades(trades []Trade)
	GetTrades() []Trade
	AddTrade(trade Trade)
}

type Chart struct {
	Id         uint   `gorm:"primary_key;AUTO_INCREMENT"`
	UserId     uint   `gorm:"foreign_key;unique_index:idx_chart"`
	Base       string `gorm:"unique_index:idx_chart"`
	Quote      string `gorm:"unique_index:idx_chart"`
	Exchange   string `gorm:"unique_index:idx_chart"`
	Period     int
	AutoTrade  uint
	Indicators []ChartIndicator `gorm:"ForeignKey:ChartId"`
	Strategies []ChartStrategy  `gorm:"ForeignKey:ChartId"`
	Trades     []Trade          `gorm:"ForeignKey:ChartId"`
	User       User
	ChartEntity
}

func (entity *Chart) GetId() uint {
	return entity.Id
}

func (entity *Chart) GetUserId() uint {
	return entity.UserId
}

func (entity *Chart) SetIndicators(indicators []ChartIndicator) {
	entity.Indicators = indicators
}

func (entity *Chart) GetIndicators() []ChartIndicator {
	return entity.Indicators
}

func (entity *Chart) AddIndicator(indicator *ChartIndicator) {
	entity.Indicators = append(entity.Indicators, *indicator)
}

func (entity *Chart) SetStrategies(strategies []ChartStrategy) {
	entity.Strategies = strategies
}

func (entity *Chart) GetStrategies() []ChartStrategy {
	return entity.Strategies
}

func (entity *Chart) AddStrategy(strategy *ChartStrategy) {
	entity.Strategies = append(entity.Strategies, *strategy)
}

func (entity *Chart) SetTrades(trades []Trade) {
	entity.Trades = trades
}

func (entity *Chart) GetTrades() []Trade {
	return entity.Trades
}

func (entity *Chart) AddTrade(trade Trade) {
	entity.Trades = append(entity.Trades, trade)
}

func (entity *Chart) GetBase() string {
	return entity.Base
}

func (entity *Chart) GetQuote() string {
	return entity.Quote
}

func (entity *Chart) GetPeriod() int {
	return entity.Period
}

func (entity *Chart) GetExchangeName() string {
	return entity.Exchange
}

func (entity *Chart) GetAutoTrade() uint {
	return entity.AutoTrade
}

func (entity *Chart) IsAutoTrade() bool {
	return entity.AutoTrade == 1
}

package dao

import (
	"github.com/jeremyhahn/tradebot/common"
)

type ChartDAO interface {
	Create(chart ChartEntity)
	Save(chart ChartEntity)
	Update(chart ChartEntity)
	Find(user *common.User) []Chart
	Get(id uint) ChartEntity
	GetIndicators(chart ChartEntity) map[string]Indicator
	GetTrades(user *common.User) []Trade
	GetLastTrade(chart ChartEntity) *Trade
	FindByCurrency(user *common.User, currencyPair *common.CurrencyPair) []Trade
}

type ChartDAOImpl struct {
	ctx   *common.Context
	Coins []Chart
	ChartDAO
}

type ChartEntity interface {
	GetId() uint
	GetBase() string
	GetQuote() string
	GetPeriod() int
	GetExchangeName() string
	IsAutoTrade() bool
	GetIndicators() []Indicator
	SetIndicators(indicators []Indicator)
	AddIndicator(indicator *Indicator)
	GetTrades() []Trade
	SetTrades(trades []Trade)
	AddTrade(trade *Trade)
}

type Chart struct {
	ID         uint   `gorm:"primary_key;AUTO_INCREMENT"`
	UserID     uint   `gorm:"foreign_key;unique_index:idx_chart"`
	Base       string `gorm:"unique_index:idx_chart"`
	Quote      string `gorm:"unique_index:idx_chart"`
	Exchange   string `gorm:"unique_index:idx_chart"`
	Period     int
	AutoTrade  uint
	Indicators []Indicator `gorm:"ForeignKey:ChartID"`
	Trades     []Trade     `gorm:"ForeignKey:ChartID"`
	ChartEntity
}

type Indicator struct {
	Id         uint   `gorm:"primary_key"`
	ChartID    uint   `gorm:"foreign_key;unique_index:idx_indicator"`
	Name       string `gorm:"unique_index:idx_indicator"`
	Parameters string `gorm:"not null"`
}

func NewChartDAO(ctx *common.Context) ChartDAO {
	ctx.DB.AutoMigrate(&Chart{})
	ctx.DB.AutoMigrate(&Indicator{})
	ctx.DB.AutoMigrate(&Trade{})
	return &ChartDAOImpl{ctx: ctx}
}

func (dao *ChartDAOImpl) Create(chart ChartEntity) {
	if err := dao.ctx.DB.Create(chart).Error; err != nil {
		dao.ctx.Logger.Errorf("[ChartDAOImpl.Create] Error:%s", err.Error())
	}
}

func (dao *ChartDAOImpl) Save(chart ChartEntity) {
	if err := dao.ctx.DB.Save(chart).Error; err != nil {
		dao.ctx.Logger.Errorf("[ChartDAOImpl.Save] Error:%s", err.Error())
	}
}

func (dao *ChartDAOImpl) Update(chart ChartEntity) {
	if err := dao.ctx.DB.Update(chart).Error; err != nil {
		dao.ctx.Logger.Errorf("[ChartDAOImpl.Update] Error:%s", err.Error())
	}
}

func (dao *ChartDAOImpl) Find(user *common.User) []Chart {
	var charts []Chart
	daoUser := &User{Id: user.Id}
	if err := dao.ctx.DB.Model(daoUser).Related(&charts).Error; err != nil {
		dao.ctx.Logger.Errorf("[ChartDAOImpl.Find] Error: %s", err.Error())
	}
	for i, chart := range charts {
		var trades []Trade
		var indicators []Indicator
		if err := dao.ctx.DB.Model(&chart).Related(&trades).Error; err != nil {
			dao.ctx.Logger.Errorf("[ChartDAOImpl.Find] Error: %s", err.Error())
		}
		if err := dao.ctx.DB.Model(&chart).Related(&indicators).Error; err != nil {
			dao.ctx.Logger.Errorf("[ChartDAOImpl.Find] Error: %s", err.Error())
		}
		charts[i].Indicators = indicators
		charts[i].Trades = trades
	}
	return charts
}

func (dao *ChartDAOImpl) GetIndicators(chart ChartEntity) map[string]Indicator {
	var indicators []Indicator
	if err := dao.ctx.DB.Model(chart).Related(&indicators).Error; err != nil {
		dao.ctx.Logger.Errorf("[ChartDAOImpl.GetIndicators] Error: %s", err.Error())
	}
	imap := make(map[string]Indicator, len(indicators))
	for _, i := range indicators {
		imap[i.Name] = i
	}
	return imap
}

func (dao *ChartDAOImpl) GetTrades(user *common.User) []Trade {
	var trades []Trade
	daoUser := &User{Id: user.Id, Username: user.Username}
	if err := dao.ctx.DB.Model(daoUser).Related(&trades).Error; err != nil {
		dao.ctx.Logger.Errorf("[ChartDAOImpl.GetTrades] Error: %s", err.Error())
	}
	return trades
}

func (dao *ChartDAOImpl) GetLastTrade(chart ChartEntity) *Trade {
	var trades []Trade
	if err := dao.ctx.DB.Order("date desc").Limit(1).Model(chart).Related(&trades).Error; err != nil {
		dao.ctx.Logger.Errorf("[ChartDAOImpl.GetLastTrade] Error: %s", err.Error())
	}
	tradeLen := len(trades)
	if tradeLen < 1 || tradeLen > 1 {
		dao.ctx.Logger.Warningf("[ChartDAOImpl.GetLastTrade] Invalid number of trades returned: %d", tradeLen)
		return &Trade{}
	}
	return &trades[0]
}

/*
func (dao *ChartDAOImpl) FindByCurrency(user *common.User, currencyPair *common.CurrencyPair) []Trade {
	var trades []Trade
	chart := &Chart{
		Base:   currencyPair.Base,
		Quote:  currencyPair.Quote,
		UserID: user.Id}
	if err := dao.ctx.DB.Model(chart).Find(&trades).Error; err != nil {
		dao.ctx.Logger.Errorf("[ChartDAOImpl.FindByCurrency] Error: %s", err.Error())
	}
	return trades
}*/

func (atc *Chart) GetId() uint {
	return atc.ID
}

func (atc *Chart) GetIndicators() []Indicator {
	return atc.Indicators
}

func (atc *Chart) SetIndicators(indicators []Indicator) {
	atc.Indicators = indicators
}

func (atc *Chart) AddIndicator(indicator *Indicator) {
	atc.Indicators = append(atc.Indicators, *indicator)
}

func (atc *Chart) GetTrades() []Trade {
	return atc.Trades
}

func (atc *Chart) SetTrades(trades []Trade) {
	atc.Trades = trades
}

func (atc *Chart) AddTrade(trade *Trade) {
	atc.Trades = append(atc.Trades, *trade)
}

func (atc *Chart) GetBase() string {
	return atc.Base
}

func (atc *Chart) GetQuote() string {
	return atc.Quote
}

func (atc *Chart) GetPeriod() int {
	return atc.Period
}

func (atc *Chart) GetExchangeName() string {
	return atc.Exchange
}

func (atc *Chart) IsAutoTrade() bool {
	return atc.AutoTrade == 1
}

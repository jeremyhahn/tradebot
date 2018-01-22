package dao

import (
	"fmt"

	"github.com/jeremyhahn/tradebot/common"
)

type IAutoTradeDAO interface {
	Create(coin IAutoTradeCoin)
	Save(coin IAutoTradeCoin)
	Update(coin IAutoTradeCoin)
	Find(user *common.User) []IAutoTradeCoin
	GetTrades(user *common.User) []Trade
	GetLastTrade(coin IAutoTradeCoin) *Trade
	FindByCurrency(user *common.User, currencyPair *common.CurrencyPair) []Trade
}

type AutoTradeDAO struct {
	ctx   *common.Context
	Coins []AutoTradeCoin
	IAutoTradeDAO
}

type IAutoTradeCoin interface {
	GetTrades() []Trade
	SetTrades(trades []Trade)
	AddTrade(trade *Trade)
	GetBase() string
	GetQuote() string
	GetPeriod() int
	GetExchange() string
}

type AutoTradeCoin struct {
	ID       uint   `gorm:"primary_key;AUTO_INCREMENT"`
	UserID   uint   `gorm:"unique_index:idx_autotrade_coin;"`
	Base     string `gorm:"unique_index:idx_autotrade_coin;"`
	Quote    string `gorm:"unique_index:idx_autotrade_coin;"`
	Exchange string `gorm:"unique_index:idx_autotrade_coin;"`
	Period   int
	Trades   []Trade `gorm:"ForeignKey:AutoTradeID"`
}

func NewAutoTradeDAO(ctx *common.Context) *AutoTradeDAO {
	ctx.DB.AutoMigrate(&AutoTradeCoin{})
	ctx.DB.AutoMigrate(&Trade{})
	return &AutoTradeDAO{ctx: ctx}
}

func (dao *AutoTradeDAO) Create(coin IAutoTradeCoin) {
	if err := dao.ctx.DB.Create(coin).Error; err != nil {
		dao.ctx.Logger.Errorf("[AutoTradeDAO.Create] Error:%s", err.Error())
	}
}

func (dao *AutoTradeDAO) Save(coin IAutoTradeCoin) {
	if err := dao.ctx.DB.Save(coin).Error; err != nil {
		dao.ctx.Logger.Errorf("[AutoTradeDAO.Save] Error:%s", err.Error())
	}
}

func (dao *AutoTradeDAO) Update(coin IAutoTradeCoin) {
	if err := dao.ctx.DB.Update(coin).Error; err != nil {
		dao.ctx.Logger.Errorf("[AutoTradeDAO.Update] Error:%s", err.Error())
	}
}

func (dao *AutoTradeDAO) Find(user *common.User) []IAutoTradeCoin {
	var coins []AutoTradeCoin
	daoUser := &User{Id: user.Id, Username: user.Username}
	if err := dao.ctx.DB.Model(daoUser).Related(&coins).Error; err != nil {
		dao.ctx.Logger.Errorf("[AutoTradeDAO.Find] Error: %s", err.Error())
	}
	fmt.Printf("%+v\n", coins)
	var icoins []IAutoTradeCoin
	for _, coin := range coins {
		icoins = append(icoins, &coin)
	}
	return icoins
}

func (dao *AutoTradeDAO) GetTrades(user *common.User) []Trade {
	var trades []Trade
	autoTradeCoin := &AutoTradeCoin{UserID: user.Id}
	if err := dao.ctx.DB.Model(autoTradeCoin).Related(&trades).Error; err != nil {
		dao.ctx.Logger.Errorf("[AutoTradeDAO.GetTrades] Error: %s", err.Error())
	}
	return trades
}

func (dao *AutoTradeDAO) GetLastTrade(coin IAutoTradeCoin) *Trade {
	var trades []Trade
	if err := dao.ctx.DB.Order("date desc").Limit(1).Model(coin).Find(&trades).Error; err != nil {
		dao.ctx.Logger.Errorf("[AutoTradeDAO.GetLastTrade] Error: %s", err.Error())
	}
	tradeLen := len(trades)
	if tradeLen < 1 || tradeLen > 1 {
		dao.ctx.Logger.Errorf("[AutoTradeDAO.GetLastTrade] Invalid number of trades returned: %d", tradeLen)
		return &Trade{}
	}
	return &trades[0]
}

func (dao *AutoTradeDAO) FindByCurrency(user *common.User, currencyPair *common.CurrencyPair) []Trade {
	var trades []Trade
	autoTradeCoin := &AutoTradeCoin{
		Base:   currencyPair.Base,
		Quote:  currencyPair.Quote,
		UserID: user.Id}
	if err := dao.ctx.DB.Model(autoTradeCoin).Find(&trades).Error; err != nil {
		dao.ctx.Logger.Errorf("[AutoTradeDAO.FindByCurrency] Error: %s", err.Error())
	}
	return trades
}

func (atc *AutoTradeCoin) GetTrades() []Trade {
	return atc.Trades
}

func (atc *AutoTradeCoin) SetTrades(trades []Trade) {
	atc.Trades = trades
}

func (atc *AutoTradeCoin) AddTrade(trade *Trade) {
	atc.Trades = append(atc.Trades, *trade)
}

func (atc *AutoTradeCoin) GetBase() string {
	return atc.Base
}

func (atc *AutoTradeCoin) GetQuote() string {
	return atc.Quote
}

func (atc *AutoTradeCoin) GetPeriod() int {
	return atc.Period
}

func (atc *AutoTradeCoin) GetExchange() string {
	return atc.Exchange
}

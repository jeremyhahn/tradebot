package dao

import "github.com/jeremyhahn/tradebot/common"

type IAutoTrade interface {
	Get(symbol string) string
}

type AutoTradeDAO struct {
	ctx   *common.Context
	Coins []AutoTradeCoin
	IAutoTrade
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

func (dao *AutoTradeDAO) Create(coin *AutoTradeCoin) {
	if err := dao.ctx.DB.Create(coin).Error; err != nil {
		dao.ctx.Logger.Errorf("[AutoTradeDAO.Create] Error:%s", err.Error())
	}
}

func (dao *AutoTradeDAO) Save(coin *AutoTradeCoin) {
	if err := dao.ctx.DB.Save(coin).Error; err != nil {
		dao.ctx.Logger.Errorf("[AutoTradeDAO.Save] Error:%s", err.Error())
	}
}

func (dao *AutoTradeDAO) Update(coin *AutoTradeCoin) {
	if err := dao.ctx.DB.Update(coin).Error; err != nil {
		dao.ctx.Logger.Errorf("[AutoTradeDAO.Update] Error:%s", err.Error())
	}
}

func (dao *AutoTradeDAO) Find(user *common.User) []AutoTradeCoin {
	var coins []AutoTradeCoin
	daoUser := &User{Id: user.Id, Username: user.Username}
	if err := dao.ctx.DB.Model(daoUser).Related(&coins).Error; err != nil {
		dao.ctx.Logger.Errorf("[AutoTradeDAO.Find] Error: %s", err.Error())
	}
	return coins
}

func (dao *AutoTradeDAO) GetTrades(user *common.User) []Trade {
	var trades []Trade
	autoTradeCoin := &AutoTradeCoin{UserID: user.Id}
	if err := dao.ctx.DB.Model(autoTradeCoin).Related(&trades).Error; err != nil {
		dao.ctx.Logger.Errorf("[AutoTradeDAO.GetTrades] Error: %s", err.Error())
	}
	return trades
}

func (dao *AutoTradeDAO) GetLastTrade(coin *AutoTradeCoin) *Trade {
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
		dao.ctx.Logger.Errorf("[TradeDAO.FindByCurrency] Error: %s", err.Error())
	}
	return trades
}

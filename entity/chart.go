package entity

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

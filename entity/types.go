package entity

import (
	"time"
)

type ProfitEntity interface {
	GetId() uint
	GetUserId() uint
	GetTradeId() uint
	GetQuantity() string
	GetBought() string
	GetSold() string
	GetFee() string
	GetTax() string
	GetTotal() string
}

type PluginEntity interface {
	GetName() string
	GetFilename() string
	GetVersion() string
	GetType() string
}

type TransactionEntity interface {
	GetId() string
	GetUserId() uint
	GetDate() time.Time
	GetMarketPair() string
	GetCurrencyPair() string
	GetType() string
	GetCategory() string
	SetCategory(category string)
	GetNetwork() string
	GetNetworkDisplayName() string
	GetQuantity() string
	GetQuantityCurrency() string
	GetFiatQuantity() string
	GetFiatQuantityCurrency() string
	GetPrice() string
	GetPriceCurrency() string
	GetFiatPrice() string
	GetFiatPriceCurrency() string
	GetQuoteFiatPrice() string
	GetQuoteFiatPriceCurrency() string
	GetFee() string
	GetFeeCurrency() string
	GetTotal() string
	GetTotalCurrency() string
	GetFiatFee() string
	GetFiatFeeCurrency() string
	GetFiatTotal() string
	GetFiatTotalCurrency() string
	IsDeleted() bool
	SetDeleted(value int)
}

type UserEntity interface {
	GetId() uint
	GetUsername() string
	GetLocalCurrency() string
	GetFiatExchange() string
	GetEtherbase() string
	GetKeystore() string
	GetWallets() []UserWallet
	GetExchanges() []UserCryptoExchange
}

type UserWalletEntity interface {
	GetUserId() uint
	GetCurrency() string
	GetAddress() string
	IsNative() bool
}

type UserTokenEntity interface {
	GetUserId() uint
	GetSymbol() string
	GetContractAddress() string
	GetWalletAddress() string
}

type UserExchangeEntity interface {
	GetUserID() uint
	GetName() string
	GetKey() string
	GetSecret() string
	GetExtra() string
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

type ChartIndicatorEntity interface {
	GetId() uint
	GetChartId() uint
	GetName() string
	GetParameters() string
}

type ChartStrategyEntity interface {
	GetId() uint
	GetChartId() uint
	GetName() string
	GetParameters() string
}

type TradeEntity interface {
	GetId() uint
	GetChartId() uint
	GetUserId() uint
	GetBase() string
	GetQuote() string
	GetExchangeName() string
	GetDate() time.Time
	GetType() string
	GetPrice() string
	GetAmount() string
	GetChartData() string
}

type PriceHistoryEntity interface {
	GetTime() int64
	GetOpen() float64
	GetHigh() float64
	GetLow() float64
	GetClose() float64
	GetVolume() float64
	GetMarketCap() int64
}

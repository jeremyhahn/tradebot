package entity

type MarketCap struct {
	Id               string `gorm:"index"`
	Name             string `gorm:"index"`
	Symbol           string `gorm:"primary_key"`
	Rank             string
	PriceUSD         string
	PriceBTC         string
	VolumeUSD24h     string
	MarketCapUSD     string
	AvailableSupply  string
	TotalSupply      string
	MaxSupply        string
	PercentChange1h  string
	PercentChange24h string
	PercentChange7d  string
	LastUpdated      string
}

type GlobalMarketCap struct {
	TotalMarketCapUSD float64
	Total24HVolumeUSD float64
	BitcoinDominance  float64
	ActiveCurrencies  float64
	ActiveMarkets     float64
	LastUpdated       int64
}

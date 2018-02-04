package dao

import (
	"github.com/jeremyhahn/tradebot/common"
	"github.com/jinzhu/gorm"
	logging "github.com/op/go-logging"
)

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

type MarketCapDAO struct {
	db     *gorm.DB
	logger *logging.Logger
}

func NewMarketCapDAO(ctx *common.Context) *MarketCapDAO {
	ctx.DB.AutoMigrate(&MarketCap{})
	ctx.DB.AutoMigrate(&GlobalMarketCap{})
	return &MarketCapDAO{
		db:     ctx.DB,
		logger: ctx.Logger}
}

func (dao *MarketCapDAO) Get(symbol string) (*common.MarketCap, error) {
	var market MarketCap
	if err := dao.db.First(&symbol).Error; err != nil {
		return nil, err
	}
	return &common.MarketCap{
		Id:               market.Id,
		Name:             market.Name,
		Symbol:           market.Symbol,
		Rank:             market.Rank,
		PriceUSD:         market.PriceUSD,
		PriceBTC:         market.PriceBTC,
		VolumeUSD24h:     market.VolumeUSD24h,
		MarketCapUSD:     market.MarketCapUSD,
		AvailableSupply:  market.AvailableSupply,
		TotalSupply:      market.TotalSupply,
		MaxSupply:        market.MaxSupply,
		PercentChange1h:  market.PercentChange1h,
		PercentChange24h: market.PercentChange24h,
		PercentChange7d:  market.PercentChange7d,
		LastUpdated:      market.LastUpdated}, nil
}

func (dao *MarketCapDAO) Delete(symbol string) error {
	return dao.db.Delete(&MarketCap{Symbol: symbol}).Error
}

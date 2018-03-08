package dao

import (
	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/entity"
	"github.com/jinzhu/gorm"
	logging "github.com/op/go-logging"
)

type MarketCapDAO struct {
	db     *gorm.DB
	logger *logging.Logger
}

func NewMarketCapDAO(ctx common.Context) *MarketCapDAO {
	return &MarketCapDAO{
		db:     ctx.GetCoreDB(),
		logger: ctx.GetLogger()}
}

func (dao *MarketCapDAO) Get(symbol string) (*common.MarketCap, error) {
	var market entity.MarketCap
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
	return dao.db.Delete(&entity.MarketCap{Symbol: symbol}).Error
}

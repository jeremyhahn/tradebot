package dao

import (
	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/entity"
	"github.com/jinzhu/gorm"
)

type Migrator interface {
	MigrateCoreDB()
	MigratePriceDB()
}

type MigratorImpl struct {
	db *gorm.DB
	Migrator
}

func NewMigrator(ctx common.Context) Migrator {
	return &MigratorImpl{db: ctx.GetCoreDB()}
}

func (migrator *MigratorImpl) MigrateCoreDB() {
	migrator.db.AutoMigrate(&entity.ChartIndicator{})
	migrator.db.AutoMigrate(&entity.ChartStrategy{})
	migrator.db.AutoMigrate(&entity.Chart{})
	migrator.db.AutoMigrate(&entity.Trade{})
	migrator.db.AutoMigrate(&entity.UserCryptoExchange{})
	migrator.db.AutoMigrate(&entity.Plugin{})
	migrator.db.AutoMigrate(&entity.MarketCap{})
	migrator.db.AutoMigrate(&entity.GlobalMarketCap{})
	migrator.db.AutoMigrate(&entity.Order{})
	migrator.db.AutoMigrate(&entity.Trade{})
	migrator.db.AutoMigrate(&entity.User{})
	migrator.db.AutoMigrate(&entity.UserWallet{})
	migrator.db.AutoMigrate(&entity.UserToken{})
	migrator.db.AutoMigrate(&entity.UserCryptoExchange{})
	migrator.db.AutoMigrate(&entity.PriceHistory{})
}

func (migrator *MigratorImpl) MigratePriceDB() {
	migrator.db.AutoMigrate(&entity.PriceHistory{})
}

package common

import (
	"fmt"
	"os"

	"github.com/jeremyhahn/tradebot/entity"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type DatabaseManager interface {
	ConnectCoreDB() *gorm.DB
	MigrateCoreDB()
	DropCoreDB()
	ConnectPriceDB() *gorm.DB
	MigratePriceDB()
	DropPriceDB()
	Close(*gorm.DB)
}

type DatabaseImpl struct {
	directory string
	prefix    string
	debugMode bool
	DatabaseManager
}

func NewDatabase() DatabaseManager {
	return CreateDatabase("", "./db", false)
}

func CreateDatabase(directory, prefix string, debugMode bool) DatabaseManager {
	return &DatabaseImpl{
		directory: directory,
		prefix:    prefix,
		debugMode: debugMode}
}

func (database *DatabaseImpl) MigrateCoreDB() {
	coreDB := database.ConnectCoreDB()
	if coreDB == nil {
		panic("Database: CoreDB database pointer is nil")
	}
	coreDB.AutoMigrate(&entity.ChartIndicator{})
	coreDB.AutoMigrate(&entity.ChartStrategy{})
	coreDB.AutoMigrate(&entity.Chart{})
	coreDB.AutoMigrate(&entity.Trade{})
	coreDB.AutoMigrate(&entity.UserCryptoExchange{})
	coreDB.AutoMigrate(&entity.Plugin{})
	coreDB.AutoMigrate(&entity.Profit{})
	coreDB.AutoMigrate(&entity.MarketCap{})
	coreDB.AutoMigrate(&entity.GlobalMarketCap{})
	coreDB.AutoMigrate(&entity.Order{})
	coreDB.AutoMigrate(&entity.Trade{})
	coreDB.AutoMigrate(&entity.User{})
	coreDB.AutoMigrate(&entity.UserWallet{})
	coreDB.AutoMigrate(&entity.UserToken{})
	coreDB.AutoMigrate(&entity.UserCryptoExchange{})
	coreDB.AutoMigrate(&entity.PriceHistory{})
}

func (database *DatabaseImpl) MigratePriceDB() {
	database.ConnectPriceDB().AutoMigrate(&entity.PriceHistory{})
}

func (database *DatabaseImpl) ConnectCoreDB() *gorm.DB {
	return database.newSQLite(fmt.Sprintf("%s/%s%s.db", database.directory, database.prefix, APPNAME))
}

func (database *DatabaseImpl) DropCoreDB() {
	os.Remove(fmt.Sprintf("%s/%s%s.db", database.directory, database.prefix, APPNAME))
}

func (database *DatabaseImpl) ConnectPriceDB() *gorm.DB {
	return database.newSQLite(fmt.Sprintf("%s/%sprices.db", database.directory, database.prefix))
}

func (database *DatabaseImpl) DropPriceDB() {
	os.Remove(fmt.Sprintf("%s/%sprices.db", database.directory, database.prefix))
}

func (database *DatabaseImpl) Close(db *gorm.DB) {
	db.Close()
}

func (database *DatabaseImpl) newSQLite(dbname string) *gorm.DB {
	db, err := gorm.Open("sqlite3", dbname)
	db.LogMode(database.debugMode)
	if err != nil {
		panic(err)
	}
	return db
}

/*
func NewMySQL() *gorm.DB {
	db, err := gorm.Open("mysql", "user:pass@tcp(ip:3306)/mydb?charset=utf8&parseTime=True")
	db.LogMode(true)
	if err != nil {
		panic(err)
	}
	return db
}
*/

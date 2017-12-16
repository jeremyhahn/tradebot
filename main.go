package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/op/go-logging"
)

const (
	APPNAME = "tradebot"
)

func main() {

	backend, _ := logging.NewSyslogBackend(APPNAME)
	logging.SetBackend(backend)
	logger := logging.MustGetLogger(APPNAME)

	sqlite := InitSQLite()
	defer sqlite.Close()

	mysql := InitMySQL()
	defer mysql.Close()

	config := NewConfiguration(sqlite, logger)
	coinbase := NewCoinbase(config, logger, "ETH-USD")
	trader := NewTrader(mysql, coinbase, logger)

	traders := make([]*Trader, 0)
	traders = append(traders, trader)

	go NewWebsocketServer(8080, traders, logger)

	trader.MakeMeRich()
}

func InitSQLite() *gorm.DB {
	db, err := gorm.Open("sqlite3", "./db/tradebot.db")
	db.LogMode(true)
	if err != nil {
		panic(err)
	}
	return db
}

func InitMySQL() *gorm.DB {
	db, err := gorm.Open("mysql", "gdaxlogger:PriceTracker007@tcp(192.168.0.10:3306)/gdaxlogger?charset=utf8&parseTime=True")
	db.LogMode(true)
	if err != nil {
		panic(err)
	}
	return db
}

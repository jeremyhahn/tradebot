package main

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/op/go-logging"
)

func main() {

	backend, _ := logging.NewSyslogBackend("tradebot")
	//format := logging.MustStringFormatter(`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,)
	//formatter := logging.NewBackendFormatter(backend, format)
	//logging.SetBackend(backend, formatter)
	logging.SetBackend(backend)
	logger := logging.MustGetLogger("coinbot")

	sqlite := InitSQLite()
	defer sqlite.Close()

	mysql := InitMySQL()
	defer mysql.Close()

	fmt.Println(time.Now().UTC())
	fmt.Println(time.Now().UTC().Hour())

	config := NewConfiguration(sqlite, logger)
	coinbase := NewCoinbase(config, logger)
	trader := NewTrader(mysql, coinbase, logger)
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

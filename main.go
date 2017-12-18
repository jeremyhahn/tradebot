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
	btc := NewCoinbase(config, logger, "BTC-USD")
	eth := NewCoinbase(config, logger, "ETH-USD")
	ltc := NewCoinbase(config, logger, "LTC-USD")

	btcChart := NewChart(mysql, btc, logger)
	ethChart := NewChart(mysql, eth, logger)
	ltcChart := NewChart(mysql, ltc, logger)

	charts := make([]*Chart, 0)
	charts = append(charts, btcChart)
	charts = append(charts, ethChart)
	charts = append(charts, ltcChart)

	ws := NewWebsocketServer(8080, charts, logger)
	go ws.Start()

	go btcChart.Stream(ws)
	go ethChart.Stream(ws)
	ltcChart.Stream(ws)
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

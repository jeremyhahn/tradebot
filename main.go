package main

import (
	"fmt"
	"time"

	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/exchange"
	"github.com/jeremyhahn/tradebot/service"
	"github.com/jeremyhahn/tradebot/websocket"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/op/go-logging"
)

func main() {

	backend, _ := logging.NewSyslogBackend(common.APPNAME)
	logging.SetBackend(backend)
	logger := logging.MustGetLogger(common.APPNAME)

	sqlite := InitSQLite()
	defer sqlite.Close()

	//mysql := InitMySQL()
	//defer mysql.Close()

	userDAO := dao.NewUserDAO(sqlite, logger)
	user := userDAO.Get(1)
	/*if user.Username == "" {
		userDAO.Create(&dao.User{
			Username: "test"})
	}*/

	//config := NewConfiguration(sqlite, logger)
	//period := 900 // seconds; 15 minutes

	gdaxCurrencyPair := &common.CurrencyPair{Base: "BTC", Quote: "USD", LocalCurrency: "USD"}
	bittrexCurrencyPair := &common.CurrencyPair{Base: "USDT", Quote: "BTC", LocalCurrency: "USDT"}
	binanceCurrencyPair := &common.CurrencyPair{Base: "BTC", Quote: "USDT", LocalCurrency: "USDT"}

	exchangeDAO := exchange.NewExchangeDAO(sqlite, logger)

	gdax := exchange.NewGDAX(exchangeDAO.Get("gdax"), logger, gdaxCurrencyPair)             // BTC-USD
	bittrex := exchange.NewBittrex(exchangeDAO.Get("bittrex"), logger, bittrexCurrencyPair) // USDT-BTC
	binance := exchange.NewBinance(exchangeDAO.Get("binance"), logger, binanceCurrencyPair) // BTCUSDT

	//btcGDAXChart := service.NewChart(sqlite, gdax, logger, period)
	//btcBittrexChart := service.NewChart(sqlite, bittrex, logger, period)
	//btcBinanceChart := service.NewChart(sqlite, binance, logger, period)

	//go btcGDAXChart.Stream()
	//go btcBittrexChart.Stream()
	//go btcBinanceChart.Stream()

	ws := websocket.NewWebsocketServer(8080, logger)
	go ws.Start()
	go ws.Run()

	exchangeMap := make(map[string]common.Exchange)
	exchangeMap["gdax"] = gdax
	exchangeMap["bittrex"] = bittrex
	exchangeMap["binance"] = binance

	for {
		//ws.PortfolioChan <- service.NewPortfolio(sqlite, logger, exchangeMap)
		//ws.ChartChan <- btcBinanceChart.GetChartData()

		portfolio := service.NewPortfolio(sqlite, logger, user, exchangeMap)
		ws.PortfolioChan <- portfolio

		fmt.Println("[Main] tick...")
		time.Sleep(1 * time.Second)
	}

}

func InitSQLite() *gorm.DB {
	db, err := gorm.Open("sqlite3", "./db/tradebot.db")
	db.LogMode(true)
	if err != nil {
		panic(err)
	}
	return db
}

/*
func InitMySQL() *gorm.DB {
	db, err := gorm.Open("mysql", "user:pass@tcp(ip:3306)/mydb?charset=utf8&parseTime=True")
	db.LogMode(true)
	if err != nil {
		panic(err)
	}
	return db
}
*/

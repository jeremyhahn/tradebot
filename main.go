package main

import (
	"github.com/jeremyhahn/tradebot/common"
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

	ctx := &common.Context{
		DB:     sqlite,
		Logger: logger}

	//userDAO := dao.NewUserDAO(ctx)
	//ctx.User = userDAO.GetById(1)
	/*if user.Username == "" {
		userDAO.Create(&dao.User{
			Username: "test"})
	}*/

	//config := NewConfiguration(sqlite, logger)
	//period := 900 // seconds; 15 minutes

	//marketcap := service.NewMarketCapService(ctx, dao.NewMarketCapDAO(ctx))
	//marketcap.GetMarkets()
	//marketcap.GetGlobalMarket("BTC")

	//wallets := userDAO.GetWallets(ctx.User)

	///gdaxCurrencyPair := &common.CurrencyPair{Base: "BTC", Quote: "USD", LocalCurrency: "USD"}
	///bittrexCurrencyPair := &common.CurrencyPair{Base: "USDT", Quote: "BTC", LocalCurrency: "USDT"}
	///binanceCurrencyPair := &common.CurrencyPair{Base: "BTC", Quote: "USDT", LocalCurrency: "USDT"}

	///exchangeDAO := exchange.NewExchangeDAO(ctx)

	///gdax := exchange.NewGDAX(exchangeDAO.Get("gdax"), logger, gdaxCurrencyPair)             // BTC-USD
	///bittrex := exchange.NewBittrex(exchangeDAO.Get("bittrex"), logger, bittrexCurrencyPair) // USDT-BTC
	///binance := exchange.NewBinance(exchangeDAO.Get("binance"), logger, binanceCurrencyPair) // BTCUSDT

	//btcGDAXChart := service.NewChart(ctx, gdax, period)
	//btcBittrexChart := service.NewChart(ctx, bittrex, period)
	//btcBinanceChart := service.NewChart(ctx, binance, period)

	//go btcGDAXChart.Stream()
	//go btcBittrexChart.Stream()
	//go btcBinanceChart.Stream()

	portfolioHub := websocket.NewPortfolioHub()
	go portfolioHub.Run()

	ws := websocket.NewWebsocketServer(ctx, 8080)
	go ws.Start(portfolioHub)
	ws.Run()

	/*
		exchangeMap := make(map[string]common.Exchange)
		exchangeMap["gdax"] = gdax
		exchangeMap["bittrex"] = bittrex
		exchangeMap["binance"] = binance
	*/

	/*
		for {

			fmt.Println("[Main] tick...")
			time.Sleep(1 * time.Second)

			//ws.PortfolioChan <- service.NewPortfolio(ctx, exchangeMap, wallets)
			//ws.ChartChan <- btcBinanceChart.GetChartData()
		}
	*/
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

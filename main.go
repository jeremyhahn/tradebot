package main

import (
	"github.com/jeremyhahn/tradebot/common"
	"github.com/jeremyhahn/tradebot/dao"
	"github.com/jeremyhahn/tradebot/service"
	"github.com/jeremyhahn/tradebot/strategy"
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

	ctx := &common.Context{
		DB:     sqlite,
		Logger: logger}

	userDAO := dao.NewUserDAO(ctx)
	ctx.User = userDAO.GetById(1)

	marketcapService := service.NewMarketCapService(logger)

	ws := websocket.NewWebsocketServer(ctx, 8080, marketcapService)
	go ws.Start()
	go ws.Run()

	/*
		cp := &common.CurrencyPair{
			Base:          "BTC",
			Quote:         "USD",
			LocalCurrency: "USD"}
		//gdax := service.NewExchangeService(ctx, dao.NewExchangeDAO(ctx)).NewExchange(ctx.User, "gdax", cp)
		//bittrex := service.NewExchangeService(ctx, dao.NewExchangeDAO(ctx)).NewExchange(ctx.User, "bittrex", cp)
		binance := service.NewExchangeService(ctx, dao.NewExchangeDAO(ctx)).NewExchange(ctx.User, "binance", cp)
		//gdax.GetTradeHistory()
		fmt.Println()
		fmt.Println()
		//bittrex.GetTradeHistory()
		fmt.Println()
		fmt.Println()
		binance.GetTradeHistory()
		os.Exit(1)*/

	//tradeService := service.NewTradeService(ctx, marketcapService)
	//tradeService.Trade()

	var services []common.ChartService
	exchangeDAO := dao.NewExchangeDAO(ctx)
	autoTradeDAO := dao.NewAutoTradeDAO(ctx)
	signalDAO := dao.NewSignalLogDAO(ctx)
	profitDAO := dao.NewProfitDAO(ctx)
	for _, autoTradeCoin := range autoTradeDAO.Find(ctx.User) {
		ctx.Logger.Debugf("[Tradebot.Main] Loading AutoTrade currency pair: %s-%s\n", autoTradeCoin.GetBase(), autoTradeCoin.GetQuote())
		currencyPair := &common.CurrencyPair{
			Base:          autoTradeCoin.GetBase(),
			Quote:         autoTradeCoin.GetQuote(),
			LocalCurrency: ctx.User.LocalCurrency}
		exchangeService := service.NewExchangeService(ctx, exchangeDAO)
		exchange := exchangeService.NewExchange(ctx.User, autoTradeCoin.GetExchange(), currencyPair)
		strategy := strategy.NewDefaultTradingStrategy(ctx, autoTradeCoin, autoTradeDAO, signalDAO, profitDAO)
		chart := service.NewChartService(ctx, exchange, strategy, autoTradeCoin.GetPeriod())
		ctx.Logger.Debugf("[Tradebot.Main] Chart: %+v\n", chart)
		services = append(services, chart)
	}

	for _, chart := range services {
		chart.Stream()
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
